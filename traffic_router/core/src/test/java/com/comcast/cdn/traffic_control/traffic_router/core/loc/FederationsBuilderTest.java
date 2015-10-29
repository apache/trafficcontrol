package com.comcast.cdn.traffic_control.traffic_router.core.loc;

import org.junit.Test;

import java.util.List;

import static org.hamcrest.CoreMatchers.not;
import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.equalTo;
import static org.hamcrest.core.IsNull.nullValue;

public class FederationsBuilderTest {
    @Test
    public void itConsumesValidJSON() throws Exception {
        FederationsBuilder federationsBuilder = new FederationsBuilder();

        String json = "{ response: [ " +
            "{ " +
                "'deliveryService' : 'kable-town-01', " +
                "'mappings' : [ " +
                    "{ 'cname' : 'cname1', " +
                        "'ttl' : '86400', " +
                        "'resolve4' : [ '192.168.56.78/24', '192.168.45.67/24' ], " +
                        "'resolve6' : [ 'fdfe:dcba:9876:5::/64', 'fd12:3456:789a:1::/64' ] " +
                    "}, " +
                    "{ 'cname' : 'cname2', 'ttl' : '86400' } " +
                "] " +
            "}, " +
            "{ " +
                "'deliveryService' : 'kable-town-02', " +
                "'mappings' : [ { 'cname' : 'cname4' , 'ttl' : '86400' } ]" +
            "} " +
        "] }";

        List<Federation> federations = federationsBuilder.fromJSON(json);

        assertThat(federations.size(), equalTo(2));
        assertThat(federations.get(0).getDeliveryService(), equalTo("kable-town-01"));
        assertThat(federations.get(0).getFederationMappings(), not(nullValue()));
    }
}
