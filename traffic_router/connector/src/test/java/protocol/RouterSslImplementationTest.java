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

package protocol;

import com.comcast.cdn.traffic_control.traffic_router.protocol.RouterSslImplementation;
import com.comcast.cdn.traffic_control.traffic_router.protocol.RouterSslServerSocketFactory;
import org.apache.tomcat.util.net.jsse.JSSEFactory;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.powermock.core.classloader.annotations.PrepareForTest;
import org.powermock.modules.junit4.PowerMockRunner;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.equalTo;
import static org.hamcrest.Matchers.instanceOf;

@RunWith(PowerMockRunner.class)
@PrepareForTest({RouterSslImplementation.class, JSSEFactory.class})
public class RouterSslImplementationTest {
	@Test
	public void itReturnsItsName() throws Exception {
		assertThat(new RouterSslImplementation().getImplementationName(), equalTo("RouterSslImplementation"));
	}

	@Test
	public void itReturnsRouterSslServerSocketFactory() throws Exception {
		assertThat(new RouterSslImplementation().getServerSocketFactory(), instanceOf(RouterSslServerSocketFactory.class));
	}
}
