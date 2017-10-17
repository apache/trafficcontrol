/*
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package com.comcast.cdn.traffic_control.traffic_monitor.health;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Set;

import org.apache.log4j.Logger;

import com.comcast.cdn.traffic_control.traffic_monitor.config.Cache;
import com.comcast.cdn.traffic_control.traffic_monitor.config.Peer;
import com.comcast.cdn.traffic_control.traffic_monitor.health.Event.EventType;

public class PeerState extends AbstractState {
	private static final Logger LOGGER = Logger.getLogger(PeerState.class);
	private static final long serialVersionUID = 1L;
	private static Map<String, PeerState> states = new HashMap<String, PeerState>();
	private static Map<String, Boolean> overrideMap = new HashMap<String, Boolean>();
	private Peer peer;
	private boolean reachable = false;

	PeerState(final String id) {
		super(id);
	}

	public static List<PeerState> getPeerStates() {
		synchronized(states) {
			return new ArrayList<PeerState>(states.values());
		}
	}

	public static PeerState getOrCreate(final Peer peer) {
		return getOrCreate(peer.getHostname(), peer);
	}

	public static PeerState getOrCreate(final String host, final Peer peer) {
		synchronized(states) {
			PeerState ps = states.get(host);

			if (ps == null) {
				ps = new PeerState(host);
				states.put(host, ps);
			}

			ps.setPeer(peer);
			peer.setState(ps);

			return ps;
		}
	}

	public static PeerState getState(final String host) {
		synchronized(states) {
			return states.get(host);
		}
	}

	private void setPeer(final Peer peer) {
		this.peer = peer;
	}

	public Peer getPeer() {
		return peer;
	}

	public static boolean has(final String host) {
		if (states.get(host) == null) {
			return false;
		}

		return true;
	}

	public static void removeAllBut(final Set<String> peerSet) {
		synchronized(states) {
			for (String key : new ArrayList<String>(states.keySet())) {
				if (!peerSet.contains(key)) {
					states.remove(key);
				}
			}
		}
	}

	public static String get(final String stateId, final String key) {
		return getState(stateId).getLastValue(key);
	}

	public static int getPeerCount() {
		return states.size();
	}

	public static boolean hasPeers() {
		if (states.isEmpty()) {
			return false;
		} else {
			return true;
		}
	}

	public static int getOnlinePeerCount() {
		int onlineCount = 0;

		for (PeerState peerState : PeerState.getPeerStates()) {
			if (!peerState.isReachable()) {
				continue;
			} else {
				onlineCount++;
			}
		}

		return onlineCount;
	}

	public static boolean hasOnlinePeers() {
		if (getOnlinePeerCount() > 0) {
			return true;
		} else {
			return false;
		}
	}

	public static boolean isCacheAvailableOnAnyPeer(final Cache c) {
		final List<Peer> onlineList = getCacheAvailableOnPeers(c);

		if (onlineList != null && !onlineList.isEmpty()) {
			return true;
		} else {
			return false;
		}
	}

	public static List<Peer> getCacheAvailableOnPeers(final Cache c) {
		final List<Peer> onlineList = new ArrayList<Peer>();

		if (PeerState.hasPeers()) {
			for (PeerState peerState : PeerState.getPeerStates()) {
				if (!peerState.isReachable()) {
					continue;
				}

				final Peer peer = peerState.getPeer();
				final String pAvailability = peerState.getLastValue(c.getHostname());

				if (pAvailability == null || pAvailability.equals("unknown")) {
					continue;
				}

				final boolean pIsAvailable = Boolean.parseBoolean(pAvailability);

				if (pIsAvailable == true) {
					LOGGER.debug(String.format("ERROR: %s - isAvailable set to %s from: %s", c.getHostname(), String.valueOf(pIsAvailable), peer.getId()));
					onlineList.add(peer);
				}
			}
		}

		return onlineList;
	}

	public static void logOverride(final Cache c) {
		final Boolean state = overrideMap.get(c.getFqdn());
		final EventType type = EventType.CACHE_STATE_CHANGE;
		type.setType(c.getType());

		if (PeerState.hasOnlinePeers()) {
			final List<Peer> onlineList = PeerState.getCacheAvailableOnPeers(c);

			if (!onlineList.isEmpty()) {
				if (state == null || !state.booleanValue()) {
					final StringBuffer msg = new StringBuffer("Health protocol override condition detected; healthy on (at least) ");
					msg.append(Arrays.toString(onlineList.toArray()).replaceAll("\\[|\\]", ""));
					Event.logStateChange(c.getHostname(), type, true, msg.toString());
					overrideMap.put(c.getFqdn(), true);
				}
			} else if (onlineList.isEmpty()) {
				if (state == null || state.booleanValue()) {
					final StringBuffer msg = new StringBuffer("Health protocol override condition irrelevant; not online on any peers");

					if (c.isAvailableKnown() && c.isAvailable()) {
						msg.append("; healthy locally");
					} else if (c.isAvailableKnown() && !c.isAvailable()) {
						msg.append("; unhealthy locally");
					} else {
						msg.append("; local state unknown");
					}

					Event.logStateChange(c.getHostname(), type, c.isAvailable(), msg.toString());
					overrideMap.put(c.getFqdn(), false);
				}
			}
		} else if (state != null && state.booleanValue()) {
			final StringBuffer msg = new StringBuffer("Health protocol override condition irrelevant; no peers online");
			Event.logStateChange(c.getHostname(), type, c.isAvailable(), msg.toString());
			overrideMap.put(c.getFqdn(), false);
		}
	}

	public static void clearOverride(final Cache c) {
		if (overrideMap.containsKey(c.getFqdn())) {
			final Boolean state = overrideMap.get(c.getFqdn());

			if (state.booleanValue()) {
				final EventType type = EventType.CACHE_STATE_CHANGE;
				type.setType(c.getType());
				Event.logStateChange(c.getHostname(), type, true, "Health protocol override condition cleared; healthy locally");
			}

			overrideMap.remove(c.getFqdn());
		}
	}

	public boolean isReachable() {
		return reachable;
	}

	public void setReachable(final boolean reachable) {
		setReachable(reachable, null);
	}

	public void setReachable(final boolean reachable, final String reason) {
		if (isReachable() != reachable) {
			final StringBuilder sb = new StringBuilder();

			if (reason == null && reachable) {
				sb.append("Peer is reachable");
			} else if (reason == null && !reachable) {
				sb.append("Peer is unreachable");
			} else if (reason != null) {
				sb.append(reason);
			}

			Event.logStateChange(peer.getHostname(), EventType.PEER_STATE_CHANGE, reachable, sb.toString());
		}

		this.reachable = reachable;
	}
}
