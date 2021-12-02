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

package org.apache.traffic_control.traffic_router.protocol;


import org.apache.coyote.http11.Http11NioProtocol;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

public class LanguidProtocol extends Http11NioProtocol implements RouterProtocolHandler {
	private static final Logger log = LogManager.getLogger(LanguidProtocol.class);
	private boolean ready = false;
	private boolean initialized = false;
	private String mbeanPath;
	private String readyAttribute;
	private String portAttribute;

	@Override
	@SuppressWarnings("PMD.SignatureDeclareThrowsException")
	public void init() throws Exception {
		if (!isReady()) {
			log.info("Init called; creating thread to monitor the state of Traffic Router");
			new LanguidPoller(this).start();
		} else {
			log.info("Traffic Router is ready; calling super.init()");
			super.init();
			setInitialized(true);
		}
	}

	@Override
	@SuppressWarnings("PMD.SignatureDeclareThrowsException")
	public void start() throws Exception {
		log.info("Start called; waiting for initialization to occur");

		while (!isInitialized()) {
			Thread.sleep(100);
		}

		log.info("Initialization complete; calling super.start()");

		super.start();
	}

	public boolean isReady() {
		return ready;
	}

	public void setReady(final boolean isReady) {
		this.ready = isReady;
	}

	public boolean isInitialized() {
		return initialized;
	}

	public void setInitialized(final boolean isInitialized) {
		this.initialized = isInitialized;
	}

	public String getMbeanPath() {
		return mbeanPath;
	}

	public void setMbeanPath(final String mbeanPath) {
		this.mbeanPath = mbeanPath;
	}

	public String getReadyAttribute() {
		return readyAttribute;
	}

	public void setReadyAttribute(final String readyAttribute) {
		this.readyAttribute = readyAttribute;
	}

	public String getPortAttribute() {
		return portAttribute;
	}

	public void setPortAttribute(final String portAttribute) {
		this.portAttribute = portAttribute;
	}
}
