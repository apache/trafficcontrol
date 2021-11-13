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

package org.apache.traffic_control.traffic_router.core.status.model;

import java.util.List;

/**
 * Model for a CacheLocation.
 */
public class LocationModel {
    private String locationID;
    private String description;
    private double latitude;
    private double longitude;
    
    List<CacheModel> caches;

    /**
     * Gets description.
     * 
     * @return the description
     */
    public String getDescription() {
        return description;
    }

    /**
     * Gets latitude.
     * 
     * @return the latitude
     */
    public double getLatitude() {
        return latitude;
    }

    /**
     * Gets locationID.
     * 
     * @return the locationID
     */
    public String getLocationID() {
        return locationID;
    }

    /**
     * Gets longitude.
     * 
     * @return the longitude
     */
    public double getLongitude() {
        return longitude;
    }

    /**
     * Sets description.
     * 
     * @param description
     *            the description to set
     */
    public void setDescription(final String description) {
        this.description = description;
    }

    /**
     * Sets latitude.
     * 
     * @param latitude
     *            the latitude to set
     */
    public void setLatitude(final double latitude) {
        this.latitude = latitude;
    }

    /**
     * Sets locationID.
     * 
     * @param locationID
     *            the locationID to set
     */
    public void setLocationID(final String locationID) {
        this.locationID = locationID;
    }

    /**
     * Sets longitude.
     * 
     * @param longitude
     *            the longitude to set
     */
    public void setLongitude(final double longitude) {
        this.longitude = longitude;
    }
}
