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

package shared;

import org.apache.traffic_control.traffic_router.shared.CertificateData;
import org.apache.traffic_control.traffic_router.shared.DeliveryServiceCertificates;
import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.mockito.ArgumentCaptor;
import org.powermock.api.mockito.PowerMockito;
import org.powermock.core.classloader.annotations.PrepareForTest;
import org.powermock.modules.junit4.PowerMockRunner;

import javax.management.AttributeChangeNotification;
import java.util.ArrayList;
import java.util.List;

import static org.hamcrest.CoreMatchers.equalTo;
import static org.hamcrest.MatcherAssert.assertThat;
import static org.mockito.Mockito.spy;
import static org.mockito.Mockito.times;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

@RunWith(PowerMockRunner.class)
@PrepareForTest({DeliveryServiceCertificates.class, System.class})
public class DeliveryServiceCertificatesTest {
	@Before
	public void before() throws Exception {
		PowerMockito.mockStatic(System.class);
		when(System.currentTimeMillis()).thenReturn(1234L);
	}

	@Test
	public void itSendsNotificationWhenNewCertData() {

		DeliveryServiceCertificates deliveryServiceCertificates = spy(new DeliveryServiceCertificates());
		ArgumentCaptor<AttributeChangeNotification> captor = ArgumentCaptor.forClass(AttributeChangeNotification.class);
		List<CertificateData> certificateDataList = new ArrayList<>();
		deliveryServiceCertificates.setCertificateDataList(certificateDataList);

		verify(deliveryServiceCertificates, times(1)).sendNotification(captor.capture());

		AttributeChangeNotification notification = captor.getValue();
		assertThat(notification.getNewValue(), equalTo(certificateDataList));
		assertThat(notification.getAttributeName(), equalTo("CertificateDataList"));
		assertThat(notification.getAttributeType(), equalTo("List<CertificateData>"));
		assertThat(notification.getMessage(), equalTo("CertificateDataList Changed"));
		assertThat(notification.getTimeStamp(), equalTo(1234L));
		assertThat(notification.getSequenceNumber(), equalTo(1L));
		assertThat(notification.getSource(), equalTo(deliveryServiceCertificates));
	}
}
