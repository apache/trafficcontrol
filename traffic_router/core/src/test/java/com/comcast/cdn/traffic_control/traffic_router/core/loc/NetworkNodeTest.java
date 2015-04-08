/*
 * Copyright 2015 Comcast Cable Communications Management, LLC
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

package com.comcast.cdn.traffic_control.traffic_router.core.loc;

import static org.junit.Assert.*;

import java.io.File;

import org.apache.log4j.Logger;
import org.junit.Before;
import org.junit.Test;

import com.comcast.cdn.traffic_control.traffic_router.core.loc.NetworkNode;
import com.comcast.cdn.traffic_control.traffic_router.core.loc.NetworkUpdater;

public class NetworkNodeTest {
	private static final Logger LOGGER = Logger.getLogger(NetworkUpdater.class);

	NetworkNode root;

	@Before
	public void setUp() throws Exception {
		//		final File file = new File(getClass().getClassLoader().getResource("comcast_ipcdn_czf.xml").toURI());
		//		final DocumentBuilderFactory dbf = DocumentBuilderFactory.newInstance();
		//		final DocumentBuilder db = dbf.newDocumentBuilder();
		//		root = NetworkNode.generateTree(db.parse(new FileInputStream(file)));
		final File file = new File(getClass().getClassLoader().getResource("czf2.json").toURI());
		root = NetworkNode.generateTree(file);
	}

	@Test
	public void testIps() {
		try {
			final String testips[][] = {
					{"192.168.8.5", "cache-group-01"},
					{"192.168.9.10", "cache-group-01"},
					{"1234:5678::2", "cache-group-01"},
					{"1234:5679::3", "cache-group-01"},
			};
			for(int i = 0; i < testips.length; i++) {
				final NetworkNode nn = root.getNetwork(testips[i][0]);
				assertNotNull(nn);
				final String loc = nn.getLoc();
				assertEquals(loc, testips[i][1]);
				LOGGER.info(String.format("result for ip=%s: %s",testips[i], loc));
			}
		} catch (Exception e) {
			e.printStackTrace();
		}
	}
}
