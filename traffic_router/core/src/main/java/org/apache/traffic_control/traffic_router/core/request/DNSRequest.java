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

package org.apache.traffic_control.traffic_router.core.request;

import org.xbill.DNS.Name;
import org.xbill.DNS.Zone;

public class DNSRequest extends Request {
	private final Name name;
	private final String zoneName;
	private final int queryType;
	private boolean dnssec = false;

	public DNSRequest(final String zoneName, final Name name, final int queryType) {
		super();

		this.queryType = queryType;
		this.name = name;
		this.zoneName = zoneName;
	}

	public DNSRequest(final Zone zone, final Name name, final int queryType) {
		super();

		this.queryType = queryType;
		this.name = name;
		this.zoneName = zone.getOrigin().toString().toLowerCase();
	}

	public int getQueryType() {
		return queryType;
	}

	@Override
	public String getType() {
		return "dns";
	}

	public boolean isDnssec() {
		return dnssec;
	}

	public void setDnssec(final boolean dnssec) {
		this.dnssec = dnssec;
	}

	public Name getName() {
		return name;
	}

	public String getZoneName() {
		return zoneName;
	}
}
