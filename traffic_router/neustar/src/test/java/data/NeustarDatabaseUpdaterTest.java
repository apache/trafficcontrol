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

package data;

import com.comcast.cdn.traffic_control.traffic_router.neustar.data.TarExtractor;
import com.comcast.cdn.traffic_control.traffic_router.neustar.files.FilesMover;
import com.quova.bff.reader.io.GPDatabaseReader;
import org.junit.Before;
import org.junit.runner.RunWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import com.comcast.cdn.traffic_control.traffic_router.neustar.data.HttpClient;
import com.comcast.cdn.traffic_control.traffic_router.neustar.data.NeustarDatabaseUpdater;
import org.apache.http.HttpEntity;
import org.apache.http.StatusLine;
import org.apache.http.client.methods.CloseableHttpResponse;
import org.apache.http.client.methods.HttpGet;
import org.apache.log4j.ConsoleAppender;
import org.apache.log4j.Level;
import org.apache.log4j.LogManager;
import org.apache.log4j.PatternLayout;
import org.junit.BeforeClass;
import org.junit.Test;
import org.powermock.core.classloader.annotations.PowerMockIgnore;
import org.powermock.core.classloader.annotations.PrepareForTest;
import org.powermock.modules.junit4.PowerMockRunner;

import java.io.File;
import java.io.InputStream;
import java.util.Date;
import java.util.zip.GZIPInputStream;

import static org.hamcrest.CoreMatchers.equalTo;
import static org.hamcrest.CoreMatchers.nullValue;
import static org.hamcrest.MatcherAssert.assertThat;
import static org.mockito.Matchers.any;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;
import static org.mockito.MockitoAnnotations.initMocks;
import static org.powermock.api.mockito.PowerMockito.whenNew;

@RunWith(PowerMockRunner.class)
@PowerMockIgnore("javax.net.ssl.*")
public class NeustarDatabaseUpdaterTest {
	@Mock
	File neustarDatabaseDirectory;

	@Mock
	File neustarTempDatabaseDirectory;

	@Mock
	File neustarOldDatabaseDirectory;

	@Mock
	TarExtractor tarExtractor;

	@Mock
	FilesMover filesMover;

	@InjectMocks
	NeustarDatabaseUpdater neustarDatabaseUpdater;

	@BeforeClass
	public static void beforeClass() {
		LogManager.getRootLogger().addAppender(new ConsoleAppender(new PatternLayout("%d %-5p [%c]: %m%n")));
		LogManager.getRootLogger().setLevel(Level.INFO);
		LogManager.getLogger("org.springframework.context").setLevel(Level.WARN);
	}

	@Before
	public void before() {
		initMocks(this);

		when(neustarTempDatabaseDirectory.isDirectory()).thenReturn(true);
		when(neustarTempDatabaseDirectory.lastModified()).thenReturn(1425236082000L);

		when(neustarOldDatabaseDirectory.isDirectory()).thenReturn(true);
		when(neustarOldDatabaseDirectory.lastModified()).thenReturn(1425236082000L);

		when(neustarDatabaseDirectory.listFiles()).thenReturn(new File[] {neustarOldDatabaseDirectory, neustarTempDatabaseDirectory});
		when(neustarTempDatabaseDirectory.mkdirs()).thenReturn(true);
	}

	@Test
	@PrepareForTest({NeustarDatabaseUpdater.class, GZIPInputStream.class, GPDatabaseReader.Builder.class})
	public void itRetrievesRemoteFileContents() throws Exception {

		StatusLine statusLine = mock(StatusLine.class);
		when(statusLine.getStatusCode()).thenReturn(200);

		InputStream remoteInputStream = mock(InputStream.class);

		HttpEntity httpEntity = mock(HttpEntity.class);
		when(httpEntity.getContent()).thenReturn(remoteInputStream);

		CloseableHttpResponse response = mock(CloseableHttpResponse.class);
		when(response.getStatusLine()).thenReturn(statusLine);
		when(response.getEntity()).thenReturn(httpEntity);

		HttpClient httpClient = mock(HttpClient.class);
		when(httpClient.execute(any(HttpGet.class))).thenReturn(response);

		GZIPInputStream gzipInputStream = mock(GZIPInputStream.class);
		whenNew(GZIPInputStream.class).withArguments(remoteInputStream).thenReturn(gzipInputStream);

		whenNew(GPDatabaseReader.Builder.class).withArguments(neustarTempDatabaseDirectory).thenReturn(mock(GPDatabaseReader.Builder.class));

		neustarDatabaseUpdater.setHttpClient(httpClient);
		neustarDatabaseUpdater.setNeustarDataUrl("http://example.com/neustardata.tgz");
		neustarDatabaseUpdater.setNeustarPollingTimeout(100);

		when(filesMover.updateCurrent(neustarDatabaseDirectory, neustarTempDatabaseDirectory, neustarOldDatabaseDirectory)).thenReturn(true);
		assertThat(neustarDatabaseUpdater.update(), equalTo(true));
		verify(httpClient).close();
		verify(response).close();
	}

	@Test
	public void itExtractsRemoteContentToNewDirectory() throws Exception {
		InputStream remoteInputStream = mock(InputStream.class);

		when(tarExtractor.extractTgzTo(neustarTempDatabaseDirectory, remoteInputStream)).thenReturn(neustarTempDatabaseDirectory);

		assertThat(neustarDatabaseUpdater.extractRemoteContent(remoteInputStream), equalTo(neustarTempDatabaseDirectory));
	}

	@Test
	public void itDeterminesLatestBuildDate() {
		assertThat(neustarDatabaseUpdater.getDatabaseBuildDate(), nullValue());

		File file = mock(File.class);
		when(file.lastModified()).thenReturn(1425236082000L);
		when(neustarDatabaseDirectory.listFiles()).thenReturn(new File[] {file});

		assertThat(neustarDatabaseUpdater.getDatabaseBuildDate(), equalTo(new Date(1425236082000L)));
	}
}
