package comment

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
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/test"
)

func TestFuncs(t *testing.T) {
	if strings.Index(selectQuery(), "SELECT") != 0 {
		t.Errorf("expected selectQuery to start with SELECT")
	}
	if strings.Index(insertQuery(), "INSERT") != 0 {
		t.Errorf("expected insertQuery to start with INSERT")
	}
	if strings.Index(updateQuery(), "UPDATE") != 0 {
		t.Errorf("expected updateQuery to start with UPDATE")
	}
	if strings.Index(deleteQuery(), "DELETE") != 0 {
		t.Errorf("expected deleteQuery to start with DELETE")
	}

}
func TestInterfaces(t *testing.T) {
	var i interface{}
	i = &TODeliveryServiceRequestComment{}

	if _, ok := i.(api.Creator); !ok {
		t.Errorf("comment must be creator")
	}
	if _, ok := i.(api.Reader); !ok {
		t.Errorf("comment must be reader")
	}
	if _, ok := i.(api.Updater); !ok {
		t.Errorf("comment must be updater")
	}
	if _, ok := i.(api.Deleter); !ok {
		t.Errorf("comment must be deleter")
	}
	if _, ok := i.(api.Identifier); !ok {
		t.Errorf("comment must be Identifier")
	}
}

func TestValidate(t *testing.T) {
	c := TODeliveryServiceRequestComment{}
	err, _ := c.Validate()
	errs := util.JoinErrsStr(test.SortErrors(test.SplitErrors(err)))

	expectedErrs := util.JoinErrsStr([]error{
		errors.New(`'deliveryServiceRequestId' is required`),
		errors.New(`'value' is required`),
	})

	if !reflect.DeepEqual(expectedErrs, errs) {
		t.Errorf("expected %s, got %s", expectedErrs, errs)
	}

	v := "the comment value"
	d := 1
	c = TODeliveryServiceRequestComment{DeliveryServiceRequestCommentNullable: tc.DeliveryServiceRequestCommentNullable{DeliveryServiceRequestID: &d, Value: &v}}

	err, _ = c.Validate()
	if err != nil {
		t.Errorf("expected nil, got %s", err)
	}

}
