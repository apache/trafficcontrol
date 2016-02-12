package data;

import com.comcast.cdn.traffic_control.traffic_monitor.data.DataPoint;
import org.junit.Test;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.equalTo;

public class DataPointTest {
	@Test
	public void itMatchesAgainstValues() {
		DataPoint dataPoint = new DataPoint("value", 1L);

		assertThat(dataPoint.matches("value"), equalTo(true));
		assertThat(dataPoint.matches("somethingelse"), equalTo(false));
		assertThat(dataPoint.matches(""), equalTo(false));
		assertThat(dataPoint.matches(null), equalTo(false));
	}

	@Test
	public void itMatchesWhenValueIsNull() {
		DataPoint dataPoint = new DataPoint(null, 1L);
		assertThat(dataPoint.matches("value"), equalTo(false));
		assertThat(dataPoint.matches(""), equalTo(false));
		assertThat(dataPoint.matches(null), equalTo(true));
	}

	@Test
	public void itUpdatesSpan() {
		DataPoint dataPoint = new DataPoint("something", 100L);
		assertThat(dataPoint.getSpan(), equalTo(1));
		assertThat(dataPoint.getIndex(), equalTo(100L));

		dataPoint.update(200L);
		assertThat(dataPoint.getIndex(), equalTo(200L));
		assertThat(dataPoint.getSpan(), equalTo(2));
	}
}
