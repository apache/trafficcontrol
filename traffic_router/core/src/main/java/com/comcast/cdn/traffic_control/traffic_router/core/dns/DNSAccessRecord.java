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

package com.comcast.cdn.traffic_control.traffic_router.core.dns;

import java.net.InetAddress;
import java.util.Date;
import java.util.TimeZone;

import org.apache.commons.lang3.time.FastDateFormat;
import org.apache.log4j.Logger;
import org.xbill.DNS.Message;
import org.xbill.DNS.Rcode;
import org.xbill.DNS.Record;
import org.xbill.DNS.Section;
import org.xbill.DNS.Type;

public final class DNSAccessRecord {
    private static final Logger ACCESS = Logger.getLogger("com.comcast.cdn.traffic_control.traffic_router.core.access");
    private static final String ACCESS_FORMAT = "DNS [%s] %s %s %s %s \"%s\"";
    private static final FastDateFormat FORMATTER = FastDateFormat.getInstance("dd/MMM/yyyy:HH:mm:ss.SSS Z",
            TimeZone.getTimeZone("GMT"));

    private Date requestDate;
    private InetAddress client;
    private Message request;
    private Message response;

    private String date;
    private String ip;
    private String type;
    private String req;
    private String rcode;
    private String resp;

    public DNSAccessRecord() {
        date = "-";
        ip = "-";
        type = "-";
        req = "-";
        rcode = "-";
        resp = "-";
    }

    public void log() {
        parseRequestDate();
        parseClient();
        parseRequest();
        parseResponse();
        ACCESS.info(String.format(ACCESS_FORMAT, date, ip, type, req, rcode, resp));
    }

    /**
     * Sets client.
     * 
     * @param client
     *            the client to set
     */
    public void setClient(final InetAddress client) {
        this.client = client;
    }

    /**
     * Sets request.
     * 
     * @param request
     *            the request to set
     */
    public void setRequest(final Message request) {
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
     * Sets response.
     * 
     * @param response
     *            the response to set
     */
    public void setResponse(final Message response) {
        this.response = response;
    }

    private void parseClient() {
        if ((client != null) && (client.getHostAddress() != null)) {
            ip = client.getHostAddress();
        }
    }

    private void parseRequest() {
        if ((request != null) && (request.getQuestion() != null)) {
            final Record question = request.getQuestion();
            type = Type.string(question.getType());
            if (question.getName() != null) {
                req = question.getName().toString();
            }
        }
    }

    private void parseRequestDate() {
        date = FORMATTER.format(requestDate);
    }

    private void parseResponse() {
        if (response != null) {
            if (response.getHeader() != null) {
                rcode = Rcode.string(response.getHeader().getRcode());
            }
            final StringBuilder tmpResp = new StringBuilder();
            final Record[] answers = response.getSectionArray(Section.ANSWER);
            if ((answers != null) && (answers.length > 0)) {
                for (final Record answer : answers) {
                    tmpResp.append(answer.rdataToString());
                    tmpResp.append(" ");
                }
                resp = tmpResp.toString().trim();
            }
        }
    }
}
