package com.comcast.cdn.traffic_control.traffic_router.keystore;

import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.security.KeyStore;

public class KeyStoreLoader {
	private final static org.apache.juli.logging.Log log = org.apache.juli.logging.LogFactory.getLog(KeyStoreLoader.class);
	private final Path keyStorePath;
	private final char[] keyPass;

	public KeyStoreLoader(final String keyStore, final char[] keyPass) {
		keyStorePath = Paths.get(keyStore);
		this.keyPass = keyPass;
	}

	public KeyStore load() {
		if (!Files.exists(keyStorePath)) {
			log.info("creating new keystore at " + keyStorePath.toAbsolutePath());
			return createKeyStore(keyStorePath, keyPass);
		}

		try (final InputStream inputStream = Files.newInputStream(keyStorePath)) {
			final KeyStore keyStore = KeyStore.getInstance(KeyStore.getDefaultType());
			keyStore.load(inputStream, keyPass);
			log.info("loaded keystore from " + keyStorePath.toAbsolutePath());
			return keyStore;
		} catch (Exception e) {
			log.error("Failed loading keystore from " + keyStorePath + " : " + e.getMessage());
		}

		return null;
	}

	public KeyStore createKeyStore(final Path path, final char[] keyPass) {
		Path existingPath = path;
		if (!Files.exists(existingPath)) {
			try {
				existingPath = Files.createFile(existingPath);
			} catch (IOException e) {
				log.error("Failed to create keystore at path " + existingPath.toAbsolutePath());
				return null;
			}
		}

		try (final OutputStream outputStream = Files.newOutputStream(existingPath)) {
			final KeyStore keyStore = KeyStore.getInstance(KeyStore.getDefaultType());
			keyStore.load(null, keyPass);
			keyStore.store(outputStream, keyPass);
			return keyStore;
		} catch (Exception e) {
			log.error("Failed initializing empty keystore at " + existingPath + " : " + e.getMessage());
		}

		return null;
	}
}
