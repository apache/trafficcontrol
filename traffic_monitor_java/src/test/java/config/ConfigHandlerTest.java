package config;

/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 * 
 *   http://www.apache.org/licenses/LICENSE-2.0
 * 
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */


import com.comcast.cdn.traffic_control.traffic_monitor.config.ConfigHandler;
import org.apache.commons.io.IOUtils;
import org.junit.Before;
import org.junit.BeforeClass;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.powermock.api.mockito.PowerMockito;
import org.powermock.core.classloader.annotations.PrepareForTest;
import org.powermock.modules.junit4.PowerMockRunner;
import org.powermock.reflect.Whitebox;

import java.io.File;
import java.io.FileReader;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.equalTo;
import static org.hamcrest.Matchers.nullValue;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;
import static org.powermock.api.mockito.PowerMockito.whenNew;

@PrepareForTest({ConfigHandler.class,File.class,IOUtils.class})
@RunWith(PowerMockRunner.class)
public class ConfigHandlerTest {
	static File mockConfigFile = mock(File.class);
	static ConfigHandler configHandler;

	@BeforeClass
	public static void beforeClass() throws Exception {
		whenNew(File.class).withArguments("/opt/traffic_monitor/conf/traffic_monitor_config.js").thenReturn(mockConfigFile);
		configHandler = ConfigHandler.getInstance();
	}

	@Before
	public void before() throws Exception {
		// Some food for thought about why "true" singletons are generally less desirable than dependency injection...
		// Without injecting a null monitor config object behind the scenes the tests don't work
		// we have to do this between every test too because once the singleton gets a hold
		// of a non-null config, it will never try to update it.
		Whitebox.setInternalState(configHandler, "config", (Object) null);
	}


	@Test
	public void itIndicatesNoConfigFileFound() {
		when(mockConfigFile.exists()).thenReturn(false);
		assertThat(configHandler.configFileExists(), equalTo(false));
		when(mockConfigFile.exists()).thenReturn(true);
		assertThat(configHandler.configFileExists(), equalTo(true));
	}

	@Test
	public void itReturnsMonitorConfigWithDefaults() {
		when(mockConfigFile.exists()).thenReturn(false);
		String configUrl = configHandler.getConfig().getCrConfigUrl();
		assertThat(configUrl, equalTo("https://${tmHostname}/CRConfig-Snapshots/${cdnName}/CRConfig.json"));
	}

	@Test
	public void itBuildsMonitorConfigFromJsonString() throws Exception {
		when(mockConfigFile.exists()).thenReturn(true);

		FileReader fileReader = mock(FileReader.class);
		whenNew(FileReader.class).withAnyArguments().thenReturn(fileReader);
		PowerMockito.mockStatic(IOUtils.class);
		when(IOUtils.toString(fileReader)).thenReturn("{\"traffic_monitor_config\":" +
			"{\"tm.crConfig.json.polling.url\": \"https://trafficops.kabletown.com/somepath/kabletown/config.json\"}" +
		"}");

		String configUrl = configHandler.getConfig().getCrConfigUrl();

		// Not trying to test that MonitorConfig does fancy stuff with string replacement
		assertThat(configUrl, equalTo("https://trafficops.kabletown.com/somepath/kabletown/config.json"));
	}

	@Test
	public void itReturnsCorrectPathToDbFile() {
		assertThat(configHandler.getDbFile("health-config.json").toString(), equalTo("/opt/traffic_monitor/db/health-config.json"));
	}
}
