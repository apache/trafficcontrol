package com.comcast.cdn.traffic_control.traffic_router.core.secure;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

@JsonIgnoreProperties(ignoreUnknown = true)
public class Certificate {
	@JsonProperty
	private String crt;

	@JsonProperty
	private String key;

	public String getCrt() {
		return crt;
	}

	public void setCrt(final String crt) {
		this.crt = crt;
	}

	public String getKey() {
		return key;
	}

	public void setKey(final String key) {
		this.key = key;
	}

	@Override
	@SuppressWarnings("PMD")
	public boolean equals(final Object o) {
		if (this == o) return true;
		if (o == null || getClass() != o.getClass()) return false;

		final Certificate that = (Certificate) o;

		if (crt != null ? !crt.equals(that.crt) : that.crt != null) return false;
		return key != null ? key.equals(that.key) : that.key == null;
	}

	@Override
	public int hashCode() {
		int result = crt != null ? crt.hashCode() : 0;
		result = 31 * result + (key != null ? key.hashCode() : 0);
		return result;
	}
}
