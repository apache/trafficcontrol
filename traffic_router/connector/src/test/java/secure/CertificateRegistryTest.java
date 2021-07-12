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

import org.apache.traffic_control.traffic_router.protocol.RouterNioEndpoint;
import org.apache.traffic_control.traffic_router.secure.CertificateDataConverter;
import org.apache.traffic_control.traffic_router.secure.CertificateRegistry;
import org.apache.traffic_control.traffic_router.secure.HandshakeData;
import org.apache.traffic_control.traffic_router.shared.CertificateData;
import org.junit.Before;
import org.junit.Test;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.equalTo;
import static org.hamcrest.Matchers.containsInAnyOrder;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.times;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;
import static org.mockito.Mockito.anyMap;

public class CertificateRegistryTest {

	private CertificateRegistry certificateRegistry = CertificateRegistry.getInstance();
	private List certificateDataList;
	private CertificateDataConverter certificateDataConverter;
	private CertificateData certificateData1;
	private CertificateData certificateData2;
	private CertificateData certificateData3;
	private HandshakeData handshakeData1;
	private HandshakeData handshakeData2;
	private HandshakeData handshakeData3;

	@Before
	public void before() throws Exception {
		certificateData1 = mock(CertificateData.class);
		certificateData2 = mock(CertificateData.class);
		certificateData3 = mock(CertificateData.class);
		when(certificateData1.alias()).thenReturn("ds-1.some-cdn.example.com");
		when(certificateData2.alias()).thenReturn("ds-2.some-cdn.example.com");
		when(certificateData3.alias()).thenReturn("ds-3.some-cdn.example.com");

		certificateDataList = new ArrayList(Arrays.asList(certificateData1, certificateData2, certificateData3));

		handshakeData1 = mock(HandshakeData.class);
		when(handshakeData1.getHostname()).thenReturn("*.ds-1.some-cdn.example.com");

		handshakeData2 = mock(HandshakeData.class);
		when(handshakeData2.getHostname()).thenReturn("*.ds-2.some-cdn.example.com");

		handshakeData3 = mock(HandshakeData.class);
		when(handshakeData3.getHostname()).thenReturn("*.ds-3.some-cdn.example.com");

		certificateDataConverter = mock(CertificateDataConverter.class);

		when(certificateDataConverter.toHandshakeData(certificateData1)).thenReturn(handshakeData1);
		when(certificateDataConverter.toHandshakeData(certificateData2)).thenReturn(handshakeData2);
		when(certificateDataConverter.toHandshakeData(certificateData3)).thenReturn(handshakeData3);

		certificateRegistry.setCertificateDataConverter(certificateDataConverter);
	}

	@Test
	public void itImportsCertificateData() throws Exception {
		certificateRegistry.importCertificateDataList(certificateDataList);

		assertThat(certificateRegistry.getHandshakeData("ds-1.some-cdn.example.com"), equalTo(handshakeData1));
		assertThat(certificateRegistry.getHandshakeData("ds-2.some-cdn.example.com"), equalTo(handshakeData2));
		assertThat(certificateRegistry.getHandshakeData("ds-3.some-cdn.example.com"), equalTo(handshakeData3));

		verify(certificateDataConverter).toHandshakeData(certificateData1);
		verify(certificateDataConverter).toHandshakeData(certificateData2);
		verify(certificateDataConverter).toHandshakeData(certificateData3);

		assertThat(certificateRegistry.getAliases(),
			containsInAnyOrder(CertificateRegistry.DEFAULT_SSL_KEY, "ds-1.some-cdn.example.com",
					"ds-2.some-cdn.example.com", "ds-3.some-cdn.example.com"));
	}

	@Test
	public void itRetrysCertificateDataOnEndpointFail() throws Exception {
		HandshakeData handshakeData3mod = mock(HandshakeData.class);
		when(handshakeData3mod.getHostname()).thenReturn("*.ds-3.some-cdn.example.com");
		CertificateData certificateData3mod = mock(CertificateData.class);
		when(certificateData3mod.alias()).thenReturn("ds-3.some-cdn.example.com");
		when(certificateDataConverter.toHandshakeData(certificateData3mod)).thenReturn(handshakeData3mod);
		certificateDataList.remove(certificateData3);
		certificateDataList.add(certificateData3mod);
		RouterNioEndpoint endpoint = mock(RouterNioEndpoint.class);
		List<String> failist = new ArrayList<>();
		failist.add("ds-3.some-cdn.example.com");
		when(endpoint.reloadSSLHosts(anyMap())).thenReturn(failist);
		certificateRegistry.setEndPoint(endpoint);
		certificateRegistry.importCertificateDataList(certificateDataList);

		assertThat(certificateRegistry.getHandshakeData("ds-1.some-cdn.example.com"), equalTo(handshakeData1));
		assertThat(certificateRegistry.getHandshakeData("ds-2.some-cdn.example.com"), equalTo(handshakeData2));
		assertThat(certificateRegistry.getHandshakeData("ds-3.some-cdn.example.com"), equalTo(handshakeData3mod));

		verify(certificateDataConverter).toHandshakeData(certificateData1);
		verify(certificateDataConverter).toHandshakeData(certificateData2);
		verify(certificateDataConverter).toHandshakeData(certificateData3mod);
		verify(endpoint).reloadSSLHosts(anyMap());

		assertThat(certificateRegistry.getAliases(),
			containsInAnyOrder(CertificateRegistry.DEFAULT_SSL_KEY, "ds-1.some-cdn.example.com", "ds-2.some-cdn.example.com", "ds-3.some-cdn.example.com"));

		// try again
		// we should see that reloadSSLHosts gets called again even though none of the inputs have changed
		when(endpoint.reloadSSLHosts(anyMap())).thenReturn(new ArrayList<>());
		certificateRegistry.importCertificateDataList(certificateDataList);
		assertThat(certificateRegistry.getHandshakeData("ds-1.some-cdn.example.com"), equalTo(handshakeData1));
		assertThat(certificateRegistry.getHandshakeData("ds-2.some-cdn.example.com"), equalTo(handshakeData2));
		assertThat(certificateRegistry.getHandshakeData("ds-3.some-cdn.example.com"), equalTo(handshakeData3mod));

		verify(certificateDataConverter, times(2)).toHandshakeData(certificateData1);
		verify(certificateDataConverter, times(2)).toHandshakeData(certificateData2);
		verify(certificateDataConverter, times(2)).toHandshakeData(certificateData3mod);
		verify(endpoint, times(2)).reloadSSLHosts(anyMap());

		assertThat(certificateRegistry.getAliases(),
			containsInAnyOrder(CertificateRegistry.DEFAULT_SSL_KEY, "ds-1.some-cdn.example.com", "ds-2.some-cdn.example.com", "ds-3.some-cdn.example.com"));
	}
}
