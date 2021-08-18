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

package org.apache.traffic_control.traffic_router.core.router;

import static org.junit.Assert.assertEquals;

import org.junit.Test;

import org.apache.traffic_control.traffic_router.core.edge.CacheLocation;
import org.apache.traffic_control.traffic_router.geolocation.Geolocation;

public class CacheLocationComparatorTest {

    @Test
    public void testCompareBothLocEqual() {
        final LocationComparator comparator = new LocationComparator(new Geolocation(1f, 1f));
        final CacheLocation loc1 = new CacheLocation("loc1", new Geolocation(0f, 0f));
        final CacheLocation loc2 = new CacheLocation("loc2", new Geolocation(0f, 0f));

        assertEquals(0, comparator.compare(loc1, loc2));
        assertEquals(0, comparator.compare(loc2, loc1));
    }

    @Test
    public void testCompareBothLocNull() {
        final LocationComparator comparator = new LocationComparator(new Geolocation(1f, 1f));
        final CacheLocation loc1 = new CacheLocation("loc1", null);
        final CacheLocation loc2 = new CacheLocation("loc2", null);

        assertEquals(0, comparator.compare(loc1, loc2));
        assertEquals(0, comparator.compare(loc2, loc1));
    }

    @Test
    public void testCompareLocsDifferent() {
        final LocationComparator comparator = new LocationComparator(new Geolocation(1f, 1f));
        final CacheLocation loc1 = new CacheLocation("loc1", new Geolocation(1f, 1f));
        final CacheLocation loc2 = new CacheLocation("loc2", new Geolocation(0f, 0f));

        assertEquals(-1, comparator.compare(loc1, loc2));
        assertEquals(1, comparator.compare(loc2, loc1));
    }

    @Test
    public void testCompareOneLocNull() {
        final LocationComparator comparator = new LocationComparator(new Geolocation(1f, 1f));
        final CacheLocation loc1 = new CacheLocation("loc1", new Geolocation(0f, 0f));
        final CacheLocation loc2 = new CacheLocation("loc2", null);

        assertEquals(-1, comparator.compare(loc1, loc2));
        assertEquals(1, comparator.compare(loc2, loc1));
    }

}
