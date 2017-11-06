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

type PhysLocationsResponse struct {
	Response []PhysLocation `json:"response"`
}

type PhysLocation struct {
	ID           int    `json:"id" db:"id"`
	Name         string `json:"name" db:"name"`
	ShortName    string `json:"short_name" db:"short_name"`
	Address      string `json:"address" db:"address"`
	City         string `json:"city" db:"city"`
	State        string `json:"state" db:"state"`
	Zip          string `json:"zip" db:"zip"`
	RegionId     int    `json:"regionId" db:"regionid"`
	POC          string `json:"poc" db:"poc"`
	Phone        string `json:"phone" db:"phone"`
	Email        string `json:"email" db:"email"`
	Comments     string `json:"comments" db:"comments"`
	RegionName   string `json:"regionName" db:"regionname"`
	LastUpdated  Time   `json:"lastUpdated" db:"last_updated"`

}
