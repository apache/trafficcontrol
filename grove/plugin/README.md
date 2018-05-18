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

To add a plugin, create a new `.go` file in the `grove/plugin` directory. This file should have a unique name, to avoid conflicts. Consider prefixing it with your company name, website, or a UUID.

The filename, sans `.go`, is the name of your plugin, and will be the key used for configuration in the remap file. For example, if your file is `f49e54fc-fd17-4e1c-92c6-67028fde8504-hello-world.go`, the name of your plugin is `f49e54fc-fd17-4e1c-92c6-67028fde8504-hello-world`.

Plugins are registered via calls to `AddPlugin` inside an `init` function in the plugin's file.

The `Funcs` object contains functions for each hook, as well as a load function for loading configuration from the remap file. The current hooks are `startup`, `onRequest`, `beforeParentRequest`, `beforeRespond`, and `afterRespond`. If your plugin does not use a hook, it may be nil.

* `startup` is called when the application starts. Examples are set global data, or start a global goroutine needed by the plugin.

* `onRequest` is called immediately when a request is received. It returns a boolean indicating whether to stop processing. Examples are IP blocking, or serving custom endpoints for statistics or to invalidate a cache entry.

* `beforeCacheLookUp` is called immedidiately before looking the object up in the cache. It can be used to modify the cacheKey to be used to for this object using the passed `CacheKeyOverrideFunc` func. Once set using that function Grove will keep using that cacheKey throughout the life of the object in the cache.

* `beforeParentRequest` is called immediately before making a request to a parent. It may manipulate the request being made to the parent. Examples are removing headers in the client request such as `Range`.

* `beforeRespond` is called immediately before responding to a client. It may manipulate the code, headers, and body being returned. Examples are header modifications, or handling if-modified-since requests.

* `afterRespond` is called immediately after responding to the client. Examples are recording stats, or writing to an access log.

* `load` is not a hook, but rather a function to load arbitrary data from the remap config file. It is given a `json.RawMessage`, and can return any object. The object it returns is then passed to this plugin's hooks.

The simplest example is the `hello_world` plugin. See `grove/plugin/hello_world.go`.

```go
func init() {
	AddPlugin(10000, Funcs{startup: hello})
}

func hello(icfg interface{}, d StartupData) {
	log.Errorf("Hello World! I'm a startup plugin! We're starting with %v bytes!\n", d.Config.CacheSizeBytes)
}
```

The plugin is initialized via `AddPlugin`, and its `hello` function is set as the `startup` hook. The `hello` function has the signature of `plugin.StartupHook`.

To pass data from one hook in your plugin to a hook that's called later, set the `Context` member of the `Data` object. For an example, see `hello_context.go`:

```go
func init() {
	AddPlugin(10000, Funcs{startup: helloCtxStart, afterRespond: helloCtxAfterResp})
}

func helloCtxStart(icfg interface{}, d StartupData) {
	*d.Context = 42
	log.Debugf("Hello World! Start set context: %+v\n", d.Context)
}

func helloCtxAfterResp(icfg interface{}, d AfterRespondData) {
	ictx := d.Context
	ctx, ok := (*ictx).(int)
	log.Debugf("Hello World! After Response got context: %+v %+v\n", ok, ctx)
}
```

The `startup` hook function `helloCtxStartup` sets the context pointer to the value `42`, and then the same plugin's `afterRespond` hook `helloCtxAfterResp` retrieves the value from its Data `.Context` pointer.

For an example of configuration, see `modify_headers.go`.

```go
func init() {
	AddPlugin(10000, Funcs{load: modRespHdrLoad, beforeRespond: modRespHdr})
}

type Hdr struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type ModHdrs struct {
	Set  []Hdr    `json:"set"`
	Drop []string `json:"drop"`
}

func modRespHdrLoad(b json.RawMessage) interface{} {
	cfg := ModHdrs{}
	err := json.Unmarshal(b, &cfg)
	if err != nil {
		log.Errorln("modifyheaders loading config, unmarshalling JSON: " + err.Error())
		return nil
	}
	log.Debugf("modifyheaders load success: %+v\n", cfg)
	return &cfg
}

func modRespHdr(icfg interface{}, d BeforeRespondData) {
	if icfg == nil {
		log.Debugln("modifyheaders has no config, returning.")
		return
	}
	cfg, ok := icfg.(*ModHdrs)
	if !ok {
		// should never happen
		log.Errorf("modifyheaders config '%v' type '%T' expected *ModHdrs\n", icfg, icfg)
		return
	}
...
}
```

The load function `modRespHdrLoad` unmarshals the given `json.RawMessage` into its expected type, and returns a pointer to the object. The beforeRespond hook function `modRespHdr` then casts the interface it's given to the type its load function returned.

Load functions are given configuration objects from the remap file, for the rule being processed, under the key `plugins`. Recall the name of the plugin is the filename, with `.go` removed. For example, for the `modify_headers.go` plugin:

```
{
  "rules": [
    {
      "name": "my-remap-rule",
      "plugins": {
        "modify_headers": {
          "set": [{"name": "Server", "value": "grove42"}],
          "drop": ["X-Powered-By"]
        }
      },
  ...
```

If a given rule does not have a plugin key, the global object will be used, as with other remap fields. For example:

```
{
  "plugins": {
    "modify_headers": {
    "set": [{"name": "Server", "value": "grove42"}],
    "drop": ["X-Powered-By"]
  },
  "rules": [
  ...
```
