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

import java.util.Comparator;

import org.apache.traffic_control.traffic_router.geolocation.Geolocation;

public class SteeringGeolocationComparator implements Comparator<SteeringResult> {

    private final Geolocation clientLocation;

    public SteeringGeolocationComparator(final Geolocation clientLocation) {
        this.clientLocation = clientLocation;
    }

    @Override
    @SuppressWarnings({"PMD.CyclomaticComplexity"})
    public int compare(final SteeringResult result1, final SteeringResult result2) {
        final Geolocation originGeo1 = result1.getSteeringTarget().getGeolocation();
        final Geolocation originGeo2 = result2.getSteeringTarget().getGeolocation();

        final Geolocation cacheGeo1 = result1.getCache().getGeolocation();
        final Geolocation cacheGeo2 = result2.getCache().getGeolocation();

        // null origin geolocations are considered greater than (i.e. farther away) than non-null origin geolocations
        if (originGeo1 != null && originGeo2 == null) {
            return -1;
        }
        if (originGeo1 == null && originGeo2 != null) {
            return 1;
        }
        if (originGeo1 == null && originGeo2 == null) {
            return 0;
        }

        // same cache and origin locations, prefer lower geoOrder
        if (cacheGeo1.equals(cacheGeo2) && originGeo1.equals(originGeo2)) {
            return Integer.compare(result1.getSteeringTarget().getGeoOrder(), result2.getSteeringTarget().getGeoOrder());
        }

        final double distanceFromClientToCache1 = clientLocation.getDistanceFrom(cacheGeo1);
        final double distanceFromClientToCache2 = clientLocation.getDistanceFrom(cacheGeo2);

        final double distanceFromCacheToOrigin1 = cacheGeo1.getDistanceFrom(originGeo1);
        final double distanceFromCacheToOrigin2 = cacheGeo2.getDistanceFrom(originGeo2);

        final double totalDistance1 = distanceFromClientToCache1 + distanceFromCacheToOrigin1;
        final double totalDistance2 = distanceFromClientToCache2 + distanceFromCacheToOrigin2;

        // different cache and origin locations, prefer shortest total distance
        if (totalDistance1 != totalDistance2) {
            // TODO: if the difference is smaller than a certain threshold/ratio, still prefer the closer edge even though distance is greater?
            return Double.compare(totalDistance1, totalDistance2);
        }

        // total distance is equal, prefer the closest edge to the client
        return Double.compare(distanceFromClientToCache1, distanceFromClientToCache2);

    }

}
