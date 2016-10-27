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

package com.comcast.cdn.traffic_control.traffic_router.core.dns;

import java.io.ByteArrayInputStream;
import java.io.IOException;
import java.io.InputStream;
import java.util.Calendar;
import java.util.Date;

import javax.xml.bind.DatatypeConverter;

import org.json.JSONException;
import org.json.JSONObject;
import org.xbill.DNS.DNSKEYRecord;
import org.xbill.DNS.Master;
import org.xbill.DNS.Name;
import org.xbill.DNS.Record;
import org.xbill.DNS.Type;

import com.verisignlabs.dnssec.security.DnsKeyPair;

public class DNSKeyPairWrapper extends DnsKeyPair implements DnsSecKeyPair {
	private long ttl;
	private Date inception;
	private Date effective;
	private Date expiration;
	private String name;

	public DNSKeyPairWrapper(final JSONObject keyPair, final long defaultTTL) throws JSONException, IOException {
		this.inception = new Date(1000L * keyPair.getLong("inceptionDate"));
		this.effective = new Date(1000L * keyPair.getLong("effectiveDate"));
		this.expiration = new Date(1000L * keyPair.getLong("expirationDate"));
		this.ttl = keyPair.optLong("ttl", defaultTTL);
		this.name = keyPair.getString("name");

		final byte[] privateKey = DatatypeConverter.parseBase64Binary(keyPair.getString("private"));
		final byte[] publicKey = DatatypeConverter.parseBase64Binary(keyPair.getString("public"));

		try (InputStream in = new ByteArrayInputStream(publicKey)) {
			final Master master = new Master(in, new Name(name), ttl);
			setPrivateKeyString(new String(privateKey));

			Record record;
			while ((record = master.nextRecord()) != null) {
				if (record.getType() == Type.DNSKEY) {
					setDNSKEYRecord((DNSKEYRecord) record);
					break;
				}
			}
		}
	}

	@Override
	public long getTTL() {
		return ttl;
	}

	@Override
	public void setTTL(final long ttl) {
		this.ttl = ttl;
	}

	@Override
	public String getName() {
		return name;
	}

	@Override
	public void setName(final String name) {
		this.name = name;
	}

	@Override
	public Date getInception() {
		return inception;
	}

	@Override
	public void setInception(final Date inception) {
		this.inception = inception;
	}

	@Override
	public Date getEffective() {
		return effective;
	}

	@Override
	public void setEffective(final Date effective) {
		this.effective = effective;
	}

	@Override
	public Date getExpiration() {
		return expiration;
	}

	@Override
	public void setExpiration(final Date expiration) {
		this.expiration = expiration;
	}

	@Override
	public boolean isKeySigningKey() {
		return ((getDNSKEYRecord().getFlags() & DNSKEYRecord.Flags.SEP_KEY) != 0);
	}

	@Override
	public boolean isExpired() {
		return getExpiration().before(Calendar.getInstance().getTime());
	}

	@Override
	public boolean isUsable() {
		final Date now = Calendar.getInstance().getTime();
		return getEffective().before(now);
	}

	@Override
	public boolean isKeyCached(final long maxTTL) {
		return getExpiration().after(new Date(System.currentTimeMillis() - (maxTTL * 1000)));
	}

	@Override
	public boolean isOlder(final DnsSecKeyPair other) {
		return getEffective().before(other.getEffective());
	}

	@Override
	public boolean isNewer(final DnsSecKeyPair other) {
		return getEffective().after(other.getEffective());
	}

	@Override
	@SuppressWarnings("PMD.OverrideBothEqualsAndHashcode")
	public boolean equals(final Object obj) {
		final DNSKeyPairWrapper okp = (DNSKeyPairWrapper) obj;

		if (!this.getDNSKEYRecord().equals(okp.getDNSKEYRecord())) {
			return false;
		} else if (!this.getPrivate().equals(okp.getPrivate())) {
			return false;
		} else if (!this.getPublic().equals(okp.getPublic())) {
			return false;
		} else if (!getEffective().equals(okp.getEffective())) {
			return false;
		} else if (!getExpiration().equals(okp.getExpiration())) {
			return false;
		} else if (!getInception().equals(okp.getInception())) {
			return false;
		} else if (!getName().equals(okp.getName())) {
			return false;
		} else if (getTTL() != okp.getTTL()) {
			return false;
		}

		return true;
	}

	@Override
	public String toString() {
		final StringBuilder sb = new StringBuilder();
		sb.append("name=").append(name)
			.append(" ttl=").append(getTTL())
			.append(" ksk=").append(isKeySigningKey())
			.append(" inception=\"");
		sb.append(getInception());
		sb.append("\" effective=\"");
		sb.append(getEffective());
		sb.append("\" expiration=\"");
		sb.append(getExpiration()).append('"');

		return sb.toString();
	}
}
