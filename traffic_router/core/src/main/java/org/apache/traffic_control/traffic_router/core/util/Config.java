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

package org.apache.traffic_control.traffic_router.core.util;

import java.io.File;

@SuppressWarnings("PMD.ClassNamingConventions")
public class Config {

	private static String confDir = null;
	private static String varDir = null;
	static {
		confDir = "src/test/resources/var/";
		if(new File("/opt/traffic_router/conf").exists()) {
			confDir = "/opt/traffic_router/conf/";
		}
		varDir = "src/test/resources/var/";
		if(new File("/opt/traffic_router").exists()) {
			varDir = "/opt/traffic_router/var/";
		}
	}
	
	public static String getConfDir() {
		return confDir;
	}
	public static String getVarDir() {
		return varDir;
	}
}
