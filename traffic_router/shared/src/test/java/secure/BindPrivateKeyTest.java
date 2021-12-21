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

package secure;

import org.apache.traffic_control.traffic_router.secure.BindPrivateKey;
import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.mockito.Mockito;
import org.powermock.api.mockito.PowerMockito;
import org.powermock.core.classloader.annotations.PowerMockIgnore;
import org.powermock.core.classloader.annotations.PrepareForTest;
import org.powermock.modules.junit4.PowerMockRunner;

import java.math.BigInteger;
import java.security.KeyFactory;
import java.security.KeyPair;
import java.security.KeyPairGenerator;
import java.security.PrivateKey;
import java.security.SecureRandom;
import java.security.interfaces.RSAPrivateCrtKey;
import java.security.spec.RSAPrivateCrtKeySpec;

import static java.util.Base64.getEncoder;
import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.equalTo;
import static org.mockito.Mockito.mock;
import static org.powermock.api.mockito.PowerMockito.doReturn;
import static org.powermock.api.mockito.PowerMockito.when;
import static org.powermock.api.mockito.PowerMockito.whenNew;

@RunWith(PowerMockRunner.class)
@PrepareForTest(BindPrivateKey.class)
@PowerMockIgnore("javax.management.*")
public class BindPrivateKeyTest {
	private String privateKeyString;
	private PrivateKey privateKey;

	String encode(BigInteger bigInteger) {
		return new String(getEncoder().encode(bigInteger.toByteArray()));
	}

	@Before
	public void before() throws Exception {
		KeyPairGenerator keyPairGenerator = KeyPairGenerator.getInstance("RSA");
		keyPairGenerator.initialize(2048, SecureRandom.getInstance("SHA1PRNG","SUN"));
		KeyPair keyPair = keyPairGenerator.generateKeyPair();

		RSAPrivateCrtKey privateCrtKey = (RSAPrivateCrtKey) keyPair.getPrivate();

		privateKeyString = "Private-key-format: v1.2\n" +
			"Algorithm: 5 (RSASHA1)\n" +
			"Modulus: " + encode(privateCrtKey.getModulus()) + "\n" +
			"PublicExponent: " + encode(privateCrtKey.getPublicExponent()) + "\n" +
			"PrivateExponent: " + encode(privateCrtKey.getPrivateExponent()) + "\n" +
			"Prime1: " + encode(privateCrtKey.getPrimeP()) + "\n" +
			"Prime2: " + encode(privateCrtKey.getPrimeQ()) + "\n" +
			"Exponent1: " + encode(privateCrtKey.getPrimeExponentP()) + "\n" +
			"Exponent2: " + encode(privateCrtKey.getPrimeExponentQ())+ "\n" +
			"Coefficient: " + encode(privateCrtKey.getCrtCoefficient())+ "\n";

		privateKey = Mockito.mock(PrivateKey.class);
		KeyFactory keyFactory = PowerMockito.mock(KeyFactory.class);

		PowerMockito.mockStatic(KeyFactory.class);
		when(KeyFactory.getInstance("RSA")).thenReturn(keyFactory);

		RSAPrivateCrtKeySpec spec = mock(RSAPrivateCrtKeySpec.class);

		whenNew(RSAPrivateCrtKeySpec.class)
			.withArguments(
				privateCrtKey.getModulus(),
				privateCrtKey.getPublicExponent(),
				privateCrtKey.getPrivateExponent(),
				privateCrtKey.getPrimeP(),
				privateCrtKey.getPrimeQ(),
				privateCrtKey.getPrimeExponentP(),
				privateCrtKey.getPrimeExponentQ(),
				privateCrtKey.getCrtCoefficient())
			.thenReturn(spec);

		doReturn(privateKey).when(keyFactory).generatePrivate(spec);
	}

	@Test
	public void itDecodesPrivateKeyString() {
		PrivateKey key = new BindPrivateKey().decode(privateKeyString);
		assertThat(key, equalTo(privateKey));
	}
}
