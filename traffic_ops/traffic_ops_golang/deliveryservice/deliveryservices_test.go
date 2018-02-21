package deliveryservice

/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

import (
	"testing"

	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/api"
)

func getInterface(req *TODeliveryService) interface{} {
	return req
}

func TestInterfaces(t *testing.T) {
	r := getInterface(&TODeliveryService{})

	if _, ok := r.(api.Creator); !ok {
		t.Errorf("DeliveryService must be Creator")
	}
	// TODO: uncomment when ds.Reader interface is implemented
	/*if _, ok := r.(api.Reader); !ok {
		t.Errorf("DeliveryService must be Reader")
	}*/
	if _, ok := r.(api.Updater); !ok {
		t.Errorf("DeliveryService must be Updater")
	}
	if _, ok := r.(api.Deleter); !ok {
		t.Errorf("DeliveryService must be Deleter")
	}
	if _, ok := r.(api.Identifier); !ok {
		t.Errorf("DeliveryService must be Identifier")
	}
	if _, ok := r.(api.Tenantable); !ok {
		t.Errorf("DeliveryService must be Tenantable")
	}
}
