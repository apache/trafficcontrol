package com.comcast.cdn.traffic_control.traffic_router.core.loc;

import com.comcast.cdn.traffic_control.traffic_router.core.util.CidrAddress;
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
            "'cname' : 'cname1', " +
            "'ttl' : '86400', " +
            "'resolve4' : [ '192.168.56.78/24', '192.168.45.67/24' ], " +
            "'resolve6' : [ 'fdfe:dcba:9876:5::/64', 'fd12:3456:789a:1::/64' ] " +
        "}";

        FederationMapping federationMapping = federationMappingBuilder.fromJSON(json);

        assertThat(federationMapping, not(nullValue()));
        assertThat(federationMapping.getCname(), equalTo("cname1"));
        assertThat(federationMapping.getTtl(), equalTo(86400));

        assertThat(federationMapping.getResolve4(),
                containsInAnyOrder(new CidrAddress("192.168.45.67/24"), new CidrAddress("192.168.56.78/24")));
        assertThat(federationMapping.getResolve6(),
                containsInAnyOrder(new CidrAddress("fd12:3456:789a:1::/64"), new CidrAddress("fdfe:dcba:9876:5::/64")));
    }
}
