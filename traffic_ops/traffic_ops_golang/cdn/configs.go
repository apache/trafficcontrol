package cdn

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
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"
	"net/http"
)

// TOCDNConf used as a type alias to define functions on to satisfy shared API REST interfaces.
type TOCDNConf struct {
	api.APIInfoImpl `json:"-"`
}

func (v *TOCDNConf) NewReadObj() interface{} { return &tc.CDNConfig{} }
func (v *TOCDNConf) SelectQuery() string     { return cdnConfSelectQuery() }
func (v *TOCDNConf) ParamColumns() map[string]dbhelpers.WhereColumnInfo {
	return map[string]dbhelpers.WhereColumnInfo{}
}

func cdnConfSelectQuery() string {
	return `SELECT
name,
id
FROM cdn`
}

func (v *TOCDNConf) Read(http.Header) ([]interface{}, error, error, int) {
	return api.GenericRead(nil, v)
}
func (v *TOCDNConf) SelectMaxLastUpdatedQuery(string, string, string, string, string, string) string {
	return ""
}                                                   //{ return selectMaxLastUpdatedQuery() }
func (v *TOCDNConf) InsertIntoDeletedQuery() string { return "" } //{return insertIntoDeletedQuery()}
func (v TOCDNConf) GetType() string {
	return "cdn_configs"
}
