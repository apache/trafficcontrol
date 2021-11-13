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

package org.apache.traffic_control.traffic_router.core.external;

import org.apache.traffic_control.traffic_router.core.util.ExternalTest;
import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.apache.catalina.LifecycleException;
import org.apache.http.client.methods.CloseableHttpResponse;
import org.apache.http.client.methods.HttpGet;
import org.apache.http.impl.client.CloseableHttpClient;
import org.apache.http.impl.client.HttpClientBuilder;
import org.apache.http.util.EntityUtils;
import org.junit.After;
import org.junit.Before;
import org.junit.Test;
import org.junit.experimental.categories.Category;

import java.util.HashMap;
import java.util.Map;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.containsInAnyOrder;

@Category(ExternalTest.class)
public class ZonesTest {
	CloseableHttpClient httpClient;

	@Before
	public void before() throws LifecycleException {
		httpClient = HttpClientBuilder.create().build();
	}

	@After
	public void after() throws Exception {
		if (httpClient != null) httpClient.close();
	}

	@Test
	public void itGetsStatsForZones() throws Exception {
		HttpGet httpGet = new HttpGet("http://localhost:3333/crs/stats/zones/caches");
		CloseableHttpResponse response = null;

		try {
			response = httpClient.execute(httpGet);
			String actual = EntityUtils.toString(response.getEntity());

			Map<String, Object> zoneStats = new ObjectMapper().readValue(actual, new TypeReference<HashMap<String, Object>>() { });

			Map<String, Object> dynamicZonesStats = (Map<String, Object>) zoneStats.get("dynamicZoneCaches");
			assertThat(dynamicZonesStats.keySet(), containsInAnyOrder("requestCount", "evictionCount", "totalLoadTime",
				"averageLoadPenalty", "hitCount", "loadSuccessCount", "missRate", "loadExceptionRate", "hitRate", "missCount", "loadCount", "loadExceptionCount"));

			Map<String, Object> staticZonesStats = (Map<String, Object>) zoneStats.get("staticZoneCaches");
			assertThat(staticZonesStats.keySet(), containsInAnyOrder("requestCount", "evictionCount", "totalLoadTime",
				"averageLoadPenalty", "hitCount", "loadSuccessCount", "missRate", "loadExceptionRate", "hitRate", "missCount", "loadCount", "loadExceptionCount"));

		} finally {
			if (response != null) response.close();
		}
	}
}
