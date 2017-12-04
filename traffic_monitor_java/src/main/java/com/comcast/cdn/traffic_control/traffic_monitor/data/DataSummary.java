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

public class DataSummary implements java.io.Serializable {
	private static final long serialVersionUID = 1L;
	private long startTime;
	private long endTime;
	private double average;
	private double high;
	private double low;
	private double start;
	private double end;
	private int dpCount;
	public DataSummary() {
	}
	public void accumulate(final DataPoint dp, final long t) {
		final double v = Double.parseDouble(dp.getValue());
		if(dpCount == 0) {
			startTime = t;
			endTime = t;
			high = v;
			low = v;
			start = v;
			end = v;
			average = v;
		} else {
			if(t > endTime) {
				endTime = t;
			} else if(t < startTime) {
				startTime = t;
			}
			if(v > high) {
				high = v;
			} else if(v < low) {
				low = v;
			}
			// a = a' + (v-a')/(c'+1)
			end = v;
			average = average + (v-average)/(dpCount+1);
		}
		dpCount++;
	}
	public long getStartTime() {
		return startTime;
	}
	public long getEndTime() {
		return endTime;
	}
	public double getAverage() {
		return average;
	}
	public double getHigh() {
		return high;
	}
	public double getLow() {
		return low;
	}
	public int getDpCount() {
		return dpCount;
	}
	public double getStart() {
		return start;
	}
	public double getEnd() {
		return end;
	}
}
