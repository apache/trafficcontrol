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

package org.apache.traffic_control.traffic_router.core.dns;

import org.apache.traffic_control.traffic_router.core.ds.DeliveryService;
import org.apache.traffic_control.traffic_router.core.router.TrafficRouter;
import org.apache.traffic_control.traffic_router.core.router.TrafficRouterManager;
import org.apache.traffic_control.traffic_router.core.edge.CacheRegister;

import org.apache.traffic_control.traffic_router.core.util.JsonUtils;
import com.fasterxml.jackson.databind.node.ArrayNode;
import com.fasterxml.jackson.databind.node.ObjectNode;
import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.powermock.core.classloader.annotations.PowerMockIgnore;
import org.powermock.core.classloader.annotations.PrepareForTest;
import org.powermock.modules.junit4.PowerMockRunner;
import org.xbill.DNS.*;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.equalTo;

import java.net.Inet4Address;
import java.net.InetAddress;

import static org.mockito.ArgumentMatchers.any;
import static org.mockito.ArgumentMatchers.eq;
import static org.mockito.Mockito.doReturn;
import static org.mockito.Mockito.mock;
import static org.powermock.api.mockito.PowerMockito.*;

import java.util.*;

import com.fasterxml.jackson.databind.node.JsonNodeFactory;
import com.fasterxml.jackson.databind.JsonNode;

@RunWith(PowerMockRunner.class)
@PrepareForTest({Header.class, NameServer.class, TrafficRouterManager.class, TrafficRouter.class, CacheRegister.class})
@PowerMockIgnore("javax.management.*")
public class NameServerTest {
    private NameServer nameServer;
    private InetAddress client;
    private TrafficRouterManager trafficRouterManager;
    private TrafficRouter trafficRouter;
    private Record ar;
    private NSRecord ns;
    
    @Before
    public void before() throws Exception {

        client = Inet4Address.getByAddress(new byte[]{(byte) 192, (byte) 168, 23, 45});
        nameServer = new NameServer();	
        trafficRouterManager = mock(TrafficRouterManager.class);
        trafficRouter = mock(TrafficRouter.class);
        CacheRegister cacheRegister = mock(CacheRegister.class);
        doReturn(cacheRegister).when(trafficRouter).getCacheRegister();
        JsonNode js = JsonNodeFactory.instance.objectNode().put("ecsEnable", true);
        when(cacheRegister.getConfig()).thenReturn(js);
        
        Name m_an, m_host, m_admin;
	    m_an = Name.fromString("dns1.example.com.");
	    m_host = Name.fromString("dns1.example.com.");
	    m_admin = Name.fromString("admin.example.com.");
	    ar = new SOARecord(m_an, DClass.IN, 0x13A8,
	    		m_host, m_admin, 0xABCDEF12L, 0xCDEF1234L,
	    		0xEF123456L, 0x12345678L, 0x3456789AL);

	    ns = new NSRecord(m_an, DClass.IN, 12345L, m_an);
    }

    @Test
    public void TestARecordQueryWithClientSubnetOption() throws Exception {
        
        Name name = Name.fromString("host1.example.com.");
        Record question = Record.newRecord(name, Type.A, DClass.IN, 12345L);
        Message query = Message.newQuery(question);
       
        //Add opt record, with client subnet option.
        int nmask = 28;
        InetAddress ipaddr = Inet4Address.getByName("192.168.33.0");
        ClientSubnetOption cso = new ClientSubnetOption(nmask, ipaddr);
        List<ClientSubnetOption> cso_list = new ArrayList<ClientSubnetOption>(1);
        cso_list.add(cso);	
        OPTRecord opt = new OPTRecord(1280, 0, 0, 0, cso_list);
        query.addRecord(opt, Section.ADDITIONAL);
        
       
	    // Add ARecord Entry in the zone
        InetAddress resolvedAddress = Inet4Address.getByName("192.168.8.9");
        Record answer = new ARecord(name, DClass.IN, 12345L, resolvedAddress);
        Record[] records = new Record[] {ar, ns, answer};
        
        Name m_an = Name.fromString("dns1.example.com.");
        Zone zone = new Zone(m_an, records);
       
        DNSAccessRecord.Builder builder = new DNSAccessRecord.Builder(1L, client);

        nameServer.setTrafficRouterManager(trafficRouterManager);
        nameServer.setEcsEnable(JsonUtils.optBoolean(trafficRouter.getCacheRegister().getConfig(), "ecsEnable", false)); // this mimics what happens in ConfigHandler

        // Following is needed to mock this call: zone = trafficRouterManager.getTrafficRouter().getZone(qname, qtype, clientAddress, dnssecRequest, builder);
        when(trafficRouterManager.getTrafficRouter()).thenReturn(trafficRouter);
        when(trafficRouter.getZone(any(Name.class), any(int.class), eq(ipaddr), any(boolean.class), any(DNSAccessRecord.Builder.class))).thenReturn(zone);

        // The function call under test:
        Message res = nameServer.query(query, client, builder);

        
        //Verification of response
        OPTRecord qopt = res.getOPT();
        assert (qopt != null);
        List<EDNSOption> list = Collections.EMPTY_LIST;
        list = qopt.getOptions(EDNSOption.Code.CLIENT_SUBNET);
        assert (list != Collections.EMPTY_LIST);
        ClientSubnetOption option = (ClientSubnetOption)list.get(0);
        assertThat(nmask, equalTo(option.getSourceNetmask()));
        assertThat(nmask, equalTo(option.getScopeNetmask()));
        assertThat(ipaddr, equalTo(option.getAddress()));
        nameServer.setEcsEnable(false);
    }
    
    @Test
    public void TestARecordQueryWithMultipleClientSubnetOption() throws Exception {
        
        Name name = Name.fromString("host1.example.com.");
        Record question = Record.newRecord(name, Type.A, DClass.IN, 12345L);
        Message query = Message.newQuery(question);
       
        //Add opt record, with multiple client subnet option.
        int nmask1 = 16;
        int nmask2 = 24;
        InetAddress ipaddr1 = Inet4Address.getByName("192.168.0.0");
        InetAddress ipaddr2 = Inet4Address.getByName("192.168.33.0");
        ClientSubnetOption cso1 = new ClientSubnetOption(nmask1, ipaddr1);
        ClientSubnetOption cso2 = new ClientSubnetOption(nmask2, ipaddr2);
        List<ClientSubnetOption> cso_list = new ArrayList<ClientSubnetOption>(1);
        cso_list.add(cso1);
        cso_list.add(cso2);
        final OPTRecord opt = new OPTRecord(1280, 0, 0, 0, cso_list);
        query.addRecord(opt, Section.ADDITIONAL);
        
       
	    // Add ARecord Entry in the zone
        InetAddress resolvedAddress = Inet4Address.getByName("192.168.8.9");
        Record answer = new ARecord(name, DClass.IN, 12345L, resolvedAddress);
        Record[] records = new Record[] {ar, ns, answer};
        
        Name m_an = Name.fromString("dns1.example.com.");
        Zone zone = new Zone(m_an, records);
       
        DNSAccessRecord.Builder builder = new DNSAccessRecord.Builder(1L, client);

        nameServer.setTrafficRouterManager(trafficRouterManager);
        nameServer.setEcsEnable(JsonUtils.optBoolean(trafficRouter.getCacheRegister().getConfig(), "ecsEnable", false)); // this mimics what happens in ConfigHandler
	
        // Following is needed to mock this call: zone = trafficRouterManager.getTrafficRouter().getZone(qname, qtype, clientAddress, dnssecRequest, builder);
        when(trafficRouterManager.getTrafficRouter()).thenReturn(trafficRouter);
        when(trafficRouter.getZone(any(Name.class), any(int.class), eq(ipaddr2), any(boolean.class), any(DNSAccessRecord.Builder.class))).thenReturn(zone);

        // The function call under test:
        Message res = nameServer.query(query, client, builder);

        
        //Verification of response
        OPTRecord qopt = res.getOPT();
        assert (qopt != null);
        List<EDNSOption> list = Collections.EMPTY_LIST;
        list = qopt.getOptions(EDNSOption.Code.CLIENT_SUBNET);
        assert (list != Collections.EMPTY_LIST);
        ClientSubnetOption option = (ClientSubnetOption)list.get(0);
        assertThat(1, equalTo(list.size()));
        assertThat(nmask2, equalTo(option.getSourceNetmask()));
        assertThat(nmask2, equalTo(option.getScopeNetmask()));
        assertThat(ipaddr2, equalTo(option.getAddress()));
        nameServer.setEcsEnable(false);
    }

    @Test
    public void TestDeliveryServiceARecordQueryWithClientSubnetOption() throws Exception {

        CacheRegister cacheRegister = mock(CacheRegister.class);
        doReturn(cacheRegister).when(trafficRouter).getCacheRegister();
        JsonNode js = JsonNodeFactory.instance.objectNode().put("ecsEnable", false);
        when(cacheRegister.getConfig()).thenReturn(js);

        ObjectNode node = JsonNodeFactory.instance.objectNode();
        ArrayNode domainNode = node.putArray("domains");
        domainNode.add("example.com");
        node.put("routingName","edge");
        node.put("coverageZoneOnly", false);
        DeliveryService ds1 = new DeliveryService("ds1", node);
        Set dses = new HashSet();
        dses.add(ds1);
        nameServer.setEcsEnabledDses(dses);

        Name name = Name.fromString("host1.example.com.");
        Record question = Record.newRecord(name, Type.A, DClass.IN, 12345L);
        Message query = Message.newQuery(question);

        //Add opt record, with client subnet option.
        int nmask = 28;
        InetAddress ipaddr = Inet4Address.getByName("192.168.33.0");
        ClientSubnetOption cso = new ClientSubnetOption(nmask, ipaddr);
        List<ClientSubnetOption> cso_list = new ArrayList<ClientSubnetOption>(1);
        cso_list.add(cso);
        OPTRecord opt = new OPTRecord(1280, 0, 0, 0, cso_list);
        query.addRecord(opt, Section.ADDITIONAL);


        // Add ARecord Entry in the zone
        InetAddress resolvedAddress = Inet4Address.getByName("192.168.8.9");
        Record answer = new ARecord(name, DClass.IN, 12345L, resolvedAddress);
        Record[] records = new Record[] {ar, ns, answer};

        Name m_an = Name.fromString("dns1.example.com.");
        Zone zone = new Zone(m_an, records);

        DNSAccessRecord.Builder builder = new DNSAccessRecord.Builder(1L, client);

        nameServer.setTrafficRouterManager(trafficRouterManager);
        nameServer.setEcsEnable(JsonUtils.optBoolean(trafficRouter.getCacheRegister().getConfig(), "ecsEnable", false)); // this mimics what happens in ConfigHandler

        // Following is needed to mock this call: zone = trafficRouterManager.getTrafficRouter().getZone(qname, qtype, clientAddress, dnssecRequest, builder);
        when(trafficRouterManager.getTrafficRouter()).thenReturn(trafficRouter);
        when(trafficRouter.getZone(any(Name.class), any(int.class), eq(ipaddr), any(boolean.class), any(DNSAccessRecord.Builder.class))).thenReturn(zone);

        // The function call under test:
        Message res = nameServer.query(query, client, builder);


        //Verification of response
        OPTRecord qopt = res.getOPT();
        assert (qopt != null);
        List<EDNSOption> list = Collections.EMPTY_LIST;
        list = qopt.getOptions(EDNSOption.Code.CLIENT_SUBNET);
        assert (list != Collections.EMPTY_LIST);
        ClientSubnetOption option = (ClientSubnetOption)list.get(0);
        assertThat(nmask, equalTo(option.getSourceNetmask()));
        assertThat(nmask, equalTo(option.getScopeNetmask()));
        assertThat(ipaddr, equalTo(option.getAddress()));
        nameServer.setEcsEnable(false);
    }

    @Test
    public void TestDeliveryServiceARecordQueryWithMultipleClientSubnetOption() throws Exception {

        CacheRegister cacheRegister = mock(CacheRegister.class);
        doReturn(cacheRegister).when(trafficRouter).getCacheRegister();
        JsonNode js = JsonNodeFactory.instance.objectNode().put("ecsEnable", false);
        when(cacheRegister.getConfig()).thenReturn(js);

        Name name = Name.fromString("host1.example.com.");
        Record question = Record.newRecord(name, Type.A, DClass.IN, 12345L);
        Message query = Message.newQuery(question);

        ObjectNode node = JsonNodeFactory.instance.objectNode();
        ArrayNode domainNode = node.putArray("domains");
        domainNode.add("example.com");
        node.put("routingName","edge");
        node.put("coverageZoneOnly", false);
        DeliveryService ds1 = new DeliveryService("ds1", node);
        Set dses = new HashSet();
        dses.add(ds1);
        nameServer.setEcsEnabledDses(dses);


        //Add opt record, with multiple client subnet option.
        int nmask1 = 16;
        int nmask2 = 24;
        InetAddress ipaddr1 = Inet4Address.getByName("192.168.0.0");
        InetAddress ipaddr2 = Inet4Address.getByName("192.168.33.0");
        ClientSubnetOption cso1 = new ClientSubnetOption(nmask1, ipaddr1);
        ClientSubnetOption cso2 = new ClientSubnetOption(nmask2, ipaddr2);
        List<ClientSubnetOption> cso_list = new ArrayList<ClientSubnetOption>(1);
        cso_list.add(cso1);
        cso_list.add(cso2);
        final OPTRecord opt = new OPTRecord(1280, 0, 0, 0, cso_list);
        query.addRecord(opt, Section.ADDITIONAL);


        // Add ARecord Entry in the zone
        InetAddress resolvedAddress = Inet4Address.getByName("192.168.8.9");
        Record answer = new ARecord(name, DClass.IN, 12345L, resolvedAddress);
        Record[] records = new Record[] {ar, ns, answer};

        Name m_an = Name.fromString("dns1.example.com.");
        Zone zone = new Zone(m_an, records);

        DNSAccessRecord.Builder builder = new DNSAccessRecord.Builder(1L, client);

        nameServer.setTrafficRouterManager(trafficRouterManager);
        nameServer.setEcsEnable(JsonUtils.optBoolean(trafficRouter.getCacheRegister().getConfig(), "ecsEnable", false)); // this mimics what happens in ConfigHandler

        // Following is needed to mock this call: zone = trafficRouterManager.getTrafficRouter().getZone(qname, qtype, clientAddress, dnssecRequest, builder);
        when(trafficRouterManager.getTrafficRouter()).thenReturn(trafficRouter);
        when(trafficRouter.getZone(any(Name.class), any(int.class), eq(ipaddr2), any(boolean.class), any(DNSAccessRecord.Builder.class))).thenReturn(zone);

        // The function call under test:
        Message res = nameServer.query(query, client, builder);


        //Verification of response
        OPTRecord qopt = res.getOPT();
        assert (qopt != null);
        List<EDNSOption> list = Collections.EMPTY_LIST;
        list = qopt.getOptions(EDNSOption.Code.CLIENT_SUBNET);
        assert (list != Collections.EMPTY_LIST);
        ClientSubnetOption option = (ClientSubnetOption)list.get(0);
        assertThat(1, equalTo(list.size()));
        assertThat(nmask2, equalTo(option.getSourceNetmask()));
        assertThat(nmask2, equalTo(option.getScopeNetmask()));
        assertThat(ipaddr2, equalTo(option.getAddress()));
        nameServer.setEcsEnable(false);
    }
   
}
