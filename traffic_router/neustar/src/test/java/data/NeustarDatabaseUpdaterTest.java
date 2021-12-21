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

package data;

import org.apache.traffic_control.traffic_router.neustar.data.HttpClient;
import org.apache.traffic_control.traffic_router.neustar.data.NeustarDatabaseUpdater;
import org.apache.traffic_control.traffic_router.neustar.data.TarExtractor;
import org.apache.traffic_control.traffic_router.neustar.files.FilesMover;
import com.quova.bff.reader.io.GPDatabaseReader;
import org.apache.http.HttpEntity;
import org.apache.http.StatusLine;
import org.apache.http.client.methods.CloseableHttpResponse;
import org.apache.http.client.methods.HttpGet;
import org.apache.logging.log4j.core.appender.ConsoleAppender;
import org.apache.logging.log4j.Level;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.core.layout.PatternLayout;
import org.junit.Before;
import org.junit.BeforeClass;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.powermock.core.classloader.annotations.PowerMockIgnore;
import org.powermock.core.classloader.annotations.PrepareForTest;
import org.powermock.modules.junit4.PowerMockRunner;

import java.io.File;
import java.io.InputStream;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.Date;
import java.util.zip.GZIPInputStream;

import static org.hamcrest.CoreMatchers.equalTo;
import static org.hamcrest.CoreMatchers.nullValue;
import static org.hamcrest.MatcherAssert.assertThat;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.ArgumentMatchers.eq;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;
import static org.mockito.MockitoAnnotations.initMocks;
import static org.powermock.api.mockito.PowerMockito.mockStatic;
import static org.powermock.api.mockito.PowerMockito.whenNew;

@RunWith(PowerMockRunner.class)
@PowerMockIgnore("javax.net.ssl.*")
@PrepareForTest({NeustarDatabaseUpdater.class, Files.class})
public class NeustarDatabaseUpdaterTest {
	@Mock
	File neustarDatabaseDirectory;

	@Mock
	File neustarOldDatabaseDirectory;

	@Mock
	TarExtractor tarExtractor;

	@Mock
	FilesMover filesMover;

	@InjectMocks
	NeustarDatabaseUpdater neustarDatabaseUpdater;
	private File mockTmpDir;

	@BeforeClass
	public static void beforeClass() {
		LogManager.getRootLogger().addAppender(new ConsoleAppender(new PatternLayout("%d %-5p [%c]: %m%n")));
		LogManager.getRootLogger().setLevel(Level.INFO);
		LogManager.getLogger("org.springframework.context").setLevel(Level.WARN);
	}

	@Before
	public void before() throws Exception {
		initMocks(this);

		when(neustarOldDatabaseDirectory.isDirectory()).thenReturn(true);
		when(neustarOldDatabaseDirectory.lastModified()).thenReturn(1425236082000L);

		mockTmpDir = mock(File.class);
		when(mockTmpDir.getName()).thenReturn("123-abc-tmp");
		when(mockTmpDir.getParentFile()).thenReturn(neustarDatabaseDirectory);

		when(neustarDatabaseDirectory.listFiles()).thenReturn(new File[] {neustarOldDatabaseDirectory, mockTmpDir});

		Path path = mock(Path.class);
		when(path.toFile()).thenReturn(mockTmpDir);

		mockStatic(Files.class);
		when(Files.createTempDirectory(any(Path.class), eq("neustar-"))).thenReturn(path);
		when(tarExtractor.extractTo(eq(mockTmpDir), any(GZIPInputStream.class))).thenReturn(true);
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

		whenNew(GPDatabaseReader.Builder.class).withArguments(any(File.class)).thenReturn(mock(GPDatabaseReader.Builder.class));

		neustarDatabaseUpdater.setHttpClient(httpClient);
		neustarDatabaseUpdater.setNeustarDataUrl("http://example.com/neustardata.tgz");
		neustarDatabaseUpdater.setNeustarPollingTimeout(100);

		when(filesMover.updateCurrent(eq(neustarDatabaseDirectory), any(File.class), eq(neustarOldDatabaseDirectory))).thenReturn(true);
		assertThat(neustarDatabaseUpdater.update(), equalTo(true));
		verify(httpClient).close();
		verify(response).close();
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
