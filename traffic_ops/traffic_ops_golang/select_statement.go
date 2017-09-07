package main

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

import "net/url"

const (
	EQUAL     = "="
	NOT_EQUAL = "!="
	OR        = "OR"
)

type Condition struct {
	Key     string
	Operand string
	Value   string
}

type SelectStatement struct {
	Select string
	Where  WhereClause
}

func (q *SelectStatement) String() string {
	if q.Where.Exists() {
		return q.Select + q.Where.String()
	} else {
		return q.Select
	}
}

type WhereClause struct {
	Condition Condition
}

func (w *WhereClause) SetCondition(c Condition) Condition {
	w.Condition = c
	return w.Condition
}

func (w *WhereClause) String() string {
	c := w.Condition
	return "\nWHERE " + c.Key + c.Operand + "$1"
}

func (w *WhereClause) Exists() bool {
	if (Condition{}) != w.Condition {
		return true
	} else {
		return false
	}
}

func newWhereClause(v url.Values) WhereClause {

	whereClause := WhereClause{}

	switch {
	case v.Get("cachegroup") != "":
		whereClause.SetCondition(Condition{"s.cachegroup", EQUAL, v.Get("cachegroup")})

	// Support what should have been the cachegroupId as well
	case v.Get("cachegroupId") != "":
		whereClause.SetCondition(Condition{"s.cachegroup", EQUAL, v.Get("cachegroupId")})

	case v.Get("cdn") != "":
		whereClause.SetCondition(Condition{"s.cdn_id", EQUAL, v.Get("cdn")})

	case v.Get("physLocation") != "":
		whereClause.SetCondition(Condition{"s.phys_location", EQUAL, v.Get("physLocation")})

	case v.Get("physLocationId") != "":
		whereClause.SetCondition(Condition{"s.phys_location", EQUAL, v.Get("physLocationId")})

	case v.Get("profileId") != "":
		whereClause.SetCondition(Condition{"s.profile", EQUAL, v.Get("profileId")})

	case v.Get("type") != "":
		whereClause.SetCondition(Condition{"s.type", EQUAL, v.Get("type")})

	case v.Get("typeId") != "":
		whereClause.SetCondition(Condition{"s.type", EQUAL, v.Get("typeId")})
	}
	return whereClause
}
