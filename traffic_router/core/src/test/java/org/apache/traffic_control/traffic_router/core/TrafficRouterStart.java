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

package org.apache.traffic_control.traffic_router.core;

import org.apache.logging.log4j.core.Appender;
import org.apache.logging.log4j.core.LoggerContext;
import org.apache.logging.log4j.core.appender.ConsoleAppender;
import org.apache.logging.log4j.core.config.Configuration;
import org.apache.logging.log4j.core.layout.PatternLayout;
import org.apache.logging.log4j.Level;
import org.apache.logging.log4j.LogManager;

import java.util.Collection;

public class TrafficRouterStart {

	public static void main(String[] args) throws Exception {
		String prefix = System.getProperty("user.dir");

		if (!prefix.endsWith("/core")) {
			prefix += "/core";
		}

		System.setProperty("dns.zones.dir", prefix + "/src/test/var/auto-zones");
		System.setProperty("deploy.dir", prefix + "/src/test");

		System.setProperty("dns.tcp.port", "1053");
		System.setProperty("dns.udp.port", "1053");

		LoggerContext ctx = (LoggerContext) LogManager.getContext(false);
		Configuration config = ctx.getConfiguration();
		config.getLoggerConfig("org.springframework").setLevel(Level.WARN);

		Collection<Appender> rootAppenders = ctx.getRootLogger().getAppenders().values();
		for (Appender a : rootAppenders) {
			ctx.getRootLogger().removeAppender(a);
		}
		ConsoleAppender consoleAppender = ConsoleAppender.newBuilder().setName("TrafficRouterStart").setLayout(PatternLayout.newBuilder().withPattern("%d{ISO8601} [%-5p] %c{4}: %m%n").build()).build();
		ctx.getRootLogger().addAppender(consoleAppender);
		config.getLoggerConfig(LogManager.ROOT_LOGGER_NAME).setLevel(Level.INFO);
		ctx.updateLoggers();

		System.out.println("[" + System.currentTimeMillis() + "] >>>>>>>>>>>>>>>> Embedded Tomcat loading Traffic Router");
		CatalinaTrafficRouter catalinaTrafficRouter = new CatalinaTrafficRouter(prefix + "/src/main/conf/server.xml", prefix + "/src/main/webapp" );
		System.out.println("[" + System.currentTimeMillis() + "] >>>>>>>>>>>>>>>> Starting Traffic Router");
		catalinaTrafficRouter.start();
		System.out.println("[" + System.currentTimeMillis() + "] >>>>>>>>>>>>>>>> Traffic Router started, press q and <ENTER> to stop");

		while ('q' != System.in.read()) {
			System.out.println("[" + System.currentTimeMillis() + "] >>>>>>>>>>>>>>> press q and <ENTER> to stop");
		}

		System.out.println("[" + System.currentTimeMillis() + "] >>>>>>>>>>>>>>>> Stopping Traffic Router");
		catalinaTrafficRouter.stop();
	}
}
