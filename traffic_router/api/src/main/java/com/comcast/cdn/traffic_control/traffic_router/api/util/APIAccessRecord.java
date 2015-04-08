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

package com.comcast.cdn.traffic_control.traffic_router.api.util;

import java.util.Date;
import java.util.TimeZone;

import javax.servlet.http.HttpServletRequest;

import org.apache.commons.lang3.time.FastDateFormat;
import org.apache.log4j.Logger;

public class APIAccessRecord {
	private static final Logger ACCESS = Logger.getLogger("com.comcast.cdn.traffic_control.traffic_router.api.access");
	private static final String ACCESS_FORMAT = "API [%s] %s %s %s";
	private static final FastDateFormat FORMATTER = FastDateFormat.getInstance("dd/MMM/yyyy:HH:mm:ss.SSS Z", TimeZone.getTimeZone("GMT"));

	public static void log(final HttpServletRequest request) {
		ACCESS.info(
			String.format(ACCESS_FORMAT, FORMATTER.format(new Date()),
			request.getRemoteAddr(),
			request.getMethod(),
			getUrl(request))
		);
	}

	private static String getUrl(final HttpServletRequest request) {
		final StringBuffer buf = request.getRequestURL();
		if (request.getQueryString() != null) {
			buf.append('?');
			buf.append(request.getQueryString());
		}
		return buf.toString();
	}
}
