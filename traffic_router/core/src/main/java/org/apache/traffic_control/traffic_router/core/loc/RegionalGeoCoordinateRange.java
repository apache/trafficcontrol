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

package org.apache.traffic_control.traffic_router.core.loc;

import com.fasterxml.jackson.annotation.JsonProperty;

public class RegionalGeoCoordinateRange {
    @JsonProperty
    private double minLat;
    @JsonProperty
    private double minLon;
    @JsonProperty
    private double maxLat;
    @JsonProperty
    private double maxLon;

    public RegionalGeoCoordinateRange() {
        minLat = 0.0;
        minLon = 0.0;
        maxLat = 0.0;
        maxLon = 0.0;
    }

    public double getMinLat() {
        return minLat;
    }

    public void setMinLat(final double minLat) {
        this.minLat = minLat;
    }

    public double getMinLon() {
        return minLon;
    }

    public void setMinLon(final double minLon) {
        this.minLon = minLon;
    }

    public double getMaxLat() {
        return maxLat;
    }

    public void setMaxLat(final double maxLat) {
        this.maxLat = maxLat;
    }

    public double getMaxLon() {
        return maxLon;
    }

    public void setMaxLon(final double maxLon) {
        this.maxLon = maxLon;
    }
}
