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

package com.comcast.cdn.traffic_control.traffic_router.logger;

import java.util.ArrayList;
import java.util.Collections;
import java.util.List;

import org.apache.log4j.AppenderSkeleton;
import org.apache.log4j.spi.LoggingEvent;

public class NoLogAppender extends AppenderSkeleton {
    private static final String NAME = "NoLogAppender";

    private final List<LoggingEvent> events = new ArrayList<LoggingEvent>();
    private boolean closed = false;

    public NoLogAppender() {
        setName(NAME);
    }

    public NoLogAppender(final boolean isActive) {
        super(isActive);
        setName(NAME);
    }

    /**
     * Clears out all previous events from the appender.
     */
    public void clear() {
        events.clear();
    }

    @Override
    public void close() {
        closed = true;
    }

    /**
     * Gets events.
     *
     * @return the events
     */
    public List<LoggingEvent> getEvents() {
        return Collections.unmodifiableList(events);
    }

    @Override
    public boolean requiresLayout() {
        return false;
    }

    @Override
    protected void append(final LoggingEvent event) {
        if (closed) {
            throw new IllegalStateException("Cannot log to a closed appender.");
        }
        events.add(event);
    }

}
