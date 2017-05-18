INSERT INTO profile (name, description, type) VALUES ('TV_PROFILE','Traffic Vault','RIAK_PROFILE') ON CONFLICT (name) DO NOTHING;
