package com.comcast.cdn.traffic_control.traffic_router.shared;

import java.util.List;

public interface DeliveryServiceCertificatesMBean {
	String OBJECT_NAME = "traffic-router:name=DeliveryServiceCertificates";
	List<CertificateData> getCertificateDataList();
	void setCertificateDataList(List<CertificateData> certificateDataList);
}
