/*
 * Copyright 2015 Comcast Cable Communications Management, LLC
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

package com.comcast.cdn.traffic_control.traffic_router.core.ds;

import org.junit.Test;

import java.util.Arrays;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.containsInAnyOrder;
import static org.hamcrest.Matchers.equalTo;

public class SteeringRegistryTest {
	@Test
	public void itConsumesValidJson() throws Exception {

		String json = "{ \"response\": [ " +
			"{ \"deliveryService\":\"steering-1\"," +
			"  \"targets\" : [" +
			"        {" +
			"            \"deliveryService\": \"ds-01\", \"weight\": 9876," +
			"            \"filters\": [ \".*/force-to-one/.*\", \".*/also-this/.*\" ]" +
			"        }," +
			"        {\"deliveryService\": \"ds-02\", \"weight\": 12345, \"filters\": []}" +
			"      ]" +
			"}, " +
			"{ \"deliveryService\":\"steering-2\"," +
			"  \"targets\" : [" +
			"        {" +
			"            \"deliveryService\": \"ds-03\", \"weight\": 1117," +
			"            \"filters\": [ \".*/three-for-me/.*\" ]" +
			"        }," +
			"        {\"deliveryService\": \"ds-02\", \"weight\": 556, \"filters\": []}" +
			"      ]" +
			"}" +
			"] }";

		SteeringRegistry steeringRegistry = new SteeringRegistry();
		steeringRegistry.update(json);
		assertThat(steeringRegistry.has("steering-1"), equalTo(true));
		assertThat(steeringRegistry.has("steering-2"), equalTo(true));

		SteeringTarget steeringTarget1 = new SteeringTarget();
		steeringTarget1.setDeliveryService("ds-01");
		steeringTarget1.setWeight(9876);
		steeringTarget1.setFilters(Arrays.asList(".*/force-to-one/.*", ".*/also-this/.*"));

		SteeringTarget steeringTarget2 = new SteeringTarget();
		steeringTarget2.setDeliveryService("ds-02");
		steeringTarget2.setWeight(12345);

		assertThat(steeringRegistry.get("steering-1").getTargets(), containsInAnyOrder(steeringTarget1, steeringTarget2));
		assertThat(steeringRegistry.get("steering-2").getTargets().get(1).getDeliveryService(), equalTo("ds-02"));
		assertThat(steeringRegistry.get("steering-2").getTargets().get(1).getFilters().isEmpty(), equalTo(true));
	}
}
