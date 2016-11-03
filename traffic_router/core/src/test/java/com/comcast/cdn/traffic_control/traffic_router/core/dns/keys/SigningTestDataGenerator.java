package com.comcast.cdn.traffic_control.traffic_router.core.dns.keys;

import com.verisignlabs.dnssec.security.DnsKeyPair;
import com.verisignlabs.dnssec.security.JCEDnsSecSigner;
import org.junit.Before;
import org.junit.Test;
import org.xbill.DNS.DClass;
import org.xbill.DNS.DNSKEYRecord;
import org.xbill.DNS.DSRecord;
import org.xbill.DNS.Name;
import org.xbill.DNS.Record;
import org.xbill.DNS.Section;
import sun.security.rsa.RSAPrivateCrtKeyImpl;

import java.io.IOException;
import java.security.Key;
import java.security.KeyPair;
import java.security.interfaces.RSAPublicKey;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Base64;
import java.util.List;

import static com.comcast.cdn.traffic_control.traffic_router.core.dns.keys.ZoneTestRecords.generateZoneRecords;
import static com.comcast.cdn.traffic_control.traffic_router.core.dns.keys.ZoneTestRecords.keySigningKeyRecord;
import static com.comcast.cdn.traffic_control.traffic_router.core.dns.keys.ZoneTestRecords.ksk1;
import static com.comcast.cdn.traffic_control.traffic_router.core.dns.keys.ZoneTestRecords.ksk2;
import static com.comcast.cdn.traffic_control.traffic_router.core.dns.keys.ZoneTestRecords.zoneSigningKeyRecord;
import static com.comcast.cdn.traffic_control.traffic_router.core.dns.keys.ZoneTestRecords.zsk1;
import static com.comcast.cdn.traffic_control.traffic_router.core.dns.keys.ZoneTestRecords.zsk2;
import static java.util.Base64.getEncoder;
import static java.util.Base64.getMimeEncoder;
import static java.util.stream.Collectors.toList;
import static org.xbill.DNS.DSRecord.SHA256_DIGEST_ID;

public class SigningTestDataGenerator {
	private Base64.Encoder encoder = getMimeEncoder(76, new byte[]{'\n'});

	byte[] encode(byte[] data) {
		return new String(encoder.encode(getEncoder().encode(data))).replaceAll("\n", "\\\\n").getBytes();
	}

	String encodeDnsKeyRecord(DNSKEYRecord dnskeyRecord) {
		return new String(getMimeEncoder(76, new byte[]{'\n'}).encode(dnskeyRecord.toString().getBytes())).replaceAll("\n", "\\\\n");
	}

	void dumpKeyPair(String varPrefix, KeyPair keyPair) throws IOException {
		dumpKey(String.format("%sPublic", varPrefix), keyPair.getPublic());
		dumpKey(String.format("%sPrivate", varPrefix), keyPair.getPrivate());
	}

	void dumpKey(String varName, Key key) throws IOException {

		byte[] base64Encoded;
		if (key instanceof RSAPrivateCrtKeyImpl) {
			String s = new BindPrivateKeyFormatter().format((RSAPrivateCrtKeyImpl) key);
			base64Encoded = new String(encoder.encode(s.getBytes())).replaceAll("\n", "\\\\n").getBytes();
		} else if (key instanceof RSAPublicKey) {
			base64Encoded = getEncoder().encode(new Pkcs1Formatter().toBytes((RSAPublicKey) key));
		} else {
			base64Encoded = encode(encode(key.getEncoded()));
		}

		System.out.println(makeBase64StringVar(varName, new String(base64Encoded)));
	}

	String makeBase64StringVar(String varName, String base64String) {
		int length = 100;
		int beginIndex = 0;
		int endIndex = length;
		StringBuilder stringBuilder = new StringBuilder("static String " + varName + " =\n");
		while (beginIndex < base64String.length()) {
			if (endIndex > base64String.length()) {
				endIndex = base64String.length();
			}
			stringBuilder.append(String.format("\t\"%s\"", base64String.substring(beginIndex, endIndex)));
			beginIndex = endIndex;
			if (beginIndex < base64String.length()) {
				stringBuilder.append(" +");
			}
			stringBuilder.append("\n");
			endIndex += length;
		}
		stringBuilder.append("\t;\n");
		return stringBuilder.toString();
	}

	@Before
	public void before() throws Exception {
		generateZoneRecords(true);
		Name origin = new Name("example.com.");

		dumpKeyPair("ksk1", ksk1);
		System.out.println();

		dumpKeyPair("ksk2", ksk2);
		System.out.println();

		dumpKeyPair("zsk1", zsk1);
		System.out.println();

		dumpKeyPair("zsk2", zsk2);
		System.out.println();

		JCEDnsSecSigner signer = new JCEDnsSecSigner(false);

		List<DnsKeyPair> kskPairs = new ArrayList<>(Arrays.asList(
			new DnsKeyPair(keySigningKeyRecord, new BindPrivateKeyFormatter().format((RSAPrivateCrtKeyImpl) ksk1.getPrivate())),
			new DnsKeyPair(keySigningKeyRecord, new BindPrivateKeyFormatter().format((RSAPrivateCrtKeyImpl) ksk2.getPrivate()))
		));

		List<DnsKeyPair> zskPairs = new ArrayList<>(Arrays.asList(
			new DnsKeyPair(zoneSigningKeyRecord, new BindPrivateKeyFormatter().format((RSAPrivateCrtKeyImpl) zsk1.getPrivate())),
			new DnsKeyPair(zoneSigningKeyRecord, new BindPrivateKeyFormatter().format((RSAPrivateCrtKeyImpl) zsk2.getPrivate()))
		));

		List<Record> signedRecords = signer.signZone(origin, ZoneTestRecords.records, kskPairs, zskPairs,
			ZoneTestRecords.sep_1_2016, ZoneTestRecords.sep_1_2026, true, SHA256_DIGEST_ID);

		ZoneTestRecords.records.forEach(rec -> {
			System.out.println("// " + rec);
			// Doesn't really matter that 'ANSWER' is totally correct, just don't use question
			String base64String = new String(getEncoder().encode(rec.toWire(Section.ANSWER)));
			String varName = String.format("postZoneRecord%d", signedRecords.indexOf(rec));
			System.out.println(makeBase64StringVar(varName, base64String));
		});

		signedRecords.forEach(rec -> {
			System.out.println("// " + rec);
			// Doesn't really matter that 'ANSWER' is totally correct, just don't use question
			String base64String = new String(getEncoder().encode(rec.toWire(Section.ANSWER)));
			String varName = String.format("signedRecord%d", signedRecords.indexOf(rec));
			System.out.println(makeBase64StringVar(varName, base64String));
		});

		List<DSRecord> dsRecords = kskPairs.stream()
			.map(pair -> new DSRecord(origin, DClass.IN, 1234000L, SHA256_DIGEST_ID, pair.getDNSKEYRecord()))
			.collect(toList());

		dsRecords.forEach(rec -> {
			System.out.println("// " + rec);
			String base64String = new String(getEncoder().encode(rec.toWire(Section.ANSWER)));
			String varName = String.format("dsRecord%d", dsRecords.indexOf(rec));
			System.out.println(makeBase64StringVar(varName, base64String));
		});

		System.out.println("// " + zoneSigningKeyRecord);
		System.out.println("// keytag " + zoneSigningKeyRecord.getFootprint());
		System.out.println(makeBase64StringVar("zoneDnsKeyRecord", encodeDnsKeyRecord(zoneSigningKeyRecord)));

		System.out.println("// " + keySigningKeyRecord);
		System.out.println("// keytag " + zoneSigningKeyRecord.getFootprint());
		System.out.println(makeBase64StringVar("keyDnsKeyRecord", encodeDnsKeyRecord(keySigningKeyRecord)));
	}

	@Test
	public void test() {
		System.out.println("ok");
	}
}
