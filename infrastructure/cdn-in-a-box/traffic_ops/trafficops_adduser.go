package main

// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

import (
	"fmt"
	"os"
	"strings"
	"text/template"
	"unicode"

	"auth"
)

const insertTemplate = `
INSERT INTO tm_user (username, role, local_passwd, confirm_local_passwd)
  VALUES ('{{ .Username }}',
          (SELECT id FROM role WHERE name='{{ .RoleName }}'),
          '{{ .EncryptedPassword }}',
          '{{ .EncryptedPassword }}')
`

// User contains only required fields to be stored in the db
type User struct {
	Username          string
	Password          string
	EncryptedPassword string
	RoleName          string
}

var (
	userEnv     = "TO_ADMIN_USER"
	passwordEnv = "TO_ADMIN_PASSWORD"
)

func main() {
	user := User{Username: "admin", RoleName: "admin"}
	v := os.Getenv("TO_ADMIN_USER")
	if len(v) > 0 {
		user.Username = v
	}

	user.Password = os.Getenv("TO_ADMIN_PASSWORD")
	if user.Password == "" {
		fmt.Println("no password supplied in ", passwordEnv)
		os.Exit(1)
	}

	// scan user input for invalid characters -- TODO: is this complete?
	strings.IndexFunc(user.Username,
		func(r rune) bool {
			return !unicode.IsPrint(r) && (unicode.IsPunct(r) || unicode.IsSpace(r))
		},
	)

	enc, err := auth.DerivePassword(user.Password)
	if err != nil {
		fmt.Printf("error encrypting password: %v\n", err)
		os.Exit(1)
	}

	user.EncryptedPassword = enc
	tmpl, err := template.New("insertuser").Parse(insertTemplate)
	if err != nil {
		fmt.Printf("error parsing template:: %v\n", err)
		os.Exit(1)
	}

	err = tmpl.Execute(os.Stdout, user)
	if err != nil {
		fmt.Printf("error executing template:: %v\n", err)
		os.Exit(1)
	}
}
