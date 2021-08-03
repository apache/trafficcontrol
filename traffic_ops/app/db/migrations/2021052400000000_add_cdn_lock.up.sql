/*
	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at
		http://www.apache.org/licenses/LICENSE-2.0
	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

CREATE TABLE IF NOT EXISTS public.cdn_lock (
                                        username text NOT NULL,
                                        cdn text NOT NULL,
                                        message text,
                                        soft boolean NOT NULL DEFAULT TRUE,
                                        last_updated timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT pk_cdn_lock PRIMARY KEY ("cdn"),
    CONSTRAINT fk_lock_cdn FOREIGN KEY ("cdn") REFERENCES cdn(name) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_lock_username FOREIGN KEY ("username") REFERENCES tm_user(username) ON DELETE CASCADE ON UPDATE CASCADE
    );

DROP TRIGGER IF EXISTS on_update_current_timestamp ON public.cdn_lock;
CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON public.cdn_lock FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();
