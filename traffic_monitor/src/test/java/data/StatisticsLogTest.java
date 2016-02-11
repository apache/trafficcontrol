package data;

import com.comcast.cdn.traffic_control.traffic_monitor.data.StatisticsLog;
import org.junit.Test;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.containsInAnyOrder;
import static org.hamcrest.Matchers.equalTo;
import static org.hamcrest.Matchers.nullValue;

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

	}
}
