package com.comcast.cdn.traffic_control.traffic_router.core.http;

import com.comcast.cdn.traffic_control.traffic_router.core.loc.Geolocation;
import com.comcast.cdn.traffic_control.traffic_router.core.router.StatTracker;
import com.comcast.cdn.traffic_control.traffic_router.core.router.StatTracker.Track.ResultType;
import com.comcast.cdn.traffic_control.traffic_router.core.router.StatTracker.Track.ResultDetails;
import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.powermock.core.classloader.annotations.PrepareForTest;
import org.powermock.modules.junit4.PowerMockRunner;

import javax.servlet.http.HttpServletRequest;
import java.net.URL;
import java.util.Date;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.equalTo;
import static org.mockito.Matchers.anyLong;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;
import static org.powermock.api.mockito.PowerMockito.mockStatic;
import static org.powermock.api.mockito.PowerMockito.whenNew;

@RunWith(PowerMockRunner.class)
@PrepareForTest({Date.class, HTTPAccessEventBuilder.class})
public class HTTPAccessEventBuilderTest {
    private HttpServletRequest request;

    @Before
    public void before() throws Exception {
        mockStatic(Date.class);
        Date startDate = mock(Date.class);
        when(startDate.getTime()).thenReturn(144140678000L);
        whenNew(Date.class).withArguments(anyLong()).thenReturn(startDate);

        Date finishDate = mock(Date.class);
        when(finishDate.getTime()).thenReturn(144140678125L);
        whenNew(Date.class).withNoArguments().thenReturn(finishDate);

        request = mock(HttpServletRequest.class);
        when(request.getRequestURL()).thenReturn(new StringBuffer("http://example.com/index.html?foo=bar"));
        when(request.getMethod()).thenReturn("GET");
        when(request.getProtocol()).thenReturn("HTTP/1.1");
        when(request.getRemoteAddr()).thenReturn("192.168.7.6");
    }

    @Test
    public void itGeneratesAccessEvents() throws Exception {
        HTTPAccessRecord.Builder builder = new HTTPAccessRecord.Builder(new Date(144140678000L), request);
        HTTPAccessRecord httpAccessRecord = builder.build();

        String httpAccessEvent = HTTPAccessEventBuilder.create(httpAccessRecord);

        assertThat(httpAccessEvent, equalTo("144140678.000 qtype=HTTP chi=192.168.7.6 url=\"http://example.com/index.html?foo=bar\" cqhm=GET cqhv=HTTP/1.1 rtype=- rloc=\"-\" rdtl=- rerr=\"-\" rurl=\"-\""));
    }

    @Test
    public void itAddsResponseData() throws Exception {

        StatTracker.Track track = new StatTracker.Track();
        HTTPAccessRecord.Builder builder = new HTTPAccessRecord.Builder(new Date(144140633999L), request)
            .resultType(track.getResult())
            .resultLocation(new Geolocation(39.7528,-104.9997))
            .responseCode(302)
            .responseURL(new URL("http://example.com/hereitis/index.html?foo=bar"));

        HTTPAccessRecord httpAccessRecord = builder.resultType(ResultType.CZ).build();
        String httpAccessEvent = HTTPAccessEventBuilder.create(httpAccessRecord);

        assertThat(httpAccessEvent, equalTo("144140678.000 qtype=HTTP chi=192.168.7.6 url=\"http://example.com/index.html?foo=bar\" cqhm=GET cqhv=HTTP/1.1 rtype=CZ rloc=\"39.75,-104.99\" rdtl=- rerr=\"-\" pssc=302 ttms=125 rurl=\"http://example.com/hereitis/index.html?foo=bar\""));
    }

    @Test
    public void itMarksTTMSLessThanMilliAsZero() throws Exception {
        Date fastFinishDate = mock(Date.class);
        when(fastFinishDate.getTime()).thenReturn(144140678000L);
        whenNew(Date.class).withNoArguments().thenReturn(fastFinishDate);

        StatTracker.Track track = new StatTracker.Track();
        HTTPAccessRecord.Builder builder = new HTTPAccessRecord.Builder(new Date(144140633999L), request)
                .resultType(track.getResult())
                .responseCode(302)
                .responseURL(new URL("http://example.com/hereitis/index.html?foo=bar"));

        HTTPAccessRecord httpAccessRecord = builder.build();
        String httpAccessEvent = HTTPAccessEventBuilder.create(httpAccessRecord);

        assertThat(httpAccessEvent, equalTo("144140678.000 qtype=HTTP chi=192.168.7.6 url=\"http://example.com/index.html?foo=bar\" cqhm=GET cqhv=HTTP/1.1 rtype=ERROR rloc=\"-\" rdtl=- rerr=\"-\" pssc=302 ttms=0 rurl=\"http://example.com/hereitis/index.html?foo=bar\""));
    }


    @Test
    public void itRecordsTrafficRouterErrors() throws Exception {
        Date fastFinishDate = mock(Date.class);
        when(fastFinishDate.getTime()).thenReturn(144140678000L);
        whenNew(Date.class).withNoArguments().thenReturn(fastFinishDate);

        StatTracker.Track track = new StatTracker.Track();
        HTTPAccessRecord.Builder builder = new HTTPAccessRecord.Builder(new Date(144140633999L), request)
                .resultType(track.getResult())
                .responseCode(302)
                .rerr("RuntimeException: you're doing it wrong")
                .responseURL(new URL("http://example.com/hereitis/index.html?foo=bar"));

        HTTPAccessRecord httpAccessRecord = builder.build();
        String httpAccessEvent = HTTPAccessEventBuilder.create(httpAccessRecord);

        assertThat(httpAccessEvent, equalTo("144140678.000 qtype=HTTP chi=192.168.7.6 url=\"http://example.com/index.html?foo=bar\" cqhm=GET cqhv=HTTP/1.1 rtype=ERROR rloc=\"-\" rdtl=- rerr=\"RuntimeException: you're doing it wrong\" pssc=302 ttms=0 rurl=\"http://example.com/hereitis/index.html?foo=bar\""));
    }
    
    @Test
    public void itRecordsMissResultDetails() throws Exception {
        Date fastFinishDate = mock(Date.class);
        when(fastFinishDate.getTime()).thenReturn(144140678000L);
        whenNew(Date.class).withNoArguments().thenReturn(fastFinishDate);

        HTTPAccessRecord.Builder builder = new HTTPAccessRecord.Builder(new Date(144140633999L), request)
                .resultType(ResultType.MISS)
                .resultDetails(ResultDetails.DS_NO_BYPASS)
                .responseCode(503);

        HTTPAccessRecord httpAccessRecord = builder.build();
        String httpAccessEvent = HTTPAccessEventBuilder.create(httpAccessRecord);

        assertThat(httpAccessEvent, equalTo("144140678.000 qtype=HTTP chi=192.168.7.6 url=\"http://example.com/index.html?foo=bar\" cqhm=GET cqhv=HTTP/1.1 rtype=MISS rloc=\"-\" rdtl=DS_NO_BYPASS rerr=\"-\" pssc=503 ttms=0 rurl=\"-\""));
    }
}
