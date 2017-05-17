INSERT INTO profile (name, description, type) VALUES ('TS_PROFILE','Traffic Stats','INFLUXDB_PROFILE') ON CONFLICT (name) DO NOTHING;
