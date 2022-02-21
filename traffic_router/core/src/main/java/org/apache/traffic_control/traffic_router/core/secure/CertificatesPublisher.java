package org.apache.traffic_control.traffic_router.core.secure;

/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 * 
 *   http://www.apache.org/licenses/LICENSE-2.0
 * 
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */


import org.apache.traffic_control.traffic_router.core.config.CertificateChecker;
import org.apache.traffic_control.traffic_router.core.ds.DeliveryService;
import org.apache.traffic_control.traffic_router.core.router.TrafficRouterManager;
import org.apache.traffic_control.traffic_router.shared.CertificateData;
import org.apache.traffic_control.traffic_router.shared.DeliveryServiceCertificatesMBean;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

import javax.management.Attribute;
import javax.management.ObjectName;
import java.lang.management.ManagementFactory;
import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.BlockingQueue;
import java.util.concurrent.TimeUnit;

public class CertificatesPublisher {
	private final static Logger LOGGER = LogManager.getLogger(CertificatesPublisher.class);
	private JsonNode deliveryServicesJson;
	private List<DeliveryService> deliveryServices = new ArrayList<>();
	private boolean updated = false;
	private boolean running = true;
	final Thread worker;


	@SuppressWarnings("PMD.AvoidCatchingThrowable")
	public CertificatesPublisher(final BlockingQueue<List<CertificateData>> certificatesQueue, final BlockingQueue<Boolean> publishStatusQueue,
	                             final CertificateChecker certificateChecker, final TrafficRouterManager trafficRouterManager) {
		worker = new Thread(() -> {
			while (running) {
				try {
					final List<CertificateData> certificateDataList = certificatesQueue.take();
					if (certificateDataList == null) {
						continue;
					}

					updated = false;
					if (certificateChecker.certificatesAreValid(certificateDataList, deliveryServicesJson)) {
						deliveryServices.forEach(ds -> {
							final boolean hasX509Cert = certificateChecker.hasCertificate(certificateDataList, ds.getId());
							ds.setHasX509Cert(hasX509Cert);
						});

						publishCertificateList(certificateDataList);

						if (updated == false) {
							publishStatusQueue.poll(5, TimeUnit.MICROSECONDS);
						}

						trafficRouterManager.trackEvent("lastHttpsCertificatesUpdate");
					} else {
						trafficRouterManager.trackEvent("lastInvalidHttpsCertificates");
					}
				} catch (Throwable e) {
					if (!running) {
						return;
					}

					LOGGER.warn("Interrupted while waiting for new certificate data list, trying again...",e);
				}
			}
		});

		worker.start();
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

	public JsonNode getDeliveryServicesJson() {
		return deliveryServicesJson;
	}

	public void setDeliveryServicesJson(final JsonNode deliveryServicesJson) {
		updated = true;
		this.deliveryServicesJson = deliveryServicesJson;
	}

	public List<DeliveryService> getDeliveryServices() {
		return deliveryServices;
	}

	public void setDeliveryServices(final List<DeliveryService> deliveryServices) {
		this.deliveryServices = deliveryServices;
	}

	public void destroy() {
		LOGGER.warn("Detected destroy setting running to false");
		running = false;
		worker.interrupt();
	}
}
