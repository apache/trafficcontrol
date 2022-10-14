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

package utils;

import org.apache.traffic_control.traffic_router.utils.HttpsProperties;
import org.junit.Test;
import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.equalTo;
import java.util.Map;

public class HttpsPropertiesTest {
    @Test
    public void checkGetHttpsProperties() throws Exception {
        final String fileName = "src/test/java/conf/https.properties";
        HttpsProperties httpsProperties = new HttpsProperties(fileName);
        Map<String, String> propsMap = httpsProperties.getHttpsPropertiesMap();
        assertThat(propsMap.get("https.certificate.location"), equalTo("/opt/traffic_router/conf/keyStore.jks"));
        assertThat(propsMap.get("https.password"), equalTo("changeit"));
        assertThat(propsMap.get("https.key.size"), equalTo("1024"));
        assertThat(propsMap.get("https.signature.algorithm"), equalTo("TestAlgorithm"));
        assertThat(propsMap.get("https.validity.years"), equalTo("TestValidity"));
        assertThat(propsMap.get("https.certificate.country"), equalTo("TestCountry"));
        assertThat(propsMap.get("https.certificate.state"), equalTo("TestState"));
        assertThat(propsMap.get("https.certificate.locality"), equalTo("TestLocality"));
        assertThat(propsMap.get("https.certificate.organization"), equalTo("TestOrg"));
        assertThat(propsMap.get("https.certificate.organizational.unit"), equalTo("; OU=Test Org Unit; OU= Test Org Unit 2"));
    }
}
