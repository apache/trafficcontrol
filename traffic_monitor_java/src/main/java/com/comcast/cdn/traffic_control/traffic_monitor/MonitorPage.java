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

package com.comcast.cdn.traffic_control.traffic_monitor;

import org.apache.wicket.markup.head.IHeaderResponse;
import org.apache.wicket.markup.head.JavaScriptHeaderItem;
import org.apache.wicket.markup.html.WebPage;

public class MonitorPage extends WebPage {
	private static final long serialVersionUID = 1L;

	public MonitorSession getMonitorSession() {
		return (MonitorSession) getSession();
	}

	@Override
	public void renderHead(final IHeaderResponse response) {
		super.renderHead(response);
		
		response.render(JavaScriptHeaderItem.forUrl("/js/main.js"));
	}
}
