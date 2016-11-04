package com.comcast.cdn.traffic_control.traffic_monitor.health;

/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 * 
 *   http://www.apache.org/licenses/LICENSE-2.0
 * 
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */


import org.apache.log4j.Logger;

import java.text.DecimalFormat;
import java.util.HashMap;
import java.util.Map;

public class DsStati implements java.io.Serializable {
	private static final Logger LOGGER = Logger.getLogger(DsStati.class);
	private static final long serialVersionUID = 1L;
	long csIndex = 0;
	long in_bytes;
	long out_bytes;
	long status_2xx;
	long status_3xx;
	long status_4xx;
	long status_5xx;
	boolean error = false;

	double kbps;
	double tps_2xx;
	double tps_3xx;
	double tps_4xx;
	double tps_5xx;
	double tps_total;

	String dsId;
	String csId;

	public static final int BITS_IN_BYTE = 8;
	public static final int MS_IN_SEC = 1000;
	public final long time;

	public DsStati(final String propBase, final CacheState cacheState, final long index, final String dsId) {
		this.csIndex = index;
		this.time = cacheState.getTime(index);
		String v = cacheState.getValue(propBase + ".in_bytes", index);
		this.in_bytes = toLong(v);
		final String k = propBase + ".out_bytes";
		v = cacheState.getValue(k, index);
		this.out_bytes = toLong(v);
		v = cacheState.getValue(propBase + ".status_2xx", index);
		this.status_2xx = toLong(v);
		this.status_3xx = toLong(cacheState.getValue(propBase + ".status_3xx", index));
		this.status_4xx = toLong(cacheState.getValue(propBase + ".status_4xx", index));
		this.status_5xx = toLong(cacheState.getValue(propBase + ".status_5xx", index));
		this.dsId = dsId;
		this.csId = cacheState.getId();
	}

	public DsStati(final DsStati stati) {
		this.error = stati.error;
		this.time = stati.time;
		this.in_bytes = stati.in_bytes;
		this.out_bytes = stati.out_bytes;
		this.status_2xx = stati.status_2xx;
		this.status_3xx = stati.status_3xx;
		this.status_4xx = stati.status_4xx;
		this.status_5xx = stati.status_5xx;

		this.kbps = stati.kbps;
		this.tps_2xx = stati.tps_2xx;
		this.tps_3xx = stati.tps_3xx;
		this.tps_4xx = stati.tps_4xx;
		this.tps_5xx = stati.tps_5xx;
		this.tps_total = stati.tps_total;
	}

	private static long toLong(final String str) {
		if (str == null) {
			return 0;
		}
		return (long) Double.parseDouble(str);
	}

	void accumulate(final DsStati ds) {
		this.in_bytes += ds.in_bytes;
		this.out_bytes += ds.out_bytes;
		this.status_2xx += ds.status_2xx;
		this.status_3xx += ds.status_3xx;
		this.status_4xx += ds.status_4xx;
		this.status_5xx += ds.status_5xx;

		this.kbps += ds.kbps;
		this.tps_2xx += ds.tps_2xx;
		this.tps_3xx += ds.tps_3xx;
		this.tps_4xx += ds.tps_4xx;
		this.tps_5xx += ds.tps_5xx;
		this.tps_total += ds.tps_total;
	}

	public boolean calculateKbps(final DsStati prior) {
		if (prior == null) {
			LOGGER.warn("why is prior null");
			return false;
		}
		if (prior.time == 0) {
			LOGGER.warn("why is prior.time==0");
		}
		if ((out_bytes == 0 || prior.out_bytes == 0) && out_bytes != prior.out_bytes) {
			LOGGER.warn(dsId + ": throwing out " + csId + ": out_bytes==0");
			if (prior.out_bytes != 0) {
				LOGGER.warn("\t prior.out_bytes=" + prior.out_bytes);
			}
			return false;
		}
		final long deltaTimeMs = time - prior.time; // / MS_IN_SEC
		if (LOGGER.isDebugEnabled()) {
			LOGGER.debug(String.format("time delta: %d, index: %d -> %d", deltaTimeMs, prior.csIndex, this.csIndex));
		}
		if (deltaTimeMs == 0) {
			LOGGER.warn("time delta 0");
			return false;
		}
		final long delta = (out_bytes - prior.out_bytes);
		// as long as the numbers are not too large, dividing both num and denom by 1000 is a waste of time
		//			rates.kbps = (delta / BITS_IN_KBPS) * BITS_IN_BYTE / deltaTime;
		kbps = ((double) delta / (double) deltaTimeMs) * BITS_IN_BYTE;
		if (kbps < 0.0) {
			LOGGER.warn(dsId + ": throwing out " + csId + ": kbps=" + kbps);
			kbps = 0.0;
			return false;
		}

		final double deltaTime = (double) deltaTimeMs / (double) MS_IN_SEC;

		tps_2xx = ((double) status_2xx - (double) prior.status_2xx) / deltaTime;
		tps_3xx = ((double) status_3xx - (double) prior.status_3xx) / deltaTime;
		tps_4xx = ((double) status_4xx - (double) prior.status_4xx) / deltaTime;
		tps_5xx = ((double) status_5xx - (double) prior.status_5xx) / deltaTime;
		tps_total = tps_2xx + tps_3xx + tps_4xx + tps_5xx;

		return true;
	}

	Map<String, String> getStati(final String base) {
		final Map<String, String> r = new HashMap<String, String>();
		r.put(base + ".in_bytes", String.valueOf(in_bytes));
		r.put(base + ".out_bytes", String.valueOf(out_bytes));
		r.put(base + ".status_2xx", String.valueOf(status_2xx));
		r.put(base + ".status_3xx", String.valueOf(status_3xx));
		r.put(base + ".status_4xx", String.valueOf(status_4xx));
		r.put(base + ".status_5xx", String.valueOf(status_5xx));

		DecimalFormat df = new DecimalFormat("0.00");
		r.put(base + ".kbps", df.format(kbps));
		r.put(base + ".tps_2xx", df.format(tps_2xx));
		r.put(base + ".tps_3xx", df.format(tps_3xx));
		r.put(base + ".tps_4xx", df.format(tps_4xx));
		r.put(base + ".tps_5xx", df.format(tps_5xx));
		r.put(base + ".tps_total", df.format(tps_total));
		return r;
	}
}
