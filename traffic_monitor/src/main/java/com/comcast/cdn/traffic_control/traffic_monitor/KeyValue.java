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

package com.comcast.cdn.traffic_control.traffic_monitor;


import org.apache.wicket.model.Model;

import com.comcast.cdn.traffic_control.traffic_monitor.health.AbstractState;

public class KeyValue extends Model<String> implements java.io.Serializable {
	private static final long serialVersionUID = 1L;
	final protected String key;
	protected final String val;
	protected final String stateId;

	public KeyValue(final String key, final String val) {
		this.key = key;
		this.val = val;
		this.stateId = null;
	}

	public KeyValue(final String key, final AbstractState cacheState) {
		this.key = key;
		this.val = null;
		this.stateId = cacheState.getId();
	}

	public String getKey() {
		return key;
	}

	@Override
	public String getObject( ) {
		return val;
	}
}
