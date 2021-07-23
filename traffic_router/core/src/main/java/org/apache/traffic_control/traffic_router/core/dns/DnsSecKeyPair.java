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
