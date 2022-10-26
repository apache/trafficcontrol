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

# dtp

DTP is the Deterministic Test Protocol. DTP works as an HTTP web server that can generate HTTP responses based on parameters provided in the request. 

### Getting Started
  ___

DTP is enabled by specifying an endpoint in the fakeOrigin configuration with a type of `testing`. The testing capabilities will then be available
at an endpoint named after the `id` specified in the configuration of that endpoint. For example:
```
	"endpoints": [
		{
			"id": "testing_endpoint",
			"type": "testing"
		}
	]
```
can then be reached at http://host:port/testing_endpoint/. Where you can then add any additional parameters to be tested as specified below.

### Usage
___

Desired behavior is requested via url string, with the following format:

```
http://endpoint/~p.<type>/~<key0>.<val0>/~<key1>.<val1> ....
```

Options also passable via custom request header:
```
X-Dtp: ~<key0>.<val0>,~<key1>.<val1>
```

Or via query string:

```
http://endpoint/~p.<type>/&~<key0>.<val0>/~<key1>.<val1> ....
```

A special header, becomes Cache-Control

```
X-Dtp-Cc: public, max-age=200
```

This header explicitly tells DTP what Cache-Control headers to set for the response.

#### Response Manipulation

Currently implemented handler (has priority) (~h.<type>):

- sc ('~sc') (status code)
- hijack ('~payload', '~payload64')

Initial modifier:
- idelay ('~idelay') golang time.Duration
- dly ('~dly') nanoseconds, can use random weighting.

Data types (~p.<type>):
- bin ('~s') - original binary
- binf ('~s') - faster binary
- txt ('~s') - plain txt file
- tex ('~s') - repeated binary texture
- gen3s ('~s') - all 3s

  

Modifiers (applies to Data types):
- rnd ('~rnd') random number for cache control
- rmhdrs ('~rmhdrs') remove specified request headers

Forwarders (requires ~p.<type>) (~f.<type>):
- posevt ('~posevt', 'etags', 'sc')
- delay ('delay') delays all but first block (and header)

Cache control:
- ui ('~ui')
- lm ('~lm')

Size:
~s -  size value (int64)

```
~s.100
~s.10M
~s.((1,30000r)25w,(100000,1000000r)45w,(1M)30w)
```

Random Seed:
~rnd - seed value (int64)

```
~rnd.42
```

Remove specified headers before processing:
~rmhdrs - List of headers
```
~rmhdrs.Range
```

Status Code:
~sc - any valid http status code value

```
~sc.502
```

HTTP Payload:
~payload - any simple string

```
~payload.Hello-World
```

Base 64 Mime Encoded Payload: 
~payload64 - a base64 mime encoded payload string

```
~payload64.SGVsbG8sIFdvcmxkCg== (Hello, World\n)
```

Specifying byte position for event (very useful for range requests, EvalNumber):
~posevt:

- byte pos for event (very useful for range requests, EvalNumber).
- close (closes connection)
- sc (send back status code if range request contains)
- todo: hang the connection
- etags
- first etag for range before
- second etag for range containining/after

```
~f.posevt.1200000.close
~f.posevt.1048577.sc.416
~f.posevt.1000001.etags.foo.bar
```

  
Delay:
~delay- golang parseable delay applied to each block Read request, NOT the First block

```
~f.delay.250ms
```
 #### Header Manipulation
Add Headers (can be combined):
  
~hdr - dot separated pairs

```
~hdr.Foo.Bar
~hdr.Foo.Faz.Bar.Baz
```

~hdr64 - dot separated pairs, with contents base64 mime encoded

```
~hdr64.Cache-Control.bWF4LWFnZT00Mg==  (max-age=42)
```

#### Cache Control
Max Age (Seconds):
~ui - Last-Modified: (UTC / ui) * ui

```
~ui.2000
```

~lm - Last-Modified in UTC seconds. This either needs ~hdr.Cache-Control, ~hdr64.CacheControl or X-Dtp-Cc: header

```
~lm.1586266180
```

~etag - Add an etag header, autogen if blank/empty or as specified

```
~etag.foo
~etag
```

~cksum_req - Standalone option which will md5sum the request header as `X-Request-Header-Cksum`

#### Examples  
500 MB plain text response:
```
curl "http://endpoint/~p.txt/~s.500M/"
```

Cacheable url that cycles every aligned 2000s (Last Modifed, Max-Age, etc):
```
curl "http://endpoint/~p.bin/~s.1G/&~ui.2000/"
```

Sample range call given (transfer 6 bytes then fail):

```
curl "http://endpoint/~p.txt/~f.posevt/~s.2M/&~posevt.1000005.close/" -r 999999-1049100
```

Slice sample call given slice block 1000000 (1 byte then 416 server response):
```
curl "http://endpoint/~p.txt/~f.posevt/~s.2M/&~posevt.1049076.sc.416/" -r 999999-1049100
```


Slice sample call given slice block 1000000 (transfer 6 bytes then change etag):
```
curl "http://endpoint/~p.txt/~s.2M/a&~f.posevt.1000005.etags.foo.bar/" -r 999999-1049100
```

Sample per 32k block delay parameter, not the first block:

```
curl "http://endpoint/~p.txt/~s.2M/&~f.delay.100ms/" -r 999999-1049100
```

Sample last modified header:

```
curl "http://endpoint/~p.txt/~s.1024/~hdr.foo.bar.Next-Modified.RnJpLCAwNyBGZWIgMjAyMCAxNTowNjo0MCBHTVQ="

results in return headers:
Foo: bar
Next-Modified: RnJpLCAwNyBGZWIgMjAyMCAxNTowNjo0MCBHTVQ=
```

  

```

curl "http://endpoint/~p.txt/~s.1024/~hdr64.Next-Modified.RnJpLCAwNyBGZWIgMjAyMCAxNTowNjo0MCBHTVQ="

results in return header:
Next-Modified: Fri, 07 Feb 2020 15:06:40 GMT

```
#### Profiling
DTP may be profiled if so desired.  If profiling is enabled, then DTP endpoints become available under:

http://endpoint/debug/pprof/

http://endpoint/debug/pprof/cmdline

http://endpoint/debug/pprof/profile

http://endpoint/debug/pprof/symbol

http://endpoint/debug/pprof/trace

Current DTP configuration may be dumped via endpoint:

http://endpoint/config
  
Header logging can be flipped on or off with the following:

http://endpoint/config?request_headers=true

http://endpoint/config?response_headers=false

http://endpoint/config?all_headers=true

### Development
___
The "type" identifies which "plugin" to use for handling a response.
This allows developers to drop in their own custom plugins.
A developer needs to include a per connection function with signature:

```
func(http.ResponseWriter, *http.Request, map[string]string)
```

the file needs to register the function as:

```
func init() {
GlobalHandlerFuncs[`type`] = TypeFunc
}
```

A simple example is in "type_hijack.go"
An example curl might be:
```
curl -v "http://endpoint/~h.hijack/&~payload.HelloWorld"
```

or

```
curl -v "http://endpoint/~h.hijack/&~payload64.SGVsbG8gV29ybGQK"
(Hello World with a \n at the the end)
```

With debug turned on you should see the message "Connection hijacked"
displayed.

Another example is type_gen3s.go
