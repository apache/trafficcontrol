package keystore;

import com.comcast.cdn.traffic_control.traffic_router.keystore.KeyStoreHelper;
import com.comcast.cdn.traffic_control.traffic_router.keystore.KeyStoreLoader;
import com.comcast.cdn.traffic_control.traffic_router.properties.PropertiesGenerator;
import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.powermock.api.mockito.PowerMockito;
import org.powermock.core.classloader.annotations.PowerMockIgnore;
import org.powermock.core.classloader.annotations.PrepareForTest;
import org.powermock.modules.junit4.PowerMockRunner;
import sun.security.tools.keytool.CertAndKeyGen;
import sun.security.x509.X500Name;

import javax.naming.ldap.LdapName;
import javax.naming.ldap.Rdn;
import javax.security.auth.x500.X500Principal;
import java.io.OutputStream;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.security.KeyStore;
import java.security.PrivateKey;
import java.security.SecureRandom;
import java.security.cert.Certificate;
import java.security.cert.X509Certificate;
import java.util.ArrayList;
import java.util.Base64;
import java.util.List;
import java.util.Vector;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.containsInAnyOrder;
import static org.hamcrest.Matchers.equalTo;
import static org.mockito.Matchers.any;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;
import static org.powermock.api.mockito.PowerMockito.mockStatic;
import static org.powermock.api.mockito.PowerMockito.whenNew;

@RunWith(PowerMockRunner.class)
@PrepareForTest({KeyStoreHelper.class, KeyStoreLoader.class, KeyStore.class})
@PowerMockIgnore("javax.security.auth.x500.X500Principal")
public class KeyStoreHelperTest {
	private KeyStore keyStore;
	private SecureRandom secureRandom;
	private Base64.Encoder encoder;
	private OutputStream outputStream;
	private X509Certificate x509Certificate;
	private PrivateKey privateKey;

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
		whenNew(KeyStoreLoader.class).withArguments("/opt/traffic_router/db/.keystore", "password".toCharArray()).thenReturn(keyStoreLoader);

		keyStore = PowerMockito.mock(KeyStore.class);
		when(keyStoreLoader.load()).thenReturn(keyStore);

		mockStatic(Paths.class);

		Path keystorePath = mock(Path.class);
		when(Paths.get("/opt/traffic_router/db/.keystore")).thenReturn(keystorePath);

		outputStream = mock(OutputStream.class);

		mockStatic(Files.class);
		when(Files.newOutputStream(keystorePath)).thenReturn(outputStream);

		CertAndKeyGen certAndKeyGen = new CertAndKeyGen("RSA", "SHA1WithRSA", null);
		certAndKeyGen.generate(2048);
		privateKey = certAndKeyGen.getPrivateKey();
		x509Certificate = certAndKeyGen.getSelfCertificate(new X500Name("CN=*.deliveryservice.cdn.example.com"), 600L);


		Rdn rdn = mock(Rdn.class);
		when(rdn.getType()).thenReturn("CN");
		when(rdn.getValue()).thenReturn("*.deliveryservice1.cdn.example.com");

		List<Rdn> rdnList = new ArrayList<>();
		rdnList.add(rdn);

		LdapName ldapName = mock(LdapName.class);
		when(ldapName.getRdns()).thenReturn(rdnList);

		whenNew(LdapName.class).withArguments("CN=*.deliveryservice1.cdn.example.com,OU=Traffic Routing,O=Traffic Control CDN,C=US").thenReturn(ldapName);

		X509Certificate x509Certificate = PowerMockito.mock(X509Certificate.class);
		X500Principal x500Principal = new X500Principal("CN=*.deliveryservice1.cdn.example.com, OU=Traffic Routing, O=Traffic Control CDN, C=US");

		when(x509Certificate.getSubjectX500Principal()).thenReturn(x500Principal);

		PowerMockito.when(keyStore.aliases()).thenAnswer(invocation -> {
			Vector<String> vector = new Vector<>();
			vector.add("alias-1");
			vector.add("alias-2");
			return vector.elements();
		});

		PowerMockito.when(keyStore.getCertificate("alias-1")).thenReturn(x509Certificate);
		PowerMockito.when(keyStore.getCertificate("alias-2")).thenReturn(mock(Certificate.class));
	}

	@Test
	public void itInitializesSingletonAndImportsCerts() throws Exception {
		KeyStoreHelper keyStoreHelper =  KeyStoreHelper.getInstance();

		verify(encoder).withoutPadding();
		verify(encoder).encodeToString(any(byte[].class));
		verify(secureRandom).nextBytes(any(byte[].class));

		assertThat(keyStoreHelper.getKeyPass(), equalTo("password".toCharArray()));
		assertThat(keyStoreHelper.getKeyStore(), equalTo(keyStore));

		boolean success = keyStoreHelper.importCertificateChain("deliveryservice.cdn.example.com",privateKey,new Certificate[] {x509Certificate});
		assertThat(success, equalTo(true));

		verify(keyStore).setKeyEntry("deliveryservice.cdn.example.com", privateKey, "password".toCharArray(), new Certificate[] {x509Certificate});

		assertThat(KeyStoreHelper.getInstance().getAllCommonNames(), containsInAnyOrder("*.deliveryservice1.cdn.example.com"));
	}
}
