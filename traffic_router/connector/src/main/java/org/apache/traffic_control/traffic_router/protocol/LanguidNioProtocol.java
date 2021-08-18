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

import org.apache.coyote.http11.AbstractHttp11JsseProtocol;
import org.apache.juli.logging.Log;
import org.apache.tomcat.util.net.NioChannel;
import org.bouncycastle.jce.provider.BouncyCastleProvider;
import java.security.Security;


public class LanguidNioProtocol extends AbstractHttp11JsseProtocol<NioChannel> implements RouterProtocolHandler {
	protected static org.apache.juli.logging.Log log = org.apache.juli.logging.LogFactory.getLog(LanguidNioProtocol.class);
	private boolean ready = false;
	private boolean initialized = false;
	private String mbeanPath;
	private String readyAttribute;
	private String portAttribute;


	//add BouncyCastle provider to support converting PKCS1 to PKCS8 since OpenSSL does not support PKCS1
	//TODO:  Figure out if we can convert from PKCS1 to PKCS8 with out BC
	static { log.warn("Adding BouncyCastle provider");
			Security.addProvider(new BouncyCastleProvider());
	}

	public LanguidNioProtocol() {
		super(new RouterNioEndpoint());
		log.warn("Serving wildcard certs for multiple domains");
	}

	@Override
	public void setSslImplementationName(final String sslClassName) {
		try {
			Class.forName(sslClassName);
			log.info("setSslImplementation: "+sslClassName);
			super.setSslImplementationName(sslClassName);
		} catch (ClassNotFoundException e) {
			log.error("LanguidNIOProtocol: Failed to set SSL implementation to " + sslClassName + " class was not found, defaulting to OpenSSL");
		}

	}

	@Override
	@SuppressWarnings("PMD.SignatureDeclareThrowsException")
	public void init() throws Exception {

		if (!isReady()) {
			log.info("Init called; creating thread to monitor the state of Traffic Router");
			new LanguidPoller(this).start();
			return;
		}

		log.info("Traffic Router SSL Protocol is ready; calling super.init()");
		getEndpoint().setBindOnInit(false);
		super.init();
		setInitialized(true);
	}

	@Override
	@SuppressWarnings("PMD.SignatureDeclareThrowsException")
	public void start() throws Exception {
		log.info("LanguidNioProtocol Handler Start called; waiting for initialization to occur");

		while (!isInitialized()) {
			try {
				Thread.sleep(100);
			} catch (InterruptedException e) {
				log.info("interrupted waiting for initialization");
			}
		}

		log.info("LanguidNioProtocol Handler Initialization complete; calling super.start()");

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

	@Override
	protected String getSslImplementationShortName() {
		return "openssl";
	}

	@Override
	protected String getNamePrefix() {
		if (isSSLEnabled()) {
			return ("https-" + getSslImplementationShortName()+ "-nio");
		} else {
			return ("http-nio");
		}
	}

	@Override
	protected Log getLog() { return log; }
}
