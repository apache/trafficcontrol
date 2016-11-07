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

package com.comcast.cdn.traffic_control.traffic_router.core.http;

import com.comcast.cdn.traffic_control.traffic_router.core.ds.DeliveryService;
import com.comcast.cdn.traffic_control.traffic_router.core.request.HTTPRequest;
import com.comcast.cdn.traffic_control.traffic_router.core.router.HTTPRouteResult;
import com.comcast.cdn.traffic_control.traffic_router.core.router.StatTracker;
import com.comcast.cdn.traffic_control.traffic_router.core.router.TrafficRouter;
import com.comcast.cdn.traffic_control.traffic_router.core.router.TrafficRouterManager;
import com.comcast.cdn.traffic_control.traffic_router.geolocation.GeolocationException;
import org.apache.log4j.Logger;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.filter.OncePerRequestFilter;

import javax.servlet.FilterChain;
import javax.servlet.ServletException;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import java.io.IOException;
import java.net.URL;
import java.util.Date;
import java.util.List;
import java.util.Map;
import java.util.Set;

public class RouterFilter extends OncePerRequestFilter {
	private static final Logger ACCESS = Logger.getLogger("com.comcast.cdn.traffic_control.traffic_router.core.access");

	@Autowired
	private TrafficRouterManager trafficRouterManager;

	@Autowired
	private StatTracker statTracker;

	private List<String> staticContentWhiteList;

	private boolean doNotLog = false;

	@Override
	public void doFilterInternal(final HttpServletRequest request, final HttpServletResponse response, final FilterChain chain) throws IOException, ServletException {
		final Date requestDate = new Date();

		if (request.getLocalPort() == trafficRouterManager.getApiPort()) {
			chain.doFilter(request, response);
			return;
		}

		if (staticContentWhiteList.contains(request.getRequestURI())) {
			chain.doFilter(request, response);

			if (doNotLog) {
				return;
			}

			final HTTPAccessRecord access = new HTTPAccessRecord.Builder(requestDate, request).build();
			ACCESS.info(HTTPAccessEventBuilder.create(access));
			return;
		}

		final HTTPAccessRecord httpAccessRecord = new HTTPAccessRecord.Builder(requestDate, request).build();
		writeHttpResponse(response, request, new HTTPRequest(request), StatTracker.getTrack(), httpAccessRecord);
	}

	private void writeHttpResponse(final HttpServletResponse response, final HttpServletRequest httpServletRequest,
	                               final HTTPRequest request, final StatTracker.Track track, final HTTPAccessRecord httpAccessRecord) throws IOException {
		final String format = httpServletRequest.getParameter("format");
		final HTTPAccessRecord.Builder httpAccessRecordBuilder = new HTTPAccessRecord.Builder(httpAccessRecord);
		DeliveryService deliveryService = null;
		try {
			final TrafficRouter trafficRouter = trafficRouterManager.getTrafficRouter();
			final HTTPRouteResult routeResult = trafficRouter.route(request, track);

			if (routeResult != null) {
				deliveryService = routeResult.getDeliveryService();
			}

			if (routeResult == null || routeResult.getUrl() == null) {
				setErrorResponseCode(response, httpAccessRecordBuilder, routeResult);
			} else {
				final URL location = routeResult.getUrl();
				final Map<String, String> responseHeaders = deliveryService.getResponseHeaders();

				for (final String key : responseHeaders.keySet()) {
					response.addHeader(key, responseHeaders.get(key));
				}

				httpAccessRecordBuilder.responseURL(location);

				if("json".equals(format)) {
					response.setContentType("application/json"); // "text/plain"
					response.getWriter().println("{\"location\": \""+location.toString()+"\" }");
					httpAccessRecordBuilder.responseCode(HttpServletResponse.SC_OK);
				} else {
					httpAccessRecordBuilder.responseCode(HttpServletResponse.SC_MOVED_TEMPORARILY);
					response.sendRedirect(location.toString());
				}
			}
		} catch (final IOException e) {
			httpAccessRecordBuilder.responseCode(-1);
			httpAccessRecordBuilder.responseURL(null);
			httpAccessRecordBuilder.rerr(e.getMessage());
			throw e;
		} catch (GeolocationException e) {
			httpAccessRecordBuilder.responseCode(-1);
			httpAccessRecordBuilder.responseURL(null);
			httpAccessRecordBuilder.rerr(e.getMessage());
		} finally {
			final Set<String> requestHeaders = trafficRouterManager.getTrafficRouter().getRequestHeaders();
			if (deliveryService != null) {
				requestHeaders.addAll(deliveryService.getRequestHeaders());
			}

			final Map<String,String> accessRequestHeaders = new HttpAccessRequestHeaders().makeMap(httpServletRequest, requestHeaders);

			final HTTPAccessRecord access = httpAccessRecordBuilder.resultType(track.getResult())
				.resultLocation(track.getResultLocation())
				.requestHeaders(accessRequestHeaders)
				.regionalGeoResult(track.getRegionalGeoResult())
				.build();
			ACCESS.info(HTTPAccessEventBuilder.create(access));
			statTracker.saveTrack(track);
		}
	}

	private void setErrorResponseCode(final HttpServletResponse response,
	                                  final HTTPAccessRecord.Builder httpAccessRecordBuilder, final HTTPRouteResult routeResult) throws IOException {

		if (routeResult != null && routeResult.getResponseCode() > 0) {
			httpAccessRecordBuilder.responseCode(routeResult.getResponseCode());
			response.sendError(routeResult.getResponseCode());
			return;
		}

		httpAccessRecordBuilder.responseCode(HttpServletResponse.SC_SERVICE_UNAVAILABLE);
		response.sendError(HttpServletResponse.SC_SERVICE_UNAVAILABLE);
	}

	public void setDoNotLog(final String logAccessString) {
		this.doNotLog = Boolean.valueOf(logAccessString);
	}

	public void setStaticContentWhiteList(final List<String> staticContentWhiteList) {
		this.staticContentWhiteList = staticContentWhiteList;
	}
}
