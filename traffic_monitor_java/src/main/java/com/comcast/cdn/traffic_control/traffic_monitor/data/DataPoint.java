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

package com.comcast.cdn.traffic_control.traffic_monitor.data;

public class DataPoint implements java.io.Serializable {
	private static final long serialVersionUID = 1L;
	private long index;
	private final String value;
	private int span;

	public DataPoint(final String value, final long index) {
		this.value = value;
		this.index = index;
		this.span = 1;
	}

	public String getValue() {
		return value;
	}

	public boolean matches(final String other) {
		return (other != null) ? other.equals(value) : value == null;
	}

	public long getIndex() {
		return index;
	}

	public void setIndex(final int index) {
		this.index = index;
	}

	public void update(final long index) {
		this.index = index;
		span++;
	}

	public int getSpan() {
		return span;
	}
}
