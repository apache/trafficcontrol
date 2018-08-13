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
package com.comcast.cdn.traffic_control.traffic_router.core.config;

import com.comcast.cdn.traffic_control.traffic_router.core.ds.DeliveryService;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.apache.commons.io.IOUtils;
import org.junit.Before;
import org.junit.Test;

import java.io.InputStream;
import java.util.List;
import java.util.Map;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.notNullValue;
import static org.junit.Assert.fail;

public class SnapshotEventsProcessorTest {
	private JsonNode newDsSnapJo;
	private JsonNode baselineJo;
	private JsonNode updateJo;

	@Before
	public void setUp() throws Exception {
		String resourcePath = "unit/ExistingConfig.json";
		InputStream inputStream = getClass().getClassLoader().getResourceAsStream(resourcePath);
		if (inputStream == null) {
			fail("Could not find file '" + resourcePath + "' needed for test from the current classpath as a resource!");
		}
		String baseDb = IOUtils.toString(inputStream);

		resourcePath = "unit/UpdateDsSnap.json";
		inputStream = getClass().getClassLoader().getResourceAsStream(resourcePath);
		if (inputStream == null) {
			fail("Could not find file '" + resourcePath + "' needed for test from the current classpath as a resource!");
		}
		String updateDb = IOUtils.toString(inputStream);

		resourcePath = "unit/NewDsSnap.json";
		inputStream = getClass().getClassLoader().getResourceAsStream(resourcePath);
		if (inputStream == null) {
			fail("Could not find file '" + resourcePath + "' needed for test from the current classpath as a resource!");
		}
		String newDsSnapDb = IOUtils.toString(inputStream);

		final ObjectMapper mapper = new ObjectMapper();
		assertThat(newDsSnapDb, notNullValue());
		assertThat(baseDb, notNullValue());
		assertThat(updateDb, notNullValue());

		newDsSnapJo = mapper.readTree(newDsSnapDb);
		assertThat(newDsSnapJo, notNullValue());
		updateJo = mapper.readTree(updateDb);
		assertThat(updateJo, notNullValue());
		baselineJo = mapper.readTree(baseDb);
		assertThat(baselineJo, notNullValue());
	}

	@Test
	public void mineEventsFromDBDiffsNoChanges() throws Exception {
		SnapshotEventsProcessor snapEvents = SnapshotEventsProcessor.diffCrConfigs(baselineJo, null);
		assertThat("Initialize should be true because the snapshot does not have a snapshot config parameter.",
				snapEvents.shouldReloadConfig());
		assertThat("18 Delivery services should have been loaded but there were only "+snapEvents.getCreationEvents().size(), snapEvents.getCreationEvents().size() == 18);
		snapEvents = SnapshotEventsProcessor.diffCrConfigs(baselineJo, baselineJo);
		assertThat("No new, updated or deleted delivery services should have been loaded beacause the snapshots were " +
						"the same: New = " + snapEvents.getCreationEvents().size() + ", updated = "+ snapEvents.getUpdateEvents().size() + ", deleted = "+ snapEvents.getDeleteEvents().size(),
				(snapEvents.getCreationEvents().size() == 0 && snapEvents.getUpdateEvents().size() == 0 && snapEvents.getDeleteEvents().size() == 0));
		assertThat("No other events should have been loaded beacause the snapshots were the same: mapping = " +
						snapEvents.getMappingEvents().size()+ ", cache = "+ snapEvents.getDeleteCacheEvents().size() +
						", ssl = "+ snapEvents.getSSLEnabledChangeEvents().size(),
				(snapEvents.getMappingEvents().size() == 0 && snapEvents.getDeleteCacheEvents().size() == 0 && snapEvents.getSSLEnabledChangeEvents().size() == 0));
	}

	@Test
	public void mineEventsFromDBNewUpdate() throws Exception {
		int changeCnt = 15;
		ConfigHandler.setLastSnapshotTimestamp(14650848001l);
		SnapshotEventsProcessor snapEvents = SnapshotEventsProcessor.diffCrConfigs(updateJo,
				baselineJo);
		assertThat("Initialize should be false.", !snapEvents.shouldReloadConfig());
		assertThat(String.valueOf(changeCnt)+" Delivery services should have been updated but there were only "+snapEvents.getChangeEvents().size(),
				snapEvents.getChangeEvents().size() == changeCnt);
		assertThat("1 links should have been updated but there were only "+snapEvents.getMappingEvents().size(),
				snapEvents.getMappingEvents().size() == 1);
	}

	@Test
	public void mineEventsFromDBDiffsNewDs() throws Exception {

		ConfigHandler.setLastSnapshotTimestamp(14650848001l);
		SnapshotEventsProcessor snapEvents = SnapshotEventsProcessor.diffCrConfigs(newDsSnapJo, baselineJo);
		assertThat("Initialize should be false.", !snapEvents.shouldReloadConfig());
		assertThat("1 Delivery services should have been added but there was "+snapEvents.getCreationEvents().size(),
				snapEvents.getCreationEvents().size() == 1);
		assertThat("4 links should have been updated but there were only "+snapEvents.getMappingEvents().size(),
				snapEvents.getMappingEvents().size() == 4);
	}

	@Test
	public void diffCrConfigNoChanges() throws Exception {
		ConfigHandler.setLastSnapshotTimestamp(14650848001l);
		SnapshotEventsProcessor snapEvents = SnapshotEventsProcessor.diffCrConfigs(updateJo, updateJo);
		assertThat("Initialize should be false.", !snapEvents.shouldReloadConfig());
		assertThat("0 Delivery services should have been added but there was "+snapEvents.getCreationEvents().size(),
				snapEvents.getCreationEvents().size() == 0);

	}

	@Test
	public void getSSLEnabledChangeEvents_updated() throws Exception {
		final SnapshotEventsProcessor sep = SnapshotEventsProcessor.diffCrConfigs(updateJo, baselineJo);
		List<DeliveryService> httpsDs = sep.getSSLEnabledChangeEvents();
		assertThat("Expected to find 3 changed https delivery services but found " + httpsDs.size(),
				httpsDs.size() == 3);
		assertThat("Did not get the expected list of Https Delivery Services " + httpsDs.toString(),
				httpsDs.toString().contains("http-only-test"));
	}

	@Test
	public void getSSLEnabledChangeEvents_new() throws Exception {
		final SnapshotEventsProcessor sep = SnapshotEventsProcessor.diffCrConfigs(newDsSnapJo, updateJo);
		List<DeliveryService> httpsDs = sep.getSSLEnabledChangeEvents();
		// Tests JSON equivalence operator for change in order of members vs. actual data value changes
		// acceptHttps and acceptHttp are swapped in one service and dispersion:limit is changed to 2 in the other
		assertThat("Expected to find 4 changed delivery services but found " + httpsDs.size(), httpsDs.size() == 4);
		assertThat("Did not get the expected list of Https Delivery Services " + httpsDs.toString(),
				httpsDs.toString().contains("http-addnew-test"));
	}

	@Test
	public void getChangeEvents() throws Exception {
		final int changeCnt = 15;
		final SnapshotEventsProcessor sep = SnapshotEventsProcessor.diffCrConfigs(updateJo, baselineJo);
		Map<String, DeliveryService> changes = sep.getChangeEvents();
		assertThat("Expected to find "+changeCnt+" changed delivery services but found " + changes.size(),
				changes.size() == changeCnt);
		assertThat("Did not get the expected list of Changed Delivery Services " + changes.toString(),
				changes.toString().contains("http-only-test"));
		assertThat("Did not get the expected list of Changed Delivery Services " + changes.toString(),
				changes.toString().contains("https-only-test"));
		assertThat("Did not get the expected list of Changed Delivery Services " + changes.toString(),
				changes.toString().contains("http-to-https-test"));
	}

	@Test
	public void getMappingEvents_update() throws Exception {
		final SnapshotEventsProcessor sep = SnapshotEventsProcessor.diffCrConfigs(updateJo, baselineJo);
		List<String> mappingChanges = sep.getMappingEvents();
		assertThat("Expected to find 1 mapping changes but found " + mappingChanges.size(),
				mappingChanges.size() == 1);
		assertThat("Did not get the expected list of mapping changes " + mappingChanges.toString(),
				mappingChanges.toString().contains("edge-cache-011"));
	}

	@Test
	public void getMappingEvents_new() throws Exception {
		final SnapshotEventsProcessor sep = SnapshotEventsProcessor.diffCrConfigs(newDsSnapJo, updateJo);
		List<String> mappingChanges = sep.getMappingEvents();
		assertThat("Expected to find 4 mapping changes but found " + mappingChanges.size(),
				mappingChanges.size() == 4);
		assertThat("Did not get the expected list of mapping changes " + mappingChanges.toString(),
				mappingChanges.contains("edge-cache-000"));
		assertThat("Did not get the expected list of mapping changes " + mappingChanges.toString(),
				mappingChanges.contains("edge-cache-001"));
		assertThat("Did not get the expected list of mapping changes " + mappingChanges.toString(),
				mappingChanges.contains("edge-cache-002"));
		assertThat("Did not get the expected list of mapping changes " + mappingChanges.toString(),
				mappingChanges.contains("edge-cache-011"));
	}

	@Test
	public void getMappingEvents_delete() throws Exception {
		final SnapshotEventsProcessor sep = SnapshotEventsProcessor.diffCrConfigs(updateJo, newDsSnapJo);
		List<String> mappingChanges = sep.getMappingEvents();
		assertThat("Expected to find 4 mapping changes but found " + mappingChanges.size(),
				mappingChanges.size() == 4);
		assertThat("Did not get the expected list of mapping changes " + mappingChanges.toString(),
				mappingChanges.contains("edge-cache-000"));
		assertThat("Did not get the expected list of mapping changes " + mappingChanges.toString(),
				mappingChanges.contains("edge-cache-001"));
		assertThat("Did not get the expected list of mapping changes " + mappingChanges.toString(),
				mappingChanges.contains("edge-cache-002"));
		assertThat("Did not get the expected list of mapping changes " + mappingChanges.toString(),
				mappingChanges.contains("edge-cache-011"));
	}
}