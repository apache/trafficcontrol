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

package org.apache.traffic_control.traffic_router.secure;

import java.security.PrivateKey;
import java.security.cert.X509Certificate;
import java.util.Arrays;
import java.util.Objects;

public class HandshakeData {
	private final String deliveryService;
	private final String hostname;
	private final X509Certificate[] certificateChain;
	private PrivateKey privateKey;

	public HandshakeData(final String deliveryService, final String hostname, final X509Certificate[] certificateChain, final PrivateKey privateKey) {
		this.deliveryService = deliveryService;
		this.hostname = hostname;
		this.certificateChain = certificateChain;
		this.privateKey = privateKey;
	}

	public String getDeliveryService() {
		return deliveryService;
	}

	public String getHostname() {
		return hostname;
	}

	public X509Certificate[] getCertificateChain() {
		return certificateChain;
	}

	public PrivateKey getPrivateKey() {
		return privateKey;
	}

	public void setPrivateKey(final PrivateKey privateKey) {
		this.privateKey = privateKey;
	}

	@Override
	public boolean equals(final Object o) {
		if (this == o) {return true;}
		if (!(o instanceof HandshakeData)) {return false;}
		final HandshakeData that = (HandshakeData) o;
		return Objects.equals(deliveryService, that.deliveryService) &&
				Objects.equals(hostname, that.hostname) &&
				Arrays.equals(certificateChain, that.certificateChain) &&
				Objects.equals(privateKey, that.privateKey);
	}

	@Override
	public int hashCode() {
		int result = Objects.hash(deliveryService, hostname, privateKey);
		result = 31 * result + Arrays.hashCode(certificateChain);
		return result;
	}
}
