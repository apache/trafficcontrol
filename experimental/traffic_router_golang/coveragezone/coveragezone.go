package coveragezone

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net"

	"github.com/apache/incubator-trafficcontrol/lib/go-tc"
)

// TODO put in lib/go-tc
type JSONCoverageZones struct {
	CoverageZones map[tc.CacheGroupName]JSONCoverageZoneCacheGroup `json:"coverageZones"`
	CustomerName  string                                           `json:"customerName"`
	Revision      string                                           `json:"revision"` // TODO change to Time with Marshal func?
}

type JSONCoverageZoneCacheGroup struct {
	Coordinates tc.CRConfigLatitudeLongitude // TODO rename CRConfigLatitudeLongitude
	Network     []string                     // TODO change to IPNet with Unmarshal func
	Network6    []string
}

type CoverageZone interface {
	Get(ip net.IP) (tc.CRConfigLatitudeLongitude, bool)
}

type netPos struct {
	net *net.IPNet
	pos tc.CRConfigLatitudeLongitude
}

type coverageZone struct {
	nets  []netPos
	net6s []netPos
}

func (c *coverageZone) Get(ip net.IP) (tc.CRConfigLatitudeLongitude, bool) {
	nets := []netPos(nil)
	if ip.To4() != nil {
		nets = c.nets
	} else {
		nets = c.net6s
	}
	for _, netPos := range nets {
		if netPos.net.Contains(ip) {
			return netPos.pos, true
		}
	}
	return tc.CRConfigLatitudeLongitude{}, false
}

func New(jcz JSONCoverageZones) (CoverageZone, error) {
	c := coverageZone{}
	for cg, jcg := range jcz.CoverageZones {
		for _, cidr := range jcg.Network {
			_, net, err := net.ParseCIDR(cidr)
			if err != nil {
				return nil, errors.New("error parsing cachegroup '" + string(cg) + "' Network '" + cidr + ": " + err.Error())
			}
			c.nets = append(c.nets, netPos{net: net, pos: jcg.Coordinates})
		}
		for _, cidr := range jcg.Network6 {
			_, network, err := net.ParseCIDR(cidr)
			if err != nil {
				return nil, errors.New("error parsing cachegroup '" + string(cg) + "' Network6 '" + cidr + ": " + err.Error())
			}
			c.net6s = append(c.net6s, netPos{net: network, pos: jcg.Coordinates})
		}
	}
	return &c, nil
}

func Load(filename string) (CoverageZone, error) {
	f, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.New("reading file '" + filename + "':" + err.Error())
	}
	jcz := JSONCoverageZones{}
	if err := json.Unmarshal(f, &jcz); err != nil {
		return nil, errors.New("unmarshalling JSON '" + filename + "':" + err.Error())
	}
	return New(jcz)
}
