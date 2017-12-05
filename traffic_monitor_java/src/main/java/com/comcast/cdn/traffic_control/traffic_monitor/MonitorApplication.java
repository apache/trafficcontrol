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

import org.apache.log4j.Logger;
import org.apache.wicket.Application;
import org.apache.wicket.Page;
import org.apache.wicket.Session;
import org.apache.wicket.protocol.http.WebApplication;
import org.apache.wicket.request.Request;
import org.apache.wicket.request.Response;
import org.apache.wicket.util.time.Duration;

import com.comcast.cdn.traffic_control.traffic_monitor.config.ConfigHandler;
import com.comcast.cdn.traffic_control.traffic_monitor.config.RouterConfig;
import com.comcast.cdn.traffic_control.traffic_monitor.health.CacheWatcher;
import com.comcast.cdn.traffic_control.traffic_monitor.health.DsWatcher;
import com.comcast.cdn.traffic_control.traffic_monitor.health.HealthDeterminer;
import com.comcast.cdn.traffic_control.traffic_monitor.health.PeerWatcher;
import com.comcast.cdn.traffic_control.traffic_monitor.health.TmWatcher;
import com.comcast.cdn.traffic_control.traffic_monitor.publish.CrStates;

public class MonitorApplication extends WebApplication {
	private static final Logger LOGGER = Logger.getLogger(MonitorApplication.class);
	private CacheWatcher cw;
	private TmWatcher tmw;
	private PeerWatcher pw;
	private DsWatcher dsw;
	private static long startTime;

	public static MonitorApplication get() {
		return (MonitorApplication) Application.get();
	}

	/**
	 * @see org.apache.wicket.Application#getHomePage()
	 */
	@Override
	public Class<? extends Page> getHomePage() {
		return Index.class;
	}

	/**
	 * @see org.apache.wicket.Application#init()
	 */
	@Override
	public void init() {
		super.init();

		if (!ConfigHandler.getInstance().configFileExists()) {
			LOGGER.fatal("Cannot find configuration file: " + ConfigHandler.getInstance().getConfigFile());
			// This will only stop Tomcat if the security manager allows it
			// https://tomcat.apache.org/tomcat-6.0-doc/security-manager-howto.html
			System.exit(1);
		}

		getResourceSettings().setResourcePollFrequency(Duration.ONE_SECOND);

		// This allows us to override the Host header sent via URLConnection
		System.setProperty("sun.net.http.allowRestrictedHeaders", "true");

		final HealthDeterminer hd = new HealthDeterminer();
		tmw = new TmWatcher(hd);
		cw = new CacheWatcher();
		cw.init(hd);
		dsw = new DsWatcher();
		dsw.init(hd);
		pw = new PeerWatcher();
		pw.init();
		tmw.addTmListener(RouterConfig.getTmListener(hd));
		tmw.init();

		CrStates.init(cw, pw, hd);

		mountPackage("/publish", CrStates.class);
		startTime = System.currentTimeMillis();
	}

	@Override
	public Session newSession(final Request request, final Response response) {
		return new MonitorSession(request);
	}

	public void onDestroy() {
		final boolean forceDown = ConfigHandler.getInstance().getConfig().shouldForceSystemExit();
		ConfigHandler.getInstance().destroy();
		LOGGER.warn("MonitorApplication: shutting down ");
		tmw.destroy();

		if (forceDown) {
			LOGGER.warn("MonitorApplication: System.exit");
			System.exit(0);
		}

		cw.destroy();
		dsw.destroy();
		pw.destroy();
	}

	public static long getUptime() {
		return System.currentTimeMillis() - startTime;
	}
}
