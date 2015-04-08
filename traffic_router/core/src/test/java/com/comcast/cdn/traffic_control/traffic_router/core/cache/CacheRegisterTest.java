/*
 * Copyright 2015 Comcast Cable Communications Management, LLC
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

package com.comcast.cdn.traffic_control.traffic_router.core.cache;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertFalse;
import static org.junit.Assert.assertNotNull;
import static org.junit.Assert.assertTrue;

import java.util.HashSet;
import java.util.Set;

import org.junit.Before;
import org.junit.Test;

import com.comcast.cdn.traffic_control.traffic_router.core.cache.CacheRegister;

public class CacheRegisterTest {

    private CacheRegister cacheRegister;

    @Before
    public void setUp() throws Exception {
        cacheRegister = new CacheRegister();
    }


    @Test
    public void testSetCacheLocations() {
//        final Set<CacheLocation> update1 = new HashSet<CacheLocation>();
//        final Set<CacheLocation> update2 = new HashSet<CacheLocation>();
//
//        final CacheLocation loc1 = new CacheLocation("loc1", null, null);
//        final CacheLocation loc2 = new CacheLocation("loc2", null, null);
//        final Cache cache = new Cache("cache");
//        cache.setAdminStatus(AdminStatus.ONLINE);
//        loc2.addCache(cache);
//        final CacheLocation loc3 = new CacheLocation("loc3", null, null);
//
//        update1.add(loc1);
//        update1.add(loc2);
//        update1.add(loc3);
//
//        final CacheLocation loc4 = new CacheLocation("loc2", null, null);
//        final CacheLocation loc5 = new CacheLocation("loc3", null, null);
//        final CacheLocation loc6 = new CacheLocation("loc4", null, null);
//
//        update2.add(loc4);
//        update2.add(loc5);
//        update2.add(loc6);
//
//        cacheRegister.setConfiguredLocations(update1);
//        assertEquals(update1, cacheRegister.getCacheLocations());
//        assertNotNull(cacheRegister.getCacheLocation("loc2").getCache("cache"));
//
//        cacheRegister.setConfiguredLocations(update2);
//        assertEquals(update2, cacheRegister.getCacheLocations());
//        assertNotNull(cacheRegister.getCacheLocation("loc2").getCache("cache"));
    }

}
