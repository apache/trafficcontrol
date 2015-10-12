package com.comcast.cdn.traffic_control.traffic_router.core.dns;

import com.comcast.cdn.traffic_control.traffic_router.core.cache.CacheRegister;
import com.comcast.cdn.traffic_control.traffic_router.core.router.StatTracker;
import com.comcast.cdn.traffic_control.traffic_router.core.router.StatTracker.Track.ResultType;
import com.comcast.cdn.traffic_control.traffic_router.core.router.TrafficRouter;
import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.powermock.api.mockito.PowerMockito;
import org.powermock.core.classloader.annotations.PrepareForTest;
import org.powermock.modules.junit4.PowerMockRunner;
import org.xbill.DNS.Name;
import org.xbill.DNS.SetResponse;
import org.xbill.DNS.Type;
import org.xbill.DNS.Zone;

import java.net.InetAddress;

import static org.mockito.Matchers.any;
import static org.mockito.Matchers.anyInt;
import static org.mockito.Mockito.*;
import static org.powermock.api.mockito.PowerMockito.whenNew;

@RunWith(PowerMockRunner.class)
@PrepareForTest({ZoneManager.class, SignatureManager.class})
public class ZoneManagerUnitTest {
    ZoneManager zoneManager;

    @Before
    public void before() throws Exception {
        TrafficRouter trafficRouter = mock(TrafficRouter.class);
        CacheRegister cacheRegister = mock(CacheRegister.class);
        when(trafficRouter.getCacheRegister()).thenReturn(cacheRegister);

        PowerMockito.spy(ZoneManager.class);
        PowerMockito.doReturn("edge").when(ZoneManager.class, "getDnsRoutingName");
        PowerMockito.doNothing().when(ZoneManager.class, "initTopLevelDomain", cacheRegister);
        PowerMockito.doNothing().when(ZoneManager.class, "initZoneCache", cacheRegister);

        SignatureManager signatureManager = PowerMockito.mock(SignatureManager.class);
        whenNew(SignatureManager.class).withArguments(any(ZoneManager.class), any(CacheRegister.class), any(KeyServer.class)).thenReturn(signatureManager);

        zoneManager = spy(new ZoneManager(trafficRouter, new StatTracker()));
    }

    @Test
    public void itMarksResultTypeInDNSAccessRecord() throws Exception {
        final Name qname = Name.fromString("edge.www.google.com.");
        final InetAddress client = InetAddress.getByName("192.168.56.78");

        SetResponse setResponse = mock(SetResponse.class);
        when(setResponse.isSuccessful()).thenReturn(false);

        Zone zone = mock(Zone.class);
        when(zone.findRecords(any(Name.class), anyInt())).thenReturn(setResponse);

        DNSAccessRecord.Builder builder = new DNSAccessRecord.Builder(1L, client);
        builder = spy(builder);

        doReturn(zone).when(zoneManager).getZone(qname, Type.A);
        when(zoneManager.getZone(qname, Type.A, client, false, builder)).thenCallRealMethod();

        zoneManager.getZone(qname, Type.A, client, false, builder);
        verify(builder).resultType(any(ResultType.class));
    }
}
