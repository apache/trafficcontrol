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
  "net/http"
  //"database/sql"
  //"github.com/jmoiron/sqlx" //sql extra
  "github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
)

func DomainsHandler (w http.ResponseWriter, r *http.Request) {

  // inf is of type APIInfo
  inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
  if userErr != nil || sysErr != nil {
    api.HandleErr(w, r, errCode, userErr, sysErr)
    return
  }
  defer inf.Close()

  var (
    cdn         int
    id          int
    name        string
    description string
    domain_name string
    resp []interface{} //really not sure about this
  )

  //how to prefetch 'cdn'?
  q := `SELECT cdn, id, name, description FROM 'Profile' WHERE name LIKE 'CCR%'`
  rows, err := inf.Tx.Query(q)
  if err != nil {
    //TODO change errCode, userErr, and sysErr
    //which one is errors.New("Error: " + err.Error()?
    //api.HandleErr(w, r, errCode, userErr, sysErr)
    return
  }
  defer rows.Close()

  for rows.Next() {

    //Do I even need to check for errors here since the perl doesn't?
    if err := rows.Scan(&cdn, &id, &name, &description); err != nil {
      //api.HandleErr(w, r, HTTP CODE,
      //  errors.New("Error scanning ...: " + err.Error())  user or system error?
    }

    err = inf.Tx.QueryRow("SELECT DOMAIN_NAME FROM CDN WHERE id = $1", 1).Scan(&domain_name)
    if err != nil {
      //api.HandleErr(w, r, HTTP CODE,
      //  errors.New("Error scanning ...: " + err.Error()") user or system error?
    }

    data := struct {
      domain_name string
      param_id    int
      id          int
      name        string
      description string
    } {
      domain_name,
      -1, // it's not a parameter anymore
      id,
      name,
      description,
    }
    resp = append(resp, data)
  }

  api.WriteResp(w, r, resp)
  /* {
    "domainName"         => $row->cdn->domain_name,
    "parameterId"        => -1,  # it's not a parameter anymore
    "profileId"          => $row->id,
    "profileName"        => $row->name,
    "profileDescription" => $row->description,
  } */
}
