package com.comcast.cdn.traffic_control.traffic_router.core.request;

import com.comcast.cdn.traffic_control.traffic_router.core.http.HTTPAccessEventBuilder;
import com.comcast.cdn.traffic_control.traffic_router.core.http.HTTPAccessRecord;
import org.apache.log4j.Logger;
import org.springframework.web.filter.OncePerRequestFilter;

import javax.servlet.FilterChain;
import javax.servlet.ServletException;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import java.io.IOException;
import java.util.Date;

public class StaticContentAccessFilter extends OncePerRequestFilter {
	private static final Logger ACCESS = Logger.getLogger("com.comcast.cdn.traffic_control.traffic_router.core.access");
	private boolean doNotLog = false;

	@Override
	protected void doFilterInternal(final HttpServletRequest request, final HttpServletResponse response, final FilterChain filterChain) throws ServletException, IOException {
		final Date requestDate = new Date();
		filterChain.doFilter(request, response);

		if (doNotLog) {
			return;
		}

		final HTTPAccessRecord access = new HTTPAccessRecord.Builder(requestDate, request).build();
		ACCESS.info(HTTPAccessEventBuilder.create(access));
	}

	public void setDoNotLog(final String logAccessString) {
		this.doNotLog = Boolean.valueOf(logAccessString);
	}
}
