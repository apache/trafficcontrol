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

import java.text.SimpleDateFormat;
import java.util.Calendar;
import java.util.Date;
import java.util.List;

import org.apache.traffic_control.traffic_router.core.util.JsonUtils;
import org.apache.traffic_control.traffic_router.core.util.JsonUtilsException;
import com.fasterxml.jackson.databind.JsonNode;
import org.xbill.DNS.Record;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

public class ZoneUtils {
	private static final Logger LOGGER = LogManager.getLogger(ZoneUtils.class);
	private static final SimpleDateFormat sdf = new SimpleDateFormat("yyyyMMddHH");

	protected static long getMaximumTTL(final List<Record> records) {
		long maximumTTL = 0;

		for (final Record record : records) {
			if (record.getTTL() > maximumTTL) {
				maximumTTL = record.getTTL();
			}
		}

		return maximumTTL;
	}

	protected static long getSerial(final JsonNode jo) {
		synchronized(sdf) {
			Date date = null;

			if (jo != null && jo.has("date")) {
				try {
					final Calendar cal = Calendar.getInstance();
					cal.setTimeInMillis(JsonUtils.getLong(jo, "date") * 1000);
					date = cal.getTime();
				} catch (JsonUtilsException ex) {
					LOGGER.error(ex, ex);
				}
			}

			if (date == null) {
				date = new Date();
			}

			return Long.parseLong(sdf.format(date)); // 2013062701
		}
	}

	protected static long getLong(final JsonNode jo, final String key, final long d) {
		if (jo == null) {
			return d;
		}

		return jo.has(key) ? jo.get(key).asLong(d) : d;
	}

	protected static String getAdminString(final JsonNode jo, final String key, final String d, final String domain) {

		if (jo == null) {
			return new StringBuffer(d).append('.').append(domain).toString();
		}

		if (!jo.has(key)) {
			return new StringBuffer(d).append('.').append(domain).toString();
		}

		// check for @ sign in string
		String admin = jo.has(key) ? jo.get(key).asText() : "";
		if (admin.contains("@")) {
			admin = admin.replace("@",".");
		} else {
			admin = new StringBuffer(admin).append('.').append(domain).toString();
		}

		return admin;

	}

}
