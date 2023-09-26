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

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
)

// CacheUrlPrefix is an unused constant of unknown purpose.
//
// Deprecated: It's likely that this is related to the removed "Cache URL"
// property of Delivery Services and the old "cacheurl" plugin for ATS which is
// not supported by any supported version of ATS.
const CacheUrlPrefix = "cacheurl_"

// RefetchSuffix is a special suffix which, if present at the end of an asset
// URL for a Content Invalidation Job, will cause ATS to treat matching URIs as
// cache "misses" - effectively this forces ATS to fetch the asset anew from an
// upstream cache server or origin.
const RefetchSuffix = "##REFETCH##"

// RefreshSuffix is a special suffix which, if present at the end of an asset
// URL for a Content Invalidation Job, will cause ATS to treat matching URIs as
// "stale" - effectively this forces ATS to retrieve confirmation from an
// upstream cache server or origin that the cached content is still valid, and
// according to their response either fetch the content anew or update its age
// and caching validity markers (ETag, If-Unmodified-Since etc.).
const RefreshSuffix = "##REFRESH##"

// RemapFile is the name of a configuration file used by ATS to map client
// requests.
//
// TODO: replace instances of "remap.config" with this constant.
const RemapFile = "remap.config"

// RegexRevalidateFileName is the name of a configuration file used by ATS to
// revalidate or invalidate cached contents based on regular expression matches.
// This is also the ConfigFile value of GLOBAL Parameters that modify how
// content invalidation is handled.
const RegexRevalidateFileName = "regex_revalidate.config"

// RegexRevalidateMaxRevalDurationDaysParamName is the Name of a GLOBAL
// Parameter that can be used to set the maximum duration (in days) of Content
// Invalidation Requests.
const RegexRevalidateMaxRevalDurationDaysParamName = "maxRevalDurationDays"

// DefaultMaxRevalDurationDays is the maximum duration (in days) of Content
// Invalidation Requests if not overridden by a Parameter.
const DefaultMaxRevalDurationDays = 90

// JobKeywordPurge is the "keyword" of "jobs" that determines if a "job"
// represents a Content Invalidation Job - "jobs" with other "keywords" are
// ignored.
//
// TODO: Figure out if this can be deprecated/removed, along with the check for
// it; all "jobs" created through the API have this "keyword", and the notion of
// generic-purpose "jobs" has been abandoned.
const JobKeywordPurge = "PURGE"

// RegexRevalidateMinTTL is a time.Duration that specifies the minimum duration
// for which Content Invalidation Jobs will remain active. Content Invalidation
// Jobs with TTLs less than this will use this duration instead.
const RegexRevalidateMinTTL = time.Hour

// A RevalType describes the type of revalidation to be performed by a content
// invalidation job.
type RevalType string

// RevalTypeMiss represents the caching condition where the asset requested
// could not be found in cache.
const RevalTypeMiss = RevalType("MISS")

// RevalTypeStale represents the caching condition where the asset requested is
// in cache, but it has expired and must be revalidated upstream.
const RevalTypeStale = RevalType("STALE")

// RevalTypeDefault is the type of revalidation that will be performed by
// default, if not otherwise specified.
const RevalTypeDefault = RevalTypeStale

// ContentTypeRegexRevalidateDotConfig is the MIME content type of the contents
// of a regex_revalidate.config ATS configuration file.
const ContentTypeRegexRevalidateDotConfig = ContentTypeTextASCII

// LineCommentRegexRevalidateDotConfig is the string that indicates the
// beginning of a line comment in the grammar of a regex_revalidate.config ATS
// configuration file.
const LineCommentRegexRevalidateDotConfig = LineCommentHash

// RegexRevalidateDotConfigOpts contains settings to configure generation options.
type RegexRevalidateDotConfigOpts struct {
	// HdrComment is the header comment to include at the beginning of the file.
	// This should be the text desired, without comment syntax (like # or //). The file's comment syntax will be added.
	// To omit the header comment, pass the empty string.
	HdrComment string
}

// MakeRegexRevalidateDotConfig constructs a regex_revalidate.config file for
// the given server, which is responsible for serving content for the given
// Delivery Services, with the given set of "global" Parameters (NOT server or
// Delivery Service Profile Parameters), and the given set of Content
// Invalidation Jobs.
func MakeRegexRevalidateDotConfig(
	server *Server,
	deliveryServices []DeliveryService,
	globalParams []tc.ParameterV5,
	jobs []InvalidationJob,
	opt *RegexRevalidateDotConfigOpts,
) (Cfg, error) {
	if opt == nil {
		opt = &RegexRevalidateDotConfigOpts{}
	}
	warnings := []string{}

	if server.CDN == "" {
		return Cfg{}, makeErr(warnings, "server CDNName missing")
	}

	params := paramsToMultiMap(filterParams(globalParams, RegexRevalidateFileName, "", "", ""))

	dsNames := map[string]struct{}{}
	for _, ds := range deliveryServices {
		if ds.XMLID == "" {
			warnings = append(warnings, "got Delivery Service from Traffic Ops with a nil xmlId! Skipping!")
			continue
		}
		dsNames[ds.XMLID] = struct{}{}
	}

	dsJobs := []InvalidationJob{}
	for _, job := range jobs {
		if job.DeliveryService == "" {
			warnings = append(warnings, "got job from Traffic Ops with an empty DeliveryService! Skipping!")
			continue
		}
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

	cfgJobs := filterJobs(dsJobs, maxReval, RegexRevalidateMinTTL)

	txt := makeHdrComment(opt.HdrComment)
	for _, job := range cfgJobs {
		txt += job.AssetURL + " " + strconv.FormatInt(job.PurgeEnd.Unix(), 10)
		if job.Type != "" && job.Type != RevalTypeDefault {
			txt += " " + string(job.Type)
		}
		txt += "\n"
	}

	return Cfg{
		Text:        txt,
		ContentType: ContentTypeRegexRevalidateDotConfig,
		LineComment: LineCommentRegexRevalidateDotConfig,
		Warnings:    warnings,
	}, nil
}

type revalJob struct {
	AssetURL string
	PurgeEnd time.Time
	Type     RevalType // RevalTypeMiss or RevalTypeStale (default)
}

type jobsSort []revalJob

func (jb jobsSort) Len() int      { return len(jb) }
func (jb jobsSort) Swap(i, j int) { jb[i], jb[j] = jb[j], jb[i] }
func (jb jobsSort) Less(i, j int) bool {
	if jb[i].AssetURL == jb[j].AssetURL {
		return jb[i].PurgeEnd.Before(jb[j].PurgeEnd)
	}
	return strings.Compare(jb[i].AssetURL, jb[j].AssetURL) < 0
}

// filterJobs returns only jobs which:
//   - have a non-empty deliveryservice
//   - have a start time later than (now + maxReval days). That is, we don't query jobs older than maxReval in the past.
//   - have a start_time+ttl > now. That is, jobs that haven't expired yet.
//
// Returns the filtered jobs.
func filterJobs(tcJobs []InvalidationJob, maxReval time.Duration, minTTL time.Duration) []revalJob {

	jobMap := map[string]revalJob{}

	for _, tcJob := range tcJobs {
		if tcJob.DeliveryService == "" {
			continue
		}

		ttl := time.Duration(tcJob.TTLHours) * time.Hour
		if ttl > maxReval {
			ttl = maxReval
		} else if ttl < minTTL {
			ttl = minTTL
		}

		if tcJob.StartTime.Add(maxReval).Before(time.Now()) {
			continue
		}

		if tcJob.StartTime.Add(ttl).Before(time.Now()) {
			continue
		}

		jobType, assetURL := processRefetch(tcJob.InvalidationType, tcJob.AssetURL)

		purgeEnd := tcJob.StartTime.Add(ttl)

		if rjob, ok := jobMap[assetURL]; !ok || purgeEnd.After(rjob.PurgeEnd) {
			jobMap[assetURL] = revalJob{AssetURL: assetURL, PurgeEnd: purgeEnd, Type: jobType}
		}
	}

	newJobs := []revalJob{}
	for _, rjob := range jobMap {
		newJobs = append(newJobs, rjob)
	}
	sort.Sort(jobsSort(newJobs))

	return newJobs
}

// processRefetch determines the type of Invalidation, returns the corresponding
// jobtype and "cleans" the regex URL for the asset to be invalidated. REFETCH
// trumps REFRESH, whether in the AssetURL or as InvalidationType.
func processRefetch(invalidationType, assetURL string) (RevalType, string) {
	if (len(invalidationType) > 0 && invalidationType == tc.REFETCH) || strings.HasSuffix(assetURL, RefetchSuffix) {
		assetURL = strings.TrimSuffix(assetURL, RefetchSuffix)
		return RevalTypeMiss, assetURL
	}

	// Default value. Either the InvalidationType == REFRESH
	// or the suffix is ##REFRESH## or neither
	assetURL = strings.TrimSuffix(assetURL, RefreshSuffix)
	return RevalTypeStale, assetURL
}
