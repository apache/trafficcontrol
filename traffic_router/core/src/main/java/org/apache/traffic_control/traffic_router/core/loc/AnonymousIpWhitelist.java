/*
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

package org.apache.traffic_control.traffic_router.core.loc;

import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.apache.traffic_control.traffic_router.core.util.JsonUtilsException;
import com.fasterxml.jackson.databind.JsonNode;

import java.util.Objects;

public class AnonymousIpWhitelist {
	private static final Logger LOGGER = LogManager.getLogger(AnonymousIpWhitelist.class);

	final private NetworkNode.SuperNode whitelist;

	public AnonymousIpWhitelist() throws NetworkNodeException {
		whitelist = new NetworkNode.SuperNode();
	}

	public void init(final JsonNode config) throws JsonUtilsException, NetworkNodeException {
		if (config.isArray()) {
			for (final JsonNode node : config) {
				final String network = node.asText();
				this.add(network);
			}
		}
	}

	public void add(final String network) throws NetworkNodeException {
		final NetworkNode node = new NetworkNode(network, AnonymousIp.WHITE_LIST_LOC);
		if (network.indexOf(':') == -1) {
			whitelist.add(node);
		} else {
			whitelist.add6(node);
		}
	}

	public boolean contains(final String address) {
		if (whitelist == null) {
			return false;
		}

		try {
			final NetworkNode nn = whitelist.getNetwork(address);
			if (Objects.equals(nn.getLoc(), AnonymousIp.WHITE_LIST_LOC)) {
				return true;
			}
		} catch (NetworkNodeException e) {
			LOGGER.warn("AnonymousIp: exception", e);
		}

		return false;
	}
}
