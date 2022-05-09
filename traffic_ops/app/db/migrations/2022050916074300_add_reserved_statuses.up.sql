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

INSERT INTO public.status ("name", description) VALUES ('ONLINE', 'Server is online.') ON CONFLICT ("name") DO NOTHING;
INSERT INTO public.status ("name", description) VALUES ('OFFLINE', 'Server is Offline. Not active in any configuration.') ON CONFLICT ("name") DO NOTHING;
INSERT INTO public.status ("name", description) VALUES ('REPORTED', 'Server is online and reported in the health protocol.') ON CONFLICT ("name") DO NOTHING;
INSERT INTO public.status ("name", description) VALUES ('ADMIN_DOWN', 'Sever is administrative down and does not receive traffic.') ON CONFLICT ("name") DO NOTHING;
INSERT INTO public.status ("name", description) VALUES ('PRE_PROD', 'Pre Production. Not active in any configuration.') ON CONFLICT ("name") DO NOTHING;
