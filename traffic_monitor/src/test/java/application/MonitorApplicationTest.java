package application;

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


import com.comcast.cdn.traffic_control.traffic_monitor.MonitorApplication;
import com.comcast.cdn.traffic_control.traffic_monitor.config.ConfigHandler;
import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.powermock.core.classloader.annotations.PrepareForTest;
import org.powermock.modules.junit4.PowerMockRunner;
import org.powermock.reflect.Whitebox;

import java.security.AccessControlException;
import java.security.Permission;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.equalTo;
import static org.junit.Assert.fail;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;
import static org.powermock.api.mockito.PowerMockito.mockStatic;

@PrepareForTest({MonitorApplication.class, ConfigHandler.class})
@RunWith(PowerMockRunner.class)
public class MonitorApplicationTest {
	private final SecurityManager originalSecurityManager = System.getSecurityManager();

	@Before
	public void before() {
		System.setSecurityManager(new SecurityManager() {
			@Override
			public void checkPermission(Permission perm)
			{
				// allow anything.
			}

			@Override
			public void checkPermission(Permission perm, Object context)
			{
				// allow anything.
			}

			@Override
			public void checkExit(int status)
			{
				super.checkExit(status);
				throw new AccessControlException("Boom");
			}
		});
	}

	@Test
	public void itSystemExitsWhenConfigFileIsMissing() {
		ConfigHandler configHandler = mock(ConfigHandler.class);
		when(configHandler.configFileExists()).thenReturn(false);

		mockStatic(ConfigHandler.class);
		when(ConfigHandler.getInstance()).thenReturn(configHandler);

		MonitorApplication monitorApplication = new MonitorApplication();
		Whitebox.setInternalState(monitorApplication, "settingsAccessible", true);

		try {
			monitorApplication.init();
			fail("Init did not do SystemExit");
		} catch (AccessControlException e) {
			assertThat(e.getMessage(), equalTo("Boom"));
		} finally {
			System.setSecurityManager(originalSecurityManager);
		}
	}
}
