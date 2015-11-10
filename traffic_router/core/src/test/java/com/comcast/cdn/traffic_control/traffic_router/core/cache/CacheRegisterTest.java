package com.comcast.cdn.traffic_control.traffic_router.core.cache;

import com.comcast.cdn.traffic_control.traffic_router.core.ds.DeliveryService;
import com.comcast.cdn.traffic_control.traffic_router.core.ds.DeliveryServiceMatcher;
import com.comcast.cdn.traffic_control.traffic_router.core.request.HTTPRequest;
import org.junit.Before;
import org.junit.Test;

import java.util.TreeSet;

import static com.comcast.cdn.traffic_control.traffic_router.core.ds.DeliveryServiceMatcher.Type.HOST;
import static com.comcast.cdn.traffic_control.traffic_router.core.ds.DeliveryServiceMatcher.Type.PATH;
import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.equalTo;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;

public class CacheRegisterTest {
	private final CacheRegister cacheRegister = new CacheRegister();
	@Before
	public void before() {
		DeliveryService deliveryService1 = mock(DeliveryService.class);
		when(deliveryService1.getId()).thenReturn("delivery service 1");

		DeliveryService deliveryService2 = mock(DeliveryService.class);
		when(deliveryService2.getId()).thenReturn("delivery service 2");

		DeliveryServiceMatcher deliveryServiceMatcher1 = new DeliveryServiceMatcher(deliveryService1);
		deliveryServiceMatcher1.addMatch(HOST, ".*\\.service01-kabletown\\..*", "");
		deliveryServiceMatcher1.addMatch(PATH, ".*abc.*", "");

		DeliveryServiceMatcher deliveryServiceMatcher2 = new DeliveryServiceMatcher(deliveryService2);
		deliveryServiceMatcher2.addMatch(HOST, ".*\\.service01-kabletown\\..*", "");
		deliveryServiceMatcher2.addMatch(PATH, ".*abcde.*", "");

		DeliveryServiceMatcher deliveryServiceMatcher3 = new DeliveryServiceMatcher(deliveryService2);
		deliveryServiceMatcher3.addMatch(HOST, ".*\\.service01-kabletown\\..*", "");
		deliveryServiceMatcher3.addMatch(PATH, ".*abcd.*", "");

		TreeSet<DeliveryServiceMatcher> deliveryServiceMatchers = new TreeSet<DeliveryServiceMatcher>();
		deliveryServiceMatchers.add(deliveryServiceMatcher1);
		deliveryServiceMatchers.add(deliveryServiceMatcher2);
		deliveryServiceMatchers.add(deliveryServiceMatcher3);

		cacheRegister.setHttpDeliveryServiceMatchers(deliveryServiceMatchers);
	}

	@Test
	public void itPicksTheMostSpecificDeliveryService() {
		HTTPRequest httpRequest = new HTTPRequest();
		httpRequest.setHostname("foo.service01-kabletown.com");
		httpRequest.setPath("foo/abcde/bar");

		assertThat(cacheRegister.getDeliveryService(httpRequest, true).getId(), equalTo("delivery service 2"));
	}
}
