package datareq

import (
	"net/url"

	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/peer"
)

func srvTRState(params url.Values, localStates peer.CRStatesThreadsafe, combinedStates peer.CRStatesThreadsafe) ([]byte, error) {
	if _, raw := params["raw"]; raw {
		return srvTRStateSelf(localStates)
	}
	return srvTRStateDerived(combinedStates)
}

func srvTRStateDerived(combinedStates peer.CRStatesThreadsafe) ([]byte, error) {
	return peer.CrstatesMarshall(combinedStates.Get())
}

func srvTRStateSelf(localStates peer.CRStatesThreadsafe) ([]byte, error) {
	return peer.CrstatesMarshall(localStates.Get())
}
