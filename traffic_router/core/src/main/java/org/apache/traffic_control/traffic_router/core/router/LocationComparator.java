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

import java.util.Comparator;

import org.apache.traffic_control.traffic_router.core.edge.Location;
import org.apache.traffic_control.traffic_router.geolocation.Geolocation;

/**
 * 
 */
public class LocationComparator implements Comparator<Location> {

    private final Geolocation sourceLocation;

    /**
     * @param sourceLocation
     */
    public LocationComparator(final Geolocation sourceLocation) {
        this.sourceLocation = sourceLocation;
    }

    /*
     * (non-Javadoc)
     * 
     * @see java.util.Comparator#compare(java.lang.Object, java.lang.Object)
     */
    @Override
    public int compare(final Location loc1, final Location loc2) {
        final Double d1 = sourceLocation.getDistanceFrom(loc1.getGeolocation());
        final Double d2 = sourceLocation.getDistanceFrom(loc2.getGeolocation());
        return d1.compareTo(d2);
    }

}
