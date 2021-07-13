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

package org.apache.traffic_control.traffic_router.core.edge;

import org.xbill.DNS.Type;

import java.net.Inet4Address;
import java.net.Inet6Address;
import java.net.InetAddress;


public class InetRecord {
	
	final private InetAddress ad;
	final private long ttl;
	final private int type;
	final private String target;

	public InetRecord(final InetAddress ad, final long ttl) {
		this.ad = ad;
		this.ttl = ttl;
		this.target = null;
		this.type = (ad instanceof Inet4Address) ? Type.A : Type.AAAA;
	}

	public InetRecord(final String alias, final long ttl) {
		this.ad = null;
		this.ttl = ttl;
		this.target = alias;
		this.type = Type.CNAME;
	}

	public InetRecord(final String target, final long ttl, final int type) {
		this.ad = null;
		this.target = target;
		this.ttl = ttl;
		this.type = type;
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
			", target='" + target + '\'' +
			", type=" + Type.string(type) +
			'}';
	}

	public boolean isAlias() {
		return (target != null && type == Type.CNAME);
	}

	public String getAlias() {
		return target;
	}

	public String getTarget() {
		return target;
	}

	public int getType() {
		return type;
	}

	@Override
	@SuppressWarnings("PMD.IfStmtsMustUseBraces")
	public boolean equals(final Object o) {
		if (this == o) return true;
		if (o == null || getClass() != o.getClass()) return false;

		final InetRecord that = (InetRecord) o;

		if (ttl != that.ttl || type != that.type) return false;
		if (ad != null ? !ad.equals(that.ad) : that.ad != null) return false;
		return !(target != null ? !target.equals(that.target) : that.target != null);

	}

	@Override
	public int hashCode() {
		int result = ad != null ? ad.hashCode() : 0;
		result = 31 * result + (int) (ttl ^ (ttl >>> 32));
		result = 31 * result + (int) (type ^ (type >>> 32));
		result = 31 * result + (target != null ? target.hashCode() : 0);
		return result;
	}
}
