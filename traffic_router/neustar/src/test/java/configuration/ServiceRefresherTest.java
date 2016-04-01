package configuration;

import com.comcast.cdn.traffic_control.traffic_router.neustar.NeustarGeolocationService;
import com.comcast.cdn.traffic_control.traffic_router.neustar.configuration.ServiceRefresher;
import com.comcast.cdn.traffic_control.traffic_router.neustar.data.NeustarDatabaseUpdater;
import org.junit.Before;
import org.junit.Test;
import org.mockito.InjectMocks;
import org.mockito.Mock;

import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;
import static org.mockito.MockitoAnnotations.initMocks;

public class ServiceRefresherTest {
	@Mock
	NeustarDatabaseUpdater neustarDatabaseUpdater;

	@Mock
	NeustarGeolocationService neustarGeolocationService;

	@InjectMocks
	ServiceRefresher serviceRefresher;

	@Before
	public void before() {
		initMocks(this);
	}

	@Test
	public void itSwallowsExceptions() {
		when(neustarDatabaseUpdater.update()).thenThrow(new RuntimeException("Boom!"));
		serviceRefresher.run();
		verify(neustarDatabaseUpdater).update();
	}
}