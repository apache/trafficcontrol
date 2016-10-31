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

package com.comcast.cdn.traffic_control.traffic_router.core.router;

import java.net.URL;

import com.comcast.cdn.traffic_control.traffic_router.core.ds.DeliveryService;

public class HTTPRouteResult implements RouteResult {
	private URL url;
	private int responseCode;
	private DeliveryService deliveryService;

	@Override
	public Object getResult() {
		return getUrl();
	}

	public URL getUrl() {
		return url;
	}

	public void setUrl(final URL url) {
		this.url = url;
	}

	public int getResponseCode() {
		return responseCode;
	}

	public void setResponseCode(final int rc) {
		this.responseCode = rc;
	}

	public DeliveryService getDeliveryService() {
		return deliveryService;
	}

	public void setDeliveryService(final DeliveryService deliveryService) {
		this.deliveryService = deliveryService;
	}
}
