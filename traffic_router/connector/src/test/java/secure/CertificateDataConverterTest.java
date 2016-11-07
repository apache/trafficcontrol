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

import com.comcast.cdn.traffic_control.traffic_router.secure.CertificateDataConverter;
import com.comcast.cdn.traffic_control.traffic_router.secure.CertificateDecoder;
import com.comcast.cdn.traffic_control.traffic_router.secure.HandshakeData;
import com.comcast.cdn.traffic_control.traffic_router.secure.PrivateKeyDecoder;
import com.comcast.cdn.traffic_control.traffic_router.shared.Certificate;
import com.comcast.cdn.traffic_control.traffic_router.shared.CertificateData;
import org.junit.Before;
import org.junit.Test;

import java.security.PrivateKey;
import java.security.cert.X509Certificate;
import java.util.Arrays;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.equalTo;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;

public class CertificateDataConverterTest {

	private CertificateDataConverter certificateDataConverter;
	private CertificateData certificateData;
	private X509Certificate x509Certificate1;
	private X509Certificate x509Certificate2;
	private X509Certificate x509Certificate3;
	private PrivateKey privateKey;

	@Before
	public void before() throws Exception {
		PrivateKeyDecoder privateKeyDecoder = mock(PrivateKeyDecoder.class);
		CertificateDecoder certificateDecoder = mock(CertificateDecoder.class);

		Certificate certificate = new Certificate();
		certificate.setCrt("encodedchaindata");
		certificate.setKey("encodedkeydata");

		certificateData = new CertificateData();
		certificateData.setCertificate(certificate);
		certificateData.setDeliveryservice("some-delivery-service");
		certificateData.setHostname("example.com");

		privateKey = mock(PrivateKey.class);
		when(privateKeyDecoder.decode("encodedkeydata")).thenReturn(privateKey);

		when(certificateDecoder.doubleDecode("encodedchaindata")).thenReturn(Arrays.asList(
			"encodedcert1", "encodedcert2", "encodedcert3"
		));

		x509Certificate1 = mock(X509Certificate.class);
		x509Certificate2 = mock(X509Certificate.class);
		x509Certificate3 = mock(X509Certificate.class);

		when(certificateDecoder.toCertificate("encodedcert1")).thenReturn(x509Certificate1);
		when(certificateDecoder.toCertificate("encodedcert2")).thenReturn(x509Certificate2);
		when(certificateDecoder.toCertificate("encodedcert3")).thenReturn(x509Certificate3);

		certificateDataConverter = new CertificateDataConverter();
		certificateDataConverter.setCertificateDecoder(certificateDecoder);
		certificateDataConverter.setPrivateKeyDecoder(privateKeyDecoder);
	}

	@Test
	public void itConvertsToHandshakeData() throws Exception {
		HandshakeData handshakeData = certificateDataConverter.toHandshakeData(certificateData);

		assertThat(handshakeData.getDeliveryService(), equalTo("some-delivery-service"));
		assertThat(handshakeData.getHostname(), equalTo("example.com"));
		assertThat(handshakeData.getPrivateKey(), equalTo(privateKey));
		assertThat(handshakeData.getCertificateChain(), equalTo(new X509Certificate[]{x509Certificate1, x509Certificate2, x509Certificate3}));
	}
}
