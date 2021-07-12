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

package geolocation;

import org.junit.Assert;

import org.junit.Before;
import org.junit.Test;

import org.apache.traffic_control.traffic_router.geolocation.Geolocation;

public class GeolocationTest {

    @Before
    public void setUp() throws Exception {
    }

    @Test
    public void testGetDistanceFrom() {
        final Geolocation l1 = new Geolocation(0f, 0f);
        final Geolocation l2 = new Geolocation(.5f, .5f);
        final double expected = 78.6;
        final double actual = l1.getDistanceFrom(l2);
        Assert.assertEquals(expected, actual, 0.1);
    }

    @Test
    public void testGetDistanceFromEquator() {
        final Geolocation l1 = new Geolocation(1f, 0f);
        final Geolocation l2 = new Geolocation(-1f, 0f);
        final double expected = 222.4;
        final double actual = l1.getDistanceFrom(l2);
        Assert.assertEquals(expected, actual, 0.1);
    }

    @Test
    public void testGetDistanceFromIntlDateLine() {
        final Geolocation l1 = new Geolocation(0f, 179f);
        final Geolocation l2 = new Geolocation(0f, -179f);
        final double expected = 222.4;
        final double actual = l1.getDistanceFrom(l2);
        Assert.assertEquals(expected, actual, 0.1);
    }

    @Test
    public void testGetDistanceFromNull() {
        final Geolocation l1 = new Geolocation(0f, 1f);
        final Geolocation l2 = null;
        final double expected = Double.POSITIVE_INFINITY;
        final double actual = l1.getDistanceFrom(l2);
        Assert.assertEquals(expected, actual, 0.1);
    }

    @Test
    public void testGetDistanceFromPrimeMeridian() {
        final Geolocation l1 = new Geolocation(0f, 1f);
        final Geolocation l2 = new Geolocation(0f, -1f);
        final double expected = 222.4;
        final double actual = l1.getDistanceFrom(l2);
        Assert.assertEquals(expected, actual, 0.1);
    }

}
