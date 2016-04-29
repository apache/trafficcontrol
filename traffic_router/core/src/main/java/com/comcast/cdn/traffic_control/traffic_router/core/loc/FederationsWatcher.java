package com.comcast.cdn.traffic_control.traffic_router.core.loc;

import com.comcast.cdn.traffic_control.traffic_router.core.util.ProtectedFetcher;
import com.comcast.cdn.traffic_control.traffic_router.core.util.TrafficOpsUtils;
import org.apache.log4j.Logger;
import org.json.JSONException;
import org.json.JSONObject;

import java.io.File;
import java.io.FileReader;
import java.io.FileWriter;
import java.io.IOException;
import java.net.URL;
import java.util.concurrent.TimeUnit;

public class FederationsWatcher extends AbstractServiceUpdater {
    private static final Logger LOGGER = Logger.getLogger(FederationsWatcher.class);

    private URL authorizationURL;
    private String postData;
    private ProtectedFetcher fetcher;
    private TrafficOpsUtils trafficOpsUtils;
    private FederationRegistry federationRegistry;

    public void configure(final URL authorizationURL, final String postData, final URL federationsURL, final long pollingInterval, final int timeout) {
        if (authorizationURL.equals(this.authorizationURL) && postData.equals(this.postData) &&
            federationsURL.equals(federationsURL) && pollingInterval == getPollingInterval()) {
            return;
        }

        // avoid recreating the fetcher if possible
        if (!authorizationURL.equals(this.authorizationURL) || !postData.equals(this.postData)) {
            this.authorizationURL = authorizationURL;
            this.postData = postData;
            fetcher = new ProtectedFetcher(authorizationURL.toString(), postData, timeout);
        }

        setDataBaseURL(federationsURL.toString(), pollingInterval);
    }

    public void configure(final JSONObject config) {
        URL authUrl;
        String jsonData;
        URL fedsUrl = null;
        long interval = -1L;

        try {
            authUrl = new URL(trafficOpsUtils.getAuthUrl());
            jsonData = trafficOpsUtils.getAuthJSON().toString();
        } catch (Exception e) {
            LOGGER.warn("Failed to update URL for TrafficOps authorization, " +
                "check the api.auth.url, and the TrafficOps username and password configuration setting: " + e.getMessage());
            // All or nothing, don't allow the watcher to be halfway misconfigured
            authUrl = this.authorizationURL;
            jsonData = this.postData;
        }
        try{
            fedsUrl = new URL(trafficOpsUtils.getUrl("federationmapping.polling.url"));
        } catch (Exception e) {
            LOGGER.warn("Invalid Federation Polling URL, check the federationmapping.polling.url configuration: " + e.getMessage());
        }

        try {
            interval = config.getLong("federationmapping.polling.interval");
        } catch (JSONException e) {
            LOGGER.warn("Bad configuration value for federationmapping.polling.interval, ignoring configuration value and keeping it at "
                + interval + " " + TimeUnit.MILLISECONDS + ": " + e.getMessage());
            interval = getPollingInterval();
        }

        final int timeout = config.optInt("federationmapping.polling.timeout", 15 * 1000); // socket timeouts are in ms

        if (authUrl != null && jsonData != null && fedsUrl != null && interval != -1L) {
            configure(authUrl, jsonData, fedsUrl, interval, timeout);
        }
    }

    @Override
    public boolean loadDatabase() throws IOException, org.apache.wicket.ajax.json.JSONException {
        final File existingDB = new File(databasesDirectory, databaseName);

        if (!existingDB.exists() || !existingDB.canRead()) {
            return false;
        }

        final char[] jsonData = new char[(int) existingDB.length()];
        final FileReader reader = new FileReader(existingDB);

        try {
            reader.read(jsonData);
        } finally {
            reader.close();
        }

        final String json = new String(jsonData);

        federationRegistry.setFederations(new FederationsBuilder().fromJSON(json));

        setLoaded(true);
        return true;
    }

    @Override
    protected File downloadDatabase(final String url, final File existingDb) {
        if (fetcher == null) {
            LOGGER.warn("Waiting for federations configuration to be processed, unable download federations");
            return null;
        }

        String jsonData = null;
        try {
            jsonData = fetcher.fetchIfModifiedSince(url, existingDb.lastModified());
        }
        catch (IOException e) {
            LOGGER.warn("Failed to fetch federations mapping from '" + url + "': " + e.getMessage());
        }

        if (jsonData == null) {
            return existingDb;
        }

        File databaseFile = null;
        FileWriter fw;
        try {
            databaseFile = File.createTempFile(tmpPrefix, tmpSuffix);
            fw = new FileWriter(databaseFile);
            fw.write(jsonData);
            fw.flush();
            fw.close();
        }
        catch (IOException e) {
            LOGGER.warn("Failed to create federations mapping file from data received from '" + url + "'");
        }

        return databaseFile;
    }

    public void setFederationRegistry(final FederationRegistry federationRegistry) {
        this.federationRegistry = federationRegistry;
    }

    public void setTrafficOpsUtils(final TrafficOpsUtils trafficOpsUtils) {
        this.trafficOpsUtils = trafficOpsUtils;
    }
}