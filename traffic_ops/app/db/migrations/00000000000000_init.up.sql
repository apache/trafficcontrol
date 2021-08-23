/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with this
 * work for additional information regarding copyright ownership.  The ASF
 * licenses this file to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.  See the
 * License for the specific language governing permissions and limitations under
 * the License.
 */

CREATE OR REPLACE FUNCTION fast_forward_schema_migrations_version()
	RETURNS TRIGGER
	LANGUAGE PLPGSQL
AS $$
BEGIN
	IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'goose_db_version') THEN
		UPDATE schema_migrations SET "version" = (SELECT version_id FROM goose_db_version
			ORDER BY TSTAMP DESC
			LIMIT 1);
		ALTER TABLE goose_db_version
			RENAME TO goose_db_version_unused;
	END IF;
	DROP TRIGGER fast_forward_schema_migrations_trigger ON schema_migrations;
	DROP FUNCTION fast_forward_schema_migrations_version;
	RETURN NULL;
END$$;

CREATE TRIGGER fast_forward_schema_migrations_trigger
	AFTER UPDATE
	ON schema_migrations
	EXECUTE PROCEDURE fast_forward_schema_migrations_version();
