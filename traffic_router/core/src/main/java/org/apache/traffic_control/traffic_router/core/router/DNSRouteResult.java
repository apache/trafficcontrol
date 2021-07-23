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

package org.apache.traffic_control.traffic_router.core.router;

import java.util.ArrayList;
import java.util.List;

import org.apache.traffic_control.traffic_router.core.ds.DeliveryService;
import org.apache.traffic_control.traffic_router.core.edge.InetRecord;

public class DNSRouteResult implements RouteResult {
	private List<InetRecord> addresses;
	private DeliveryService deliveryService;

	public Object getResult() {
		return getAddresses();
	}

	public List<InetRecord> getAddresses() {
		return addresses;
	}

	public void setAddresses(final List<InetRecord> addresses) {
		this.addresses = addresses;
	}

	public void addAddresses(final List<InetRecord> addresses) {
		if (this.addresses == null) {
			this.addresses = new ArrayList<>();
		}

		this.addresses.addAll(addresses);
	}

	public DeliveryService getDeliveryService() {
		return deliveryService;
	}

	public void setDeliveryService(final DeliveryService deliveryService) {
		this.deliveryService = deliveryService;
	}
}
