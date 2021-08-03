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

ALTER TABLE public.server_server_capability DROP CONSTRAINT fk_server_capability;
ALTER TABLE public.deliveryservices_required_capability DROP CONSTRAINT fk_required_capability;
ALTER TABLE public.server_server_capability ADD CONSTRAINT fk_server_capability FOREIGN KEY (server_capability) REFERENCES server_capability(name) ON DELETE RESTRICT;
ALTER TABLE public.deliveryservices_required_capability ADD CONSTRAINT fk_required_capability FOREIGN KEY (required_capability) REFERENCES server_capability(name) ON DELETE RESTRICT;
