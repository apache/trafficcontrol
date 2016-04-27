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

package com.comcast.cdn.traffic_control.traffic_router.neustar.data;

import org.apache.commons.compress.archivers.tar.TarArchiveEntry;
import org.apache.commons.compress.archivers.tar.TarArchiveInputStream;
import org.apache.log4j.Logger;

import java.io.File;
import java.io.FileOutputStream;
import java.io.IOException;
import java.io.InputStream;
import java.util.zip.GZIPInputStream;

public class TarExtractor {
	private final Logger LOGGER = Logger.getLogger(TarExtractor.class);

	public File extractTgzTo(File directory, InputStream inputStream) {
		try {
			return extractTo(directory, new GZIPInputStream(inputStream));
		} catch (IOException e) {
			LOGGER.error("Failed to extract gzip tar file to " + directory.getAbsolutePath() + ": " + e.getMessage());
			return null;
		}
	}

	public File extractTo(File directory, InputStream inputStream) {
		TarArchiveInputStream tarArchiveInputStream = new TarArchiveInputStream(inputStream);
		TarArchiveEntry tarArchiveEntry;
		try {
			while ((tarArchiveEntry = tarArchiveInputStream.getNextTarEntry()) != null) {
				if (tarArchiveEntry.isDirectory()) {
					continue;
				}

				File file = new File(directory, tarArchiveEntry.getName());
				if (!file.createNewFile()) {
					LOGGER.warn("Failed to extract file to " + file.getAbsolutePath());
					continue;
				}

				copyInputStreamToFile(tarArchiveInputStream, file);

			}
		} catch (IOException e) {
			LOGGER.error("Failed extracting tar archive to directory " + directory.getAbsolutePath() + " : " + e.getMessage());
		}

		return directory;
	}

	protected void copyInputStreamToFile(InputStream inputStream, File file) throws IOException {
		byte[] buffer = new byte[50 * 1024 * 1024];
		int bytesRead;

		FileOutputStream outputStream = new FileOutputStream(file);

		try {
			while ((bytesRead = inputStream.read(buffer)) != -1) {
				outputStream.write(buffer, 0, bytesRead);
			}
		} finally {
			outputStream.close();
		}
	}
}
