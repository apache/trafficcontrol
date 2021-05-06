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

# Adding a Plugin

To add a plugin, create a new `.go` file in the `cache-config/t3c-generate/plugin` directory. This file should have a unique name, to avoid conflicts. Consider prefixing it with your company name, website, or a UUID.

The filename, sans `.go`, is the name of your plugin, and will be the key used for configuration in the remap file. For example, if your file is `f49e54fc-fd17-4e1c-92c6-67028fde8504-hello-world.go`, the name of your plugin is `f49e54fc-fd17-4e1c-92c6-67028fde8504-hello-world`.

Plugins are registered via calls to `AddPlugin` inside an `init` function in the plugin's file. The `AddPlugin` function takes a priority, and a set of hook functions. The priority is the order in which plugins are called, starting from 0. Note the priority of plugins included with Traffic Control use a base priority of 10000, unless priority order matters for them.

The `Funcs` object contains functions for each hook. The current hooks are `startup` and `onRequest`. If your plugin does not use a hook, it may be nil.

* `startup` is called when the application starts.

* `onRequest` is called immediately when a request is received. It returns a boolean indicating whether to stop processing.

The simplest example is the `hello_world` plugin. See `plugin/hello_world.go`.

```go
import (
	"strings"
)
func init() {
	AddPlugin(10000, Funcs{onRequest: hello})
}
const HelloPath = "/_hello_world"
func hello(d OnRequestData) IsRequestHandled {
	if d.Cfg.TOURL.Path != HelloPath {
		return RequestUnhandled
	}
	cfgFile := "Hello World!\n"
	fmt.Println(cfgFile)
	os.Exit(42)
	return RequestHandled
}

```

The plugin is initialized via `AddPlugin`, and its `hello` function is set as the `onRequest` hook. The `hello` function has the signature of `plugin.OnRequestFunc`.

# Examples

Example plugins are included in the `/plugin` directory

*hello_world*: Example of a simple HTTP endpoint.

# Glossary

Definitions of terms used in this document.

*Plugin*: A self-contained component whose code is executed when certain events in the main application occur.
*Hook*: A plugin function which is called when a certain event happens in the main application.
*Plugin Data*: Application data given to a plugin, as a function parameter passed to a hook function, including configuration data, running state, and HTTP request state.
