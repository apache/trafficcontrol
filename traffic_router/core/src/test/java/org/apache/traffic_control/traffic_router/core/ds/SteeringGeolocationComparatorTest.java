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

package org.apache.traffic_control.traffic_router.core.ds;

import org.apache.traffic_control.traffic_router.core.edge.Cache;
import org.apache.traffic_control.traffic_router.geolocation.Geolocation;

import org.junit.Before;
import org.junit.Test;

import static org.junit.Assert.assertEquals;

public class SteeringGeolocationComparatorTest {
    /*
    This test class assumes some knowledge of United States geography. For reference,
    here are the rough distances looking at a map from left to right:

    Seattle <--- 1300 mi ---> Denver <--- 2000 mi ---> Boston

    */

    private Geolocation seattleGeolocation;
    private Geolocation denverGeolocation;
    private Geolocation bostonGeolocation;

    private Cache seattleCache;
    private Cache denverCache;
    private Cache bostonCache;

    private SteeringTarget seattleTarget;
    private SteeringTarget seattleTarget2;
    private SteeringTarget denverTarget;
    private SteeringTarget bostonTarget;

    private SteeringResult seattleResult;
    private SteeringResult seattleResult2;
    private SteeringResult denverResult;
    private SteeringResult bostonResult;

    private SteeringGeolocationComparator seattleComparator;
    private SteeringGeolocationComparator denverComparator;
    private SteeringGeolocationComparator bostonComparator;

    @Before
    public void before() {
        seattleGeolocation = new Geolocation(47.0, -122.0);
        denverGeolocation = new Geolocation(39.0, -104.0);
        bostonGeolocation = new Geolocation(42.0, -71.0);

        seattleCache = new Cache("seattle-id", "seattle-hash-id", 1, seattleGeolocation);
        denverCache = new Cache("denver-id", "denver-hash-id", 1, denverGeolocation);
        bostonCache = new Cache("boston-id", "boston-hash-id", 1, bostonGeolocation);

        seattleTarget = new SteeringTarget();
        seattleTarget.setGeolocation(seattleGeolocation);
        seattleResult = new SteeringResult(seattleTarget, null);
        seattleResult.setCache(seattleCache);

        seattleTarget2 = new SteeringTarget();
        seattleTarget2.setGeolocation(seattleGeolocation);
        seattleResult2 = new SteeringResult(seattleTarget2, null);
        seattleResult2.setCache(seattleCache);

        denverTarget = new SteeringTarget();
        denverTarget.setGeolocation(denverGeolocation);
        denverResult = new SteeringResult(denverTarget, null);
        denverResult.setCache(denverCache);

        bostonTarget = new SteeringTarget();
        bostonTarget.setGeolocation(bostonGeolocation);
        bostonResult = new SteeringResult(bostonTarget, null);
        bostonResult.setCache(bostonCache);

        seattleComparator = new SteeringGeolocationComparator(seattleGeolocation);
        denverComparator = new SteeringGeolocationComparator(denverGeolocation);
        bostonComparator = new SteeringGeolocationComparator(bostonGeolocation);

    }

    @Test
    public void testLeftNullOriginGeo() {
        seattleResult.getSteeringTarget().setGeolocation(null);
        assertEquals(1, seattleComparator.compare(seattleResult, bostonResult));
    }

    @Test
    public void testRightNullOriginGeo() {
        denverResult.getSteeringTarget().setGeolocation(null);
        assertEquals(-1, seattleComparator.compare(seattleResult, denverResult));
    }

    @Test
    public void testBothNullOriginGeo() {
        seattleResult.getSteeringTarget().setGeolocation(null);
        denverResult.getSteeringTarget().setGeolocation(null);
        assertEquals(0, seattleComparator.compare(seattleResult, denverResult));
    }

    @Test
    public void testSameCacheAndOriginGeo() {
        assertEquals(0, seattleComparator.compare(seattleResult, seattleResult));
    }

    @Test
    public void testSameCacheAndOriginGeoWithGeoOrder() {
        seattleTarget.setGeoOrder(1);
        seattleTarget2.setGeoOrder(2);
        assertEquals(-1, seattleComparator.compare(seattleResult, seattleResult2));
        assertEquals(1, seattleComparator.compare(seattleResult2, seattleResult));
        seattleTarget.setGeoOrder(2);
        assertEquals(0, seattleComparator.compare(seattleResult, seattleResult2));
    }

    @Test
    public void testDifferentCacheAndOriginGeo() {
        assertEquals(-1, seattleComparator.compare(seattleResult, denverResult));
        assertEquals(-1, seattleComparator.compare(denverResult, bostonResult));
        assertEquals(1, seattleComparator.compare(bostonResult, seattleResult));
    }

    @Test
    public void testCacheGeoDifferentFromOriginGeo() {
        seattleResult.setCache(denverCache);
        // seattle -> denver -> seattle || seattle -> denver
        assertEquals(1, seattleComparator.compare(seattleResult, denverResult));
        // denver -> seattle || denver
        assertEquals(1, denverComparator.compare(seattleResult, denverResult));
        // boston -> denver -> seattle || boston -> denver
        assertEquals(1, bostonComparator.compare(seattleResult, denverResult));
        // seattle -> denver -> seattle || seattle -> boston
        assertEquals(-1, seattleComparator.compare(seattleResult, bostonResult));
        seattleResult.setCache(bostonCache);
        bostonResult.setCache(denverCache);
        // seattle -> boston -> seattle || seattle -> denver -> boston
        assertEquals(1, seattleComparator.compare(seattleResult, bostonResult));
    }

}
