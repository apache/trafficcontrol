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
import java.util.Iterator;

import org.apache.log4j.Logger;
import org.xbill.DNS.CNAMERecord;
import org.xbill.DNS.DClass;
import org.xbill.DNS.Flags;
import org.xbill.DNS.Message;
import org.xbill.DNS.Name;
import org.xbill.DNS.OPTRecord;
import org.xbill.DNS.RRset;
import org.xbill.DNS.Rcode;
import org.xbill.DNS.Record;
import org.xbill.DNS.Section;
import org.xbill.DNS.SetResponse;
import org.xbill.DNS.Zone;

import com.comcast.cdn.traffic_control.traffic_router.core.router.TrafficRouterManager;

public class NameServer {
	private static final int MAX_SUPPORTED_EDNS_VERS = 0;
	private static final int MAX_ITERATIONS = 6;
	private static final int NUM_SECTIONS = 4;

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
	public Message query(final Message request, final InetAddress clientAddress) {
		final Message response = new Message();
		try {
			addQuestion(request, response);
			addAnswers(request, response, clientAddress);
//		} catch (final DNSException e) {
//			LOGGER.error(e.getMessage(), e);
//			response.getHeader().setRcode(e.getRcode());
		} catch (final RuntimeException e) {
			LOGGER.error(e.getMessage(), e);
			response.getHeader().setRcode(Rcode.SERVFAIL);
		}

		return response;
	}

	private void addAnswers(final Message request, final Message response, final InetAddress clientAddress) {
		final Record question = request.getQuestion();
		final int qclass = question.getDClass();
		final int qtype = question.getType();
		final Name qname = question.getName();

		if ((request.getOPT() != null) && (request.getOPT().getVersion() > MAX_SUPPORTED_EDNS_VERS)) {
			response.getHeader().setRcode(Rcode.NOTIMP);
			final OPTRecord opt = new OPTRecord(0, Rcode.BADVERS, MAX_SUPPORTED_EDNS_VERS);
			response.addRecord(opt, Section.ADDITIONAL);
			return;
		}

		if ((qclass != DClass.IN) && (qclass != DClass.ANY)) {
			response.getHeader().setRcode(Rcode.REFUSED);
			return;
		}

		final Zone zone = trafficRouterManager.getTrafficRouter().getDynamicZone(qname, qtype, clientAddress);
		if(zone == null) {
			response.getHeader().setRcode(Rcode.REFUSED);
			return;
		}
		lookup(qname, qtype, zone, response, 0);
	}

	private static void addAuthority(final Zone zone, final Message response) {
		final RRset authority = zone.getNS();
		addRRset(authority.getName(), response, authority, Section.AUTHORITY);
		response.getHeader().setFlag(Flags.AA);
	}

	private static void addQuestion(final Message request, final Message response) {
		response.getHeader().setID(request.getHeader().getID());
		response.getHeader().setFlag(Flags.QR);
		if (request.getHeader().getFlag(Flags.RD)) {
			response.getHeader().setFlag(Flags.RD);
		}
		response.addRecord(request.getQuestion(), Section.QUESTION);
	}

	private static void addRRset(final Name name, final Message response, final RRset rrset, final int section) {
		for (int s = 1; s < NUM_SECTIONS; s++) {
			if (response.findRRset(name, rrset.getType(), s)) {
				return;
			}
		}
		@SuppressWarnings("unchecked")
		final Iterator<Record> it = rrset.rrs();
		while (it.hasNext()) {
			// TODO randomize if NS... or always?
			Record r = it.next();
			if (r.getName().isWild() && !name.isWild()) {
				r = r.withName(name);
			}
			response.addRecord(r, section);
		}
	}

	private static void lookup(final Name qname, final int qtype, final Zone zone, final Message response, final int iteration) {
		if (iteration > MAX_ITERATIONS) {
			return;
		}

		final SetResponse sr = zone.findRecords(qname, qtype);

		if (sr.isSuccessful()) {
			final RRset[] answers = sr.answers();

			for (final RRset answer : answers) {
				addRRset(qname, response, answer, Section.ANSWER);
			}

			addAuthority(zone, response);
		} else if (sr.isNXDOMAIN()) {
			response.getHeader().setRcode(Rcode.NXDOMAIN);
			response.addRecord(zone.getSOA(), Section.AUTHORITY);
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

			response.addRecord(zone.getSOA(), Section.AUTHORITY);
		} else if (sr.isCNAME()) {
			final CNAMERecord cname = sr.getCNAME();
			final RRset cnameSet = new RRset(cname);
			addRRset(qname, response, cnameSet, Section.ANSWER);
			lookup(cname.getTarget(), qtype, zone, response, iteration + 1);
		}
	}

	public TrafficRouterManager getTrafficRouterManager() {
		return trafficRouterManager;
	}

	public void setTrafficRouterManager(final TrafficRouterManager trafficRouterManager) {
		this.trafficRouterManager = trafficRouterManager;
	}
}
