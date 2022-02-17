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

package org.apache.traffic_control.traffic_router.core.edge;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.HashMap;
import java.util.HashSet;
import java.util.List;
import java.util.Map;
import java.util.Set;

import org.apache.commons.lang3.builder.EqualsBuilder;
import org.apache.commons.lang3.builder.HashCodeBuilder;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.apache.traffic_control.traffic_router.geolocation.Geolocation;

/**
 * A physical location that has caches.
 */
public class CacheLocation extends Location {

	public static final Logger LOGGER = LogManager.getLogger(CacheLocation.class);

	private final Map<String, Cache> caches;
	private List<String> backupCacheGroups = null;
	private boolean useClosestGeoOnBackupFailure = true;
	private final Set<LocalizationMethod> enabledLocalizationMethods;

	public enum LocalizationMethod {
		DEEP_CZ,
		CZ,
		GEO
	}

	/**
	 * Creates a CacheLocation with the specified ID at the specified location.
	 * 
	 * @param id
	 *            the id of the location
	 * @param geolocation
	 *            the coordinates of this location
	 */
	public CacheLocation(final String id, final Geolocation geolocation) {
		this(id, geolocation, null, true, new HashSet<>());
	}

	public CacheLocation(final String id, final Geolocation geoLocation, final Set<LocalizationMethod> enabledLocalizationMethods) {
		this(id, geoLocation, null, true, enabledLocalizationMethods);
	}

	/**
	 * Creates a CacheLocation with the specified ID at the specified location.
	 * 
	 * @param id
	 *            the id of the location
	 * @param geolocation
	 *            the coordinates of this location
	 *
	 * @param backupCacheGroups
	 *            the backup cache groups for this id
	 *
	 * @param useClosestGeoOnBackupFailure
	 *            the backup fallback setting for this id
	 */
	public CacheLocation(
			final String id,
			final Geolocation geolocation,
			final List<String> backupCacheGroups,
			final boolean useClosestGeoOnBackupFailure,
			final Set<LocalizationMethod> enabledLocalizationMethods
	) {
		super(id, geolocation);
		this.backupCacheGroups = backupCacheGroups;
		this.useClosestGeoOnBackupFailure = useClosestGeoOnBackupFailure;
		this.enabledLocalizationMethods = enabledLocalizationMethods;
		if (this.enabledLocalizationMethods.isEmpty()) {
			this.enabledLocalizationMethods.addAll(Arrays.asList(LocalizationMethod.values()));
		}
		caches = new HashMap<String, Cache>();
	}

	public boolean isEnabledFor(final LocalizationMethod localizationMethod) {
		return enabledLocalizationMethods.contains(localizationMethod);
	}

	/**
	 * Adds the specified cache to this location.
	 * 
	 * @param cache
	 *            the cache to add
	 */
	public void addCache(final Cache cache) {
	    synchronized (caches) {
			caches.put(cache.getId(), cache);
		}
	}

	public void clearCaches() {
		synchronized (caches) {
			caches.clear();
		}
	}

	public void loadDeepCaches(final Set<String> deepCacheNames, final CacheRegister cacheRegister) {
	    synchronized (caches) {
			if (caches.isEmpty() && deepCacheNames != null) {
				for (final String deepCacheName : deepCacheNames) {
					final Cache deepCache = cacheRegister.getCacheMap().get(deepCacheName);
					if (deepCache != null) {
						LOGGER.debug("DDC: Adding " + deepCacheName + " to " + getId());
						caches.put(deepCache.getId(), deepCache);
					}
				}
			}
		}
	}

	@Override
	public boolean equals(final Object obj) {
		if (this == obj) {
			return true;
		} else if (obj instanceof CacheLocation) {
			final CacheLocation rhs = (CacheLocation) obj;
			return new EqualsBuilder()
			.append(getId(), rhs.getId())
			.isEquals();
		} else {
			return false;
		}
	}

	/**
	 * Retrieves the specified {@link Cache} from the location.
	 * 
	 * @param id
	 *            the ID for the desired <code>Cache</code>
	 * @return the cache or <code>null</code> if the cache doesn't exist
	 */
	public Cache getCache(final String id) {
		return caches.get(id);
	}

	/**
	 * Retrieves the {@link Set} of caches at this location.
	 * 
	 * @return the caches
	 */
	public List<Cache> getCaches() {
		return new ArrayList<Cache>(caches.values());
	}

	/**
	 * Gets backupCacheGroups.
	 * 
	 * @return the backupCacheGroups
	 */
	public List<String> getBackupCacheGroups() {
		return backupCacheGroups;
	}

	/**
	 * Tests useClosestGeoOnBackupFailure.
	 * 
	 * @return useClosestGeoOnBackupFailure
	 */
	public boolean isUseClosestGeoLoc() {
		return useClosestGeoOnBackupFailure;
	}

	/**
	 * Determines if the specified {@link Cache} exists at this location.
	 * 
	 * @param id
	 *            the <code>Cache</code> to check
	 * @return true if the <code>Cache</code> is at this location, false otherwise
	 */
	public boolean hasCache(final String id) {
		return caches.containsKey(id);
	}

	@Override
	public int hashCode() {
		return new HashCodeBuilder(1, 31)
		.append(getId())
		.toHashCode();
	}
}
