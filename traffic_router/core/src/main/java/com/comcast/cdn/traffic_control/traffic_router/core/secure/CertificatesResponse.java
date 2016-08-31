package com.comcast.cdn.traffic_control.traffic_router.core.secure;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

import java.util.List;

@JsonIgnoreProperties(ignoreUnknown = true)
public class CertificatesResponse {
	@JsonProperty
	private List<CertificateData> response;

	public List<CertificateData> getResponse() {
		return response;
	}

	public void setResponse(final List<CertificateData> response) {
		this.response = response;
	}
}
