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

import org.bouncycastle.jce.provider.BouncyCastleProvider;

import java.io.IOException;
import java.security.GeneralSecurityException;
import java.security.KeyFactory;
import java.security.PrivateKey;
import java.security.PublicKey;
import java.security.Security;
import java.security.spec.KeySpec;

@SuppressWarnings("PMD.AbstractNaming")
public abstract class Pkcs {
	private final String data;
	private final PrivateKey privateKey;
	private PublicKey publicKey;
	private KeySpec keySpec;
	private KeySpec publicKeySpec;

	public Pkcs(final String data) throws IOException, GeneralSecurityException {
		this.data = data;
		keySpec = toKeySpec(data);
		Security.addProvider(new BouncyCastleProvider());
		privateKey = KeyFactory.getInstance("RSA", "BC").generatePrivate(keySpec);
	}

	public Pkcs(final String privateData, final String publicData) throws IOException, GeneralSecurityException {
		this.data = privateData;
		keySpec = toKeySpec(data);
		privateKey = KeyFactory.getInstance("RSA", "BC").generatePrivate(keySpec);
		publicKeySpec = toKeySpec(publicData);
		publicKey = KeyFactory.getInstance("RSA", "BC").generatePublic(publicKeySpec);
	}

	public String getData() {
		return data;
	}

	public KeySpec getKeySpec() {
		return keySpec;
	}

	public KeySpec getPublicKeySpec() {
		return publicKeySpec;
	}

	public void setKeySpec(final KeySpec keySpec) {
		this.keySpec = keySpec;
	}

	public PrivateKey getPrivateKey() {
		return privateKey;
	}

	public PublicKey getPublicKey() {
		return publicKey;
	}

	public abstract String getHeader();

	public abstract String getFooter();

	private String stripHeaderAndFooter(final String data) {
		return data.replaceAll(getHeader(), "").replaceAll(getFooter(), "").replaceAll("\\s", "");
	}

	protected abstract KeySpec decodeKeySpec(final String data) throws IOException, GeneralSecurityException;

	private KeySpec toKeySpec(final String data) throws IOException, GeneralSecurityException {
		return decodeKeySpec(stripHeaderAndFooter(data));
	}
}
