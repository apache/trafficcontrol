package com.comcast.cdn.traffic_control.traffic_router.core.util;

public class JsonUtilsException extends Exception {

    public JsonUtilsException(final String reason) {
        super(reason);
    }

    public JsonUtilsException(final String message, final Throwable cause) {
        super(message, cause);
    }
}
