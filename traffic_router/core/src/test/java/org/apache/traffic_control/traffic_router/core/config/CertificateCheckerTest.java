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

package org.apache.traffic_control.traffic_router.core.config;

import org.apache.traffic_control.traffic_router.shared.Certificate;
import org.apache.traffic_control.traffic_router.shared.CertificateData;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.junit.Before;
import org.junit.Test;

import java.io.File;
import java.util.Arrays;
import java.util.List;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.equalTo;

public class CertificateCheckerTest {

	private JsonNode deliveryServicesJson;
	private List<CertificateData> certificateDataList;
	private CertificateData certificateData;

	@Before
	public void before() throws Exception {
		Certificate certificate = new Certificate();
		certificate.setCrt("the-crt");
		certificate.setKey("the-key");

		certificateData = new CertificateData();
		certificateData.setHostname("https-delivery-service.thecdn.example.com");
		certificateData.setDeliveryservice("https-delivery-service");
		certificateData.setCertificate(certificate);

		certificateDataList = Arrays.asList(
			certificateData
		);

	}

	@Test
	public void itReturnsFalseWhenDeliveryServiceNameIsNull() throws Exception {
		final File file = new File("src/test/resources/deliveryServices_missingDSName.json");
		final ObjectMapper mapper = new ObjectMapper();
		deliveryServicesJson = mapper.readTree(file);
		CertificateChecker certificateChecker = new CertificateChecker();
		certificateData.setDeliveryservice(null);

		assertThat(certificateChecker.certificatesAreValid(certificateDataList, deliveryServicesJson), equalTo(false));
	}

	@Test
	public void itReturnsFalseWhenDeliveryServiceNameIsBlank() throws Exception {
		final File file = new File("src/test/resources/deliveryServices_missingDSName.json");
		final ObjectMapper mapper = new ObjectMapper();
		deliveryServicesJson = mapper.readTree(file);
		CertificateChecker certificateChecker = new CertificateChecker();
		certificateData.setDeliveryservice("");

		assertThat(certificateChecker.certificatesAreValid(certificateDataList, deliveryServicesJson), equalTo(false));
	}

	@Test
	public void itReturnsTrueWhenAllHttpsDeliveryServicesHaveCertificates() throws Exception {
		final File file = new File("src/test/resources/deliveryServices.json");
		final ObjectMapper mapper = new ObjectMapper();
		deliveryServicesJson = mapper.readTree(file);
		CertificateChecker certificateChecker = new CertificateChecker();

		assertThat(certificateChecker.certificatesAreValid(certificateDataList, deliveryServicesJson), equalTo(true));
	}

	@Test
	public void itReturnsFalseWhenAnyHttpsDeliveryServiceMissingCertificates() throws Exception {

		final File file = new File("src/test/resources/deliveryServices_missingCert.json");
		final ObjectMapper mapper = new ObjectMapper();
		deliveryServicesJson = mapper.readTree(file);

		assertThat(new CertificateChecker().certificatesAreValid(certificateDataList, deliveryServicesJson), equalTo(false));
	}
}
