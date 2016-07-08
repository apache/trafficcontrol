package keystore;

import com.comcast.cdn.traffic_control.traffic_router.keystore.KeyStoreLoader;
import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.powermock.api.mockito.PowerMockito;
import org.powermock.core.classloader.annotations.PrepareForTest;
import org.powermock.modules.junit4.PowerMockRunner;

import java.io.InputStream;
import java.io.OutputStream;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.security.KeyStore;

import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;
import static org.powermock.api.mockito.PowerMockito.mockStatic;

@RunWith(PowerMockRunner.class)
@PrepareForTest(KeyStoreLoader.class)
public class KeyStoreLoaderTest {

	private KeyStore keyStore;
	private Path keyStorePath;
	private InputStream inputStream;

	@Before
	public void before() throws Exception {
		keyStore = PowerMockito.mock(KeyStore.class);

		keyStorePath = mock(Path.class);

		mockStatic(Paths.class);
		when(Paths.get("/some/path/.keystore")).thenReturn(keyStorePath);

		inputStream = mock(InputStream.class);
		mockStatic(Files.class);
		when(Files.newInputStream(keyStorePath)).thenReturn(inputStream);

		mockStatic(KeyStore.class);
		when(KeyStore.getDefaultType()).thenReturn("secret-sauce");
		when(KeyStore.getInstance("secret-sauce")).thenReturn(keyStore);
	}

	@Test
	public void itLoadsExistingKeyStore() throws Exception {
		when(Files.exists(keyStorePath)).thenReturn(true);
		new KeyStoreLoader("/some/path/.keystore", "password".toCharArray()).load();
		verify(keyStore).load(inputStream, "password".toCharArray());
	}

	@Test
	public void itCreatesKeyStore() throws Exception {
		when(Files.exists(keyStorePath)).thenReturn(false);

		Path freshKeyStorePath = mock(Path.class);
		when(Files.createFile(keyStorePath)).thenReturn(freshKeyStorePath);

		OutputStream outputStream = mock(OutputStream.class);
		when(Files.newOutputStream(freshKeyStorePath)).thenReturn(outputStream);

		new KeyStoreLoader("/some/path/.keystore", "password".toCharArray()).load();
		verify(keyStore).load(null, "password".toCharArray());
		verify(keyStore).store(outputStream, "password".toCharArray());
	}
}