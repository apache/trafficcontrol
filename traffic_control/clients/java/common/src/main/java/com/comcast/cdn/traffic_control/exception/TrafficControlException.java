package com.comcast.cdn.traffic_control.exception;

public class TrafficControlException extends Exception {
	private static final long serialVersionUID = 914940907727369814L;

	public TrafficControlException() {
		super();
	}

	public TrafficControlException(String message, Throwable cause, boolean enableSuppression, boolean writableStackTrace) {
		super(message, cause, enableSuppression, writableStackTrace);
	}

	public TrafficControlException(String message, Throwable cause) {
		super(message, cause);
	}

	public TrafficControlException(String message) {
		super(message);
	}

	public TrafficControlException(Throwable cause) {
		super(cause);
	}
	
}
