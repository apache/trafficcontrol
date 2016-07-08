package com.comcast.cdn.traffic_control.traffic_router.keystore;

import com.comcast.cdn.traffic_control.traffic_router.properties.PropertiesGenerator;

import javax.naming.ldap.LdapName;
import javax.naming.ldap.Rdn;
import java.io.ByteArrayInputStream;
import java.io.IOException;
import java.io.OutputStream;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.security.KeyFactory;
import java.security.KeyStore;
import java.security.KeyStoreException;
import java.security.PrivateKey;
import java.security.SecureRandom;
import java.security.cert.Certificate;
import java.security.cert.CertificateFactory;
import java.security.cert.X509Certificate;
import java.security.spec.PKCS8EncodedKeySpec;
import java.util.ArrayList;
import java.util.Base64;
import java.util.Enumeration;
import java.util.List;

public class KeyStoreHelper {
	protected static org.apache.juli.logging.Log log = org.apache.juli.logging.LogFactory.getLog(KeyStoreHelper.class);
	public static final String KEYSTORE_PROPERTIES_PATH = "/conf/keystore.properties";
	public static final String KEYPASS_PROPERTY = "keypass";
	private KeyStore keyStore;
	private char[] keyPass;
	private long lastLoaded;

	// Recommended Singleton Pattern implementation
	// https://community.oracle.com/docs/DOC-918906

	private KeyStoreHelper() {
		getKeyPass();
		getKeyStore();
	}

	public static KeyStoreHelper getInstance() {
		return KeyStoreHelperHolder.HELPER;
	}

	@SuppressWarnings("PMD.AccessorClassGeneration")
	private static class KeyStoreHelperHolder {
		private static final KeyStoreHelper HELPER = new KeyStoreHelper();
	}

	public String createKeypass() {
		final byte[] bytes = new byte[20];
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
			return reload();
		}
		return keyStore;
	}

	public boolean importCertificate(final String alias, final String encodedKey, final String encodedCertificate) {
		try {
			X509Certificate x509Certificate = (X509Certificate) CertificateFactory.getInstance("X.509")
				.generateCertificate(new ByteArrayInputStream(Base64.getDecoder().decode(encodedCertificate)));

			byte[] keyBytes = Base64.getDecoder().decode(encodedKey.getBytes());
			PKCS8EncodedKeySpec keySpec = new PKCS8EncodedKeySpec(keyBytes);
			KeyFactory fact = KeyFactory.getInstance("RSA");
			PrivateKey key = fact.generatePrivate(keySpec);

			return importCertificate(alias, key, x509Certificate);
		} catch (Exception e) {
			e.printStackTrace();
		}

		return false;
	}

	public boolean importCertificate(final String alias, final PrivateKey privateKey, final Certificate certificate) {
		try (final OutputStream outputStream = Files.newOutputStream(Paths.get(getKeystorePath()))) {
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
		final String keyStorePath = System.getProperty("deploy.dir", "/opt/traffic_router") + "/db/.keystore";
		System.setProperty("javax.net.ssl.keyStore", keyStorePath);
		return keyStorePath;
	}

	public List<String> getAllCommonNames() {
		final List<String> commonNames = new ArrayList<>();

		try {
			final Enumeration<String> aliases = keyStore.aliases();

			while (aliases.hasMoreElements()) {
				final String alias = aliases.nextElement();
				final Certificate certificate = keyStore.getCertificate(alias);

				if (!(certificate instanceof X509Certificate)) {
					continue;
				}

				final X509Certificate x509cert = (X509Certificate) certificate;
				final LdapName ldapDN = new LdapName(x509cert.getSubjectX500Principal().getName());

				for (final Rdn rdn: ldapDN.getRdns()) {
					if ("CN".equals(rdn.getType())) {
						commonNames.add(rdn.getValue().toString());
					}
				}
			}

		} catch (Exception e) {
			log.error("Failed retrieving name stuff from the keystore: " + e.getClass().getSimpleName() + " " + e.getMessage());
		}

		return commonNames;
	}

	public boolean save() {
		try (final OutputStream outputStream = Files.newOutputStream(Paths.get(getKeystorePath()))) {
			keyStore.store(outputStream, keyPass);
		} catch (Exception e) {
			log.error("Failed saving new data to keystore at " + getKeystorePath() + " : " + e.getMessage());
			return false;
		}

		return true;
	}

	public long getLastModified() {
		try {
			return Files.getLastModifiedTime(Paths.get(getKeystorePath())).toMillis();
		} catch (IOException e) {
			e.printStackTrace();
			return(0L);
		}
	}

	public long getLastLoaded() {
		return lastLoaded;
	}

	public KeyStore reload() {
		keyStore = new KeyStoreLoader(getKeystorePath(), getKeyPass()).load();
		lastLoaded = System.currentTimeMillis();
		return keyStore;
	}

	public boolean clearCertificates() {
		try {
			Enumeration<String> aliases = keyStore.aliases();
			while (aliases.hasMoreElements()) {
				keyStore.deleteEntry(aliases.nextElement());
			}
		} catch (KeyStoreException e) {
			log.error("Failed to clear certificates from keystore!");
		}

		return false;
	}
}
