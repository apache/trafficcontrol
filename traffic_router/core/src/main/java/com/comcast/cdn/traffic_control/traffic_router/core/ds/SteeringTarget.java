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

import com.comcast.cdn.traffic_control.traffic_router.core.hash.DefaultHashable;
import com.fasterxml.jackson.annotation.JsonProperty;

import java.util.ArrayList;
import java.util.List;
import java.util.regex.Pattern;

public class SteeringTarget extends DefaultHashable {
	@JsonProperty
	private String deliveryService;
	@JsonProperty
	private int weight;
	@JsonProperty
	private List<String> filters = new ArrayList<String>();
	final private List<Pattern> filterPatterns = new ArrayList<Pattern>();

	public DefaultHashable generateHashes() {
		return generateHashes(deliveryService, weight);
	}

	public void setDeliveryService(final String deliveryService) {
		this.deliveryService = deliveryService;
	}

	public String getDeliveryService() {
		return deliveryService;
	}

	public void setWeight(final int weight) {
		this.weight = weight;
	}

	public int getWeight() {
		return weight;
	}

	public List<String> getFilters() {
		return filters;
	}

	public void setFilters(final List<String> filters) {
		this.filters = filters;
		filterPatterns.clear();
		for (String filter : filters) {
			filterPatterns.add(Pattern.compile(filter));
		}
	}

	public boolean hasMatchingFilter(final String path) {
		for (Pattern filterPattern : filterPatterns) {
			if (filterPattern.matcher(path).matches()) {
				return true;
			}
		}

		return false;
	}

	@SuppressWarnings("PMD")
	@Override
	public boolean equals(Object o) {
		if (this == o) return true;
		if (o == null || getClass() != o.getClass()) return false;

		SteeringTarget that = (SteeringTarget) o;

		if (weight != that.weight) return false;
		if (deliveryService != null ? !deliveryService.equals(that.deliveryService) : that.deliveryService != null)
			return false;
		return filters != null ? filters.equals(that.filters) : that.filters == null;

	}

	@Override
	public int hashCode() {
		int result = deliveryService != null ? deliveryService.hashCode() : 0;
		result = 31 * result + weight;
		result = 31 * result + (filters != null ? filters.hashCode() : 0);
		return result;
	}
}
