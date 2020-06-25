package com.comcast.cdn.traffic_control.traffic_router.core.loc;

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
