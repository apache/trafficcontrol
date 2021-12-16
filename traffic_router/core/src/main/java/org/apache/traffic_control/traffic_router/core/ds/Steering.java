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

package org.apache.traffic_control.traffic_router.core.ds;

import com.fasterxml.jackson.annotation.JsonProperty;

import java.util.ArrayList;
import java.util.List;
import java.util.Objects;

public class Steering {
	@JsonProperty
	private String deliveryService;
	@JsonProperty
	private boolean clientSteering;
	@JsonProperty
	private List<SteeringTarget> targets = new ArrayList<SteeringTarget>();
	@JsonProperty
	private List<SteeringFilter> filters = new ArrayList<SteeringFilter>();

	public List<SteeringTarget> getTargets() {
		return targets;
	}

	public void setTargets(final List<SteeringTarget> targets) {
		this.targets = targets;
	}

	public String getDeliveryService() {
		return deliveryService;
	}

	public void setDeliveryService(final String id) {
		this.deliveryService = id;
	}

	public boolean isClientSteering() {
		return clientSteering;
	}

	public void setClientSteering(final boolean clientSteering) {
		this.clientSteering = clientSteering;
	}

	public List<SteeringFilter> getFilters() {
		return filters;
	}

	public void setFilters(final List<SteeringFilter> filters) {
		this.filters = filters;
	}

	public String getBypassDestination(final String requestPath) {
		for (final SteeringFilter filter : filters) {
			if (filter.matches(requestPath) && hasTarget(filter.getDeliveryService())) {
				return filter.getDeliveryService();
			}
		}

		return null;
	}

	public boolean hasTarget(final String deliveryService) {
		for (final SteeringTarget target : targets) {
			if (deliveryService.equals(target.getDeliveryService())) {
				return true;
			}
		}

		return false;
	}

	@Override
	public boolean equals(final Object o) {
		if (this == o) {
			return true;
		}
		if (o == null || getClass() != o.getClass()) {
			return false;
		}

		final Steering steering = (Steering) o;
		return Objects.equals(deliveryService, steering.deliveryService)
				&& Objects.equals(targets, steering.targets)
				&& Objects.equals(filters, steering.filters);
	}

	@Override
	public int hashCode() {
		int result = deliveryService != null ? deliveryService.hashCode() : 0;
		result = 31 * result + (targets != null ? targets.hashCode() : 0);
		result = 31 * result + (filters != null ? filters.hashCode() : 0);
		return result;
	}
}
