package health;

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


import com.comcast.cdn.traffic_control.traffic_monitor.config.Cache;
import com.comcast.cdn.traffic_control.traffic_monitor.health.CacheStateUpdater;
import com.comcast.cdn.traffic_control.traffic_monitor.health.CacheStatisticsClient;
import com.ning.http.client.AsyncHttpClient;
import com.ning.http.client.ListenableFuture;
import com.ning.http.client.ProxyServer;
import com.ning.http.client.Request;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.powermock.core.classloader.annotations.PrepareForTest;
import org.powermock.modules.junit4.PowerMockRunner;

import static org.mockito.Matchers.any;
import static org.mockito.Mockito.doReturn;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.spy;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;
import static org.powermock.api.mockito.PowerMockito.whenNew;

@PrepareForTest({CacheStatisticsClient.class, AsyncHttpClient.class, ProxyServer.class})
@RunWith(PowerMockRunner.class)
public class CacheStatisticsClientTest {
	@Test
	public void itExecutesAsynchronousRequest() throws Exception {

		ListenableFuture listenableFuture = mock(ListenableFuture.class);
		AsyncHttpClient asyncHttpClient = spy(new AsyncHttpClient());
		doReturn(listenableFuture).when(asyncHttpClient).executeRequest(any(Request.class), any(CacheStateUpdater.class));

		whenNew(AsyncHttpClient.class).withNoArguments().thenReturn(asyncHttpClient);

		Cache cache = mock(Cache.class);
		when(cache.getQueryIp()).thenReturn("192.168.99.100");
		when(cache.getQueryPort()).thenReturn(0);
		when(cache.getStatisticsUrl()).thenReturn("http://cache1.example.com/astats");

		CacheStateUpdater cacheStateUpdater = mock(CacheStateUpdater.class);
		CacheStatisticsClient cacheStatisticsClient = new CacheStatisticsClient();

		cacheStatisticsClient.fetchCacheStatistics(cache, cacheStateUpdater);
		verify(cacheStateUpdater).setFuture(listenableFuture);
	}
}
