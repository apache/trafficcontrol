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

package org.apache.traffic_control.traffic_router.core.util;

public abstract class AbstractUpdatable {
	private long lastUpdated = 0;

	public abstract boolean update(String newDB);
	public abstract boolean noChange();

	public void complete() {
		// override if you wish to exec code after the download is complete
	}

	public long getLastUpdated() {
		return lastUpdated;
	}

	public void setLastUpdated(final long lastUpdated) {
		this.lastUpdated = lastUpdated;
	}
	public void cancelUpdate() {}
}
