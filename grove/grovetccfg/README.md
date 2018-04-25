# grovetccfg

Traffic Control configuration generator for the Grove HTTP caching proxy.

# Building

1. Install and set up a Golang development environment.
    * See https://golang.org/doc/install
2. Clone this repository into your GOPATH.
```bash
mkdir -p $GOPATH/src/github.com/apache
cd $GOPATH/src/github.com/apache
git clone https://github.com/apach/incubator-trafficcontrol
```
3. Build the application
```bash
cd $GOPATH/src/github.com/apache/incubator-trafficcontrol/grove/grovetccfg
go build
```
5. Install and configure an RPM development environment
   * See https://wiki.centos.org/HowTos/SetupRpmBuildEnvironment
4. Build the RPM
```bash
./build/build_rpm.sh
```

# Running

The `grovetccfg` tool has an RPM, but no service or config files. It must be run manually, even after installing the RPM. Consider running the tool in a cron job.

Example:

`./grovetccfg -api=1.2 -host my-http-cache -insecure -touser carpenter -topass 'walrus' -tourl https://cdn.example.net -pretty > remap.json`

Flags:

| Flag | Description |
| --- | --- |
| `api` | The Traffic Ops API version to use. The default is 1.2. If 1.3 is passed, it will use a newer and more efficient endpoint. |
| `host` | The Traffic Ops server to create configuration from. This must be a cache server in Traffic Ops. |
| `insecure` | Whether to ignore certificate errors when connecting to Traffic Ops |
| `touser` | The Traffic Ops user to use. |
| `topass` | The Traffic Ops user password. |
| `tourl` | The Traffic Ops URL, including the scheme and fully qualified domain name. |
| `pretty` | Whether to pretty-print JSON |
