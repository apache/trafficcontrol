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


import org.apache.wicket.model.Model;

public class CacheDataModel extends Model<String> {
	private static final long serialVersionUID = 1L;
	private final String label;
	long i = 0;

	public CacheDataModel(final String label) {
		this.label = label;

		if (label == null) {
			super.setObject(null);
		} else {
			super.setObject(label + ": ");
		}
	}

	public String getKey() {
		return label;
	}

	public String getValue() {
		return String.valueOf(i);
	}

	public long getRawValue() {
		return i;
	}

	public void inc() {
		synchronized (this) {
			i++;
			this.set(i);
		}
	}

	public void setObject(final String o) {
		if (label == null) {
			super.setObject(o);
		} else {
			super.setObject(label + ": " + o);
		}
	}

	public void set(final long arg) {
		synchronized (this) {
			i = arg;
			this.setObject(String.valueOf(arg));
		}
	}
}
