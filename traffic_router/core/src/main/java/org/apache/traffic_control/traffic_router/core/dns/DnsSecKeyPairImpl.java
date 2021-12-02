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

import org.apache.traffic_control.traffic_router.core.util.JsonUtils;
import org.apache.traffic_control.traffic_router.core.util.JsonUtilsException;
import org.apache.traffic_control.traffic_router.secure.BindPrivateKey;
import com.fasterxml.jackson.databind.JsonNode;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.xbill.DNS.DNSKEYRecord;
import org.xbill.DNS.DNSSEC;
import org.xbill.DNS.Master;
import org.xbill.DNS.Name;
import org.xbill.DNS.Record;
import org.xbill.DNS.Type;

import java.io.ByteArrayInputStream;
import java.io.IOException;
import java.io.InputStream;
import java.security.PrivateKey;
import java.security.PublicKey;
import java.util.Base64.Decoder;
import java.util.Calendar;
import java.util.Date;

import static java.util.Base64.getMimeDecoder;

public class DnsSecKeyPairImpl implements DnsSecKeyPair {
	private static final Logger LOGGER = LogManager.getLogger(DnsSecKeyPairImpl.class);
	private long ttl;
	private Date inception;
	private Date effective;
	private Date expiration;
	private String name;
	private DNSKEYRecord dnskeyRecord;
	private PrivateKey privateKey;

	public DnsSecKeyPairImpl(final JsonNode keyPair, final long defaultTTL) throws JsonUtilsException, IOException {
		this.inception = new Date(1000L * JsonUtils.getLong(keyPair, "inceptionDate"));
		this.effective = new Date(1000L * JsonUtils.getLong(keyPair, "effectiveDate"));
		this.expiration = new Date(1000L * JsonUtils.getLong(keyPair, "expirationDate"));
		this.ttl = JsonUtils.optLong(keyPair, "ttl", defaultTTL);
		this.name = JsonUtils.getString(keyPair, "name").toLowerCase();

		final Decoder mimeDecoder = getMimeDecoder();
		try {
			privateKey = new BindPrivateKey().decode(new String(mimeDecoder.decode(JsonUtils.getString(keyPair, "private"))));
		} catch (Exception e) {
			LOGGER.error("Failed to decode PKCS1 key from json data!: " + e.getMessage(), e);
		}

		final byte[] publicKey = mimeDecoder.decode(JsonUtils.getString(keyPair, "public"));

		try (InputStream in = new ByteArrayInputStream(publicKey)) {
			final Master master = new Master(in, new Name(name), ttl);

			Record record;
			while ((record = master.nextRecord()) != null) {
				if (record.getType() == Type.DNSKEY) {
					this.dnskeyRecord = (DNSKEYRecord) record;
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
	public DNSKEYRecord getDNSKEYRecord() {
		return dnskeyRecord;
	}

	@Override
	public PrivateKey getPrivate() {
		return privateKey;
	}

	@Override
	public PublicKey getPublic() {
		try {
			return dnskeyRecord.getPublicKey();
		} catch (DNSSEC.DNSSECException e) {
			LOGGER.error("Failed to extract public key from DNSKEY record for " + name + " : " + e.getMessage(), e);
		}
		return null;
	}

	@SuppressWarnings("PMD.OverrideBothEqualsAndHashcode")
	public boolean equals(final Object obj) {
		final DnsSecKeyPairImpl okp = (DnsSecKeyPairImpl) obj;

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
