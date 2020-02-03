package ats

import (
	"database/sql"
	"errors"
	"github.com/apache/trafficcontrol/lib/go-atscfg"
)

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

func GetNameVersionString(tx *sql.Tx) (string, error) {
	toolName, url, err := GetToolNameAndURL(tx)
	if err != nil {
		return "", errors.New("getting toolname and url parameters: " + err.Error())
	}
	return atscfg.GetNameVersionStringFromToolNameAndURL(toolName, url), nil
}

func HeaderComment(tx *sql.Tx, name string) (string, error) {
	nameVersionStr, err := GetNameVersionString(tx)
	if err != nil {
		return "", errors.New("getting name version string: " + err.Error())
	}
	return atscfg.HeaderCommentWithTOVersionStr(name, nameVersionStr), nil
}
