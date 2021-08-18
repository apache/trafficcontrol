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

package org.apache.traffic_control.traffic_router.core.dns.keys;

import org.apache.traffic_control.traffic_router.core.dns.RRSetsBuilder;
import org.apache.traffic_control.traffic_router.shared.ZoneTestRecords;
import org.junit.Before;
import org.junit.Test;
import org.xbill.DNS.RRset;
import org.xbill.DNS.Type;

import java.util.List;
import java.util.Optional;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.equalTo;
import static org.hamcrest.Matchers.notNullValue;

public class RRSetsBuilderTest {
	@Before
	public void before() throws Exception {
		ZoneTestRecords.generateZoneRecords(false);
	}

	RRset findRRSet(List<RRset> rRsets, String name, int type) {
		Optional<RRset> option = rRsets.stream()
			.filter(rRset -> name.equals(rRset.getName().toString()) && rRset.getType() == type)
			.findFirst();

		return option.isPresent() ? option.get() : null;
	}

	@Test
	public void itGroupsResourceRecordsAccordingToRfc4034() throws Exception {
		List<RRset> rRsets = new RRSetsBuilder().build(ZoneTestRecords.records);
		assertThat(rRsets.size(), equalTo(9));
		assertThat(findRRSet(rRsets, "mirror.www.example.com.", Type.CNAME), notNullValue());
		assertThat(findRRSet(rRsets, "ftp.example.com.", Type.AAAA), notNullValue());
		assertThat(findRRSet(rRsets, "ftp.example.com.", Type.A), notNullValue());
		assertThat(findRRSet(rRsets, "www.example.com.", Type.A), notNullValue());
		assertThat(findRRSet(rRsets, "www.example.com.", Type.TXT), notNullValue());
		assertThat(findRRSet(rRsets, "example.com.", Type.NS), notNullValue());
		assertThat(findRRSet(rRsets, "mirror.ftp.example.com.", Type.CNAME), notNullValue());
		assertThat(findRRSet(rRsets, "www.example.com.", Type.AAAA), notNullValue());
		assertThat(findRRSet(rRsets, "example.com.", Type.SOA), notNullValue());
	}
}
