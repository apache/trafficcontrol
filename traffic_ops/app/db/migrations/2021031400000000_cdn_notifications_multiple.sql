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

-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
DROP TABLE IF EXISTS cdn_notification;

CREATE TABLE cdn_notification (
    id BIGSERIAL PRIMARY KEY,
    cdn text NOT NULL,
    "user" text NOT NULL,
    notification text NOT NULL,
    last_updated timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT fk_notification_cdn FOREIGN KEY (cdn) REFERENCES cdn(name) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_notification_user FOREIGN KEY ("user") REFERENCES tm_user(username) ON DELETE CASCADE ON UPDATE CASCADE
);
DROP TRIGGER IF EXISTS on_update_current_timestamp ON cdn_notification;
CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON cdn_notification FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS cdn_notification;

CREATE TABLE cdn_notification (
    cdn text NOT NULL,
    "user" text NOT NULL,
    notification text,
    last_updated timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT pk_cdn_notification PRIMARY KEY (cdn),
    CONSTRAINT fk_notification_cdn FOREIGN KEY (cdn) REFERENCES cdn(name) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_notification_user FOREIGN KEY ("user") REFERENCES tm_user(username) ON DELETE CASCADE ON UPDATE CASCADE
);
DROP TRIGGER IF EXISTS on_update_current_timestamp ON cdn_notification;
CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON cdn_notification FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();

