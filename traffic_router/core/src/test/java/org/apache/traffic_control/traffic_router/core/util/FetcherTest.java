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

package org.apache.traffic_control.traffic_router.core.util;

import org.junit.Test;
import org.junit.runner.RunWith;
import org.powermock.core.classloader.annotations.PowerMockIgnore;
import org.powermock.core.classloader.annotations.PrepareForTest;
import org.powermock.modules.junit4.PowerMockRunner;


import java.io.BufferedReader;
import java.io.InputStream;
import java.io.InputStreamReader;
import java.net.HttpURLConnection;
import java.net.SocketTimeoutException;
import java.net.URL;
import java.net.URLConnection;

import static org.junit.Assert.assertTrue;
import static org.junit.Assert.assertEquals;
import static org.mockito.Mockito.doThrow;
import static org.mockito.Mockito.verify;
import static org.powermock.api.mockito.PowerMockito.method;
import static org.powermock.api.mockito.PowerMockito.mock;
import static org.powermock.api.mockito.PowerMockito.when;
import static org.powermock.api.mockito.PowerMockito.whenNew;
import static org.powermock.api.support.membermodification.MemberModifier.stub;

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
        stub(method(URL.class, "openConnection")).toReturn(httpURLConnection);
        whenNew(URL.class).withArguments("http://www.example.com").thenReturn(url);

        Fetcher fetcher = new Fetcher();
        fetcher.fetchIfModifiedSince("http://www.example.com", 123456L);
        verify(httpURLConnection).setIfModifiedSince(123456L);
    }

    @Test
    public void itChecksIfSocketTimeoutExceptionThrown() throws Exception {
        // Test that an IOException (ex connection-failure) is caught, re-thrown, and
        // that the embedded error is correctly encapsulated in the exception
        //
        final String mockedException = "Mocked Connection Failure";
        URL url = mock(URL.class);
        HttpURLConnection httpURLConnection = mock(HttpURLConnection.class);

        stub(method(URL.class, "openConnection")).toReturn(httpURLConnection);
        doThrow(new SocketTimeoutException(mockedException)).when(httpURLConnection).connect();
        whenNew(URL.class).withArguments("http://www.example.com").thenReturn(url);

        Fetcher fetcher = new Fetcher();
        try {
            fetcher.fetchIfModifiedSince("http://www.example.com", 123456L);
            assertTrue(false);
        } catch (SocketTimeoutException e) {
            assertEquals(e.toString(), "java.net.SocketTimeoutException: " + mockedException);
        }
    }

    @Test
    public void itChecksIfOtherThrown() throws Exception {
        // Test that other exceptions (ex CastClassException) is caught and
        // squelched. This matches existing functionality.
        //
        // This test relies on the CastClassException which is thrown when "connection"
        // is cast to type "HttpURLConnection". This will abort execution of
        // Fetcher::getConnection() during the cast of "connection" to "http". This will
        // result in the callstack unwinding normally with null return codes.
        //
        URLConnection urlConnection = mock(URLConnection.class);
        stub(method(URL.class, "openConnection")).toReturn(urlConnection);

        URL url = mock(URL.class);
        Fetcher fetcher = new Fetcher();
        assertEquals(fetcher.fetchIfModifiedSince("http://www.example.com", 123456L), null);
    }
}
