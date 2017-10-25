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


import java.util.HashMap;
import java.util.Map;

public class EmbeddedStati implements java.io.Serializable {
	private static final long serialVersionUID = 1L;
	private DsStati currentDtati;
	private final String id;
	private StatType statType;

	public enum StatType {
		LOCATION,
		CACHE,
		TYPE
	};

	public EmbeddedStati(final StatType statType, final String id, final String delimiter) {
		final StringBuilder statId = new StringBuilder();

		statId.append(statType.toString().toLowerCase());
		statId.append(delimiter);
		statId.append(id);

		this.id = statId.toString();
		this.statType = statType;
	}

	public EmbeddedStati(final StatType statType, final String id) {
		this(statType, id, ".");
	}

	public void accumulate(final DsStati stati) {
		if (currentDtati == null) {
			currentDtati = new DsStati(stati);
		} else {
			currentDtati.accumulate(stati);
		}
	}

	public Map<String, String> completeRound() {
		if (currentDtati == null) {
			return null;
		}

		final Map<String, String> r = new HashMap<String, String>();

		r.putAll(currentDtati.getStati(this.getId()));
		currentDtati = null;

		return r;
	}

	public String getId() {
		return id;
	}

	public boolean isHidden() {
		return (statType == StatType.CACHE) ? true : false;
	}

	public StatType getStatType() {
		return statType;
	}
}
