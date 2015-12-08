/*
 * Copyright 2015 Comcast Cable Communications Management, LLC
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

import java.io.IOException;
import java.net.URL;
import java.util.Date;
import java.util.Enumeration;
import java.util.HashMap;
import java.util.Map;
import java.util.Set;

import javax.servlet.ServletConfig;
import javax.servlet.ServletException;
import javax.servlet.http.HttpServlet;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;

import org.apache.log4j.Logger;
import org.springframework.context.ApplicationContext;
import org.springframework.web.context.support.WebApplicationContextUtils;

import com.comcast.cdn.traffic_control.traffic_router.core.ds.DeliveryService;
import com.comcast.cdn.traffic_control.traffic_router.core.loc.GeolocationException;
import com.comcast.cdn.traffic_control.traffic_router.core.request.HTTPRequest;
import com.comcast.cdn.traffic_control.traffic_router.core.router.TrafficRouter;
import com.comcast.cdn.traffic_control.traffic_router.core.router.TrafficRouterManager;
import com.comcast.cdn.traffic_control.traffic_router.core.router.HTTPRouteResult;
import com.comcast.cdn.traffic_control.traffic_router.core.router.StatTracker;
import com.comcast.cdn.traffic_control.traffic_router.core.router.StatTracker.Track;

/**
 * Servlet to handle content routing requests.
 */
@SuppressWarnings("PMD.MoreThanOneLogger")
public class TRServlet extends HttpServlet {
	private static final long serialVersionUID = 1L;

	private static final Logger LOGGER = Logger.getLogger(TRServlet.class);
	private static final Logger ACCESS = Logger.getLogger("com.comcast.cdn.traffic_control.traffic_router.core.access");


	private TrafficRouterManager trafficRouterManager;
	private StatTracker statTracker;

	/*
	 * (non-Javadoc)
	 * 
	 * @see javax.servlet.GenericServlet#init(javax.servlet.ServletConfig)
	 */
	@Override
	public void init(final ServletConfig config) throws ServletException {

		final ApplicationContext context = WebApplicationContextUtils.getWebApplicationContext(config
				.getServletContext());
		setTrafficRouterManager(context.getBean(TrafficRouterManager.class));
		setStatTracker(context.getBean(StatTracker.class));

		super.init(config);
	}

	/**
	 * Sets trafficRouter.
	 * 
	 * @param trafficRouterManager
	 *            the trafficRouterManager to set
	 */
	public void setTrafficRouterManager(final TrafficRouterManager trafficRouterManager) {
		this.trafficRouterManager = trafficRouterManager;
	}

	/*
	 * (non-Javadoc)
	 * 
	 * @see javax.servlet.http.HttpServlet#doGet(javax.servlet.http.HttpServletRequest,
	 * javax.servlet.http.HttpServletResponse)
	 */
	@Override
	protected void doGet(final HttpServletRequest request, final HttpServletResponse response) 
			throws ServletException, IOException {
		final Date requestDate = new Date();

		final HTTPRequest req = new HTTPRequest();
		req.setClientIP(request.getRemoteAddr());
		req.setPath(request.getPathInfo());
		req.setQueryString(request.getQueryString());
		req.setHostname(request.getServerName());
		req.setRequestedUrl(request.getRequestURL().toString());

		final StatTracker.Track track = StatTracker.getTrack();
		final String xmm = request.getHeader("X-MM-Client-IP");
		final String fip = request.getParameter("fakeClientIpAddress");

		if (xmm != null) {
			LOGGER.info("X-MM-Client-IP value (header): " + xmm + ", for " + req.getHostname());
			req.setClientIP(xmm);
		} else if (fip != null) {
			LOGGER.info("Fake IP Address (param): " + fip + ", for " + req.getHostname());
			req.setClientIP(fip);
		}

		final Map<String, String> headers = new HashMap<String, String>();
		final Enumeration<?> headerNames = request.getHeaderNames();
		while (headerNames.hasMoreElements()) {
			final String name = (String) headerNames.nextElement();
			final String value = request.getHeader(name);
			headers.put(name, value);
		}
		req.setHeaders(headers);

		final HTTPAccessRecord httpAccessRecord = new HTTPAccessRecord.Builder(requestDate, request).build();
		writeHttpResponse(response, request, req, track, httpAccessRecord);
	}

	@SuppressWarnings("PMD.CyclomaticComplexity")
	private void writeHttpResponse(final HttpServletResponse response, final HttpServletRequest httpServletRequest,
			final HTTPRequest request, final Track track, final HTTPAccessRecord httpAccessRecord) throws IOException {
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
				if (routeResult != null && routeResult.getResponseCode() > 0) {
					httpAccessRecordBuilder.responseCode(routeResult.getResponseCode());
					response.sendError(routeResult.getResponseCode());
				} else {
					httpAccessRecordBuilder.responseCode(HttpServletResponse.SC_SERVICE_UNAVAILABLE);
					response.sendError(HttpServletResponse.SC_SERVICE_UNAVAILABLE);
				}
			} else {
				final URL location = routeResult.getUrl();
				final Map<String, String> responseHeaders = deliveryService.getResponseHeaders();

				for (String key : responseHeaders.keySet()) {
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
			httpAccessRecordBuilder.rerr(track.getResultInfo());
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
			requestHeaders.addAll(deliveryService.getRequestHeaders());

			final Map<String,String> accessRequestHeaders = new HttpAccessRequestHeaders().makeMap(httpServletRequest, requestHeaders);

			final HTTPAccessRecord access = httpAccessRecordBuilder.resultType(track.getResult())
				.resultLocation(track.getResultLocation())
				.requestHeaders(accessRequestHeaders)
				.build();
			ACCESS.info(HTTPAccessEventBuilder.create(access));
			statTracker.saveTrack(track);
		}
	}

	public StatTracker getStatTracker() {
		return statTracker;
	}

	public void setStatTracker(final StatTracker statTracker) {
		this.statTracker = statTracker;
	}

}
