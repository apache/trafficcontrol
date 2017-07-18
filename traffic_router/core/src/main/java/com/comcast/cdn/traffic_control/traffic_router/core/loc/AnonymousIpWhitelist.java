package com.comcast.cdn.traffic_control.traffic_router.core.loc;

import org.apache.log4j.Logger;
import org.apache.wicket.ajax.json.JSONArray;
import org.apache.wicket.ajax.json.JSONException;

public class AnonymousIpWhitelist {
	private static final Logger LOGGER = Logger.getLogger(AnonymousIpWhitelist.class);

	final private NetworkNode.SuperNode whitelist;

	public AnonymousIpWhitelist() throws NetworkNodeException {
		whitelist = new NetworkNode.SuperNode();
	}

	public void init(final JSONArray config) throws JSONException, NetworkNodeException {
		for (int i = 0; i < config.length(); i++) {
			final String network = config.getString(i);
			this.add(network);
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
			if (nn.getLoc() == AnonymousIp.WHITE_LIST_LOC) {
				return true;
			}
		} catch (NetworkNodeException e) {
			LOGGER.warn("AnonymousIp: exception", e);
		}

		return false;
	}
}
