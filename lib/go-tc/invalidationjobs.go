package tc

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
	"database/sql"
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-log"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

// MaxTTL is the maximum value of TTL representable as a time.Duration object, which is used
// internally by InvalidationJobInput objects to store the TTL.
const MaxTTL = math.MaxInt64 / 3600000000000

const twoDays = time.Hour * 48

// ValidJobRegexPrefix matches the only valid prefixes for a relative-path
// Content Invalidation Job regular expression.
var ValidJobRegexPrefix = regexp.MustCompile(`^\?/.*$`)

// InvalidationJob represents a content invalidation job as returned by the API.
type InvalidationJob struct {
	AssetURL        *string `json:"assetUrl"`
	CreatedBy       *string `json:"createdBy"`
	DeliveryService *string `json:"deliveryService"`
	ID              *uint64 `json:"id"`
	Keyword         *string `json:"keyword"`
	Parameters      *string `json:"parameters"`

	// StartTime is the time at which the job will come into effect. Must be in the future, but will
	// fail to Validate if it is further in the future than two days.
	StartTime *Time `json:"startTime"`
}

// InvalidationJobsResponse is the type of a response from Traffic Ops to a
// request made to its /jobs API endpoint.
type InvalidationJobsResponse struct {
	Response []InvalidationJob `json:"response"`
	Alerts
}

// InvalidationJobInput represents user input intending to create or modify a content invalidation job.
type InvalidationJobInput struct {

	// DeliveryService needs to be an identifier for a Delivery Service. It can be either a string - in which
	// case it is treated as an XML_ID - or a float64 (because that's the type used by encoding/json
	// to represent all JSON numbers) - in which case it's treated as an integral, unique identifier
	// (and any fractional part is discarded, i.e. 2.34 -> 2)
	DeliveryService *interface{} `json:"deliveryService"`

	// Regex is a regular expression which not only must be valid, but should also start with '/'
	// (or escaped: '\/')
	Regex *string `json:"regex"`

	// StartTime is the time at which the job will come into effect. Must be in the future.
	StartTime *Time `json:"startTime"`

	// TTL indicates the Time-to-Live of the job. This can be either a valid string for
	// time.ParseDuration, or a float64 indicating the number of hours. Note that regardless of the
	// actual value here, Traffic Ops will only consider it rounded down to the nearest natural
	// number
	TTL *interface{} `json:"ttl"`

	dsid *uint
	ttl  *time.Duration
}

// UserInvalidationJobInput Represents legacy-style user input to the /user/current/jobs API endpoint.
// This is much less flexible than InvalidationJobInput, which should be used instead when possible.
type UserInvalidationJobInput struct {
	DSID  *uint   `json:"dsId"`
	Regex *string `json:"regex"`

	// StartTime is the time at which the job will come into effect. Must be in the future, but will
	// fail to Validate if it is further in the future than two days.
	StartTime *Time   `json:"startTime"`
	TTL       *uint64 `json:"ttl"`
	Urgent    *bool   `json:"urgent"`
}

// UserInvalidationJob is a full representation of content invalidation jobs as stored in the
// database, including several unused fields.
type UserInvalidationJob struct {

	// Agent is unused, and developers should never count on it containing or meaning anything.
	Agent    *uint   `json:"agent"`
	AssetURL *string `json:"assetUrl"`

	// AssetType is unused, and developers should never count on it containing or meaning anything.
	AssetType       *string `json:"assetType"`
	DeliveryService *string `json:"deliveryService"`
	EnteredTime     *Time   `json:"enteredTime"`
	ID              *uint   `json:"id"`
	Keyword         *string `json:"keyword"`

	// ObjectName is unused, and developers should never count on it containing or meaning anything.
	ObjectName *string `json:"objectName"`

	// ObjectType is unused, and developers should never count on it containing or meaning anything.
	ObjectType *string `json:"objectType"`
	Parameters *string `json:"parameters"`
	Username   *string `json:"username"`
}

// DSID gets the integral, unique identifier of the Delivery Service identified by
// InvalidationJobInput.DeliveryService
//
// This requires a transaction connected to a Traffic Ops database, because if DeliveryService is
// an xml_id, a database lookup will be necessary to get the unique, integral identifier. Thus,
// this method also checks for the existence of the identified Delivery Service, and will return
// an error if it does not exist.
func (j *InvalidationJobInput) DSID(tx *sql.Tx) (uint, error) {
	if j.dsid != nil {
		return *j.dsid, nil
	}

	if j.DeliveryService == nil {
		return 0, errors.New("Attempted to turn a nil DeliveryService into a DSID")
	}
	if tx == nil {
		return 0, errors.New("Attempted to turn a DeliveryService into a DSID with no DB connection")
	}

	var ret uint
	switch t := (*j.DeliveryService).(type) {
	case float64:
		v := (*j.DeliveryService).(float64)
		if v < 0 {
			return 0, errors.New("Delivery Service ID cannot be negative")
		}

		u := uint(v)
		var exists bool
		row := tx.QueryRow(`SELECT EXISTS(SELECT * FROM deliveryservice WHERE id=$1)`, u)
		if err := row.Scan(&exists); err != nil {
			log.Errorf("Error checking for deliveryservice existence in DSID: %v\n", err)
			return 0, errors.New("Unknown error occurred")
		} else if !exists {
			return 0, fmt.Errorf("No Delivery Service exists matching identifier: %v", *j.DeliveryService)
		}

		j.dsid = &u
		return u, nil

	case string:
		row := tx.QueryRow(`SELECT id FROM deliveryservice WHERE xml_id=$1`, *j.DeliveryService)
		if err := row.Scan(&ret); err != nil {
			if err == sql.ErrNoRows {
				return 0, fmt.Errorf("No DeliveryService exists matching identifier: %v", *j.DeliveryService)
			}
			return 0, errors.New("Unknown error occurred")
		}
		j.dsid = &ret
		return ret, nil

	default:
		log.Errorf("unsupported DS key type: %T\n", t)
		return 0, errors.New("Unknown error occurred")

	}
}

// Validate validates that the user input is correct, given a transaction
// connected to the Traffic Ops database. In particular, it enforces the
// constraints described on each field, as well as ensuring they actually exist.
// This method calls InvalidationJobInput.DSID to validate the DeliveryService
// field.
//
// This returns an error describing any and all problematic fields encountered during validation.
func (job *InvalidationJobInput) Validate(tx *sql.Tx) error {
	errs := []string{}
	err := validation.ValidateStruct(job,
		validation.Field(&job.DeliveryService, validation.Required),
		validation.Field(&job.Regex, validation.Required, validation.NewStringRule(func(s string) bool {
			return strings.HasPrefix(s, `\/`) || strings.HasPrefix(s, "/")
		}, `must start with '/' (or '\/')`)),
		validation.Field(&job.TTL, validation.Required),
	)

	if err != nil {
		errs = append(errs, err.Error())
	}

	if job.DeliveryService != nil {
		if _, err = job.DSID(tx); err != nil {
			errs = append(errs, err.Error())
		}
	}

	if job.Regex != nil && *job.Regex != "" {
		if _, err := regexp.Compile(*job.Regex); err != nil {
			errs = append(errs, "regex: is not a valid Regular Expression: "+err.Error())
		}
	}

	if job.StartTime == nil {
		errs = append(errs, "startTime: cannot be blank")
	} else if job.StartTime.Time.Before(time.Now()) {
		errs = append(errs, "startTime: must be in the future")
	}

	if job.TTL != nil {
		hours, err := job.TTLHours()
		if err != nil {
			errs = append(errs, "ttl: must be a number of hours, or a duration string e.g. '48h'")
		}
		var maxDays uint
		err = tx.QueryRow(`SELECT value FROM parameter WHERE name='maxRevalDurationDays' AND config_file='regex_revalidate.config'`).Scan(&maxDays)
		maxHours := maxDays * 24
		if err == nil && hours > maxHours { // silently ignore other errors too
			errs = append(errs, "ttl: cannot exceed "+strconv.FormatUint(uint64(maxHours), 10)+"!")
		}
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, ", "))
	}

	return nil
}

type compareJob struct {
	AssetURL  string
	TTLHours  uint
	StartTime time.Time
}

// ValidateJobUniqueness returns a message describing each overlap between
// existing content invalidation jobs for the same assetURL as the one passed.
//
// TODO: This doesn't belong in the lib, and it swallows errors because it
// can't log them.
func ValidateJobUniqueness(tx *sql.Tx, dsID uint, startTime time.Time, assetURL string, ttlHours uint) []string {
	var errs []string

	const readQuery = `
SELECT asset_url,
	   ttl_hr,
       start_time
FROM job
WHERE job.job_deliveryservice = $1
`
	rows, err := tx.Query(readQuery, dsID)
	if err != nil {
		errs = append(errs, "unable to query for invalidation jobs while validating job uniqueness")
	} else {
		defer rows.Close()
		jobStart := startTime
		for rows.Next() {
			testJob := compareJob{}
			err = rows.Scan(
				&testJob.AssetURL,
				&testJob.TTLHours,
				&testJob.StartTime)
			if err != nil {
				continue
			}
			if !strings.HasSuffix(testJob.AssetURL, assetURL) {
				continue
			}
			if testJob.TTLHours == 0 {
				continue
			}
			testJobStart := testJob.StartTime
			testJobEnd := testJobStart.Add(time.Hour * time.Duration(testJob.TTLHours))
			jobEnd := jobStart.Add(time.Hour * time.Duration(ttlHours))
			// jobStart in testJob range
			if (testJobStart.Before(jobStart) && jobStart.Before(testJobEnd)) ||
				// jobEnd in testJob range
				(testJobStart.Before(jobEnd) && jobEnd.Before(testJobEnd)) ||
				// job range encaspulates testJob range
				(testJobEnd.Before(jobEnd) && jobStart.Before(jobStart)) {
				errs = append(errs, fmt.Sprintf("Invalidation request duplicate found for %v, start:%v end:%v",
					testJob.AssetURL, testJobStart, testJobEnd))
			}
		}
	}

	return errs
}

// TTLHours gets the number of hours of the job's TTL - rounded down to the nearest natural number,
// or an error if it is an invalid value.
func (j *InvalidationJobInput) TTLHours() (uint, error) {
	if j.ttl != nil {
		return uint((*j.ttl).Hours()), nil
	}
	if j.TTL == nil {
		return 0, errors.New("Attempted to convert a nil TTL into hours")
	}

	var ret uint
	switch t := (*j.TTL).(type) {
	case float64:
		v := (*j.TTL).(float64)
		if v < 0 {
			return 0, errors.New("TTL cannot be negative!")
		}
		if v >= MaxTTL {
			return 0, fmt.Errorf("TTL cannot exceed %d hours!", MaxTTL)
		}
		ttl := time.Duration(int64(v * 3600000000000))
		j.ttl = &ttl
		ret = uint(ttl.Hours())

	case string:
		d, err := time.ParseDuration((*j.TTL).(string))
		if err != nil || d.Hours() < 1 {
			return 0, fmt.Errorf("Invalid duration entered for TTL! Must be at least one hour, but no more than %d hours!", MaxTTL)
		}
		j.ttl = &d
		ret = uint(d.Hours())

	default:
		log.Errorf("unsupported TTL key type: %T\n", t)
		return 0, errors.New("Unknown error occurred")
	}

	return ret, nil
}

// TTLHours will parse job.Parameters to find TTL, returns an int representing
// number of hours. Returns 0 in case of issue (0 is an invalid TTL).
func (job *InvalidationJob) TTLHours() uint {
	if job.Parameters == nil {
		return 0
	}
	ttl := strings.Split(*job.Parameters, ":")
	if len(ttl) != 2 {
		return 0
	}

	hours, err := strconv.Atoi(ttl[1][:len(ttl[1])-1])
	if err != nil {
		return 0
	}
	return uint(hours)
}

// Validate checks that the InvalidationJob is valid, by ensuring all of its fields are well-defined.
//
// This returns an error describing any and all problematic fields encountered during validation.
func (job *InvalidationJob) Validate() error {
	errs := []string{}
	err := validation.ValidateStruct(job,
		validation.Field(&job.AssetURL, validation.Required, is.URL),
		validation.Field(&job.CreatedBy, validation.Required),
		validation.Field(&job.DeliveryService, validation.Required),
		validation.Field(&job.ID, validation.Required),
		validation.Field(&job.Keyword, validation.Required),
		validation.Field(&job.Parameters, validation.Required),
	)

	if err != nil {
		errs = append(errs, err.Error())
	}

	if job.StartTime == nil {
		return errors.New(strings.Join(append(errs, "startTime: cannot be blank"), ", "))
	}

	if job.StartTime.After(time.Now().Add(twoDays)) {
		errs = append(errs, "startTime: must be within two days from now")
	}

	if job.StartTime.Before(time.Now()) {
		errs = append(errs, "startTime: cannot be in the past")
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, ", "))
	}

	return nil
}

// Validate validates that the user input is correct, given a transaction
// connected to the Traffic Ops database.
//
// This requires a database transaction to check that the DSID is a valid
// identifier for an existing Delivery Service.
//
// Returns an error describing any and all problematic fields encountered during
// validation.
func (job *UserInvalidationJobInput) Validate(tx *sql.Tx) error {
	errs := []string{}
	err := validation.ValidateStruct(job,
		validation.Field(&job.Regex, validation.Required, validation.NewStringRule(func(s string) bool {
			return strings.HasPrefix(s, `\/`) || strings.HasPrefix(s, "/")
		}, `must start with '/' (or '\/')`)),
		validation.Field(&job.DSID, validation.Required),
		validation.Field(&job.TTL, validation.Required),
	)
	if err != nil {
		errs = append(errs, err.Error())
	}

	if job.StartTime == nil {
		errs = append(errs, "startTime: cannot be blank")
	} else if job.StartTime.After(time.Now().Add(twoDays)) {
		errs = append(errs, "startTime: must be within two days")
	}

	if job.Regex != nil && *(job.Regex) != "" {
		if _, err := regexp.Compile(*(job.Regex)); err != nil {
			errs = append(errs, "regex: is not a valid regular expression: "+err.Error())
		}
	}

	if job.DSID != nil {
		row := tx.QueryRow(`SELECT id FROM deliveryservice WHERE id = $1::bigint`, job.DSID)
		var id uint
		if err := row.Scan(&id); err != nil {
			log.Errorln(err.Error())
			errs = append(errs, "no Delivery Service corresponding to 'dsId'")
		}
	}

	if job.TTL != nil {
		row := tx.QueryRow(`SELECT value FROM parameter WHERE name='maxRevalDurationDays' AND config_file='regex_revalidate.config'`)
		var maxDays uint64
		err := row.Scan(&maxDays)
		maxHours := maxDays * 24
		if err == sql.ErrNoRows && MaxTTL < *(job.TTL) {
			errs = append(errs, "ttl: cannot exceed "+strconv.FormatUint(MaxTTL, 10))
		} else if err == nil && maxHours < *(job.TTL) { // silently ignore other errors
			errs = append(errs, "ttl: cannot exceed "+strconv.FormatUint(maxHours, 10))
		} else if *(job.TTL) < 1 {
			errs = append(errs, "ttl: must be at least 1")
		}
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, ", "))
	}
	return nil
}

// These are the allowed values for the InvalidationType of an
// InvalidationJobCreateV4/InvalidationJobV4.
const (
	REFRESH = "REFRESH"
	REFETCH = "REFETCH"
)

// InvalidationJobsResponseV4 is the type of a response from Traffic Ops to a
// request made to its /jobs API endpoint for API major version 4.
type InvalidationJobsResponseV4 struct {
	Response []InvalidationJobV4 `json:"response"`
	Alerts
}

// InvalidationJobCreateV4 is an alias for the InvalidationJobCreateV40 struct used for the latest minor version associated with api major version 4.
type InvalidationJobCreateV4 InvalidationJobCreateV40

// InvalidationJobCreateV40 represents user input intending to create a content invalidation job.
type InvalidationJobCreateV40 struct {
	// The Delivery Service XML-ID for which the Invalidation Job is to be applied.
	DeliveryService string `json:"deliveryService"`

	// Regex is a regular expression which not only must be valid, but should also start with '/'
	// (or escaped: '\/')
	Regex string `json:"regex"`

	// StartTime is the time at which the job will come into effect. Must be in the future.
	StartTime time.Time `json:"startTime"`

	// TTLHours indicates the Time-to-Live of the job in hours. Must be a positive integer value.
	TTLHours uint32 `json:"ttlHours"`

	// InvalidationType must be either REFRESH (default behavior) or REFETCH. If REFETCH, must
	// also comply with global parameter setting
	InvalidationType string `json:"invalidationType"`
}

// InvalidationJobV4 is an alias for the InvalidationJobV4 struct used for the latest minor version associated with api major version 4.
type InvalidationJobV4 InvalidationJobV40

// InvalidationJobV40 represents a content invalidation job as returned by the API.
// Also used for Update calls.
type InvalidationJobV40 struct {
	ID               uint64    `json:"id"`
	AssetURL         string    `json:"assetUrl"`
	CreatedBy        string    `json:"createdBy"`
	DeliveryService  string    `json:"deliveryService"`
	TTLHours         uint      `json:"ttlHours"`
	InvalidationType string    `json:"invalidationType"`
	StartTime        time.Time `json:"startTime"`
}

// String implements the fmt.Stringer interface by providing a textual
// representation of the InvalidationJobV4.
func (job InvalidationJobV4) String() string {
	return fmt.Sprintf(`InvalidationJobV4{ID: %d, AssetURL: "%s", CreatedBy: "%s", DeliveryService: "%s", TTLHours: %d, InvalidationType: "%s", StartTime: "%s"}`,
		job.ID,
		job.AssetURL,
		job.CreatedBy,
		job.DeliveryService,
		job.TTLHours,
		job.InvalidationType,
		job.StartTime.Format(time.RFC3339),
	)
}
