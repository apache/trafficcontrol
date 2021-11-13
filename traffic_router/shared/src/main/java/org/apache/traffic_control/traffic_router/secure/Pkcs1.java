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

import java.io.IOException;
import java.security.GeneralSecurityException;
import java.security.spec.KeySpec;

public class Pkcs1 extends Pkcs {

	static public final String HEADER = "-----BEGIN RSA PRIVATE KEY-----";
	static public final String FOOTER = "-----END RSA PRIVATE KEY-----";

	public Pkcs1(final String data) throws IOException, GeneralSecurityException {
		super(data);
	}

	public Pkcs1(final String privateData, final String publicData) throws IOException, GeneralSecurityException {
		super(privateData,publicData);
	}

	@Override
	public String getHeader() {
		return HEADER;
	}

	@Override
	public String getFooter() {
		return FOOTER;
	}

	@Override
	protected KeySpec decodeKeySpec(final String data) throws IOException, GeneralSecurityException {
		return new Pkcs1KeySpecDecoder().decode(data);
	}
}
