package com.comcast.cdn.traffic_control.traffic_router.keystore;

import com.comcast.cdn.traffic_control.traffic_router.properties.PropertiesGenerator;

import javax.naming.ldap.LdapName;
import javax.naming.ldap.Rdn;
import java.io.ByteArrayInputStream;
import java.io.IOException;
import java.io.OutputStream;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.security.Key;
import java.security.KeyFactory;
import java.security.KeyStore;
import java.security.KeyStoreException;
import java.security.Principal;
import java.security.PrivateKey;
import java.security.SecureRandom;
import java.security.cert.Certificate;
import java.security.cert.CertificateException;
import java.security.cert.CertificateFactory;
import java.security.cert.X509Certificate;
import java.security.spec.PKCS8EncodedKeySpec;
import java.util.ArrayList;
import java.util.Base64;
import java.util.Date;
import java.util.Enumeration;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Vector;

public class KeyStoreHelper {
	protected static org.apache.juli.logging.Log log = org.apache.juli.logging.LogFactory.getLog(KeyStoreHelper.class);
	public static final String KEYSTORE_PROPERTIES_PATH = "/conf/keystore.properties";
	public static final String KEYPASS_PROPERTY = "keypass";
	private volatile KeyStore keyStore;
	private char[] keyPass;
	private long lastLoaded;
	private final Map<String, PrivateKey> privateKeyMap = new HashMap<>();
	private final Map<String, Boolean> aliasCertMap = new HashMap<>();

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

	X509Certificate toCertificate(final String encodedCertificate) throws IOException, CertificateException {
		final byte[] encodedBytes = Base64.getDecoder().decode(encodedCertificate);

		try (ByteArrayInputStream stream = new ByteArrayInputStream(encodedBytes)) {
			return (X509Certificate) CertificateFactory.getInstance("X.509").generateCertificate(stream);
		}
	}

	public boolean importCertificateChain(final String alias, final String encodedKey, final String[] encodedCertificateChain) {
		try {
			final X509Certificate x509Chain[] = new X509Certificate[encodedCertificateChain.length];

			for (int i = 0; i < encodedCertificateChain.length; i++) {
				x509Chain[i] = toCertificate(encodedCertificateChain[i]);
				final Date notAfter = x509Chain[i].getNotAfter();
				final Date notBefore = x509Chain[i].getNotBefore();
				final Principal subject = x509Chain[i].getSubjectDN();
				final Principal issuer = x509Chain[i].getIssuerDN();
				log.info("Import [" + alias + "][" + i + "] not before " + notBefore + " and not after " + notAfter);
				log.info("Import [" + alias + "][" + i + "] subject " + subject);
				log.info("Import [" + alias + "][" + i + "] issuer " + issuer);
			}

			byte[] keyBytes = Base64.getDecoder().decode(encodedKey.getBytes());
			PKCS8EncodedKeySpec keySpec = new PKCS8EncodedKeySpec(keyBytes);
			KeyFactory fact = KeyFactory.getInstance("RSA");
			PrivateKey key = fact.generatePrivate(keySpec);

			return importCertificateChain(alias, key, x509Chain);
		} catch (Exception e) {
			log.error("Failed importing certificates for alias '" + alias + "'");
			log.error(e);
		}

		return false;
	}

	public boolean importCertificateChain(final String alias, final PrivateKey privateKey, final Certificate[] certificateChain) {
		try {
			keyStore.setKeyEntry(alias, privateKey, keyPass, certificateChain);
			privateKeyMap.put(alias, privateKey);
			log.info("Imported certificate chain into keystore for " + alias);
		} catch (Exception e) {
			log.error("Failed importing certificate chain with alias '" + alias + "' to keystore : " + e.getMessage());
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
			log.error("Failed retrieving common names data from the keystore: " + e.getClass().getSimpleName() + " " + e.getMessage());
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

		if (keyStore == null) {
			log.error("Failed reloading keystore from " + getKeystorePath());
			return null;
		}

		lastLoaded = System.currentTimeMillis();

		try {
			final Enumeration<String> aliases = keyStore.aliases();
			while (aliases.hasMoreElements()) {
				final String alias = aliases.nextElement();
				final Key key = keyStore.getKey(alias, getKeyPass());
				if (key instanceof PrivateKey) {
					privateKeyMap.put(alias, (PrivateKey) key);
				}
			}
			log.info("Reloaded keystore path from " + getKeystorePath() + " " + keyStore.size() + " entries");
		} catch (Exception e) {
			log.error("Cannot get size of keystore ",e);
		}
		return keyStore;
	}

	public boolean clearCertificates() {
		aliasCertMap.clear();
		try {
			Enumeration<String> aliases = keyStore.aliases();
			while (aliases.hasMoreElements()) {
				final String alias = aliases.nextElement();
				keyStore.deleteEntry(alias);
				privateKeyMap.remove(alias);
			}
		} catch (KeyStoreException e) {
			log.error("Failed to clear certificates from keystore!");
		}

		return false;
	}

	public PrivateKey getPrivateKey(String alias) {
		if (!privateKeyMap.containsKey(alias)) {
			log.warn("No private key exists for " + alias);
		}
		return privateKeyMap.get(alias);
	}

	@SuppressWarnings("PMD.UseArrayListInsteadOfVector")
	public Enumeration<String> getAliases() {
		try {
			return keyStore.aliases();
		} catch (Exception e) {
			log.warn("Failed to get aliases from keystore!: " + e.getMessage());
		}
		return new Vector<String>().elements();
	}

	public boolean hasCertificate(String prefix) {

		if (aliasCertMap.containsKey(prefix)) {
			return aliasCertMap.get(prefix);
		}

		aliasCertMap.put(prefix, false);

		try {
			Enumeration<String> aliasIterator = keyStore.aliases();
			while (aliasIterator.hasMoreElements()) {
				if (aliasIterator.nextElement().startsWith(prefix)) {
					aliasCertMap.put(prefix, true);
					break;
				}
			}
		} catch (Exception e) {
			log.error("Failed to search keystore aliases for prefix " + prefix + " : " + e.getMessage());
		}

		return aliasCertMap.get(prefix);
	}
}
