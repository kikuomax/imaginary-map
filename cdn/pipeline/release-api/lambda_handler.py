# -*- coding: utf-8 -*-

import boto3
import io
import json
import logging
import zipfile


logger = logging.getLogger(__name__)
logger.setLevel(level=logging.DEBUG)

pipeline = boto3.client('codepipeline')

s3 = boto3.client('s3')

gateway = boto3.client('apigateway')


class StreamingBodyCloser(object):
    """Makes a StreamingBody a ContextManager.
    """
    def __init__(self, underlying):
        self.underlying = underlying

    def __enter__(self):
        return self.underlying

    def __exit__(self, exc_type, exc_value, traceback):
        self.underlying.close()


def get_job_id(event):
    """Returns the job ID in a given event.

    :type event: dict
    :param event: event from CodePipeline.

    :rtype: str
    :return: job ID in ``event``.

    :raises KeyError: if ``event`` does not have a necessary property.
    """
    return event['CodePipeline.job']['id']


def get_data(event):
    """Returns the data property in a given event.

    :type event: dict
    :param event: event from CodePipeline.

    :rtype: dict
    :return: data property in ``event``.

    :raises KeyError: if ``event`` does not have a necessary property.
    """
    return event['CodePipeline.job']['data']


def get_user_parameters(data):
    """Returns the user parameters in given data.

    :type data: dict
    :param data: data in a CodePipeline event.

    :rtype: dict
    :return: user parameters in ``data``. a JSON object.

    :raises KeyError: if ``data`` does not have a necessary property.
    """
    return json.loads(data['actionConfiguration']['configuration']['UserParameters'])


def get_input_artifacts(data):
    """Returns input artifacts in given data.

    :type data: dict
    :param data: data property in a CodePipeline event.

    :rtype: list
    :return: input artifacts in ``data``.

    :raises KeyError: if ``data`` does not have a necessary property.
    """
    return data['inputArtifacts']


def load_zip_file_on_s3(bucket, object_key):
    """Loads a zip file on S3.

    :type bucket: str
    :param bucket: name of the bucket where a zip file is stored.

    :type object_key: str
    :param object_key: object name associated with a zip file.

    :rtype: zipfile.ZipFile
    :return: ZipFile on S3.
    """
    obj = s3.get_object(Bucket=bucket, Key=object_key)
    with StreamingBodyCloser(obj['Body']) as body_in:
        body = body_in.read()
    return zipfile.ZipFile(io.BytesIO(body), mode='r')


def load_stack_outputs(input_artifacts, stack_outputs_file_name):
    """Loads stack outputs.

    :type input_artifacts: list
    :param input_artifacts: input artifacts from CodePipeline.

    :type stack_outputs_file_name: str
    :param stack_outputs_file_name: name of the stack outputs file given as
    an input artifact.

    :raises AssertionError: if ``input_artifacts`` is invalid.

    :raises KeyError: if ``input_artifacts`` does not have a necessary property.

    :raises json.JSONDecodeError:
    if stack outputs file is not a valid JSON file.
    """
    assert len(input_artifacts) == 1
    input_location = input_artifacts[0]['location']
    assert input_location['type'] == 'S3'
    s3_location = input_location['s3Location']
    zip_file = load_zip_file_on_s3(
        bucket=s3_location['bucketName'],
        object_key=s3_location['objectKey'])
    with zip_file as zip_in:
        return json.loads(zip_in.read(stack_outputs_file_name).decode())


def get_api_id(stack_outputs, api_logical_id):
    """Obtains an API ID from given stack outputs.

    :type stack_outputs: str
    :param stack_outputs: stack outputs.

    :type api_logical_id: str
    :param api_logical_id: logical ID of the API.

    :rtype: str
    :return: API ID.
    """
    return stack_outputs[f'{api_logical_id}Id']


def method_setting_to_patch_operations(method_setting):
    """Converts a given method setting into equivalent patch operations.

    A patch operation can be specified to ``gateway.update_stage``.

    :type method_setting: dict
    :param method_setting: method setting to convert.

    :rtype: list
    :return: list of path operations.

    :raises KeyError: if ``method_setting`` does not have a necessary property.
    """
    resource_path = method_setting['resourcePath']
    http_method = method_setting['httpMethod']
    base_path = f'{resource_path}/{http_method}'
    return [
        {
            'op': 'replace',
            'path': f'{base_path}/logging/loglevel',
            'value': method_setting.get('loggingLevel', 'OFF')
        },
        {
            'op': 'replace',
            'path': f'{base_path}/throttling/burstLimit',
            # allows a number in UserParameters
            'value': str(method_setting.get('throttlingBurstLimit', 5000))
        },
        {
            'op': 'replace',
            'path': f'{base_path}/throttling/rateLimit',
            # allows a number in UserParameters
            'value': str(method_setting.get('throttlingRateLimit', 10000))
        }
    ]
    return ops


def update_method_settings(api_id, stage_name, method_settings):
    """Updates method settings of a given stage.

    :type api_id: str
    :param api_id: ID of an API whose stage is to be updated.

    :type stage_name: str
    :param stage_name: name of a stage to be updated.

    :type method_settings: list
    :param method_settings: method settings to be applied to a stage.
    """
    logger.info(f'updating method settings: api={api_id}, stage={stage_name}')
    # TODO: consider stage's current settings
    patch_ops = [
        op for ms in method_settings
            for op in method_setting_to_patch_operations(ms)
    ]
    logger.debug(patch_ops)
    if len(patch_ops) == 0:
        return
    gateway.update_stage(
        restApiId=api_id,
        stageName=stage_name,
        patchOperations=patch_ops)


def create_stage(api_id, stage_name, configuration):
    """Creates a new stage.

    Call this function only when a stage does not exist.

    :type api_id: str
    :param api_id: ID of an API whose stage is to be created.

    :type stage_name: str
    :param stage_name: name of a stage to be created.

    :type configuration: dict
    :param configuration: configuration of a stage to be created.
    """
    logger.info(f'creating deployment: api={api_id}')
    description = f'deployment for {configuration.get("description", "stage")}'
    deployment = gateway.create_deployment(
        restApiId=api_id,
        description=description)
    logger.info(f'creating stage: api={api_id}, stage={stage_name}, deployment={deployment["id"]}')
    gateway.create_stage(
        restApiId=api_id,
        stageName=stage_name,
        description=configuration.get('description', 'stage'),
        deploymentId=deployment['id'],
        variables=configuration.get('variables', {}),
        tags=configuration.get('tags', {}))
    update_method_settings(
        api_id=api_id,
        stage_name=stage_name,
        method_settings=configuration.get('methodSettings', []))


def update_stage_deployment(api_id, stage_name, configuration):
    """Updates the deployment of a given stage.

    :type api_id: str
    :param api_id: ID of an API whose stage is to be updated.

    :type stage_name: str
    :param stage_name: name of a stage to be updated.
    """
    logger.info(f'creating deployment: api={api_id}')
    description = f'deployment for {configuration.get("description", "stage")}'
    deployment = gateway.create_deployment(
        restApiId=api_id,
        description=description)
    logger.info(f'updating stage deployment: api={api_id}, stage={stage_name}, deployment={deployment["id"]}')
    gateway.update_stage(
        restApiId=api_id,
        stageName=stage_name,
        patchOperations=[
            {
                'op': 'replace',
                'path': '/deploymentId',
                'value': deployment['id']
            }
        ])


def update_stage(api_id, stage_name, configuration):
    """Updates or creates an API stage.

    If a stage associated with ``stage_name`` does not exist, a new stage is
    created.

    :type api_id: str
    :param api_id: ID of an API whose stage is to be updated.

    :type stage_name: str
    :param stage_name: name of a stage to be updated.

    :type configuration: dict
    :param configuration: configuration of a stage to be updated.
    """
    try:
        logger.info(f'getting stage: api={api_id}, stage={stage_name}')
        stage = gateway.get_stage(restApiId=api_id, stageName=stage_name)
        logger.debug(stage)
        update_stage_deployment(
            api_id=api_id,
            stage_name=stage_name,
            configuration=configuration)
        update_method_settings(
            api_id=api_id,
            stage_name=stage_name,
            method_settings=configuration.get('methodSettings', []))
    except gateway.exceptions.NotFoundException:
        create_stage(
            api_id=api_id,
            stage_name=stage_name,
            configuration=configuration)


def put_success_result(job_id):
    """Tells CodePipeline that this job has succeeded.

    :type job_id: str
    :param job_id: ID of this job.
    """
    pipeline.put_job_success_result(jobId=job_id)


def put_failure_result(job_id, message, failure_type='JobFailed'):
    """Tells CodePipeline that this job has failed.

    :type job_id: str
    :param job_id: ID of this job.

    :type message: str
    :param message: message about the failure.

    :type failure_type: str
    :param failure_type: type of the failure.
    """
    pipeline.put_job_failure_result(
        jobId=job_id,
        failureDetails={
            'type': failure_type,
            'message': message
        })


def lambda_handler(event, context):
    try:
        logger.debug(event)
        job_id = get_job_id(event)
        logger.info(f'job ID: {job_id}')
        data = get_data(event)
        params = get_user_parameters(data)
        logger.debug(params)
        input_artifacts = get_input_artifacts(data)
        stack_outputs = load_stack_outputs(
            input_artifacts,
            params['stackOutputsFileName'])
        logger.debug(stack_outputs)
        api_id = get_api_id(stack_outputs, params['apiLogicalId'])
        logger.debug(f'API ID: {api_id}')
        update_stage(
            api_id=api_id,
            stage_name=params['stageName'],
            configuration=params['stageConfiguration'])
        put_success_result(job_id)
    except (
            KeyError,
            IndexError,
            AssertionError,
            json.JSONDecodeError,
            gateway.exceptions.NotFoundException) as err:
        logger.error(err)
        put_failure_result(
            job_id,
            message=str(err),
            failure_type='ConfigurationError')
    except BaseException as err:
        # is it OK to catch all kinds of errors?
        logger.error(err)
        put_failure_result(job_id, message=str(err))
