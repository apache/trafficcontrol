<!--- vim:fo=tc:tw=0:
--->

Release Notes
=============
* `1.9.1` - [Milestone](https://github.com/basho/riak-go-client/issues?q=milestone%3Ariak-go-client-1.9.1)
* `1.9.0` - [Milestone](https://github.com/basho/riak-go-client/issues?q=milestone%3Ariak-go-client-1.9.0)
* `1.8.0` - [Milestone](https://github.com/basho/riak-go-client/issues?q=milestone%3Ariak-go-client-1.8.0)
* `1.7.0` - Following PRs included:
  * [TimeSeries support](https://github.com/basho/riak-go-client/pull/68)
* `1.6.0` - Following PRs included:
  * [Add `NodeOptions` test and concurrent commands test for `MaxConnections`](https://github.com/basho/riak-go-client/pull/61)
  * [Security auth fixes](https://github.com/basho/riak-go-client/pull/60)
  * [Use request or command timeout as applicable](https://github.com/basho/riak-go-client/pull/57)
  * [Re-tryable vs non-re-tryable commands](https://github.com/basho/riak-go-client/pull/56)
  * [Re-try reads on temporary network errors, bug fixes in error situations](https://github.com/basho/riak-go-client/pull/52)
* `1.5.1` (DEPRECATED RELEASE) - Following PRs addressed:
  * [Improve connection error handling](https://github.com/basho/riak-go-client/pull/48)
* `1.5.0` - Following PRs addressed:
  * [Add `FetchBucketTypePropsCommand` and `StoreBucketTypePropsCommand`](https://github.com/basho/riak-go-client/pull/42)
* `1.4.0` - Following issues / PRs addressed:
  * [Add `ResetBucketCommand`](https://github.com/basho/riak-go-client/pull/35)
  * [Legacy Counter support](https://github.com/basho/riak-go-client/pull/33)
* `1.3.0` - Following issues / PRs addressed:
  * [Add `NoDefaultNode` option to `ClusterOptions`](https://github.com/basho/riak-go-client/pull/28)
  * [`ConnectionManager` / `NodeManager` fixes](https://github.com/basho/riak-go-client/pull/25)
  * [`ConnectionManager` expiration fix](https://github.com/basho/riak-go-client/issues/23)
* `1.2.0` - Following issues / PRs addressed:
  * [Conflict resolver not being passed to Fetch/Store-ValueCommand](https://github.com/basho/riak-go-client/issues/21)
  * [Reduce exported API](https://github.com/basho/riak-go-client/pull/20)
  * [Modify ClientError to trap an inner error if necessary](https://github.com/basho/riak-go-client/pull/19)
* `1.1.0` - Following issues / PRs addressed:
  * [Issues with incrementing counters within Maps](https://github.com/basho/riak-go-client/issues/17)
  * [Extra goroutine in Execute](https://github.com/basho/riak-go-client/issues/16)
  * [Execute does not return error correctly](https://github.com/basho/riak-go-client/isues/15)
* `1.0.0` - Initial release with Riak 2.0 support.
* `1.0.0-beta11 - Initial beta release with Riak 2 support. Command queuing and retrying not implemented yet.

