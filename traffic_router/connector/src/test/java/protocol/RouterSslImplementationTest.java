package protocol;

import com.comcast.cdn.traffic_control.traffic_router.protocol.RouterSslImplementation;
import com.comcast.cdn.traffic_control.traffic_router.protocol.RouterSslServerSocketFactory;
import org.apache.tomcat.util.net.jsse.JSSEFactory;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.powermock.core.classloader.annotations.PrepareForTest;
import org.powermock.modules.junit4.PowerMockRunner;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.equalTo;
import static org.hamcrest.Matchers.instanceOf;

@RunWith(PowerMockRunner.class)
@PrepareForTest({RouterSslImplementation.class, JSSEFactory.class})
public class RouterSslImplementationTest {
	@Test
	public void itReturnsItsName() {
		assertThat(new RouterSslImplementation().getImplementationName(), equalTo("RouterSslImplementation"));
	}

	@Test
	public void itReturnsRouterSslServerSocketFactory() {
		assertThat(new RouterSslImplementation().getServerSocketFactory(), instanceOf(RouterSslServerSocketFactory.class));
	}
}
