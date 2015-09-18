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

import org.xbill.DNS.Message;

import java.net.InetAddress;

// Using Josh Bloch Builder pattern so suppress these warnings.
@SuppressWarnings({"PMD.MissingStaticMethodInNonInstantiatableClass",
        "PMD.AccessorClassGeneration",
        "PMD.CyclomaticComplexity"})
public final class DNSAccessRecord {
    private final long queryInstant;
    private final InetAddress client;
    private final Message query;
    private final Message response;

    public long getQueryInstant() {
        return queryInstant;
    }

    public InetAddress getClient() {
        return client;
    }

    public Message getQuery() {
        return query;
    }

    public Message getResponse() {
        return response;
    }

    public static class Builder {
        private final long queryInstant;
        private final InetAddress client;
        private Message query;
        private Message response;

        public Builder(final long queryInstant, final InetAddress client) {
            this.queryInstant = queryInstant;
            this.client = client;
        }

        public Builder query(final Message query) {
            this.query = query;
            return this;
        }

        public Builder response(final Message response) {
            this.response = response;
            return this;
        }

        public DNSAccessRecord build() {
            return new DNSAccessRecord(this);
        }
    }

    private DNSAccessRecord(final Builder builder) {
        queryInstant = builder.queryInstant;
        client = builder.client;
        query = builder.query;
        response = builder.response;
    }

}
