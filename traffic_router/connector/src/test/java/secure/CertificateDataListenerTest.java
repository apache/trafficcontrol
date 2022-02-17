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

package secure;

import org.apache.traffic_control.traffic_router.secure.CertificateDataListener;
import org.apache.traffic_control.traffic_router.secure.CertificateRegistry;
import org.apache.traffic_control.traffic_router.shared.CertificateData;
import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.powermock.api.mockito.PowerMockito;
import org.powermock.core.classloader.annotations.PowerMockIgnore;
import org.powermock.core.classloader.annotations.PrepareForTest;
import org.powermock.modules.junit4.PowerMockRunner;

import javax.management.AttributeChangeNotification;
import javax.management.Notification;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;

import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.times;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

@RunWith(PowerMockRunner.class)
@PrepareForTest(CertificateRegistry.class)
@PowerMockIgnore("javax.management.*")
public class CertificateDataListenerTest {

	private CertificateRegistry certificateRegistry;

	@Before
	public void before() throws Exception {
		certificateRegistry = mock(CertificateRegistry.class);
		PowerMockito.mockStatic(CertificateRegistry.class);
		when(CertificateRegistry.getInstance()).thenReturn(certificateRegistry);
	}

	@Test
	public void itImportsCertificateDataToRegistry() throws Exception {
		List<CertificateData> oldList = new ArrayList<>();
		List<CertificateData> newList = new ArrayList<>();

		Object notifier = "notifier";

		Notification notification = new AttributeChangeNotification(notifier, 1L, System.currentTimeMillis(),
			"CertificateDataList Changed", "CertificateDataList", "List<CertificateDataList>", oldList, newList);

		CertificateDataListener certificateDataListener = new CertificateDataListener();
		certificateDataListener.handleNotification(notification, null);
		verify(certificateRegistry).importCertificateDataList(newList);
	}

	@Test
	public void itIgnoresBadInput() throws Exception {
		Notification notification = new Notification("notifier", "source", 1L, "hello world");
		CertificateDataListener certificateDataListener = new CertificateDataListener();
		certificateDataListener.handleNotification(notification, null);
		verify(certificateRegistry, times(0)).importCertificateDataList(any());

		List<String> badData = Arrays.asList("foo", "bar", "baz");

		notification = new AttributeChangeNotification("notifier", 1L, System.currentTimeMillis(),
			"CertificateDataList Changed", "CertificateDataList", "List<CertificateDataList>", null, badData);

		certificateDataListener.handleNotification(notification, null);
		verify(certificateRegistry, times(0)).importCertificateDataList(any());
	}
}
