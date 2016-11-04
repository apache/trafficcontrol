package data;

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


import com.comcast.cdn.traffic_control.traffic_monitor.data.DataPoint;
import com.comcast.cdn.traffic_control.traffic_monitor.data.StatisticsLog;
import com.comcast.cdn.traffic_control.traffic_monitor.health.DeliveryServiceStateRegistry;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.powermock.core.classloader.annotations.PrepareForTest;
import org.powermock.modules.junit4.PowerMockRunner;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.Deque;
import java.util.HashSet;
import java.util.List;
import java.util.concurrent.CyclicBarrier;
import java.util.concurrent.atomic.AtomicInteger;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.contains;
import static org.hamcrest.Matchers.containsInAnyOrder;
import static org.hamcrest.Matchers.equalTo;
import static org.hamcrest.Matchers.nullValue;
import static org.powermock.api.mockito.PowerMockito.mockStatic;
import static org.powermock.api.mockito.PowerMockito.when;

@PrepareForTest({StatisticsLog.class, System.class})
@RunWith(PowerMockRunner.class)
public class StatisticsLogTest {
	@Test
	public void itLogsDataPoints() {
		StatisticsLog statisticsLog = new StatisticsLog();
		assertThat(statisticsLog.get("foo"), nullValue());

		statisticsLog.putDataPoint("foo", "okay");
		assertThat(statisticsLog.get("foo").getLast().getValue(), equalTo("okay"));
		assertThat(statisticsLog.get("foo").getLast().getIndex(), equalTo(0L));
		assertThat(statisticsLog.get("foo").getLast().getSpan(), equalTo(1));

		statisticsLog.putDataPoint("foo", "okay");
		assertThat(statisticsLog.get("foo").size(), equalTo(1));
		assertThat(statisticsLog.get("foo").getLast().getIndex(), equalTo(0L));
		assertThat(statisticsLog.get("foo").getLast().getSpan(), equalTo(2));

		statisticsLog.putDataPoint("foo", "bar");
		assertThat(statisticsLog.get("foo").size(), equalTo(2));

		statisticsLog.putDataPoint("foo2", "okay");
		assertThat(statisticsLog.getKeys(), containsInAnyOrder("foo2", "foo"));

		assertThat(statisticsLog.hasValue("foo"), equalTo(true));
		assertThat(statisticsLog.hasValue("bar"), equalTo(false));

		assertThat(statisticsLog.getLastValue("foo"), equalTo("bar"));
		assertThat(statisticsLog.getLastValue("baz"), nullValue());

		assertThat(statisticsLog.getValue("foo", 0), equalTo("bar"));
		assertThat(statisticsLog.getValue("baz", 0), nullValue());

		assertThat(statisticsLog.getBool("foo"), equalTo(false));
		assertThat(statisticsLog.getLong("foo"), equalTo(0L));
		assertThat(statisticsLog.getDouble("foo"), equalTo(0.0));

		assertThat(statisticsLog.getTime(0L), equalTo(0L));
	}

	@Test
	public void itIndexesPerUpdate() {
		mockStatic(System.class);
		when(System.currentTimeMillis()).thenReturn(1455224271177L);

		StatisticsLog statisticsLog = new StatisticsLog();
		statisticsLog.putDataPoint("kbps","1234.5");

		statisticsLog.prepareForUpdate("cache1", 1000L);

		assertThat(statisticsLog.getValue("kbps", 0), equalTo("1234.5"));
		assertThat(statisticsLog.getValue("kbps", 1), nullValue());

		assertThat(statisticsLog.getTime(0L), equalTo(0L));
		assertThat(statisticsLog.getTime(1L), equalTo(1455224271177L));

		statisticsLog.putDataPoint("kbps", "5432.1");

		assertThat(statisticsLog.getValue("kbps", 1), equalTo("5432.1"));
	}

	@Test
	public void itKeepsAtLeastOneStatistic() {
		mockStatic(System.class);
		when(System.currentTimeMillis()).thenReturn(1455224271177L);

		StatisticsLog statisticsLog = new StatisticsLog();
		statisticsLog.prepareForUpdate("cache1", 1000L);
		statisticsLog.putDataPoint("kbps","1234.5");
		statisticsLog.putDataPoint("kbps","5432.1");

		when(System.currentTimeMillis()).thenReturn(1455224271177L + 1001L);

		statisticsLog.prepareForUpdate("cache1", 1000L);

		assertThat(statisticsLog.getValue("kbps", 1), equalTo("5432.1"));
	}

	@Test
	public void itAddsPlaceHolder() {
		mockStatic(System.class);
		when(System.currentTimeMillis()).thenReturn(1455224271177L);

		StatisticsLog statisticsLog = new StatisticsLog();
		statisticsLog.prepareForUpdate("cache1", 1000L);
		statisticsLog.putDataPoint("kbps","1234.5");

		when(System.currentTimeMillis()).thenReturn(1455224271177L + 100L);
		statisticsLog.prepareForUpdate("cache1", 1000L);
		statisticsLog.putDataPoint("tps", "123000");

		when(System.currentTimeMillis()).thenReturn(1455224271177L + 200L);
		statisticsLog.prepareForUpdate("cache1", 1000L);

		assertThat(statisticsLog.getLastValue("kbps"), nullValue());
		assertThat(statisticsLog.getLastValue("tps"), equalTo("123000"));
	}


	@Test
	public void itFilters() {
		mockStatic(System.class);
		when(System.currentTimeMillis()).thenReturn(1455224271177L);

		StatisticsLog statisticsLog = new StatisticsLog();
		statisticsLog.prepareForUpdate("cache1", 1000L);
		statisticsLog.putDataPoint("kbps","1234.5");

		assertThat(statisticsLog.filter(10, new String[] {"tps"},false, false).size(), equalTo(0));
		assertThat(statisticsLog.filter(10, new String[] {"kbps", "tps"},false, false).get("kbps").getFirst().getValue(), equalTo("1234.5"));
		assertThat(statisticsLog.filter(10, new String[] {"KBPS"},true, false).get("kbps").getFirst().getValue(), equalTo("1234.5"));

		statisticsLog.addHiddenStats(new HashSet(Arrays.asList("kbps", "errors")));
		assertThat(statisticsLog.filter(10, new String[] {"kbps", "tps"},false, false).size(), equalTo(0));
		assertThat(statisticsLog.filter(10, new String[] {"kbps", "tps"},false, true).get("kbps").getFirst().getValue(), equalTo("1234.5"));

		statisticsLog.putDataPoint("kbps", "5432.1");
		assertThat(statisticsLog.filter(1, new String[] {"kbps", "tps"},false, true).get("kbps").size(), equalTo(1));
		assertThat(statisticsLog.filter(0, new String[] {"kbps", "tps"},false, true).get("kbps").size(), equalTo(2));

		statisticsLog.putDataPoint("rxBytes", "1000");

		assertThat(statisticsLog.filter(0, new String[] {"kbps"},false, true).keySet().size(), equalTo(1));
		assertThat(statisticsLog.filter(0, null,false, true).keySet().size(), equalTo(2));
		assertThat(statisticsLog.filter(0, new String[] {},false, true).keySet().size(), equalTo(0));
	}

	@Test
	public void itReturnsNullForIndexesOutsideDataPointSpan() {
		mockStatic(System.class);

		StatisticsLog statisticsLog = new StatisticsLog();

		when(System.currentTimeMillis()).thenReturn(1455224271177L);
		statisticsLog.prepareForUpdate("cache1", 1000L);

		when(System.currentTimeMillis()).thenReturn(1455224271177L + 100L);
		statisticsLog.prepareForUpdate("cache1", 1000L);

		when(System.currentTimeMillis()).thenReturn(1455224271177L + 200L);
		statisticsLog.prepareForUpdate("cache1", 1000L);

		statisticsLog.putDataPoint("kbps","1234.5");

		when(System.currentTimeMillis()).thenReturn(1455224271177L + 300L);
		statisticsLog.prepareForUpdate("cache1", 1000L);
		statisticsLog.putDataPoint("kbps","1234.5");

		assertThat(statisticsLog.getValue("kbps", 1), nullValue());
		assertThat(statisticsLog.getValue("kbps", 2), nullValue());
		assertThat(statisticsLog.getValue("kbps", 3), equalTo("1234.5"));
		assertThat(statisticsLog.getValue("kbps", 4), equalTo("1234.5"));
		assertThat(statisticsLog.getValue("kbps", 5), nullValue());
	}

	@Test
	public void itIsThreadSafe() throws Exception {
		StatisticsLog statisticsLog = new StatisticsLog();

		final CyclicBarrier cyclicBarrier = new CyclicBarrier(4);

		Publisher publisher = new Publisher(statisticsLog, cyclicBarrier);
		Prepper prepper = new Prepper(statisticsLog, cyclicBarrier);
		TimeGetter timeGetter = new TimeGetter(statisticsLog, cyclicBarrier);

		Thread publisherThread = new Thread(publisher);
		Thread prepperThread = new Thread(prepper);
		Thread getterThread = new Thread(timeGetter);

		getterThread.start();
		prepperThread.start();
		publisherThread.start();

		cyclicBarrier.await();
		List<Integer> exceptions = new ArrayList<Integer>();
		exceptions.add(timeGetter.exceptionCount.get());
		exceptions.add(publisher.exceptionCount.get());
		exceptions.add(prepper.exceptionCount.get());
		assertThat(exceptions, contains(0, 0, 0));
	}
}

class Publisher implements Runnable {
	static String[] keys = new String[] {"one", "two", "three", "four", "five", "six", "seven"};
	static String[] values = new String[] {"aardvark", "bear", "crocodile", "dog", "elephant", "fox", "gorilla"};
	StatisticsLog statisticsLog;
	CyclicBarrier cyclicBarrier;
	AtomicInteger exceptionCount = new AtomicInteger(0);

	public Publisher(StatisticsLog statisticsLog, CyclicBarrier cyclicBarrier) {
		this.statisticsLog = statisticsLog;
		this.cyclicBarrier = cyclicBarrier;
	}

	@Override
	public void run() {
		for (int i = 0; i < 20000; i++) {
			try {
				statisticsLog.putDataPoint(keys[i%keys.length], values[i%values.length]);
			} catch (Throwable t) {
				exceptionCount.incrementAndGet();
			}
		}
		try {
			cyclicBarrier.await();
		} catch (Exception e) {
			e.printStackTrace();
		}
	}
}

class Prepper implements Runnable {
	StatisticsLog statisticsLog;
	CyclicBarrier cyclicBarrier;
	AtomicInteger exceptionCount = new AtomicInteger(0);

	public Prepper(StatisticsLog statisticsLog, CyclicBarrier cyclicBarrier) {
		this.statisticsLog = statisticsLog;
		this.cyclicBarrier = cyclicBarrier;
	}

	@Override
	public void run() {
		for (int i = 0; i < 500; i++) {
			try {
				statisticsLog.prepareForUpdate("state id", 5*60*1000);
			} catch (Throwable t) {
				t.printStackTrace();
				exceptionCount.incrementAndGet();
			}
		}

		try {
			cyclicBarrier.await();
		} catch (Exception e) {
			e.printStackTrace();
		}
	}
}

// This loosely mimics DeliveryServiceStateRegistry.createStati that is having concurrency problems
class TimeGetter implements Runnable {
	StatisticsLog statisticsLog;
	CyclicBarrier cyclicBarrier;
	AtomicInteger exceptionCount = new AtomicInteger(0);
	DeliveryServiceStateRegistry deliveryServiceStateRegistry = DeliveryServiceStateRegistry.getInstance();

	public TimeGetter(StatisticsLog statisticsLog, CyclicBarrier cyclicBarrier) {
		this.statisticsLog = statisticsLog;
		this.cyclicBarrier = cyclicBarrier;
	}

	@Override
	public void run() {
		for (int i = 0; i < 500; i++) {
			try {
				Deque<DataPoint> dataPoints = statisticsLog.get(Publisher.keys[i % Publisher.keys.length]);
				if (dataPoints != null && !dataPoints.isEmpty()) {
					long lastIndex = dataPoints.getLast().getIndex();
					lastIndex = deliveryServiceStateRegistry.getLastGoodIndex(dataPoints, lastIndex);
					if (lastIndex < 0) {
						continue;
					}

					statisticsLog.getTime(lastIndex);

				}
			} catch (Throwable t) {
				t.printStackTrace();
				exceptionCount.incrementAndGet();
			}
		}
		try {
			cyclicBarrier.await();
		} catch (Exception e) {
			e.printStackTrace();
		}
	}
}
