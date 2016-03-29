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

import org.apache.log4j.Logger;
import org.springframework.context.ApplicationContext;
import org.springframework.context.support.FileSystemXmlApplicationContext;

public class TestBase {
	private static final Logger LOGGER = Logger.getLogger(TestBase.class);
	private static ApplicationContext context;

	public static ApplicationContext getContext() {
		System.setProperty("deploy.dir", "src/test");
		System.setProperty("dns.zones.dir", "src/test/var/auto-zones");

		if (context != null) {
			return context;
		}

		LOGGER.warn("Initializing context before running integration tests");
		context = new FileSystemXmlApplicationContext("src/main/webapp/WEB-INF/applicationContext.xml");
		LOGGER.warn("Context initialized integration tests will now start running");
		return context;
	}

}
