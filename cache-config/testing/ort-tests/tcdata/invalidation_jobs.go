package tcdata

import (
	"net/http"
	"testing"
)

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

func (r *TCData) CreateTestInvalidationJobs(t *testing.T) {
	// loop through invalidation jobs and create
	for _, job := range r.TestData.InvalidationJobs {
		resp, _, err := TOSession.CreateInvalidationJob(job)
		t.Log("Response: ", job.Regex, " ", resp)
		if err != nil {
			t.Errorf("could not CREATE jobs: %v", err)
		}
	}
}

func (r *TCData) DeleteTestInvalidationJobs(t *testing.T) {
	for _, job := range r.TestData.InvalidationJobs {

		jobs, _, err := TOSession.GetInvalidationJobsWithHdr(job.DeliveryService, nil, http.Header{})
		if err != nil {
			t.Errorf("could not request jobs for DeliveryService: %v err: %v", job.DeliveryService, err)
		}
		for _, j := range jobs {
			if _, _, err = TOSession.DeleteInvalidationJob(*j.ID); err != nil {
				t.Errorf("failed to delete invalidation job %d", *j.ID)
			}
		}
	}
}
