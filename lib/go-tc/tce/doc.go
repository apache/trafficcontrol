// Package tce (Traffic Control enums) contains enumerations and strongly typed names.
//
// These enums should be treated as enumerables, and MUST NOT be cast as anything else (integer, strings, etc). Enums MUST NOT be compared to strings or integers via casting. Enumerable data SHOULD be stored as the enumeration, not as a string or number. The *only* reason they are internally represented as strings, is to make them implicitly serialize to human-readable JSON. They should not be treated as strings. Casting or storing strings or numbers defeats the purpose of enum safety and conveniences.
//
// When storing enumumerable data in memory, it SHOULD be converted to and stored as an enum via the corresponding `FromString` function, checked whether the conversion failed and Invalid values handled, and valid data stored as the enum. This guarantees stored data is valid, and catches invalid input as soon as possible.
//
// When adding new enum types, enums should be internally stored as strings, so they implicitly serialize as human-readable JSON, unless the performance or memory of integers is necessary (it almost certainly isn't). Enums should always have the "invalid" value as the empty string (or 0), so default-initialized enums are invalid.
// Enums should always have a FromString() conversion function, to convert input data to enums. Conversion functions should usually be case-insensitive, and may ignore underscores or hyphens, depending on the use case.
//
package tce

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
