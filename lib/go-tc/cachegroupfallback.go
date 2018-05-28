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

// A List of cachegroupFallbacks Response
// swagger:response cachegroupFallbacksResponse
// in: body
type cachegroupFallbacksResponse struct {
	// in: body
	Response []cachegroupFallback `json:"response"`
}

// A Single cachegroupFallback Response for Update and Create to depict what changed
// swagger:response cachegroupFallbackResponse
// in: body
type cachegroupFallbackResponse struct {
	// in: body
	Response cachegroupFallback `json:"response"`
}

// cachegroupFallback ...
type cachegroupFallback struct {

	PrimaryCgId int `json:"primaryId" db:"primary_cg"`
	BackupCgId  int `json:"backupId" db:"backup_cg"`
	SetOrder    int `json:"setOrder" db:"set_order"`

}

// cachegroupFallbackNullable ...
type cachegroupFallbackNullable struct {

	PrimaryCgId *int `json:"primaryId" db:"primary_cg"`
	BackupCgId  *int `json:"backupId" db:"backup_cg"`
	SetOrder    *int `json:"setOrder" db:"set_order"`
}

