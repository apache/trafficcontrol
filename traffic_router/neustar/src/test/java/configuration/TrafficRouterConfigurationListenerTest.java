/*
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package configuration;

import org.apache.traffic_control.traffic_router.neustar.configuration.ServiceRefresher;
import org.apache.traffic_control.traffic_router.neustar.configuration.TrafficRouterConfigurationListener;
import org.junit.Before;
import org.junit.Test;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.invocation.InvocationOnMock;
import org.mockito.stubbing.Answer;
import org.springframework.core.env.Environment;

import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.ScheduledFuture;
import java.util.concurrent.TimeUnit;

import static org.mockito.ArgumentMatchers.any;
import static org.mockito.ArgumentMatchers.eq;
import static org.mockito.Mockito.doReturn;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.times;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.verifyZeroInteractions;
import static org.mockito.Mockito.when;
import static org.mockito.MockitoAnnotations.initMocks;

public class TrafficRouterConfigurationListenerTest {
	@Mock
	ScheduledExecutorService scheduledExecutorService;

	@Mock
	Environment environment;

	@Mock
	ServiceRefresher serviceRefresher;

	@InjectMocks
	TrafficRouterConfigurationListener trafficRouterConfigurationListener;

	@Before
	public void before() {
		initMocks(this);

		when(environment.getProperty("neustar.polling.interval", Long.class, 86400000L)).thenReturn(86400000L);
	}

	@Test
	public void itCancelsExistingTaskBeforeStartingAnother() {

		ScheduledFuture scheduledFuture = mock(ScheduledFuture.class);

		when(scheduledFuture.isDone()).thenAnswer(new Answer<Boolean>() {
			int doneCheckCount = 0;

			@Override
			public Boolean answer(InvocationOnMock invocation) throws Throwable {
				doneCheckCount++;
				return doneCheckCount > 3;
			}
		});

		doReturn(scheduledFuture).when(scheduledExecutorService).scheduleAtFixedRate(any(Runnable.class), eq(0L), eq(86400000L), eq(TimeUnit.MILLISECONDS));

		trafficRouterConfigurationListener.configurationChanged();
		verifyZeroInteractions(scheduledFuture);

		trafficRouterConfigurationListener.configurationChanged();
		verify(scheduledFuture).cancel(true);
		verify(scheduledFuture, times(4)).isDone();
	}
}