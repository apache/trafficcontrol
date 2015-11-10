package com.comcast.cdn.traffic_control.traffic_router.core.util;

import org.junit.Test;

import java.util.Iterator;
import java.util.Set;
import java.util.TreeSet;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.equalTo;

public class ComparableStringByLengthTest {
	@Test
	public void itSortsAscendingToShorterStrings() {
		String[] strings = new String[] {
			"a", "ba", "b", "bac", "ab", "abc"
		};

		Set set = new TreeSet();
		for (String string : strings) {
			set.add(new ComparableStringByLength(string));
		}

		Iterator<ComparableStringByLength> iterator = set.iterator();

		assertThat(iterator.next().toString(), equalTo("abc"));
		assertThat(iterator.next().toString(), equalTo("bac"));
		assertThat(iterator.next().toString(), equalTo("ab"));
		assertThat(iterator.next().toString(), equalTo("ba"));
		assertThat(iterator.next().toString(), equalTo("a"));
		assertThat(iterator.next().toString(), equalTo("b"));
	}
}
