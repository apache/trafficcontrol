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

package org.apache.traffic_control.traffic_router.core.request;

import static org.apache.traffic_control.traffic_router.core.ds.DeliveryServiceMatcher.Type.HOST;
import static org.apache.traffic_control.traffic_router.core.ds.DeliveryServiceMatcher.Type.PATH;
import static org.apache.traffic_control.traffic_router.core.ds.DeliveryServiceMatcher.Type.HEADER;
import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.equalTo;
import static org.junit.Assert.fail;

import org.junit.Test;

import java.util.HashMap;
import java.util.Iterator;
import java.util.Map;
import java.util.Set;
import java.util.TreeSet;

public class RequestMatcherTest {

	@Test
	public void itDoesNotAllowHEADERmatchesWithoutHeaderName() {
		try {
			new RequestMatcher(HEADER, ".*kabletown.*");
			fail("Should have thrown IllegalArgumentException");
		} catch (IllegalArgumentException iae) {
			assertThat(iae.getMessage(), equalTo("Request Header name must be supplied for type HEADER"));
		}

		try {
			new RequestMatcher(HEADER, ".*kabletown.*", "");
			fail("Should have thrown IllegalArgumentException");
		} catch (IllegalArgumentException iae) {
			assertThat(iae.getMessage(), equalTo("Request Header name must be supplied for type HEADER"));
		}
	}

	@Test
	public void itMatchesByHost() {
		Request request = new Request();
		RequestMatcher requestMatcher = new RequestMatcher(HOST, ".*\\.host\\..*", "");

		assertThat(requestMatcher.matches(request), equalTo(false));

		request.setHostname("foo.host.bar");
		assertThat(requestMatcher.matches(request), equalTo(true));
	}

	@Test
	public void itMatchesByPath() {
		RequestMatcher requestMatcher = new RequestMatcher(PATH, ".*path.*");
		Request request = new Request();

		assertThat(requestMatcher.matches(request), equalTo(false));

		HTTPRequest httpRequest = new HTTPRequest();
		httpRequest.setPath("/foo/path/bar");

		assertThat(requestMatcher.matches(httpRequest), equalTo(true));
	}

	@Test
	public void itMatchesByQuery() {
		RequestMatcher requestMatcher = new RequestMatcher(PATH, ".*car=red.*");
		Request request = new Request();

		assertThat(requestMatcher.matches(request), equalTo(false));

		HTTPRequest httpRequest = new HTTPRequest();
		httpRequest.setPath("/foo/path/bar");
		httpRequest.setQueryString("car=red");

		assertThat(requestMatcher.matches(httpRequest), equalTo(true));
	}

	@Test
	public void itMatchesByPathAndQuery() {
		RequestMatcher requestMatcher = new RequestMatcher(PATH, "\\/foo\\/path\\/bar\\?car=red");
		Request request = new Request();

		assertThat(requestMatcher.matches(request), equalTo(false));

		HTTPRequest httpRequest = new HTTPRequest();
		httpRequest.setPath("/foo/path/bar");
		httpRequest.setQueryString("car=red");

		assertThat(requestMatcher.matches(httpRequest), equalTo(true));
	}

	@Test
	public void itMatchesByRequestHeader() {
		RequestMatcher requestMatcher = new RequestMatcher(HEADER, ".*kabletown.*", "Host");
		Request request = new Request();
		assertThat(requestMatcher.matches(request), equalTo(false));

		Map<String, String> headers = new HashMap<String, String>();
		headers.put("Host", "www.kabletown.com");

		HTTPRequest httpRequest = new HTTPRequest();
		httpRequest.setHeaders(headers);

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
