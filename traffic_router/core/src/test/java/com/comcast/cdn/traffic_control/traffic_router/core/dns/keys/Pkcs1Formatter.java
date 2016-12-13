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

package com.comcast.cdn.traffic_control.traffic_router.core.dns.keys;

import sun.security.rsa.RSAPrivateCrtKeyImpl;
import sun.security.util.DerOutputStream;
import sun.security.util.DerValue;

import java.io.IOException;
import java.security.interfaces.RSAPublicKey;

public class Pkcs1Formatter {

	// https://tools.ietf.org/html/rfc3447#appendix-A.1.1

	public byte[] toBytes(RSAPrivateCrtKeyImpl key) throws IOException {
		byte tag = 2;
		DerValue[] outputSequence = new DerValue[] {
			new DerValue(tag, new byte[]{0}),
			new DerValue(tag, key.getModulus().toByteArray()),
			new DerValue(tag, key.getPublicExponent().toByteArray()),
			new DerValue(tag, key.getPrivateExponent().toByteArray()),
			new DerValue(tag, key.getPrimeP().toByteArray()),
			new DerValue(tag, key.getPrimeQ().toByteArray()),
			new DerValue(tag, key.getPrimeExponentP().toByteArray()),
			new DerValue(tag, key.getPrimeExponentQ().toByteArray()),
			new DerValue(tag, key.getCrtCoefficient().toByteArray()),
		};

		DerOutputStream outputStream = new DerOutputStream();

		outputStream.putSequence(outputSequence);
		outputStream.flush();

		return outputStream.toByteArray();
	}

	public byte[] toBytes(RSAPublicKey key) throws IOException {
		byte tag = 2;
		DerValue[] outputSequence = new DerValue[] {
			new DerValue(tag, key.getModulus().toByteArray()),
			new DerValue(tag, key.getPublicExponent().toByteArray())
		};

		DerOutputStream outputStream = new DerOutputStream();

		outputStream.putSequence(outputSequence);
		outputStream.flush();

		return outputStream.toByteArray();
	}
}
