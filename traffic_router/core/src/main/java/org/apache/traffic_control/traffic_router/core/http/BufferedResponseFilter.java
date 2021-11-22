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

package org.apache.traffic_control.traffic_router.core.http;

import com.google.common.net.HttpHeaders;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.springframework.web.filter.OncePerRequestFilter;

import javax.servlet.FilterChain;
import javax.servlet.ServletException;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import java.io.IOException;

public class BufferedResponseFilter extends OncePerRequestFilter {
	public static final Logger LOGGER = LogManager.getLogger(BufferedResponseFilter.class);

	public void doFilterInternal(final HttpServletRequest request, final HttpServletResponse response, final FilterChain chain) throws IOException, ServletException {
		final BufferedResponse responseWrapper = new BufferedResponse(response);

		chain.doFilter(request, responseWrapper);

		// Close the connection without waiting for the 10-second connect timeout,
		// in case the client does not close the connection. Even though this is
		// the only case for which we are interested in sending Connection: close,
		// sending it sometimes means we must always send it. From RFC 2616:
		// > HTTP/1.1 applications that do not support persistent connections MUST
		// > include the "close" connection option in every message.
		response.addHeader(HttpHeaders.CONNECTION, "close");

		responseWrapper.copyBodyToResponse();
	}
}
