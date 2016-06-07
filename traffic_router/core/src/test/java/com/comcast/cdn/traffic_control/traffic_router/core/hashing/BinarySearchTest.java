package com.comcast.cdn.traffic_control.traffic_router.core.hashing;

import org.junit.Test;

import java.util.Arrays;

import static org.hamcrest.core.IsEqual.equalTo;
import static org.junit.Assert.assertThat;

public class BinarySearchTest {
	@Test
	public void itReturnsMatchingIndex() {
		double[] hashes = new double[] {1.0, 2.0, 3.0, 4.0};
		assertThat(Arrays.binarySearch(hashes, 3.0), equalTo(2));
	}

	@Test
	public void itReturnsInsertionPoint() {
		double[] hashes = new double[] {1.0, 2.0, 3.0, 4.0};
		assertThat(Arrays.binarySearch(hashes, 3.5), equalTo(-4));
		assertThat(Arrays.binarySearch(hashes, 4.01), equalTo(-5));
	}
}
