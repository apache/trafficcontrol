package com.comcast.cdn.traffic_control.traffic_router.core.secure;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

@JsonIgnoreProperties(ignoreUnknown = true)
public class CertificateData {
	@JsonProperty
	private String deliveryservice;

	@JsonProperty
	private Certificate certificate;

	@JsonProperty
	private String hostname;

	public String getDeliveryservice() {
		return deliveryservice;
	}

	public void setDeliveryservice(final String deliveryservice) {
		this.deliveryservice = deliveryservice;
	}

	public Certificate getCertificate() {
		return certificate;
	}

	public void setCertificate(final Certificate certificate) {
		this.certificate = certificate;
	}

	public String getHostname() {
		return hostname;
	}

	public void setHostname(final String hostname) {
		this.hostname = hostname;
	}

	@SuppressWarnings("PMD.IfStmtsMustUseBraces")
	@Override
	public boolean equals(final Object o) {
		if (this == o) return true;
		if (o == null || getClass() != o.getClass()) return false;

		final CertificateData that = (CertificateData) o;

		if (deliveryservice != null ? !deliveryservice.equals(that.deliveryservice) : that.deliveryservice != null)
			return false;
		if (certificate != null ? !certificate.equals(that.certificate) : that.certificate != null) return false;
		return hostname != null ? hostname.equals(that.hostname) : that.hostname == null;

	}

	@Override
	public int hashCode() {
		int result = deliveryservice != null ? deliveryservice.hashCode() : 0;
		result = 31 * result + (certificate != null ? certificate.hashCode() : 0);
		result = 31 * result + (hostname != null ? hostname.hashCode() : 0);
		return result;
	}
}
