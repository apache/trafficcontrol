package data;

/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 * 
 *   http://www.apache.org/licenses/LICENSE-2.0
 * 
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */


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
