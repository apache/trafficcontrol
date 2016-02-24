package health;

import com.comcast.cdn.traffic_control.traffic_monitor.config.Cache;
import com.comcast.cdn.traffic_control.traffic_monitor.health.CacheStateUpdater;
import com.comcast.cdn.traffic_control.traffic_monitor.health.CacheStatisticsClient;
import com.ning.http.client.AsyncHttpClient;
import com.ning.http.client.ListenableFuture;
import com.ning.http.client.ProxyServer;
import com.ning.http.client.Request;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.powermock.core.classloader.annotations.PrepareForTest;
import org.powermock.modules.junit4.PowerMockRunner;

import static org.mockito.Matchers.any;
import static org.mockito.Mockito.doReturn;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.spy;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;
import static org.powermock.api.mockito.PowerMockito.whenNew;

@PrepareForTest({CacheStatisticsClient.class, AsyncHttpClient.class, ProxyServer.class})
@RunWith(PowerMockRunner.class)
public class CacheStatisticsClientTest {
	@Test
	public void itExecutesAsynchronousRequest() throws Exception {

		ListenableFuture listenableFuture = mock(ListenableFuture.class);
		AsyncHttpClient asyncHttpClient = spy(new AsyncHttpClient());
		doReturn(listenableFuture).when(asyncHttpClient).executeRequest(any(Request.class), any(CacheStateUpdater.class));

		whenNew(AsyncHttpClient.class).withNoArguments().thenReturn(asyncHttpClient);

		Cache cache = mock(Cache.class);
		when(cache.getQueryIp()).thenReturn("192.168.99.100");
		when(cache.getQueryPort()).thenReturn(0);
		when(cache.getStatisticsUrl()).thenReturn("http://cache1.example.com/astats");

		CacheStateUpdater cacheStateUpdater = mock(CacheStateUpdater.class);
		CacheStatisticsClient cacheStatisticsClient = new CacheStatisticsClient();

		cacheStatisticsClient.fetchCacheStatistics(cache, cacheStateUpdater);
		verify(cacheStateUpdater).setFuture(listenableFuture);
	}
}
