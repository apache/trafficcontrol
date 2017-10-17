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
import com.comcast.cdn.traffic_control.traffic_monitor.health.CacheState;
import com.comcast.cdn.traffic_control.traffic_monitor.health.CacheStateUpdater;
import com.comcast.cdn.traffic_control.traffic_monitor.health.HealthDeterminer;
import com.comcast.cdn.traffic_control.traffic_monitor.wicket.models.CacheDataModel;
import com.ning.http.client.Response;
import org.apache.wicket.ajax.json.JSONObject;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.powermock.core.classloader.annotations.PrepareForTest;
import org.powermock.modules.junit4.PowerMockRunner;

import java.util.concurrent.CancellationException;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.Future;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.TimeoutException;
import java.util.concurrent.atomic.AtomicInteger;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.equalTo;
import static org.hamcrest.Matchers.nullValue;
import static org.junit.Assert.fail;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.spy;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.verifyZeroInteractions;
import static org.mockito.Mockito.when;
import static org.powermock.api.mockito.PowerMockito.mockStatic;

@PrepareForTest({CacheStateUpdater.class, System.class})
@RunWith(PowerMockRunner.class)
public class CacheStateUpdaterTest {
	@Test
	public void itBlowsUpWhenUpdateNotCalledBeforeOnCompleted() throws Exception {

		CacheState cacheState = new CacheState("cache1");
		CacheDataModel cacheDataModel = mock(CacheDataModel.class);
		CacheStateUpdater cacheStateUpdater = new CacheStateUpdater(cacheState, cacheDataModel);

		try {
			cacheStateUpdater.onCompleted(mock(Response.class));
			fail("Should have caught NPE!");
		} catch (NullPointerException e) {
			// expected
		}
	}

	@Test
	public void itHandlesResponsesAfterUpdateCalled() throws Exception {
		// Having to create json for always creating a cache is cumbersome
		JSONObject cacheJson = new JSONObject();
		cacheJson.put("ip", "192.168.99.100");
		cacheJson.put("status", "online");
		cacheJson.put("locationId", "cache-location");
		cacheJson.put("profile", "edge");
		cacheJson.put("fqdn", "cache1.example.com");
		cacheJson.put("type", "cache-type");
		cacheJson.put("port", 9876);

		Cache cache = spy(new Cache("cache1", cacheJson));
		when(cache.getStatisticsUrl()).thenReturn("http://cache1.example.com/astats");

		CacheState cacheState = spy(new CacheState("cache1"));

		when(cacheState.getCache()).thenReturn(cache);
		CacheDataModel cacheDataModel = mock(CacheDataModel.class);

		CacheStateUpdater cacheStateUpdater = new CacheStateUpdater(cacheState, cacheDataModel);

		mockStatic(System.class);
		when(System.currentTimeMillis()).thenReturn(1455559116000L);

		HealthDeterminer healthDeterminer = new HealthDeterminer();
		cacheStateUpdater.update(healthDeterminer, new AtomicInteger(0), System.currentTimeMillis() + 1000L);
		Response response = mock(Response.class);

		when(response.getStatusCode()).thenReturn(200);
		when(response.getResponseBody()).thenReturn("{" +
			"\"global\":{\"kbps\":\"1234.56\",\"tps\":\"3000\"}, " +
			"\"system\":{\"loadavg\":\"0.75\"} " +
			"}");

		when(System.currentTimeMillis()).thenReturn(1455559117000L);

		assertThat(cacheStateUpdater.onCompleted(response), equalTo(200));

		assertThat(cacheState.getLastValue("queryTime"), equalTo("1000"));
		assertThat(cacheState.getLastValue("stateUrl"), equalTo("http://cache1.example.com/astats"));
		assertThat(cacheState.getLastValue("ats.kbps"), equalTo("1234.56"));
		assertThat(cacheState.getLastValue("ats.tps"), equalTo("3000"));
		assertThat(cacheState.getLastValue("ats.tps"), equalTo("3000"));
		assertThat(cacheState.getLastValue("system.loadavg"), equalTo("0.75"));
		assertThat(cacheState.getLastValue("status"), equalTo("online"));
		assertThat(cacheState.getLastValue("error-string"), nullValue());
		verify(cache).setState(cacheState, healthDeterminer);

		when(response.getResponseBody()).thenReturn("{" +
			"\"ats\":{\"kbps\":\"987.65\",\"tps\":\"5555\"}, " +
			"\"system\":{\"loadavg\":\"0.75\"} " +
			"}");

		cacheStateUpdater.onCompleted(response);

		assertThat(cacheState.getLastValue("ats.kbps"), equalTo("987.65"));
		assertThat(cacheState.getLastValue("ats.tps"), equalTo("5555"));
	}

	@Test
	public void itSetsErrorWhenResponseIsNot200() throws Exception {
		CacheState cacheState = mock(CacheState.class);

		CacheStateUpdater cacheStateUpdater = new CacheStateUpdater(cacheState, mock(CacheDataModel.class));
		cacheStateUpdater.update(mock(HealthDeterminer.class), mock(AtomicInteger.class), 0);

		Response response = mock(Response.class);
		when(response.getStatusCode()).thenReturn(400);
		when(response.getStatusText()).thenReturn("Bad Request");

		cacheStateUpdater.onCompleted(response);
		verify(cacheState).setError("400 - Bad Request");
	}

	@Test
	public void itHandlesExceptions() throws Exception {
		Cache cache = mock(Cache.class);
		when(cache.getStatisticsUrl()).thenReturn("http://cache1.example.com/astats");

		CacheState cacheState = mock(CacheState.class);
		when(cacheState.getCache()).thenReturn(cache);

		CacheDataModel errorCount = new CacheDataModel("error count");
		CacheStateUpdater cacheStateUpdater = new CacheStateUpdater(cacheState, errorCount);

		mockStatic(System.class);
		when(System.currentTimeMillis()).thenReturn(1455559116000L);

		AtomicInteger failCount = new AtomicInteger(0);
		cacheStateUpdater.update(mock(HealthDeterminer.class), failCount, 0);

		when(System.currentTimeMillis()).thenReturn(1455559117111L);
		cacheStateUpdater.onThrowable(new RuntimeException("boom"));

		verify(cacheState).setError("java.lang.RuntimeException: boom");
		verify(cacheState).putDataPoint("queryTime", "1111");
		assertThat(failCount.get(), equalTo(1));
		assertThat(errorCount.getValue(), equalTo("1"));
	}

	@Test
	public void itHandlesCancellations() throws Exception {
		Cache cache = mock(Cache.class);
		when(cache.getStatisticsUrl()).thenReturn("http://cache1.example.com/astats");

		CacheState cacheState = mock(CacheState.class);
		when(cacheState.getCache()).thenReturn(cache);

		CacheDataModel errorCount = new CacheDataModel("error count");
		CacheStateUpdater cacheStateUpdater = new CacheStateUpdater(cacheState, errorCount);

		mockStatic(System.class);
		when(System.currentTimeMillis()).thenReturn(1455559116000L);

		AtomicInteger failCount = new AtomicInteger(0);
		cacheStateUpdater.update(mock(HealthDeterminer.class), failCount, 0);

		when(System.currentTimeMillis()).thenReturn(1455559117111L);
		cacheStateUpdater.onThrowable(new CancellationException("timed out"));

		verify(cacheState).setError("java.util.concurrent.CancellationException: timed out");
		verify(cacheState).putDataPoint("queryTime", "1111");
		assertThat(failCount.get(), equalTo(0));
		assertThat(errorCount.getValue(), equalTo("1"));
	}

	@Test
	public void itDoesNotCancelCompletedOrPreviouslyCancelledUpdates() {
		CacheState cacheState = new CacheState("cache1");

		CacheStateUpdater cacheStateUpdater = new CacheStateUpdater(cacheState, mock(CacheDataModel.class));
		mockStatic(System.class);
		when(System.currentTimeMillis()).thenReturn(1455559116000L);

		cacheStateUpdater.update(mock(HealthDeterminer.class), mock(AtomicInteger.class), 0);

		AtomicInteger cancelCount = mock(AtomicInteger.class);
		assertThat(cacheStateUpdater.completeFetchStatistics(cancelCount), equalTo(true));
		verifyZeroInteractions(cancelCount);

		Future future = mock(Future.class);

		when(future.cancel(true)).thenThrow(new RuntimeException("Test failed, should not have called this"));

		cacheStateUpdater.setFuture(future);
		when(future.isDone()).thenReturn(true);

		assertThat(cacheStateUpdater.completeFetchStatistics(cancelCount), equalTo(true));
		verifyZeroInteractions(cancelCount);

		when(future.isDone()).thenReturn(false);
		when(future.isCancelled()).thenReturn(true);

		assertThat(cacheStateUpdater.completeFetchStatistics(cancelCount), equalTo(true));
		verifyZeroInteractions(cancelCount);
	}

	@Test
	public void itCancelsUpdatesThatHaveNotMetDeadline() {
		CacheState cacheState = new CacheState("cache1");

		CacheStateUpdater cacheStateUpdater = new CacheStateUpdater(cacheState, mock(CacheDataModel.class));
		FutureExample future = new FutureExample();
		cacheStateUpdater.setFuture(future);

		mockStatic(System.class);
		cacheStateUpdater.update(new HealthDeterminer(), new AtomicInteger(0), 1455559116000L);

		when(System.currentTimeMillis()).thenReturn(1455559117000L);
		AtomicInteger cancelCount = mock(AtomicInteger.class);
		assertThat(cacheStateUpdater.completeFetchStatistics(cancelCount), equalTo(true));
		verify(cancelCount).incrementAndGet();
		assertThat(future.cancelCalls, equalTo(1));
		assertThat(future.cancelLastCalledWith, equalTo(true));
	}

	@Test
	public void itDoesNotCancelUpdatesThatAreEarlierThanDeadline() {
		CacheState cacheState = new CacheState("cache1");

		CacheStateUpdater cacheStateUpdater = new CacheStateUpdater(cacheState, mock(CacheDataModel.class));
		FutureExample future = new FutureExample();
		cacheStateUpdater.setFuture(future);

		mockStatic(System.class);
		cacheStateUpdater.update(new HealthDeterminer(), new AtomicInteger(0), 1455559116000L);

		when(System.currentTimeMillis()).thenReturn(1455559115000L);
		AtomicInteger cancelCount = mock(AtomicInteger.class);
		assertThat(cacheStateUpdater.completeFetchStatistics(cancelCount), equalTo(false));
		verifyZeroInteractions(cancelCount);
		assertThat(future.cancelCalls, equalTo(0));
	}

	class FutureExample implements Future {
		int cancelCalls = 0;
		boolean cancelLastCalledWith = false;

		@Override
		public boolean cancel(boolean mayInterruptIfRunning) {
			cancelCalls++;
			cancelLastCalledWith = mayInterruptIfRunning;
			return true;
		}

		@Override
		public boolean isCancelled() {
			return false;
		}

		@Override
		public boolean isDone() {
			return false;
		}

		@Override
		public Object get() throws InterruptedException, ExecutionException {
			return null;
		}

		@Override
		public Object get(long timeout, TimeUnit unit) throws InterruptedException, ExecutionException, TimeoutException {
			return null;
		}
	}
}
