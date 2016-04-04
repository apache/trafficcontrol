/*
 * Copyright 2015 Comcast Cable Communications Management, LLC
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

package com.comcast.cdn.traffic_control.traffic_router.core;

import org.apache.log4j.ConsoleAppender;
import org.apache.log4j.Level;
import org.apache.log4j.LogManager;
import org.apache.log4j.PatternLayout;

public class TrafficRouterStart {

	public static void main(String[] args) throws Exception {
		System.setProperty("deploy.dir", "src/test");
		System.setProperty("dns.zones.dir", "src/test/var/auto-zones");

		System.setProperty("dns.tcp.port", "1053");
		System.setProperty("dns.udp.port", "1053");

		LogManager.getLogger("org.eclipse.jetty").setLevel(Level.WARN);
		LogManager.getLogger("org.springframework").setLevel(Level.WARN);

		ConsoleAppender consoleAppender = new ConsoleAppender(new PatternLayout("%d{ISO8601} [%-5p] %c{4}: %m%n"));
		LogManager.getRootLogger().addAppender(consoleAppender);
		LogManager.getRootLogger().setLevel(Level.INFO);

		EmbeddedTrafficRouter embeddedTrafficRouter = new EmbeddedTrafficRouter("core/src/main");
		System.out.println(">>>>>>>>>>>>>>>> Starting Traffic Router in Embedded Tomcat");
		embeddedTrafficRouter.start();

		System.in.read();
		System.out.println(">>>>>>>>>>>>>>>> Stopping Traffic Router");
		embeddedTrafficRouter.stop();
	}
}
