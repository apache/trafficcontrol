import pytest
import json
import logging


"""
PyTest Fixture to store keys for cdns endpoint
"""
@pytest.fixture
def get_cdn_keys():
    # Response keys for cdns endpoint
    with open('post_data.json', 'r') as f:
         data = json.load(f)
    cdn_keys = list(data["cdns"].keys())
    return cdn_keys

"""
Test step to validate keys from cdns endpoint response in TO V4
"""
def test_get_cdn(to_login, get_cdn_keys, cdn_post_data):
    #validate CDN keys from cdns get response
    logging.info("Accessing Cdn endpoint through Traffic ops session")
    cdn_name = cdn_post_data[0]["cdns"]["name"]
    cdn_get_response = to_login.get_cdns(query_params={"name":str(cdn_name)})
    try:
       cdn_data = cdn_get_response[0] 
       cdn_keys = list(cdn_data[0].keys())
       logging.info("CDN Keys from cdns endpoint response {}".format(cdn_keys))
       #validate cdn values from post data in cdns get response
       post_data = [cdn_post_data[0]["cdns"]['name'], cdn_post_data[0]["cdns"]['domainName'], cdn_post_data[0]["cdns"]['dnssecEnabled']]
       get_data = [cdn_data[0]['name'], cdn_data[0]['domainName'], cdn_data[0]['dnssecEnabled']]
       assert cdn_keys.sort() == get_cdn_keys.sort()
       assert get_data == post_data
    except IndexError:
        logging.error("No CDN data from cdns get request")
        pytest.fail("Response is empty from get request , Failing test_get_cdn")


"""
Test step to validate POST method for cdns endpoint
"""
def test_post_cdn(cdn_post_data):

    post_data=cdn_post_data[0]['cdns']
    cdn_data = [post_data['name'], post_data['domainName'], post_data['dnssecEnabled']]
    print("Print pre post data for cdn endpoint {}".format(post_data))
    #Accessing cdns POST method
    try:
        cdn_response = cdn_post_data[1]
        cdn_response_data = [cdn_response['name'], cdn_response['domainName'], cdn_response['dnssecEnabled']]
        assert cdn_response_data == cdn_data
    except IndexError:
        logging.error("No CDN response data from cdns POST request")
        pytest.fail("Response is empty from POST request , Failing test_post_cdn")


def pytest_sessionfinish(session, exitstatus, cdn_post_data, to_login):
    try:
        cdn_response = cdn_post_data[1]
        id = cdn_response["id"]
        to_login.delete_cdn_by_id(cdn_id=id)
    except IndexError:
        logging.error("CDN doesn't created")

    
    

