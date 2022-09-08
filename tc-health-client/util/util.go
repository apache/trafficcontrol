package util

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
	"io"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"unsafe"
)

type ConfigFile struct {
	Filename       string
	LastModifyTime int64
}

// get the file modification times for a configuration file.
func GetFileModificationTime(fn string) (int64, error) {
	f, err := os.Open(fn)
	if err != nil {
		return 0, errors.New("opening " + fn + ": " + err.Error())
	}
	defer f.Close()

	finfo, err := f.Stat()
	if err != nil {
		return 0, errors.New("unable to get file status for " + fn + ": " + err.Error())
	}
	return finfo.ModTime().UnixNano(), nil
}

// NopWriter is a no-op io.WriteCloser.
// It always returns the length of the given bytes as written.
// Writes to this go nowhere. Close may be called, and is also a no-op.
type NopWriter struct{}

// Write implements io.Writer.
// It always returns the length of the given bytes as written.
func (nw NopWriter) Write(p []byte) (n int, err error) { return len(p), nil }
func (nw NopWriter) Close() error                      { return nil }

// WriteNoCloser wraps an io.Writer and returns an io.WriteCloser for which Close will be a no-op.
//
// This allows for passing a WriteCloser to something that ordinarily takes ownership and closes,
// Without that Close() closing the real Writer. For example, to prevent closing os.Stdout.
//
// This can also be used to turn an io.Writer into an io.WriteCloser.
type NoCloser struct {
	wr io.Writer
}

func MakeNoCloser(wr io.Writer) io.WriteCloser        { return NoCloser{wr: wr} }
func (nc NoCloser) Write(p []byte) (n int, err error) { return nc.wr.Write(p) }
func (nw NoCloser) Close() error                      { return nil }

// AtomicPtr provides atomic getting and setting of a pointer.
// It may be default-constructed or via NewAtomicPtr.
// It must not be copied after first use.
//
// No further synchronization is provided, besides the atomic pointer exchange.
// If value pointed to is modified concurrently, users must perform their own synchronization.
type AtomicPtr[T any] struct {
	up *unsafe.Pointer
}

// NewAtomicPtr creates a new AtomicPtr[T] initialized with ptr.
func NewAtomicPtr[T any](ptr *T) *AtomicPtr[T] {
	up := (unsafe.Pointer)(ptr)
	return &AtomicPtr[T]{up: &up}
}

// Get returns the pointer.
// This is safe for multiple concurrent goroutines.
//
// No further synchronization is provided, besides the atomic pointer exchange.
// If the returned object is modified concurrently, users must perform their own synchronization.
func (ag *AtomicPtr[T]) Get() *T {
	return (*T)(atomic.LoadPointer(ag.up))
}

// Set sets the AtomicPtr's pointer.
// This is safe for multiple concurrent goroutines.
//
// No further synchronization is provided, besides the atomic pointer exchange.
// If the set object is modified concurrently, users must perform their own synchronization
func (ag *AtomicPtr[T]) Set(ptr *T) {
	up := (unsafe.Pointer)(ptr)
	atomic.StorePointer(ag.up, up)
}

// SyncMap is a strongly-typed sync.Map.
type SyncMap[KT any, VT any] struct {
	m sync.Map
}

func (sm *SyncMap[KT, VT]) Load(key KT) (VT, bool) {
	iVal, ok := sm.m.Load(key)
	if !ok {
		var vt VT
		return vt, false
	}
	return iVal.(VT), true
}

func (sm *SyncMap[KT, VT]) Store(key KT, val VT) {
	sm.m.Store(key, val)
}

func (sm *SyncMap[KT, VT]) Range(f func(key KT, value VT) bool) {
	sm.m.Range(func(iKey, iVal interface{}) bool {
		key := iKey.(KT)
		val := iVal.(VT)
		return f(key, val)
	})
}

func (sm *SyncMap[KT, VT]) LoadOrStore(key KT, val VT) (actual VT, loaded bool) {
	iVal, loaded := sm.m.LoadOrStore(key, val)
	return iVal.(VT), loaded
}

// HostNameToShort takes a hostname which may be a Fully Qualified Domain Name or a short hostname,
// and returns the short hostname.
// If host is already a short hostname, it is returned unmodified.
func HostNameToShort(host string) string {
	dotPos := strings.Index(host, ".")
	if dotPos < 0 {
		return host
	}
	return host[:dotPos]
}
