package data;

import com.comcast.cdn.traffic_control.traffic_monitor.data.StatisticsLog;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.powermock.core.classloader.annotations.PrepareForTest;
import org.powermock.modules.junit4.PowerMockRunner;

import java.util.Arrays;
import java.util.HashSet;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.containsInAnyOrder;
import static org.hamcrest.Matchers.equalTo;
import static org.hamcrest.Matchers.nullValue;
import static org.powermock.api.mockito.PowerMockito.mockStatic;
import static org.powermock.api.mockito.PowerMockito.when;

@PrepareForTest({StatisticsLog.class, System.class})
@RunWith(PowerMockRunner.class)
public class StatisticsLogTest {
	@Test
	public void itLogsDataPoints() {
		StatisticsLog statisticsLog = new StatisticsLog();
		assertThat(statisticsLog.get("foo"), nullValue());

		statisticsLog.putDataPoint("foo", "okay");
		assertThat(statisticsLog.get("foo").getLast().getValue(), equalTo("okay"));
		assertThat(statisticsLog.get("foo").getLast().getIndex(), equalTo(0L));
		assertThat(statisticsLog.get("foo").getLast().getSpan(), equalTo(1));

		statisticsLog.putDataPoint("foo", "okay");
		assertThat(statisticsLog.get("foo").size(), equalTo(1));
		assertThat(statisticsLog.get("foo").getLast().getIndex(), equalTo(0L));
		assertThat(statisticsLog.get("foo").getLast().getSpan(), equalTo(2));

		statisticsLog.putDataPoint("foo", "bar");
		assertThat(statisticsLog.get("foo").size(), equalTo(2));

		statisticsLog.putDataPoint("foo2", "okay");
		assertThat(statisticsLog.getKeys(), containsInAnyOrder("foo2", "foo"));

		assertThat(statisticsLog.hasValue("foo"), equalTo(true));
		assertThat(statisticsLog.hasValue("bar"), equalTo(false));

		assertThat(statisticsLog.getLastValue("foo"), equalTo("bar"));
		assertThat(statisticsLog.getLastValue("baz"), nullValue());

		assertThat(statisticsLog.getValue("foo", 0), equalTo("bar"));
		assertThat(statisticsLog.getValue("baz", 0), nullValue());

		assertThat(statisticsLog.getBool("foo"), equalTo(false));
		assertThat(statisticsLog.getLong("foo"), equalTo(0L));
		assertThat(statisticsLog.getDouble("foo"), equalTo(0.0));

		assertThat(statisticsLog.getTime(0L), equalTo(0L));
	}

	@Test
	public void itIndexesPerUpdate() {
		mockStatic(System.class);
		when(System.currentTimeMillis()).thenReturn(1455224271177L);

		StatisticsLog statisticsLog = new StatisticsLog();
		statisticsLog.putDataPoint("kbps","1234.5");

		statisticsLog.prepareForUpdate("cache1", 1000L);

		assertThat(statisticsLog.getValue("kbps", 0), equalTo("1234.5"));
		assertThat(statisticsLog.getValue("kbps", 1), nullValue());

		assertThat(statisticsLog.getTime(0L), equalTo(0L));
		assertThat(statisticsLog.getTime(1L), equalTo(1455224271177L));

		statisticsLog.putDataPoint("kbps", "5432.1");

		assertThat(statisticsLog.getValue("kbps", 1), equalTo("5432.1"));
	}

	@Test
	public void itKeepsAtLeastOneStatistic() {
		mockStatic(System.class);
		when(System.currentTimeMillis()).thenReturn(1455224271177L);

		StatisticsLog statisticsLog = new StatisticsLog();
		statisticsLog.prepareForUpdate("cache1", 1000L);
		statisticsLog.putDataPoint("kbps","1234.5");
		statisticsLog.putDataPoint("kbps","5432.1");

		when(System.currentTimeMillis()).thenReturn(1455224271177L + 1001L);

		statisticsLog.prepareForUpdate("cache1", 1000L);

		assertThat(statisticsLog.getValue("kbps", 1), equalTo("5432.1"));
	}

	@Test
	public void itAddsPlaceHolder() {
		mockStatic(System.class);
		when(System.currentTimeMillis()).thenReturn(1455224271177L);

		StatisticsLog statisticsLog = new StatisticsLog();
		statisticsLog.prepareForUpdate("cache1", 1000L);
		statisticsLog.putDataPoint("kbps","1234.5");

		when(System.currentTimeMillis()).thenReturn(1455224271177L + 100L);
		statisticsLog.prepareForUpdate("cache1", 1000L);
		statisticsLog.putDataPoint("tps", "123000");

		when(System.currentTimeMillis()).thenReturn(1455224271177L + 200L);
		statisticsLog.prepareForUpdate("cache1", 1000L);

		assertThat(statisticsLog.getLastValue("kbps"), nullValue());
		assertThat(statisticsLog.getLastValue("tps"), equalTo("123000"));
	}


	@Test
	public void itFilters() {
		mockStatic(System.class);
		when(System.currentTimeMillis()).thenReturn(1455224271177L);

		StatisticsLog statisticsLog = new StatisticsLog();
		statisticsLog.prepareForUpdate("cache1", 1000L);
		statisticsLog.putDataPoint("kbps","1234.5");

		assertThat(statisticsLog.filter(10, new String[] {"tps"},false, false).size(), equalTo(0));
		assertThat(statisticsLog.filter(10, new String[] {"kbps", "tps"},false, false).get("kbps").getFirst().getValue(), equalTo("1234.5"));
		assertThat(statisticsLog.filter(10, new String[] {"KBPS"},true, false).get("kbps").getFirst().getValue(), equalTo("1234.5"));

		statisticsLog.addHiddenStats(new HashSet(Arrays.asList("kbps", "errors")));
		assertThat(statisticsLog.filter(10, new String[] {"kbps", "tps"},false, false).size(), equalTo(0));
		assertThat(statisticsLog.filter(10, new String[] {"kbps", "tps"},false, true).get("kbps").getFirst().getValue(), equalTo("1234.5"));

		statisticsLog.putDataPoint("kbps", "5432.1");
		assertThat(statisticsLog.filter(1, new String[] {"kbps", "tps"},false, true).get("kbps").size(), equalTo(1));
		assertThat(statisticsLog.filter(0, new String[] {"kbps", "tps"},false, true).get("kbps").size(), equalTo(2));

		statisticsLog.putDataPoint("rxBytes", "1000");

		assertThat(statisticsLog.filter(0, new String[] {"kbps"},false, true).keySet().size(), equalTo(1));
		assertThat(statisticsLog.filter(0, null,false, true).keySet().size(), equalTo(2));
		assertThat(statisticsLog.filter(0, new String[] {},false, true).keySet().size(), equalTo(0));
	}

	@Test
	public void itGetsValuesForIndexGoo() {
		mockStatic(System.class);

		StatisticsLog statisticsLog = new StatisticsLog();

		when(System.currentTimeMillis()).thenReturn(1455224271177L);
		statisticsLog.prepareForUpdate("cache1", 1000L);

		when(System.currentTimeMillis()).thenReturn(1455224271177L + 100L);
		statisticsLog.prepareForUpdate("cache1", 1000L);

		when(System.currentTimeMillis()).thenReturn(1455224271177L + 200L);
		statisticsLog.prepareForUpdate("cache1", 1000L);

		statisticsLog.putDataPoint("kbps","1234.5");

		when(System.currentTimeMillis()).thenReturn(1455224271177L + 300L);
		statisticsLog.prepareForUpdate("cache1", 1000L);
		statisticsLog.putDataPoint("kbps","1234.5");

		assertThat(statisticsLog.getValue("kbps", 1), nullValue());
		assertThat(statisticsLog.getValue("kbps", 2), nullValue());
		assertThat(statisticsLog.getValue("kbps", 3), equalTo("1234.5"));
		assertThat(statisticsLog.getValue("kbps", 4), equalTo("1234.5"));
		assertThat(statisticsLog.getValue("kbps", 5), nullValue());
	}
}
