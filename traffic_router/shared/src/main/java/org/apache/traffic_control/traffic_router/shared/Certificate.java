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

package org.apache.traffic_control.traffic_router.shared;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

@JsonIgnoreProperties(ignoreUnknown = true)
public class Certificate {
	@JsonProperty
	private String crt;

	@JsonProperty
	private String key;

	public String getCrt() {
		return crt;
	}

	public void setCrt(final String crt) {
		this.crt = crt;
	}

	public String getKey() {
		return key;
	}

	public void setKey(final String key) {
		this.key = key;
	}

	@Override
	public boolean equals(final Object o) {
		if (this == o) return true;
		if (o == null || getClass() != o.getClass()) return false;

		final Certificate that = (Certificate) o;

		if (crt != null ? !crt.equals(that.crt) : that.crt != null) return false;
		return key != null ? key.equals(that.key) : that.key == null;
	}

	@Override
	public int hashCode() {
		int result = crt != null ? crt.hashCode() : 0;
		result = 31 * result + (key != null ? key.hashCode() : 0);
		return result;
	}
}
