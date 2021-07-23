/*
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package org.apache.traffic_control.traffic_router.shared;

import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.ObjectMapper;

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
	public void setCertificateDataList(final List<CertificateData> certificateDataList) {
		final List<CertificateData> oldCertificateDataList = this.certificateDataList;
		this.certificateDataList = certificateDataList;

		sendNotification(new AttributeChangeNotification(this, sequenceNumber, System.currentTimeMillis(), "CertificateDataList Changed",
			"CertificateDataList", "List<CertificateData>", oldCertificateDataList, this.certificateDataList));
		sequenceNumber++;
	}

	@Override
	@SuppressWarnings("PMD.AvoidThrowingRawExceptionTypes")
	public void setCertificateDataListString(final String certificateDataListString) {
		try {
			final List<CertificateData> certificateDataList = new ObjectMapper().
				readValue(certificateDataListString, new TypeReference<List<CertificateData>>() { });
			setCertificateDataList(certificateDataList);
		} catch (Exception e) {
			throw new RuntimeException("Failed to convert json certificate data list to list of CertificateData objects", e);
		}
	}

	@Override
	public boolean equals(final Object o) {
		if (this == o) {
			return true;
		}

		if (o == null || getClass() != o.getClass()) {
			return false;
		}

		final DeliveryServiceCertificates that = (DeliveryServiceCertificates) o;

		if (sequenceNumber != that.sequenceNumber) {
			return false;
		}

		return certificateDataList != null ? certificateDataList.equals(that.certificateDataList) : that.certificateDataList == null;

	}

	@Override
	public int hashCode() {
		int result = certificateDataList != null ? certificateDataList.hashCode() : 0;
		result = 31 * result + (int) (sequenceNumber ^ (sequenceNumber >>> 32));
		return result;
	}
}
