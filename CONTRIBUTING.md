**Contributing to Traffic Control**
=================

Thanks again for your time and interest in this project!

The following document is a set of guidelines to streamline the contribution process for our contributors. Please feel free to suggest changes to this document in a pull request!

Things to know before getting started
-------------------------------------
#### Code of Conduct
Please review the [code of conduct](). By participating, you agree to it and are expected to follow this code. Please report any issues or unacceptable behavior to [admin email]

#### Design Decisions
When we need to make changes to the project, we first discuss it on the [users@trafficcontrol.incubator.apache.org](mailto:users@trafficcontrol.incubator.apache.org) mailing list. We document our decisions as well as brainstorm project related ideas in the [wiki](https://cwiki.apache.org/confluence/display/TC/Traffic+Control+Home).

How to Contribute
-------------------------------------

We love pull requests! We simply don't have the time or resources to add every feature and support every platform. If you have improvements (enhancements or bug fixes) for Traffic Control, start by creating an [issue](https://issues.apache.org/jira/browse/TC) and discussing it with us first on the [users@trafficcontrol.incubator.apache.org](mailto:users@trafficcontrol.incubator.apache.org) mailing list. We might already be working on it, or there might be an existing way to do it.

Once your issue has been approved and you're ready to start slinging code, we have a few [guidelines](https://github.com/at9418/incubator-trafficcontrol/edit/master/CONTRIBUTING.md#guidelines) to help maintain code quality and ensure the pull request process goes smoothly.

Remember, your code doesn't have to be perfect. Hack it together and submit a [pull request](https://help.github.com/articles/using-pull-requests/). We'll work with you to make sure it fits properly into the project.

#### Making a pull request
If you've never made a pull request, it's super-easy. Github has a great tutorial [here](https://help.github.com/articles/using-pull-requests/). In a nutshell, you click the _fork_ button to make a fork, clone it and make your change, then click the green _New pull request_ button on your repository's Github page and follow the instructions. That's it! We'll look at it and get back to you.

Guidelines
----------
Following the project conventions will make the pull request process go faster and smoother. If making changes to the Traffic Ops API, please consult the [Traffic Ops API Guidelines](https://cwiki.apache.org/confluence/display/TC/API+Guidelines).

#### Create an issue

If you want to add a new feature, make a [JIRA issue](https://issues.apache.org/jira/browse/TC) and discuss it with us first on the [users@trafficcontrol.incubator.apache.org](mailto:users@trafficcontrol.incubator.apache.org) mailing list. We might already be working on it, or there might be an existing way to do it.

If it's a bug fix, make a [JIRA issue](https://issues.apache.org/jira/browse/TC) and discuss it with us first on the [users@trafficcontrol.incubator.apache.org](mailto:users@trafficcontrol.incubator.apache.org) mailing list. We need to know what the problem is and how to reproduce it so please create a [JIRA issue](https://issues.apache.org/jira/browse/TC) for that as well.

#### Documentation

If your pull request changes the user interface or API, make sure the relevant documentation is updated.

#### Code formatting

Keep functions small. Big functions are hard to read, and hard to review. Try to make your changes look like the surrounding code, and follow language conventions. For Go, run `gofmt` and `go vet`. For Perl, `perltidy`. For Java, [PMD](https://pmd.github.io).

#### One pull request per feature

Like big functions, big pull requests are just hard to review. Make each pull request as small as possible. For example, if you're adding ten independent API endpoints, make each a separate pull request. If you're adding interdependent functions or endpoints to multiple components, make a pull request for each, starting at the lowest level.

#### Tests

Make sure all existing tests pass. If you change the way something works, be sure to update tests to reflect the change. Add unit tests for new functions, and add integration tests for new interfaces.

Tests that fail if your feature doesn't work are much more useful than tests which only validate best-case scenarios.

We're in the process of adding more tests and testing frameworks, so if a testing framework doesn't exist for the component you're changing, don't worry about it.

#### Commit messages

Try to make your commit messages follow [git best practices](http://chris.beams.io/posts/git-commit/).
1. Separate subject from body with a blank line
2. Limit the subject line to 50 characters
3. Capitalize the subject line
4. Do not end the subject line with a period
5. Use the imperative mood in the subject line
6. Wrap the body at 72 characters
7. Use the body to explain what and why vs. how

This makes it easier for people to read and understand what each commit does, on both the command line interface and Github.com.

---

Don't let all these guidelines discourage you, we're more interested in community involvement than perfection.

What are you waiting for? Get hacking!
