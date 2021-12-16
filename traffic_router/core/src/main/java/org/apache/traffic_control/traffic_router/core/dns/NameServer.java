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

package org.apache.traffic_control.traffic_router.core.dns;

import java.net.InetAddress;
import java.util.*;

import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.xbill.DNS.DClass;
import org.xbill.DNS.ExtendedFlags;
import org.xbill.DNS.Flags;
import org.xbill.DNS.Message;
import org.xbill.DNS.Name;
import org.xbill.DNS.OPTRecord;
import org.xbill.DNS.RRset;
import org.xbill.DNS.Rcode;
import org.xbill.DNS.Record;
import org.xbill.DNS.SOARecord;
import org.xbill.DNS.Section;
import org.xbill.DNS.SetResponse;
import org.xbill.DNS.Type;
import org.xbill.DNS.Zone;
import org.xbill.DNS.EDNSOption;
import org.xbill.DNS.ClientSubnetOption;

import org.apache.traffic_control.traffic_router.core.ds.DeliveryService;
import org.apache.traffic_control.traffic_router.core.router.TrafficRouterManager;

@SuppressWarnings("PMD.CyclomaticComplexity")
public class NameServer {
	private static final int MAX_SUPPORTED_EDNS_VERS = 0;
	private static final int MAX_ITERATIONS = 6;
	private static final int NUM_SECTIONS = 4;
	private static final int FLAG_DNSSECOK = 1;
	private static final int FLAG_SIGONLY = 2;

	private static final Logger LOGGER = LogManager.getLogger(NameServer.class);
	private boolean ecsEnable = false;
	private Set<DeliveryService> ecsEnabledDses = new HashSet<>();
	/**
	 * 
	 */
	private TrafficRouterManager trafficRouterManager;

	/**
	 * Queries the zones based on the request and returns the appropriate response.
	 * 
	 * @param request
	 *            the query message
	 * @param clientAddress
	 *            the IP address of the client
	 * @return a response message
	 */
	public Message query(final Message request, final InetAddress clientAddress, final DNSAccessRecord.Builder builder) {
		final Message response = new Message();
		try {
			addQuestion(request, response);
			addAnswers(request, response, clientAddress, builder);
		} catch (final RuntimeException e) {
			LOGGER.error(e.getMessage(), e);
			response.getHeader().setRcode(Rcode.SERVFAIL);
		}

		return response;
	}

	@SuppressWarnings({"PMD.CyclomaticComplexity", "PMD.NPathComplexity", "PMD.AvoidDeeplyNestedIfStmts"})
	private void addAnswers(final Message request, final Message response, final InetAddress clientAddress, final DNSAccessRecord.Builder builder) {
		final Record question = request.getQuestion();

		if (question != null) {
			final int qclass = question.getDClass();
			final Name qname = question.getName();
			final OPTRecord qopt = request.getOPT();
			List<EDNSOption> list = Collections.EMPTY_LIST;
			boolean dnssecRequest = false;
			int qtype = question.getType();
			int flags = 0;

			if ((qopt != null) && (qopt.getVersion() > MAX_SUPPORTED_EDNS_VERS)) {
				response.getHeader().setRcode(Rcode.NOTIMP);
				final OPTRecord opt = new OPTRecord(0, Rcode.BADVERS, MAX_SUPPORTED_EDNS_VERS);
				response.addRecord(opt, Section.ADDITIONAL);
				return;
			}

			if ((qclass != DClass.IN) && (qclass != DClass.ANY)) {
				response.getHeader().setRcode(Rcode.REFUSED);
				return;
			}

			if (qopt != null && (qopt.getFlags() & ExtendedFlags.DO) != 0) {
				flags = FLAG_DNSSECOK;
				dnssecRequest = true;
			}

			if (qtype == Type.SIG || qtype == Type.RRSIG) {
				qtype = Type.ANY;
				flags |= FLAG_SIGONLY;
			}
			// Get list of options matching client subnet option code (8)
			if (qopt != null ){
				list = qopt.getOptions(EDNSOption.Code.CLIENT_SUBNET);
			}
			InetAddress ipaddr = null;
			int nmask = 0;
			if (isEcsEnable(qname)) {
				for (final EDNSOption option : list) {
					assert (option instanceof ClientSubnetOption);
					// If there are multiple ClientSubnetOptions in the Option RR, then
					// choose the one with longest source prefix. RFC 7871
					if (((ClientSubnetOption)option).getSourceNetmask() > nmask) {
						nmask = ((ClientSubnetOption)option).getSourceNetmask();
						ipaddr = ((ClientSubnetOption)option).getAddress();
					}
				}
			}
			if ((ipaddr!= null) && (isEcsEnable(qname))) {
				builder.client(ipaddr);

				LOGGER.debug("DNS: Using Client IP Address from ECS Option" + ipaddr.getHostAddress() + "/" 
						+ nmask);
				lookup(qname, qtype, ipaddr, response, flags, dnssecRequest, builder);
			} else {
				lookup(qname, qtype, clientAddress, response, flags, dnssecRequest, builder);
			}
			
			if (response.getHeader().getRcode() == Rcode.REFUSED) {
				return;
			}

			// Check if we had incoming ClientSubnetOption in Option RR, then we need
			// to return with the response, setting the scope subnet as well
			if ((nmask != 0) && (isEcsEnable(qname))) {
				final ClientSubnetOption cso = new ClientSubnetOption(nmask, nmask, ipaddr);
				final List<ClientSubnetOption> csoList = new ArrayList<ClientSubnetOption>(1);
				csoList.add(cso);	
				// OptRecord Arguments: payloadSize = 1280, xrcode = 0, version=0, flags=0, option List
				final OPTRecord opt = new OPTRecord(1280, 0, 0, 0, csoList);
				response.addRecord(opt, Section.ADDITIONAL);
			}
		
			if (qopt != null && flags == FLAG_DNSSECOK) {
				final int optflags = ExtendedFlags.DO;
				final OPTRecord opt = new OPTRecord(1280, (byte) 0, (byte) 0, optflags);
				response.addRecord(opt, Section.ADDITIONAL);
			}
		}
	}

	@SuppressWarnings({"PMD.UseStringBufferForStringAppends"})
	private boolean isDeliveryServiceEcsEnabled(final Name name) {
		boolean isEnabled = false;

		for (final DeliveryService ds : ecsEnabledDses) {
			String domain = ds.getDomain();

			if (domain == null) {
				continue;
			}

			if (domain.endsWith("+")) {
				domain = domain.replaceAll("\\+\\z", ".") + ZoneManager.getTopLevelDomain();
			}

			if (name.relativize(Name.root).toString().contains(domain)) {
				isEnabled = true;
				break;
			}
		}

		return isEnabled;
	}

	private static void addAuthority(final Zone zone, final Message response, final int flags) {
		final RRset authority = zone.getNS();
		addRRset(authority.getName(), response, authority, Section.AUTHORITY, flags);
		response.getHeader().setFlag(Flags.AA);
	}

	private static void addSOA(final Zone zone, final Message response, final int section, final int flags) {
		// we locate the SOA this way so that we can ensure we get the RRSIGs rather than just the one SOA Record
		final SetResponse fsoa = zone.findRecords(zone.getOrigin(), Type.SOA);

		if (!fsoa.isSuccessful()) {
			return;
		}

		for (final RRset answer : fsoa.answers()) {
			addRRset(zone.getOrigin(), response, setNegativeTTL(answer, flags), section, flags);
		}
	}

	@SuppressWarnings({"unchecked", "PMD.CyclomaticComplexity", "PMD.NPathComplexity"})
	private static void addDenialOfExistence(final Name qname, final Zone zone, final Message response, final int flags) {
		// The requirements for this are described in RFC 7129
		if ((flags & (FLAG_SIGONLY | FLAG_DNSSECOK)) == 0) {
			return;
		}

		RRset nsecSpan = null;
		Name candidate = null;

		final Iterator<RRset> zi = zone.iterator();

		while (zi.hasNext()) {
			final RRset rrset = zi.next();

			if (rrset.getType() != Type.NSEC) {
				continue;
			}

			final Iterator<Record> it = rrset.rrs();

			while (it.hasNext()) {
				final Record r = it.next();
				final Name name = r.getName();

				if (name.compareTo(qname) < 0 || (candidate != null && name.compareTo(candidate) < 0)) {
					candidate = name;
					nsecSpan = rrset;
				} else if (name.compareTo(qname) > 0 && candidate != null) {
					break;
				}
			}
		}

		if (candidate != null && nsecSpan != null) {
			addRRset(candidate, response, nsecSpan, Section.AUTHORITY, flags);
		}

		final SetResponse nxsr = zone.findRecords(zone.getOrigin(), Type.NSEC);
		if (nxsr.isSuccessful()) {
			for (final RRset answer : nxsr.answers()) {
				addRRset(qname, response, answer, Section.AUTHORITY, flags);
			}
		}
	}

	private static void addQuestion(final Message request, final Message response) {
		response.getHeader().setID(request.getHeader().getID());
		response.getHeader().setFlag(Flags.QR);
		if (request.getHeader().getFlag(Flags.RD)) {
			response.getHeader().setFlag(Flags.RD);
		}
		response.addRecord(request.getQuestion(), Section.QUESTION);
	}

	@SuppressWarnings({"unchecked", "PMD.CyclomaticComplexity"})
	private static void addRRset(final Name name, final Message response, final RRset rrset, final int section, final int flags) {
		for (int s = 1; s < NUM_SECTIONS; s++) {
			if (response.findRRset(name, rrset.getType(), s)) {
				return;
			}
		}

		final List<Record> recordList = new ArrayList<Record>();

		if ((flags & FLAG_SIGONLY) == 0) {
			final Iterator<Record> it = rrset.rrs();
			while (it.hasNext()) {
				Record r = it.next();
				if (r.getName().isWild() && !name.isWild()) {
					r = r.withName(name);
				}
				recordList.add(r);
			}
		}

		// We prefer to shuffle the list over "cycling" as we could with rrset.rrs(true) above.
		Collections.shuffle(recordList);

		for (final Record r : recordList) {
			response.addRecord(r, section);
		}

		if ((flags & (FLAG_SIGONLY | FLAG_DNSSECOK)) != 0) {
			final Iterator<Record> it = rrset.sigs();
			while (it.hasNext()) {
				Record r = it.next();
				if (r.getName().isWild() && !name.isWild()) {
					r = r.withName(name);
				}
				response.addRecord(r, section);
			}
		}
	}

	@SuppressWarnings("unchecked")
	private static RRset setNegativeTTL(final RRset original, final int flags) {
		/*
		 * If DNSSEC is enabled/requested, use the SOA and sigs, otherwise
		 * lower the TTL on the SOA record to the minimum/ncache TTL,
		 * using whichever is lower. Behavior is defined in RFC 2308.
		 * In practice we see Vantio using the minimum from the SOA, while BIND
		 * uses the lowest TTL in the RRset in the authority section. When DNSSEC
		 * is enabled, the TTL for the RRsigs is derived from the minimum of the
		 * SOA via the jdnssec library, hence only modifying the TTL of the SOA
		 * itself in the non-DNSSEC use case below. We would invalidate the existing
		 * RRsigs if we modified the TTL of a signed RRset.
		 */

		// signed RRset and DNSSEC requested; return unmodified
		if (original.sigs().hasNext() && (flags & (FLAG_SIGONLY | FLAG_DNSSECOK)) != 0) {
			return original;
		}

		final RRset rrset = new RRset();
		final Iterator<Record> it = original.rrs();

		while (it.hasNext()) {
			Record record = it.next();

			if (record instanceof SOARecord) {
				final SOARecord soa = (SOARecord) record;

				// the value of the minimum field is less than the actual TTL; adjust
				if (soa.getMinimum() != 0 || soa.getTTL() > soa.getMinimum()) {
					record = new SOARecord(soa.getName(), DClass.IN, soa.getMinimum(), soa.getHost(), soa.getAdmin(),
							soa.getSerial(), soa.getRefresh(), soa.getRetry(), soa.getExpire(),
							soa.getMinimum());
				} // else use the unmodified record
			}

			rrset.addRR(record);
		}

		return rrset;
	}

	private void lookup(final Name qname, final int qtype, final InetAddress clientAddress, final Message response, final int flags, final boolean dnssecRequest, final DNSAccessRecord.Builder builder) {
		lookup(qname, qtype, clientAddress, null, response, 0, flags, dnssecRequest, builder);
	}

	@SuppressWarnings({"unchecked", "PMD.CyclomaticComplexity", "PMD.NPathComplexity"})
	private void lookup(final Name qname, final int qtype, final InetAddress clientAddress, final Zone incomingZone, final Message response, final int iteration, final int flags, final boolean dnssecRequest, final DNSAccessRecord.Builder builder) {
		if (iteration > MAX_ITERATIONS) {
			return;
		}

		Zone zone = incomingZone;

		// this allows us to locate zones for which we are authoritative
		if (zone == null || !qname.subdomain(zone.getOrigin())) {
			zone = trafficRouterManager.getTrafficRouter().getZone(qname, qtype, clientAddress, dnssecRequest, builder);
		}

		// null means we did not find a zone for which we are authoritative
		if (zone == null) {
			if (iteration == 0) {
				// refuse the query if we're not authoritative and we're not recursing
				response.getHeader().setRcode(Rcode.REFUSED);
			}

			return;
		}

		final SetResponse sr = zone.findRecords(qname, qtype);

		if (sr.isSuccessful()) {
			for (final RRset answer : sr.answers()) {
				addRRset(qname, response, answer, Section.ANSWER, flags);
			}

			addAuthority(zone, response, flags);
		} else if (sr.isCNAME()) {
			/*
			 * This is an ugly hack to work around the answers() method not working for CNAMEs.
			 * A CNAME results in isSuccessful() being false, and answers() requires isSuccessful()
			 * to be true. Because of this, we can either use reflection (slow) or use the getNS() method, which
			 * returns the RRset stored internally in "data" and is not actually specific to NS records.
			 * Our CNAME and RRSIGs are in this RRset, so use getNS() despite its name.
			 * Refer to the dnsjava SetResponse code for more information.
			 */
			final RRset rrset = sr.getNS();
			addRRset(qname, response, rrset, Section.ANSWER, flags);

			/*
			 * Allow recursive lookups for CNAME targets; the logic above allows us to
			 * ensure that we only recurse for domains for which we are authoritative.
			 */
			lookup(sr.getCNAME().getTarget(), qtype, clientAddress, zone, response, iteration + 1, flags, dnssecRequest, builder);
		} else if (sr.isNXDOMAIN()) {
			response.getHeader().setRcode(Rcode.NXDOMAIN);
			response.getHeader().setFlag(Flags.AA);
			addDenialOfExistence(qname, zone, response, flags);
			addSOA(zone, response, Section.AUTHORITY, flags);
		} else if (sr.isNXRRSET()) {
			/*
			 * Per RFC 2308 NODATA is inferred by having no records;
			 * NXRRSET is discussed in RFC 2136, but that RFC is for Dynamic DNS updates.
			 * We'll ignore the NXRRSET from the API, and allow the client resolver to
			 * deal with NODATA per RFC 2308:
			 *   "NODATA" - a pseudo RCODE which indicates that the name is valid, for
			 *   the given class, but are no records of the given type.
			 *   A NODATA response has to be inferred from the answer.
			 */

			// The requirements for this are described in RFC 7129
			if ((flags & (FLAG_SIGONLY | FLAG_DNSSECOK)) != 0) {
				final SetResponse ndsr = zone.findRecords(qname, Type.NSEC);
				if (ndsr.isSuccessful()) {
					for (final RRset answer : ndsr.answers()) {
						addRRset(qname, response, answer, Section.AUTHORITY, flags);
					}
				}
			}

			addSOA(zone, response, Section.AUTHORITY, flags);
			response.getHeader().setFlag(Flags.AA);
		}
	}

	public TrafficRouterManager getTrafficRouterManager() {
		return trafficRouterManager;
	}

	public void setTrafficRouterManager(final TrafficRouterManager trafficRouterManager) {
		this.trafficRouterManager = trafficRouterManager;
	}

	public void destroy() {
		/*
		 * Yes, this is odd. We need to call destroy on ZoneManager, but it's static, so
		 * we don't have a Spring bean ref; we do for NameServer, so this method is called.
		 * Given that we know we're shutting down and NameServer relies on ZoneManager,
		 * we'll call destroy while we can without hacking Spring too hard.
		 */
		ZoneManager.destroy();
	}

	public boolean isEcsEnable(final Name qname) {
		return ecsEnable || isDeliveryServiceEcsEnabled(qname);
	}

	public void setEcsEnable(final boolean ecsEnable) {
		this.ecsEnable = ecsEnable;
	}

	public Set<DeliveryService> getEcsEnabledDses() {
		return ecsEnabledDses;
	}

	public void setEcsEnabledDses(final Set<DeliveryService> ecsEnabledDses) {
		this.ecsEnabledDses = ecsEnabledDses;
	}
}
