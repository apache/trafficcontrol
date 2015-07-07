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

package com.comcast.cdn.traffic_control.traffic_router.core.dns;

import java.io.File;
import java.io.FileWriter;
import java.io.IOException;
import java.net.Inet6Address;
import java.net.InetAddress;
import java.net.UnknownHostException;
import java.text.SimpleDateFormat;
import java.util.ArrayList;
import java.util.Date;
import java.util.HashMap;
import java.util.Iterator;
import java.util.List;
import java.util.Map;

import org.apache.commons.io.IOUtils;
import org.apache.log4j.Logger;
import org.json.JSONArray;
import org.json.JSONException;
import org.json.JSONObject;
import org.xbill.DNS.AAAARecord;
import org.xbill.DNS.ARecord;
import org.xbill.DNS.CNAMERecord;
import org.xbill.DNS.DClass;
import org.xbill.DNS.NSRecord;
import org.xbill.DNS.Name;
import org.xbill.DNS.RRset;
import org.xbill.DNS.Record;
import org.xbill.DNS.SOARecord;
import org.xbill.DNS.SetResponse;
import org.xbill.DNS.TextParseException;
import org.xbill.DNS.Type;
import org.xbill.DNS.Zone;

import com.comcast.cdn.traffic_control.traffic_router.core.cache.Cache;
import com.comcast.cdn.traffic_control.traffic_router.core.cache.CacheRegister;
import com.comcast.cdn.traffic_control.traffic_router.core.cache.InetRecord;
import com.comcast.cdn.traffic_control.traffic_router.core.cache.Resolver;
import com.comcast.cdn.traffic_control.traffic_router.core.cache.Cache.DeliveryServiceReference;
import com.comcast.cdn.traffic_control.traffic_router.core.ds.DeliveryService;
import com.comcast.cdn.traffic_control.traffic_router.core.request.DNSRequest;
import com.comcast.cdn.traffic_control.traffic_router.core.router.TrafficRouter;
import com.comcast.cdn.traffic_control.traffic_router.core.router.StatTracker;
import com.comcast.cdn.traffic_control.traffic_router.core.router.StatTracker.Track;
import com.comcast.cdn.traffic_control.traffic_router.core.util.Config;

public class ZoneManager extends Resolver {
	private static final Logger LOGGER = Logger.getLogger(ZoneManager.class);

	private final TrafficRouter trafficRouter;
	private final List<Zone> _zones;
	private static String dnsRoutingName;
	private static String httpRoutingName;

	private final StatTracker statTracker;

	public ZoneManager(final TrafficRouter tr, final StatTracker statTracker, final String drn, final String hrn) throws IOException {
		dnsRoutingName = drn.toLowerCase();
		httpRoutingName = hrn.toLowerCase();

		_zones = generateZones(tr.getCacheRegister());
		try {
			writeZones(_zones, Config.getVarDir());
		} catch(IOException e) {
			LOGGER.warn(e,e);
		}
		trafficRouter = tr;
		this.statTracker = statTracker;
	}

	private static void writeZones(final List<Zone> zones, final String confDir) throws IOException {
		synchronized(LOGGER) {
			final String zoneDirStr = confDir+"auto-zones";
			final File dir = new File(zoneDirStr);
			if(dir.exists()) {
				dir.delete();
			}
			final boolean result = dir.mkdirs();
			LOGGER.info(zoneDirStr+": "+result);
			for(Zone z : zones) {
				final String fileName = zoneDirStr+"/"+z.getOrigin().toString();
				LOGGER.warn("writing: "+fileName);
				//				File f = new File(fileName);
				//				result = f.createNewFile();
				final String file = z.toMasterFile();
				final FileWriter w = new FileWriter(fileName);
				IOUtils.write(file, w);
				w.flush();
				w.close();
			}
		}
	}

	public StatTracker getStatTracker() {
		return statTracker;
	}

	private static List<Zone> generateZones(final CacheRegister data) throws IOException {
		final Map<String, List<Record>> zoneMap = new HashMap<String, List<Record>>();
		final Map<String, DeliveryService> dsMap = new HashMap<String, DeliveryService>();
		final String tld = data.getConfig().optString("domain_name");
		for(DeliveryService ds : data.getDeliveryServices().values()) {
			final JSONArray domains = ds.getDomains();
			if(domains == null) {
				continue;
			}
			for (int j = 0; j < domains.length(); j++) {
				String domain = domains.optString(j);
				if(domain.endsWith("+")) {
					domain = domain.replaceAll("\\+\\z", ".") + tld;
				}
				dsMap.put(domain, ds);
			}
		}
		final Map<String, List<Record>> superDomains = populateZoneMap(zoneMap, dsMap, data);

		final List<Zone> zones = new ArrayList<Zone>();
		zones.addAll(fillZones(zoneMap, dsMap, data));
		// ensure that the superDomains go onto the end of the list
		zones.addAll(fillZones(superDomains, dsMap, data));
		return zones;
	}

	private static List<Zone> fillZones(final Map<String, List<Record>> zoneMap, final Map<String, DeliveryService> dsMap, final CacheRegister data) 
			throws IOException {
		final List<Zone> zones = new ArrayList<Zone>();
		final String hostname = InetAddress.getLocalHost().getHostName().replaceAll("\\..*", "");
		for(String domain : zoneMap.keySet()) {
			zones.add( createZone(domain, zoneMap, dsMap, hostname, data) );
		}
		return zones;
	}
	private static Zone createZone(final String domain, final Map<String, List<Record>> zoneMap, final Map<String, DeliveryService> dsMap, 
			final String hostname, final CacheRegister data) throws IOException {
		final DeliveryService ds = dsMap.get(domain);
		final JSONObject trafficRouters = data.getTrafficRouters();
		final JSONObject config = data.getConfig();

		JSONObject ttl = null;
		JSONObject soa = null;
		if(ds != null) {
			ttl = ds.getTtls();
			soa = ds.getSoa();
		} else {
			ttl = config.optJSONObject("ttls");
			soa = config.optJSONObject("soa");
		}

		final Name name = new Name(domain+".");
		final List<Record> list = zoneMap.get(domain);
		final Name meTr = newName(hostname, domain);
		final Name admin = newName(getString(soa, "admin", "twelve_monkeys"), domain);
		list.add(new SOARecord(name, DClass.IN, 
				getLong(ttl, "SOA", 86400), meTr, admin, 
				getLong(soa, "serial", getSerial()), 
				getLong(soa, "refresh", 28800), 
				getLong(soa, "retry", 7200), 
				getLong(soa, "expire", 604800), 
				getLong(soa, "minimum", 60)));

		addTrafficRouters(list, trafficRouters, name, ttl, domain, ds);
		addStaticDnsEntries(list, ds, domain);

		return new Zone(name, list.toArray(new Record[list.size()]));
	}

	private static void addStaticDnsEntries(final List<Record> list, final DeliveryService ds, final String domain)
			throws TextParseException, UnknownHostException {
		if (ds != null && ds.getStaticDnsEntries() != null) {
			final JSONArray entryList = ds.getStaticDnsEntries();

			for (int j = 0; j < entryList.length(); j++) {
				try {
					final JSONObject staticEntry = entryList.getJSONObject(j);
					final String type = staticEntry.getString("type").toUpperCase();
					final String jsName = staticEntry.getString("name");
					final String value = staticEntry.getString("value");
					final Name name = newName(jsName, domain);
					long ttl = staticEntry.optInt("ttl");

					if (ttl == 0) {
						ttl = getLong(ds.getTtls(), type, 60);
					}

					if ("A".equals(type)) {
						list.add(new ARecord(name, DClass.IN, ttl, InetAddress.getByName(value)));
					} else if ("AAAA".equals(type)) {
						list.add(new AAAARecord(name, DClass.IN, ttl, InetAddress.getByName(value)));
					} else if ("CNAME".equals(type)) {
						list.add(new CNAMERecord(name, DClass.IN, ttl, new Name(value)));
					}
				} catch (JSONException e) {
					LOGGER.error(e);
				}
			}
		}
	}

	private static void addTrafficRouters(final List<Record> list, final JSONObject trafficRouters, final Name name, 
			final JSONObject ttl, final String domain, final DeliveryService ds) 
					throws TextParseException, UnknownHostException {
		final boolean addTrafficRouters = (ds != null && ds.isDns()) ? false : true;
		final boolean useTrafficRouterStr = (ds != null) ? true : false;
		final boolean ip6RoutingEnabled = (ds == null || (ds != null && ds.isIp6RoutingEnabled())) ? true : false;

		for(String key : JSONObject.getNames(trafficRouters)) {
			final JSONObject trJo = trafficRouters.optJSONObject(key);
			if(trJo.has("status") && "OFFLINE".equals(trJo.optString("status"))) {
				// if "status": "OFFLINE"
				continue;
			}
			final Name trName = newName(key,domain);
			String ip6 = trJo.optString("ip6");

			list.add(new NSRecord(name, DClass.IN, getLong(ttl, "NS", 60), trName));
			list.add(new ARecord(trName,
					DClass.IN, getLong(ttl, "A", 60), 
					InetAddress.getByName(trJo.optString("ip"))));

			if (ip6 != null && !ip6.isEmpty() && ip6RoutingEnabled) {
				ip6 = ip6.replaceAll("/.*", "");
				list.add(new AAAARecord(trName,
						DClass.IN,
						getLong(ttl, "AAAA", 60),
						Inet6Address.getByName(ip6)));
			}

			if (!addTrafficRouters) {
				continue;
			}

			addTrafficRouterIps(list, domain, key, trJo, ttl, ip6RoutingEnabled, useTrafficRouterStr);
		}
	}

	private static void addTrafficRouterIps(final List<Record> list, final String domain, final String key,
			final JSONObject trJo, final JSONObject ttl, final boolean addTrafficRoutersAAAA, final boolean useTrafficRouterStr) 
					throws TextParseException, UnknownHostException {
		final Name trName = (useTrafficRouterStr)? newName(getHttpRoutingName(), domain):newName(key,domain);
		list.add(new ARecord(trName,
				DClass.IN,
				getLong(ttl, "A", 60),
				InetAddress.getByName(trJo.optString("ip"))));
		String ip6 = trJo.optString("ip6");
		if (addTrafficRoutersAAAA && ip6 != null && !ip6.isEmpty()) {
			ip6 = ip6.replaceAll("/.*", "");
			list.add(new AAAARecord(trName,
					DClass.IN,
					getLong(ttl, "AAAA", 60),
					Inet6Address.getByName(ip6)));
		}
	}

	private static Name newName(final String hostname, final String domain) throws TextParseException {
		return new Name(hostname+"."+domain+".");
	}

	private static final Map<String, List<Record>> populateZoneMap(final Map<String, List<Record>> zoneMap,
			final Map<String, DeliveryService> dsMap, final CacheRegister data) throws IOException {
		final Map<String, List<Record>> superDomains = new HashMap<String, List<Record>>();
		for(Cache c : data.getCacheMap().values()) {
			for(DeliveryServiceReference dsr : c.getDeliveryServices()) {
				final DeliveryService ds = data.getDeliveryService(dsr.getDeliveryServiceId());
				final String fqdn = dsr.getFqdn();
				final String[] parts = fqdn.split("\\.", 2);
				final String host = parts[0];
				final String domain = parts[1];
				dsMap.put(domain, ds);
				List<Record> zholder = zoneMap.get(domain);
				if(zholder == null) {
					zholder = new ArrayList<Record>();
					zoneMap.put(domain, zholder);
				}
				if (host.equalsIgnoreCase(getDnsRoutingName())) {
					continue;
				}
				final Name name = new Name(fqdn+".");
				final JSONObject ttl = ds.getTtls();
				try {
				zholder.add(new ARecord(name,
						DClass.IN,
						getLong(ttl, "A", 60),
						c.getIp4()));
				} catch(java.lang.IllegalArgumentException e) {
					LOGGER.warn(e+" : "+c.getIp4());
				}
				final InetAddress ip6 = c.getIp6();
				if(ip6 != null && ds != null && ds.isIp6RoutingEnabled()) {
					zholder.add(new AAAARecord(name,
							DClass.IN,
							getLong(ttl, "AAAA", 60),
							ip6));
				}
				final String superdomain = domain.split("\\.", 2)[1];
				zholder = superDomains.get(superdomain);
				if(zholder == null) {
					zholder = new ArrayList<Record>();
					superDomains.put(superdomain, zholder);
				}
			}
		}
		return superDomains;
	}

	private static final SimpleDateFormat sdf = new SimpleDateFormat("yyyyMMddHH");
	private static long getSerial() {
		synchronized(sdf) {
			return Long.parseLong(sdf.format(new Date())); // 2013062701
		}
	}

	private static String getString(final JSONObject jo, final String key, final String d) {
		if(jo == null) { return d; }
		if(!jo.has(key)) { return d; }
		return jo.optString(key);
	}
	private static long getLong(final JSONObject jo, final String key, final long d) {
		if(jo == null) { return d; }
		if(!jo.has(key)) { return d; }
		return jo.optLong(key);
	}

	/**
	 * Gets trafficRouter.
	 * 
	 * @return the trafficRouter
	 */
	public TrafficRouter getTrafficRouter() {
		return trafficRouter;
	}

	/**
	 * Attempts to find a {@link Zone} that would contain the specified {@link Name}.
	 * 
	 * @param name
	 *            the Name to use to attempt to find the Zone
	 * @param clientAddress
	 *            the IP address of the client
	 * @return the Zone to use to resolve the specified Name
	 * @throws DNSException
	 *             if no Zone is found that supports the specified Name
	 */
	public Zone getZone(final Name name) {
		Zone result = null;

		//		zonesLock.readLock().lock();
		for (final Zone zone : _zones) {
			final Name origin = zone.getOrigin();
			if (name.subdomain(origin)) {
				result = zone;
				break;
			}
		}
		//		zonesLock.readLock().unlock();

		if (result == null) {
			LOGGER.warn(String.format("subdomain NOT found: '%s'", name));
			//			throw new DNSException(Rcode.REFUSED);
		}
		return result;
	}

	/**
	 * Gets zoneDirectory.
	 * 
	 * @return the zoneDirectory
	 */
	//	public File getZoneDirectory() {
	//		return zoneDirectory;
	//	}


	/**
	 * Creates a dynamic zone that serves a set of A and AAAA records for the specified {@link Name}
	 * .
	 * 
	 * @param staticZone
	 *            The Zone that would normally serve this request
	 * @param name
	 *            the Name that is being requested
	 * @param clientAddress
	 *            the IP address of the requestor
	 * @return the new Zone to serve the request or null if the static Zone should be used
	 */
	private Zone createDynamicZone(final Zone staticZone, final Name name, final int qtype, final InetAddress clientAddress) {
		if (clientAddress == null) {
			return staticZone;
		}

		final DNSRequest request = new DNSRequest();
		request.setClientIP(clientAddress.getHostAddress());
		request.setHostname(name.relativize(Name.root).toString());
		request.setQtype(qtype);
		final Track track = StatTracker.getTrack();
		List<InetRecord> addresses = null;
		try {
			addresses = trafficRouter.route(request, track);
			return fillZone(staticZone, name, addresses);
		} catch (final Exception e) {
			LOGGER.error(e.getMessage(), e);
		} finally {
			statTracker.saveTrack(track);
		}
		return null;
	}

	private static Zone fillZone(final Zone staticZone, final Name name, final List<InetRecord> addresses) {
		if(addresses == null) { return null; }
		try {
			final Zone result = createZone(staticZone);
			int recordsAdded = 0;
			for (final InetRecord address : addresses) {
				final Record record = createRecord(name, address);
				if (record != null) {
					result.addRecord(record);
					recordsAdded++;
				}
			}
			if (recordsAdded > 0) {
				return result;
			}
		} catch (final IOException e) {
			LOGGER.error(e.getMessage(), e);
		}
		return null;
	}

	private static Record createRecord(final Name name, final InetRecord address) throws TextParseException {
		Record record = null;
		if(address.isAlias()) {
			record = new CNAMERecord(name, DClass.IN, address.getTTL(), new Name(address.getAlias() +"."));			
		} else if (address.isInet4()) { // address instanceof Inet4Address
			record = new ARecord(name, DClass.IN, address.getTTL(), address.getAddress());
		} else if (address.isInet6()) {
			record = new AAAARecord(name, DClass.IN, address.getTTL(), address.getAddress());
		}
		return record;
	}

	private static Zone createZone(final Zone staticZone) throws IOException {
		final List<Record> records = new ArrayList<Record>();
		records.add(staticZone.getSOA());
		@SuppressWarnings("unchecked")
		final Iterator<Record> ns = staticZone.getNS().rrs();
		while (ns.hasNext()) {
			records.add(ns.next());
		}
		return new Zone(staticZone.getOrigin(), records.toArray(new Record[records.size()]));
	}


	private List<InetRecord> lookup(final Name qname, final Zone zone, final int type) {
		final List<InetRecord> ipAddresses = new ArrayList<InetRecord>();
		final SetResponse sr = zone.findRecords(qname, type);
		if (LOGGER.isDebugEnabled()) {
			LOGGER.debug("SetResponse: " + sr);
		}
		if (sr.isSuccessful()) {
			final RRset[] answers = sr.answers();
			for (final RRset answer : answers) {
				LOGGER.debug("addRRset: " + answer);
				@SuppressWarnings("unchecked")
				final Iterator<Record> it = answer.rrs();
				while (it.hasNext()) {
					final Record r = it.next();
					if(r instanceof ARecord) {
						final ARecord ar = (ARecord)r;
						ipAddresses.add(new InetRecord(ar.getAddress(), ar.getTTL()));
					} else if(r instanceof AAAARecord) {
						final AAAARecord ar = (AAAARecord)r;
						ipAddresses.add(new InetRecord(ar.getAddress(), ar.getTTL()));
					} else {
						LOGGER.debug("record not ARecord or AAAARecord: " + r);
					}
				}
			}
			return ipAddresses;
		} else {
			LOGGER.debug(String.format("failed: zone.findRecords(%s, %d)", qname, type));
		}
		return null;
	}
	public List<InetRecord> resolve(final String fqdn) {
		try {
			final Name name = new Name(fqdn);
			final Zone zone = getZone(name);
			if(zone == null) {
				LOGGER.debug("No zone - Defaulting to system resolver: "+fqdn);
				return super.resolve(fqdn);
			}
			return lookup(name, zone, Type.A);
		} catch (TextParseException e) {
			LOGGER.warn("TextParseException from: "+fqdn,e);
		}

		LOGGER.debug(String.format("resolved from zone: %s", fqdn));
		return null;
	}
	public List<InetRecord> resolve(final String fqdn, final String address) throws UnknownHostException {
		try {
			final Name name = new Name(fqdn);
			Zone zone = getZone(name);
			final InetAddress addr = InetAddress.getByName(address);
			final int qtype = (addr instanceof Inet6Address) ? Type.AAAA : Type.A;
			final Zone dynamicZone = createDynamicZone(zone, name, qtype, addr);
			if (dynamicZone != null) { 
				zone = dynamicZone; 
			}
			if(zone == null) {
				LOGGER.debug("No zone - Defaulting to system resolver: "+fqdn);
				return super.resolve(fqdn);
			}
			return lookup(name, zone, Type.A);
		} catch (TextParseException e) {
			LOGGER.warn("TextParseException from: "+fqdn,e);
		}

		LOGGER.debug(String.format("resolved from zone: %s", fqdn));
		return null;
	}

	public Zone getDynamicZone(final Name qname, final int qtype, final InetAddress clientAddress) {
		final Zone zone = getZone(qname);
		if (zone == null) {
			return null;
		}

		final SetResponse sr = zone.findRecords(qname, qtype);

		if (sr.isSuccessful()) {
			return zone;
		} else if (qname.toString().toLowerCase().matches(getDnsRoutingName() + "\\..*")) {
			final Zone dynamicZone = createDynamicZone(zone, qname, qtype, clientAddress);

			if (dynamicZone != null) {
				return dynamicZone;
			}
		}

		return zone;
	}

	private static String getDnsRoutingName() {
		return dnsRoutingName;
	}

	private static String getHttpRoutingName() {
		return httpRoutingName;
	}
}
