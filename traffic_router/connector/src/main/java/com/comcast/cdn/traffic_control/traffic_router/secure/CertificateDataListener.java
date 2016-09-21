package com.comcast.cdn.traffic_control.traffic_router.secure;

import com.comcast.cdn.traffic_control.traffic_router.shared.CertificateData;

import javax.management.AttributeChangeNotification;
import javax.management.Notification;
import javax.management.NotificationListener;
import java.util.ArrayList;
import java.util.List;

public class CertificateDataListener implements NotificationListener {
	@Override
	public void handleNotification(final Notification notification, final Object handback) {
		if (!(notification instanceof AttributeChangeNotification)) {
			return;
		}

		List<CertificateData> certificateDataList = new ArrayList<>();

		final Object newValue = ((AttributeChangeNotification) notification).getNewValue();
		if (certificateDataList.getClass().isInstance(newValue)) {
			certificateDataList = (List<CertificateData>) newValue;
			CertificateRegistry.getInstance().importCertificateDataList(certificateDataList);
		}
	}
}
