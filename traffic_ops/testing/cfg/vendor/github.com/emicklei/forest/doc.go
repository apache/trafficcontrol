/*
Package forest has functions for REST Api testing in Go

This package provides a few simple helper types and functions to create functional tests that call HTTP services.
A test uses a forest Client which encapsulates a standard http.Client and a base url.
Such a client can be created inside a function or by initializing a package variable for shared access.
Using a client, you can send http requests and call multiple expectation functions on each response.

Most functions of the forest package take the *testing.T variable as an argument to report any error.

Example
		// setup a shared client to your API
		var chatter = forest.NewClient("http://api.chatter.com", new(http.Client))


		func TestGetMessages(t *testing.T) {
			r := chatter.GET(t, forest.Path("/v1/messages").Query("user","zeus"))
			ExpectStatus(t,r,200)
			ExpectJSONArray(t,r,func(messages []interface{}){

				// in the callback you can validate the response structure
				if len(messages) == 0 {
					t.Error("expected messages, got none")
				}
			})
		}

To compose http requests, you create a RequestConfig value which as a Builder interface for setting
the path,query,header and body parameters. The ProcessTemplate function can be useful to create textual payloads.
To inspect http responses, you use the Expect functions that perform the unmarshalling or use XMLPath or JSONPath functions directly on the response.


If needed, implement the standard TestMain to do global setup and teardown.

	func TestMain(m *testing.M) {
		// there is no *testing.T available, use an stdout implementation
		t := forest.TestingT

		// setup
		chatter.PUT(t, forest.Path("/v1/messages/{id}",1).Body("<payload>"))
		ExpectStatus(t,r,204)

		exitCode := m.Run()

		// teardown
		chatter.DELETE(t, forest.Path("/v1/messages/{id}",1))
		ExpectStatus(t,r,204)

		os.Exit(exitCode)
	}

Special features

In contrast to the standard behavior, the Body of a http.Response is made re-readable.
This means one can apply expectations to a response as well as dump the full contents.

The function XMLPath provides XPath expression support. It uses the [https://godoc.org/launchpad.net/xmlpath] package.
The similar function JSONPath can be used on JSON documents.

Colorizes error output (can be configured using package vars).

All functions can also be used in a setup and teardown as part of TestMain.

(c) 2015, http://ernestmicklei.com. MIT License
*/
package forest
