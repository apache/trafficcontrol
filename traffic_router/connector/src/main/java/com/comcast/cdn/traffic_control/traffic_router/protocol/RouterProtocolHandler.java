package com.comcast.cdn.traffic_control.traffic_router.protocol;

import org.apache.coyote.ProtocolHandler;

public interface RouterProtocolHandler extends ProtocolHandler {
	boolean isReady();

	void setReady(final boolean isReady);

	boolean isInitialized();

	void setInitialized(final boolean isInitialized);

	String getMbeanPath();

	void setMbeanPath(final String mbeanPath);

	String getReadyAttribute();

	void setReadyAttribute(final String readyAttribute);

	String getPortAttribute();

	void setPortAttribute(final String portAttribute);

	int getPort();

	void setPort(int port);
}
