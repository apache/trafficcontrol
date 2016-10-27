package com.comcast.cdn.traffic_control.traffic_router.core.dns;

import org.xbill.DNS.RRset;
import org.xbill.DNS.Record;
import org.xbill.DNS.Type;

import java.util.Comparator;
import java.util.List;
import java.util.Map;
import java.util.function.Function;
import java.util.stream.Collectors;

public class RRSetsBuilder {
	final private Function<List<Record>, RRset> recordsToRRSet = (Function<List<Record>, RRset>) records -> {
		final RRset rrSet = new RRset();
		records.forEach(rrSet::addRR);
		return rrSet;
	};

	private static String qualifer(final Record record) {
		return String.format("%s %d %d %d", record.getName().toString(), record.getDClass(), record.getType(), record.getTTL());
	}

	final private Comparator<RRset> rrSetComparator = (rrSet1, rrSet2) -> {
		int x = rrSet1.getName().compareTo(rrSet2.getName());

		if (x != 0) {
			return x;
		}

		x = rrSet1.getDClass() - rrSet2.getDClass();
		if (x != 0) {
			return x;
		}

		if (rrSet1.getType() == Type.SOA) {
			return -1;
		}

		return rrSet1.getType() - rrSet2.getType();
	};

	public List<RRset> build(final List<Record> records) {
		final Map<String, List<Record>> map = records.stream().sorted().collect(
			Collectors.groupingBy(RRSetsBuilder::qualifer, Collectors.toList())
		);

		return map.values().stream().map(recordsToRRSet).sorted(rrSetComparator).collect(Collectors.toList());
	}
}
