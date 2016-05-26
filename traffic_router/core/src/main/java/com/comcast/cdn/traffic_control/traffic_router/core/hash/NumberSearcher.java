package com.comcast.cdn.traffic_control.traffic_router.core.hash;

import java.util.Arrays;

public class NumberSearcher {
	public static int findClosest(final Double[] numbers, final double target) {
		final int index = Arrays.binarySearch(numbers, target);
		if (index >= 0) {
			return index;
		}

		final int biggerThanIndex = -(index + 1);
		if (biggerThanIndex == numbers.length) {
			return numbers.length - 1;
		}

		if (biggerThanIndex == 0) {
			return 0;
		}

		final int smallerThanIndex = biggerThanIndex - 1;

		final double biggerThanDelta = Math.abs(numbers[biggerThanIndex] - target);
		final double smallerThanDelta = Math.abs(numbers[smallerThanIndex] - target);

		return (biggerThanDelta < smallerThanDelta) ? biggerThanIndex : smallerThanIndex;
	}
}
