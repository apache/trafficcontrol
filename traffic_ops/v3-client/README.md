# Traffic Ops Go Client

## Deprecated
The version of the Traffic Ops API supported by this client is deprecated.
Please switch to the `github.com/apache/trafficcontrol/traffic_ops/v4-client`
package.

## Getting Started
1. Obtain the latest version of the library

`go get github.com/apache/trafficcontrol/traffic_ops/client`

2. Get a basic TO session started and fetch a list of CDNs
```go
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/apache/trafficcontrol/lib/go-tc"
	toclient "github.com/apache/trafficcontrol/traffic_ops/v3-client"
)

const TOURL = "http://localhost"
const TOUser = "user"
const TOPassword = "password"
const AllowInsecureConnections = true
const UserAgent = "MySampleApp"
const UseClientCache = false
const TrafficOpsRequestTimeout = time.Second * time.Duration(10)

func main() {
	session, remoteaddr, err := toclient.LoginWithAgent(
		TOURL,
		TOUser,
		TOPassword,
		AllowInsecureConnections,
		UserAgent,
		UseClientCache,
		TrafficOpsRequestTimeout)
	if err != nil {
		fmt.Printf("An error occurred while logging in:\n\t%v\n", err)
		os.Exit(1)
	}
	fmt.Println("Connected to: " + remoteaddr.String())
	var cdns []tc.CDN
	cdns, _, err = session.GetCDNsWithHdr(nil)
	if err != nil {
		fmt.Printf("An error occurred while getting cdns:\n\t%v\n", err)
		os.Exit(1)
	}
	for _, cdn := range cdns {
		fmt.Println(cdn.Name)
	}
}
```
