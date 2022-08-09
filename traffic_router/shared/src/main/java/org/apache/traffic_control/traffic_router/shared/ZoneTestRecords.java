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

package org.apache.traffic_control.traffic_router.shared;

import org.apache.traffic_control.traffic_router.secure.BindPrivateKey;
import org.apache.traffic_control.traffic_router.secure.Pkcs1KeySpecDecoder;
import org.xbill.DNS.*;

import java.io.IOException;
import java.net.Inet6Address;
import java.net.InetAddress;
import java.security.GeneralSecurityException;
import java.security.KeyFactory;
import java.security.KeyPair;
import java.security.KeyPairGenerator;
import java.security.NoSuchAlgorithmException;
import java.security.NoSuchProviderException;
import java.security.PrivateKey;
import java.security.PublicKey;
import java.security.SecureRandom;
import java.time.Duration;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Date;
import java.util.List;

import static java.util.Base64.getMimeDecoder;

@SuppressWarnings("PMD.ClassNamingConventions")
public class ZoneTestRecords {
	public static List<Record> records;

	public static Date start;
	public static Date expiration;
	public static Name origin;
	public static Date sep_1_2016 = new Date(1472688000000L);
	public static Date sep_1_2026 = new Date(1788220800000L);
	public static DNSKEYRecord zoneSigningKeyRecord;
	public static DNSKEYRecord keySigningKeyRecord;

	public static KeyPair ksk1;
	public static KeyPair zsk1;
	public static KeyPair ksk2;
	public static KeyPair zsk2;

	static List<KeyPair> generateKeyPairs() throws NoSuchAlgorithmException, NoSuchProviderException {
		final KeyPairGenerator keyPairGenerator = KeyPairGenerator.getInstance("RSA");
		keyPairGenerator.initialize(2048, SecureRandom.getInstance("SHA1PRNG","SUN"));
		final List<KeyPair> keyPairs = new ArrayList<>();
		keyPairs.add(keyPairGenerator.generateKeyPair());
		keyPairs.add(keyPairGenerator.generateKeyPair());
		keyPairs.add(keyPairGenerator.generateKeyPair());
		keyPairs.add(keyPairGenerator.generateKeyPair());
		return keyPairs;
	}

	private static KeyPair recreateKeyPair(final String publicKey, final String privateKey) throws GeneralSecurityException, IOException {
		final PrivateKey privateKeyCopy = new BindPrivateKey().decode(new String(getMimeDecoder().decode(privateKey)));
		final PublicKey publicKeyCopy = KeyFactory.getInstance("RSA").generatePublic(new Pkcs1KeySpecDecoder().decode(publicKey));
		return new KeyPair(publicKeyCopy, privateKeyCopy);
	}

	public static List<Record> generateZoneRecords(final boolean makeNewKeyPairs) throws IOException, GeneralSecurityException {
		start = new Date(System.currentTimeMillis() - (24 * 3600 * 1000));
		expiration = new Date(System.currentTimeMillis() + (7 * 24 * 3600 * 1000));

		origin = new Name("example.com.");

		final Duration tenYears = Duration.ofDays(3650);
		final Duration oneDay = Duration.ofDays(1);
		final Duration threeDays = Duration.ofDays(3);
		final Duration threeWeeks = Duration.ofDays(21);

		final long oneHour = 3600;
		final Name nameServer1 = new Name("ns1.example.com.");
		final Name nameServer2 = new Name("ns2.example.com.");

		final Name adminEmail = new Name("admin.example.com.");

		final Name webServer = new Name("www.example.com.");
		final Name ftpServer = new Name("ftp.example.com.");

		final Name webMirror = new Name("mirror.www.example.com.");
		final Name ftpMirror = new Name("mirror.ftp.example.com.");

		final String txtRecord = "dead0123456789";

		records = new ArrayList<>(Arrays.asList(
			new AAAARecord(webServer, DClass.IN, threeDays.getSeconds(), Inet6Address.getByAddress(new byte[]{32, 1, 13, -72, 0, 0, 0, 0, 0, 5, 0, 6, 0, 7, 0, 8})), // 2001:db8::5:6:7:8
			new AAAARecord(ftpServer, DClass.IN, threeDays.getSeconds(), Inet6Address.getByAddress(new byte[]{32, 1, 13, -72, 0, 0, 0, 0, 0, 18, 0, 52, 0, 86, 0, 120})), // 2001:db8::12:34:56:78
			new NSRecord(origin, DClass.IN, tenYears.getSeconds(), nameServer1),
			new NSRecord(origin, DClass.IN, tenYears.getSeconds(), nameServer2),
			new ARecord(webServer, DClass.IN, threeWeeks.getSeconds(), InetAddress.getByAddress(new byte[] {11, 22, 33, 44})),
			new ARecord(webServer, DClass.IN, threeWeeks.getSeconds(), InetAddress.getByAddress(new byte[] {55, 66, 77, 88})),
			new ARecord(ftpServer, DClass.IN, threeWeeks.getSeconds(), InetAddress.getByAddress(new byte[] {12, 34, 56, 78})),
			new ARecord(ftpServer, DClass.IN, threeWeeks.getSeconds(), InetAddress.getByAddress(new byte[] {21, 43, 65, 87})),
			new AAAARecord(webServer, DClass.IN, threeDays.getSeconds(), Inet6Address.getByAddress(new byte[]{32, 1, 13, -72, 0, 0, 0, 0, 0, 4, 0, 3, 0, 2, 0, 1})), // 2001:db8::4:3:2:1
			new SOARecord(origin, DClass.IN, tenYears.getSeconds(), nameServer1,
				adminEmail, 2016091400L, oneDay.getSeconds(), oneHour, threeWeeks.getSeconds(), threeDays.getSeconds()),
			new AAAARecord(ftpServer, DClass.IN, threeDays.getSeconds(), Inet6Address.getByAddress(new byte[]{32, 1, 13, -72, 0, 0, 0, 0, 0, 33, 0, 67, 0, 101, 0, -121})), // 2001:db8::21:43:65:87
			new CNAMERecord(webMirror, DClass.IN, tenYears.getSeconds(), webServer),
			new CNAMERecord(ftpMirror, DClass.IN, tenYears.getSeconds(), ftpServer),
			new TXTRecord(webServer, DClass.IN, tenYears.getSeconds(), txtRecord)
		));

		if (makeNewKeyPairs) {
			final List<KeyPair> keyPairs = generateKeyPairs();
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
				DNSKEYRecord.Flags.ZONE_KEY, DNSKEYRecord.Protocol.DNSSEC, DNSSEC.Algorithm.RSASHA1, zsk1.getPublic().getEncoded());

		keySigningKeyRecord = new DNSKEYRecord(origin, DClass.IN, 315569520L,
				DNSKEYRecord.Flags.ZONE_KEY | DNSKEYRecord.Flags.SEP_KEY, DNSKEYRecord.Protocol.DNSSEC, DNSSEC.Algorithm.RSASHA1, ksk1.getPublic().getEncoded());
		return records;
	}
}
