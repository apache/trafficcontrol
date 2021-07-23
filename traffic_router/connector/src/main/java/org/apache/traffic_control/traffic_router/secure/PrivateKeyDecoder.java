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
import java.security.PrivateKey;
import java.util.Base64;

public class PrivateKeyDecoder {
	public PrivateKey decode(final String data) throws IOException, GeneralSecurityException {
		final String decodedData = new String(Base64.getMimeDecoder().decode(data.getBytes()));
		return decodedData.contains(Pkcs1.HEADER) ? new Pkcs1(decodedData).getPrivateKey() : new Pkcs8(decodedData).getPrivateKey();
	}
}
