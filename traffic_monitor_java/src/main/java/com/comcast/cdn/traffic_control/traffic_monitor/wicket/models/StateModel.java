package com.comcast.cdn.traffic_control.traffic_monitor.wicket.models;

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


import com.comcast.cdn.traffic_control.traffic_monitor.health.AbstractState;
import org.apache.wicket.model.Model;

public class StateModel extends Model<String> {
	String stateName;
	String key;

	public StateModel(final String stateName, final String key) {
		this.stateName = stateName;
		this.key = key;
	}

	protected String getObject(AbstractState state) {
		if (state == null) {
			return "err";
		}
		if ("_status_string_".equals(key)) {
			return state.getStatusString();
		}
		final boolean clearData = state.getBool("clearData");
		if (clearData) {
			return "-";
		}
		return state.getLastValue(key);
	}
}
