package com.comcast.cdn.traffic_control.traffic_router.core.request;

import static com.comcast.cdn.traffic_control.traffic_router.core.ds.DeliveryServiceMatcher.Type.HOST;
import static com.comcast.cdn.traffic_control.traffic_router.core.ds.DeliveryServiceMatcher.Type.PATH;
import static com.comcast.cdn.traffic_control.traffic_router.core.ds.DeliveryServiceMatcher.Type.HEADER;
import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.equalTo;
import static org.junit.Assert.fail;

import org.junit.Test;

import java.util.HashMap;
import java.util.HashSet;
import java.util.Iterator;
import java.util.Map;
import java.util.Set;
import java.util.TreeSet;

public class RequestMatcherTest {
	@Test
	public void itMatchesByHost() {
		Request request = new Request();
		request.setHostname("foo.host.bar");

		RequestMatcher requestMatcher = new RequestMatcher(HOST, ".*\\.host\\..*", "");
		assertThat(requestMatcher.matches(request), equalTo(true));
	}

	@Test
	public void itMatchesByPath() {
		HTTPRequest httpRequest = new HTTPRequest();
		httpRequest.setPath("/foo/path/bar");

		RequestMatcher requestMatcher = new RequestMatcher(PATH, ".*path.*");
		assertThat(requestMatcher.matches(httpRequest), equalTo(true));
	}

	@Test
	public void itMatchesByRequestHeader() {
		Map<String, String> headers = new HashMap<String, String>();
		headers.put("foo", "something.abcd.else");

		HTTPRequest httpRequest = new HTTPRequest();
		httpRequest.setHeaders(headers);

		RequestMatcher requestMatcher = new RequestMatcher(HEADER, ".*abcd.*", "foo");
		assertThat(requestMatcher.matches(httpRequest), equalTo(true));

	}

	@Test
	public void itSupportsOrderingByItsRegex() {
		RequestMatcher requestMatcher1 = new RequestMatcher(PATH, ".*abcd.*");
		RequestMatcher requestMatcher2 = new RequestMatcher(PATH, ".*abcde.*");
		RequestMatcher requestMatcher3 = new RequestMatcher(PATH, ".*bcd.*");
		RequestMatcher requestMatcher4 = new RequestMatcher(PATH, ".*bcdef.*");

		Set<RequestMatcher> set = new TreeSet<RequestMatcher>();
		set.add(requestMatcher1);
		set.add(requestMatcher2);
		set.add(requestMatcher3);
		set.add(requestMatcher4);

		Iterator<RequestMatcher> iterator = set.iterator();

		assertThat(iterator.next(), equalTo(requestMatcher2));
		assertThat(iterator.next(), equalTo(requestMatcher4));
		assertThat(iterator.next(), equalTo(requestMatcher1));
		assertThat(iterator.next(), equalTo(requestMatcher3));
	}

	@Test
	public void itThrowsIllegalArgumentException() {
		try {
			new RequestMatcher(HEADER, "a-regex");
			fail("Should have caught Illegal Argument Exception!");
		} catch (IllegalArgumentException e) {
			assertThat(e.getMessage(), equalTo("Request Header name must be supplied for type HEADER"));
		}
	}

	@Test
	public void itSupportsEquals() {
		RequestMatcher requestMatcher1 = new RequestMatcher(HOST, ".*abc.*");
		RequestMatcher requestMatcher2 = new RequestMatcher(HOST, ".*abc.*");

		assertThat(requestMatcher1, equalTo(requestMatcher2));
		assertThat(requestMatcher2, equalTo(requestMatcher1));
	}

	@Test
	public void itSupportsHashCode() {
		RequestMatcher requestMatcher1 = new RequestMatcher(HOST, ".*abc.*");
		RequestMatcher requestMatcher2 = new RequestMatcher(HOST, ".*abc.*");
		assertThat(requestMatcher1.hashCode(), equalTo(requestMatcher2.hashCode()));
	}
}
