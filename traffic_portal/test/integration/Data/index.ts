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

interface GetRequest {
	route: string;
	queryKey: string;
	queryValue: string,
	replace: "route"
}

interface CleanupData {
	route: string;
	getRequest?: Array<GetRequest>;
}

interface Cleanup {
	action: string;
	route: string;
	method: "post" | "get" | "delete";
	data: Array<CleanupData>;
}

type SetupData<T = unknown> = Record<string | symbol, T> & {
	getRequest?: Array<GetRequest>;
}

export interface LoginData {
	description?: string;
	password: string;
	username: string;
	validationMessage?: string;
}

interface TestData extends Record<string | symbol, string> {
	description: string;
	validationMessage: string;
}

interface TestCases {
	logins: Array<LoginData>;
	add: Array<TestData>;
	remove: Array<TestData>;
	update: Array<TestData>;
}

export interface Test {
	cleanup: Array<Cleanup>;
	setup: Array<SetupData>;
	tests: Array<TestCases>;
}

export * from "./asns";
export * from "./cachegroup";
export * from "./cdn";
export * from "./coordinates";
export * from "./divisions";
export * from "./login";
export * from "./origins";
export * from "./parameters";
export * from "./physlocations";
