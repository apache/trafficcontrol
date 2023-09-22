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

# Implementing a new Traffic Vault backend (e.g. Foo)

1. Create a new directory in ./backends which will contain/define your new package (`foo`) which will provide all the functionality to support a new `Foo` backend for Traffic Vault.
2. In this new `./backends/foo` directory, create a new file: `foo.go`, with `package foo` to define the package name.
3. In `foo.go`, define a struct (e.g. `type Foo struct`) which will act as the method receiver for all the required `TrafficVault` methods. This struct should contain any fields necessary to provide the required functionality. For instance, it should most likely contain all the required configuration to connect to and use the `Foo` data store.
```go
type Foo struct {
    cfg Config
}

type Config struct {
    user     string
    password string
}
```
4. Implement all the methods required by the `TrafficVault` interface on your new `Foo` struct. Initially, you may want to simply stub out the methods and implement them later:
```go
func (f *Foo) GetDeliveryServiceSSLKeys(xmlID string, version string, tx *sql.Tx, ctx context.Context) (tc.DeliveryServiceSSLKeysV15, bool, error) {
	return tc.DeliveryServiceSSLKeysV15{}, false, nil
}

... (snip)

func (f *Foo) Ping(tx *sql.Tx, ctx context.Context) (tc.TrafficVaultPingResponse, error) {
	return tc.TrafficVaultPingResponse{}, nil
}
```
5. Define a `trafficvault.LoadFunc` which will parse the given JSON config (from cdn.conf's `traffic_vault_config` option) and return a pointer to an instance of the `Foo` type:
```go
func loadFoo(b json.RawMessage) (trafficvault.TrafficVault, error) {
    // unmarshal the given JSON, validate it, return an error if any
    // fooCfg, err := parseAndValidateConfig(b)
	return &Foo{cfg: fooCfg}, nil
}
```
6. Define a package `init` function which calls `trafficvault.AddBackend` with your backend's name and `LoadFunc` in order to register your new Traffic Vault `Foo` backend for use:
```go
func init() {
	trafficvault.AddBackend("foo", loadFoo)
}
```
7. In `./backends/backends.go`, import your new package:
```go
import (
    _ "github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/trafficvault/backends/foo"
)
```
This is required for the package `init()` function to run and register the new backend.
8. You are now able to test your new Traffic Vault `Foo` backend. First, in `cdn.conf`, you need to set `traffic_vault_backend` to `"foo"` and include your desired `Foo` configuration in `traffic_vault_config`. Once that is done, Traffic Vault is enabled, and you can use Traffic Ops API routes that require Traffic Vault. At this point, you should go back and fully implement the stubbed out `TrafficVault` interface methods on your `Foo` type.
