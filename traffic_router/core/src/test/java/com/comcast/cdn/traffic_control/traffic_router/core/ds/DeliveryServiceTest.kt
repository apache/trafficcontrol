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
package com.comcast.cdn.traffic_control.traffic_router.core.ds

import kotlin.Throws
import java.lang.Exception
import com.comcast.cdn.traffic_control.traffic_router.core.request.HTTPRequest
import com.fasterxml.jackson.databind.ObjectMapper
import org.hamcrest.MatcherAssert
import org.hamcrest.Matchers
import org.junit.Test
import org.powermock.reflect.Whitebox

class DeliveryServiceTest {
    @Test
    @Throws(Exception::class)
    fun itHandlesLackOfRequestHeaderNamesInJSON() {
        val mapper = ObjectMapper()
        val jsonStr = "{\"routingName\":\"edge\",\"coverageZoneOnly\":false}"
        val jsonConfiguration = mapper.readTree(jsonStr)
        val deliveryService = DeliveryService("a-delivery-service", jsonConfiguration)
        MatcherAssert.assertThat(deliveryService.requestHeaders.size, Matchers.equalTo(0))
    }

    @Test
    @Throws(Exception::class)
    fun itHandlesLackOfConsistentHashQueryParamsInJSON() {
        val mapper = ObjectMapper()
        val json = mapper.readTree("{\"routingName\":\"edge\",\"coverageZoneOnly\":false}")
        val d = DeliveryService("test", json)
        assert(d.consistentHashQueryParams != null)
        assert(d.consistentHashQueryParams.size == 0)
    }

    @Test
    @Throws(Exception::class)
    fun itHandlesDuplicatesInConsistentHashQueryParams() {
        val mapper = ObjectMapper()
        val json =
            mapper.readTree("{\"routingName\":\"edge\",\"coverageZoneOnly\":false,\"consistentHashQueryParams\":[\"test\", \"quest\", \"test\"]}")
        val d = DeliveryService("test", json)
        assert(d.consistentHashQueryParams != null)
        assert(d.consistentHashQueryParams.size == 2)
        assert(d.consistentHashQueryParams.contains("test"))
        assert(d.consistentHashQueryParams.contains("quest"))
    }

    @Test
    @Throws(Exception::class)
    fun itExtractsQueryParams() {
        val json =
            ObjectMapper().readTree("{\"routingName\":\"edge\",\"coverageZoneOnly\":false,\"consistentHashQueryParams\":[\"test\", \"quest\"]}")
        val r = HTTPRequest()
        r.path = "/path1234/some_stream_name1234/some_other_info.m3u8"
        r.queryString = "test=value&foo=fizz&quest=oth%20ervalue&bar=buzz"
        assert(DeliveryService("test", json).extractSignificantQueryParams(r) == "quest=oth ervaluetest=value")
    }

    @Test
    @Throws(Exception::class)
    fun itConfiguresRequestHeadersFromJSON() {
        val mapper = ObjectMapper()
        val jsonStr =
            "{\"routingName\":\"edge\",\"requestHeaders\":[\"Cookie\",\"Cache-Control\",\"If-Modified-Since\",\"Content-Type\"],\"coverageZoneOnly\":false}"
        val jsonConfiguration = mapper.readTree(jsonStr)
        val deliveryService = DeliveryService("a-delivery-service", jsonConfiguration)
        MatcherAssert.assertThat(
            deliveryService.requestHeaders,
            Matchers.containsInAnyOrder("Cache-Control", "Cookie", "Content-Type", "If-Modified-Since")
        )
    }

    @Test
    @Throws(Exception::class)
    fun itAddsRequiredCapabilities() {
        val mapper = ObjectMapper()
        val jsonConfiguration =
            mapper.readTree("{\"requiredCapabilities\":[\"all-read\",\"all-write\",\"cdn-read\"],\"routingName\":\"edge\",\"coverageZoneOnly\":false}")
        val deliveryService = DeliveryService("has-required-capabilities", jsonConfiguration)
        MatcherAssert.assertThat(
            Whitebox.getInternalState(deliveryService, "requiredCapabilities"),
            Matchers.containsInAnyOrder("all-read", "all-write", "cdn-read")
        )
    }
}