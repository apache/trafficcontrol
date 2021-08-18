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

package org.apache.traffic_control.traffic_router.core.ds;

import org.junit.Test;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.containsInAnyOrder;
import static org.hamcrest.Matchers.equalTo;
import static org.hamcrest.Matchers.nullValue;

public class SteeringRegistryTest {
	@Test
	public void itConsumesValidJson() throws Exception {

		String json = "{ \"response\": [ " +
			"{ \"deliveryService\":\"steering-1\"," +
				"  \"targets\" : [" +
				"        {\"deliveryService\": \"steering-target-01\", \"weight\": 9876}," +
				"        {\"deliveryService\": \"steering-target-02\", \"weight\": 12345}" +
				"      ]," +
				"  \"filters\" : [" +
				"      { \"pattern\" : \".*/force-to-one/.*\", \"deliveryService\" : \"steering-target-01\" }," +
				"      { \"pattern\" : \".*/also-this/.*\", \"deliveryService\" : \"steering-target-01\" }" +
				"   ]"+
				"}, " +
				"{ \"deliveryService\":\"steering-2\"," +
				"  \"targets\" : [" +
				"        {\"deliveryService\": \"steering-target-3\", \"weight\": 1117}," +
				"        {\"deliveryService\": \"steering-target-02\", \"weight\": 556}" +
				"      ]" +
				"}" +

			"] }";

		SteeringRegistry steeringRegistry = new SteeringRegistry();
		steeringRegistry.update(json);
		assertThat(steeringRegistry.has("steering-1"), equalTo(true));
		assertThat(steeringRegistry.has("steering-2"), equalTo(true));

		SteeringTarget steeringTarget1 = new SteeringTarget();
		steeringTarget1.setDeliveryService("steering-target-01");
		steeringTarget1.setWeight(9876);

		SteeringTarget steeringTarget2 = new SteeringTarget();
		steeringTarget2.setDeliveryService("steering-target-02");
		steeringTarget2.setWeight(12345);

		assertThat(steeringRegistry.get("steering-1").getTargets(), containsInAnyOrder(steeringTarget1, steeringTarget2));
		assertThat(steeringRegistry.get("steering-2").getTargets().get(1).getDeliveryService(), equalTo("steering-target-02"));

		assertThat(steeringRegistry.get("steering-1").getFilters().get(0).getPattern(), equalTo(".*/force-to-one/.*"));
		assertThat(steeringRegistry.get("steering-1").getFilters().get(0).getDeliveryService(), equalTo("steering-target-01"));
		assertThat(steeringRegistry.get("steering-1").getFilters().get(1).getPattern(), equalTo(".*/also-this/.*"));
		assertThat(steeringRegistry.get("steering-1").getFilters().get(1).getDeliveryService(), equalTo("steering-target-01"));

		assertThat(steeringRegistry.get("steering-1").getBypassDestination("/stuff/force-to-one/more/stuff"), equalTo("steering-target-01"));
		assertThat(steeringRegistry.get("steering-1").getBypassDestination("/should/not/match"), nullValue());

		assertThat(steeringRegistry.get("steering-2").getFilters().isEmpty(), equalTo(true));
	}
}
