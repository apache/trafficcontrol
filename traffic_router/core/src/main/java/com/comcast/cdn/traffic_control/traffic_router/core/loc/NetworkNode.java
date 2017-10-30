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
import java.io.IOException;
import java.net.InetAddress;
import java.net.UnknownHostException;
import java.util.Iterator;
import java.util.Map;
import java.util.TreeMap;
import java.util.List;
import java.util.ArrayList;
import java.util.concurrent.CountDownLatch;

import com.comcast.cdn.traffic_control.traffic_router.core.util.CidrAddress;
import com.comcast.cdn.traffic_control.traffic_router.core.util.JsonUtils;
import com.comcast.cdn.traffic_control.traffic_router.core.util.JsonUtilsException;
import com.comcast.cdn.traffic_control.traffic_router.geolocation.Geolocation;

import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.apache.log4j.Logger;

import com.comcast.cdn.traffic_control.traffic_router.core.cache.Cache;
import com.comcast.cdn.traffic_control.traffic_router.core.cache.CacheLocation;
import com.comcast.cdn.traffic_control.traffic_router.core.cache.CacheRegister;

public class NetworkNode implements Comparable<NetworkNode> {
    private static final Logger LOGGER = Logger.getLogger(NetworkNode.class);
    private static final String DEFAULT_SUB_STR = "0.0.0.0/0";

    private static NetworkNode instance;
    private static NetworkNode deepInstance;

    private static CacheRegister cacheRegister;
    private static final CountDownLatch cacheRegisterLatch = new CountDownLatch(1);

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

    public static NetworkNode getDeepInstance() {
        if (deepInstance != null) {
            return deepInstance;
        }

        try {
            deepInstance = new NetworkNode(DEFAULT_SUB_STR);
        } catch (NetworkNodeException e) {
            LOGGER.warn(e);
        }

        return deepInstance;
    }

    public static void setCacheRegister(final CacheRegister cr) {
        cacheRegister = cr;
        cacheRegisterLatch.countDown();
    }

    public static CacheRegister getCacheRegisterBlocking() {
        try {
            cacheRegisterLatch.await();
        } catch (InterruptedException e) {
            LOGGER.warn(e);
        } finally {
            return cacheRegister;
        }
    }

    public static NetworkNode generateTree(final File f, final boolean verifyOnly, final boolean useDeep) throws IOException  {
        final ObjectMapper mapper = new ObjectMapper();
        return generateTree(mapper.readTree(f), verifyOnly, useDeep);
    }

    public static NetworkNode generateTree(final File f, final boolean verifyOnly) throws IOException  {
        return generateTree(f, verifyOnly, false);
    }

    public static NetworkNode generateTree(final JsonNode json, final boolean verifyOnly) {
        return generateTree(json, verifyOnly, false);
    }

    @SuppressWarnings({"PMD.CyclomaticComplexity", "PMD.NPathComplexity"})
    public static NetworkNode generateTree(final JsonNode json, final boolean verifyOnly, final boolean useDeep) {
        try {
            final JsonNode coverageZones = JsonUtils.getJsonNode(json, "coverageZones");

            final SuperNode root = new SuperNode();

            final Iterator<String> czIter = coverageZones.fieldNames();
            while (czIter.hasNext()) {
                final String loc = czIter.next();
                final JsonNode locData = JsonUtils.getJsonNode(coverageZones, loc);
                final JsonNode coordinates = locData.get("coordinates");
                Geolocation geolocation = null;

                if (coordinates != null && coordinates.has("latitude") && coordinates.has("longitude")) {
                    final double latitude = coordinates.get("latitude").asDouble();
                    final double longitude = coordinates.get("longitude").asDouble();
                    geolocation = new Geolocation(latitude, longitude);
                }
                CacheLocation deepLoc = null;
                if (useDeep) {
                    try {
                        final JsonNode caches = JsonUtils.getJsonNode(locData, "caches");
                        for (final JsonNode cacheJson : caches) {
                            final String cacheHostname = cacheJson.asText();
                            if (deepLoc == null) {
                                deepLoc = new CacheLocation( "deep." + loc, new Geolocation(0.0, 0.0));  // TODO JvD
                            }
                            // Get the cache from the cacheregister here - don't create a new cache due to the deep file, only reuse the
                            // ones we already know about.
                            final Cache cache = getCacheRegisterBlocking().getCacheMap().get(cacheHostname);
                            if (cache == null) {
                                LOGGER.warn("DDC: deep cache entry " + cacheHostname + " not found in crconfig server list (it might not belong to this CDN)");
                            } else {
                                LOGGER.info("DDC: Adding " + cacheHostname + " to " + deepLoc.getId() + ".");
                                deepLoc.addCache(cache);
                            }
                        }
                    } catch (JsonUtilsException ex) {
                        LOGGER.warn("An exception was caught while accessing the caches key of " + loc + " in the incoming coverage zone file: " + ex.getMessage());
                    }
                }

                if (!addNetworkNodesToRoot(root, locData, loc, deepLoc, geolocation, useDeep)) {
                    return null;
                }
            }

            if (!verifyOnly) {
                if (useDeep) {
                    deepInstance = root;
                } else {
                    instance = root;
                }
            }

            return root;
        } catch (JsonUtilsException ex) {
            LOGGER.warn(ex, ex);
        } catch (NetworkNodeException ex) {
            LOGGER.fatal(ex, ex);
        }

        return null;
    }

    private static boolean addNetworkNodesToRoot(final SuperNode root, final JsonNode locData, final String loc,
                                                 final CacheLocation deepLoc, final Geolocation geolocation, final boolean useDeep) {
        for (final String key : new String[]{"network6", "network"}) {
            try {
                for (final JsonNode network : JsonUtils.getJsonNode(locData, key)) {
                    final String ip = network.asText();

                    try {
                        final NetworkNode nn = new NetworkNode(ip, loc, geolocation);
                        if (useDeep && deepLoc != null) { // for deepLoc, we add the location here; normally it gets added by setLocation.
                            nn.setCacheLocation(deepLoc);
                        }
                        if ("network6".equals(key)) {
                            root.add6(nn);
                        } else {
                            root.add(nn);
                        }
                    } catch (NetworkNodeException ex) {
                        LOGGER.error(ex, ex);
                        return false;
                    }
                }
            } catch (JsonUtilsException ex) {
                LOGGER.warn("An exception was caught while accessing the " + key + " key of " + loc + " in the incoming coverage zone file: " + ex.getMessage());
            }
        }
        return true;
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
