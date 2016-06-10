package com.comcast.cdn.traffic_control.traffic_router.keystore;

import com.comcast.cdn.traffic_control.traffic_router.properties.PropertiesGenerator;

import java.io.OutputStream;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.security.KeyStore;
import java.security.PrivateKey;
import java.security.SecureRandom;
import java.security.cert.Certificate;
import java.util.Base64;

public class KeyStoreHelper {
	protected static org.apache.juli.logging.Log log = org.apache.juli.logging.LogFactory.getLog(KeyStoreHelper.class);
	public static final String KEYSTORE_PROPERTIES_PATH = "/conf/keystore.properties";
	public static final String KEYPASS_PROPERTY = "keypass";
	private KeyStore keyStore;
	private char[] keyPass;

	// Recommended Singleton Pattern implementation
	// https://community.oracle.com/docs/DOC-918906

	private KeyStoreHelper() {
		getKeyPass();
		getKeyStore();
	}

	public static KeyStoreHelper getInstance() {
		return KeyStoreHelperHolder.HELPER;
	}

	private static class KeyStoreHelperHolder {
		private static final KeyStoreHelper HELPER = new KeyStoreHelper();
	}

	public String createKeypass() {
		byte[] bytes = new byte[20];
		new SecureRandom().nextBytes(bytes);
		return Base64.getEncoder().withoutPadding().encodeToString(bytes);
	}

	public char[] getKeyPass() {
		if (keyPass == null) {
			keyPass = new PropertiesGenerator(getKeystorePropertiesPath()).getProperty(KEYPASS_PROPERTY, createKeypass()).toCharArray();
		}

		return keyPass;
	}

	public KeyStore getKeyStore() {
		if (keyStore == null) {
			String keystorePath = getKeystorePath();
			char[] keyPass = getKeyPass();
			keyStore = new KeyStoreLoader(keystorePath, keyPass).load();
		}
		return keyStore;
	}

	public boolean importCertificate(String alias, PrivateKey privateKey, Certificate certificate) {
		try (OutputStream outputStream = Files.newOutputStream(Paths.get(getKeystorePath()))) {
			keyStore.setKeyEntry(alias, privateKey, keyPass, new Certificate[] {certificate});
			keyStore.store(outputStream, keyPass);
		} catch (Exception e) {
			log.error("Failed importing certificate with alias '" + alias + "' to keystore at " + getKeystorePath() + " : " + e.getMessage());
			return false;
		}

		return true;
	}

	public String getKeystorePropertiesPath() {
		return System.getProperty("deploy.dir", "/opt/traffic_router") + KEYSTORE_PROPERTIES_PATH;
	}

	public String getKeystorePath() {
		final String keyStorePath = System.getProperty("deploy.dir", "/opt/traffic_router") + "/.keystore";
		System.setProperty("javax.net.ssl.trustStore", keyStorePath);
		return keyStorePath;
	}
}
