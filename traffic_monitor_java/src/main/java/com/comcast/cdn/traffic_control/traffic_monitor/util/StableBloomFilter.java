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

package com.comcast.cdn.traffic_control.traffic_monitor.util;

import java.io.Serializable;
import java.nio.charset.Charset;
import java.security.MessageDigest;
import java.security.NoSuchAlgorithmException;
import java.util.Arrays;
import java.security.SecureRandom;
import java.math.BigInteger;

import org.apache.log4j.Logger;

public class StableBloomFilter implements Serializable {
	private static final Logger LOGGER = Logger.getLogger(StableBloomFilter.class);
	private static final long serialVersionUID = 1L;
	final private byte[] byteSet;
	final private int k; // number of hash functions
	final private int p; // cells to decrement
	final private int m; // cell count
	// int d; number of bits per cell: 8
	final private int max; // value a cell is set to 

	static final Charset charset = Charset.forName("UTF-8"); // encoding used for storing hash values as strings
	static final MessageDigest digestFunction;
	static { // The digest method is reused between instances
		MessageDigest tmp;
		try {
			tmp = MessageDigest.getInstance("SHA1"); // SHA1, MD5
		} catch (NoSuchAlgorithmException e) {
			tmp = null;
		}
		digestFunction = tmp;
	}

	public StableBloomFilter(final int cellCount, final int k, final int p, final int max) {
		this.k = k;
		this.p = p;
		this.m = cellCount;
		this.max = max;
		this.byteSet = new byte[cellCount];
	}

	private static int[] generateHashes(final byte[] data, final int hashes) {
		final int[] result = new int[hashes];

		byte salt = 0;
		for (int k = 0; k < hashes; ) {
			byte[] digest;
			synchronized (digestFunction) {
				digestFunction.update(salt);
				salt++;
				digest = digestFunction.digest(data);                
			}

			for (int i = 0; i < digest.length/4 && k < hashes; i++) {
				int h = 0;
				for (int j = (i*4); j < (i*4)+4; j++) {
					h <<= 8;
					h |= ((int) digest[j]) & 0xFF;
				}
				result[k] = h;
				k++;
			}
		}
		return result;
	}

	public int getK() {
		return k;
	}
	public int getP() {
		return p;
	}

	public void clear() {
		Arrays.fill(byteSet, (byte)0);
	}

	public void add(final String element) {
		add(element.getBytes(charset));
	}
	public boolean add(final byte[] bytes) {
		final boolean hit = contains(bytes);
		for (int i = 0; i < p; i++) {
			final int pi = (int) (Math.random()*m);
			if(byteSet[pi] != 0) {
				byteSet[pi]--;
			}
		}
		for (int hash : generateHashes(bytes, k)) {
			final int index = Math.abs(hash % m);
			byteSet[index] = (byte)this.max;
		}
		return hit;
	}

	public boolean contains(final String element) {
		return contains(element.getBytes(charset));
	}

	public boolean contains(final byte[] bytes) {
		for (int hash : generateHashes(bytes, k)) {
			if (byteSet[Math.abs(hash % m)] == 0) {
				return false;
			}
		}
		return true;
	}
	public double getFPS() {
		final double d11 = 1/(double)k;
		final double d12 = 1/(double)m;
		final double d1 = d11 - d12;
		final double d2 = (p*d1);
		final double d4 = Math.pow((1.0/(1.0 + (1.0/ d2))),max);
		return Math.pow((1 - d4), k);
	}

	public static void main(final String[] args) {
		final SessionIdentifierGenerator generator = new SessionIdentifierGenerator();
		final int elementCount = 800000;
		final StableBloomFilter sbf = new StableBloomFilter(elementCount, 2, 4, 2);
		LOGGER.warn("FPS: "+sbf.getFPS());
		final int strCount = 100000;
		LOGGER.warn("strCount: "+strCount);
		final String[] strList = new String[strCount];
		int fnCnt = 0;
		int fpCnt = 0;
		for(int i = 0; i < strCount; i++) {
			strList[i] = generator.nextSessionId();
//			LOGGER.warn(strList[i]);
			sbf.add(strList[i]);
		}
		LOGGER.warn(" --- ");
		for(int i = 0; i < strCount; i++) {
			final boolean result = sbf.contains(strList[i]);
			if(!result) {
				fnCnt++;
			}
//			LOGGER.warn(strList[i] + " : "+result);
		}
		LOGGER.warn("fnCnt: "+fnCnt);
		for(int i = 0; i < strCount; i++) {
			final String rstr = generator.nextSessionId();
			final boolean result = sbf.contains(rstr);
			if(result) {
				fpCnt++;
			}
//			LOGGER.warn(rstr + " : "+result);
		}
		LOGGER.warn("fpCnt: "+fpCnt);
	}

	static class SessionIdentifierGenerator {
		final private SecureRandom random = new SecureRandom();

		public String nextSessionId() {
			return new BigInteger(130, random).toString(32);
		}

	}
}