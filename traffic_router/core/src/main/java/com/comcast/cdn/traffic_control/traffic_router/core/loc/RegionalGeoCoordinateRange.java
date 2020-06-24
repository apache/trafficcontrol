package com.comcast.cdn.traffic_control.traffic_router.core.loc;

import com.fasterxml.jackson.annotation.JsonProperty;

import java.io.Serializable;

public class RegionalGeoCoordinateRange implements Serializable {
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

    public void setMinLat(double minLat) {
        this.minLat = minLat;
    }

    public double getMinLon() {
        return minLon;
    }

    public void setMinLon(double minLon) {
        this.minLon = minLon;
    }

    public double getMaxLat() {
        return maxLat;
    }

    public void setMaxLat(double maxLat) {
        this.maxLat = maxLat;
    }

    public double getMaxLon() {
        return maxLon;
    }

    public void setMaxLon(double maxLon) {
        this.maxLon = maxLon;
    }
}
