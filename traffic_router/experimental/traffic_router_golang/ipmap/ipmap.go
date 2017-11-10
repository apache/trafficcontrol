package ipmap

import (
	"errors"
	"net"

	"github.com/apache/incubator-trafficcontrol/traffic_router/experimental/traffic_router_golang/coveragezone"

	"github.com/apache/incubator-trafficcontrol/lib/go-tc"
)

// TODO implement, with Coverage Zone File + Maxmind DB

var ErrNotFound = errors.New("not found")

type LatLon struct {
	Lat float64
	Lon float64
}

// DummyLocations returns dummy locations for IPs.
// TODO remove when real geo lookup is implemented
func DummyLocations() map[string]LatLon {
	return map[string]LatLon{
		"127.0.0.1":   {39.579244, -104.934282},
		"192.168.0.1": {39.579244, -104.934282},
		"::1":         {39.579244, -104.934282},
	}
}

func New() coveragezone.CoverageZone {
	return &ipmap{}
}

type ipmap struct{}

// Get takes an IP and returns the Latitude and Longitude.
func (i *ipmap) Get(ip net.IP) (tc.CRConfigLatitudeLongitude, bool) {
	// TODO: get IP via req.RemoteAddr => net.SplitHostPort => net.ParseIP
	locations := DummyLocations()
	ll, ok := locations[ip.String()]
	if !ok {
		return tc.CRConfigLatitudeLongitude{}, false
	}
	return tc.CRConfigLatitudeLongitude{Lat: ll.Lat, Lon: ll.Lon}, true
}
