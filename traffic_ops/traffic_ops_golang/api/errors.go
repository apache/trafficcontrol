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

// Errs is the concrete implementation of Errors, which is used so that we can
// deal with Errors everywhere without having to handle a bunch of
// (de/)referencing everywhere.
type Errs struct {
	// code is the HTTP response code. If no error has occurred, this *should*
	// be http.StatusOK.
	code int
	// systemError is an error that occurred internally; not safe for exposure
	// to the user.
	systemError error
	// userError is an error that should be shown to the user to explain why
	// their request failed.
	userError error
}

// Code returns the HTTP response code. If no error has occurred, this
// *should* return http.StatusOK.
func (e *Errs) Code() int {
	if e == nil {
		return http.StatusOK
	}
	return e.code
}

// SetCode sets the HTTP response code returned by `Code`.
func (e *Errs) SetCode(code int) {
	e.code = code
}

// SystemError returns an error that occurred internally; not safe for
// exposure to the user.
func (e *Errs) SystemError() error {
	if e == nil {
		return nil
	}
	return e.systemError
}

// UserError returns an error that should be shown to the user to explain
// why their request failed. This **must** return `nil` whenever `Code`
// returns a value above 499.
func (e *Errs) UserError() error {
	if e == nil {
		return nil
	}
	return e.userError
}

// SetSystemErrorString sets the Errs' system-level error to a new error
// containing the passed message.
func (e *Errs) SetSystemErrorString(err string) {
	e.systemError = errors.New(err)
}

// SetSystemError sets the Errs' system-level error to the passed error.
func (e *Errs) SetSystemError(err error) {
	e.systemError = err
}

// SetSystemErrorf sets the Errs' system-level error to a new error containing
// the passed formatted message, in the fmt package's syntax (allows wrapping
// via %w).
func (e *Errs) SetSystemErrorf(format string, args ...any) {
	e.systemError = fmt.Errorf(format, args...)
}

// SetUserError sets the Errs' user-level error to the passed error.
func (e *Errs) SetUserError(err error) {
	e.userError = err
}

// SetUserErrorString sets the Errs' user-level error to a new error containing
// the passed message.
func (e *Errs) SetUserErrorString(err string) {
	e.userError = errors.New(err)
}

// SetUserErrorf sets the Errs' user-level error to a new error containing the
// passed formatted message, in the fmt package's syntax (allows wrapping via
// %w).
func (e *Errs) SetUserErrorf(format string, args ...any) {
	e.userError = fmt.Errorf(format, args...)
}

// String implements the fmt.Stringer interface.
func (e *Errs) String() string {
	return fmt.Sprintf("%d %s, SystemError='%v', UserError='%v'", e.code, http.StatusText(e.code), e.systemError, e.userError)
}

// Error implements the error interface. This will return the user error, if
// there is one, otherwise it will return the system error. In the case that no
// actual error has occurred, it will neturn an empty string. This panics if the
// Errs it is called on is nil (just like a normal error would).
func (e *Errs) Error() string {
	if e.userError != nil {
		return e.userError.Error()
	}
	if e.systemError != nil {
		return e.systemError.Error()
	}
	return ""
}

// Errors represents a set of zero to two errors to be handled by the API.
type Errors interface {
	// Code returns the HTTP response code. If no error has occurred, this
	// *should* return http.StatusOK.
	Code() int
	// SetCode sets the HTTP response code returned by `Code`. In general, you
	// probably won't need to use this.
	SetCode(int)
	// SystemError returns an error that occurred internally; not safe for
	// exposure to the user.
	SystemError() error
	// SetSystemErrorString sets the Errors' system-level error to the passed
	// error. In general, you probably won't need to use this.
	SetSystemError(err error)
	// SetSystemErrorf sets the Errors' system-level error to a new error
	// containing the passed formatted message, in the fmt package's syntax
	// (allows wrapping via %w). In general, you probably won't need to use
	// this.
	SetSystemErrorf(format string, args ...any)
	// SetSystemErrorString sets the Errors' system-level error to a new error
	// containing the passed message. In general, you probably won't need to use
	// this.
	SetSystemErrorString(err string)
	// UserError returns an error that should be shown to the user to explain
	// why their request failed. This **must** return `nil` whenever `Code`
	// returns a value above 499.
	UserError() error
	// SetUserErrorString sets the Errors' user-level error to the passed error.
	// In general, you probably won't need to use this.
	SetUserError(err error)
	// SetUserErrorf sets the Errors' user-level error to a new error containing
	// the passed formatted message, in the fmt package's syntax (allows
	// wrapping via %w). In general, you probably won't need to use this.
	SetUserErrorf(format string, args ...any)
	// SetUserErrorString sets the Errors' user-level error to a new error
	// containing the passed message. In general, you probably won't need to use
	// this.
	SetUserErrorString(err string)

	// String implements the fmt.Stringer interface.
	String() string
	// Error implements the error interface. This will return the user error, if
	// there is one, otherwise it will return the system error. In the case that
	// no actual error has occurred, it will neturn an empty string. This panics
	// if the Errors it is called on is nil (just like a normal error would).
	Error() string
}

// NewErrors directly constructs a new Errors using the passed information.
func NewErrors(code int, userError, systemError error) Errors {
	return &Errs{
		code:        code,
		systemError: systemError,
		userError:   userError,
	}
}

// NewSystemError creates an Errors that only contains the given system error,
// and has the appropriate response code.
func NewSystemError(err error) Errors {
	return &Errs{
		code:        http.StatusInternalServerError,
		systemError: err,
		userError:   nil,
	}
}

// NewSystemErrorString creates an Errors that only contains a system error,
// having the appropriate response code and containing the given message.
func NewSystemErrorString(err string) Errors {
	return &Errs{
		code:        http.StatusInternalServerError,
		systemError: errors.New(err),
		userError:   nil,
	}
}

// NewSystemErrorf creates an Errors that only contains a system error, having
// the appropriate response code, and containing a message from the given
// format arguments (using the fmt package's syntax, with %w allowed for
// wrapping).
func NewSystemErrorf(format string, args ...any) Errors {
	return &Errs{
		code:        http.StatusInternalServerError,
		systemError: fmt.Errorf(format, args...),
		userError:   nil,
	}
}

// NewUserError creates an Errors that only contains the given user error,
// and has the appropriate response code.
func NewUserError(err error) Errors {
	return &Errs{
		code:        http.StatusInternalServerError,
		systemError: nil,
		userError:   err,
	}
}

// NewUserErrorString creates an Errors that only contains a user error,
// having the appropriate response code and containing the given message.
func NewUserErrorString(err string) Errors {
	return &Errs{
		code:        http.StatusInternalServerError,
		systemError: nil,
		userError:   errors.New(err),
	}
}

// NewUserErrorf creates an Errors that only contains a user error, having
// the appropriate response code, and containing a message from the given
// format arguments (using the fmt package's syntax, with %w allowed for
// wrapping).
func NewUserErrorf(format string, args ...any) Errors {
	return &Errs{
		code:        http.StatusInternalServerError,
		systemError: nil,
		userError:   fmt.Errorf(format, args...),
	}
}

// NewResourceModifiedError creates an Errors that only contains an HTTP
// Precondition Failed status code and associated error message.
func NewResourceModifiedError() Errors {
	return &Errs{
		code:        http.StatusPreconditionFailed,
		systemError: nil,
		userError:   ResourceModifiedError,
	}
}
