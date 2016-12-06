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

import com.verisignlabs.dnssec.security.DnsKeyPair;
import com.verisignlabs.dnssec.security.JCEDnsSecSigner;
import com.verisignlabs.dnssec.security.SignUtils;
import org.apache.log4j.Logger;
import org.xbill.DNS.DNSKEYRecord;
import org.xbill.DNS.DSRecord;
import org.xbill.DNS.Name;
import org.xbill.DNS.Record;

import java.io.IOException;
import java.security.GeneralSecurityException;
import java.util.ArrayList;
import java.util.Date;
import java.util.List;

public class JDnsSecSigner implements ZoneSigner {
	private static final Logger LOGGER = Logger.getLogger(JDnsSecSigner.class);
	@Override
	public List<Record> signZone(final Name name, final List<Record> records, final List<DnsSecKeyPair> kskPairs, final List<DnsSecKeyPair> zskPairs,
		final Date inception, final Date expiration, final boolean fullySignKeySet, final int digestId) throws IOException, GeneralSecurityException {
		LOGGER.info("Signing records, name for first record is " + records.get(0).getName());
		final List<DnsKeyPair> kPairs = new ArrayList<>();
		final List<DnsKeyPair> zPairs = new ArrayList<>();

		for (final DnsSecKeyPair keyPair : kskPairs) {
			if (keyPair instanceof DnsKeyPair) {
				kPairs.add((DnsKeyPair) keyPair);
			} else {
				throw new IllegalArgumentException("kskPairs contains non jdnssec object!");
			}
		}

		for (final DnsSecKeyPair keyPair : zskPairs) {
			if (keyPair instanceof DnsKeyPair) {
				zPairs.add((DnsKeyPair) keyPair);
			} else {
				throw new IllegalArgumentException("zskPairs contains non jdnssec object!");
			}
		}

		final JCEDnsSecSigner signer = new JCEDnsSecSigner(false);

		return signer.signZone(name, records, kPairs, zPairs, inception, expiration, fullySignKeySet, digestId);
	}

	@Override
	public DSRecord calculateDSRecord(final DNSKEYRecord dnskeyRecord, final int digestId, final long ttl) {
		LOGGER.info("Calculating DS Records for " + dnskeyRecord.getName());
		return SignUtils.calculateDSRecord(dnskeyRecord, DSRecord.SHA256_DIGEST_ID, ttl);
	}
}
