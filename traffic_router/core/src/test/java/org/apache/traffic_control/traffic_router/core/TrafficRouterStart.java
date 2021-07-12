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

import org.apache.log4j.ConsoleAppender;
import org.apache.log4j.Level;
import org.apache.log4j.LogManager;
import org.apache.log4j.PatternLayout;

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

		LogManager.getLogger("org.springframework").setLevel(Level.WARN);

		ConsoleAppender consoleAppender = new ConsoleAppender(new PatternLayout("%d{ISO8601} [%-5p] %c{4}: %m%n"));
		LogManager.getRootLogger().addAppender(consoleAppender);
		LogManager.getRootLogger().setLevel(Level.INFO);

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
