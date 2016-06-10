package keystore;

import com.comcast.cdn.traffic_control.traffic_router.keystore.KeyStoreHelper;
import com.comcast.cdn.traffic_control.traffic_router.keystore.KeyStoreLoader;
import com.comcast.cdn.traffic_control.traffic_router.properties.PropertiesGenerator;
import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.powermock.api.mockito.PowerMockito;
import org.powermock.core.classloader.annotations.PrepareForTest;
import org.powermock.modules.junit4.PowerMockRunner;
import sun.security.tools.keytool.CertAndKeyGen;
import sun.security.x509.X500Name;

import java.io.OutputStream;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.security.KeyStore;
import java.security.PrivateKey;
import java.security.SecureRandom;
import java.security.cert.Certificate;
import java.security.cert.X509Certificate;
import java.util.Base64;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.equalTo;
import static org.mockito.Matchers.any;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;
import static org.powermock.api.mockito.PowerMockito.mockStatic;
import static org.powermock.api.mockito.PowerMockito.whenNew;

@RunWith(PowerMockRunner.class)
@PrepareForTest({KeyStoreHelper.class, KeyStoreLoader.class, KeyStore.class})
public class KeyStoreHelperTest {
	private KeyStore keyStore;
	private SecureRandom secureRandom;
	private Base64.Encoder encoder;
	private OutputStream outputStream;

	@Before
	public void before() throws Exception {
		secureRandom = mock(SecureRandom.class);
		whenNew(SecureRandom.class).withNoArguments().thenReturn(secureRandom);

		encoder = mock(Base64.Encoder.class);
		when(encoder.withoutPadding()).thenReturn(encoder);

		when(encoder.encodeToString(any(byte[].class))).thenReturn("password");

		mockStatic(Base64.class);
		when(Base64.getEncoder()).thenReturn(encoder);

		PropertiesGenerator propertiesGenerator = mock(PropertiesGenerator.class);
		when(propertiesGenerator.getProperty(KeyStoreHelper.KEYPASS_PROPERTY, "password")).thenReturn("password");

		whenNew(PropertiesGenerator.class).withArguments("/opt/traffic_router/conf/keystore.properties").thenReturn(propertiesGenerator);

		KeyStoreLoader keyStoreLoader = mock(KeyStoreLoader.class);
		whenNew(KeyStoreLoader.class).withArguments("/opt/traffic_router/.keystore", "password".toCharArray()).thenReturn(keyStoreLoader);

		keyStore = PowerMockito.mock(KeyStore.class);
		when(keyStoreLoader.load()).thenReturn(keyStore);

		mockStatic(Paths.class);

		Path keystorePath = mock(Path.class);
		when(Paths.get("/opt/traffic_router/.keystore")).thenReturn(keystorePath);

		outputStream = mock(OutputStream.class);

		mockStatic(Files.class);
		when(Files.newOutputStream(keystorePath)).thenReturn(outputStream);
	}

	@Test
	public void itInitializesSingletonAndImportsCerts() throws Exception {
		KeyStoreHelper keyStoreHelper =  KeyStoreHelper.getInstance();

		verify(encoder).withoutPadding();
		verify(encoder).encodeToString(any(byte[].class));
		verify(secureRandom).nextBytes(any(byte[].class));

		assertThat(keyStoreHelper.getKeyPass(), equalTo("password".toCharArray()));
		assertThat(keyStoreHelper.getKeyStore(), equalTo(keyStore));

		CertAndKeyGen certAndKeyGen = new CertAndKeyGen("RSA", "SHA1WithRSA", null);
		certAndKeyGen.generate(2048);
		PrivateKey privateKey = certAndKeyGen.getPrivateKey();
		X509Certificate x509Certificate = certAndKeyGen.getSelfCertificate(new X500Name("CN=*.deliveryservice.cdn.example.com"), 600L);

		boolean success = keyStoreHelper.importCertificate("deliveryservice.cdn.example.com",privateKey,x509Certificate);
		assertThat(success, equalTo(true));

		verify(keyStore).setKeyEntry("deliveryservice.cdn.example.com", privateKey, "password".toCharArray(), new Certificate[] {x509Certificate});
		verify(keyStore).store(outputStream, "password".toCharArray());
	}
}
