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

package com.comcast.cdn.traffic_control.traffic_router.core.router;


import com.comcast.cdn.traffic_control.traffic_router.core.cache.Cache;
import com.comcast.cdn.traffic_control.traffic_router.core.cache.CacheLocation;
import com.comcast.cdn.traffic_control.traffic_router.core.cache.CacheRegister;
import com.comcast.cdn.traffic_control.traffic_router.core.config.ParseException;
import com.comcast.cdn.traffic_control.traffic_router.core.ds.DeliveryService;
import com.comcast.cdn.traffic_control.traffic_router.core.ds.DeliveryServiceMatcher;
import org.json.JSONArray;
import org.json.JSONException;
import org.json.JSONObject;

import java.net.UnknownHostException;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.TreeSet;

public class CacheRegisterBuilder {

    public static void parseCacheConfig(final JSONObject contentServers, final CacheRegister cacheRegister) throws JSONException, ParseException {
        final Map<String,Cache> map = new HashMap<String,Cache>();
        final Map<String, List<String>> statMap = new HashMap<String, List<String>>();
        for (final String node : JSONObject.getNames(contentServers)) {
            final JSONObject jo = contentServers.getJSONObject(node);
            final CacheLocation loc = cacheRegister.getCacheLocation(jo.getString("locationId"));
            if (loc != null) {
                String hashId = node;
                if(jo.has("hashId")) {
                    hashId = jo.optString("hashId");
                }
                final Cache cache = new Cache(node, hashId, jo.optInt("hashCount"));
                cache.setFqdn(jo.getString("fqdn"));
                cache.setPort(jo.getInt("port"));
                final String ip = jo.getString("ip");
                final String ip6 = jo.optString("ip6");
                try {
                    cache.setIpAddress(ip, ip6, 0);
                } catch (UnknownHostException e) {
                    System.out.println(e + ": " + ip);
                }

                if(jo.has("deliveryServices")) {
                    final List<Cache.DeliveryServiceReference> references = new ArrayList<Cache.DeliveryServiceReference>();
                    final JSONObject dsJos = jo.optJSONObject("deliveryServices");
                    for(String ds : JSONObject.getNames(dsJos)) {
                        final Object dso = dsJos.get(ds);

                        List<String> dsNames = statMap.get(ds);

                        if (dsNames == null) {
                            dsNames = new ArrayList<String>();
                        }

                        if (dso instanceof JSONArray) {
                            final JSONArray fqdnList = (JSONArray) dso;

                            if (fqdnList != null && fqdnList.length() > 0) {
                                for (int i = 0; i < fqdnList.length(); i++) {
                                    final String name = fqdnList.getString(i).toLowerCase();

                                    if (i == 0) {
                                        references.add(new Cache.DeliveryServiceReference(ds, name));
                                    }

                                    final String tld = cacheRegister.getConfig().optString("domain_name").toLowerCase();

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
                                }
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

    public static void parseDeliveryServiceConfig(final JSONObject deliveryServices, final CacheRegister cacheRegister) throws JSONException {
        final TreeSet<DeliveryServiceMatcher> dnsServiceMatchers = new TreeSet<DeliveryServiceMatcher>();
        final TreeSet<DeliveryServiceMatcher> httpServiceMatchers = new TreeSet<DeliveryServiceMatcher>();
        final Map<String,DeliveryService> dsMap = new HashMap<String,DeliveryService>();

        for(String dsId : JSONObject.getNames(deliveryServices)) {
            final JSONObject dsJo = deliveryServices.getJSONObject(dsId);
            final JSONArray matchsets = dsJo.getJSONArray("matchsets");
            final DeliveryService ds = new DeliveryService(dsId, dsJo);
            boolean isDns = false;

            dsMap.put(dsId, ds);

            for (int i = 0; i < matchsets.length(); i++) {
                final JSONObject matchset = matchsets.getJSONObject(i);
                final String protocol = matchset.getString("protocol");

                if ("DNS".equals(protocol)) {
                    isDns = true;
                }

                final DeliveryServiceMatcher m = new DeliveryServiceMatcher(ds);

                if ("HTTP".equals(protocol)) {
                    httpServiceMatchers.add(m);
                } else if("DNS".equals(protocol)) {
                    dnsServiceMatchers.add(m);
                }

                final JSONArray list = matchset.getJSONArray("matchlist");
                for (int j = 0; j < list.length(); j++) {
                    final JSONObject matcherJo = list.getJSONObject(j);
                    final DeliveryServiceMatcher.Type type = DeliveryServiceMatcher.Type.valueOf(matcherJo.getString("match-type"));
                    final String target = matcherJo.optString("target");
                    m.addMatch(type, matcherJo.getString("regex"), target);
                }
            }

            ds.setDns(isDns);
        }

        cacheRegister.setDeliveryServiceMap(dsMap);
        cacheRegister.setDnsDeliveryServiceMatchers(dnsServiceMatchers);
        cacheRegister.setHttpDeliveryServiceMatchers(httpServiceMatchers);
    }

}
