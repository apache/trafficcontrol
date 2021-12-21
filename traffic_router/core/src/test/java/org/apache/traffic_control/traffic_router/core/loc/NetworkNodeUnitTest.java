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

import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.junit.Test;

import java.util.HashSet;
import java.util.Iterator;
import java.util.Map;
import java.util.Set;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.equalTo;
import static org.hamcrest.Matchers.not;

public class NetworkNodeUnitTest {
    
    @Test
    public void itSupportsARootNode() throws Exception {
        NetworkNode root = new NetworkNode("0.0.0.0/0");
        NetworkNode network = new NetworkNode("192.168.1.0/24");
        assertThat(root.add(network), equalTo(true));
        assertThat(root.children.entrySet().iterator().next().getKey(), equalTo(network));
    }

    @Test
    public void itDoesNotAddANodeOutsideOfNetwork() throws Exception {
        NetworkNode network = new NetworkNode("192.168.0.0/16");
        NetworkNode subnetwork = new NetworkNode("10.10.0.0/16");

        assertThat(network.add(subnetwork), equalTo(false));
    }

    @Test
    public void itFindsIpBelongingToNetwork() throws Exception {
        NetworkNode network = new NetworkNode("192.168.1.0/24");
        assertThat(network.getNetwork("192.168.1.1"), equalTo(network));
        assertThat(network.getNetwork("192.168.2.1"), not(equalTo(network)));
    }

    @Test
    public void itDoesNotAddDuplicates() throws Exception {
        NetworkNode supernet = new NetworkNode("192.168.0.0/16");
        NetworkNode network1 = new NetworkNode("192.168.1.0/24");
        NetworkNode duplicate = new NetworkNode("192.168.1.0/24");

        assertThat(supernet.add(network1), equalTo(true));
        assertThat(supernet.children.size(), equalTo(1));

        assertThat(supernet.add(duplicate), equalTo(false));
        assertThat(supernet.children.size(), equalTo(1));
    }

    @Test
    public void itPutsNetworksIntoOrderedHierarchy() throws Exception {
        NetworkNode root = new NetworkNode("0.0.0.0/0");

        NetworkNode subnet1 = new NetworkNode("192.168.6.0/24");
        NetworkNode subnet2 = new NetworkNode("192.168.55.0/24");

        NetworkNode net = new NetworkNode("192.168.0.0/16");

        root.add(net);

        assertThat(root.children.entrySet().iterator().next().getKey(),equalTo(net));

        root.add(subnet2);
        root.add(subnet1);

        final Iterator<Map.Entry<NetworkNode, NetworkNode>> iterator = net.children.entrySet().iterator();
        assertThat(iterator.next().getKey(),equalTo(subnet1));
        assertThat(iterator.next().getKey(),equalTo(subnet2));
    }

    @Test
    public void itSupportsDeepCaches() throws Exception {
        String czmapString = "{" +
                "\"revision\": \"Mon Dec 21 15:04:01 2015\"," +
                "\"customerName\": \"Kabletown\"," +
                "\"deepCoverageZones\": {" +
                "\"us-co-denver\": {" +
                "\"network\": [\"192.168.55.0/24\",\"192.168.6.0/24\",\"192.168.0.0/16\"]," +
                "\"network6\": [\"1234:5678::/64\",\"1234:5679::/64\"]," +
                "\"caches\": [\"host1\",\"host2\"]" +
                "}" +
                "}" +
                "}";

        final ObjectMapper mapper = new ObjectMapper();
        final JsonNode json = mapper.readTree(czmapString);
        NetworkNode networkNode = NetworkNode.generateTree(json, false, true);
        NetworkNode foundNetworkNode = networkNode.getNetwork("192.168.55.100");

        Set<String> expected = new HashSet<String>();
        expected.add("host1");
        expected.add("host2");
        assertThat(foundNetworkNode.getDeepCacheNames(), equalTo(expected));
    }

    @Test
    public void itDoesIpV6() throws Exception {
        String czmapString = "{" +
            "\"revision\": \"Mon Dec 21 15:04:01 2015\"," +
            "\"customerName\": \"Kabletown\"," +
            "\"coverageZones\": {" +
            "\"us-co-denver\": {" +
            "\"network\": [\"192.168.55.0/24\",\"192.168.6.0/24\",\"192.168.0.0/16\"]," +
            "\"network6\": [\"1234:5678::/64\",\"1234:5679::/64\"]" +
            "}" +
            "}" +
            "}";

        final ObjectMapper mapper = new ObjectMapper();
        final JsonNode json = mapper.readTree(czmapString);
        NetworkNode networkNode = NetworkNode.generateTree(json, false);
        NetworkNode foundNetworkNode = networkNode.getNetwork("1234:5678::1");

        assertThat(foundNetworkNode.getLoc(), equalTo("us-co-denver"));
    }

    @Test
    public void itPutsAllSubnetsUnderSuperNet() throws Exception {
        NetworkNode root = new NetworkNode("0.0.0.0/0");

        NetworkNode subnet1 = new NetworkNode("192.168.6.0/24");
        root.add(subnet1);

        NetworkNode subnet2 = new NetworkNode("192.168.55.0/24");
        root.add(subnet2);

        NetworkNode net = new NetworkNode("192.168.0.0/16");
        root.add(net);

        assertThat(root.children.isEmpty(), equalTo(false));
        NetworkNode generation1Node = root.children.values().iterator().next();
        assertThat(generation1Node.toString(), equalTo("[192.168.0.0/16] - location:null"));

        final Iterator<Map.Entry<NetworkNode, NetworkNode>> iterator = generation1Node.children.entrySet().iterator();
        NetworkNode generation2FirstNode = iterator.next().getKey();
        NetworkNode generation2SecondNode = iterator.next().getKey();


        assertThat(generation2FirstNode.toString(), equalTo("[192.168.6.0/24] - location:null"));
        assertThat(generation2SecondNode.toString(), equalTo("[192.168.55.0/24] - location:null"));
    }

    @Test
    public void itMatchesIpsInOverlappingSubnets() throws Exception {
        String czmapString = "{" +
            "\"revision\": \"Mon Dec 21 15:04:01 2015\"," +
            "\"customerName\": \"Kabletown\"," +
            "\"coverageZones\": {" +
            "\"us-co-denver\": {" +
            "\"network\": [\"192.168.55.0/24\",\"192.168.6.0/24\",\"192.168.0.0/16\"]," +
            "\"network6\": [\"0:0:0:0:0:ffff:a4f:3700/24\"]" +
            "}" +
            "}" +
            "}";

        final ObjectMapper mapper = new ObjectMapper();
        final JsonNode json = mapper.readTree(czmapString);
        NetworkNode networkNode = NetworkNode.generateTree(json, false);
        NetworkNode foundNetworkNode = networkNode.getNetwork("192.168.55.2");

        assertThat(foundNetworkNode.getLoc(), equalTo("us-co-denver"));
    }

    @Test
    public void itRejectsInvalidIpV4Network() throws Exception {
        String czmapString = "{" +
            "\"revision\": \"Mon Dec 21 15:04:01 2015\"," +
            "\"customerName\": \"Kabletown\"," +
            "\"coverageZones\": {" +
            "\"us-co-denver\": {" +
            "\"network\": [\"192.168.55.0/40\",\"192.168.6.0/24\",\"192.168.0.0/16\"]," +
            "\"network6\": [\"1234:5678::/64\",\"1234:5679::/64\"]" +
            "}" +
            "}" +
            "}";

        final ObjectMapper mapper = new ObjectMapper();
        final JsonNode json = mapper.readTree(czmapString);
        assertThat(NetworkNode.generateTree(json, false), equalTo(null));
    }

    @Test
    public void itRejectsInvalidIpV6Network() throws Exception {
        String czmapString = "{" +
            "\"revision\": \"Mon Dec 21 15:04:01 2015\"," +
            "\"customerName\": \"Kabletown\"," +
            "\"coverageZones\": {" +
            "\"us-co-denver\": {" +
            "\"network\": [\"192.168.55.0/24\",\"192.168.6.0/24\",\"192.168.0.0/16\"]," +
            "\"network6\": [\"1234:5678::/64\",\"zyx:5679::/64\"]" +
            "}" +
            "}" +
            "}";

        final ObjectMapper mapper = new ObjectMapper();
        final JsonNode json = mapper.readTree(czmapString);
        assertThat(NetworkNode.generateTree(json, false), equalTo(null));
    }
}