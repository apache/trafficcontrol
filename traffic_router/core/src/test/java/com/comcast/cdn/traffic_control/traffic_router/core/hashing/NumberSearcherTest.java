package com.comcast.cdn.traffic_control.traffic_router.core.hashing;

import com.comcast.cdn.traffic_control.traffic_router.core.hash.NumberSearcher;
import org.junit.Test;

import static org.hamcrest.core.IsEqual.equalTo;
import static org.junit.Assert.assertThat;

public class NumberSearcherTest {
	@Test
	public void itFindsClosest() {
		Double[] numbers = { 1.2, 2.3, 3.4, 4.5, 5.6 };

		NumberSearcher numberSearcher = new NumberSearcher();
		assertThat(numberSearcher.findClosest(numbers,3.4), equalTo(2));
		assertThat(numberSearcher.findClosest(numbers,1.9), equalTo(1));
		assertThat(numberSearcher.findClosest(numbers,1.3), equalTo(0));
		assertThat(numberSearcher.findClosest(numbers,6.7), equalTo(4));
		assertThat(numberSearcher.findClosest(numbers,0.1), equalTo(0));
	}
}
