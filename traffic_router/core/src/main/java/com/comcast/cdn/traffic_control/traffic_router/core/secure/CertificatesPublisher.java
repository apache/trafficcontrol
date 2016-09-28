package com.comcast.cdn.traffic_control.traffic_router.core.secure;

import com.comcast.cdn.traffic_control.traffic_router.core.config.CertificateChecker;
import com.comcast.cdn.traffic_control.traffic_router.core.ds.DeliveryService;
import com.comcast.cdn.traffic_control.traffic_router.shared.CertificateData;
import com.comcast.cdn.traffic_control.traffic_router.shared.DeliveryServiceCertificatesMBean;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.apache.log4j.Logger;
import org.json.JSONObject;

import javax.management.Attribute;
import javax.management.ObjectName;
import java.lang.management.ManagementFactory;
import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.BlockingQueue;
import java.util.concurrent.TimeUnit;

public class CertificatesPublisher {
	private final static Logger LOGGER = Logger.getLogger(CertificatesPublisher.class);
	private JSONObject deliveryServicesJson;
	private List<DeliveryService> deliveryServices = new ArrayList<>();

	@SuppressWarnings("PMD.AvoidCatchingThrowable")
	public CertificatesPublisher(final BlockingQueue<List<CertificateData>> certificatesQueue, final BlockingQueue<Boolean> publishStatusQueue, final CertificateChecker certificateChecker) {

		new Thread(() -> {
			while (true) {
				try {
					final List<CertificateData> certificateDataList = certificatesQueue.take();
					if (certificateDataList == null) {
						continue;
					}

					if (certificateChecker.certificatesAreValid(certificateDataList, deliveryServicesJson)) {
						deliveryServices.forEach(ds -> {
							final boolean hasX509Cert = certificateChecker.hasCertificate(certificateDataList, ds.getId());
							ds.setHasX509Cert(hasX509Cert);
						});
						publishCertificateList(certificateDataList);
						publishStatusQueue.poll(2, TimeUnit.SECONDS);
					}
				} catch (Throwable e) {
					LOGGER.warn("Interrupted while waiting for new certificate data list, trying again...",e);
				}
			}
		}).start();
	}

	private void publishCertificateList(final List<CertificateData> certificateDataList) {
		try {
			final ObjectName objectName = new ObjectName(DeliveryServiceCertificatesMBean.OBJECT_NAME);
			ManagementFactory.getPlatformMBeanServer().setAttribute(objectName,
				new Attribute("CertificateDataListString", new ObjectMapper().writeValueAsString(certificateDataList)));
		} catch (Exception e) {
			LOGGER.error("Failed to add certificate data list as management MBean! " + e.getClass().getSimpleName() + ": " + e.getMessage(), e);
		}
	}

	public JSONObject getDeliveryServicesJson() {
		return deliveryServicesJson;
	}

	public void setDeliveryServicesJson(final JSONObject deliveryServicesJson) {
		this.deliveryServicesJson = deliveryServicesJson;
	}

	public List<DeliveryService> getDeliveryServices() {
		return deliveryServices;
	}

	public void setDeliveryServices(final List<DeliveryService> deliveryServices) {
		this.deliveryServices = deliveryServices;
	}
}
