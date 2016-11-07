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

package com.comcast.cdn.traffic_control.traffic_router.core.cache;

import java.net.Inet4Address;
import java.net.Inet6Address;
import java.net.InetAddress;


public class InetRecord {
	
	final private InetAddress ad;
	final private long ttl;
	final private String alias;
	
	public InetRecord(final InetAddress ad, final long ttl) {
		this.ad = ad;
		this.ttl = ttl;
		this.alias = null;
	}

	public InetRecord(final String alias, final long ttl) {
		this.ad = null;
		this.ttl = ttl;
		this.alias = alias;
	}

	public boolean isInet4() {
		return ad instanceof Inet4Address;
	}
	public boolean isInet6() {
		return ad instanceof Inet6Address;
	}

	public long getTTL() {
		return ttl;
	}

	public InetAddress getAddress() {
		return ad;
	}

	@Override
	public String toString() {
		return "InetRecord{" +
			"ad=" + ad +
			", ttl=" + ttl +
			", alias='" + alias + '\'' +
			'}';
	}

	public boolean isAlias() {
		return (alias != null);
	}

	public String getAlias() {
		return alias;
	}

	@Override
	@SuppressWarnings("PMD.IfStmtsMustUseBraces")
	public boolean equals(final Object o) {
		if (this == o) return true;
		if (o == null || getClass() != o.getClass()) return false;

		final InetRecord that = (InetRecord) o;

		if (ttl != that.ttl) return false;
		if (ad != null ? !ad.equals(that.ad) : that.ad != null) return false;
		return !(alias != null ? !alias.equals(that.alias) : that.alias != null);

	}

	@Override
	public int hashCode() {
		int result = ad != null ? ad.hashCode() : 0;
		result = 31 * result + (int) (ttl ^ (ttl >>> 32));
		result = 31 * result + (alias != null ? alias.hashCode() : 0);
		return result;
	}
}
