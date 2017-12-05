package com.comcast.cdn.traffic_control.exception;

public class OperationException extends TrafficControlException {
	private static final long serialVersionUID = 8799467021892976240L;

	public OperationException() {
		super();
	}

	public OperationException(String message, Throwable cause, boolean enableSuppression, boolean writableStackTrace) {
		super(message, cause, enableSuppression, writableStackTrace);
	}

	public OperationException(String message, Throwable cause) {
		super(message, cause);
	}

	public OperationException(String message) {
		super(message);
	}

	public OperationException(Throwable cause) {
		super(cause);
	}
	
}
