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

import org.apache.traffic_control.traffic_router.core.edge.InetRecord;
import org.apache.traffic_control.traffic_router.core.util.CidrAddress;
import org.apache.traffic_control.traffic_router.core.util.ComparableTreeSet;
import org.junit.Before;
import org.junit.Test;

import java.util.ArrayList;
import java.util.List;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.containsInAnyOrder;

public class FederationRegistryTest {

    private List<Federation> federations;

    @Before
    public void before() throws Exception {
        CidrAddress cidrAddress1 = CidrAddress.fromString("192.168.10.11/16");
        CidrAddress cidrAddress2 = CidrAddress.fromString("192.168.20.22/24");

        ComparableTreeSet<CidrAddress> cidrAddressesIpV4 = new ComparableTreeSet<CidrAddress>();
        cidrAddressesIpV4.add(cidrAddress1);
        cidrAddressesIpV4.add(cidrAddress2);

        CidrAddress cidrAddress3 = CidrAddress.fromString("fdfe:dcba:9876:5::/64");
        CidrAddress cidrAddress4 = CidrAddress.fromString("fd12:3456:789a:1::/64");

        ComparableTreeSet<CidrAddress> cidrAddressesIpV6 = new ComparableTreeSet<CidrAddress>();
        cidrAddressesIpV6.add(cidrAddress3);
        cidrAddressesIpV6.add(cidrAddress4);

        FederationMapping federationMapping = new FederationMapping("cname1", 1234, cidrAddressesIpV4, cidrAddressesIpV6);
        List<FederationMapping> federationMappings = new ArrayList<FederationMapping>();
        federationMappings.add(federationMapping);

        Federation federation = new Federation("kable-town-01", federationMappings);

        federations = new ArrayList<Federation>();
        federations.add(federation);
    }


    @Test
    public void itFindsMapping() throws Exception {
        FederationRegistry federationRegistry = new FederationRegistry();
        federationRegistry.setFederations(federations);

        List<InetRecord> inetRecords = federationRegistry.findInetRecords("kable-town-01", CidrAddress.fromString("192.168.10.11/24"));
        assertThat(inetRecords, containsInAnyOrder(new InetRecord("cname1", 1234)));

        inetRecords = federationRegistry.findInetRecords("kable-town-01", CidrAddress.fromString("192.168.10.11/16"));
        assertThat(inetRecords, containsInAnyOrder(new InetRecord("cname1", 1234)));
    }
}