package health;

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


import com.comcast.cdn.traffic_control.traffic_monitor.config.Cache;
import com.comcast.cdn.traffic_control.traffic_monitor.health.HealthDeterminer;
import org.apache.wicket.ajax.json.JSONObject;
import org.junit.Test;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.equalTo;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;

public class HealthDeterminerTest {
	@Test
	public void itMarksCacheWithoutStateAsHavingNoError() throws Exception {
		HealthDeterminer healthDeterminer = new HealthDeterminer();
		Cache cache = mock(Cache.class);
		when(cache.getStatus()).thenReturn("okay");
		JSONObject statsJson = healthDeterminer.getJSONStats(cache, true, true);

		assertThat(statsJson.getString(HealthDeterminer.ERROR_STRING), equalTo(HealthDeterminer.NO_ERROR_FOUND));
		assertThat(statsJson.getString(HealthDeterminer.STATUS), equalTo("okay"));
	}
}
