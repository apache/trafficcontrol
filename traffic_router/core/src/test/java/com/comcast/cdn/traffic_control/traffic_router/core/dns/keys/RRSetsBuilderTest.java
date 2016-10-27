package com.comcast.cdn.traffic_control.traffic_router.core.dns.keys;

import com.comcast.cdn.traffic_control.traffic_router.core.dns.RRSetsBuilder;
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
		assertThat(rRsets.size(), equalTo(8));
		assertThat(findRRSet(rRsets, "mirror.www.example.com.", Type.CNAME), notNullValue());
		assertThat(findRRSet(rRsets, "ftp.example.com.", Type.AAAA), notNullValue());
		assertThat(findRRSet(rRsets, "ftp.example.com.", Type.A), notNullValue());
		assertThat(findRRSet(rRsets, "www.example.com.", Type.A), notNullValue());
		assertThat(findRRSet(rRsets, "example.com.", Type.NS), notNullValue());
		assertThat(findRRSet(rRsets, "mirror.ftp.example.com.", Type.CNAME), notNullValue());
		assertThat(findRRSet(rRsets, "www.example.com.", Type.AAAA), notNullValue());
		assertThat(findRRSet(rRsets, "example.com.", Type.SOA), notNullValue());
	}
}
