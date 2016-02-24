package com.comcast.cdn.traffic_control.traffic_monitor.wicket.models;

import org.apache.wicket.model.Model;

public class CacheDataModel extends Model<String> {
	private static final long serialVersionUID = 1L;
	private final String label;
	long i = 0;

	public CacheDataModel(final String label) {
		this.label = label;

		if (label == null) {
			super.setObject(null);
		} else {
			super.setObject(label + ": ");
		}
	}

	public String getKey() {
		return label;
	}

	public String getValue() {
		return String.valueOf(i);
	}

	public long getRawValue() {
		return i;
	}

	public void inc() {
		synchronized (this) {
			i++;
			this.set(i);
		}
	}

	public void setObject(final String o) {
		if (label == null) {
			super.setObject(o);
		} else {
			super.setObject(label + ": " + o);
		}
	}

	public void set(final long arg) {
		synchronized (this) {
			i = arg;
			this.setObject(String.valueOf(arg));
		}
	}
}
