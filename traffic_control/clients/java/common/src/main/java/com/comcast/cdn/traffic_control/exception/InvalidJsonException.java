package com.comcast.cdn.traffic_control.exception;

public class InvalidJsonException extends TrafficControlException {
	private static final long serialVersionUID = 1884362711438565843L;

	public InvalidJsonException() {
		super();
	}

	public InvalidJsonException(String message, Throwable cause, boolean enableSuppression, boolean writableStackTrace) {
		super(message, cause, enableSuppression, writableStackTrace);
	}

	public InvalidJsonException(String message, Throwable cause) {
		super(message, cause);
	}

	public InvalidJsonException(String message) {
		super(message);
	}

	public InvalidJsonException(Throwable cause) {
		super(cause);
	}
	
}
