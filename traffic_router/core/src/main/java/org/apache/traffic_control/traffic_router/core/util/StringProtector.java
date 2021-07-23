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

package org.apache.traffic_control.traffic_router.core.util;

import java.io.IOException;
import java.io.UnsupportedEncodingException;
import java.security.GeneralSecurityException;

import javax.crypto.Cipher;
import javax.crypto.SecretKey;
import javax.crypto.SecretKeyFactory;
import javax.crypto.spec.PBEKeySpec;
import javax.crypto.spec.PBEParameterSpec;

import org.apache.commons.codec.binary.Base64;

public class StringProtector {

	final private Base64 base64 = new Base64(true);
	private final Cipher encryptor;
	private final Cipher decryptor;
	private static final byte[] SALT = {
		(byte) 0xde, (byte) 0x33, (byte) 0x10, (byte) 0x12,
		(byte) 0xde, (byte) 0x33, (byte) 0x10, (byte) 0x12,
	};

	public StringProtector(final String passwd) throws GeneralSecurityException {
		final SecretKeyFactory keyFactory = SecretKeyFactory.getInstance("PBEWithMD5AndDES");
		final SecretKey key = keyFactory.generateSecret(new PBEKeySpec(passwd.toCharArray()));
		encryptor = Cipher.getInstance("PBEWithMD5AndDES");
		encryptor.init(Cipher.ENCRYPT_MODE, key, new PBEParameterSpec(SALT, 20));

		decryptor = Cipher.getInstance("PBEWithMD5AndDES");
		decryptor.init(Cipher.DECRYPT_MODE, key, new PBEParameterSpec(SALT, 20));
	}

	public byte[] encrypt(final byte[] property) throws GeneralSecurityException, UnsupportedEncodingException {
		return encryptor.doFinal(property);
	}
	public String encrypt(final String property) throws UnsupportedEncodingException, GeneralSecurityException {
		return base64.encodeAsString(encrypt(property.getBytes("UTF-8")));
	}
	public String encryptForUrl(final byte[] data) throws UnsupportedEncodingException, GeneralSecurityException {
		return base64.encodeAsString(encrypt(data));
	}
	public String encodeForUrl(final byte[] data) throws UnsupportedEncodingException, GeneralSecurityException {
		return base64.encodeAsString(data);
	}

	public byte[] decrypt(final byte[] property) throws GeneralSecurityException, IOException {
		return decryptor.doFinal(property);
	}
	public String decrypt(final String property) throws GeneralSecurityException, IOException {
		final byte[] bytes = decrypt(base64.decode(property));
		return new String(bytes, "UTF-8");
	}

	//	public static void main(final String[] args) throws Exception {
	//		StringProtector sp = new StringProtector("my passwd");
	//		String originalPassword = "secret";
	////		System.out.println("Original password: " + originalPassword);
	//		String encryptedPassword = sp.encrypt(originalPassword);
	////		System.out.println("Encrypted password: " + encryptedPassword);
	//		String decryptedPassword = sp.decrypt(encryptedPassword);
	////		System.out.println("Decrypted password: " + decryptedPassword);
	//	}
}
