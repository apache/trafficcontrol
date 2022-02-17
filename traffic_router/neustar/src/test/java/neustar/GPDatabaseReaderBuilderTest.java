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

package neustar;

import com.maxmind.db.Reader;
import com.quova.bff.reader.io.GPDatabaseReader;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.powermock.core.classloader.annotations.PrepareForTest;
import org.powermock.modules.junit4.PowerMockRunner;

import java.io.File;
import java.io.IOException;

import static org.hamcrest.CoreMatchers.containsString;
import static org.hamcrest.CoreMatchers.equalTo;
import static org.hamcrest.CoreMatchers.notNullValue;
import static org.hamcrest.MatcherAssert.assertThat;
import static org.junit.Assert.fail;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;
import static org.powermock.api.mockito.PowerMockito.whenNew;

@RunWith(PowerMockRunner.class)
@PowerMockIgnore("javax.management.*")
public class GPDatabaseReaderBuilderTest {

	@Test
	public void buildThrowsExceptionForNullFile() throws Exception {
		try {
			new GPDatabaseReader.Builder(null).build();
			fail("Should have thrown exception!");
		} catch (IllegalArgumentException e) {
			assertThat(e.getMessage(), equalTo("The directory is null."));
		}
	}

	@Test
	public void buildThrowsExceptionForNondirectory() throws Exception {
		try {
			File file = mock(File.class);
			new GPDatabaseReader.Builder(file).build();
			fail("Should have thrown exception!");
		} catch (IllegalArgumentException e) {
			assertThat(e.getMessage(), containsString("is not a directory."));
		}
	}

	@Test
	public void buildThrowsExceptionForEmptyDirectory() throws Exception {
		try {
			File file = mock(File.class);
			when(file.isDirectory()).thenReturn(true);
			new GPDatabaseReader.Builder(file).build();
			fail("Should have thrown exception!");
		} catch (IOException e) {
			assertThat(e.getMessage(), equalTo("Error to load the gpdb files."));
		}
	}

	@Test
	@PrepareForTest({GPDatabaseReader.class, Reader.class})
	public void buildReturnsReaderForDirectoryWithGpdbFiles() throws Exception {
		File gpdbFile = mock(File.class);
		when(gpdbFile.getName()).thenReturn("dataV4.gpdb");
		when(gpdbFile.getPath()).thenReturn("/tmp/dataV4.gpdb");

		File file = mock(File.class);
		whenNew(File.class).withArguments("/tmp/dataV4.gpdb").thenReturn(file);

		Reader reader = mock(Reader.class);
		whenNew(Reader.class).withArguments(gpdbFile).thenReturn(reader);

		File directory = mock(File.class);
		when(directory.isDirectory()).thenReturn(true);
		when(directory.listFiles()).thenReturn(new File[] {gpdbFile});
		assertThat(new GPDatabaseReader.Builder(directory).build(), notNullValue());
	}
}
