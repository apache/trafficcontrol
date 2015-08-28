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

import java.util.Calendar;
import java.util.List;

import org.apache.log4j.Logger;
import org.xbill.DNS.Name;
import org.xbill.DNS.Record;

public class SignedZoneKey extends ZoneKey {
	private static final Logger LOGGER = Logger.getLogger(SignedZoneKey.class);
	private Calendar expiration;

	public SignedZoneKey(final Name name, final List<Record> records) {
		// sorting of records takes place in the ZoneKey constructor
		super(name, records);
	}

	public Calendar getExpiration() {
		return expiration;
	}

	public void setExpiration(final Calendar expiration) {
		this.expiration = expiration;
	}

	public long getSignatureDuration() {
		return this.expiration.getTimeInMillis() - getTimestamp();
	}

	public long getRefreshHorizon() {
		return getTimestamp() + Math.round((double) getSignatureDuration() / 2.0); // force a refresh when we're halfway through our validity period
	}
}
