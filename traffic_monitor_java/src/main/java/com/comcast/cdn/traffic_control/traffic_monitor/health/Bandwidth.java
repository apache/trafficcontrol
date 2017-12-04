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

package com.comcast.cdn.traffic_control.traffic_monitor.health;

public class Bandwidth {
	public static final int BITS_IN_BYTE = 8;
	public static final int BITS_IN_KBPS = 1000;
//	public static final int MS_IN_SEC = 1000;
	public final long timeInMS;
	public final long bits;

	public Bandwidth(final String bytes) {
		timeInMS = System.currentTimeMillis();
		bits = Long.parseLong(bytes) * BITS_IN_BYTE;
	}
	public Bandwidth(final long bytes) {
		timeInMS = System.currentTimeMillis();
		bits = bytes * BITS_IN_BYTE;
	}

	public double calculateKbps(final Bandwidth current) {
		double result = 0.0;
		// as long as the numbers are not too large, dividing both num and denom by 1000 is a waste of time
		final double tDelta = ((current.timeInMS - timeInMS));// / MS_IN_SEC);
		if (tDelta > 0.0) {
			final double bitDiff = (current.bits - bits);
			result = bitDiff / tDelta;
		}
		return Math.max(0.0, result);
	}
}
