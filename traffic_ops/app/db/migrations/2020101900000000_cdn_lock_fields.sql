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

ALTER TABLE cdn
    ADD COLUMN locked_by text,
    ADD COLUMN locked_reason text,
    ADD CONSTRAINT cdn_locked_by_fkey FOREIGN KEY (locked_by) REFERENCES tm_user (username) ON UPDATE CASCADE ON DELETE RESTRICT;


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

ALTER TABLE cdn
    DROP COLUMN locked_by,
    DROP COLUMN locked_reason;
