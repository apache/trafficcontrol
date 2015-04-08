/*
 * Copyright 2015 Comcast Cable Communications Management, LLC
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

package com.comcast.cdn.traffic_control.traffic_router.core.request;

import java.util.Map;

import org.apache.commons.lang3.builder.EqualsBuilder;
import org.apache.commons.lang3.builder.HashCodeBuilder;

public class HTTPRequest extends Request {

    private String requestedUrl;
    private String path;
    private String queryString;
    private Map<String, String> headers;

    @Override
    public boolean equals(final Object obj) {
        if (this == obj) {
            return true;
        } else if (obj instanceof HTTPRequest) {
            final HTTPRequest rhs = (HTTPRequest) obj;
            return new EqualsBuilder()
                    .appendSuper(super.equals(obj))
                    .append(getHeaders(), rhs.getHeaders())
                    .append(getPath(), rhs.getPath())
                    .append(getQueryString(), rhs.getQueryString())
                    .isEquals();
        } else {
            return false;
        }
    }

    public Map<String, String> getHeaders() {
        return headers;
    }

    public String getPath() {
        return path;
    }

    public String getQueryString() {
        return queryString;
    }

    /**
     * Gets the requested URL. This URL will not include the query string if the client provided
     * one.
     * 
     * @return the requestedUrl
     */
    public String getRequestedUrl() {
        return requestedUrl;
    }

    @Override
    public int hashCode() {
        return new HashCodeBuilder(1, 31)
                .appendSuper(super.hashCode())
                .append(getHeaders())
                .append(getPath())
                .append(getQueryString())
                .toHashCode();
    }

    public void setHeaders(final Map<String, String> headers) {
        this.headers = headers;
    }

    public void setPath(final String path) {
        this.path = path;
    }

    public void setQueryString(final String queryString) {
        this.queryString = queryString;
    }

    /**
     * Sets the requested URL. This URL SHOULD NOT include the query string if the client provided
     * one.
     * 
     * @param requestedUrl
     *            the requestedUrl to set
     */
    public void setRequestedUrl(final String requestedUrl) {
        this.requestedUrl = requestedUrl;
    }

}
