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

package org.apache.traffic_control.traffic_router.core.monitor;

import org.apache.traffic_control.traffic_router.core.util.ResourceUrl;

class TrafficMonitorResourceUrl implements ResourceUrl {
	private final TrafficMonitorWatcher trafficMonitorWatcher;
	private final String urlTemplate;
	private int i = 0;
	public TrafficMonitorResourceUrl(final TrafficMonitorWatcher trafficMonitorWatcher, final String urlTemplate) {
		this.trafficMonitorWatcher = trafficMonitorWatcher;
		this.urlTemplate = urlTemplate;
	}
	@Override
	public String nextUrl() {
		final String[] hosts = trafficMonitorWatcher.getHosts();

		if (hosts == null || hosts.length == 0) {
			return urlTemplate;
		}

		i %= hosts.length;
		final String host = hosts[i];
		i++;
		return urlTemplate.replace("[host]", host);
	}
}
