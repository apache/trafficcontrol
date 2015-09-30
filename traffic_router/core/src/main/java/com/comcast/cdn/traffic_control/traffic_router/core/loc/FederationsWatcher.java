package com.comcast.cdn.traffic_control.traffic_router.core.loc;

import com.comcast.cdn.traffic_control.traffic_router.core.util.ProtectedFetcher;
import org.json.JSONException;
import org.json.JSONObject;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.*;
import java.net.URL;
import java.util.List;

public class FederationsWatcher extends AbstractServiceUpdater {
    private static final Logger LOGGER = LoggerFactory.getLogger(FederationsWatcher.class);

    private URL authorizationURL;
    private String postData;
    private URL federationsURL;
    private long pollingInterval;
    private ProtectedFetcher fetcher;

    private List<Federation> federations;

    public void configure(final URL authorizationURL, final String postData, final URL federationsURL, final long pollingInterval) {

        if (authorizationURL.equals(this.authorizationURL) && postData.equals(this.postData) &&
            federationsURL.equals(this.federationsURL) && pollingInterval == this.pollingInterval) {
            return;
        }

        this.federationsURL = federationsURL;
        this.pollingInterval = pollingInterval;

        // avoid recreating the fetcher if possible
        if (!authorizationURL.equals(this.authorizationURL) || !postData.equals(this.postData)) {
            this.authorizationURL = authorizationURL;
            this.postData = postData;
            fetcher = new ProtectedFetcher(authorizationURL.toString(), postData, 0);
        }
    }

    public void configure(final JSONObject config) {
        URL authUrl;
        String jsonData;
        URL fedsUrl;
        long interval = -1L;

        try {
            authUrl = new URL(config.getString("keystore.auth.url"));
            final String username = config.getString("dns.sec.keyserver.username");
            final String password = config.getString("dns.sec.keyserver.password");
            jsonData = "{u:\"" + username + "\",p:\"" + password + "\"}";
        } catch (Exception e) {
            LOGGER.warn("Failed Getting Configuration for ProtectedFetcher for FederationsWatcher: " + e.getMessage());
            // All or nothing, don't allow the watcher to be halfway misconfigured
            authUrl = this.authorizationURL;
            jsonData = this.postData;
        }

        try {
            fedsUrl = new URL(config.getString("federationmapping.polling.federationsURL"));
        } catch (Exception e) {
            LOGGER.warn("Failed Getting Configuration for FederationsWatcher URL: " + e.getMessage());
            fedsUrl = this.federationsURL;
        }

        try {
            interval = config.getLong("federationmapping.polling.interval");
        } catch (JSONException e) {
            LOGGER.warn("Failed getting configuration for FederationsWatcher Polling Interval " + e.getMessage());
            interval = this.pollingInterval;
        }

        if (authUrl != null && jsonData != null && fedsUrl != null && interval != -1L) {
            configure(authUrl, jsonData, fedsUrl, interval);
        }
    }

    public List<Federation> getFederations() {
        return federations;
    }

    @Override
    public void verifyDatabase(final File dbFile) throws IOException {

    }

    @Override
    public boolean loadDatabase() throws IOException, org.apache.wicket.ajax.json.JSONException {
        final File existingDB = new File(databaseLocation);

        if (!existingDB.exists() || !existingDB.canRead()) {
            return false;
        }

        final char[] jsonData = new char[(int) existingDB.length()];
        new FileReader(existingDB).read(jsonData);
        final String json = new String(jsonData);

        federations = new FederationsBuilder().fromJSON(json);

        setLoaded(true);
        return true;
    }

    @Override
    protected File downloadDatabase(final String url, final File existingDb) {
        final String jsonData = null;

        try {
            fetcher.fetchIfModifiedSince(url, existingDb.lastModified());
        }
        catch (IOException e) {
            LOGGER.warn("Failed to fetch federations mapping from '" + url + "': ", e.getMessage());
        }

        if (jsonData == null) {
            return existingDb;
        }

        File databaseFile = null;
        try {
            databaseFile = File.createTempFile(tmpPrefix, tmpSuffix);
            new FileWriter(databaseFile).write(jsonData);
        }
        catch (IOException e) {
            LOGGER.warn("Failed to create federations mapping file from data received from '" + url + "'");
        }

        return databaseFile;
    }
}