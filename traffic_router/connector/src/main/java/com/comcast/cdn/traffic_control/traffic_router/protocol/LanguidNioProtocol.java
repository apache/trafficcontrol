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

import org.apache.coyote.http11.Http11NioProtocol;
import org.apache.tomcat.util.net.SSLImplementation;
import org.apache.tomcat.util.net.jsse.JSSEImplementation;

public class LanguidNioProtocol extends Http11NioProtocol implements RouterProtocolHandler {
	protected static org.apache.juli.logging.Log log = org.apache.juli.logging.LogFactory.getLog(LanguidNioProtocol.class);
	private boolean ready = false;
	private boolean initialized = false;
	private String mbeanPath;
	private String readyAttribute;
	private String portAttribute;
	private String sslClassName = JSSEImplementation.class.getCanonicalName();

	public LanguidNioProtocol() {
		log.warn("Serving wildcard certs for multiple domains");
		ep = new RouterNioEndpoint();
		setSSLImplementation(RouterSslImplementation.class.getCanonicalName());
	}

	public boolean setSSLImplementation(final String sslClassName) {
		try {
			Class.forName(sslClassName);
			this.sslClassName = sslClassName;
			return true;
		} catch (ClassNotFoundException e) {
			log.error("Failed to set SSL implementation to " + sslClassName + " class was not found, defaulting to JSSE");
			this.sslClassName = JSSEImplementation.class.getCanonicalName();
		}

		return false;
	}

	@Override
	@SuppressWarnings("PMD.SignatureDeclareThrowsException")
	public void init() throws Exception {
		if (!isReady()) {
			log.info("Init called; creating thread to monitor the state of Traffic Router");
			new LanguidPoller(this).start();
			return;
		}

		log.info("Traffic Router is ready; calling super.init()");
		super.init();

		sslImplementation = (JSSEImplementation) SSLImplementation.getInstance(sslClassName);
		setInitialized(true);
	}

	@Override
	@SuppressWarnings("PMD.SignatureDeclareThrowsException")
	public void start() throws Exception {
		log.info("Start called; waiting for initialization to occur");

		while (!isInitialized()) {
			try {
				Thread.sleep(100);
			} catch (InterruptedException e) {
				log.info("interrupted waiting for initialization");
			}
		}

		log.info("Initialization complete; calling super.start()");

		super.start();
	}

	@Override
	public boolean isReady() {
		return ready;
	}

	@Override
	public void setReady(final boolean isReady) {
		this.ready = isReady;
	}

	@Override
	public boolean isInitialized() {
		return initialized;
	}

	@Override
	public void setInitialized(final boolean isInitialized) {
		this.initialized = isInitialized;
	}

	@Override
	public String getMbeanPath() {
		return mbeanPath;
	}

	@Override
	public void setMbeanPath(final String mbeanPath) {
		this.mbeanPath = mbeanPath;
	}

	@Override
	public String getReadyAttribute() {
		return readyAttribute;
	}

	@Override
	public void setReadyAttribute(final String readyAttribute) {
		this.readyAttribute = readyAttribute;
	}

	@Override
	public String getPortAttribute() {
		return portAttribute;
	}

	@Override
	public void setPortAttribute(final String portAttribute) {
		this.portAttribute = portAttribute;
	}
}
