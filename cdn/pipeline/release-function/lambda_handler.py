# -*- coding: utf-8 -*-

import boto3
import json
import logging


LOGGER = logging.getLogger(__name__)
LOGGER.setLevel(logging.DEBUG)


def get_data(event):
    """Returns the data object in a given event.

    :type event: dict
    :param event: event from CodePipeline.

    :rtype: dict
    :return: data object in ``event``.

    :raises KeyError: if ``event`` does not have a necessary key.
    """
    return event['CodePipeline.job']['data']


def get_user_parameters_string(data):
    """Returns the UserParameters string in a given data object.

    :type data: dict
    :param data: data object in a CodePipeline event.

    :rtype: dict
    :return: UserParameters string in ``data``.

    :raises KeyError: if ``data`` does not have a necessary key.
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


def lambda_handler(event, context):
    """Supposed to be called from CodePipeline.
    https://docs.aws.amazon.com/codepipeline/latest/userguide/action-reference-Lambda.html
    """
    LOGGER.debug(event)
    data = get_data(event)
    params_str = get_user_parameters_string(data)
    params = parse_user_parameters(params_str)
    LOGGER.debug(params)
    input_artifacts = get_input_artifacts(data)
    output_artifacts = get_output_artifacts(data)
