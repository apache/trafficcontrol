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

package org.apache.traffic_control.traffic_router.neustar.files;

import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

import java.io.File;
import java.io.IOException;
import java.nio.file.FileVisitResult;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.nio.file.SimpleFileVisitor;
import java.nio.file.StandardCopyOption;
import java.nio.file.attribute.BasicFileAttributes;

public class FilesMover {
	private static final Logger LOGGER = LogManager.getLogger(FilesMover.class);

	public boolean purgeDirectory(File directory) {
		try {
			Files.walkFileTree(Paths.get(directory.getAbsolutePath()), new SimpleFileVisitor<Path>() {
				@Override
				public FileVisitResult visitFile(Path file, BasicFileAttributes basicFileAttributes) throws IOException {
					Files.delete(file);
					return FileVisitResult.CONTINUE;
				}
			});

			return true;
		} catch (IOException e) {
			LOGGER.error("Failed purging directory " + directory.getAbsolutePath() + ": " + e.getMessage());
			return false;
		}
	}

	public boolean moveFiles(File sourceDirectory, File destinationDirectory) {
		if (!destinationDirectory.exists() && !destinationDirectory.mkdirs()) {
			return false;
		}

		if (!destinationDirectory.canWrite()) {
			return false;
		}

		for (File file : sourceDirectory.listFiles()) {
			if (file.isDirectory()) {
				continue;
			}
			Path source = Paths.get(file.getAbsolutePath());
			Path destination = Paths.get(destinationDirectory.getAbsolutePath(), file.getName());
			try {
				Files.move(source, destination, StandardCopyOption.REPLACE_EXISTING);
			} catch (IOException e) {
				return false;
			}
		}

		return true;
	}

	public boolean updateCurrent(File currentDirectory, File newDirectory, File oldDirectory) {
		if (!currentDirectory.canWrite() || !newDirectory.canWrite()) {
			return false;
		}

		if (oldDirectory.exists() && !purgeDirectory(oldDirectory)) {
			return false;
		}

		if (!moveFiles(currentDirectory, oldDirectory)) {
			return false;
		}

		if (!moveFiles(newDirectory, currentDirectory)) {
			moveFiles(oldDirectory, currentDirectory);
			return false;
		}

		purgeDirectory(oldDirectory);
		return true;
	}
}
