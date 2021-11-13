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

package org.apache.traffic_control.traffic_router.core.loc;

import org.apache.traffic_control.traffic_router.core.util.CidrAddress;
import org.junit.Test;

import static org.hamcrest.CoreMatchers.not;
import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.containsInAnyOrder;
import static org.hamcrest.Matchers.equalTo;
import static org.hamcrest.core.IsNull.nullValue;

public class FederationMappingBuilderTest {
    @Test
    public void itConsumesValidJSON() throws Exception {
        FederationMappingBuilder federationMappingBuilder = new FederationMappingBuilder();

        String json = "{ " +
            "\"cname\" : \"cname1\", " +
            "\"ttl\" : \"86400\", " +
            "\"resolve4\" : [ \"192.168.56.78/24\", \"192.168.45.67/24\" ], " +
            "\"resolve6\" : [ \"fdfe:dcba:9876:5::/64\", \"fd12:3456:789a:1::/64\" ] " +
        "}";

        FederationMapping federationMapping = federationMappingBuilder.fromJSON(json);

        assertThat(federationMapping, not(nullValue()));
        assertThat(federationMapping.getCname(), equalTo("cname1"));
        assertThat(federationMapping.getTtl(), equalTo(86400));

        assertThat(federationMapping.getResolve4(),
                containsInAnyOrder(CidrAddress.fromString("192.168.45.67/24"), CidrAddress.fromString("192.168.56.78/24")));
        assertThat(federationMapping.getResolve6(),
                containsInAnyOrder(CidrAddress.fromString("fd12:3456:789a:1::/64"), CidrAddress.fromString("fdfe:dcba:9876:5::/64")));
    }

    @Test
    public void itConsumesJSONWithoutResolvers() throws Exception {
        FederationMappingBuilder federationMappingBuilder = new FederationMappingBuilder();

        String json = "{ " +
                "\"cname\" : \"cname1\", " +
                "\"ttl\" : \"86400\" " +
                "}";

        FederationMapping federationMapping = federationMappingBuilder.fromJSON(json);

        assertThat(federationMapping, not(nullValue()));
        assertThat(federationMapping.getCname(), equalTo("cname1"));
        assertThat(federationMapping.getTtl(), equalTo(86400));

        assertThat(federationMapping.getResolve4(), not(nullValue()));
        assertThat(federationMapping.getResolve6(), not(nullValue()));
    }

}
