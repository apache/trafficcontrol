package v4

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

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
	client "github.com/apache/trafficcontrol/traffic_ops/v4-client"
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
		GetTestJobsByValidData(t)
		GetTestJobsByInvalidData(t)
		CreateTestJobsInvalidDS(t)
		CreateTestJobsAlreadyExistTTL(t)
		CreateTestJobsWithPastDate(t)
		CreateTestJobsWithFutureDate(t)
		CreateJobsMissingDate(t)
		CreateJobsMissingRegex(t)
		CreateJobsMissingTtl(t)
		UpdateTestJobsInvalidDS(t)
		DeleteTestJobs(t)
		DeleteTestJobsByInvalidId(t)
	})
}

func CreateTestJobs(t *testing.T) {
	toDSes, _, err := TOSession.GetDeliveryServices(client.RequestOptions{})
	if err != nil {
		t.Fatalf("cannot get Delivery Services: %v - alerts: %+v", err, toDSes.Alerts)
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
		resp, _, err := TOSession.CreateInvalidationJob(request, client.RequestOptions{})
		if err != nil {
			t.Errorf("could not create job: %v - alerts: %+v", err, resp.Alerts)
		}
	}
}

func JobCollisionWarningTest(t *testing.T) {
	if len(testData.DeliveryServices) < 1 {
		t.Fatal("Need at least one Delivery Service to test Invalidation Job collisions")
	}
	if testData.DeliveryServices[0].XMLID == nil {
		t.Fatal("Found a Delivery Service in the testing data with null or undefined XMLID")
	}
	xmlID := *testData.DeliveryServices[0].XMLID

	startTime := tc.Time{
		Time:  time.Now().Add(time.Hour),
		Valid: true,
	}
	firstJob := tc.InvalidationJobInput{
		DeliveryService: util.InterfacePtr(&xmlID),
		Regex:           util.StrPtr(`/\.*([A-Z]0?)`),
		TTL:             util.InterfacePtr(16),
		StartTime:       &startTime,
	}

	resp, _, err := TOSession.CreateInvalidationJob(firstJob, client.RequestOptions{})
	if err != nil {
		t.Fatalf("Unexpected error creating a content invalidation Job: %v - alerts: %+v", err, resp.Alerts)
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

	alerts, _, err := TOSession.CreateInvalidationJob(newJob, client.RequestOptions{})
	if err != nil {
		t.Fatalf("expected invalidation job create to succeed: %v - %+v", err, alerts.Alerts)
	}

	if len(alerts.Alerts) < 2 {
		t.Fatalf("expected at least 2 alerts on creation, got %v", len(alerts.Alerts))
	}

	found := false
	for _, alert := range alerts.Alerts {
		if alert.Level == tc.WarnLevel.String() && strings.Contains(alert.Text, *firstJob.Regex) {
			found = true
		}
	}
	if !found {
		t.Error("Expected a warning-level error about the regular expression, but couldn't find one")
	}

	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("deliveryService", xmlID)
	jobs, _, err := TOSession.GetInvalidationJobs(opts)
	if err != nil {
		t.Fatalf("unable to get invalidation jobs: %v - alerts: %+v", err, jobs.Alerts)
	}

	var realJob *tc.InvalidationJob
	for i, job := range jobs.Response {
		if job.StartTime == nil || job.DeliveryService == nil || job.CreatedBy == nil {
			t.Error("Traffic Ops returned a representation of a content invalidation Job that had null or undefined Start Time and/or Delivery Service and/or Created By")
			continue
		}
		diff := newJob.StartTime.Time.Sub(job.StartTime.Time)
		if *job.DeliveryService == xmlID && *job.CreatedBy == "admin" && diff < time.Second {
			realJob = &jobs.Response[i]
			break
		}
	}

	if realJob == nil || realJob.ID == nil || *realJob.ID == 0 {
		t.Fatal("could not find new job")
	}

	newTime.Time = startTime.Time.Add(time.Hour * 2)
	realJob.StartTime = &newTime
	alerts, _, err = TOSession.UpdateInvalidationJob(*realJob, client.RequestOptions{})
	if err != nil {
		t.Fatalf("expected invalidation job update to succeed: %v - alerts: %+v", err, alerts.Alerts)
	}

	if len(alerts.Alerts) < 2 {
		t.Fatalf("expected at least 2 alerts on update, got %v", len(alerts.Alerts))
	}

	found = false
	for _, alert := range alerts.Alerts {
		if alert.Level == tc.WarnLevel.String() && strings.Contains(alert.Text, *firstJob.Regex) {
			found = true
		}
	}
	if !found {
		t.Error("Expected a warning-level error about the regular expression, but couldn't find one")
	}
}

func CreateTestInvalidationJobs(t *testing.T) {
	toDSes, _, err := TOSession.GetDeliveryServices(client.RequestOptions{})
	if err != nil {
		t.Fatalf("cannot get Delivery Services: %v - alerts: %+v", err, toDSes.Alerts)
	}
	dsNameIDs := make(map[string]int64, len(toDSes.Response))
	for _, ds := range toDSes.Response {
		if ds.XMLID == nil || ds.ID == nil {
			t.Error("Traffic Ops returned a representation of a Delivery Service that had null or undefined XMLID and/or ID")
			continue
		}
		dsNameIDs[*ds.XMLID] = int64(*ds.ID)
	}

	for _, job := range testData.InvalidationJobs {
		if job.DeliveryService == nil {
			t.Error("Found a Job in the test data that has null or undefined Delivery Service")
		}
		_, ok := dsNameIDs[(*job.DeliveryService).(string)]
		if !ok {
			t.Fatalf("can't create test data job: delivery service '%v' not found in Traffic Ops", job.DeliveryService)
		}
		if alerts, _, err := TOSession.CreateInvalidationJob(job, client.RequestOptions{}); err != nil {
			t.Errorf("could not create job: %v - alerts: %+v", err, alerts)
		}
	}
}

func CreateTestInvalidJob(t *testing.T) {
	toDSes, _, err := TOSession.GetDeliveryServices(client.RequestOptions{})
	if err != nil {
		t.Fatalf("cannot get Delivery Services: %v - alerts: %+v", err, toDSes.Alerts)
	}
	dsNameIDs := make(map[string]int64, len(toDSes.Response))
	for _, ds := range toDSes.Response {
		if ds.XMLID == nil || ds.ID == nil {
			t.Error("Traffic Ops returned a representation of a Delivery Service that had null or undefined XMLID and/or ID")
			continue
		}
		dsNameIDs[*ds.XMLID] = int64(*ds.ID)
	}

	if len(testData.InvalidationJobs) < 1 {
		t.Fatal("Need at least one Invalidation Job to test creating an invalid Job")
	}
	job := testData.InvalidationJobs[0]
	if job.DeliveryService == nil {
		t.Fatal("Found a Job in the testing data that has null or undefined Delivery Service")
	}
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
	_, reqInf, err := TOSession.CreateInvalidationJob(job, client.RequestOptions{})
	if err == nil {
		t.Error("creating invalid job (TTL higher than maxRevalDurationDays) - expected: error, actual: nil error")
	}
	if reqInf.StatusCode < http.StatusBadRequest || reqInf.StatusCode >= http.StatusInternalServerError {
		t.Errorf("creating invalid job (TTL higher than maxRevalDurationDays) - expected: 400-level status code, actual: %d", reqInf.StatusCode)
	}
}

func GetTestJobsQueryParams(t *testing.T) {
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("deliveryService", "ds2")
	toJobs, _, err := TOSession.GetInvalidationJobs(opts)
	if err != nil {
		t.Fatalf("error getting jobs for Delivery Service 'ds2': %v - alerts: %+v", err, toJobs.Alerts)
	}
	foundOne := false
	for _, j := range toJobs.Response {
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
	toJobs, _, err := TOSession.GetInvalidationJobs(client.RequestOptions{})
	if err != nil {
		t.Fatalf("error getting jobs: %v - alerts: %+v", err, toJobs.Alerts)
	}

	toDSes, _, err := TOSession.GetDeliveryServices(client.RequestOptions{})
	if err != nil {
		t.Fatalf("cannot get Delivery Services: %v - alerts: %+v", err, toDSes.Alerts)
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
		for j, toJob := range toJobs.Response {
			if toJob.DeliveryService == nil {
				t.Errorf("to job (index %v) has nil delivery service", j)
				continue
			}
			if toJob.AssetURL == nil {
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
				t.Errorf("test job ds %v regex %s start time expected '%+v' actual '%+v'", *testJob.DeliveryService, *testJob.Regex, testJobTime, toJobTime)
				continue
			}
			found = true
			break
		}
		if !found {
			t.Errorf("test job ds %v regex %s expected: exists, actual: not found", *testJob.DeliveryService, *testJob.Regex)
		}
	}
}

func GetTestInvalidationJobs(t *testing.T) {
	jobs, _, err := TOSession.GetInvalidationJobs(client.RequestOptions{})
	if err != nil {
		t.Fatalf("error getting invalidation jobs: %v - alerts: %+v", err, jobs.Alerts)
	}

	toDSes, _, err := TOSession.GetDeliveryServices(client.RequestOptions{})
	if err != nil {
		t.Fatalf("cannot get Delivery Services: %v - alerts: %+v", err, toDSes.Alerts)
	}

	for _, ds := range toDSes.Response {
		if ds.ID == nil {
			t.Fatal("Erroneous Delivery Service - has invalid ID: <nil>")
		}
		if *ds.ID <= 0 {
			t.Fatalf("Erroneous Delivery Service - has invalid ID: %d", *ds.ID)
		}
	}

	for _, testJob := range testData.InvalidationJobs {
		found := false
		for _, toJob := range jobs.Response {
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

func GetTestJobsByValidData(t *testing.T) {
	toJobs, _, err := TOSession.GetInvalidationJobs(client.RequestOptions{})
	if err != nil {
		t.Fatalf("error getting jobs %v - alerts: %+v", err, toJobs.Alerts)
	}
	if len(toJobs.Response) < 1 {
		t.Fatal("Need at least one Jobs to test GET Jobs scenario")
	}
	jobs := toJobs.Response[0]

	assetUrl := jobs.AssetURL
	createdBy := jobs.CreatedBy
	id := jobs.ID
	dsName := jobs.DeliveryService
	keyword := jobs.Keyword

	//Get Jobs by Asset URL
	if len(*assetUrl) > 1 {
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("assetUrl", *assetUrl)
		toJobs, _, _ = TOSession.GetInvalidationJobs(opts)
		if len(toJobs.Response) < 1 {
			t.Errorf("Expected atleast one Jobs response for GET Jobs by Asset URL, but found %d ", len(toJobs.Response))
		}
	} else {
		t.Errorf("Asset URL Field is Empty, so can't test get jobs")
	}

	//Get Jobs by CreatedBy
	if len(*createdBy) > 1 {
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("createdBy", *createdBy)
		toJobs, _, _ = TOSession.GetInvalidationJobs(opts)
		if len(toJobs.Response) < 1 {
			t.Errorf("Expected atleast one Jobs response for GET Jobs by CreatedBy, but found %d ", len(toJobs.Response))
		}
	} else {
		t.Errorf("CreatedBy Field is empty, so can't test get jobs")
	}

	//Get Jobs by ID
	if *id > 1 {
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("id", strconv.FormatUint(uint64(*id), 10))
		toJobs, _, _ = TOSession.GetInvalidationJobs(opts)
		if len(toJobs.Response) != 1 {
			t.Errorf("Expected only one Jobs response for GET Jobs by ID, but found %d ", len(toJobs.Response))
		}
	} else {
		t.Errorf("ID Field is empty, so can't test get jobs %d", *id)
	}

	//Get Jobs by Keyword
	if len(*keyword) > 1 {
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("Keyword", *keyword)
		toJobs, _, _ = TOSession.GetInvalidationJobs(opts)
		if len(toJobs.Response) < 1 {
			t.Errorf("Expected atleast one Jobs response for GET Jobs by keyword, but found %d ", len(toJobs.Response))
		}
	} else {
		t.Errorf("Keyword field is empty, so can't test get jobs")
	}

	//Get Delivery Service ID by Name
	if len(*dsName) > 0 {
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("xmlId", *dsName)
		toDSes, _, _ := TOSession.GetDeliveryServices(opts)
		if len(toDSes.Response) > 0 {
			dsId := toDSes.Response[0].ID
			if *dsId > 0 {
				//Get Jobs by DSID
				opts := client.NewRequestOptions()
				opts.QueryParameters.Set("dsId", strconv.Itoa(*dsId))
				toJobs, _, _ = TOSession.GetInvalidationJobs(opts)
				if len(toJobs.Response) < 1 {
					t.Errorf("Expected atleast one Jobs response for GET Jobs by delivery service, but found %d ", len(toJobs.Response))
				}
			} else {
				t.Error("Delivery service id is empty")
			}
		} else {
			t.Error("No responses for get delivery service by name")
		}
	} else {
		t.Error("Delivery Service Name field is empty, so can't retrive ID from name")
	}

	//Get UserID ID by Username
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("username", "admin")
	userResp, _, _ := TOSession.GetUsers(opts)
	if len(userResp.Response) > 0 {
		userId := userResp.Response[0].ID
		if *userId > 0 {
			//Get Jobs by userID
			opts := client.NewRequestOptions()
			opts.QueryParameters.Set("userId", strconv.Itoa(*userId))
			toJobs, _, _ = TOSession.GetInvalidationJobs(opts)
			if len(toJobs.Response) < 1 {
				t.Errorf("Expected atleast one Jobs response for GET Jobs by users, but found %d ", len(toJobs.Response))
			}
		} else {
			t.Error("User id is empty")
		}
	} else {
		t.Error("No user response available for get user by name")
	}

	// get maxRevalDurationDays parameter value
	opts = client.NewRequestOptions()
	opts.QueryParameters.Set("name", "maxRevalDurationDays")
	params, _, err := TOSession.GetParameters(opts)
	if err != nil {
		t.Fatalf("unexpected error getting maxRevalDurationDays parameter")
	}
	if len(params.Response) != 1 {
		t.Fatalf("expected 1 maxRevalDurationDays parameter, got %d", len(params.Response))
	}
	maxRevalDurationDays, err := strconv.Atoi(params.Response[0].Value)
	if err != nil {
		t.Fatalf("unexpected error converting maxRevalDurationDays value to int: %v", err)
	}

	// get jobs created within maxRevalDurationDays window
	opts = client.NewRequestOptions()
	opts.QueryParameters.Set("maxRevalDurationDays", "")
	maxRevalJobs, _, err := TOSession.GetInvalidationJobs(opts)
	if err != nil {
		t.Errorf("unexpected error getting jobs by maxRevalDurationDays: %v", err)
	} else if len(maxRevalJobs.Response) < 1 {
		t.Errorf("GET /jobs by maxRevalDurationDays - expected at least 1 job")
	}
	for _, j := range maxRevalJobs.Response {
		if time.Since((*j.StartTime).Time) > time.Duration(maxRevalDurationDays)*24*time.Hour {
			t.Errorf("GET /jobs by maxRevalDurationDays returned job that is older than %d days: {%s, %s, %v}", maxRevalDurationDays, *j.DeliveryService, *j.AssetURL, *j.StartTime)
		}
	}

	// create DS xml_id -> cdn_id lookup map
	dses, _, err := TOSession.GetDeliveryServices(client.NewRequestOptions())
	if err != nil {
		t.Fatalf("unexpectd error getting delivery services: %v", err)
	}
	dsToCDN := make(map[string]string, len(dses.Response))
	for _, ds := range dses.Response {
		dsToCDN[*ds.XMLID] = *ds.CDNName
	}

	cdn := "cdn2"
	// get jobs by CDN ID
	opts = client.NewRequestOptions()
	opts.QueryParameters.Set("cdn", cdn)
	cdnJobs, _, err := TOSession.GetInvalidationJobs(opts)
	if err != nil {
		t.Errorf("unexpected error getting jobs by cdn: %v", err)
	} else if len(cdnJobs.Response) < 1 {
		t.Errorf("GET /jobs by cdn - expected at least 1 job")
	}
	for _, j := range cdnJobs.Response {
		if dsToCDN[*j.DeliveryService] != cdn {
			t.Errorf("GET /jobs by cdn returned job that does not belong to CDN %s: {%s, %s, %v}", cdn, *j.DeliveryService, *j.AssetURL, *j.StartTime)
		}
	}
}

func GetTestJobsByInvalidData(t *testing.T) {

	//Get Jobs by Invalid Asset URL
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("assetUrl", "abcd")
	toJobs, _, _ := TOSession.GetInvalidationJobs(opts)
	if len(toJobs.Response) != 0 {
		t.Errorf("Expected no response from Get Jobs by Invalid Asset URL, but found %d ", len(toJobs.Response))
	}

	//Get Jobs by Invalid CreatedBy
	opts = client.NewRequestOptions()
	opts.QueryParameters.Set("createdBy", "abcd")
	toJobs, _, _ = TOSession.GetInvalidationJobs(opts)
	if len(toJobs.Response) != 0 {
		t.Errorf("Expected no response from Get Jobs by Invalid CreatedBy, but found %d ", len(toJobs.Response))
	}

	//Get Jobs by Invalid ID
	opts = client.NewRequestOptions()
	opts.QueryParameters.Set("id", "11111")
	toJobs, _, _ = TOSession.GetInvalidationJobs(opts)
	if len(toJobs.Response) != 0 {
		t.Errorf("Expected no response from Get Jobs by Invalid ID, but found %d ", len(toJobs.Response))
	}

	//Get Jobs by Invalid Keyword
	opts = client.NewRequestOptions()
	opts.QueryParameters.Set("keyword", "invalid")
	toJobs, _, _ = TOSession.GetInvalidationJobs(opts)
	if len(toJobs.Response) != 0 {
		t.Errorf("Expected no response from Get Jobs by Invalid Keyword, but found %d ", len(toJobs.Response))
	}

	//Get Jobs by Invalid DSID
	opts = client.NewRequestOptions()
	opts.QueryParameters.Set("dsId", "11111")
	toJobs, _, _ = TOSession.GetInvalidationJobs(opts)
	if len(toJobs.Response) != 0 {
		t.Errorf("Expected no response from Get Jobs by Invalid DSID, but found %d ", len(toJobs.Response))
	}

	//Get Jobs by Invalid DSName
	opts = client.NewRequestOptions()
	opts.QueryParameters.Set("deliveryService", "abcd")
	toJobs, _, _ = TOSession.GetInvalidationJobs(opts)
	if len(toJobs.Response) != 0 {
		t.Errorf("Expected no response from Get Jobs by Invalid DSName, but found %d ", len(toJobs.Response))
	}

	//Get Jobs by Invalid userID
	opts = client.NewRequestOptions()
	opts.QueryParameters.Set("userId", "11111")
	toJobs, _, _ = TOSession.GetInvalidationJobs(opts)
	if len(toJobs.Response) != 0 {
		t.Errorf("Expected no response from Get Jobs by Invalid userID, but found %d ", len(toJobs.Response))
	}
}

func CreateTestJobsInvalidDS(t *testing.T) {
	if len(testData.InvalidationJobs) < 1 {
		t.Error("Need at least one Invalidation Jobs to test invalid ds")
	}
	job := testData.InvalidationJobs[0]
	job.StartTime = &tc.Time{
		Time:  time.Now().Add(time.Minute).UTC(),
		Valid: true,
	}
	testData.InvalidationJobs[0] = job

	//Invalid DS
	request := tc.InvalidationJobInput{
		DeliveryService: util.InterfacePtr("invalid"),
		Regex:           job.Regex,
		StartTime:       job.StartTime,
		TTL:             job.TTL,
	}
	resp, reqInf, err := TOSession.CreateInvalidationJob(request, client.RequestOptions{})
	if err == nil {
		t.Errorf("Expected No DeliveryService exists matching identifier: %v - alerts: %v", request.DeliveryService, resp.Alerts)
	}
	if reqInf.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code 400, got %v", reqInf.StatusCode)
	}

	//Missing DS
	request = tc.InvalidationJobInput{
		Regex:     job.Regex,
		StartTime: job.StartTime,
		TTL:       job.TTL,
	}
	resp, reqInf, err = TOSession.CreateInvalidationJob(request, client.RequestOptions{})
	if err == nil {
		t.Errorf("Expected deliveryService: cannot be blank - alerts: %v", resp.Alerts)
	}
	if reqInf.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code 400, got %v", reqInf.StatusCode)
	}

	//Empty DS
	request = tc.InvalidationJobInput{
		DeliveryService: util.InterfacePtr(""),
		Regex:           job.Regex,
		StartTime:       job.StartTime,
		TTL:             job.TTL,
	}
	resp, reqInf, err = TOSession.CreateInvalidationJob(request, client.RequestOptions{})
	if err == nil {
		t.Errorf("Expected deliveryService: cannot be blank., No DeliveryService exists matching identifier: - alerts: %v", resp.Alerts)
	}
	if reqInf.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code 400, got %v", reqInf.StatusCode)
	}
}

func CreateTestJobsAlreadyExistTTL(t *testing.T) {
	if len(testData.InvalidationJobs) < 1 {
		t.Error("Need at least one Invalidation Jobs to create duplicate data")
	}
	job := testData.InvalidationJobs[0]
	job.StartTime = &tc.Time{
		Time:  time.Now().Add(time.Minute).UTC(),
		Valid: true,
	}
	testData.InvalidationJobs[0] = job

	request := tc.InvalidationJobInput{
		DeliveryService: job.DeliveryService,
		Regex:           job.Regex,
		StartTime:       job.StartTime,
		TTL:             job.TTL,
	}
	resp, _, err := TOSession.CreateInvalidationJob(request, client.RequestOptions{})
	if err != nil {
		t.Errorf("Expected Invalidation request created, but found error %v - Alert %v", err, resp.Alerts)
	}
}

func CreateTestJobsWithPastDate(t *testing.T) {
	if len(testData.InvalidationJobs) < 1 {
		t.Fatal("Need at least one Invalidation Job to test creating an invalid Job")
	}
	//past start date
	dt := time.Now()
	dt.Format("2019-06-18 21:28:31")
	job := testData.InvalidationJobs[0]
	job.StartTime = &tc.Time{
		Time:  dt.AddDate(0, 0, -1),
		Valid: true,
	}
	testData.InvalidationJobs[0] = job
	request := tc.InvalidationJobInput{
		DeliveryService: job.DeliveryService,
		Regex:           job.Regex,
		StartTime:       job.StartTime,
		TTL:             job.TTL,
	}
	resp, reqInf, err := TOSession.CreateInvalidationJob(request, client.RequestOptions{})
	if err == nil {
		t.Errorf("Expected startTime: must be in the future - Alert %v", resp.Alerts)
	}
	if reqInf.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code 400, got %v", reqInf.StatusCode)
	}

	//RFC Format past start date
	dt = time.Now()
	dt.Format("2019-10-12T07:20:50.52Z")
	job = testData.InvalidationJobs[0]
	job.StartTime = &tc.Time{
		Time:  dt.AddDate(0, 0, -1),
		Valid: true,
	}
	testData.InvalidationJobs[0] = job
	request = tc.InvalidationJobInput{
		DeliveryService: job.DeliveryService,
		Regex:           job.Regex,
		StartTime:       job.StartTime,
		TTL:             job.TTL,
	}
	resp, reqInf, err = TOSession.CreateInvalidationJob(request, client.RequestOptions{})
	if err == nil {
		t.Errorf("Expected startTime: must be in the future - Alert %v", resp.Alerts)
	}
	if reqInf.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code 400, got %v", reqInf.StatusCode)
	}

	//Non standard Format past start date
	dt = time.Now()
	dt.Format("2020-03-11 14:12:20-06")
	job = testData.InvalidationJobs[0]
	job.StartTime = &tc.Time{
		Time:  dt.AddDate(0, 0, -5),
		Valid: true,
	}
	testData.InvalidationJobs[0] = job
	request = tc.InvalidationJobInput{
		DeliveryService: job.DeliveryService,
		Regex:           job.Regex,
		StartTime:       job.StartTime,
		TTL:             job.TTL,
	}
	resp, reqInf, err = TOSession.CreateInvalidationJob(request, client.RequestOptions{})
	if err == nil {
		t.Errorf("Expected startTime: must be in the future - Alert %v", resp.Alerts)
	}
	if reqInf.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code 400, got %v", reqInf.StatusCode)
	}

	//unix standard format past start date
	job = testData.InvalidationJobs[0]
	job.StartTime = &tc.Time{
		Time:  time.Unix(1, 0),
		Valid: true,
	}
	testData.InvalidationJobs[0] = job
	request = tc.InvalidationJobInput{
		DeliveryService: job.DeliveryService,
		Regex:           job.Regex,
		StartTime:       job.StartTime,
		TTL:             job.TTL,
	}
	resp, reqInf, err = TOSession.CreateInvalidationJob(request, client.RequestOptions{})
	if err == nil {
		t.Errorf("Expected startTime: must be in the future - Alert %v", resp.Alerts)
	}
	if reqInf.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code 400, got %v", reqInf.StatusCode)
	}
}

func CreateTestJobsWithFutureDate(t *testing.T) {
	if len(testData.InvalidationJobs) < 1 {
		t.Fatal("Need at least one Invalidation Job to test creating an invalid Job")
	}
	//RFC Future start date
	dt := time.Now()
	dt.Format("2019-10-12T07:20:50.52Z")
	job := testData.InvalidationJobs[0]
	job.StartTime = &tc.Time{
		Time:  dt.AddDate(0, 0, 1),
		Valid: true,
	}
	testData.InvalidationJobs[0] = job
	request := tc.InvalidationJobInput{
		DeliveryService: job.DeliveryService,
		Regex:           job.Regex,
		StartTime:       job.StartTime,
		TTL:             job.TTL,
	}
	resp, reqInf, err := TOSession.CreateInvalidationJob(request, client.RequestOptions{})
	if err != nil {
		t.Errorf("Expected Invalidation request created, but found error %v - Alert %v", err, resp.Alerts)
	}
	if reqInf.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %v", reqInf.StatusCode)
	}

	//Non standard format Future start date
	dt = time.Now()
	dt.Format("2020-03-11 14:12:20-06")
	job = testData.InvalidationJobs[0]
	job.StartTime = &tc.Time{
		Time:  dt.AddDate(0, 0, 1),
		Valid: true,
	}
	testData.InvalidationJobs[0] = job
	request = tc.InvalidationJobInput{
		DeliveryService: job.DeliveryService,
		Regex:           job.Regex,
		StartTime:       job.StartTime,
		TTL:             job.TTL,
	}
	resp, reqInf, err = TOSession.CreateInvalidationJob(request, client.RequestOptions{})
	if err != nil {
		t.Errorf("Expected Invalidation request created, but found error %v - Alert %v", err, resp.Alerts)
	}
	if reqInf.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %v", reqInf.StatusCode)
	}

	//UNIX format Future start date
	dt = time.Now()
	dt.Format(".000")
	job = testData.InvalidationJobs[0]
	job.StartTime = &tc.Time{
		Time:  dt.AddDate(0, 0, 1),
		Valid: true,
	}
	testData.InvalidationJobs[0] = job
	request = tc.InvalidationJobInput{
		DeliveryService: job.DeliveryService,
		Regex:           job.Regex,
		StartTime:       job.StartTime,
		TTL:             job.TTL,
	}
	resp, reqInf, err = TOSession.CreateInvalidationJob(request, client.RequestOptions{})
	if err != nil {
		t.Errorf("Expected Invalidation request created, but found error %v - Alert %v", err, resp.Alerts)
	}
	if reqInf.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %v", reqInf.StatusCode)
	}
}

func CreateJobsMissingDate(t *testing.T) {
	if len(testData.InvalidationJobs) < 1 {
		t.Fatal("Need at least one Invalidation Job to test creating an invalid Job")
	}
	//Missing date
	job := testData.InvalidationJobs[0]
	request := tc.InvalidationJobInput{
		DeliveryService: job.DeliveryService,
		Regex:           job.Regex,
		TTL:             job.TTL,
	}
	resp, reqInf, err := TOSession.CreateInvalidationJob(request, client.RequestOptions{})
	if err == nil {
		t.Errorf("Expected startTime: cannot be blank, but no error found %v - Alert %v", resp, resp.Alerts)
	}
	if reqInf.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code 400, got %v", reqInf.StatusCode)
	}
}

func CreateJobsMissingRegex(t *testing.T) {
	if len(testData.InvalidationJobs) < 1 {
		t.Fatal("Need at least one Invalidation Job to test creating an invalid Job")
	}
	//Missing Regex
	//Future start date
	dt := time.Now()
	dt.Format("2019-10-12T07:20:50.52Z")
	job := testData.InvalidationJobs[0]
	job.StartTime = &tc.Time{
		Time:  dt.AddDate(0, 0, 1),
		Valid: true,
	}
	testData.InvalidationJobs[0] = job
	request := tc.InvalidationJobInput{
		DeliveryService: job.DeliveryService,
		TTL:             job.TTL,
		StartTime:       job.StartTime,
	}
	resp, reqInf, err := TOSession.CreateInvalidationJob(request, client.RequestOptions{})
	if err == nil {
		t.Errorf("Expected regex: cannot be blank, but no error found %v - Alert %v", resp, resp.Alerts)
	}
	if reqInf.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code 400, got %v", reqInf.StatusCode)
	}

	//Empty Regex
	job.Regex = nil
	request = tc.InvalidationJobInput{
		DeliveryService: job.DeliveryService,
		Regex:           job.Regex,
		TTL:             job.TTL,
		StartTime:       job.StartTime,
	}
	resp, _, err = TOSession.CreateInvalidationJob(request, client.RequestOptions{})
	if err == nil {
		t.Errorf("Expected regex: cannot be blank, but no error found %v - Alert %v", resp, resp.Alerts)
	}
	if reqInf.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code 400, got %v", reqInf.StatusCode)
	}
}

func CreateJobsMissingTtl(t *testing.T) {
	if len(testData.InvalidationJobs) < 1 {
		t.Fatal("Need at least one Invalidation Job to test creating an invalid Job")
	}
	//Missing TTL
	//Future start date
	dt := time.Now()
	dt.Format("2019-10-12T07:20:50.52Z")
	job := testData.InvalidationJobs[0]
	job.StartTime = &tc.Time{
		Time:  dt.AddDate(0, 0, 1),
		Valid: true,
	}
	testData.InvalidationJobs[0] = job
	request := tc.InvalidationJobInput{
		DeliveryService: job.DeliveryService,
		Regex:           job.Regex,
		StartTime:       job.StartTime,
	}
	resp, reqInf, err := TOSession.CreateInvalidationJob(request, client.RequestOptions{})
	if err == nil {
		t.Errorf("Expected ttl: cannot be blank., but no error found %v - Alert %v", resp, resp.Alerts)
	}
	if reqInf.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code 400, got %v", reqInf.StatusCode)
	}

	//Empty TTL
	job.TTL = nil
	request = tc.InvalidationJobInput{
		DeliveryService: job.DeliveryService,
		Regex:           job.Regex,
		TTL:             job.TTL,
		StartTime:       job.StartTime,
	}
	resp, _, err = TOSession.CreateInvalidationJob(request, client.RequestOptions{})
	if err == nil {
		t.Errorf("Expected ttl: cannot be blank., ttl: must be a number of hours, or a duration string e.g. '48h', but no error found %v - Alert %v", resp, resp.Alerts)
	}
	if reqInf.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code 400, got %v", reqInf.StatusCode)
	}
}

func UpdateTestJobsInvalidDS(t *testing.T) {
	if len(testData.DeliveryServices) < 2 {
		t.Fatal("Need at least two Delivery Service to update Invalidation Job")
	}
	if testData.DeliveryServices[0].XMLID == nil || testData.DeliveryServices[1].XMLID == nil {
		t.Fatal("Found a Delivery Service in the testing data with null or undefined XMLID")
	}
	xmlID := *testData.DeliveryServices[0].XMLID
	startTime := tc.Time{
		Time:  time.Now().Add(time.Hour),
		Valid: true,
	}
	firstJob := tc.InvalidationJobInput{
		DeliveryService: util.InterfacePtr(&xmlID),
		Regex:           util.StrPtr(`/\.*([A-Z]0?)`),
		TTL:             util.InterfacePtr(16),
		StartTime:       &startTime,
	}

	resp, _, err := TOSession.CreateInvalidationJob(firstJob, client.RequestOptions{})
	if err != nil {
		t.Fatalf("Unexpected error creating a content invalidation Job: %v - alerts: %+v", err, resp.Alerts)
	}
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("deliveryService", xmlID)
	jobs, _, err := TOSession.GetInvalidationJobs(opts)
	if err != nil {
		t.Fatalf("unable to get invalidation jobs: %v - alerts: %+v", err, jobs.Alerts)
	}

	var realJob tc.InvalidationJob
	for i, job := range jobs.Response {
		if job.StartTime == nil || job.DeliveryService == nil || job.CreatedBy == nil {
			t.Error("Traffic Ops returned a representation of a content invalidation Job that had null or undefined Start Time and/or Delivery Service and/or Created By")
			continue
		}
		diff := firstJob.StartTime.Time.Sub(job.StartTime.Time)
		if *job.DeliveryService == xmlID && *job.CreatedBy == "admin" && diff < time.Second {
			realJob = jobs.Response[i]
			break
		}
	}
	if realJob.ID == nil || *realJob.ID == 0 {
		t.Fatal("could not find new job")
	}

	//update existing jobs with new ds id
	originalJob := realJob
	newTime := tc.Time{
		Time:  startTime.Time.Add(time.Hour * 2),
		Valid: true,
	}
	originalJob.StartTime = &newTime
	originalJob.DeliveryService = testData.DeliveryServices[1].XMLID
	alerts, reqInf, err := TOSession.UpdateInvalidationJob(originalJob, client.RequestOptions{})
	if err == nil {
		t.Fatalf("Expected Cannot change 'deliveryService' of existing invalidation job! - alerts: %+v", alerts.Alerts)
	}
	if reqInf.StatusCode != http.StatusConflict {
		t.Errorf("Expected status code 409, got %v", reqInf.StatusCode)
	}

	//update existing jobs with invalid ds id
	invalidDsIdJob := realJob
	invalidDsId := "abcd"
	invalidDsIdJob.DeliveryService = &invalidDsId
	alerts, reqInf, err = TOSession.UpdateInvalidationJob(invalidDsIdJob, client.RequestOptions{})
	if err == nil {
		t.Fatalf("Expected Cannot change 'deliveryService' of existing invalidation job! - alerts: %+v", alerts.Alerts)
	}
	if reqInf.StatusCode != http.StatusConflict {
		t.Errorf("Expected status code 409, got %v", reqInf.StatusCode)
	}

	//update existing jobs with blank ds id
	blankDsIdJob := realJob
	blankDsId := ""
	blankDsIdJob.DeliveryService = &blankDsId
	alerts, reqInf, err = TOSession.UpdateInvalidationJob(blankDsIdJob, client.RequestOptions{})
	if err == nil {
		t.Fatalf("Expected deliveryService: cannot be blank. - alerts: %+v", alerts.Alerts)
	}
	if reqInf.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code 400, got %v", reqInf.StatusCode)
	}

	//update existing jobs with asset url not starts with origin.infra
	invalidAssetURLJob := realJob
	assetURL := "http://google.com"
	invalidAssetURLJob.AssetURL = &assetURL
	alerts, reqInf, err = TOSession.UpdateInvalidationJob(invalidAssetURLJob, client.RequestOptions{})
	if err == nil {
		t.Fatalf("Expected Cannot set asset URL that does not start with Delivery Service origin URL: http://origin.infra.ciab.test. - alerts: %+v", alerts.Alerts)
	}
	if reqInf.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code 400, got %v", reqInf.StatusCode)
	}

	//update existing jobs with blank asset url
	blankAssetURLJob := realJob
	assetURL = ""
	blankAssetURLJob.AssetURL = &assetURL
	alerts, reqInf, err = TOSession.UpdateInvalidationJob(blankAssetURLJob, client.RequestOptions{})
	if err == nil {
		t.Fatalf("Expected assetUrl: cannot be blank. alerts: %+v", alerts.Alerts)
	}
	if reqInf.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code 400, got %v", reqInf.StatusCode)
	}

	//update existing jobs with blank created by
	blankCreatedByJob := realJob
	createdBy := ""
	blankCreatedByJob.CreatedBy = &createdBy
	alerts, reqInf, err = TOSession.UpdateInvalidationJob(blankCreatedByJob, client.RequestOptions{})
	if err == nil {
		t.Fatalf("Expected createdBy: cannot be blank. alerts: %+v", alerts.Alerts)
	}
	if reqInf.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code 400, got %v", reqInf.StatusCode)
	}

	//update existing jobs created by
	createdByJob := realJob
	createdBy = "operator"
	createdByJob.CreatedBy = &createdBy
	alerts, reqInf, err = TOSession.UpdateInvalidationJob(createdByJob, client.RequestOptions{})
	if err == nil {
		t.Fatalf("Expected Cannot change 'createdBy' of existing invalidation jobs!. alerts: %+v", alerts.Alerts)
	}
	if reqInf.StatusCode != http.StatusConflict {
		t.Errorf("Expected status code 409, got %v", reqInf.StatusCode)
	}

	//update existing jobs with blank parameters
	blankParametersJob := realJob
	parameters := ""
	blankParametersJob.Parameters = &parameters
	alerts, reqInf, err = TOSession.UpdateInvalidationJob(blankParametersJob, client.RequestOptions{})
	if err == nil {
		t.Fatalf("Expected parameters: cannot be blank. alerts: %+v", alerts.Alerts)
	}
	if reqInf.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code 400, got %v", reqInf.StatusCode)
	}

	//update existing jobs start date after 2 days
	startDateFutureJob := realJob
	dt := time.Now()
	dt.Format("2019-10-12T07:20:50.52Z")
	startDateFutureJob.StartTime = &tc.Time{
		Time:  dt.AddDate(0, 0, 3),
		Valid: true,
	}
	alerts, reqInf, err = TOSession.UpdateInvalidationJob(startDateFutureJob, client.RequestOptions{})
	if err == nil {
		t.Fatalf("Expected startTime: must be within two days from now. alerts: %+v", alerts.Alerts)
	}
	if reqInf.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code 400, got %v", reqInf.StatusCode)
	}

	//update jobs with past start date
	pastStartDateJob := realJob
	dt = time.Now()
	dt.Format("2019-06-18 21:28:31")
	pastStartDateJob.StartTime = &tc.Time{
		Time:  dt.AddDate(0, 0, -3),
		Valid: true,
	}
	alerts, reqInf, err = TOSession.UpdateInvalidationJob(pastStartDateJob, client.RequestOptions{})
	if err == nil {
		t.Fatalf("Expected startTime: cannot be in the past. alerts: %+v", alerts.Alerts)
	}
	if reqInf.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code 400, got %v", reqInf.StatusCode)
	}

	//update jobs with RFC Format past start date
	pastStartDateJob = realJob
	dt = time.Now()
	dt.Format("2019-10-12T07:20:50.52Z")
	pastStartDateJob.StartTime = &tc.Time{
		Time:  dt.AddDate(0, 0, -1),
		Valid: true,
	}
	alerts, reqInf, err = TOSession.UpdateInvalidationJob(pastStartDateJob, client.RequestOptions{})
	if err == nil {
		t.Fatalf("Expected startTime: cannot be in the past. alerts: %+v", alerts.Alerts)
	}
	if reqInf.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code 400, got %v", reqInf.StatusCode)
	}

	//update jobs with UNIX Format past start date
	pastStartDateJob = realJob
	pastStartDateJob.StartTime = &tc.Time{
		Time:  time.Unix(1, 0),
		Valid: true,
	}
	alerts, reqInf, err = TOSession.UpdateInvalidationJob(pastStartDateJob, client.RequestOptions{})
	if err == nil {
		t.Fatalf("Expected startTime: cannot be in the past. alerts: %+v", alerts.Alerts)
	}
	if reqInf.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code 400, got %v", reqInf.StatusCode)
	}

	//update jobs with non standard Format past start date
	pastStartDateJob = realJob
	dt = time.Now()
	dt.Format("2020-03-11 14:12:20-06")
	pastStartDateJob.StartTime = &tc.Time{
		Time:  dt.AddDate(0, 0, -1),
		Valid: true,
	}
	alerts, reqInf, err = TOSession.UpdateInvalidationJob(pastStartDateJob, client.RequestOptions{})
	if err == nil {
		t.Fatalf("Expected startTime: cannot be in the past. alerts: %+v", alerts.Alerts)
	}
	if reqInf.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code 400, got %v", reqInf.StatusCode)
	}

	//update jobs with RFC Format Future start date
	startDateFutureJob = realJob
	dt = time.Now()
	dt.Format("2019-10-12T07:20:50.52Z")
	startDateFutureJob.StartTime = &tc.Time{
		Time:  dt.AddDate(0, 0, 1),
		Valid: true,
	}
	alerts, reqInf, err = TOSession.UpdateInvalidationJob(startDateFutureJob, client.RequestOptions{})
	if err != nil {
		t.Fatalf("Expected Content invalidation job updated. alerts: %+v", alerts.Alerts)
	}
	if reqInf.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %v", reqInf.StatusCode)
	}

	//update jobs with UNIX Format Future start date
	startDateFutureJob = realJob
	dt = time.Now()
	dt.Format(".000")
	startDateFutureJob.StartTime = &tc.Time{
		Time:  dt.AddDate(0, 0, 1),
		Valid: true,
	}
	alerts, reqInf, err = TOSession.UpdateInvalidationJob(startDateFutureJob, client.RequestOptions{})
	if err != nil {
		t.Fatalf("Expected Content invalidation job updated. alerts: %+v", alerts.Alerts)
	}
	if reqInf.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %v", reqInf.StatusCode)
	}

	//update jobs with non standard Format Future start date
	startDateFutureJob = realJob
	dt = time.Now()
	dt.Format("2020-03-11 14:12:20-06")
	startDateFutureJob.StartTime = &tc.Time{
		Time:  dt.AddDate(0, 0, 1),
		Valid: true,
	}
	alerts, reqInf, err = TOSession.UpdateInvalidationJob(startDateFutureJob, client.RequestOptions{})
	if err != nil {
		t.Fatalf("Expected Content invalidation job updated. alerts: %+v", alerts.Alerts)
	}
	if reqInf.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %v", reqInf.StatusCode)
	}

	//update existing jobs with new id
	newIdJob := realJob
	var b uint64 = 1111
	var a *uint64 = &b
	newIdJob.ID = a
	alerts, reqInf, err = TOSession.UpdateInvalidationJob(newIdJob, client.RequestOptions{})
	if err == nil {
		t.Fatalf("Expected Cannot change an invalidation job 'id'! - alerts: %+v", alerts.Alerts)
	}
	if reqInf.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status code 404, got %v", reqInf.StatusCode)
	}
}

func DeleteTestJobs(t *testing.T) {

	//Get all Jobs
	toJobs, _, err := TOSession.GetInvalidationJobs(client.RequestOptions{})
	if err != nil {
		t.Fatalf("error getting jobs %v - alerts: %+v", err, toJobs.Alerts)
	}
	if len(toJobs.Response) < 1 {
		t.Fatal("Need at least one Jobs to test GET Jobs scenario")
	}
	jobs := toJobs.Response[0]
	id := jobs.ID

	//Delete Jobs by valid id
	alerts, reqInf, err := TOSession.DeleteInvalidationJob(uint64(*id), client.RequestOptions{})
	if err != nil {
		t.Errorf("Expected Content invalidation job was deleted Error - %v, Alerts %v", err, alerts.Alerts)
	}
	if reqInf.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %v", reqInf.StatusCode)
	}
}

func DeleteTestJobsByInvalidId(t *testing.T) {

	//Delete Jobs by invalid id
	var b uint64 = 1111
	var a *uint64 = &b
	alerts, reqInf, err := TOSession.DeleteInvalidationJob(uint64(*a), client.RequestOptions{})
	if err == nil {
		t.Errorf("Expected No job by id. Error - %v, Alerts %v", err, alerts.Alerts)
	}
	if reqInf.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status code 404, got %v", reqInf.StatusCode)
	}
}
