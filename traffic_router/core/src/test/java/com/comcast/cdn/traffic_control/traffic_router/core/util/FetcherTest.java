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

package com.comcast.cdn.traffic_control.traffic_router.core.util;

import org.junit.Test;
import org.junit.runner.RunWith;
import org.powermock.core.classloader.annotations.PowerMockIgnore;
import org.powermock.core.classloader.annotations.PrepareForTest;
import org.powermock.core.classloader.annotations.SuppressStaticInitializationFor;
import org.powermock.modules.junit4.PowerMockRunner;

import java.io.BufferedReader;
import java.io.InputStream;
import java.io.InputStreamReader;
import java.net.HttpURLConnection;
import java.net.URL;

import static org.mockito.Mockito.*;
import static org.powermock.api.mockito.PowerMockito.mock;
import static org.powermock.api.mockito.PowerMockito.whenNew;

@RunWith(PowerMockRunner.class)
@PrepareForTest({Fetcher.class, URL.class, InputStreamReader.class})
@PowerMockIgnore({"javax.net.ssl.*", "javax.management.*"})
public class FetcherTest {

    @Test
    public void itChecksIfDataHasChangedSinceLastFetch() throws Exception {
        InputStream inputStream = mock(InputStream.class);

        InputStreamReader inputStreamReader = mock(InputStreamReader.class);
        whenNew(InputStreamReader.class).withArguments(inputStream).thenReturn(inputStreamReader);

        BufferedReader bufferedReader = mock(BufferedReader.class);
        when(bufferedReader.readLine()).thenReturn(null);

        whenNew(BufferedReader.class).withArguments(inputStreamReader).thenReturn(bufferedReader);

        HttpURLConnection httpURLConnection = mock(HttpURLConnection.class);
        when(httpURLConnection.getInputStream()).thenReturn(inputStream);

        URL url = mock(URL.class);
        when(url.openConnection()).thenReturn(httpURLConnection);
        whenNew(URL.class).withArguments("http://www.example.com").thenReturn(url);

        Fetcher fetcher = new Fetcher();
        fetcher.fetchIfModifiedSince("http://www.example.com", 123456L);
        verify(httpURLConnection).setIfModifiedSince(123456L);
    }
}