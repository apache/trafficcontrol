package v3

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
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v6/lib/go-util"

	"github.com/apache/trafficcontrol/v6/lib/go-tc"
)

func TestJobs(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, DeliveryServices}, func() {
		CreateTestJobs(t)
		CreateTestInvalidationJobs(t)
		CreateTestInvalidJob(t)
		GetTestJobsQueryParams(t)
		GetTestJobs(t)
		GetTestInvalidationJobs(t)
		JobCollisionWarningTest(t)
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

func JobCollisionWarningTest(t *testing.T) {
	startTime := tc.Time{
		Time:  time.Now().Add(time.Hour),
		Valid: true,
	}
	firstJob := tc.InvalidationJobInput{
		DeliveryService: util.InterfacePtr(testData.DeliveryServices[0].XMLID),
		Regex:           util.StrPtr(`/\.*([A-Z]0?)`),
		TTL:             util.InterfacePtr(16),
		StartTime:       &startTime,
	}

	_, _, err := TOSession.CreateInvalidationJob(firstJob)
	if err != nil {
		t.Fatal(err)
	}

	newTime := tc.Time{
		Time:  startTime.Time.Add(time.Hour),
		Valid: true,
	}
	newJob := tc.InvalidationJobInput{
		DeliveryService: firstJob.DeliveryService,
		Regex:           firstJob.Regex,
		TTL:             firstJob.TTL,
		StartTime:       &newTime,
	}

	alerts, _, err := TOSession.CreateInvalidationJob(newJob)
	if err != nil {
		t.Fatalf("expected invalidation job create to succeed: %v", err)
	}

	if len(alerts.Alerts) != 2 {
		t.Fatalf("expected 2 alerts on creation, got %v", len(alerts.Alerts))
	}

	if alerts.Alerts[0].Level != tc.WarnLevel.String() {
		t.Fatalf("expected first alert to be a warning, got %v", alerts.Alerts[0].Level)
	}

	if !strings.Contains(alerts.Alerts[0].Text, *firstJob.Regex) {
		t.Fatalf("expected first alert to be about the first job, got: %v", alerts.Alerts[0].Text)
	}

	jobs, _, err := TOSession.GetInvalidationJobs(util.InterfacePtr(*testData.DeliveryServices[0].XMLID), nil)
	if err != nil {
		t.Fatalf("unable to get invalidation jobs: %v", err)
	}

	var realJob *tc.InvalidationJob
	for i, job := range jobs {
		d := (*newJob.DeliveryService)
		y := d.(*string)
		diff := newJob.StartTime.Time.Sub(job.StartTime.Time)
		if *job.DeliveryService == *y && *job.CreatedBy == "admin" &&
			diff < time.Second {
			realJob = &jobs[i]
			break
		}
	}

	if realJob == nil || *realJob.ID == 0 {
		t.Fatal("could not find new job")
	}

	newTime.Time = startTime.Time.Add(time.Hour * 2)
	realJob.StartTime = &newTime
	alerts, _, err = TOSession.UpdateInvalidationJob(*realJob)
	if err != nil {
		t.Fatalf("expected invalidation job update to succeed: %v", err)
	}

	if len(alerts.Alerts) != 2 {
		t.Fatalf("expected 2 alerts on update, got %v", len(alerts.Alerts))
	}

	if alerts.Alerts[0].Level != tc.WarnLevel.String() {
		t.Fatalf("expected first alert to be a warning, got %v", alerts.Alerts[0].Level)
	}

	if !strings.Contains(alerts.Alerts[0].Text, *firstJob.Regex) {
		t.Fatalf("expected first alert to be about the first job, got: %v", alerts.Alerts[0].Text)
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

func CreateTestInvalidJob(t *testing.T) {
	toDSes, _, err := TOSession.GetDeliveryServicesNullable()
	if err != nil {
		t.Fatalf("cannot GET Delivery Services: %v - %v", err, toDSes)
	}
	dsNameIDs := map[string]int64{}
	for _, ds := range toDSes {
		dsNameIDs[*ds.XMLID] = int64(*ds.ID)
	}

	job := testData.InvalidationJobs[0]
	_, ok := dsNameIDs[(*job.DeliveryService).(string)]
	if !ok {
		t.Fatalf("can't create test data job: delivery service '%v' not found in Traffic Ops", job.DeliveryService)
	}
	maxRevalDays := 0
	foundMaxRevalDays := false
	for _, p := range testData.Parameters {
		if p.Name != "maxRevalDurationDays" {
			continue
		}
		maxRevalDays, err = strconv.Atoi(p.Value)
		if err != nil {
			t.Fatalf("unable to parse maxRevalDurationDays value '%s' to int", p.Value)
		}
		foundMaxRevalDays = true
	}
	if !foundMaxRevalDays {
		t.Fatalf("expected: parameter named maxRevalDurationDays, actual: not found")
	}
	tooHigh := interface{}((maxRevalDays * 24) + 1)
	job.TTL = &tooHigh
	_, reqInf, err := TOSession.CreateInvalidationJob(job)
	if err == nil {
		t.Error("creating invalid job (TTL higher than maxRevalDurationDays) - expected: error, actual: nil error")
	}
	if reqInf.StatusCode < http.StatusBadRequest || reqInf.StatusCode >= http.StatusInternalServerError {
		t.Errorf("creating invalid job (TTL higher than maxRevalDurationDays) - expected: 400-level status code, actual: %d", reqInf.StatusCode)
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
