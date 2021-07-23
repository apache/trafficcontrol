/*
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package org.apache.traffic_control.traffic_router.core.hashing;

import org.apache.traffic_control.traffic_router.core.hash.DefaultHashable;
import org.apache.traffic_control.traffic_router.core.hash.NumberSearcher;
import org.junit.Before;
import org.junit.Test;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.MockitoAnnotations;

import static org.hamcrest.Matchers.equalTo;
import static org.hamcrest.core.IsNot.not;
import static org.junit.Assert.assertThat;

public class HashableTest {
	@Mock
	private NumberSearcher numberSearcher = new NumberSearcher();

	@InjectMocks
	private DefaultHashable defaultHashable;

	@Before
	public void before() {
		MockitoAnnotations.openMocks(this);
	}

	@Test
	public void itReturnsClosestHash() {
		defaultHashable.generateHashes("hash id", 100);
		double hash = defaultHashable.getClosestHash(1.23);

		assertThat(hash, not(equalTo(0.0)));
		assertThat(defaultHashable.getClosestHash(1.23), equalTo(hash));
	}

}
