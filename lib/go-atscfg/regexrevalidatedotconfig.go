package atscfg

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
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
)

const RegexRemapPrefix = "regex_remap_"
const CacheUrlPrefix = "cacheurl_"

const RemapFile = "remap.config"

const RegexRevalidateFileName = "regex_revalidate.config"
const RegexRevalidateMaxRevalDurationDaysParamName = "maxRevalDurationDays"
const DefaultMaxRevalDurationDays = 90
const JobKeywordPurge = "PURGE"
const RegexRevalidateMinTTL = time.Hour

type Job struct {
	AssetURL string
	PurgeEnd time.Time
}

type Jobs []Job

func (jb Jobs) Len() int      { return len(jb) }
func (jb Jobs) Swap(i, j int) { jb[i], jb[j] = jb[j], jb[i] }
func (jb Jobs) Less(i, j int) bool {
	if jb[i].AssetURL == jb[j].AssetURL {
		return jb[i].PurgeEnd.Before(jb[j].PurgeEnd)
	}
	return strings.Compare(jb[i].AssetURL, jb[j].AssetURL) < 0
}

func MakeRegexRevalidateDotConfig(
	cdnName tc.CDNName,
	params map[string][]string, // params on profile GLOBAL fileName RegexRevalidateFileName
	toToolName string, // tm.toolname global parameter (TODO: cache itself?)
	toURL string, // tm.url global parameter (TODO: cache itself?)
	jobs []tc.Job, // jobs should be jobs on DSes on this cdn
) string {

	// TODO: add cdn, startTime query params to /jobs endpoint
	err := error(nil)

	maxDays := DefaultMaxRevalDurationDays
	if maxDaysStrs := params[RegexRevalidateMaxRevalDurationDaysParamName]; len(maxDaysStrs) > 0 {
		if maxDays, err = strconv.Atoi(maxDaysStrs[0]); err != nil { // just use the first, if there were multiple params
			log.Warnln("making regex revalidate config: max days param '" + maxDaysStrs[0] + "' is not an integer, using default value!")
			maxDays = DefaultMaxRevalDurationDays
		}
	}

	maxReval := time.Duration(maxDays) * time.Hour * 24

	cfgJobs := filterJobs(jobs, maxReval, RegexRevalidateMinTTL)

	txt := GenericHeaderComment(string(cdnName), toToolName, toURL)
	for _, job := range cfgJobs {
		txt += job.AssetURL + " " + strconv.FormatInt(job.PurgeEnd.Unix(), 10) + "\n"
	}

	return txt
}

// filterJobs returns only jobs which:
//   - have a non-null deliveryservice
//   - have parameters of the form TTL:%dh
//   - have a start time later than (now + maxReval days). That is, we don't query jobs older than maxReval in the past.
//   - are "purge" jobs
//   - have a start_time+ttl > now. That is, jobs that haven't expired yet.
func filterJobs(jobs []tc.Job, maxReval time.Duration, minTTL time.Duration) []Job {
	jobMap := map[string]time.Time{}
	for _, job := range jobs {
		if job.DeliveryService == "" {
			continue
		}
		if !strings.HasPrefix(job.Parameters, `TTL:`) {
			continue
		}
		if !strings.HasSuffix(job.Parameters, `h`) {
			continue
		}

		ttlHoursStr := job.Parameters
		ttlHoursStr = strings.TrimPrefix(ttlHoursStr, `TTL:`)
		ttlHoursStr = strings.TrimSuffix(ttlHoursStr, `h`)
		ttlHours, err := strconv.Atoi(ttlHoursStr)
		if err != nil {
			log.Errorf("job %+v has unexpected parameters ttl format, config generation skipping!\n", job)
			continue
		}

		ttl := time.Duration(ttlHours) * time.Hour
		if ttl > maxReval {
			ttl = maxReval
		} else if ttl < minTTL {
			ttl = minTTL
		}

		jobStartTime, err := time.Parse(tc.JobTimeFormat, job.StartTime)
		if err != nil {
			log.Errorf("job %+v has unexpected time format, config generation skipping!\n", job)
			continue
		}

		if jobStartTime.Add(maxReval).Before(time.Now()) {
			continue
		}

		if jobStartTime.Add(ttl).Before(time.Now()) {
			continue
		}
		if job.Keyword != JobKeywordPurge {
			continue
		}

		purgeEnd := jobStartTime.Add(ttl)

		if existingPurgeEnd, ok := jobMap[job.AssetURL]; !ok || purgeEnd.After(existingPurgeEnd) {
			jobMap[job.AssetURL] = purgeEnd
		}
	}

	newJobs := []Job{}
	for assetURL, purgeEnd := range jobMap {
		newJobs = append(newJobs, Job{AssetURL: assetURL, PurgeEnd: purgeEnd})
	}
	sort.Sort(Jobs(newJobs))

	return newJobs
}
