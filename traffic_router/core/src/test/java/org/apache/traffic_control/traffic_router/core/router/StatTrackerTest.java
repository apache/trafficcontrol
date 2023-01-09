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

import org.junit.Test;
import java.util.HashMap;
import java.util.Map;
import static org.powermock.api.mockito.PowerMockito.spy;
import static org.powermock.api.mockito.PowerMockito.when;
import static org.springframework.test.util.AssertionErrors.assertEquals;

public class StatTrackerTest {
    @Test
    public void testIncTally() {
        StatTracker tracker = spy(new StatTracker());
        StatTracker.Tallies tallies = new StatTracker.Tallies();
        StatTracker.Track track = StatTracker.getTrack();

        tallies.setCzCount(Long.MAX_VALUE);

        Map<String, StatTracker.Tallies> map = new HashMap<>();
        map.put("blah", tallies);
        when(tracker.getDnsMap()).thenReturn(map);

        track.setRouteType(StatTracker.Track.RouteType.DNS, "blah");
        track.setResult(StatTracker.Track.ResultType.CZ);
        tracker.saveTrack(track);
        assertEquals("expected czCount to be max long value but got " + tallies.getCzCount(), Long.MAX_VALUE, tallies.getCzCount());
    }
}
