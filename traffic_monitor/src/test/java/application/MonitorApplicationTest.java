package application;

import com.comcast.cdn.traffic_control.traffic_monitor.MonitorApplication;
import com.comcast.cdn.traffic_control.traffic_monitor.config.ConfigHandler;
import org.junit.Before;
import org.junit.Ignore;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.powermock.api.mockito.PowerMockito;
import org.powermock.core.classloader.annotations.PrepareForTest;
import org.powermock.modules.junit4.PowerMockRunner;

import java.security.AccessControlException;
import java.security.Permission;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.equalTo;
import static org.junit.Assert.fail;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;
import static org.powermock.api.mockito.PowerMockito.mockStatic;

@PrepareForTest(ConfigHandler.class)
@RunWith(PowerMockRunner.class)
public class MonitorApplicationTest {
	private final SecurityManager originalSecurityManager = System.getSecurityManager();

	@Before
	public void before() {
		System.setSecurityManager(new SecurityManager() {
			@Override
			public void checkPermission(Permission perm)
			{
				// allow anything.
			}

			@Override
			public void checkPermission(Permission perm, Object context)
			{
				// allow anything.
			}

			@Override
			public void checkExit(int status)
			{
				super.checkExit(status);
				throw new AccessControlException("Boom");
			}


		});
	}

	@Test
	public void itSystemExitsWhenConfigFileIsMissing() {
		try {
			ConfigHandler configHandler = mock(ConfigHandler.class);
			when(configHandler.configFileExists()).thenReturn(false);
			mockStatic(ConfigHandler.class);
			when(ConfigHandler.getInstance()).thenReturn(configHandler);

			new MonitorApplication().init();
			fail("Init did not do SystemExit");
		} catch (AccessControlException e) {
			assertThat(e.getMessage(), equalTo("Boom"));
		} finally {
			System.setSecurityManager(originalSecurityManager);
		}
	}
}
