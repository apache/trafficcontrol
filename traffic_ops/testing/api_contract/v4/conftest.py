#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

"""This module is used to create a Traffic Ops session 
and to store prerequisite data for endpoints."""
import json
import logging
import sys
from random import randint
from urllib.parse import urlparse
import pytest
from trafficops.tosession import TOSession
from trafficops.restapi import OperationError


# Create and configure logger
logger = logging.getLogger()


def pytest_addoption(parser):
    """Passing in Traffic Ops arguments [Username, Password, Url and Hostname] from command line.
    :param parser: Parser to parse command line arguments
    :type parser: object
    """
    parser.addoption(
        "--to-user", action="store", help="User name for Traffic Ops Session."
    )
    parser.addoption(
        "--to-password", action="store", help="Password for Traffic Ops Session."
    )
    parser.addoption(
        "--to-url", action="store", help="Traffic Ops URL."
    )


@pytest.fixture(name="to_args")
def to_data(pytestconfig):
    """PyTest fixture to store Traffic ops arguments passed from command line.
    :param pytestconfig: Session-scoped fixture that returns the session's pytest.Config object
    :type pytestconfig: pytest.Config object
    :returns args: Return Traffic Ops arguments
    :rtype: dict
    """
    args = {}
    with open("to_data.json", encoding="utf-8", mode="r") as session_file:
        data = json.load(session_file)
    session_data = data["test"]
    args["api_version"] = urlparse(session_data.get("url")).path.strip('/').split('/')[1]
    args["port"] = session_data.get("port")

    to_user = pytestconfig.getoption("--to-user")
    to_password = pytestconfig.getoption("--to-password")
    to_url = pytestconfig.getoption("--to-url")

    if not all([to_user, to_password, to_url]):
        logger.info(
            "Traffic Ops session data were not passed from Command line Args.")
        args["user"] = session_data.get("user")
        args["password"] = session_data.get("password")
        args["url"] = session_data.get("url")
    else:
        args["user"] = to_user
        args["password"] = to_password
        args["url"] = to_url
        logger.info("Parsed Traffic ops session data from args %s", args)
    return args


@pytest.fixture(name="to_session")
def to_login(to_args):
    """PyTest Fixture to create a Traffic Ops session from Traffic Ops Arguments
    passed as command line arguments in to_args fixture in conftest.
    :param to_args: Fixture to get Traffic ops session arguments
    :type to_args: dict
    :returns to_session: Return Traffic ops session
    :rtype: TOSession object
    """
    # Create a Traffic Ops V4 session and login
    to_url = urlparse(to_args["url"])
    to_host = to_url.hostname
    try:
        to_session = TOSession(host_ip=to_host, host_port=to_args["port"],
                               api_version=to_args["api_version"], ssl=True, verify_cert=False)
        logger.info("Established Traffic Ops Session.")
    except OperationError:
        sys.exit(-1)

    # Login To TO_API
    to_session.login(to_args["user"], to_args["password"])
    logger.info("Successfully logged into Traffic Ops.")
    return to_session


@pytest.fixture()
def cdn_prereq(to_session, get_cdn_data):
    """PyTest Fixture to create POST data for cdns endpoint.
    :param to_session: Fixture to get Traffic ops session 
    :type to_session: TOsession
    :param get_cdn_data: Fixture to get cdn data from a prereq file
    :type get_cdn_data: dict
    :returns prerequisite_data: Returns sample Post data and actual api response
    :rtype: list
    """

    # Return new post data and post response from cdns POST request
    get_cdn_data["name"] = get_cdn_data["name"][:4]+str(randint(0, 1000))
    get_cdn_data["domainName"] = get_cdn_data["domainName"][:5] + \
        str(randint(0, 1000))
    logger.info("New cdn data to hit POST method %s", get_cdn_data)
    # Hitting cdns POST methed
    response = to_session.create_cdn(data=get_cdn_data)
    prerequisite_data = None
    try:
        cdn_response = response[0]
        prerequisite_data = [get_cdn_data, cdn_response]
    except IndexError:
        logger.error("No CDN response data from cdns POST request.")
    return prerequisite_data
