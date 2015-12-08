/*
 * Copyright 2015 Comcast Cable Communications Management, LLC
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

package com.comcast.cdn.traffic_control.traffic_router.core.loc;

import java.io.File;
import java.io.FileNotFoundException;
import java.io.FileReader;
import java.net.InetAddress;
import java.net.UnknownHostException;
import java.util.Map;
import java.util.TreeMap;

import com.comcast.cdn.traffic_control.traffic_router.core.util.CidrAddress;
import org.apache.log4j.Logger;
import org.apache.wicket.ajax.json.JSONArray;
import org.apache.wicket.ajax.json.JSONException;
import org.apache.wicket.ajax.json.JSONObject;
import org.apache.wicket.ajax.json.JSONTokener;
import org.w3c.dom.Document;
import org.w3c.dom.Node;
import org.w3c.dom.NodeList;

import com.comcast.cdn.traffic_control.traffic_router.core.cache.CacheLocation;

public class NetworkNode implements Comparable<NetworkNode> {
	private static final Logger LOGGER = Logger.getLogger(NetworkNode.class);

	private static final String DEFAULT_SUB_STR = "0.0.0.0/0";


	private CidrAddress cidrAddress;
	private String loc;
	String source = "";
	protected Map<NetworkNode,NetworkNode> children;

	public NetworkNode(final String str) throws NetworkNodeException {
		this(str, null);
	}

	public NetworkNode(final String str, final String loc) throws NetworkNodeException {
		this.source = str;
		this.loc = loc;

		cidrAddress = CidrAddress.fromString(str);
	}

	public String toString() {
		String str = "";
		try {
			str = InetAddress.getByAddress(cidrAddress.getHostBytes()).toString().replace("/", "");
		} catch (UnknownHostException e) {
			LOGGER.warn(e,e);
		}
		return "["+str+"/"+ cidrAddress.getNetmaskLength()+"] - location:" + this.getLoc();
	}

	public NetworkNode getNetwork(final String ip) throws NetworkNodeException {
		return getNetwork(new NetworkNode(ip));
	}
	public NetworkNode getNetwork(final NetworkNode ipnn) {
		if(this.compareTo(ipnn)!=0) { return null; }// not a match
		if(children == null) { return this; }

		final NetworkNode c = children.get(ipnn);
		if(c==null) { return this; }
		return c.getNetwork(ipnn);
	}

	@Override
	public int compareTo(final NetworkNode o) {
		return cidrAddress.compareTo(o.cidrAddress);
	}

	public Boolean add(final NetworkNode nn) {
		synchronized(this) {
			if(children == null) {
				children = new TreeMap<NetworkNode,NetworkNode>();
			}
			return add(children, nn);
		}
	}
	protected Boolean add(final Map<NetworkNode,NetworkNode> children, final NetworkNode nn) {
		if(compareTo(nn)!=0) {
			return false;
		}
		final NetworkNode child = children.get(nn);
		if(child == null) {
			children.put(nn,nn);
			return true;
		}

		if (child.cidrAddress.getNetmaskLength() == nn.cidrAddress.getNetmaskLength()) {
			return false;
		}

		// one is a subnet of another...

		if (child.cidrAddress.getNetmaskLength() < nn.cidrAddress.getNetmaskLength()) {
			child.add(nn);
			return true;
		}

		// swap
		nn.add(child);
		children.remove(child);
		children.put(nn, nn);
		return true;
	}

	public static NetworkNode generateTree(final File f) 
			throws NetworkNodeException, FileNotFoundException, JSONException  {

			final JSONObject json = new JSONObject(new JSONTokener(new FileReader(f)));
			return generateTree(json);
	}
	public static class SuperNode extends NetworkNode {
		private Map<NetworkNode, NetworkNode> children6;

		public SuperNode() throws NetworkNodeException {
			super(DEFAULT_SUB_STR);
		}

		public Boolean add6(final NetworkNode nn) {
			if(children6 == null) {
				children6 = new TreeMap<NetworkNode,NetworkNode>();
			}
			return add(children6, nn);
		}
		public NetworkNode getNetwork(final String ip) throws NetworkNodeException {
			final NetworkNode nn = new NetworkNode(ip);
			if (nn.cidrAddress.isIpV6()) {
				return getNetwork6(nn);
			}
			return getNetwork(nn);
		}
		public NetworkNode getNetwork6(final NetworkNode ipnn) {
			if(children6 == null) { return this; }

			final NetworkNode c = children6.get(ipnn);
			if(c==null) { return this; }
			return c.getNetwork(ipnn);
		}
	}

	@SuppressWarnings("PMD.CyclomaticComplexity")
	private static NetworkNode generateTree(final JSONObject json) {
		try {
			final JSONObject coverageZones = json.getJSONObject("coverageZones");

			final SuperNode root = new SuperNode();
			instance = root;

			for (String loc : JSONObject.getNames(coverageZones)) {
				final JSONObject locData = coverageZones.getJSONObject(loc);

				try {
					final JSONArray network6 = locData.getJSONArray("network6");

					for (int i = 0; i < network6.length(); i++) {
						final String ip = network6.getString(i);

						try {
							root.add6(new NetworkNode(ip, loc));
						} catch (NetworkNodeException ex) {
							LOGGER.error(ex, ex);
						}
					}
				} catch (JSONException ex) {
					LOGGER.warn("An exception was caught while accessing the network6 key of " + loc + " in the incoming coverage zone file: " + ex.getMessage());
				}

				try {
					final JSONArray network = locData.getJSONArray("network");

					for (int i = 0; i < network.length(); i++) {
						final String ip = network.getString(i);

						try {
							root.add(new NetworkNode(ip, loc));
						} catch (NetworkNodeException ex) {
							LOGGER.error(ex, ex);
						}
					}
				} catch (JSONException ex) {
					LOGGER.warn("An exception was caught while accessing the network key of " + loc + " in the incoming coverage zone file: " + ex.getMessage());
				}
			}

			return root;
		} catch (JSONException e) {
			LOGGER.warn(e,e);
		} catch (NetworkNodeException ex) {
			LOGGER.fatal(ex, ex);
		}

		return null;
	}
	public static NetworkNode generateTree(final Document doc) throws NetworkNodeException {

		final NetworkNode root = new NetworkNode(DEFAULT_SUB_STR, null);
		instance = root;
		final NodeList nl = doc.getElementsByTagName("coverageZone");

		// loop coverageZone
		for(int i = 0; i < nl.getLength(); i++) {
			final Node zone = nl.item(i);
			final NodeList nl2 = zone.getChildNodes();
			String loc = null;
			for(int j = 0; j < nl2.getLength(); j++) {
				final Node n = nl2.item(j);
				if(n.getNodeName().equals("location")) {
					loc = n.getTextContent();
					break;
				}
			}
			// loop network
			for(int j = 0; j < nl2.getLength(); j++) {
				final Node n = nl2.item(j);
				if(n.getNodeName().equals("network")) {
					root.add(new NetworkNode(n.getTextContent(), loc));
				}
			}
		}
		return root;
	}

	private static NetworkNode instance;
	public static NetworkNode getInstance() {
		if(instance!=null) { return instance; }
		try {
			instance = new NetworkNode(DEFAULT_SUB_STR);
		} catch (NetworkNodeException e) {
			LOGGER.warn(e);
		}
		return instance;
	}
	public String getLoc() {
		return loc;
	}
	public void setLoc(final String loc) {
		this.loc = loc;
	}

	CacheLocation cacheLocation = null;
	public CacheLocation getCacheLocation() {
		return cacheLocation;
	}
	public void setCacheLocation(final CacheLocation cl2) {
		cacheLocation = cl2;
	}
	public int size() {
		if(children==null) { return 1; }
		int size = 1;
		for(NetworkNode n : children.keySet()) {
			size += n.size();
		}
		return size;
	}
	public void clearCacheCache() {
		synchronized(this) {
			cacheLocation = null;

			if (this instanceof SuperNode) {
				final SuperNode sn = (SuperNode) this;

				if (sn.children6 != null) {
					for (NetworkNode n : sn.children6.keySet()) {
						n.clearCacheCache();
					}
				}
			}

			if (children != null) {
				for(NetworkNode n : children.keySet()) {
					n.clearCacheCache();
				}
			}
		}
	}

}