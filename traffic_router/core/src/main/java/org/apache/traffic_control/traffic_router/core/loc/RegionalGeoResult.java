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


public class RegionalGeoResult {
    public enum RegionalGeoResultType {
        ALLOWED, ALTERNATE_WITH_CACHE, ALTERNATE_WITHOUT_CACHE, DENIED
    }

    public final static int REGIONAL_GEO_DENIED_HTTP_CODE = 520;

    private String url;
    private int httpResponseCode;
    private RegionalGeoResultType resultType;
    private RegionalGeoRule.PostalsType ruleType;
    private String postal;
    private boolean usingFallbackConfig;
    private boolean allowedByWhiteList;

    public String getUrl() {
        return url;
    }

    public void setUrl(final String url) {
        this.url = url;
    }

    public int getHttpResponseCode() {
        return httpResponseCode;
    }

    public void setHttpResponseCode(final int rc) {
        this.httpResponseCode = rc;
    }

    public RegionalGeoResultType getType() {
        return resultType;
    }

    public void setType(final RegionalGeoResultType resultType) {
        this.resultType = resultType;
    }

    public RegionalGeoRule.PostalsType getRuleType() {
        return ruleType;
    }

    public void setRuleType(final RegionalGeoRule.PostalsType ruleType) {
        this.ruleType = ruleType;
    }

    public String getPostal() {
        return postal;
    }

    public void setPostal(final String postal) {
        this.postal = postal;
    }

    public void setUsingFallbackConfig(final boolean usingFallbackConfig) {
        this.usingFallbackConfig = usingFallbackConfig;
    }

    public void setAllowedByWhiteList(final boolean allowedByWhiteList) {
        this.allowedByWhiteList = allowedByWhiteList;
    }

    public String toString() {
        final StringBuilder sb = new StringBuilder();

        if (postal == null) {
            sb.append('-');
        } else {
            sb.append(postal);
        }
        sb.append(':');

        // allow:1; disallow:0
        if (resultType == RegionalGeoResultType.ALLOWED) {
            sb.append('1');
        } else {
            sb.append('0');
        }
        sb.append(':');

        // include rule: I, exclude rule: X, no rule matches: -
        if (resultType == RegionalGeoResultType.DENIED) {
            sb.append('-');
        } else {
            if (ruleType == null) {
                sb.append('-');
            } else if (ruleType == RegionalGeoRule.PostalsType.INCLUDE) {
                sb.append('I');
            } else {
                sb.append('X');
            }
        }
        sb.append(':');

        if (usingFallbackConfig) {
            sb.append('1');
        } else {
            sb.append('0');
        }
        sb.append(':');

        if (allowedByWhiteList) {
            sb.append('1');
        } else {
            sb.append('0');
        }

        return sb.toString();
    }

}

