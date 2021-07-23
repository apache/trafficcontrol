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

package org.apache.traffic_control.traffic_router.core.router;


import org.apache.traffic_control.traffic_router.core.edge.Cache;
import org.apache.traffic_control.traffic_router.core.edge.CacheLocation;
import org.apache.traffic_control.traffic_router.core.edge.CacheRegister;
import org.apache.traffic_control.traffic_router.core.config.ParseException;
import org.apache.traffic_control.traffic_router.core.ds.DeliveryService;
import org.apache.traffic_control.traffic_router.core.ds.DeliveryServiceMatcher;
import org.apache.traffic_control.traffic_router.core.util.JsonUtils;
import org.apache.traffic_control.traffic_router.core.util.JsonUtilsException;
import com.fasterxml.jackson.databind.JsonNode;

import java.net.UnknownHostException;
import java.util.HashMap;
import java.util.Iterator;
import java.util.List;
import java.util.Map;
import java.util.ArrayList;
import java.util.TreeSet;

public class CacheRegisterBuilder {

    public static void parseCacheConfig(final JsonNode contentServers, final CacheRegister cacheRegister) throws JsonUtilsException, ParseException {
        final Map<String,Cache> map = new HashMap<String,Cache>();
        final Map<String, List<String>> statMap = new HashMap<String, List<String>>();

        final Iterator<String> contentServersIter = contentServers.fieldNames();
        while (contentServersIter.hasNext()) {
            final String node = contentServersIter.next();
            final JsonNode jo = JsonUtils.getJsonNode(contentServers, node);
            final CacheLocation loc = cacheRegister.getCacheLocation(JsonUtils.getString(jo, "locationId"));
            if (loc != null) {
                String hashId = node;
                if(jo.has("hashId")) {
                    hashId = jo.get("hashId").asText();
                }
                final int hashCount = JsonUtils.optInt(jo, "hashCount");
                final Cache cache = new Cache(node, hashId, hashCount);
                cache.setFqdn(JsonUtils.getString(jo, "fqdn"));
                cache.setPort(JsonUtils.getInt(jo, "port"));
                final String ip = JsonUtils.getString(jo, "ip");
                final String ip6 = JsonUtils.optString(jo, "ip6");
                try {
                    cache.setIpAddress(ip, ip6, 0);
                } catch (UnknownHostException e) {
                    System.out.println(e + ": " + ip);
                }

                if (jo.has("deliveryServices")) {
                    final List<Cache.DeliveryServiceReference> references = new ArrayList<Cache.DeliveryServiceReference>();
                    final JsonNode dsJos = jo.get("deliveryServices");
                    final Iterator<String> dsIter = dsJos.fieldNames();
                    while (dsIter.hasNext()) {
                        final String ds = dsIter.next();
                        final JsonNode dso = dsJos.get(ds);

                        List<String> dsNames = statMap.get(ds);

                        if (dsNames == null) {
                            dsNames = new ArrayList<String>();
                        }

                        if (dso.isArray()) {
                            int i = 0;
                            for (JsonNode fqdn : dso) {
                                final String name = fqdn.asText().toLowerCase();

                                if (i == 0) {
                                    references.add(new Cache.DeliveryServiceReference(ds, name));
                                }

                                final String tld = JsonUtils.optString(cacheRegister.getConfig(), "domain_name").toLowerCase();

                                if (name.contains(tld)) {
                                    final String reName = name.replaceAll("^.*?\\.", "");

                                    if (!dsNames.contains(reName)) {
                                        dsNames.add(reName);
                                    }
                                } else {
                                    if (!dsNames.contains(name)) {
                                        dsNames.add(name);
                                    }
                                }
                                i++;
                            }
                        } else {
                            references.add(new Cache.DeliveryServiceReference(ds, dso.toString()));

                            if (!dsNames.contains(dso.toString())) {
                                dsNames.add(dso.toString());
                            }
                        }
                        statMap.put(ds, dsNames);
                    }
                    cache.setDeliveryServices(references);
                }
                loc.addCache(cache);
                map.put(cache.getId(), cache);
            }
        }
        cacheRegister.setCacheMap(map);
    }

    public static void parseDeliveryServiceConfig(final JsonNode deliveryServices, final CacheRegister cacheRegister) throws JsonUtilsException {
        final TreeSet<DeliveryServiceMatcher> deliveryServiceMatchers = new TreeSet<DeliveryServiceMatcher>();
        final Map<String,DeliveryService> dsMap = new HashMap<String,DeliveryService>();

        final Iterator<String> keyIter = deliveryServices.fieldNames();
        while (keyIter.hasNext()) {
            final String dsId = keyIter.next();
            final JsonNode dsJo = JsonUtils.getJsonNode(deliveryServices, dsId);
            final JsonNode matchsets = JsonUtils.getJsonNode(dsJo, "machsets");
            final DeliveryService ds = new DeliveryService(dsId, dsJo);
            boolean isDns = false;

            dsMap.put(dsId, ds);

            for (final JsonNode matchset : matchsets) {
                final String protocol = JsonUtils.getString(matchset, "protocol");

                final DeliveryServiceMatcher m = new DeliveryServiceMatcher(ds);
                deliveryServiceMatchers.add(m);

                for (JsonNode matchlist : matchset.get("matchlist")) {
                    final DeliveryServiceMatcher.Type type = DeliveryServiceMatcher.Type.valueOf(JsonUtils.getString(matchlist, "match-type"));
                    final String target = JsonUtils.optString(matchlist, "target");
                    m.addMatch(type, JsonUtils.getString(matchlist, "regex"), target);
                }
            }
            ds.setDns(isDns);
        }

        cacheRegister.setDeliveryServiceMap(dsMap);
        cacheRegister.setDeliveryServiceMatchers(deliveryServiceMatchers);
    }

}
