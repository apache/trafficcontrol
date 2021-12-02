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

import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

import java.io.IOException;
import java.security.GeneralSecurityException;
import java.security.spec.KeySpec;
import java.security.spec.PKCS8EncodedKeySpec;
import java.util.Base64;

public class Pkcs8 extends Pkcs {
	private static final Logger LOGGER = LogManager.getLogger(Pkcs8.class);
	public static final String HEADER = "-----BEGIN PRIVATE KEY-----";
	public static final String FOOTER = "-----END PRIVATE KEY-----";

	public Pkcs8(final String data) throws IOException, GeneralSecurityException {
		super(data);
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
		try {
			return new PKCS8EncodedKeySpec(Base64.getDecoder().decode((data.getBytes())));
		} catch (Exception e) {
			LOGGER.error("Failed to create PKCS8 Encoded Key Spec " + e.getClass().getCanonicalName() + ": " + e.getMessage(), e);
		}
		return null;
	}
}
