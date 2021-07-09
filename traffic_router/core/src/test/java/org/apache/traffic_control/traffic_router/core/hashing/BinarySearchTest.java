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
