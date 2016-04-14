package com.comcast.cdn.traffic_control.traffic_router.neustar.files;

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
