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

import com.fasterxml.jackson.annotation.JsonProperty;

import java.util.List;

public class Steering {
	@JsonProperty
	private String deliveryService;
	@JsonProperty
	private List<SteeringTarget> targets;

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

	public String getBypassDestination(final String requestPath) {
		for (SteeringTarget target : targets) {
			if (target.hasMatchingFilter(requestPath)) {
				return target.getDeliveryService();
			}
		}

		return null;
	}

	@Override
	@SuppressWarnings("PMD")
	public boolean equals(Object o) {
		if (this == o) return true;
		if (o == null || getClass() != o.getClass()) return false;

		Steering steering = (Steering) o;

		if (deliveryService != null ? !deliveryService.equals(steering.deliveryService) : steering.deliveryService != null)
			return false;
		return targets != null ? targets.equals(steering.targets) : steering.targets == null;

	}

	@Override
	@SuppressWarnings("PMD")
	public int hashCode() {
		int result = deliveryService != null ? deliveryService.hashCode() : 0;
		result = 31 * result + (targets != null ? targets.hashCode() : 0);
		return result;
	}
}
