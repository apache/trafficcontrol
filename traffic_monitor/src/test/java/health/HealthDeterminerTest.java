package health;

import com.comcast.cdn.traffic_control.traffic_monitor.config.Cache;
import com.comcast.cdn.traffic_control.traffic_monitor.health.HealthDeterminer;
import org.apache.wicket.ajax.json.JSONObject;
import org.junit.Test;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.equalTo;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;

public class HealthDeterminerTest {
	@Test
	public void itMarksCacheWithoutStateAsHavingNoError() throws Exception {
		HealthDeterminer healthDeterminer = new HealthDeterminer();
		Cache cache = mock(Cache.class);
		when(cache.getStatus()).thenReturn("okay");
		JSONObject statsJson = healthDeterminer.getJSONStats(cache, true, true);

		assertThat(statsJson.getString(HealthDeterminer.ERROR_STRING), equalTo(HealthDeterminer.NO_ERROR_FOUND));
		assertThat(statsJson.getString(HealthDeterminer.STATUS), equalTo("okay"));
	}
}
