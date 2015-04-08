/*
 * Copyright 2015 Comcast Cable Communications Management, LLC
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

package com.comcast.cdn.traffic_control.traffic_router.core.dns;

import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.ExecutorService;

import org.jmock.Expectations;
import org.jmock.Mockery;
import org.jmock.integration.junit4.JMock;
import org.jmock.integration.junit4.JUnit4Mockery;
import org.jmock.lib.legacy.ClassImposteriser;
import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;

import com.comcast.cdn.traffic_control.traffic_router.core.dns.NameServerMain;
import com.comcast.cdn.traffic_control.traffic_router.core.dns.protocol.Protocol;

@RunWith(JMock.class)
public class NameServerMainTest {

    private final Mockery context = new JUnit4Mockery() {
        {
            setImposteriser(ClassImposteriser.INSTANCE);
        }
    };

    private ExecutorService executorService;
    private List<Protocol> protocols;
    private Protocol p1;
    private Protocol p2;
    private NameServerMain main;

    @Before
    public void setUp() throws Exception {
        p1 = context.mock(Protocol.class, "p1");
        p2 = context.mock(Protocol.class, "p2");
        protocols = new ArrayList<Protocol>();
        protocols.add(p1);
        protocols.add(p2);
        executorService = context.mock(ExecutorService.class);

        main = new NameServerMain();
        main.setProtocols(protocols);
        main.setProtocolService(executorService);
    }

    @Test
    public void testDestroy() throws Exception {
        context.checking(new Expectations() {
            {
                one(p1).shutdown();
                one(p2).shutdown();
                one(executorService).shutdownNow();
            }
        });
        main.destroy();
    }

    @Test
    public void testInit() throws Exception {
        context.checking(new Expectations() {
            {
                one(executorService).submit(p1);
                one(executorService).submit(p2);
            }
        });
        main.init();
    }

}
