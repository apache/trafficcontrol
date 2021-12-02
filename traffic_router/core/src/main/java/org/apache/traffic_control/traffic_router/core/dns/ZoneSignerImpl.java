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

package org.apache.traffic_control.traffic_router.core.dns;

import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.xbill.DNS.DNSKEYRecord;
import org.xbill.DNS.DNSSEC;
import org.xbill.DNS.DSRecord;
import org.xbill.DNS.NSECRecord;
import org.xbill.DNS.Name;
import org.xbill.DNS.RRSIGRecord;
import org.xbill.DNS.RRset;
import org.xbill.DNS.Record;
import org.xbill.DNS.SOARecord;
import org.xbill.DNS.Type;

import java.io.IOException;
import java.security.GeneralSecurityException;
import java.security.PrivateKey;
import java.util.ArrayList;
import java.util.Collections;
import java.util.Date;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Optional;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.ConcurrentMap;
import java.util.stream.Collectors;
import java.util.stream.Stream;
import java.util.stream.StreamSupport;

import static java.util.stream.Collectors.toList;
import static org.xbill.DNS.DClass.IN;

public class ZoneSignerImpl implements ZoneSigner {
	private final static Logger LOGGER = LogManager.getLogger(ZoneSignerImpl.class);

	private Stream<Record> toRRStream(final RRset rrSet) {
		final Iterable<Record> iterable = () -> rrSet.rrs(false);
		return StreamSupport.stream(iterable.spliterator(), false);
	}

	private Stream<Record> toRRSigStream(final RRset rrSset) {
		final Iterable<Record> iterable = rrSset::sigs;
		return StreamSupport.stream(iterable.spliterator(), false);
	}

	private RRSIGRecord sign(final RRset rrset, final DNSKEYRecord dnskeyRecord, final PrivateKey privateKey, final Date inception, final Date expiration) {
		try {
			return DNSSEC.sign(rrset, dnskeyRecord, privateKey, inception, expiration);
		} catch (DNSSEC.DNSSECException e) {
			final String message = String.format("Failed to sign Resource Record Set for %s %d %d %d : %s",
					dnskeyRecord.getName(), dnskeyRecord.getDClass(), dnskeyRecord.getType(), dnskeyRecord.getTTL(), e.getMessage());
			LOGGER.error(message, e);
			return null;
		}
	}

	private boolean isSignatureAlmostExpired(final Date inception, final Date expiration, final Date now) {
		// now is over halfway through validity period
		return now.getTime() > inception.getTime() + ((expiration.getTime() - inception.getTime())/2);
	}

	private RRset signRRset(final RRset rrSet, final List<DnsSecKeyPair> kskPairs, final List<DnsSecKeyPair> zskPairs,
							final Date inception, final Date expiration,
							final ConcurrentMap<RRSIGCacheKey, ConcurrentMap<RRsetKey, RRSIGRecord>> RRSIGCache) {
		final List<RRSIGRecord> signatures = new ArrayList<>();
		final List<DnsSecKeyPair> pairs = rrSet.getType() == Type.DNSKEY ? kskPairs : zskPairs;
		final Date now = new Date();

		pairs.forEach(pair -> {
			final DNSKEYRecord dnskeyRecord = pair.getDNSKEYRecord();
			final PrivateKey privateKey = pair.getPrivate();
			RRSIGRecord signature = null;
			try {
				if (RRSIGCache == null) {
					signature = sign(rrSet, dnskeyRecord, privateKey, inception, expiration);
				} else {
					final ConcurrentMap<RRsetKey, RRSIGRecord> sigMap = RRSIGCache.computeIfAbsent(new RRSIGCacheKey(privateKey.getEncoded(), dnskeyRecord.getAlgorithm()), rrsigCacheKey -> new ConcurrentHashMap<>());
					signature = sigMap.computeIfAbsent(new RRsetKey(rrSet), k -> sign(rrSet, dnskeyRecord, privateKey, inception, expiration));

					if (signature != null && isSignatureAlmostExpired(signature.getTimeSigned(), signature.getExpire(), now)) {
						signature = sigMap.compute(new RRsetKey(rrSet), (k, v) -> sign(rrSet, dnskeyRecord, privateKey, inception, expiration));
					}
				}
			} catch (Exception e) {
				final String message = String.format("Failed to sign Resource Record Set for %s %d %d %d : %s",
					dnskeyRecord.getName(), dnskeyRecord.getDClass(), dnskeyRecord.getType(), dnskeyRecord.getTTL(), e.getMessage());

				LOGGER.error(message, e);
			}
			if (signature != null) {
				signatures.add(signature);
			}
		});

		final RRset signedRRset = new RRset();

		toRRStream(rrSet).forEach(signedRRset::addRR);
		signatures.forEach(signedRRset::addRR);

		return signedRRset;
	}

	private SOARecord findSoaRecord(final List<Record> records) {
		final Optional<Record> soaRecordOptional = records.stream().filter(record -> record instanceof SOARecord).findFirst();
		if (soaRecordOptional.isPresent()) {
			return (SOARecord) soaRecordOptional.get();
		}
		return null;
	}

	private List<NSECRecord> createNsecRecords(final List<Record> records) {
		final Map<Name, List<Record>> recordMap = records.stream().collect(Collectors.groupingBy(Record::getName));
		final List<Name> names = recordMap.keySet().stream().sorted().collect(toList());

		final Map<Name, Name> nextNameTuples = new HashMap<>();

		for (int i = 0; i < names.size(); i++) {
			final Name k = names.get(i);
			final Name v = names.get((i + 1) % names.size());
			nextNameTuples.put(k,v);
		}

		final SOARecord soaRecord = findSoaRecord(records);
		if (soaRecord == null) {
			LOGGER.warn("No SOA record found, this extremely likely to produce DNSSEC errors");
		}

		final long minimumSoaTtl = soaRecord != null ? soaRecord.getMinimum() : 0L;

		final List<NSECRecord> nsecRecords = new ArrayList<>();
		names.forEach(name -> {
			final int[] mostTypes = recordMap.get(name).stream().mapToInt(Record::getType).toArray();
			final int[] allTypes = new int[mostTypes.length + 2];
			System.arraycopy(mostTypes, 0, allTypes, 0, mostTypes.length);
			allTypes[mostTypes.length] = Type.NSEC;
			allTypes[mostTypes.length + 1] = Type.RRSIG;
			nsecRecords.add(new NSECRecord(name, IN, minimumSoaTtl, nextNameTuples.get(name), allTypes));
		});

		return nsecRecords;
	}


	@Override
	public List<Record> signZone(final List<Record> records, final List<DnsSecKeyPair> kskPairs,
								 final List<DnsSecKeyPair> zskPairs, final Date inception, final Date expiration,
								 final ConcurrentMap<RRSIGCacheKey, ConcurrentMap<RRsetKey, RRSIGRecord>> RRSIGCache) throws IOException, GeneralSecurityException {

		final List<NSECRecord> nsecRecords = createNsecRecords(records);
		records.addAll(nsecRecords);

		Collections.sort(records, (record1, record2) -> {
			if (record1.getType() != Type.SOA && record2.getType() != Type.SOA) {
				return record1.compareTo(record2);
			}

			int x = record1.getName().compareTo(record2.getName());

			if (x != 0) {
				return x;
			}

			x = record1.getDClass() - record2.getDClass();

			if (x != 0) {
				return x;
			}

			if (record1.getType() != record2.getType()) {
				return record1.getType() == Type.SOA ? -1 : 1;
			}

			return record1.compareTo(record2);
		});

		final List<RRset> rrSets = new RRSetsBuilder().build(records);

		final List<RRset> signedRrSets = rrSets.stream()
			.map(rRset -> signRRset(rRset, kskPairs, zskPairs, inception, expiration, RRSIGCache))
			.sorted((rRset1, rRset2) -> rRset1.getName().compareTo(rRset2.getName()))
			.collect(toList());

		final List<Record> signedZoneRecords = new ArrayList<>();

		signedRrSets.forEach(rrSet -> {
			signedZoneRecords.addAll(toRRStream(rrSet).collect(toList()));
			signedZoneRecords.addAll(toRRSigStream(rrSet).collect(toList()));
		});

		return signedZoneRecords;
	}

	@Override
	public DSRecord calculateDSRecord(final DNSKEYRecord dnskeyRecord, final int digestId, final long ttl) {
		LOGGER.info("Calculating DS Records for " + dnskeyRecord.getName());
		return new DSRecord(dnskeyRecord.getName(), IN, ttl, digestId, dnskeyRecord);
	}

}
