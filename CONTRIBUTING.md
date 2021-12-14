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

Contributing to Traffic Control
===============================
Thanks again for your time and interest in this project!

The following document is a set of guidelines to streamline the contribution
process for our contributors. Please feel free to suggest changes to this
document in a pull request!

Things to know before getting started
-------------------------------------
### Code of Conduct
Please try to keep discussions respectful and follow the
[Apache Software Foundation Code of Conduct](CODE_OF_CONDUCT.md) when
interacting with others.

How to Contribute
-----------------
We love pull requests! We simply don't have the time or resources to add every
feature, fix every bug and support every platform. If you have improvements
(enhancements or bug fixes) for Traffic Control, start by creating a
[Github issue](https://github.com/apache/trafficcontrol/issues) or a
[feature blueprint](blueprints/BLUEPRINT_TEMPLATE.md) and discussing it with us
first on the
[dev@trafficcontrol.apache.org](mailto:dev@trafficcontrol.apache.org) mailing
list. We might already be working on it, or there might be an existing way to do
it.

#### Design Decisions

##### Discussion Email
When we need to make changes to the project, we first discuss it on the
[dev@trafficcontrol.apache.org](mailto:dev@trafficcontrol.apache.org) mailing
list. Use your best judgement here. Small changes (i.e. bug fixes or small
improvements) may not warrant an email discussion. However, larger changes (i.e.
new features or changes to existing functionality) should be discussed on the
email list prior to development.

##### Feature Blueprints (optional)
Some contributors may choose to create a feature blueprint to accompany their
email to better articulate an idea and solicit feedback in the form of a PR. The
process involves the following:

1. Create a new PR that includes a markdown file that utilizes the
[BLUEPRINT_TEMPLATE.md](blueprints/BLUEPRINT_TEMPLATE.md)
template. For example, submit a PR for blueprints/my-feature.md.
1. Send an email to the
[dev@trafficcontrol.apache.org](mailto:dev@trafficcontrol.apache.org) mailing
list with a short description of your feature plus a link to the blueprint PR.
1. Wait for feedback to be applied to your blueprint PR. Because it's a PR,
line-specific feedback can be given.
1. Just like any PR, once all the concerns have been addressed, the blueprint is
merged into the blueprints directory (if accepted) or closed (if rejected). Note: Blueprint authors have the option of invoking [lazy consensus](https://community.apache.org/committers/lazyConsensus.html) to facilitate the merge of the blueprint if community feedback is not being provided or feedback has stalled. To invoke lazy consensus, please make the community aware via an email to the [dev@trafficcontrol.apache.org](mailto:dev@trafficcontrol.apache.org) mailing list that you intend to invoke lazy consensus after X hours with no feedback.
1. Start work on the feature. Optionally, you can open a draft PR if you want
feedback during development.
1. Submit your PR for review and reference the blueprint in the PR description.

#### Pull Requests
Once your issue has been approved or your feature blueprint has been merged and
you're ready to start slinging code, we have a few [guidelines](#guidelines)
to help maintain code quality and ensure the pull request process goes smoothly.

If you've never made a pull request, it's super-easy. Github has a great
tutorial [here](https://help.github.com/articles/using-pull-requests/). In a
nutshell, you click the _fork_ button to make a fork, clone it and make your
change, then click the green _New pull request_ button on your repository's
Github page and follow the instructions.

That's it! We'll look at it and get back to you.

Guidelines
----------
Following the project conventions will make the pull request process go faster
and smoother.

#### Component-specific guidelines
Some "components" or pieces of Traffic Control have their own guidelines that
don't generally apply elsewhere.

##### Traffic Ops
If making changes to the Traffic Ops API, please consult the
[Traffic Ops API Guidelines](https://traffic-control-cdn.readthedocs.io/en/latest/development/api_guidelines.html).

##### Traffic Portal v2 (experimental)
The contribution guidelines for the experimental Traffic Portal v2 Angular
project are detailed in
[that project's README](experimental/traffic-portal/README.md).

#### Create an issue or feature blueprint
If you want to add a new feature, make a
[GitHub issue](https://github.com/apache/trafficcontrol/issues) or
[feature blueprint](blueprints/BLUEPRINT_TEMPLATE.md) and discuss it with us
first on the
[dev@trafficcontrol.apache.org](mailto:dev@trafficcontrol.apache.org) mailing
list. We might already be working on it, or there might be an existing way to do
it.

If it's a bug fix, make a
[GitHub issue](https://github.com/apache/trafficcontrol/issues) and optionally
discuss it with us first on the
[dev@trafficcontrol.apache.org](mailto:dev@trafficcontrol.apache.org) mailing
list.

#### Documentation
If your pull request changes the user interface or API, make sure the relevant
[documentation](http://trafficcontrol.apache.org/docs/latest/index.html) is
updated. Documentation [source code](docs/source) is written using
[reStructuredText](https://en.wikipedia.org/wiki/ReStructuredText). Please
verify any document changes by installing
[Sphinx](http://www.sphinx-doc.org/en/stable/) and running 'make' from the
[root of the docs directory](docs).

#### Code formatting
Keep functions small. Big functions are hard to read, and hard to review. Try to
make your changes look like the surrounding code, and follow language
conventions. For Go, run `go fmt` and `go vet`. For Perl, `perltidy`. For Java,
[PMD](https://pmd.github.io).

#### One pull request per feature
Like big functions, big pull requests are just hard to review. Make each pull
request as small as possible. For example, if you're adding ten independent API
endpoints, make each a separate pull request. If you're adding interdependent
functions or endpoints to multiple components, make a pull request for each,
starting at the lowest level.

#### Tests
Make sure all existing tests pass. If you change the way something works, be
sure to update tests to reflect the change. Add unit tests for new functions,
and add integration tests for new interfaces.

Tests that fail if your feature doesn't work are much more useful than tests
which only validate best-case scenarios.

We're in the process of adding more tests and testing frameworks, so if a
testing framework doesn't exist for the component you're changing, don't worry
about it.

#### Commit messages
Try to make your commit messages follow
[git best practices](http://chris.beams.io/posts/git-commit/).
1. Separate subject from body with a blank line
1. Limit the subject line to 50 characters
1. Capitalize the subject line
1. Do not end the subject line with a period
1. Use the imperative mood in the subject line
1. Wrap the body at 72 characters
1. Use the body to explain what and why vs. how

This makes it easier for people to read and understand what each commit does, on
both the command line interface and GitHub.

---

Don't let all these guidelines discourage you, we're more interested in
community involvement than perfection, in general.

What are you waiting for? Get hacking!
