# Go management extensions

## Installation
  
	go get github.com/davecheney/gmx

## Getting started

Instrumenting your application with gmx is as simple as importing the `gmx` package in your `main` package via the side effect operator.

	package main

	import _ "github.com/davecheney/gmx"

By default gmx opens a unix socket in `/tmp`, the name of the socket is

	/tmp/.gmx.$PID.0

## Protocol version 0

The current protocol version is 0, which is a simple JSON based protocol. You can communicate with the gmx socket using a tool like socat.

	% socat UNIX-CONNECT:/tmp/.gmx.$(pgrep godoc).0 stdin
	["runtime.version", "runtime.numcpu"]
	{"runtime.numcpu":4,"runtime.version":"weekly.2012-01-27 11688+"}
     
The request is a json array of strings representing keys that you wish to query. The result is a json map, the keys of that map are keys that matched the keys in your request. The value of the entry will be the result of the published function, encoded in json. If there is no matching key registered, no entry will appear in the result map.

For convenience a client is included in the gmxc sub directory. Please consult the `README` in that directory for more details.

## Registering gmx keys

New keys can be registered using the `Publish` function

	gmx.Publish(key string, f func() interface{})

`f` can be any function that returns a json encodable result. `f` is executed whenever its key is invoked, responsibility for ensuring the function is thread safe rests with the author of `f`.

## Runtime instrumentation

By default gmx instruments selected values from the  `runtime` and `os` packages, refer to the `runtime.go` and `os.go` source for more details.

## Changelog

6/Feb/2012 

+	gmx now honors the value of os.TempDir() when opening the unix socket
+	gmxc now accepts regexps for key names

5/Feb/2012 

+	Initial release
