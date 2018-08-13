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

package com.comcast.cdn.traffic_control.traffic_router.core.edge;

import java.util.Collection;
import java.util.Map;

/**
 * Provides access to the configured {@link CacheLocation}s and the {@link Cache}s that are a part
 * of them.
 */
public interface CacheLocationManager {
	/**
	 * Gets the {@link CacheLocation} specified by the provided ID.
	 *
	 * @param id
	 *            the ID for the desired <code>CacheLocation</code>
	 * @return the <code>CacheLocation</code> or null if no location exists for the specified ID.
	 */
	CacheLocation getCacheLocation(final String id);

	/**
	 * Returns the configured {@link CacheLocation}s.
	 * 
	 * @return the configured <code>CacheLocations</code> or an empty {@link Collection} if no
	 *         locations are configured
	 */
	Collection<CacheLocation> getCacheLocations();

	void setCacheMap(Map<String,Cache> map);
	Map<String,Cache> getCacheMap();
}
