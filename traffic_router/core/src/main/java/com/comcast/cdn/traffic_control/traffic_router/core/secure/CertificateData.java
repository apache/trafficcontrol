package com.comcast.cdn.traffic_control.traffic_router.core.secure;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

@JsonIgnoreProperties(ignoreUnknown = true)
public class CertificateData {
	@JsonProperty
	private String deliveryservice;

	@JsonProperty
	private Certificate certificate;

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

	@Override
	@SuppressWarnings("PMD")
	public boolean equals(final Object o) {
		if (this == o) return true;
		if (o == null || getClass() != o.getClass()) return false;

		final CertificateData that = (CertificateData) o;

		if (deliveryservice != null ? !deliveryservice.equals(that.deliveryservice) : that.deliveryservice != null)
			return false;
		return certificate != null ? certificate.equals(that.certificate) : that.certificate == null;

	}

	@Override
	public int hashCode() {
		int result = deliveryservice != null ? deliveryservice.hashCode() : 0;
		result = 31 * result + (certificate != null ? certificate.hashCode() : 0);
		return result;
	}
}
