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

package org.apache.traffic_control.traffic_router.core.loc;

import java.util.ArrayList;
import java.util.List;

public class RegionalGeoDsvc {
    private final String id;
    private final List<RegionalGeoRule> urlRules = new ArrayList<RegionalGeoRule>();

    public RegionalGeoDsvc(final String id) {
        this.id = id;
    }

    public String getId() {
        return id;
    }

    public void addRule(final RegionalGeoRule urlRule) {
        urlRules.add(urlRule);
    }

    public RegionalGeoRule matchRule(final String url) {
        for (final RegionalGeoRule rule : urlRules) {
            if (rule.matchesUrl(url)) {
                return rule;
            }
        }

        return null;
    }
}

