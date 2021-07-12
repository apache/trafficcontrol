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

package org.apache.traffic_control.traffic_router.core.dns;

import java.util.Arrays;

public class RRSIGCacheKey {
    final private byte[] privateKeyBytes;
    final private int algorithm;

    public RRSIGCacheKey(final byte[] privateKeyBytes, final int algorithm) {
        this.privateKeyBytes = privateKeyBytes;
        this.algorithm = algorithm;
    }

    @Override
    public boolean equals(final Object o) {
        if (this == o) {
            return true;
        }
        if (o == null || getClass() != o.getClass()) {
            return false;
        }
        final RRSIGCacheKey that = (RRSIGCacheKey) o;
        return algorithm == that.algorithm && Arrays.equals(privateKeyBytes, that.privateKeyBytes);
    }

    @Override
    public int hashCode() {
        int result = algorithm;
        result = 31 * result + Arrays.hashCode(privateKeyBytes);
        return result;
    }

}
