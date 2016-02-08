package application;

import com.comcast.cdn.traffic_control.traffic_monitor.MonitorApplication;
import org.junit.Before;
import org.junit.Ignore;
import org.junit.Test;

import java.security.AccessControlException;
import java.security.Permission;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.equalTo;
import static org.junit.Assert.fail;

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
			new MonitorApplication().init();
			fail("Init did not do SystemExit");
		} catch (AccessControlException e) {
			assertThat(e.getMessage(), equalTo("Boom"));
		} finally {
			System.setSecurityManager(originalSecurityManager);
		}
	}
}
