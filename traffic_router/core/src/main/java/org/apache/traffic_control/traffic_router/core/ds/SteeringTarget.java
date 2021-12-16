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
import org.apache.traffic_control.traffic_router.core.hash.DefaultHashable;
import org.apache.traffic_control.traffic_router.geolocation.Geolocation;

import java.util.Objects;

public class SteeringTarget extends DefaultHashable {

	private static final double DEFAULT_LAT = 0.0;
	private static final double DEFAULT_LON = 0.0;

	@JsonProperty
	private String deliveryService;
	@JsonProperty
	private int weight;
	@JsonProperty
	private int order = 0;
	@JsonProperty
	private int geoOrder = 0;
	@JsonProperty
	private double latitude = DEFAULT_LAT;
	@JsonProperty
	private double longitude = DEFAULT_LON;

	private Geolocation geolocation;

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

	public void setOrder(final int order) {
		this.order = order;
	}

	public int getOrder() {
		return order;
	}

	public void setGeoOrder(final int geoOrder) {
		this.geoOrder = geoOrder;
	}

	public int getGeoOrder() {
		return geoOrder;
	}

	public void setLatitude(final double latitude) {
		this.latitude = latitude;
	}

	public double getLatitude() {
		return latitude;
	}

	public void setLongitude(final double longitude) {
		this.longitude = longitude;
	}

	public double getLongitude() {
		return longitude;
	}

	public Geolocation getGeolocation() {
		if (geolocation != null) {
			return geolocation;
		}
		if (latitude != DEFAULT_LAT && longitude != DEFAULT_LON) {
			geolocation = new Geolocation(latitude, longitude);
		}
		return geolocation;
	}

	public void setGeolocation(final Geolocation geolocation) {
		this.geolocation = geolocation;
	}

	@Override
	public boolean equals(final Object o) {
		if (this == o) return true;
		if (o == null || getClass() != o.getClass()) return false;

		final SteeringTarget target = (SteeringTarget) o;

		if (weight != target.weight ||
				order != target.order ||
				geoOrder != target.geoOrder ||
				latitude != target.latitude ||
				longitude != target.longitude) return false;
		return Objects.equals(deliveryService, target.deliveryService);

	}

	@Override
	public int hashCode() {
		int result = deliveryService != null ? deliveryService.hashCode() : 0;
		result = 31 * result + weight;
		result = 31 * result + order;
		result = 31 * result + geoOrder;
		result = 31 * result + (int) latitude;
		result = 31 * result + (int) longitude;
		return result;
	}
}
