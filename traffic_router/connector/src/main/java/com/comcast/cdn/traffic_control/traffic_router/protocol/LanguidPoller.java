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

package com.comcast.cdn.traffic_control.traffic_router.protocol;

import java.lang.management.ManagementFactory;

import javax.management.MBeanServer;
import javax.management.ObjectName;

public class LanguidPoller extends Thread {
	protected static org.apache.juli.logging.Log log = org.apache.juli.logging.LogFactory.getLog(LanguidPoller.class);
	final private LanguidProtocol languidProtocol;

	public LanguidPoller(final LanguidProtocol languidProtocol) {
		this.languidProtocol = languidProtocol;
	}

	@Override
	public void run() {
		log.info("Waiting for state from mbean path " + languidProtocol.getMbeanPath());

		boolean firstTime = true;
		while (true) {
			try {
				final MBeanServer mbs = ManagementFactory.getPlatformMBeanServer();
				// See src/main/opt/conf/server.xml
				// This is calling traffic-router:name=languidState
				final ObjectName languidState = new ObjectName(languidProtocol.getMbeanPath());
				final Object readyValue = mbs.getAttribute(languidState, languidProtocol.getReadyAttribute());
				final Object portValue = mbs.getAttribute(languidState, languidProtocol.getPortAttribute());
				final boolean ready = Boolean.parseBoolean(readyValue.toString());
				final int port = Integer.parseInt(portValue.toString());

				if (firstTime) {
					log.info("Waiting for ready state from Traffic Router before accepting connections on port " + port);
				}
				
				if (ready) {
					if (port > 0) {
						languidProtocol.setPort(port);
					}

					log.info("Traffic Router published the ready state; calling init() on our reference to Connector with a listen port of " + languidProtocol.getPort());
					languidProtocol.setReady(true);
					languidProtocol.init();
					break;
				}
			} catch (Exception ex) {
				// the above will throw an exception if the mbean has yet to be published
				log.debug(ex);
			}

			try {
				Thread.sleep(100);
			} catch (InterruptedException ex) {
				log.fatal(ex);
			}
			firstTime = false;
		}
	}
}
