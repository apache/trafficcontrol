package health;

import com.comcast.cdn.traffic_control.traffic_monitor.config.Cache;
import com.comcast.cdn.traffic_control.traffic_monitor.health.ProxiedRequest;
import com.ning.http.client.AsyncHttpClient;
import org.junit.Test;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.equalTo;
import static org.hamcrest.Matchers.not;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;

public class ProxiedRequestTest {
	@Test
	public void itBuildsARequest() {
		ProxiedRequest proxiedRequest = new ProxiedRequest("192.168.1.100", 9876, "http://example.com/index.html", new AsyncHttpClient());

		assertThat(proxiedRequest.getRequest().getProxyServer().getURI().toString(), equalTo("http://192.168.1.100:9876/"));
		assertThat(proxiedRequest.getRequest().getMethod(), equalTo("GET"));
		assertThat(proxiedRequest.getRequest().getUrl(), equalTo("http://example.com/index.html"));
	}

	@Test
	public void itUpdatesItselfFromCache() {
		ProxiedRequest proxiedRequest = new ProxiedRequest("192.168.1.100", 9876, "http://example.com/index.html", new AsyncHttpClient());

		Cache cache = mock(Cache.class);
		when(cache.getQueryIp()).thenReturn("192.168.1.100");
		when(cache.getQueryPort()).thenReturn(9876);
		when(cache.getStatisticsUrl()).thenReturn("http://example.com/index.html");

		ProxiedRequest updatedProxiedRequest = proxiedRequest.updateForCache(cache, new AsyncHttpClient());
		assertThat(updatedProxiedRequest.hashCode(), equalTo(proxiedRequest.hashCode()));

		when(cache.getQueryIp()).thenReturn("192.168.1.222");

		updatedProxiedRequest = proxiedRequest.updateForCache(cache, new AsyncHttpClient());
		assertThat(updatedProxiedRequest.hashCode(), not(equalTo(proxiedRequest.hashCode())));

		proxiedRequest = updatedProxiedRequest;

		assertThat(proxiedRequest.getIpAddress(), equalTo("192.168.1.222"));

		when(cache.getQueryPort()).thenReturn(2222);
		proxiedRequest = proxiedRequest.updateForCache(cache, new AsyncHttpClient());
		assertThat(proxiedRequest.getPort(), equalTo(2222));

		when(cache.getStatisticsUrl()).thenReturn("http://example.com/statistics.html");

		proxiedRequest = proxiedRequest.updateForCache(cache, new AsyncHttpClient());
		assertThat(proxiedRequest.getUrl(), equalTo("http://example.com/statistics.html"));

		when(cache.getQueryPort()).thenReturn(0);

		proxiedRequest = proxiedRequest.updateForCache(cache, new AsyncHttpClient());
		assertThat(proxiedRequest.getPort(), equalTo(80));
	}
}
