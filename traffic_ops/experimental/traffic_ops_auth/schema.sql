CREATE TABLE IF NOT EXISTS traffic_ops_auth_users (
	username text PRIMARY KEY,
	hash text NOT NULL,
	salt text NOT NULL,
	role text NOT NULL,
	created timestamp without time zone DEFAULT now() NOT NULL
);
