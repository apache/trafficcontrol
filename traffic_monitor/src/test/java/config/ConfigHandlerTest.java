package config;

import com.comcast.cdn.traffic_control.traffic_monitor.config.ConfigHandler;
import com.comcast.cdn.traffic_control.traffic_monitor.config.MonitorConfig;
import org.apache.commons.io.IOUtils;
import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.powermock.api.mockito.PowerMockito;
import org.powermock.core.classloader.annotations.PrepareForTest;
import org.powermock.modules.junit4.PowerMockRunner;
import org.powermock.reflect.Whitebox;

import java.io.File;
import java.io.FileReader;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.equalTo;
import static org.hamcrest.Matchers.nullValue;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;
import static org.powermock.api.mockito.PowerMockito.whenNew;

@PrepareForTest({ConfigHandler.class,File.class,IOUtils.class})
@RunWith(PowerMockRunner.class)
public class ConfigHandlerTest {
	private File mockConfigFile;
	private File mockVarDirectory;
	private ConfigHandler configHandler;

	@Before
	public void before() throws Exception {
		mockConfigFile = mock(File.class);
		whenNew(File.class).withArguments("/opt/traffic_monitor/conf/traffic_monitor_config.js").thenReturn(mockConfigFile);

		mockVarDirectory = mock(File.class);
		whenNew(File.class).withArguments("/opt/traffic_monitor/var").thenReturn(mockVarDirectory);

		configHandler = ConfigHandler.getInstance();

		// Some food for thought about why "true" singletons are generally less desirable than dependency injection...
		// Without injecting a null monitor config object behind the scenes the tests don't work
		MonitorConfig monitorConfig = null;
		Whitebox.setInternalState(configHandler, "config", monitorConfig);
		Whitebox.setInternalState(configHandler, "configFile", mockConfigFile);
		Whitebox.setInternalState(configHandler, "varDirectory", mockVarDirectory);
	}

	@Test
	public void itChecksTheFileSystemForTrafficMonitorConfig() throws Exception {
		assertThat(configHandler.getConfFile(),nullValue());

		when(mockConfigFile.exists()).thenReturn(true);
		assertThat(configHandler.getConfFile(),equalTo("/opt/traffic_monitor/conf/traffic_monitor_config.js"));
	}

	@Test
	public void itReturnsMonitorConfigWithDefaults() {
		String configUrl = configHandler.getConfig().getCrConfigUrl();
		assertThat(configUrl, equalTo("https://${tmHostname}/CRConfig-Snapshots/${cdnName}/CRConfig.json"));
	}

	@Test
	public void itBuildsMonitorConfigFromJsonString() throws Exception {
		when(mockConfigFile.exists()).thenReturn(true);

		FileReader fileReader = mock(FileReader.class);
		whenNew(FileReader.class).withAnyArguments().thenReturn(fileReader);
		PowerMockito.mockStatic(IOUtils.class);
		when(IOUtils.toString(fileReader)).thenReturn("{\"traffic_monitor_config\":" +
			"{\"tm.crConfig.json.polling.url\": \"https://trafficops.kabletown.com/somepath/kabletown/config.json\"}" +
		"}");

		String configUrl = configHandler.getConfig().getCrConfigUrl();

		// Not trying to test that MonitorConfig does fancy stuff with string replacement
		assertThat(configUrl, equalTo("https://trafficops.kabletown.com/somepath/kabletown/config.json"));
	}

	@Test
	public void itReportsTheDatabaseDirectory() {
		assertThat(configHandler.getDbDir(), nullValue());

		when(mockVarDirectory.exists()).thenReturn(true);
		assertThat(configHandler.getDbDir(), equalTo("/opt/traffic_monitor/db/"));
	}
}