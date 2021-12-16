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

package org.apache.traffic_control.traffic_router.core.dns;

import java.util.Calendar;
import java.util.List;
import java.util.OptionalLong;

import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.xbill.DNS.Name;
import org.xbill.DNS.RRSIGRecord;
import org.xbill.DNS.Record;

public class SignedZoneKey extends ZoneKey {
	private static final Logger LOGGER = LogManager.getLogger(SignedZoneKey.class);

	private Calendar minimumSignatureExpiration;
	private Calendar kskExpiration;
	private Calendar zskExpiration;

	public SignedZoneKey(final Name name, final List<Record> records) {
		// sorting of records takes place in the ZoneKey constructor
		super(name, records);
	}

	public Calendar getMinimumSignatureExpiration() {
		return minimumSignatureExpiration;
	}

	public void setMinimumSignatureExpiration(final List<Record> signedRecords, final Calendar defaultExpiration) {
		final OptionalLong minSignatureExpiration = signedRecords.stream()
				.filter(r -> r instanceof RRSIGRecord)
				.mapToLong(r -> ((RRSIGRecord) r).getExpire().getTime())
				.min();
		if (!minSignatureExpiration.isPresent()) {
			LOGGER.error("unable to calculate minimum signature expiration: no RRSIG records given");
			this.minimumSignatureExpiration = defaultExpiration;
			return;
		}
		final Calendar tmp = Calendar.getInstance();
		tmp.setTimeInMillis(minSignatureExpiration.getAsLong());
		this.minimumSignatureExpiration = tmp;
	}

	public long getSignatureDuration() {
		return this.minimumSignatureExpiration.getTimeInMillis() - getTimestamp();
	}

	public long getRefreshHorizon() {
		return getTimestamp() + Math.round((double) getSignatureDuration() / 2.0); // force a refresh when we're halfway through our validity period
	}

	public long getEarliestSigningKeyExpiration() {
		if (getKSKExpiration().before(getZSKExpiration())) {
			return getKSKExpiration().getTimeInMillis();
		} else {
			return getZSKExpiration().getTimeInMillis();
		}
	}

	public Calendar getKSKExpiration() {
		return kskExpiration;
	}

	public void setKSKExpiration(final Calendar kskExpiration) {
		this.kskExpiration = kskExpiration;
	}

	public Calendar getZSKExpiration() {
		return zskExpiration;
	}

	public void setZSKExpiration(final Calendar zskExpiration) {
		this.zskExpiration = zskExpiration;
	}
}
