package staticdnsentry

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
	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"
)

type TOStaticDNSEntry struct{
	ReqInfo *api.APIInfo `json:"-"`
	tc.StaticDNSEntry
}

func GetReaderSingleton() func(reqInfo *api.APIInfo)api.Reader {
	return func(reqInfo *api.APIInfo)api.Reader {
		toReturn := TOStaticDNSEntry{reqInfo, tc.StaticDNSEntry{}}
		return &toReturn
	}
}

func (staticDNSEntry *TOStaticDNSEntry) Read(parameters map[string]string) ([]interface{}, []error, tc.ApiErrorType) {
	queryParamsToQueryCols := map[string]dbhelpers.WhereColumnInfo{
		"deliveryservice": dbhelpers.WhereColumnInfo{"deliveryservice", nil}, // order by
	}
	where, orderBy, queryValues, errs := dbhelpers.BuildWhereAndOrderBy(parameters, queryParamsToQueryCols)
	if len(errs) > 0 {
		log.Errorf("Data Conflict Error")
		return nil, errs, tc.DataConflictError
	}
	query := selectQuery() + where + orderBy
	log.Debugln("Query is ", query)
	rows, err := staticDNSEntry.ReqInfo.Tx.NamedQuery(query, queryValues)
	if err != nil {
		log.Errorf("Error querying StaticDNSEntries: %v", err)
		return nil, []error{tc.DBError}, tc.SystemError
	}
	defer rows.Close()
	staticDNSEntries := []interface{}{}
	for rows.Next() {
		s := tc.StaticDNSEntry{}
		if err = rows.StructScan(&s); err != nil {
			log.Errorln("error parsing StaticDNSEntry rows: " + err.Error())
			return nil, []error{tc.DBError}, tc.SystemError
		}
		staticDNSEntries = append(staticDNSEntries, s)
	}
	return staticDNSEntries, []error{}, tc.NoError
}

func selectQuery() string {
	return `
SELECT
ds.xml_id as dsname,
sde.host,
sde.ttl,
sde.address,
tp.name as type,
cg.name as cachegroup
FROM staticdnsentry as sde
JOIN type as tp on sde.type = tp.id
JOIN cachegroup as cg ON sde.cachegroup = cg.id
JOIN deliveryservice as ds on sde.deliveryservice = ds.id
`
}
