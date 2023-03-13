import pytest
import json
import logging


"""
PyTest Fixture to store keys for cdns endpoint
"""
@pytest.fixture
def get_cdn_keys():
    # Response keys for cdns endpoint
    with open('prerequisite_data.json', 'r') as f:
        data = json.load(f)
    cdn_keys = list(data["cdns"].keys())
    return cdn_keys


"""
Test step to validate keys from cdns endpoint response and POST method
"""
def test_get_cdn(to_login, get_cdn_keys, cdn_prereq):
    # validate CDN keys from cdns get response
    logging.info("Accessing Cdn endpoint through Traffic ops session")
    cdn_name = cdn_prereq[0]["cdns"]["name"]
    cdn_get_response = to_login.get_cdns(query_params={"name": str(cdn_name)})
    try:
        cdn_data = cdn_get_response[0]
        cdn_keys = list(cdn_data[0].keys())
        logging.info(
            "CDN Keys from cdns endpoint response {}".format(cdn_keys))
        # validate cdn values from prereq data in cdns get response
        prereq_data = [cdn_prereq[0]["cdns"]['name'], cdn_prereq[0]["cdns"]['domainName'], cdn_prereq[0]["cdns"]['dnssecEnabled']]
        get_data = [cdn_data[0]['name'], cdn_data[0]['domainName'], cdn_data[0]['dnssecEnabled']]
        assert cdn_keys.sort() == get_cdn_keys.sort()
        assert get_data == prereq_data
    except IndexError:
        logging.error("No CDN data from cdns get request")
        pytest.fail("Response is empty from get request , Failing test_get_cdn")


"""
Delete CDN after test execution to avoid redundancy"""
def pytest_sessionfinish(cdn_prereq, to_login):
    try:
        cdn_response = cdn_prereq[1]
        id = cdn_response["id"]
        to_login.delete_cdn_by_id(cdn_id=id)
    except IndexError:
        logging.error("CDN doesn't created")
