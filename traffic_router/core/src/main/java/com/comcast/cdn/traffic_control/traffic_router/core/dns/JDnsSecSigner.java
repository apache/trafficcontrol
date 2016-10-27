package com.comcast.cdn.traffic_control.traffic_router.core.dns;

import com.verisignlabs.dnssec.security.DnsKeyPair;
import com.verisignlabs.dnssec.security.JCEDnsSecSigner;
import com.verisignlabs.dnssec.security.SignUtils;
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
	@Override
	public List<Record> signZone(final Name name, final List<Record> records, final List<DnsSecKeyPair> kskPairs, final List<DnsSecKeyPair> zskPairs,
		final Date inception, final Date expiration, final boolean fullySignKeySet, final int digestId) throws IOException, GeneralSecurityException {

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
		return SignUtils.calculateDSRecord(dnskeyRecord, DSRecord.SHA256_DIGEST_ID, ttl);
	}
}
