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

package com.comcast.cdn.traffic_control.traffic_router.core.dns;

import java.text.SimpleDateFormat;
import java.util.Calendar;
import java.util.Date;
import java.util.List;

import org.apache.log4j.Logger;
import org.json.JSONException;
import org.json.JSONObject;
import org.xbill.DNS.Record;

public class ZoneUtils {
	private static final Logger LOGGER = Logger.getLogger(ZoneUtils.class);
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

	protected static long getSerial(final JSONObject jo) {
		synchronized(sdf) {
			Date date = null;

			if (jo != null && jo.has("date")) {
				try {
					final Calendar cal = Calendar.getInstance();
					cal.setTimeInMillis(jo.getLong("date") * 1000);
					date = cal.getTime();
				} catch (JSONException ex) {
					LOGGER.error(ex, ex);
				}
			}

			if (date == null) {
				date = new Date();
			}

			return Long.parseLong(sdf.format(date)); // 2013062701
		}
	}

	protected static long getLong(final JSONObject jo, final String key, final long d) {
		if (jo == null) {
			return d;
		}

		if (!jo.has(key)) {
			return d;
		}

		return jo.optLong(key);
	}

	protected static String getAdminString(final JSONObject jo, final String key, final String d, final String domain) {

		if (jo == null) {
			return new StringBuffer(d).append(".").append(domain).toString();
		}

		if (!jo.has(key)) {
			return new StringBuffer(d).append(".").append(domain).toString();
		}

		// check for @ sign in string
		String admin = jo.optString(key);
		if (admin.contains("@")) {
			admin = admin.replace("@",".");
		} else {
			admin = new StringBuffer(admin).append(".").append(domain).toString();
		}

		return admin;

	}

}
