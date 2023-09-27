This PR uses [`.asf.yaml`](https://s.apache.org/asfyamltriage) to assign the GitHub [Triage role](https://docs.github.com/en/organizations/managing-access-to-your-organizations-repositories/repository-roles-for-an-organization#permissions-for-each-role) to non-committer contributors who fixed {ISSUE_THRESHOLD} or more Issues in the past {SINCE_DAYS_AGO} days:

{CONTRIB_LIST_LIST}

{CONGRATS} For the month of {MONTH}, {LIST_OF_CONTRIBUTORS} will have the ability to
* Apply labels to Issues and Pull Requests
* Assign a user to an Issue (note that non-committers can be assigned to an Issue after commenting on it)
* Add a user as a Reviewer of a Pull Request, which will send a request to them to review it
* Mark Issues and Pull Requests as a duplicate
<hr/>
{EXPIRE} If you want to be an Apache Traffic Control collaborator next month:

1. Read our [contribution guidelines]({REPO_URL}/blob/master/CONTRIBUTING.md)
2. Find an Issue to work on (recommended issues have the [good first issue]({REPO_URL}/issues?q=is:issue+is:open+label:"good+first+issue"+no:assignee) label) and ask to be assigned
3. Get coding! For questions on how to contribute, you can reach the ATC community on
    - The `#traffic-control` channel of the ASF Slack ([invite link](https://s.apache.org/tc-slack-request))
    - The ATC Dev [mailing list](https://trafficcontrol.apache.org/mailing_lists) ([archives](https://lists.apache.org/list?dev@trafficcontrol.apache.org:lte=5y:))
<hr/>

## Which Traffic Control components are affected by this PR?
- Other: [`.asf.yaml`]({REPO_URL}/blob/master/.asf.yaml)

## What is the best way to verify this PR?
Verify that the fixed Issues listed above are linked to [PRs from the past {SINCE_DAYS_AGO} days]({REPO_URL}/pulls?q=is:pr+linked:issue+merged:{SINCE_DAY}..{TODAY})

## PR submission checklist
- [ ] This PR has tests
- [ ] This PR has documentation
- [ ] This PR has a CHANGELOG.md entry
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
