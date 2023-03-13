import pytest
from trafficops.tosession import TOSession
from urllib.parse import urlparse
import sys 
import json
from random import randint
import logging


""" 
Passing in Traffic Ops Arguments [Username, Password, Url and Hostname] from Command Line
"""
def pytest_addoption(parser):
    parser.addoption(
        '--to_user', action='store', help='User name for Traffic Ops Session'
    )
    parser.addoption(
        '--to_password', action='store',  help='Password for Traffic Ops Session'
    )
    parser.addoption(
        '--to_url', action='store',  help='Traffic Ops URL'
    )
    parser.addoption(
        '--hostname', action='store',  help='Traffic Ops hostname'
    )


""" 
PyTest fixture to store Traffic ops Arguments passed from command line
"""
@pytest.fixture
def to_data(pytestconfig):
    args = {}
    args['user'] = pytestconfig.getoption('--to_user')
    args['password'] = pytestconfig.getoption('--to_password')
    args['url'] = pytestconfig.getoption('--to_url')
    args['hostname'] = pytestconfig.getoption('--hostname')
    return args


"""
PyTest Fixture to create a Traffic Ops session from Traffic Ops Arguments
passed as command line arguments in to_data fixture in conftest
"""
@pytest.fixture()
def to_login(to_data):

    # Create a Traffic Ops V4 session and login
    print("Parsed TO args {}".format(to_data))
    if to_data["user"] == None:
        f = open('to_data.json')
        data = json.load(f)
        to_data = data["test"]
    to_url = urlparse(to_data["url"])
    to_host = to_url.hostname

    TO=TOSession(host_ip=to_host,host_port=443,api_version='4.0',ssl=True,verify_cert=False)
    print("established connection")

    # Login To TO_API
    TO.login(to_data["user"], to_data["password"])

    if not TO.logged_in:
        print("Failure Logging into Traffic Ops")
        sys.exit(-1)
    else:
        print("Successfully logged into Traffic Ops")
    return TO


"""
PyTest Fixture to create POST data for cdns endpoint
"""
@pytest.fixture()
def cdn_post_data(to_login):

    #Return new post data and post response from cdns POST request
    f = open('prerequisite_data.json')
    data = json.load(f)

    data["cdns"]["name"] =data["cdns"]["name"][:4]+str(randint(0,1000))
    data["cdns"]["domainName"] = data["cdns"]["domainName"][:5] + str(randint(0,1000))
    logging.info("New post data to hit POST method {}".format(data))
    with open('prerequisite_data.json', 'w') as f:
         json.dump(data, f)
    f.close()
    #Hitting cdns POST methed
    response = to_login.create_cdn(data=data["cdns"])
    try:
        cdn_response = response[0]
        prerequisite_data = [data, cdn_response]
        return prerequisite_data
    except IndexError:
        logging.error("No CDN response data from cdns POST request")
