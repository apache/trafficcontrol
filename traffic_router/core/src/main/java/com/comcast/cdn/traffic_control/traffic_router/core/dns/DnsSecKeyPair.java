package com.comcast.cdn.traffic_control.traffic_router.core.dns;

import org.xbill.DNS.DNSKEYRecord;

import java.security.PrivateKey;
import java.security.PublicKey;
import java.util.Date;

public interface DnsSecKeyPair {
	long getTTL();

	void setTTL(long ttl);

	String getName();

	void setName(String name);

	Date getInception();

	void setInception(Date inception);

	Date getEffective();

	void setEffective(Date effective);

	Date getExpiration();

	void setExpiration(Date expiration);

	boolean isKeySigningKey();

	boolean isExpired();

	boolean isUsable();

	boolean isKeyCached(long maxTTL);

	boolean isOlder(DnsSecKeyPair other);

	boolean isNewer(DnsSecKeyPair other);

	PrivateKey getPrivate();

	PublicKey getPublic();

	DNSKEYRecord getDNSKEYRecord();

	@Override
	@SuppressWarnings("PMD.OverrideBothEqualsAndHashcode")
	boolean equals(Object obj);

	@Override
	String toString();
}
