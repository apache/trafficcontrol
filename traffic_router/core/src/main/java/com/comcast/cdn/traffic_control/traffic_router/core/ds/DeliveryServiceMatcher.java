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

package com.comcast.cdn.traffic_control.traffic_router.core.ds;

import java.util.ArrayList;
import java.util.List;
import java.util.regex.Pattern;

import com.comcast.cdn.traffic_control.traffic_router.core.request.HTTPRequest;
import com.comcast.cdn.traffic_control.traffic_router.core.request.Request;

public class DeliveryServiceMatcher {

    public enum Type {
        HOST, HEADER, PATH
    }
	private DeliveryService deliveryService;
	final private List<Matcher> matches = new ArrayList<Matcher>();

	public DeliveryServiceMatcher(final DeliveryService ds) {
		this.deliveryService = ds;
	}

	static class Matcher {
		private final Type type;
		private final String regex;
		private final Pattern pattern;
		final private String target;

		public Matcher(final Type type, final String regex, final String target) {
			this.type = type;
			this.regex = regex;
			pattern = Pattern.compile(regex, Pattern.CASE_INSENSITIVE);
			this.target = target;
		}
		public boolean matches(final Request request) {
			boolean result = false;
			String target = null;
			switch (type) {
			case HOST:
				target = request.getHostname();
				break;
			case HEADER:
				if (request instanceof HTTPRequest) {
					final HTTPRequest httpRequest = (HTTPRequest) request;
					target = httpRequest.getHeaders().get(this.target);
				}
				break;
			case PATH:
				if (request instanceof HTTPRequest) {
					final HTTPRequest httpRequest = (HTTPRequest) request;
					target = httpRequest.getPath();
					if (httpRequest.getQueryString() != null) {
						final StringBuilder sb = new StringBuilder(target);
						sb.append('?');
						sb.append(httpRequest.getQueryString());
						target = sb.toString();
					}
				}
				break;
			default:
				break;
			}

			if (target != null) {
				result = pattern.matcher(target).matches();
			}
			return result;
		}
		public String getRegex() {
			return regex;
		}
	}
	

	public DeliveryService getDeliveryService() {
		return deliveryService;
	}

	public void setDeliveryService(final DeliveryService deliveryService) {
		this.deliveryService = deliveryService;
	}

	public void addMatch(final Type type, final String string, final String target) {
		matches .add(new Matcher(type,string,target));
	}

	public boolean matches(final Request request) {
		for(Matcher m : matches) {
			if(!m.matches(request)) {
				return false;
			}
		}
		return true;
	}
}
