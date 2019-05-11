package version

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

// These values are irrelevant during an rpm build as they will be forced to compile time values via ldflags in build_rpm.sh
// If you need to update the defaults, do so in VERSION

// VerMajor should be be incremented whenever backward breaking changes are made
var VerMajor = "0"

// VerMinor should be updated whenever new features and bugfix rollups are performed, a typical minor release
var VerMinor = "1"

// VerPatch should be incremented if patches are backported to old releases
var VerPatch = "0"

// VerDesc should be an arbitrary string prefix to the git commit to distinguish what release pipeline stage this is from
var VerDesc = "dev"

// VerCommit should represent the git hash that was compiled with the binary
var VerCommit = "0"

// VerFull should be the full version string
var VerFull = VerMajor + "." + VerMinor + "." + VerPatch + "_" + VerDesc + "_" + VerCommit
