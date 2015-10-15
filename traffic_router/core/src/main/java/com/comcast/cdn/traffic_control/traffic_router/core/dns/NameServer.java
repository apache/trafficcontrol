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

import java.net.InetAddress;
import java.util.ArrayList;
import java.util.Collections;
import java.util.Iterator;
import java.util.List;

import org.apache.log4j.Logger;
import org.xbill.DNS.CNAMERecord;
import org.xbill.DNS.DClass;
import org.xbill.DNS.ExtendedFlags;
import org.xbill.DNS.Flags;
import org.xbill.DNS.Message;
import org.xbill.DNS.Name;
import org.xbill.DNS.OPTRecord;
import org.xbill.DNS.RRset;
import org.xbill.DNS.Rcode;
import org.xbill.DNS.Record;
import org.xbill.DNS.Section;
import org.xbill.DNS.SetResponse;
import org.xbill.DNS.Type;
import org.xbill.DNS.Zone;

import com.comcast.cdn.traffic_control.traffic_router.core.router.TrafficRouterManager;

public class NameServer {
	private static final int MAX_SUPPORTED_EDNS_VERS = 0;
	private static final int MAX_ITERATIONS = 6;
	private static final int NUM_SECTIONS = 4;
	private static final int FLAG_DNSSECOK = 1;
	private static final int FLAG_SIGONLY = 2;

	private static final Logger LOGGER = Logger.getLogger(NameServer.class);
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

	@SuppressWarnings({"PMD.CyclomaticComplexity", "PMD.NPathComplexity"})
	private void addAnswers(final Message request, final Message response, final InetAddress clientAddress, final DNSAccessRecord.Builder builder) {
		final Record question = request.getQuestion();
		final int qclass = question.getDClass();
		final Name qname = question.getName();
		final OPTRecord qopt = request.getOPT();
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

		final Zone zone = trafficRouterManager.getTrafficRouter().getZone(qname, qtype, clientAddress, dnssecRequest, builder);

		if (zone == null) {
			response.getHeader().setRcode(Rcode.REFUSED);
			return;
		}

		lookup(qname, qtype, zone, response, 0, flags);

		if (qopt != null && flags == FLAG_DNSSECOK) {
			final int optflags = ExtendedFlags.DO;
			final OPTRecord opt = new OPTRecord(1280, (byte) 0, (byte) 0, optflags);
			response.addRecord(opt, Section.ADDITIONAL);
		}
	}

	private static void addAuthority(final Zone zone, final Message response, final int flags) {
		final RRset authority = zone.getNS();
		addRRset(authority.getName(), response, authority, Section.AUTHORITY, flags);
		response.getHeader().setFlag(Flags.AA);
	}

	private static void addSOA(final Zone zone, final Message response, final int section, final int flags) {
		// we locate the SOA this way so that we can ensure we get the RRSIGs rather than just the one SOA Record
		final SetResponse fsoa = zone.findRecords(zone.getOrigin(), Type.SOA);

		for (final RRset answer : fsoa.answers()) {
			addRRset(zone.getOrigin(), response, answer, section, flags);
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

	@SuppressWarnings({"unchecked", "PMD.CyclomaticComplexity", "PMD.NPathComplexity"})
	private static void lookup(final Name qname, final int qtype, final Zone zone, final Message response, final int iteration, final int flags) {
		if (iteration > MAX_ITERATIONS) {
			return;
		}

		final SetResponse sr = zone.findRecords(qname, qtype);

		if (sr.isSuccessful()) {
			for (final RRset answer : sr.answers()) {
				addRRset(qname, response, answer, Section.ANSWER, flags);
			}

			addAuthority(zone, response, flags);
		} else if (sr.isNXDOMAIN()) {
			response.getHeader().setRcode(Rcode.NXDOMAIN);

			// The requirements for this are described in RFC 7129
			if ((flags & (FLAG_SIGONLY | FLAG_DNSSECOK)) != 0) {
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

			addSOA(zone, response, Section.AUTHORITY, flags);
			response.getHeader().setFlag(Flags.AA);
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
		} else if (sr.isCNAME()) {
			final CNAMERecord cname = sr.getCNAME();
			final RRset cnameSet = new RRset(cname);
			addRRset(qname, response, cnameSet, Section.ANSWER, flags);
			lookup(cname.getTarget(), qtype, zone, response, iteration + 1, flags);
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
		LOGGER.info("Calling destroy on ZoneManager");
		ZoneManager.destroy();
	}
}
