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

import java.math.BigInteger;
import java.security.KeyFactory;
import java.security.PrivateKey;
import java.security.spec.RSAPrivateCrtKeySpec;
import java.util.Arrays;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

import static java.util.Base64.getDecoder;

public class BindPrivateKey {
	private static final Logger LOGGER = LogManager.getLogger(BindPrivateKey.class);

	private BigInteger decodeBigInt(final String s) {
		return new BigInteger(1, getDecoder().decode(s.getBytes()));
	}

	private Map<String, BigInteger> decodeBigIntegers(final String s) {

		final List<String> bigIntKeys = Arrays.asList(
			"Modulus", "PublicExponent", "PrivateExponent", "Prime1", "Prime2", "Exponent1", "Exponent2", "Coefficient"
		);

		final Map<String, BigInteger>  bigIntegerMap = new HashMap<>();

		for (final String line : s.split("\n")) {
			final String[] tokens = line.split(": ");

			if (bigIntKeys.stream().filter(k -> k.equals(tokens[0])).findFirst().isPresent()) {
				bigIntegerMap.put(tokens[0], decodeBigInt(tokens[1]));
			}
		}

		return bigIntegerMap;
	}

	public PrivateKey decode(final String data) {
		final Map<String, BigInteger> map = decodeBigIntegers(data);
		final BigInteger modulus = map.get("Modulus");
		final BigInteger publicExponent = map.get("PublicExponent");
		final BigInteger privateExponent = map.get("PrivateExponent");
		final BigInteger prime1 = map.get("Prime1");
		final BigInteger prime2 = map.get("Prime2");
		final BigInteger exp1 = map.get("Exponent1");
		final BigInteger exp2 = map.get("Exponent2");
		final BigInteger coeff = map.get("Coefficient");

		final RSAPrivateCrtKeySpec keySpec = new RSAPrivateCrtKeySpec(modulus,publicExponent,privateExponent,prime1,prime2,exp1,exp2,coeff);

		try {
			return KeyFactory.getInstance("RSA").generatePrivate(keySpec);
		} catch (Exception e) {
			LOGGER.error("Failed to decode Bind Private Key data: " + e.getMessage(), e);
		}

		return null;
	}
}
