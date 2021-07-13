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

import org.apache.traffic_control.traffic_router.core.ds.DeliveryServiceMatcher.Type;
import org.apache.traffic_control.traffic_router.core.util.ComparableStringByLength;

import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class RequestMatcher implements Comparable<RequestMatcher> {
	// This "meta" pattern is used to strip away all leading and trailing non-word characters except '.' and '-' from the original regex
	private static final String META_REGEX = "([\\W])*([\\w-\\./]+).*";
	private static final Pattern metaPattern = Pattern.compile(META_REGEX);

	private final Type type;
	private final Pattern pattern;
	private String requestHeader = "";
	private final ComparableStringByLength comparableRegex;

	public RequestMatcher(final Type type, final String regex, final String requestHeader) {
		if (type == Type.HEADER && (requestHeader == null || requestHeader.isEmpty())) {
			throw new IllegalArgumentException("Request Header name must be supplied for type HEADER");
		}

		this.type = type;
		this.requestHeader = requestHeader;
		pattern = Pattern.compile(regex, Pattern.CASE_INSENSITIVE);

		final Matcher matcher = metaPattern.matcher(regex);
		matcher.matches();
		comparableRegex = new ComparableStringByLength(matcher.group(2));
	}

	public RequestMatcher(final Type type, final String regex) {
		this(type, regex, "");
	}

	public boolean matches(final Request request) {
		final String target = getTarget(request);

		if (target == null) {
			return false;
		}

		return pattern.matcher(target).matches();
	}

	public Type getType() {
		return type;
	}

	public Pattern getPattern() {
		return pattern;
	}

	private String getTarget(final Request request) {
		if (type == Type.HOST) {
			return request.getHostname();
		}

		if (!(request instanceof HTTPRequest)) {
			return null;
		}

		final HTTPRequest httpRequest = (HTTPRequest) request;
		if (type == Type.HEADER) {
		   if (httpRequest.getHeaders() != null) {
			   return httpRequest.getHeaders().get(requestHeader);
		   }
		   return null;
		}

		if (type == Type.PATH) {
			if (httpRequest.getQueryString() == null) {
				return httpRequest.getPath();
			}

			return httpRequest.getPath() + "?" + httpRequest.getQueryString();
		}

		return null;
	}

	@Override
	public int compareTo(final RequestMatcher other) {
		if (this == other || this.equals(other)) {
			return 0;
		}
		return this.comparableRegex.compareTo(other.comparableRegex);
	}

	@Override
	@SuppressWarnings({"PMD.IfStmtsMustUseBraces" , "PMD.CyclomaticComplexity", "PMD.NPathComplexity"})
	public boolean equals(final Object other) {
		if (this == other) return true;
		if (other == null || getClass() != other.getClass()) return false;

		final RequestMatcher that = (RequestMatcher) other;

		if (type != that.type) return false;

		if (pattern != null ? !pattern.pattern().equals(that.pattern.pattern()) : that.pattern != null) return false;
		if (requestHeader != null ? !requestHeader.equals(that.requestHeader) : that.requestHeader != null)
			return false;
		return !(comparableRegex != null ? !comparableRegex.equals(that.comparableRegex) : that.comparableRegex != null);

	}

	@Override
	@SuppressWarnings({"PMD.CyclomaticComplexity", "PMD.NPathComplexity"})
	public int hashCode() {
		int result = type != null ? type.hashCode() : 0;
		result = 31 * result + (pattern != null ? pattern.pattern().hashCode() : 0);
		result = 31 * result + (requestHeader != null ? requestHeader.hashCode() : 0);
		result = 31 * result + (comparableRegex != null ? comparableRegex.hashCode() : 0);
		return result;
	}

	@Override
	public String toString() {
		return "RequestMatcher{" +
			"type=" + type +
			", pattern=" + pattern +
			", requestHeader='" + requestHeader + '\'' +
			", comparableRegex=" + comparableRegex +
			'}';
	}
}
