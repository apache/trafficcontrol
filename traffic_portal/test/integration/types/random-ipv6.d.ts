// Type definitions for random-ipv6
// Project: https://github.com/mock-end/random-ipv6

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
interface Options {
	compressed?: boolean;
	max?: number | string;
	min?: number | string;
	padded?: boolean;
}

declare module "random-ipv6" {
	function randomIPv6(schema?: Options): string;
	function randomIPv6(schema: string, options?: Options): string;
	export = randomIPv6;
}
