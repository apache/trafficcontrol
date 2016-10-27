package com.comcast.cdn.traffic_control.traffic_router.core.dns;

import org.xbill.DNS.DNSKEYRecord;
import org.xbill.DNS.DSRecord;
import org.xbill.DNS.Name;
import org.xbill.DNS.Record;

import java.io.IOException;
import java.security.GeneralSecurityException;
import java.util.Date;
import java.util.List;

public class ZoneSignerImpl implements ZoneSigner {
	@Override
	public List<Record> signZone(final Name name, final List<Record> records, final List<DnsSecKeyPair> kskPairs, final List<DnsSecKeyPair> zskPairs,
		final Date inception, final Date expiration, final boolean fullySignKeySet, final int digestId) throws IOException, GeneralSecurityException {
		return null;
	}

	@Override
	public DSRecord calculateDSRecord(final DNSKEYRecord dnskeyRecord, final int digestId, final long ttl) {
		return null;
	}
}
