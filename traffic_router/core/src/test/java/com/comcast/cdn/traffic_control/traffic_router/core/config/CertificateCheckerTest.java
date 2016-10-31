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

package com.comcast.cdn.traffic_control.traffic_router.core.config;

import com.comcast.cdn.traffic_control.traffic_router.shared.Certificate;
import com.comcast.cdn.traffic_control.traffic_router.shared.CertificateData;
import org.json.JSONArray;
import org.json.JSONObject;
import org.junit.Before;
import org.junit.Test;

import java.util.Arrays;
import java.util.List;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.equalTo;

public class CertificateCheckerTest {

	private JSONObject deliveryServicesJson;
	private List<CertificateData> certificateDataList;

	@Before
	public void before() throws Exception {
		Certificate certificate = new Certificate();
		certificate.setCrt("the-crt");
		certificate.setKey("the-key");

		CertificateData certificateData = new CertificateData();
		certificateData.setHostname("https-delivery-service.thecdn.example.com");
		certificateData.setDeliveryservice("https-delivery-service");
		certificateData.setCertificate(certificate);

		certificateDataList = Arrays.asList(
			certificateData
		);

		JSONObject matchItem1 = new JSONObject().put("regex", ".*\\.https-delivery-service\\..*");
		JSONArray matchListArray1 = new JSONArray().put(0, matchItem1);
		JSONObject matchSetItem1 = new JSONObject()
			.put("protocol", "HTTP")
			.put("matchlist", matchListArray1);
		JSONArray domainsArray1 = new JSONArray().put(0, "*.https-delivery-service.thecdn.example.com");

		JSONArray matchsetsArray1 = new JSONArray();
		matchsetsArray1.put(0, matchSetItem1);

		JSONObject protocol1 = new JSONObject().put("acceptHttps", "true");

		JSONObject httpsDeliveryServiceJson = new JSONObject()
			.put("sslEnabled", "true")
			.put("protocol", protocol1)
			.put("matchsets", matchsetsArray1)
			.put("domains", domainsArray1);

		JSONObject matchItem2 = new JSONObject().put("regex", ".*\\.http-delivery-service\\..*");
		JSONArray matchListArray2 = new JSONArray().put(0, matchItem2);
		JSONObject matchSetItem2 = new JSONObject()
			.put("protocol", "HTTP")
			.put("matchlist", matchListArray2);
		JSONArray domainsArray2 = new JSONArray().put(0, "*.http-delivery-service.thecdn.example.com");

		JSONArray matchsetsArray2 = new JSONArray().put(0, matchSetItem2);
		JSONObject protocol2 = new JSONObject().put("acceptHttps", "false");

		JSONObject httpDeliveryServiceJson = new JSONObject()
			.put("sslEnabled", "false")
			.put("protocol", protocol2)
			.put("matchsets", matchsetsArray2)
			.put("domains", domainsArray2);


		JSONObject matchItem3 = new JSONObject().put("regex", ".*\\.dnssec-delivery-service\\..*");
		JSONArray matchListArray3 = new JSONArray().put(0, matchItem3);

		JSONObject matchSetItem3 = new JSONObject()
			.put("protocol", "DNS")
			.put("matchlist", matchListArray3);
		JSONArray domainsArray3 = new JSONArray().put(0, "*.dnssec-delivery-service.thecdn.example.com");

		JSONObject dnssecDeliveryServiceJson = new JSONObject()
			.put("sslEnabled", "true")
			.put("protocol", new JSONObject().put("acceptHttps", "true"))
			.put("matchsets", new JSONArray().put(0, matchSetItem3))
			.put("domains", domainsArray3);

		deliveryServicesJson = new JSONObject()
			.put("https-delivery-service", httpsDeliveryServiceJson)
			.put("http-delivery-service", httpDeliveryServiceJson)
			.put("dnssec-delivery-service", dnssecDeliveryServiceJson);
	}

	@Test
	public void itReturnsTrueWhenAllHttpsDeliveryServicesHaveCertificates() throws Exception {
		CertificateChecker certificateChecker = new CertificateChecker();

		assertThat(certificateChecker.certificatesAreValid(certificateDataList, deliveryServicesJson), equalTo(true));
	}

	@Test
	public void itReturnsFalseWhenAnyHttpsDeliveryServiceMissingCertificates() throws Exception {
		JSONObject matchItem = new JSONObject().put("regex", ".*\\.bad-https-delivery-service\\..*");
		JSONArray matchListArray = new JSONArray().put(0, matchItem);
		JSONObject matchSetItem = new JSONObject()
			.put("protocol", "HTTP")
			.put("matchlist", matchListArray);
		JSONArray domainsArray = new JSONArray().put(0, "*.bad-https-delivery-service.thecdn.example.com");

		JSONArray matchsetsArray = new JSONArray().put(0, matchSetItem);
		JSONObject protocol = new JSONObject().put("acceptHttps", "true");

		JSONObject httpsDeliveryServiceJson = new JSONObject()
			.put("sslEnabled", "true")
			.put("protocol", protocol)
			.put("matchsets", matchsetsArray)
			.put("domains", domainsArray);

		deliveryServicesJson.put("bad-https-delivery-service", httpsDeliveryServiceJson);

		assertThat(new CertificateChecker().certificatesAreValid(certificateDataList, deliveryServicesJson), equalTo(false));
	}
}
