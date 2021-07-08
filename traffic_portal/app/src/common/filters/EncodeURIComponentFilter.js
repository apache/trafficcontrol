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

/**
 * This is a factory function for an AngularJS filter meant to allow easy
 * sanitization of interpolated input into URL query strings.
 *
 * @example
 *
 * <a ng-href="https://example.test/something?foo={{someString | encodeURIComponent}}">
 *   A safe link
 * </a>
 *
 * @returns {(input: string) => string} A filter that sanitizes input text for safe insertion into URL query strings.
 */
var EncodeURIComponentFilter = function() {
    return function(input) {
        return encodeURIComponent(input);
    };
};

module.exports = EncodeURIComponentFilter;
