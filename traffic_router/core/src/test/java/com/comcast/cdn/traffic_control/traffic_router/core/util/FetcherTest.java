package com.comcast.cdn.traffic_control.traffic_router.core.util;

import org.junit.Test;
import org.junit.runner.RunWith;
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
@PrepareForTest({URL.class, InputStreamReader.class})
@SuppressStaticInitializationFor("com.comcast.cdn.traffic_control.traffic_router.core.util.Fetcher")
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