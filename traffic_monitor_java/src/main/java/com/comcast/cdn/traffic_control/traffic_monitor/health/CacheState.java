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

package com.comcast.cdn.traffic_control.traffic_monitor.health;

import com.comcast.cdn.traffic_control.traffic_monitor.config.Cache;
import static com.comcast.cdn.traffic_control.traffic_monitor.health.HealthDeterminer.AdminStatus;
import static com.comcast.cdn.traffic_control.traffic_monitor.health.HealthDeterminer.AdminStatus.ADMIN_DOWN;
import static com.comcast.cdn.traffic_control.traffic_monitor.health.HealthDeterminer.AdminStatus.OFFLINE;
import static com.comcast.cdn.traffic_control.traffic_monitor.health.HealthDeterminer.AdminStatus.REPORTED;
import static com.comcast.cdn.traffic_control.traffic_monitor.health.HealthDeterminer.AdminStatus.STANDBY;

public class CacheState extends AbstractState {
	transient private Cache cache;
	public static final String STATUS = "status";
	public static final String ERROR_STRING = "error-string";

	public CacheState(final String id) {
		super(id);
	}

	public void setCache(final Cache cache) {
		this.cache = cache;
	}

	public Cache getCache() {
		return cache;
	}

	public void setError(final String error) {
		putDataPoint(STATUS, cache.getStatus());
		putDataPoint(ERROR_STRING, error);
		final Event.EventType type = Event.EventType.CACHE_STATE_CHANGE;
		type.setType(cache.getType());
		setAvailable(type, getIsAvailable(false), error);
	}

	public boolean getIsAvailable(final boolean isHealthy) {
		final AdminStatus status;
		try {
			status = AdminStatus.valueOf(cache.getStatus());
		} catch (IllegalArgumentException e) {
			return false;
		}

		return getIsAvailable(status, isHealthy);
	}

	public boolean getIsAvailable(AdminStatus status, final boolean isHealthy) {
		return (status == REPORTED) ? isHealthy : (status != ADMIN_DOWN && status != OFFLINE && status != STANDBY);
	}
}