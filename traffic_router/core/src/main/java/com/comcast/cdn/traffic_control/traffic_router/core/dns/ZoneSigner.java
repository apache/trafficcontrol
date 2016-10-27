package com.comcast.cdn.traffic_control.traffic_router.core.dns;

import org.xbill.DNS.DNSKEYRecord;
import org.xbill.DNS.DSRecord;
import org.xbill.DNS.Name;
import org.xbill.DNS.Record;

import java.io.IOException;
import java.security.GeneralSecurityException;
import java.util.Date;
import java.util.List;

public interface ZoneSigner {
	List<Record> signZone(Name name, List<Record> records, List<DnsSecKeyPair> kskPairs, List<DnsSecKeyPair> zskPairs,
	                      Date inception, Date expiration, boolean fullySignKeySet, int digestId) throws IOException, GeneralSecurityException;
	DSRecord calculateDSRecord(DNSKEYRecord dnskeyRecord, int digestId, long ttl);
}
