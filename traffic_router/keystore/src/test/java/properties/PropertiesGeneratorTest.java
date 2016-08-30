package properties;

import com.comcast.cdn.traffic_control.traffic_router.properties.PropertiesGenerator;
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
import java.util.Properties;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.equalTo;
import static org.mockito.Mockito.doAnswer;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.verify;
import static org.powermock.api.mockito.PowerMockito.mockStatic;
import static org.powermock.api.mockito.PowerMockito.whenNew;

@RunWith(PowerMockRunner.class)
@PrepareForTest({PropertiesGenerator.class})
public class PropertiesGeneratorTest {

	public static final String TMP_EXAMPLE_BAZ_PROPERTIES = "/tmp/example/baz.properties";
	private Properties properties;
	private OutputStream newFileOutputStream;
	private OutputStream existingFileOutputStream;
	private InputStream inputStream;
	private Path propertiesPath;

	@Before
	public void before() throws Exception {
		propertiesPath = mock(Path.class);
		final Path createdPropertiesPath = mock(Path.class);

		mockStatic(Paths.class);
		PowerMockito.when(Paths.get(TMP_EXAMPLE_BAZ_PROPERTIES)).thenReturn(propertiesPath);

		mockStatic(Files.class);
		PowerMockito.when(Files.createFile(propertiesPath)).thenReturn(createdPropertiesPath);

		inputStream = mock(InputStream.class);
		PowerMockito.when(Files.newInputStream(propertiesPath)).thenReturn(inputStream);

		properties = mock(Properties.class);
		whenNew(Properties.class).withNoArguments().thenReturn(properties);

		existingFileOutputStream = mock(OutputStream.class);
		PowerMockito.when(Files.newOutputStream(propertiesPath)).thenReturn(existingFileOutputStream);

		newFileOutputStream = mock(OutputStream.class);
		PowerMockito.when(Files.newOutputStream(createdPropertiesPath)).thenReturn(newFileOutputStream);
	}

	@Test
	public void itReturnsValueFromExistingPropertiesFile() throws Exception {
		PowerMockito.when(Files.exists(propertiesPath)).thenReturn(true);

		doAnswer(invocation -> {
			PowerMockito.when(properties.getProperty("foo")).thenReturn("existingValue");
			return null;
		}).when(properties).load(inputStream);

		assertThat(new PropertiesGenerator(TMP_EXAMPLE_BAZ_PROPERTIES).getProperty("foo", "bar"), equalTo("existingValue"));
	}

	@Test
	public void itStoresDefaultToExistingPropertiesFile() throws Exception {
		PowerMockito.when(Files.exists(propertiesPath)).thenReturn(true);
		PowerMockito.when(properties.getProperty("foo")).thenReturn(null);

		assertThat(new PropertiesGenerator(TMP_EXAMPLE_BAZ_PROPERTIES).getProperty("foo", "bar"), equalTo("bar"));

		verify(properties).setProperty("foo", "bar");
		verify(properties).store(existingFileOutputStream, null);
	}

	@Test
	public void itStoresDefaultToNewPropertiesFile() throws Exception {
		PowerMockito.when(Files.exists(propertiesPath)).thenReturn(false);

		assertThat(new PropertiesGenerator(TMP_EXAMPLE_BAZ_PROPERTIES).getProperty("foo", "bar"), equalTo("bar"));

		verify(properties).setProperty("foo", "bar");
		verify(properties).store(newFileOutputStream, null);
	}
}
