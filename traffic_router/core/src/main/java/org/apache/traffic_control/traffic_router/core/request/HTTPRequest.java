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

package org.apache.traffic_control.traffic_router.core.request;

import java.net.URL;
import java.util.Enumeration;
import java.util.HashMap;
import java.util.Map;

import org.apache.commons.lang3.builder.EqualsBuilder;
import org.apache.commons.lang3.builder.HashCodeBuilder;

import javax.servlet.http.HttpServletRequest;

public class HTTPRequest extends Request {
    public static final String X_MM_CLIENT_IP = "X-MM-Client-IP";
    public static final String FAKE_IP = "fakeClientIpAddress";

    private String requestedUrl;
    private String path;
    private String uri;
    private String queryString;
    private Map<String, String> headers;
    private boolean secure = false;

    public HTTPRequest() { }

    @SuppressWarnings("PMD.ConstructorCallsOverridableMethod")
    public HTTPRequest(final HttpServletRequest request) {
        applyRequest(request);
    }

    @SuppressWarnings("PMD.ConstructorCallsOverridableMethod")
    public HTTPRequest(final HttpServletRequest request, final URL url) {
        applyRequest(request);
        applyUrl(url);
    }

    @SuppressWarnings("PMD.ConstructorCallsOverridableMethod")
    public HTTPRequest(final URL url) {
        applyUrl(url);
    }

    public void applyRequest(final HttpServletRequest request) {
        setClientIP(request.getRemoteAddr());
        setPath(request.getPathInfo());
        setQueryString(request.getQueryString());
        setHostname(request.getServerName());
        setRequestedUrl(request.getRequestURL().toString());
        setUri(request.getRequestURI());

        final String xmm = request.getHeader(X_MM_CLIENT_IP);
        final String fip = request.getParameter(FAKE_IP);

        if (xmm != null) {
            setClientIP(xmm);
        } else if (fip != null) {
            setClientIP(fip);
        }

        final Map<String, String> headers = new HashMap<String, String>();
        final Enumeration<?> headerNames = request.getHeaderNames();
        while (headerNames.hasMoreElements()) {
            final String name = (String) headerNames.nextElement();
            final String value = request.getHeader(name);
            headers.put(name, value);
        }
        setHeaders(headers);
        secure = request.isSecure();
    }

    public void applyUrl(final URL url) {
        setPath(url.getPath());
        setQueryString(url.getQuery());
        setHostname(url.getHost());
        setRequestedUrl(url.toString());
    }

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
                    .append(getUri(), rhs.getUri())
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
                .append(getUri())
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

    @Override
    public String getType() {
        return "http";
    }

    public String getUri() {
        return uri;
    }

    public void setUri(final String uri) {
       this.uri = uri;
    }

    public String getHeaderValue(final String name) {
        if (headers != null && headers.containsKey(name)) {
            return headers.get(name);
        }

        return null;
    }

    public boolean isSecure() {
        return secure;
    }
}
