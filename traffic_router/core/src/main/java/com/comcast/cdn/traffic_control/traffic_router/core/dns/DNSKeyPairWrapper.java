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

import java.io.ByteArrayInputStream;
import java.io.IOException;
import java.io.InputStream;
import java.util.Date;

import javax.xml.bind.DatatypeConverter;

import org.apache.log4j.Logger;
import org.json.JSONException;
import org.json.JSONObject;
import org.xbill.DNS.DNSKEYRecord;
import org.xbill.DNS.Master;
import org.xbill.DNS.Name;
import org.xbill.DNS.Record;
import org.xbill.DNS.Type;

import com.verisignlabs.dnssec.security.DnsKeyPair;

public class DNSKeyPairWrapper {
	private static final Logger LOGGER = Logger.getLogger(DNSKeyPairWrapper.class);

	private DnsKeyPair dnsKeyPair;
	private long ttl;
	private Date inception;
	private Date effective;
	private Date expiration;
	private String name;

	public DNSKeyPairWrapper(final JSONObject keyPair) throws JSONException, IOException {
		this.inception = new Date (1000L * keyPair.getLong("inceptionDate"));
		this.effective = new Date (1000L * keyPair.getLong("effectiveDate"));
		this.expiration = new Date (1000L * keyPair.getLong("expirationDate"));
		this.ttl = keyPair.getLong("ttl");
		this.name = keyPair.getString("name");
		//this.status = keyPair.getString("status"); // this field is used by Traffic Ops; we detect expiration by using the above dates

		final byte[] privateKey = DatatypeConverter.parseBase64Binary(keyPair.getString("private"));
		final byte[] publicKey = DatatypeConverter.parseBase64Binary(keyPair.getString("public"));
		final InputStream in = new ByteArrayInputStream(publicKey);

		final Master master = new Master(in, new Name(name), ttl);
		Record record = null;
		final DnsKeyPair dkp = new DnsKeyPair();
		dkp.setPrivateKeyString(new String(privateKey));

		while ((record = master.nextRecord()) != null) {
			if (record.getType() == Type.DNSKEY) {
				dkp.setDNSKEYRecord((DNSKEYRecord) record);
				LOGGER.debug("record name is " + record.getName() + "; domain is " + name);
				break;
			}
		}

		this.dnsKeyPair = dkp;
	}

	public DnsKeyPair getDnsKeyPair() {
		return dnsKeyPair;
	}

	public void setDnsKeyPair(final DnsKeyPair dnsKeyPair) {
		this.dnsKeyPair = dnsKeyPair;
	}

	public long getTtl() {
		return ttl;
	}

	public void setTtl(final long ttl) {
		this.ttl = ttl;
	}

	public String getName() {
		return name;
	}

	public void setName(final String name) {
		this.name = name;
	}

	public Date getInception() {
		return inception;
	}

	public void setInception(final Date inception) {
		this.inception = inception;
	}

	public Date getEffective() {
		return effective;
	}

	public void setEffective(final Date effective) {
		this.effective = effective;
	}

	public Date getExpiration() {
		return expiration;
	}

	public void setExpiration(final Date expiration) {
		this.expiration = expiration;
	}
}
