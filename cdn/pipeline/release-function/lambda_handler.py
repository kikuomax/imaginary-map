# -*- coding: utf-8 -*-

import boto3
import io
import json
import logging
import zipfile


logger = logging.getLogger(__name__)
logger.setLevel(logging.DEBUG)

session = boto3.session.Session()

lambda_api = session.client('lambda')

pipeline = session.client('codepipeline')


class StreamingBodyCloser(object):
    """Makes a StreamingBody a Context Manager.

    Makes sure that the ``close`` method of a StreamingBody is called.
    """
    def __init__(self, underlying):
        """Makes a given StreamingBody a Context Manager.

        :type underlying: botocore.StreamingBody
        :param underlying: a StreamingBody to be a Context Manager.
        """
        self.underlying = underlying

    def __enter__(self):
        return self.underlying

    def __exit__(self, exc_type, exc_value, traceback):
        self.underlying.close()


def get_job_id(event):
    """Results the job ID in a given event.

    :type event: dict
    :param event: event from CodePipeline.

    :rtype: str
    :return: job ID in ``event``.

    :raises KeyError: if ``event`` does not have a necessary property.
    """
    return event['CodePipeline.job']['id']


def get_account_id(event):
    """Returns the account ID in a given event.

    :type event: dict
    :param event: event from CodePipeline.

    :rtype: str
    :return: account ID in ``event``.

    :raises KeyError: if ``event`` does not have a necessary property.
    """
    return event['CodePipeline.job']['accountId']


def get_data(event):
    """Returns the data object in a given event.

    :type event: dict
    :param event: event from CodePipeline.

    :rtype: dict
    :return: data object in ``event``.

    :raises KeyError: if ``event`` does not have a necessary property.
    """
    return event['CodePipeline.job']['data']


def get_user_parameters_string(data):
    """Returns the UserParameters string in a given data object.

    :type data: dict
    :param data: data object in a CodePipeline event.

    :rtype: dict
    :return: UserParameters string in ``data``.

    :raises KeyError: if ``data`` does not have a necessary property.
    """
    return data['actionConfiguration']['configuration']['UserParameters']


def get_input_artifacts(data):
    """Returns the inputArtifacts array in a given data object.

    :type data: dict
    :param data: data object in a CodePipeline event.

    :rtype: dict
    :return: inputArtifacts object in ``data``.

    :raises KeyError: if ``data`` does not have ``inputArtifacts``.
    """
    return data['inputArtifacts']


def get_output_artifacts(data):
    """Returns the outputArtifacts array in a given data object.

    :type data: dict
    :param data: data object in a CodePipeline event.

    :rtype: dict
    :return: outputArtifacts object in ``data``.

    :raises KeyError: if ``data`` does not have ``outputArtifacts``.
    """
    return data['outputArtifacts']


def parse_user_parameters(params_str):
    """Parses a given UserParameters string.

    A UserParameters string must be a JSON object.

    :type params_str: str
    :param params_str: UserParameters string to be parsed.

    :rtype: dict
    :return: parsed user parameters.

    :raises JSONDecodeError: if ``params_str`` is not a valid JSON text.

    :raises ValueError: if ``params_str`` is not a JSON object.
    """
    params = json.loads(params_str)
    if not isinstance(params, dict):
        raise ValueError(f'UserParameters must be a JSON object not {type(params)}')
    return params


def load_zip_file_on_s3(bucket, object_key):
    """Loads a zip file stored as a given S3 object.

    :type bucket: str
    :param bucket: bucket name.

    :type object_key: str
    :param object_key: object key.

    :rtype: zipfile.ZipFile
    :return: zip file on S3.
    """
    s3 = session.client('s3')
    obj = s3.get_object(Bucket=bucket, Key=object_key)
    with StreamingBodyCloser(obj['Body']) as body_in:
        body_bytes = body_in.read()
    return zipfile.ZipFile(io.BytesIO(body_bytes), mode='r')


def load_stack_outputs(input_artifacts, stack_output_name):
    """Loads stack outputs from given input artifacts.

    :raises IndexError: if ``input_artifacts`` is empty.

    :raises KeyError: if ``input_artifacts`` does not have a necessary property.

    :raises JSONDecodeError: if the artifact is not a valid JSON file.
    """
    artifact_info = input_artifacts[0]
    logger.debug(artifact_info)
    location = artifact_info['location']
    assert location['type'] == 'S3'
    s3_location = location['s3Location']
    zip_file = load_zip_file_on_s3(
        bucket=s3_location['bucketName'],
        object_key=s3_location['objectKey'])
    with zip_file as zip_in:
        artifact_json = zip_in.read(stack_output_name).decode()
    return json.loads(artifact_json)


def get_function_name(stack_outputs, function_logical_id):
    """Obtains the function name from given stack outputs.

    :type stack_outputs: dict
    :param stack_outputs: CloudFormation stack outputs.

    :type function_logical_id: str
    :param function_logical_id: logical ID of the function resource.

    :rtype: str
    :return: function name.
    """
    return stack_outputs[f'{function_logical_id}Name']


def get_api_id(stack_outputs, api_logical_id):
    """Obtains the API ID from given stack outputs.


    :type stack_outputs: dict
    :param stack_outputs: CloudFormation stack outputs.

    :type api_logical_id: str
    :param api_logical_id: logical ID of the API resource.

    :rtype: str
    :return: API ID.
    """
    return stack_outputs[f'{api_logical_id}Id']


def get_function_alias(function_name, alias):
    """Obtains the information of a specified alias of a given function.

    :type function_name: str
    :param function_name: name of the function from which alias information is
    to be obtained.

    :type alias: str
    :param alias: name of an alias to be queried.

    :rtype: dict
    :return: alias information.

    :raises labmda_api.exceptions.ResourceNotFoundException:
    if the function does not have an alias specified by ``alias``.
    """
    logger.info(f'getting alias: function={function_name}, alias={alias}')
    return lambda_api.get_alias(
        FunctionName=function_name,
        Name=alias)


def update_function_alias(function_name, alias, version, description):
    """Updates a given function alias.

    A new function alias is created if no alias is defined.

    :type function_name: str
    :param function_name: name of the function whose alias is to be updated.

    :type alias: str
    :param alias: name of an alias to be updated.

    :type version: str
    :param version: version to be associated with ``alias``.

    :type description: str
    :param description: description of the function.
    """
    logger.info(f'updating alias: function={function_name}, alias={alias}, version={version}, description={description}')
    try:
        lambda_api.update_alias(
            FunctionName=function_name,
            Name=alias,
            FunctionVersion=version,
            Description=description)
    except lambda_api.exceptions.ResourceNotFoundException:
        # creates a new alias
        lambda_api.create_alias(
            FunctionName=function_name,
            Name=alias,
            FunctionVersion=version,
            Description=description)


def add_permission_to_api(
        function_name,
        alias,
        statement_id,
        api_id,
        account_id):
    """Adds a permission to an API.

    :type function_name: str
    :param function_name: function to grant an API access.

    :type alias: str
    :param alias: function alias to grant an API access.

    :type api_id: str
    :param api_id: API ID to be granted access.

    :type account_id: str
    :param account_id: account ID to be granted access.
    """
    function_arn = f'arn:aws:lambda:{session.region_name}:{account_id}:function:{function_name}:{alias}'
    action = 'lambda:InvokeFunction'
    principal = 'apigateway.amazonaws.com'
    source_arn = f'arn:aws:execute-api:{session.region_name}:{account_id}:{api_id}/*/GET/*'
    logger.debug(f'adding permission: function={function_name}, alias={alias}, source ARN={source_arn}')
    try:
        lambda_api.add_permission(
            FunctionName=function_arn,
            StatementId=statement_id,
            Action=action,
            Principal=principal,
            SourceArn=source_arn)
    except lambda_api.exceptions.ResourceConflictException:
        # verifies the existing policy
        logger.warn('permission already exists')
        policy = lambda_api.get_policy(FunctionName=function_arn)
        policy_document = json.loads(policy['Policy'])
        statement = [s for s in policy_document['Statement'] if s['Sid'] == statement_id][0]
        if statement['Effect'] != 'Allow':
            raise
        if statement['Principal']['Service'] != principal:
            raise
        if statement['Action'] != action:
            raise
        if statement['Resource'] != function_arn:
            raise
        if statement['Condition']['ArnLike']['AWS:SourceArn'] != source_arn:
            raise
        logger.info('existing permission is verified')


def put_success_result(job_id):
    """Tells CodePipeline this job has succeeded.

    :type job_id: str
    :param job_id: job ID.
    """
    pipeline.put_job_success_result(jobId=job_id)


def put_failure_result(job_id, message, failure_type='JobFailed'):
    """Tells CodePipeline this job has failed.

    :type job_id: str
    :param job_id: job ID.

    :type message: str
    :param message: message about the failure.

    :type failure_type: str
    :param failure_type: type of the failure.
    Must be one of values defined in https://docs.aws.amazon.com/codepipeline/latest/APIReference/API_FailureDetails.html
    """
    pipeline.put_job_failure_result(
        jobId=job_id,
        failureDetails={
            'type': failure_type,
            'message': message
        })


def lambda_handler(event, context):
    """Supposed to be called from CodePipeline.
    https://docs.aws.amazon.com/codepipeline/latest/userguide/action-reference-Lambda.html
    """
    logger.debug(event)
    try:
        job_id = get_job_id(event)
        logger.info(f'job ID: {job_id}')
        account_id = get_account_id(event)
        logger.info(f'account ID: {account_id}')
        data = get_data(event)
        params_str = get_user_parameters_string(data)
        params = parse_user_parameters(params_str)
        logger.debug(params)
        input_artifacts = get_input_artifacts(data)
        logger.debug(input_artifacts)
        output_artifacts = get_output_artifacts(data)
        logger.debug(output_artifacts)
        stack_outputs = load_stack_outputs(
            input_artifacts,
            stack_output_name=params['stackOutputFileName'])
        logger.debug(stack_outputs)
        function_name = get_function_name(
            stack_outputs,
            params['functionLogicalId'])
        api_id = get_api_id(
            stack_outputs,
            params['apiLogicalId'])
        src_alias_info = get_function_alias(function_name, params['srcAlias'])
        logger.debug(src_alias_info)
        function_version = src_alias_info['FunctionVersion']
        dest_alias = params['destAlias']
        update_function_alias(
            function_name,
            alias=dest_alias,
            version=function_version,
            description=params.get('description', ''))
        add_permission_to_api(
            function_name,
            alias=dest_alias,
            statement_id=params['permissionStatementId'],
            api_id=api_id,
            account_id=account_id)
        put_success_result(job_id)
    except (
            IndexError,
            KeyError,
            TypeError,
            ValueError,
            json.JSONDecodeError,
            lambda_api.exceptions.ResourceNotFoundException,
            lambda_api.exceptions.ResourceConflictException) as err:
        logger.error(err)
        put_failure_result(
            job_id,
            message=str(err),
            failure_type='ConfigurationError')
    except BaseException as err:
        # is it OK to catch all kinds of errors?
        logger.error(err)
        put_failure_result(
            job_id,
            message=str(err))
