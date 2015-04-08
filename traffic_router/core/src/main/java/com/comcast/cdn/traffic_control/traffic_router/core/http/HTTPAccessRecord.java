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

package com.comcast.cdn.traffic_control.traffic_router.core.http;

import java.net.URL;
import java.util.Date;
import java.util.TimeZone;

import javax.servlet.http.HttpServletRequest;

import org.apache.commons.lang3.time.FastDateFormat;
import org.apache.log4j.Logger;

public class HTTPAccessRecord {
    private static final Logger ACCESS = Logger.getLogger("com.comcast.cdn.traffic_control.traffic_router.core.access");
    private static final String ACCESS_FORMAT = "HTTP [%s] %s %s %s %s";
    private static final FastDateFormat FORMATTER = FastDateFormat.getInstance("dd/MMM/yyyy:HH:mm:ss.SSS Z",
            TimeZone.getTimeZone("GMT"));
    private static final String EMPTY_FIELD = "-";

    private Date requestDate;
    private HttpServletRequest request;
    private int responseCode;
    private URL responseURL;

    /**
     * Logs the access record.
     */
    public void log() {
        ACCESS.info(String.format(ACCESS_FORMAT, formatAccessDate(), formatRequestIP(), formatRequestURL(),
                formatResponseCode(), formatResponseURL()));
    }

    /**
     * Sets request.
     * 
     * @param request
     *            the request to set
     */
    public void setRequest(final HttpServletRequest request) {
        this.request = request;
    }

    /**
     * Sets requestDate.
     * 
     * @param requestDate
     *            the requestDate to set
     */
    public void setRequestDate(final Date requestDate) {
        this.requestDate = new Date(requestDate.getTime());
    }

    /**
     * Sets responseCode.
     * 
     * @param responseCode
     *            the responseCode to set
     */
    public void setResponseCode(final int responseCode) {
        this.responseCode = responseCode;
    }

    /**
     * Sets responseURL.
     * 
     * @param responseURL
     *            the responseURL to set
     */
    public void setResponseURL(final URL responseURL) {
        this.responseURL = responseURL;
    }

    /**
     * Formats the access date using the request date.
     * 
     * @return the formatted date
     */
    private String formatAccessDate() {
        String accessDate = EMPTY_FIELD;
        if (requestDate != null) {
            accessDate = FORMATTER.format(requestDate);
        }
        return accessDate;
    }

    /**
     * Formats the requestor's IP.
     * 
     * @return the formatted IP
     */
    private String formatRequestIP() {
        String ip = EMPTY_FIELD;
        if ((request != null) && (request.getRemoteAddr() != null)) {
            ip = request.getRemoteAddr();
        }
        return ip;
    }

    /**
     * Formats the request URL.
     * 
     * @return the formatted request URL
     */
    private String formatRequestURL() {
        String url = EMPTY_FIELD;
        if (request != null) {
            final StringBuffer buf = request.getRequestURL();
            if (request.getQueryString() != null) {
                buf.append('?');
                buf.append(request.getQueryString());
            }
            url = buf.toString();
        }
        return url;
    }

    /**
     * Formats the response code.
     * 
     * @return the formatted response code
     */
    private String formatResponseCode() {
        String code = EMPTY_FIELD;
        if (responseCode > 0) {
            code = String.valueOf(responseCode);
        }
        return code;
    }

    /**
     * Formats the response URL.
     * 
     * @return the formatted response URL
     */
    private String formatResponseURL() {
        String url = EMPTY_FIELD;
        if (responseURL != null) {
            url = responseURL.toString();
        }
        return url;
    }
}
