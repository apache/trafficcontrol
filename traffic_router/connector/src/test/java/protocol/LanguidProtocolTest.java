package protocol;

import com.comcast.cdn.traffic_control.traffic_router.keystore.KeyStoreHelper;
import com.comcast.cdn.traffic_control.traffic_router.protocol.LanguidProtocol;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.powermock.core.classloader.annotations.PrepareForTest;
import org.powermock.modules.junit4.PowerMockRunner;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.equalTo;
import static org.junit.Assert.fail;
import static org.powermock.api.mockito.PowerMockito.mockStatic;
import static org.powermock.api.mockito.PowerMockito.when;

@PrepareForTest({KeyStoreHelper.class})
@RunWith(PowerMockRunner.class)
public class LanguidProtocolTest {

	@Test
	public void itUsesKeyStoreHelperOnInit() throws Exception {
		mockStatic(KeyStoreHelper.class);
		when(KeyStoreHelper.getInstance()).thenThrow(new RuntimeException("boom"));
		try {
			LanguidProtocol languidProtocol = new LanguidProtocol();
			languidProtocol.setSSLEnabled(true);
			languidProtocol.init();
			fail("KeyStoreHelper.getInstance() not called");
		}catch (RuntimeException e) {
			assertThat(e.getMessage(), equalTo("boom"));
		}
	}
}
