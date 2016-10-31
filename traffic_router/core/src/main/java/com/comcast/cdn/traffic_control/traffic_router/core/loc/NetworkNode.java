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

package com.comcast.cdn.traffic_control.traffic_router.core.loc;

import java.io.File;
import java.io.FileNotFoundException;
import java.io.FileReader;
import java.net.InetAddress;
import java.net.UnknownHostException;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import java.util.TreeMap;

import com.comcast.cdn.traffic_control.traffic_router.core.util.CidrAddress;
import com.comcast.cdn.traffic_control.traffic_router.geolocation.Geolocation;

import org.apache.log4j.Logger;
import org.apache.wicket.ajax.json.JSONArray;
import org.apache.wicket.ajax.json.JSONException;
import org.apache.wicket.ajax.json.JSONObject;
import org.apache.wicket.ajax.json.JSONTokener;

import com.comcast.cdn.traffic_control.traffic_router.core.cache.CacheLocation;

public class NetworkNode implements Comparable<NetworkNode> {
    private static final Logger LOGGER = Logger.getLogger(NetworkNode.class);
    private static final String DEFAULT_SUB_STR = "0.0.0.0/0";

    private static NetworkNode instance;

    private CidrAddress cidrAddress;
    private String loc;
    private CacheLocation cacheLocation = null;
    private Geolocation geolocation = null;
    protected Map<NetworkNode,NetworkNode> children;

    public static NetworkNode getInstance() {
        if (instance != null) {
            return instance;
        }

        try {
            instance = new NetworkNode(DEFAULT_SUB_STR);
        } catch (NetworkNodeException e) {
            LOGGER.warn(e);
        }

        return instance;
    }

    public static NetworkNode generateTree(final File f, final boolean verifyOnly) throws NetworkNodeException, FileNotFoundException, JSONException  {
        return generateTree(new JSONObject(new JSONTokener(new FileReader(f))), verifyOnly);
    }

    @SuppressWarnings("PMD.CyclomaticComplexity")
    public static NetworkNode generateTree(final JSONObject json, final boolean verifyOnly) {
        try {
            final JSONObject coverageZones = json.getJSONObject("coverageZones");

            final SuperNode root = new SuperNode();

            for (final String loc : JSONObject.getNames(coverageZones)) {
                final JSONObject locData = coverageZones.getJSONObject(loc);
                final JSONObject coordinates = locData.optJSONObject("coordinates");
                Geolocation geolocation = null;

                if (coordinates != null && coordinates.has("latitude") && coordinates.has("longitude")) {
                    final double latitude = coordinates.optDouble("latitude");
                    final double longitude = coordinates.optDouble("longitude");
                    geolocation = new Geolocation(latitude, longitude);
                }

                try {
                    final JSONArray network6 = locData.getJSONArray("network6");

                    for (int i = 0; i < network6.length(); i++) {
                        final String ip = network6.getString(i);

                        try {
                            root.add6(new NetworkNode(ip, loc, geolocation));
                        } catch (NetworkNodeException ex) {
                            LOGGER.error(ex, ex);
                            return null;
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
                            root.add(new NetworkNode(ip, loc, geolocation));
                        } catch (NetworkNodeException ex) {
                            LOGGER.error(ex, ex);
                            return null;
                        }
                    }
                } catch (JSONException ex) {
                    LOGGER.warn("An exception was caught while accessing the network key of " + loc + " in the incoming coverage zone file: " + ex.getMessage());
                }
            }

            if (!verifyOnly) {
                instance = root;
            }

            return root;
        } catch (JSONException e) {
            LOGGER.warn(e,e);
        } catch (NetworkNodeException ex) {
            LOGGER.fatal(ex, ex);
        }

        return null;
    }

    public NetworkNode(final String str) throws NetworkNodeException {
        this(str, null);
    }

    public NetworkNode(final String str, final String loc) throws NetworkNodeException {
        this(str, loc, null);
    }

    public NetworkNode(final String str, final String loc, final Geolocation geolocation) throws NetworkNodeException {
        this.loc = loc;
        this.geolocation = geolocation;
        cidrAddress = CidrAddress.fromString(str);
    }

    public NetworkNode getNetwork(final String ip) throws NetworkNodeException {
        return getNetwork(new NetworkNode(ip));
    }

    public NetworkNode getNetwork(final NetworkNode ipnn) {
        if (this.compareTo(ipnn) != 0) {
            return null;
        }

        if (children == null) {
            return this;
        }

        final NetworkNode c = children.get(ipnn);

        if (c == null) {
            return this;
        }

        return c.getNetwork(ipnn);
    }

    public Boolean add(final NetworkNode nn) {
        synchronized(this) {
            if (children == null) {
                children = new TreeMap<NetworkNode,NetworkNode>();
            }

            return add(children, nn);
        }
    }

    protected Boolean add(final Map<NetworkNode,NetworkNode> children, final NetworkNode networkNode) {
        if (compareTo(networkNode) != 0) {
            return false;
        }

        for (final NetworkNode child : children.values()) {
            if (child.cidrAddress.equals(networkNode.cidrAddress)) {
                return false;
            }
        }

        final List<NetworkNode> movedChildren = new ArrayList<NetworkNode>();

        for (final NetworkNode child : children.values()) {
            if (networkNode.cidrAddress.includesAddress(child.cidrAddress)) {
                movedChildren.add(child);
                networkNode.add(child);
            }
        }

        for (final NetworkNode movedChild : movedChildren) {
            children.remove(movedChild);
        }

        for (final NetworkNode child : children.values()) {
            if (child.cidrAddress.includesAddress(networkNode.cidrAddress)) {
                return child.add(networkNode);
            }
        }

        children.put(networkNode, networkNode);
        return true;
    }

    public String getLoc() {
        return loc;
    }

    public Geolocation getGeolocation() {
        return geolocation;
    }

    public CacheLocation getCacheLocation() {
        return cacheLocation;
    }

    public void setCacheLocation(final CacheLocation cacheLocation) {
        this.cacheLocation = cacheLocation;
    }

    public int size() {
        if (children == null) {
            return 1;
        }

        int size = 1;

        for (final NetworkNode child : children.keySet()) {
            size += child.size();
        }

        return size;
    }

    public void clearCacheLocations() {
        synchronized(this) {
            cacheLocation = null;

            if (this instanceof SuperNode) {
                final SuperNode superNode = (SuperNode) this;

                if (superNode.children6 != null) {
                    for (final NetworkNode child : superNode.children6.keySet()) {
                        child.clearCacheLocations();
                    }
                }
            }

            if (children != null) {
                for (final NetworkNode child : children.keySet()) {
                    child.clearCacheLocations();
                }
            }
        }
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

        public NetworkNode getNetwork6(final NetworkNode networkNode) {
            if (children6 == null) {
                return this;
            }

            final NetworkNode c = children6.get(networkNode);

            if (c == null) {
                return this;
            }

            return c.getNetwork(networkNode);
        }
    }

    @Override
    public int compareTo(final NetworkNode other) {
        return cidrAddress.compareTo(other.cidrAddress);
    }

    public String toString() {
        String str = "";
        try {
            str = InetAddress.getByAddress(cidrAddress.getHostBytes()).toString().replace("/", "");
        } catch (UnknownHostException e) {
            LOGGER.warn(e,e);
        }

        return "[" + str + "/" + cidrAddress.getNetmaskLength() + "] - location:" + this.getLoc();
    }
}
