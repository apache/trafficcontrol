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

package org.apache.traffic_control.traffic_router.core.dns;

import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.ExecutorService;

import org.junit.Before;
import org.junit.Test;

import org.apache.traffic_control.traffic_router.core.dns.protocol.Protocol;

import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.verify;

public class NameServerMainTest {

    private ExecutorService executorService;
    private List<Protocol> protocols;
    private Protocol p1;
    private Protocol p2;
    private NameServerMain main;

    @Before
    public void setUp() throws Exception {
        p1 = mock(Protocol.class);
        p2 = mock(Protocol.class);
        protocols = new ArrayList<>();
        protocols.add(p1);
        protocols.add(p2);
        executorService = mock(ExecutorService.class);

        main = new NameServerMain();
        main.setProtocols(protocols);
        main.setProtocolService(executorService);
    }

    @Test
    public void testDestroy() throws Exception {
        main.destroy();
        verify(p1).shutdown();
        verify(p2).shutdown();
    }

    @Test
    public void testInit() throws Exception {
        main.init();
        verify(executorService).submit(p1);
        verify(executorService).submit(p2);
    }

}
