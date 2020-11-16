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
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

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

const ContentTypeRegexRevalidateDotConfig = ContentTypeTextASCII
const LineCommentRegexRevalidateDotConfig = LineCommentHash

func MakeRegexRevalidateDotConfig(
	server *Server,
	deliveryServices []DeliveryService,
	globalParams []tc.Parameter,
	jobs []tc.Job,
	hdrComment string,
) (Cfg, error) {
	warnings := []string{}

	if server.CDNName == nil {
		return Cfg{}, makeErr(warnings, "server CDNName missing")
	}

	params := paramsToMultiMap(filterParams(globalParams, RegexRevalidateFileName, "", "", ""))

	dsNames := map[string]struct{}{}
	for _, ds := range deliveryServices {
		if ds.XMLID == nil {
			warnings = append(warnings, "got Delivery Service from Traffic Ops with a nil xmlId! Skipping!")
			continue
		}
		dsNames[*ds.XMLID] = struct{}{}
	}

	dsJobs := []tc.Job{}
	for _, job := range jobs {
		if _, ok := dsNames[job.DeliveryService]; !ok {
			continue
		}
		dsJobs = append(dsJobs, job)
	}

	// TODO: add cdn, startTime query params to /jobs endpoint

	maxDays := DefaultMaxRevalDurationDays
	if maxDaysStrs := params[RegexRevalidateMaxRevalDurationDaysParamName]; len(maxDaysStrs) > 0 {
		sort.Strings(maxDaysStrs)
		err := error(nil)
		if maxDays, err = strconv.Atoi(maxDaysStrs[0]); err != nil { // just use the first, if there were multiple params
			warnings = append(warnings, "max days param '"+maxDaysStrs[0]+"' is not an integer, using default value!")
			maxDays = DefaultMaxRevalDurationDays
		}
	}

	maxReval := time.Duration(maxDays) * time.Hour * 24

	cfgJobs, jobWarns := filterJobs(dsJobs, maxReval, RegexRevalidateMinTTL)
	warnings = append(warnings, jobWarns...)

	txt := makeHdrComment(hdrComment)
	for _, job := range cfgJobs {
		txt += job.AssetURL + " " + strconv.FormatInt(job.PurgeEnd.Unix(), 10) + "\n"
	}

	return Cfg{
		Text:        txt,
		ContentType: ContentTypeRegexRevalidateDotConfig,
		LineComment: LineCommentRegexRevalidateDotConfig,
		Warnings:    warnings,
	}, nil
}

type job struct {
	AssetURL string
	PurgeEnd time.Time
}

type jobsSort []job

func (jb jobsSort) Len() int      { return len(jb) }
func (jb jobsSort) Swap(i, j int) { jb[i], jb[j] = jb[j], jb[i] }
func (jb jobsSort) Less(i, j int) bool {
	if jb[i].AssetURL == jb[j].AssetURL {
		return jb[i].PurgeEnd.Before(jb[j].PurgeEnd)
	}
	return strings.Compare(jb[i].AssetURL, jb[j].AssetURL) < 0
}

// filterJobs returns only jobs which:
//   - have a non-null deliveryservice
//   - have parameters of the form TTL:%dh
//   - have a start time later than (now + maxReval days). That is, we don't query jobs older than maxReval in the past.
//   - are "purge" jobs
//   - have a start_time+ttl > now. That is, jobs that haven't expired yet.
// Returns the filtered jobs, and any warnings.
func filterJobs(jobs []tc.Job, maxReval time.Duration, minTTL time.Duration) ([]job, []string) {
	warnings := []string{}

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
			warnings = append(warnings, fmt.Sprintf("job %+v has unexpected parameters ttl format, config generation skipping!\n", job))
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
			warnings = append(warnings, fmt.Sprintf("job %+v has unexpected time format, config generation skipping!\n", job))
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

	newJobs := []job{}
	for assetURL, purgeEnd := range jobMap {
		newJobs = append(newJobs, job{AssetURL: assetURL, PurgeEnd: purgeEnd})
	}
	sort.Sort(jobsSort(newJobs))

	return newJobs, warnings
}
