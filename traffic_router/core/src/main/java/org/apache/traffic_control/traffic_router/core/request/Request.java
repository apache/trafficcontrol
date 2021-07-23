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

import org.apache.commons.lang3.builder.EqualsBuilder;
import org.apache.commons.lang3.builder.HashCodeBuilder;

public class Request {

	private String clientIP;
	private String hostname;

	@Override
	public boolean equals(final Object obj) {
		if (this == obj) {
			return true;
		} else if (obj instanceof Request) {
			final Request rhs = (Request) obj;
			return new EqualsBuilder()
			.append(getClientIP(), rhs.getClientIP())
			.append(getHostname(), rhs.getHostname())
			.isEquals();
		} else {
			return false;
		}
	}

	public String getClientIP() {
		return clientIP;
	}

	public String getHostname() {
		return hostname;
	}

	@Override
	public int hashCode() {
		return new HashCodeBuilder(1, 31)
		.append(getClientIP())
		.append(getHostname())
		.toHashCode();
	}

	public void setClientIP(final String clientIP) {
		this.clientIP = clientIP;
	}

	public void setHostname(final String hostname) {
		if(hostname==null) {
			this.hostname = null;
			return;
		}
		this.hostname = hostname.toLowerCase();
	}

	public String getType() {
		return "unknown";
	}
}
