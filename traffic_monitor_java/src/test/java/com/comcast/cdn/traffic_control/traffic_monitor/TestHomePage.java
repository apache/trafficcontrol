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

import com.comcast.cdn.traffic_control.traffic_monitor.config.ConfigHandler;
import com.comcast.cdn.traffic_control.traffic_monitor.config.MonitorConfig;
import com.comcast.cdn.traffic_control.traffic_monitor.health.CacheWatcher;
import org.apache.wicket.util.tester.WicketTester;
import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.powermock.core.classloader.annotations.PrepareForTest;
import org.powermock.modules.junit4.PowerMockRunner;

import javax.net.ssl.SSLContext;

import static org.powermock.api.mockito.PowerMockito.mock;
import static org.powermock.api.mockito.PowerMockito.mockStatic;
import static org.powermock.api.mockito.PowerMockito.when;

import java.io.File;

import static org.mockito.Matchers.anyString;

@PrepareForTest({ConfigHandler.class, CacheWatcher.class, SSLContext.class})
@RunWith(PowerMockRunner.class)
public class TestHomePage {
	private WicketTester tester;

	@Before
	public void setUp() throws Exception {
		mockStatic(SSLContext.class);

		when(SSLContext.getInstance("TLS")).thenReturn(mock(SSLContext.class));

		MonitorConfig monitorConfig = mock(MonitorConfig.class);
		when(monitorConfig.getHealthPollingInterval()).thenReturn(10 * 1000);
		when(monitorConfig.getTmFrequency()).thenReturn(10 * 1000L);
		when(monitorConfig.getPeerPollingInterval()).thenReturn(10 * 1000L);
		when(monitorConfig.getHeathUrl()).thenReturn("http://example.com/healthParams");
		when(monitorConfig.getCrConfigUrl()).thenReturn("http://example.com/crConfig");
		when(monitorConfig.getPeerThreadPool()).thenReturn(1);
		when(monitorConfig.getPeerUrl()).thenReturn("http://example.com/publish/CrStates");

		ConfigHandler configHandler = mock(ConfigHandler.class);

		when(configHandler.getConfig()).thenReturn(monitorConfig);
		when(configHandler.configFileExists()).thenReturn(true);
		when(configHandler.getDbFile(anyString())).thenReturn(mock(File.class));

		mockStatic(ConfigHandler.class);
		when(ConfigHandler.getInstance()).thenReturn(configHandler);

		tester = new WicketTester(new MonitorApplication());
	}

	@Test
	public void homepageRendersSuccessfully() {
		//start and render the test page
		tester.startPage(Index.class);

		//assert rendered page class
		tester.assertRenderedPage(Index.class);
	}
}
