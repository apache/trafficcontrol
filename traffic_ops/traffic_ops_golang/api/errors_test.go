package api

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
	"errors"
	"fmt"
	"net/http"
)

func ExampleErrors_String() {
	fmt.Println(NewErrors(http.StatusOK, nil, nil).String())

	// Output: 200 OK, SystemError='<nil>', UserError='<nil>'
}
func ExampleErrors_Error() {
	fmt.Println(NewErrors(http.StatusAccepted, errors.New("user error"), errors.New("system error")).Error())
	// Output: user error
}
func ExampleErrors_Code() {
	e := NewErrors(http.StatusOK, nil, nil)
	fmt.Println(e.Code())
	// Output: 200
}
func ExampleErrors_SetCode() {
	e := NewErrors(http.StatusOK, errors.New("something happened"), nil)
	e.SetCode(http.StatusBadRequest)
	fmt.Println(e.String())
	// Output: 400 Bad Request, SystemError='<nil>', UserError='something happened'
}
func ExampleErrors_SystemError() {
	e := NewErrors(http.StatusAccepted, nil, errors.New("testquest"))
	fmt.Println(e.SystemError())
	// Output: testquest
}
func ExampleErrors_UserError() {
	e := NewErrors(http.StatusAccepted, errors.New("testquest"), nil)
	fmt.Println(e.UserError())
	// Output: testquest
}
func ExampleErrors_SetSystemError() {
	e := NewErrors(http.StatusInternalServerError, nil, nil)
	e.SetSystemError(errors.New("testquest"))
	fmt.Println(e.String())
	// Output: 500 Internal Server Error, SystemError='testquest', UserError='<nil>'
}
func ExampleErrors_SetSystemErrorString() {
	e := NewErrors(http.StatusInternalServerError, nil, nil)
	e.SetSystemErrorString("testquest")
	fmt.Println(e.String())
	// Output: 500 Internal Server Error, SystemError='testquest', UserError='<nil>'
}
func ExampleErrors_SetSystemErrorf() {
	e := NewErrors(http.StatusInternalServerError, nil, nil)
	e.SetSystemErrorf("test: %w", errors.New("quest"))
	fmt.Println(e.String())
	// Output: 500 Internal Server Error, SystemError='test: quest', UserError='<nil>'
}
func ExampleErrors_SetUserError() {
	e := NewErrors(http.StatusInternalServerError, nil, nil)
	e.SetUserError(errors.New("testquest"))
	fmt.Println(e.String())
	// Output: 500 Internal Server Error, SystemError='<nil>', UserError='testquest'
}
func ExampleErrors_SetUserErrorString() {
	e := NewErrors(http.StatusInternalServerError, nil, nil)
	e.SetUserErrorString("testquest")
	fmt.Println(e.String())
	// Output: 500 Internal Server Error, SystemError='<nil>', UserError='testquest'
}
func ExampleErrors_SetUserErrorf() {
	e := NewErrors(http.StatusInternalServerError, nil, nil)
	e.SetUserErrorf("test: %w", errors.New("quest"))
	fmt.Println(e.String())
	// Output: 500 Internal Server Error, SystemError='<nil>', UserError='test: quest'
}
func ExampleErrors() {
	handler := func(fail bool) Errors {
		if fail {
			return NewSystemErrorString("failed")
		}
		return nil
	}

	var errs error = handler(true)
	fmt.Println(errs)

	errs = handler(false)
	fmt.Println(errs)

	// Output: failed
	// <nil>
}

func ExampleErrs_Code() {
	var e *Errs
	fmt.Println(e.Code())
	// Output: 200
}
func ExampleErrs_SystemError() {
	var e *Errs
	fmt.Println(e.SystemError())
	// Output: <nil>
}
func ExampleErrs_UserError() {
	var e *Errs
	fmt.Println(e.UserError())
	// Output: <nil>
}
func ExampleErrs_Error() {
	e := &Errs{}
	fmt.Println(e.Error())

	e.SetSystemError(errors.New("testquest"))
	fmt.Println(e.Error())

	// Output:
	// testquest
}

func ExamlpleNewErrors() {
	fmt.Println(
		NewErrors(
			http.StatusForbidden,
			errors.New("you don't have permission to do that in that Tenant"),
			errors.New("user tried to access a forbidden tenant"),
		),
	)

	// Output: 403 Forbidden, SystemError='user tried to access a forbidden tenant', UserError='you don't have permission to do that in that Tenant'
}
func ExampleNewSystemError() {
	fmt.Println(NewSystemError(errors.New("testquest")).String())
	// Output: 500 Internal Server Error, SystemError='testquest', UserError='<nil>'
}
func ExampleNewSystemErrorString() {
	fmt.Println(NewSystemErrorString("testquest").String())
	// Output: 500 Internal Server Error, SystemError='testquest', UserError='<nil>'
}
func ExampleNewSystemErrorf() {
	fmt.Println(NewSystemErrorf("test: %w", errors.New("quest")).String())
	// Output: 500 Internal Server Error, SystemError='test: quest', UserError='<nil>'
}
func ExampleNewUserError() {
	fmt.Println(NewUserError(errors.New("testquest")).String())
	// Output: 500 Internal Server Error, SystemError='<nil>', UserError='testquest'
}
func ExampleNewUserErrorString() {
	fmt.Println(NewUserErrorString("testquest").String())
	// Output: 500 Internal Server Error, SystemError='<nil>', UserError='testquest'
}
func ExampleNewUserErrorf() {
	fmt.Println(NewUserErrorf("test: %w", errors.New("quest")).String())
	// Output: 500 Internal Server Error, SystemError='<nil>', UserError='test: quest'
}
func ExampleNewResourceModifiedError() {
	fmt.Println(NewResourceModifiedError().String())
	// Output: 412 Precondition Failed, SystemError='<nil>', UserError='resource was modified since the time specified by the request headers'
}
