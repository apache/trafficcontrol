package com.comcast.cdn.traffic_control.traffic_router.core.hashing;

import com.comcast.cdn.traffic_control.traffic_router.core.hash.DefaultHashable;
import com.comcast.cdn.traffic_control.traffic_router.core.hash.NumberSearcher;
import org.junit.Before;
import org.junit.Test;
import org.mockito.InjectMocks;
import org.mockito.Mock;

import static org.hamcrest.Matchers.equalTo;
import static org.hamcrest.core.IsNot.not;
import static org.junit.Assert.assertThat;
import static org.mockito.MockitoAnnotations.initMocks;

public class HashableTest {
	@Mock
	private NumberSearcher numberSearcher = new NumberSearcher();

	@InjectMocks
	private DefaultHashable defaultHashable;

	@Before
	public void before() {
		initMocks(this);
	}

	@Test
	public void itReturnsClosestHash() {
		defaultHashable.generateHashes("hash id", 100);
		double hash = defaultHashable.getClosestHash(1.23);

		assertThat(hash, not(equalTo(0.0)));
		assertThat(defaultHashable.getClosestHash(1.23), equalTo(hash));
	}

}
