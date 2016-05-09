package com.comcast.cdn.traffic_control.traffic_router.core.hashing;

import com.comcast.cdn.traffic_control.traffic_router.core.hash.ConsistentHasher;
import com.comcast.cdn.traffic_control.traffic_router.core.hash.DefaultHashable;
import com.comcast.cdn.traffic_control.traffic_router.core.hash.MD5HashFunction;
import com.comcast.cdn.traffic_control.traffic_router.core.hash.NumberSearcher;
import org.junit.Before;
import org.junit.Test;
import org.mockito.InjectMocks;
import org.mockito.Mock;

import java.util.ArrayList;
import java.util.List;

import static org.hamcrest.core.AnyOf.anyOf;
import static org.hamcrest.core.IsEqual.equalTo;
import static org.junit.Assert.assertThat;
import static org.mockito.MockitoAnnotations.initMocks;

public class ConsistentHasherTest {
	@Mock
	MD5HashFunction md5HashFunction = new MD5HashFunction();

	@Mock
	NumberSearcher numberSearcher = new NumberSearcher();

	@InjectMocks
	DefaultHashable hashable1 = new DefaultHashable();

	@InjectMocks
	DefaultHashable hashable2 = new DefaultHashable();

	@InjectMocks
	DefaultHashable hashable3 =  new DefaultHashable();

	List<DefaultHashable> hashables = new ArrayList<DefaultHashable>();

	@InjectMocks
	ConsistentHasher consistentHasher;

	@Before
	public void before() {
		hashable1.generateHashes("hashId1", 100);
		hashable2.generateHashes("hashId2", 100);
		hashable3.generateHashes("hashId3", 100);

		hashables.add(hashable1);
		hashables.add(hashable2);
		hashables.add(hashable3);

		initMocks(this);
	}

	@Test
	public void itHashes() throws Exception {
		DefaultHashable hashable = consistentHasher.selectHashable(hashables, "some-string", false);
		assertThat(hashable, anyOf(equalTo(hashable1), equalTo(hashable2), equalTo(hashable3)));
		DefaultHashable nextHashable = consistentHasher.selectHashable(hashables, "some-string", false);
		assertThat(nextHashable, equalTo(hashable));
	}

	@Test
	public void itHashesMoreThanOne() throws Exception {
		List<DefaultHashable> results = consistentHasher.selectHashables(hashables, 2, "some-string", false);
		assertThat(results.size(), equalTo(2));
		assertThat(results.get(0), anyOf(equalTo(hashable1), equalTo(hashable2), equalTo(hashable3)));
		assertThat(results.get(1), anyOf(equalTo(hashable1), equalTo(hashable2), equalTo(hashable3)));

		assertThat(consistentHasher.selectHashables(hashables, 2, "some-string", false), equalTo(results));

	}
}
