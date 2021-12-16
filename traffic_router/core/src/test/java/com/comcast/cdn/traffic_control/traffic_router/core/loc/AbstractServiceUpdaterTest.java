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

package com.comcast.cdn.traffic_control.traffic_router.core.loc;

import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.powermock.api.mockito.PowerMockito;
import org.powermock.core.classloader.annotations.PowerMockIgnore;
import org.powermock.core.classloader.annotations.PrepareForTest;
import org.powermock.modules.junit4.PowerMockRunner;

import java.io.File;
import java.io.IOException;
import java.net.HttpURLConnection;
import java.net.URL;
import java.nio.file.Files;
import java.nio.file.Path;

import static org.mockito.Matchers.any;
import static org.mockito.Matchers.anyString;
import static org.mockito.Matchers.eq;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.times;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;
import static org.powermock.api.mockito.PowerMockito.mockStatic;
import static org.powermock.api.mockito.PowerMockito.whenNew;

@RunWith(PowerMockRunner.class)
@PrepareForTest({AbstractServiceUpdater.class, HttpURLConnection.class, URL.class, Files.class})
@PowerMockIgnore("javax.management.*")
public class AbstractServiceUpdaterTest {

	private HttpURLConnection connection;
	private Path databasesDirectory;
	private Path databasePath;
	private File databaseFile;

	@Before
	public void before() throws Exception {
		databaseFile = mock(File.class);
		when(databaseFile.exists()).thenReturn(true);
		when(databaseFile.lastModified()).thenReturn(1L);

		databasePath = mock(Path.class);
		when(databasePath.toFile()).thenReturn(databaseFile);

		databasesDirectory = mock(Path.class);
		when(databasesDirectory.resolve(anyString())).thenReturn(databasePath);

		mockStatic(Files.class);
		PowerMockito.when(Files.exists(any(Path.class))).thenReturn(true);


		connection = PowerMockito.mock(HttpURLConnection.class);
		when(connection.getHeaderField("ETag")).thenReturn("version-1");
		when(connection.getResponseCode()).thenReturn(304);

		URL url = PowerMockito.mock(URL.class);
		when(url.openConnection()).thenReturn(connection);

		whenNew(URL.class).withAnyArguments().thenReturn(url);
	}

	@Test
	public void itUsesETag() throws Exception {
		Updater updater = new Updater();
		updater.setDatabasesDirectory(databasesDirectory);
		updater.dataBaseURL = "http://www.example.com";
		updater.updateDatabase();

		verify(connection, times(0)).setRequestProperty(eq("If-None-Match"), anyString());
		verify(connection).getHeaderField("ETag");

		updater.updateDatabase();
		verify(connection).setRequestProperty(eq("If-None-Match"), anyString());
		verify(connection, times(2)).getHeaderField("ETag");
	}

	class Updater extends AbstractServiceUpdater {
		@Override
		public boolean verifyDatabase(File dbFile) throws IOException {
			return false;
		}

		@Override
		public boolean loadDatabase() throws IOException {
			return false;
		}

		@Override
		public boolean isLoaded() {
			return true;
		}


	}
}
