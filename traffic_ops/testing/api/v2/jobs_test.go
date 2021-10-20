package v2

/*

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

import (
	"strings"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v6/lib/go-tc"
)

func TestJobs(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, DeliveryServices}, func() {
		CreateTestJobs(t)
		CreateTestInvalidationJobs(t)
		GetTestJobsQueryParams(t)
		GetTestJobs(t)
		GetTestInvalidationJobs(t)
	})
}

func CreateTestJobs(t *testing.T) {
	toDSes, _, err := TOSession.GetDeliveryServicesNullable()
	if err != nil {
		t.Fatalf("cannot GET DeliveryServices: %v - %v", err, toDSes)
	}

	for i, job := range testData.InvalidationJobs {
		job.StartTime = &tc.Time{
			Time:  time.Now().Add(time.Minute).UTC(),
			Valid: true,
		}
		testData.InvalidationJobs[i] = job
	}

	for _, job := range testData.InvalidationJobs {
		request := tc.InvalidationJobInput{
			DeliveryService: job.DeliveryService,
			Regex:           job.Regex,
			StartTime:       job.StartTime,
			TTL:             job.TTL,
		}
		_, _, err := TOSession.CreateInvalidationJob(request)
		if err != nil {
			t.Errorf("could not CREATE job: %v", err)
		}
	}
}

func CreateTestInvalidationJobs(t *testing.T) {
	toDSes, _, err := TOSession.GetDeliveryServicesNullable()
	if err != nil {
		t.Fatalf("cannot GET Delivery Services: %v - %v", err, toDSes)
	}
	dsNameIDs := map[string]int64{}
	for _, ds := range toDSes {
		dsNameIDs[*ds.XMLID] = int64(*ds.ID)
	}

	for _, job := range testData.InvalidationJobs {
		_, ok := dsNameIDs[(*job.DeliveryService).(string)]
		if !ok {
			t.Fatalf("can't create test data job: delivery service '%v' not found in Traffic Ops", job.DeliveryService)
		}
		if _, _, err := TOSession.CreateInvalidationJob(job); err != nil {
			t.Errorf("could not CREATE job: %v", err)
		}
	}
}

func GetTestJobsQueryParams(t *testing.T) {
	var xmlId interface{} = "ds2"
	toJobs, _, err := TOSession.GetInvalidationJobs(&xmlId, nil)
	if err != nil {
		t.Fatalf("error getting jobs: %v", err)
	}
	foundOne := false
	for _, j := range toJobs {
		if j.DeliveryService == nil {
			t.Error("expected: non-nil DeliveryService pointer, actual: nil")
		} else if *j.DeliveryService != "ds2" {
			t.Errorf("expected: DeliveryService == ds2, actual: DeliveryService == %s", *j.DeliveryService)
		} else {
			foundOne = true
		}
	}
	if !foundOne {
		t.Error("expected: to find at least one job with deliveryService == ds2, actual: found none")
	}
}

func GetTestJobs(t *testing.T) {
	toJobs, _, err := TOSession.GetInvalidationJobs(nil, nil)
	if err != nil {
		t.Fatalf("error getting jobs: %v", err)
	}

	toDSes, _, err := TOSession.GetDeliveryServicesNullable()
	if err != nil {
		t.Fatalf("cannot GET DeliveryServices: %v - %v", err, toDSes)
	}

	for i, testJob := range testData.InvalidationJobs {
		found := false
		if testJob.DeliveryService == nil {
			t.Errorf("test job (index %v) has nil delivery service", i)
			continue
		} else if testJob.Regex == nil {
			t.Errorf("test job (index %v) has nil regex", i)
			continue
		}
		for j, toJob := range toJobs {
			if toJob.DeliveryService == nil {
				t.Errorf("to job (index %v) has nil delivery service", j)
				continue
			} else if toJob.AssetURL == nil {
				t.Errorf("to job (index %v) has nil asset url", j)
				continue
			}
			if *toJob.DeliveryService != *testJob.DeliveryService {
				continue
			}
			if !strings.HasSuffix(*toJob.AssetURL, *testJob.Regex) {
				continue
			}
			toJobTime := toJob.StartTime.Round(time.Minute)
			testJobTime := testJob.StartTime.Round(time.Minute)
			if !toJobTime.Equal(testJobTime) {
				t.Errorf("test job ds %v regex %v start time expected '%+v' actual '%+v'", *testJob.DeliveryService, *testJob.Regex, testJobTime, toJobTime)
				continue
			}
			found = true
			break
		}
		if !found {
			t.Errorf("test job ds %v regex %v expected: exists, actual: not found", *testJob.DeliveryService, *testJob.Regex)
		}
	}
}

func GetTestInvalidationJobs(t *testing.T) {
	jobs, _, err := TOSession.GetInvalidationJobs(nil, nil)
	if err != nil {
		t.Fatalf("error getting invalidation jobs: %v", err)
	}

	toDSes, _, err := TOSession.GetDeliveryServicesNullable()
	if err != nil {
		t.Fatalf("cannot GET DeliveryServices: %v - %v", err, toDSes)
	}

	for _, ds := range toDSes {
		if *ds.ID <= 0 {
			t.Fatalf("Erroneous Delivery Service - has invalid ID: %+v", ds)
		}
	}

	for _, testJob := range testData.InvalidationJobs {
		found := false
		for _, toJob := range jobs {
			if *toJob.DeliveryService != (*testJob.DeliveryService).(string) {
				continue
			}
			if !strings.HasSuffix(*toJob.AssetURL, *testJob.Regex) {
				continue
			}
			if !toJob.StartTime.Round(time.Minute).Equal(testJob.StartTime.Round(time.Minute)) {
				t.Errorf("test invalidation job start time expected '%+v' actual '%+v'", testJob.StartTime, toJob.StartTime)
				continue
			}
			found = true
			break
		}
		if !found {
			t.Errorf("expected a test job %+v to exist, but it didn't", testJob)
		}
	}
}
