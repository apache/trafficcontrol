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
 * Model for a Cache.
 */
public class CacheModel {
    private String cacheId;
    private String fqdn;
    private List<String> ipAddresses;
    private int port;
    private String adminStatus;
    private boolean lastUpdateHealthy;
    private long lastUpdateTime;
	private long connections;
	private long currentBW;
	private long availBW;
	boolean cacheOnline;

    /**
     * Gets adminStatus.
     * 
     * @return the adminStatus
     */
    public String getAdminStatus() {
        return adminStatus;
    }

    /**
     * Gets cacheId.
     * 
     * @return the cacheId
     */
    public String getCacheId() {
        return cacheId;
    }

    /**
     * Gets fqdn.
     * 
     * @return the fqdn
     */
    public String getFqdn() {
        return fqdn;
    }

    /**
     * Gets ipAddresses.
     * 
     * @return the ipAddresses
     */
    public List<String> getIpAddresses() {
        return ipAddresses;
    }

    /**
     * Gets lastUpdateTime.
     * 
     * @return the lastUpdateTime
     */
    public long getLastUpdateTime() {
        return lastUpdateTime;
    }

    /**
     * Gets port.
     * 
     * @return the port
     */
    public int getPort() {
        return port;
    }

    /**
     * Gets lastUpdateHealth.
     * 
     * @return the lastUpdateHealth
     */
    public boolean isLastUpdateHealthy() {
        return lastUpdateHealthy;
    }

    /**
     * Sets adminStatus.
     * 
     * @param adminStatus
     *            the adminStatus to set
     */
    public void setAdminStatus(final String adminStatus) {
        this.adminStatus = adminStatus;
    }

    /**
     * Sets cacheId.
     * 
     * @param cacheId
     *            the cacheId to set
     */
    public void setCacheId(final String cacheId) {
        this.cacheId = cacheId;
    }

    /**
     * Sets fqdn.
     * 
     * @param fqdn
     *            the fqdn to set
     */
    public void setFqdn(final String fqdn) {
        this.fqdn = fqdn;
    }

    /**
     * Sets lastUpdateHealthy.
     * 
     * @param lastUpdateHealthy
     *            the lastUpdateHealthy to set
     */
    public void setLastUpdateHealthy(final boolean lastUpdateHealthy) {
        this.lastUpdateHealthy = lastUpdateHealthy;
    }

    /**
     * Sets ipAddresses.
     * 
     * @param ipAddresses
     *            the ipAddresses to set
     */
    public void setIpAddresses(final List<String> ipAddresses) {
        this.ipAddresses = ipAddresses;
    }

    /**
     * Sets lastUpdateTime.
     * 
     * @param lastUpdateTime
     *            the lastUpdateTime to set
     */
    public void setLastUpdateTime(final long lastUpdateTime) {
        this.lastUpdateTime = lastUpdateTime;
    }

    /**
     * Sets port.
     * 
     * @param port
     *            the port to set
     */
    public void setPort(final int port) {
        this.port = port;
    }

	public void setConnections(final long numConn) {
		this.connections = numConn;
	}

	public long getConnections() {
		return connections;
	}

	public long getCurrentBW() {
		return currentBW;
	}

	public long getAvailBW() {
		return availBW;
	}

	public void setCurrentBW(final long currBW) {
		this.currentBW = currBW;
	}

	public void setAvailBW(final long availBW) {
		this.availBW = availBW;
	}

	public void setCacheOnline(final boolean cacheOnline) {
		this.cacheOnline = cacheOnline;
	}
	public boolean isCacheOnline() {
		return cacheOnline;
	}
}
