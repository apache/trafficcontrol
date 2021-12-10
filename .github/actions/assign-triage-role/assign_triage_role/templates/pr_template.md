<!--
Thank you for contributing! Please be sure to read our contribution guidelines: https://github.com/apache/trafficcontrol/blob/master/CONTRIBUTING.md
If this closes or relates to an existing issue, please reference it using one of the following:

Closes: #ISSUE
Related: #ISSUE

If this PR fixes a security vulnerability, DO NOT submit! Instead, contact
the Apache Traffic Control Security Team at security@trafficcontrol.apache.org and follow the
guidelines at https://apache.org/security regarding vulnerability disclosure.
-->
This PR uses [`.asf.yaml`](https://s.apache.org/asfyamltriage) to assign the GitHub [Triage role](https://docs.github.com/en/organizations/managing-access-to-your-organizations-repositories/repository-roles-for-an-organization#permissions-for-each-role) to non-committer contributors who fixed {ISSUE_THRESHOLD} or more Issues in the past {SINCE_DAYS_AGO} days:

{CONTRIB_LIST_LIST}

{CONGRATS} For the month of {MONTH}, {LIST_OF_CONTRIBUTORS} will have the ability to
* Apply labels to Issues and Pull Requests
* Assign a user to an Issue (note that non-committers must first comment on an Issue before they can be assigned to it)
* Add a user as a Reviewer of a Pull Request, which will send a request to them to review it
* Mark Issues and Pull Requests as a duplicate
<hr>
{EXPIRE} If you want to be an Apache Traffic Control collaborator next month:

1. Read our [contribution guidelines](https://github.com/apache/trafficcontrol/blob/master/CONTRIBUTING.md)
2. Find an Issue to work on (recommended issues have the [good first issue](https://github.com/apache/trafficcontrol/issues?q=is:issue+is:open+label:"good+first+issue"+no:assignee) label) and ask to be assigned
3. Get coding! For questions on how to contribute, you can reach the ATC community on
    - The `#traffic-control` channel of the ASF Slack ([invite link](https://s.apache.org/tc-slack-request))
    - The ATC Dev [mailing list](https://trafficcontrol.apache.org/mailing_lists) ([archives](https://lists.apache.org/list?dev@trafficcontrol.apache.org:lte=5y:))
<!-- **^ Add meaningful description above** --><hr>

## Which Traffic Control components are affected by this PR?
<!-- Please delete all components from this list that are NOT affected by this PR.
Feel free to add the name of a tool or script that is affected but not on the list.
-->
- Other: [`.asf.yaml`](https://github.com/apache/trafficcontrol/blob/master/.asf.yaml)

## What is the best way to verify this PR?
<!-- Please include here ALL the steps necessary to test your PR.
If your PR has tests (and most should), provide the steps needed to run the tests.
If not, please provide step-by-step instructions to test the PR manually and explain why your PR does not need tests. -->
Verify that the fixed Issues listed above are linked to [PRs from the past {SINCE_DAYS_AGO} days](https://github.com/apache/trafficcontrol/pulls?q=is:pr+linked:issue+merged:{SINCE_DAY}..{TODAY})

## PR submission checklist
- [ ] This PR has tests <!-- If not, please delete this text and explain why this PR does not need tests. -->
- [ ] This PR has documentation <!-- If not, please delete this text and explain why this PR does not need documentation. -->
- [ ] This PR has a CHANGELOG.md entry <!-- A fix for a bug from an ATC release, an improvement, or a new feature should have a changelog entry. -->
- [x] This PR **DOES NOT FIX A SERIOUS SECURITY VULNERABILITY** (see [the Apache Software Foundation's security guidelines](https://apache.org/security) for details)

<!--
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
-->
