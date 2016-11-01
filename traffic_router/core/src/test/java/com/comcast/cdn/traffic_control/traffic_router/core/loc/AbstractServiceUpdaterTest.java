package com.comcast.cdn.traffic_control.traffic_router.core.loc;

import org.apache.wicket.ajax.json.JSONException;
import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.powermock.api.mockito.PowerMockito;
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
		updater.updateDatabase();

		verify(connection, times(0)).setRequestProperty(eq("If-None-Match"), anyString());
		verify(connection).getHeaderField("ETag");

		updater.updateDatabase();
		verify(connection).setRequestProperty(eq("If-None-Match"), anyString());
		verify(connection, times(2)).getHeaderField("ETag");
	}

	class Updater extends AbstractServiceUpdater {
		@Override
		public boolean verifyDatabase(File dbFile) throws IOException, JSONException {
			return false;
		}

		@Override
		public boolean loadDatabase() throws IOException, JSONException {
			return false;
		}

		@Override
		public boolean isLoaded() {
			return true;
		}


	}
}
