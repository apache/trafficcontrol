#!/usr/bin/env python
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#


# Takes an XML file (created by JIRA issue export) and
#  1) Performs dry-run to list labels and milestones that need to be manually created
#  2) creates equivalent issues in github
#

import sys
import argparse
import xml.etree.ElementTree as ET
import json
import requests
import logging
import time

logging.basicConfig(level=logging.INFO)
repo_name = "apache/trafficcontrol"

#
# List of relevant JIRA field names
#
field_names = ['title', 'link', 'description', 'type', 'priority', 'status', 'assignee',
               'reporter', 'labels', 'fixVersion', 'component', 'resolution', 'version']

assignee_map = {
    'Jan van Doorn' :'knutsel',
    'Jeff Elsloo': 'elsloo',
    'Jifeng Yang': 'jifyang',
    'Dan Kirkwood' :'dangogh',
    'David Neuman' :'dneuman64',
    'Nir Sopher' :'nir-sopher',
    'Jeremy Mitchell' :'mitchell852',
    'Robert Butts' :'rob05c',
    'Zhilin Huang' :'zhilhuan',
    'Dewayne Richardson' :'dewrich',
    'Rawlin Peters' :'rawlinp',
    'Dylan Volz' :'DylanVolz',
    'Matt Mills' :'MattMills',
    'John Shen' :'weifensh',
    'Derek Gelinas' :'dg4prez',
    'Steve Malenfant': 'smalenfant',
    'Hank Beatty': 'hbeatty',
}

milestone_map = {'2.1.0': 1,
                 '2.2.0': 2,
                 '3.x': 3}

class GithubAPI(object):
    ''' Base Class for callers of the Github v3API
        
        Handles authentication, rate limiting and sending proper headers
    '''
    def __init__(self, username, api_key, repo, gh_server="api.github.com"):
        self.username = username
        self.api_key = api_key
        self.repo = repo
        self.gh_server = gh_server
        self._set_base_url()
        self.abuse_delay = 3.5 # Github API allows only 20 requests per 60 seconds

    def set_github_server(self, server_name):
        self.gh_server = server_name
        self._set_base_url()

    def _set_base_url(self):
        self.base_url = str.format("https://{gh_server}/repos/{repo}",
                                   gh_server=self.gh_server,
                                   repo=self.repo)

        
    def make_request(self, path, json_body, use_abuse_delay=False):
        url = self.base_url + path
        logging.debug("POST %s\n%s" % (url, json.dumps(json_body)))
        rsp = requests.post(url,
                            headers = {"Accept": "application/vnd.github.v3+json"},
                            auth=requests.auth.HTTPBasicAuth(self.username,
                                                             self.api_key),
                            json=json_body)
        for (name, val) in rsp.headers.iteritems():
            logging.debug("%s:%s" % (name, val))
        logging.debug(json.dumps(rsp.json()))
                          
        rsp.raise_for_status()
        if int(rsp.headers['X-RateLimit-Remaining']) == 0:
            wait_time = time.time() - int(rsp.headers['X-RateLimit-Reset'])
            logging.warn("Pausing %d seconds for rate-limit" % (wait_time))
            time.sleep(wait_time)

        if use_abuse_delay:
            time.sleep(self.abuse_delay)
            
        return rsp
        
    
class GithubIssues(GithubAPI):
    def __init__(self, username, api_key, repo):
        super(GithubIssues, self).__init__(username, api_key, repo)

        
    def create_issue(self, issue):
        rsp = self.make_request("/issues", issue, use_abuse_delay=True)
        return rsp
    

class JiraIssueImporter(object):
    def __init__(self, input_file):
        self.input_file = input_file

        
    def parse(self):
        self.tree = ET.parse(self.input_file)
        self.root = self.tree.getroot()

        
    def collect_issues(self):
        self.issue_list = list()
        for item in self.root[0].findall('item'):
            issue = dict()
            for field in field_names:
                if field == "labels":
                    jira_list = item.find('labels').findall('label')
                    issue['labels'] = [label.text for label in jira_list]

                else:
                    try:
                        issue[field] = item.find(field).text
                    except AttributeError:
                        issue[field] = None
                        logging.debug("Missing %s for item %s" % (field, item.find('title').text))

            self.issue_list.append(issue)

    def get_unique_values(self, field_name):
        ''' Iterates over issues and returns list of unqiue values in a field'''
        vals = {str(x[field_name]) for x in self.issue_list if x[field_name] is not None}
        return vals

    
    def translate_issue_to_github(self, jira_issue):
        ''' Translate JIRA format to Github format'''
        new_issue = {}
        new_issue['title'] = jira_issue['title']
        new_issue['body'] = "%s\n\nAuthor: %s\nJIRA Link: <a href=\"%s\">%s</a>" % \
                            (jira_issue['description'],
                             jira_issue['reporter'],
                             jira_issue['link'],
                             jira_issue['link'])
        if jira_issue['version'] is not None:
            new_issue['body'] += "\nFound Version: %s" % (jira_issue['version'])

        #if jira_issue['assignee'] != "Unassigned":
        #    new_issue['assignees'] = [assignee_map[jira_issue['assignee']]]

        if jira_issue['fixVersion'] is not None:
                new_issue['milestone'] = milestone_map[jira_issue['fixVersion']]
                
        new_issue['labels'] = jira_issue['labels']
        for extra_label in ['type', 'priority', 'component']:
            if jira_issue[extra_label] is not None:
                new_issue['labels'].append(jira_issue[extra_label])
        return new_issue

    
    def translate_all_issues(self):
        gh_issues = []
        for i in self.issue_list:
            gh_issues.append(self.translate_issue_to_github(i))
            
        return gh_issues

def list_all(importer):
    ''' Print list of unique values in JIRA fields'''
    print "Listing all components and milestones"
    print "Fields:\n ",
    print '\n  '.join(field_names)
    print

    for field in field_names:
        if field in ["description", "title", "link"]:
            continue
            
        print "%s:\n " %(field),
        print '\n  '.join(importer.get_unique_values(field))
        print    
    
def main():
    args = parse_command_line()
    importer = JiraIssueImporter(args.input_file)
    importer.parse()
    importer.collect_issues()

    if args.list_all:
        list_all(importer)

    else:
        gh_issues =  importer.translate_all_issues()
        gh_api = GithubIssues(args.username, args.api_key, repo_name)
        for (idx, issue) in enumerate(gh_issues):
            logging.info("Creating Issue (%d/%d): %s" % (idx, len(gh_issues), issue['title']))
            try:
                gh_api.create_issue(issue)
            except requests.exceptions.HTTPError as e:
                
                raise e

            
def parse_command_line():
  parser = argparse.ArgumentParser(description='Import JIRA issues from XML file into Github')
  parser.add_argument('input_file', help='path to input XML file')
  parser.add_argument('username', help='Github Username')
  parser.add_argument('api_key', help='Github username/Personal Access Token')  
  parser.add_argument('--list-all',
                      action='store_true',
                      help='List all issues and milestones in JIRA')


  args = parser.parse_args()
  print "Reading input from: %s" % (args.input_file)  
  return args    

    
if __name__ == '__main__':
  main()
