package manager

import (
	"github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/log"
	"github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/peer"
)

// StartPeerManager listens for peer results, and when it gets one, it adds it to the peerStates list, and optimistically combines the good results into combinedStates
func StartPeerManager(peerChan <-chan peer.Result, localStates peer.CRStatesThreadsafe, peerStates peer.CRStatesPeersThreadsafe) peer.CRStatesThreadsafe {
	combinedStates := peer.NewCRStatesThreadsafe()
	go func() {
		for {
			select {
			case crStatesResult := <-peerChan:
				peerStates.Set(crStatesResult.Id, crStatesResult.PeerStats)
				combinedStates.Set(combineCrStates(peerStates.Get(), localStates.Get()))
				crStatesResult.PollFinished <- crStatesResult.PollID
			}
		}
	}()
	return combinedStates
}

// TODO JvD: add deliveryservice stuff
func combineCrStates(peerStates map[string]peer.Crstates, localStates peer.Crstates) peer.Crstates {
	combinedStates := peer.NewCrstates()
	for cacheName, localCacheState := range localStates.Caches { // localStates gets pruned when servers are disabled, it's the source of truth
		downVotes := 0 // TODO JvD: change to use parameter when deciding to be optimistic or pessimistic.
		if localCacheState.IsAvailable {
			// log.Infof(cacheName, " is available locally - setting to IsAvailable: true")
			combinedStates.Caches[cacheName] = peer.IsAvailable{IsAvailable: true} // we don't care about the peers, we got a "good one", and we're optimistic
		} else {
			downVotes++ // localStates says it's not happy
			for _, peerCrStates := range peerStates {
				if peerCrStates.Caches[cacheName].IsAvailable {
					// log.Infoln(cacheName, "- locally we think it's down, but", peerName, "says IsAvailable: ", peerCrStates.Caches[cacheName].IsAvailable, "trusting the peer.")
					combinedStates.Caches[cacheName] = peer.IsAvailable{IsAvailable: true} // we don't care about the peers, we got a "good one", and we're optimistic
					break                                                                  // one peer that thinks we're good is all we need.
				} else {
					// log.Infoln(cacheName, "- locally we think it's down, and", peerName, "says IsAvailable: ", peerCrStates.Caches[cacheName].IsAvailable, "down voting")
					downVotes++ // peerStates for this peer doesn't like it
				}
			}
		}
		if downVotes > len(peerStates) {
			// log.Infoln(cacheName, "-", downVotes, "down votes, setting to IsAvailable: false")
			combinedStates.Caches[cacheName] = peer.IsAvailable{IsAvailable: false}
		}
	}

	for deliveryServiceName, localDeliveryService := range localStates.Deliveryservice {
		deliveryService := peer.Deliveryservice{}
		if localDeliveryService.IsAvailable {
			deliveryService.IsAvailable = true
		}
		deliveryService.DisabledLocations = localDeliveryService.DisabledLocations

		for peerName, iPeerStates := range peerStates {
			peerDeliveryService, ok := iPeerStates.Deliveryservice[deliveryServiceName]
			if !ok {
				log.Warnf("local delivery service %s not found in peer %s\n", deliveryServiceName, peerName)
				continue
			}
			if peerDeliveryService.IsAvailable {
				deliveryService.IsAvailable = true
			}
			deliveryService.DisabledLocations = intersection(deliveryService.DisabledLocations, peerDeliveryService.DisabledLocations)
		}
		combinedStates.Deliveryservice[deliveryServiceName] = deliveryService
	}

	return combinedStates
}
