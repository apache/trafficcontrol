package com.comcast.cdn.traffic_control.traffic_router.shared;

import javax.management.AttributeChangeNotification;
import javax.management.NotificationBroadcasterSupport;
import java.util.List;

public class DeliveryServiceCertificates extends NotificationBroadcasterSupport implements DeliveryServiceCertificatesMBean {
	private List<CertificateData> certificateDataList;
	private long sequenceNumber = 1L;

	@Override
	public List<CertificateData> getCertificateDataList() {
		return certificateDataList;
	}

	@Override
	public void setCertificateDataList(List<CertificateData> certificateDataList) {
		List<CertificateData> oldCertificateDataList = this.certificateDataList;
		this.certificateDataList = certificateDataList;

		sendNotification(new AttributeChangeNotification(this, sequenceNumber, System.currentTimeMillis(), "CertificateDataList Changed",
			"CertificateDataList", "List<CertificateData>", oldCertificateDataList, this.certificateDataList));
		sequenceNumber++;
	}

	@Override
	public boolean equals(Object o) {
		if (this == o) return true;
		if (o == null || getClass() != o.getClass()) return false;

		DeliveryServiceCertificates that = (DeliveryServiceCertificates) o;

		if (sequenceNumber != that.sequenceNumber) return false;
		return certificateDataList != null ? certificateDataList.equals(that.certificateDataList) : that.certificateDataList == null;

	}

	@Override
	public int hashCode() {
		int result = certificateDataList != null ? certificateDataList.hashCode() : 0;
		result = 31 * result + (int) (sequenceNumber ^ (sequenceNumber >>> 32));
		return result;
	}
}
