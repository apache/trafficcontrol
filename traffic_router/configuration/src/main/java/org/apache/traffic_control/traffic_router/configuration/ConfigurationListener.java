/*
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package org.apache.traffic_control.traffic_router.configuration;

public interface ConfigurationListener {

    /**
     * Called when the configuration has changed.
     * 
     * @param username The username of the user triggering the change.
     */
    void configurationChanged()

    /**
     * Called when an error occurs while processing the configuration change.
     * 
     * @param username The username of the user triggering the change.
     * @param error    The error message describing the issue.
     */
    void configurationError();

    /**
     * Called when the configuration change is successfully applied.
     * 
     * @param username The username of the user triggering the change.
     */
    void configurationApplied();

}
