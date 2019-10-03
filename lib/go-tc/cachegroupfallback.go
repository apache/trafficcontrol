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

// A List of CacheGroupFallbacks Response
// swagger:response CacheGroupFallbacksResponse
// in: body
type CacheGroupFallbacksResponse struct {
	// in: body
	Response []CacheGroupFallback `json:"response"`
}

// A Single CacheGroupFallback Response for Update and Create to depict what changed
// swagger:response CacheGroupFallbackResponse
// in: body
type CacheGroupFallbackResponse struct {
	// in: body
	Response CacheGroupFallback `json:"response"`
}

// CacheGroupFallback ...
type CacheGroupFallback struct {
	PrimaryCgId int `json:"primaryId" db:"primary_cg"`
	BackupCgId  int `json:"backupId" db:"backup_cg"`
	SetOrder    int `json:"setOrder" db:"set_order"`
}

// CacheGroupFallbackNullable ...
type CacheGroupFallbackNullable struct {
	PrimaryCgId *int `json:"primaryId" db:"primary_cg"`
	BackupCgId  *int `json:"backupId" db:"backup_cg"`
	SetOrder    *int `json:"setOrder" db:"set_order"`
}
