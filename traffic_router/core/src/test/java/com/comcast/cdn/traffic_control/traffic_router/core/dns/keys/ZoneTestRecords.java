package com.comcast.cdn.traffic_control.traffic_router.core.dns.keys;

import com.comcast.cdn.traffic_control.traffic_router.secure.Pkcs1;
import org.xbill.DNS.AAAARecord;
import org.xbill.DNS.ARecord;
import org.xbill.DNS.CNAMERecord;
import org.xbill.DNS.DClass;
import org.xbill.DNS.DNSKEYRecord;
import org.xbill.DNS.NSRecord;
import org.xbill.DNS.Name;
import org.xbill.DNS.Record;
import org.xbill.DNS.SOARecord;

import java.net.Inet6Address;
import java.net.InetAddress;
import java.security.KeyPair;
import java.security.KeyPairGenerator;
import java.security.PrivateKey;
import java.security.PublicKey;
import java.security.SecureRandom;
import java.time.Duration;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Date;
import java.util.List;

import static org.xbill.DNS.DNSKEYRecord.Flags.SEP_KEY;
import static org.xbill.DNS.DNSKEYRecord.Flags.ZONE_KEY;
import static org.xbill.DNS.DNSKEYRecord.Protocol.DNSSEC;
import static org.xbill.DNS.DNSSEC.Algorithm.RSASHA1;

public class ZoneTestRecords {
	static List<Record> records;

	static Date start;
	static Date expiration;
	static Name origin;
	static Date sep_1_2016 = new Date(1472688000000L);
	static Date sep_1_2026 = new Date(1788220800000L);
	static DNSKEYRecord zoneSigningKeyRecord;
	static DNSKEYRecord keySigningKeyRecord;

	static KeyPair ksk1;
	static KeyPair zsk1;
	static KeyPair ksk2;
	static KeyPair zsk2;

	static List<KeyPair> generateKeyPairs() throws Exception {
		KeyPairGenerator keyPairGenerator = KeyPairGenerator.getInstance("RSA");
		keyPairGenerator.initialize(2048, SecureRandom.getInstance("SHA1PRNG","SUN"));
		List<KeyPair> keyPairs = new ArrayList<>();
		keyPairs.add(keyPairGenerator.generateKeyPair());
		keyPairs.add(keyPairGenerator.generateKeyPair());
		keyPairs.add(keyPairGenerator.generateKeyPair());
		keyPairs.add(keyPairGenerator.generateKeyPair());
		return keyPairs;
	}

	private static KeyPair recreateKeyPair(String publicKey, String privateKey) throws Exception {
		Pkcs1 pkcs1 = new Pkcs1(privateKey, publicKey);

		PrivateKey privateKeyCopy = pkcs1.getPrivateKey();
		PublicKey publicKeyCopy = pkcs1.getPublicKey();

		return new KeyPair(publicKeyCopy, privateKeyCopy);
	}

	static List<Record> generateZoneRecords(boolean makeNewKeyPairs) throws Exception {
		start = new Date(System.currentTimeMillis() - (24 * 3600 * 1000));
		expiration = new Date(System.currentTimeMillis() + (7 * 24 * 3600 * 1000));

		origin = new Name("example.com.");

		Duration tenYears = Duration.ofDays(3650);
		Duration oneDay = Duration.ofDays(1);
		Duration threeDays = Duration.ofDays(3);
		Duration threeWeeks = Duration.ofDays(21);

		long oneHour = 3600;
		Name nameServer1 = new Name("ns1.example.com.");
		Name nameServer2 = new Name("ns2.example.com.");

		Name adminEmail = new Name("admin.example.com.");

		Name webServer = new Name("www.example.com.");
		Name ftpServer = new Name("ftp.example.com.");

		Name webMirror = new Name("mirror.www.example.com.");
		Name ftpMirror = new Name("mirror.ftp.example.com.");

		records = new ArrayList<>(Arrays.asList(
			new AAAARecord(webServer, DClass.IN, threeDays.getSeconds(), Inet6Address.getByName("2001:db8::5:6:7:8")),
			new AAAARecord(ftpServer, DClass.IN, threeDays.getSeconds(), Inet6Address.getByName("2001:db8::12:34:56:78")),
			new NSRecord(origin, DClass.IN, tenYears.getSeconds(), nameServer1),
			new NSRecord(origin, DClass.IN, tenYears.getSeconds(), nameServer2),
			new ARecord(webServer, DClass.IN, threeWeeks.getSeconds(), InetAddress.getByAddress(new byte[] {11, 22, 33, 44})),
			new ARecord(webServer, DClass.IN, threeWeeks.getSeconds(), InetAddress.getByAddress(new byte[] {55, 66, 77, 88})),
			new ARecord(ftpServer, DClass.IN, threeWeeks.getSeconds(), InetAddress.getByAddress(new byte[] {12, 34, 56, 78})),
			new ARecord(ftpServer, DClass.IN, threeWeeks.getSeconds(), InetAddress.getByAddress(new byte[] {21, 43, 65, 87})),
			new AAAARecord(webServer, DClass.IN, threeDays.getSeconds(), Inet6Address.getByName("2001:db8::4:3:2:1")),
			new SOARecord(origin, DClass.IN, tenYears.getSeconds(), nameServer1,
				adminEmail, 2016091400L, oneDay.getSeconds(), oneHour, threeWeeks.getSeconds(), threeDays.getSeconds()),
			new AAAARecord(ftpServer, DClass.IN, threeDays.getSeconds(), Inet6Address.getByName("2001:db8::21:43:65:87")),
			new CNAMERecord(webMirror, DClass.IN, tenYears.getSeconds(), webServer),
			new CNAMERecord(ftpMirror, DClass.IN, tenYears.getSeconds(), ftpServer)
		));

		if (makeNewKeyPairs) {
			List<KeyPair> keyPairs = generateKeyPairs();
			ksk1 = keyPairs.get(0);
			zsk1 = keyPairs.get(1);
			ksk2 = keyPairs.get(2);
			zsk2 = keyPairs.get(3);
		} else {
			ksk1 = recreateKeyPair(SigningData.ksk1Public, SigningData.ksk1Private);
			zsk1 = recreateKeyPair(SigningData.zsk1Public, SigningData.zsk1Private);
			ksk2 = recreateKeyPair(SigningData.ksk2Public, SigningData.ksk2Private);
			zsk2 = recreateKeyPair(SigningData.zsk2Public, SigningData.zsk2Private);
		}

		zoneSigningKeyRecord = new DNSKEYRecord(origin, DClass.IN, 31556952L,
			ZONE_KEY, DNSSEC, RSASHA1, zsk1.getPublic().getEncoded());

		keySigningKeyRecord = new DNSKEYRecord(origin, DClass.IN, 315569520L,
			ZONE_KEY | SEP_KEY, DNSSEC, RSASHA1, ksk1.getPublic().getEncoded());
		return records;
	}
}
