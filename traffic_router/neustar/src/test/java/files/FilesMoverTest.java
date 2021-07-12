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

package files;

import org.apache.traffic_control.traffic_router.neustar.files.FilesMover;
import org.junit.Before;
import org.junit.Rule;
import org.junit.Test;
import org.junit.rules.TemporaryFolder;

import java.io.File;
import java.io.FileInputStream;
import java.io.FileOutputStream;
import java.util.Arrays;

import static org.hamcrest.CoreMatchers.equalTo;
import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.containsInAnyOrder;

public class FilesMoverTest {
	@Rule
	public TemporaryFolder parentFolder = new TemporaryFolder();

	File currentFile1;
	File currentFile2;

	File tmpFolder;
	File newFile1;
	File newFile2;

	File oldFolder;

	FilesMover filesMover = new FilesMover();

	@Before
	public void before() throws Exception {
		tmpFolder = parentFolder.newFolder("tmp");
		oldFolder = parentFolder.newFolder("old");

		currentFile1 = parentFolder.newFile("data1.txt");

		FileOutputStream fileOutputStream = new FileOutputStream(currentFile1);
		fileOutputStream.write("currentFile1".getBytes());
		fileOutputStream.close();

		currentFile2 = parentFolder.newFile("current2.txt");
		fileOutputStream = new FileOutputStream(currentFile2);
		fileOutputStream.write("currentFile2".getBytes());
		fileOutputStream.close();

		newFile1 = new File(tmpFolder, "data1.txt");
		fileOutputStream = new FileOutputStream(newFile1);
		fileOutputStream.write("new file 1".getBytes());
		fileOutputStream.close();

		newFile2 = new File(tmpFolder, "new2.txt");
		fileOutputStream = new FileOutputStream(newFile2);
		fileOutputStream.write("new file 2".getBytes());
		fileOutputStream.close();
	}

	@Test
	public void itPurgesDirectory() {
		assertThat(filesMover.purgeDirectory(tmpFolder), equalTo(true));
		assertThat(tmpFolder.list().length, equalTo(0));
	}

	@Test
	public void itMovesContents() {
		assertThat(filesMover.moveFiles(parentFolder.getRoot(), oldFolder), equalTo(true));
		assertThat(Arrays.asList(oldFolder.list()), containsInAnyOrder("data1.txt", "current2.txt"));
	}

	@Test
	public void itUpdatesCurrent() throws Exception {
		boolean updated = filesMover.updateCurrent(parentFolder.getRoot(), tmpFolder, oldFolder);
		assertThat(updated, equalTo(true));
		assertThat(Arrays.asList(parentFolder.getRoot().list()), containsInAnyOrder("data1.txt", "new2.txt", "tmp", "old"));
		File file = new File(parentFolder.getRoot(), "data1.txt");
		FileInputStream fileInputStream = new FileInputStream(file);
		byte[] buffer = new byte[64];
		int numBytes = fileInputStream.read(buffer);
		assertThat(new String(buffer, 0 , numBytes), equalTo("new file 1"));
		assertThat(oldFolder.list().length, equalTo(0));
	}
}
