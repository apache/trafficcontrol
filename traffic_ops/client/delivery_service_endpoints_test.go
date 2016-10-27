/*
   Copyright 2015 Comcast Cable Communications Management, LLC

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package client

import (
	"testing"

	"github.com/jheitz200/test_helper"
)

func TestDeliveryServicesEp(t *testing.T) {
	testHelper.Context(t, "Given the need to test that DeliveryServices uses the correct URL")

	ep := deliveryServicesEp()
	expected := "/api/1.2/deliveryservices.json"
	if ep != expected {
		testHelper.Error(t, "Should get back %s for \"deliveryServicesEp\", got: %s", expected, ep)
	} else {
		testHelper.Success(t, "Should be able to get the correct delivery services endpoint")
	}
}

func TestDeliveryServiceEp(t *testing.T) {
	testHelper.Context(t, "Given the need to test that DeliveryService uses the correct URL")

	ep := deliveryServiceEp("123")
	expected := "/api/1.2/deliveryservices/123.json"
	if ep != expected {
		testHelper.Error(t, "Should get back %s for \"deliveryServiceEp\", got: %s", expected, ep)
	} else {
		testHelper.Success(t, "Should be able to get the correct delivery service endpoint")
	}
}

func TestDeliveryServiceStateEp(t *testing.T) {
	testHelper.Context(t, "Given the need to test that DeliveryServiceState uses the correct URL")

	ep := deliveryServiceStateEp("123")
	expected := "/api/1.2/deliveryservices/123/state.json"
	if ep != expected {
		testHelper.Error(t, "Should get back %s for \"deliveryServiceStateEp\", got: %s", expected, ep)
	} else {
		testHelper.Success(t, "Should be able to get the correct delivery service state endpoint")
	}
}

func TestDeliveryServiceHealthEp(t *testing.T) {
	testHelper.Context(t, "Given the need to test that DeliveryServiceHealth uses the correct URL")

	ep := deliveryServiceHealthEp("123")
	expected := "/api/1.2/deliveryservices/123/health.json"
	if ep != expected {
		testHelper.Error(t, "Should get back %s for \"deliveryServiceHealthEp\", got: %s", expected, ep)
	} else {
		testHelper.Success(t, "Should be able to get the correct delivery service health endpoint")
	}
}

func TestDeliveryServiceCapacityEp(t *testing.T) {
	testHelper.Context(t, "Given the need to test that DeliveryServiceCapacity uses the correct URL")

	ep := deliveryServiceCapacityEp("123")
	expected := "/api/1.2/deliveryservices/123/capacity.json"
	if ep != expected {
		testHelper.Error(t, "Should get back %s for \"deliveryServiceCapacityEp\", got: %s", expected, ep)
	} else {
		testHelper.Success(t, "Should be able to get the correct delivery service capacity endpoint")
	}
}

func TestDeliveryServiceRoutingEp(t *testing.T) {
	testHelper.Context(t, "Given the need to test that DeliveryServiceRouting uses the correct URL")

	ep := deliveryServiceRoutingEp("123")
	expected := "/api/1.2/deliveryservices/123/routing.json"
	if ep != expected {
		testHelper.Error(t, "Should get back %s for \"deliveryServiceRoutingEp\", got: %s", expected, ep)
	} else {
		testHelper.Success(t, "Should be able to get the correct delivery service routing endpoint")
	}
}

func TestDeliveryServiceServerEp(t *testing.T) {
	testHelper.Context(t, "Given the need to test that DeliveryServiceServer uses the correct URL")

	ep := deliveryServiceServerEp("1", "2")
	expected := "/api/1.2/deliveryserviceserver.json?page=1&limit=2"
	if ep != expected {
		testHelper.Error(t, "Should get back %s for \"deliveryServiceServerEp\", got: %s", expected, ep)
	} else {
		testHelper.Success(t, "Should be able to get the correct delivery service server endpoint")
	}
}
