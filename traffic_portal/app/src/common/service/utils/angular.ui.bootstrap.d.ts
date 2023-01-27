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
 * @file
 *
 * Much of the content here has been copied from the
 * `@types/angular-ui-bootstrap` NPM package, which is provided by the Microsoft
 * Corporation under the MIT license.
 *
 * Those typings appear to be broken, or at the very least this author could not
 * figure out how to import them in a JSDoc type decoration. Therefore, the
 * parts used by this project were extracted, slightly modified, and pasted
 * here.
 *
 * For details and the original source, refer to:
 * https://www.npmjs.com/package/@types/angular-ui-bootstrap
 */

import type { IPromise, IScope, IAugmentedJQuery } from "angular";

interface IModalScope extends IScope {
	$dismiss(reason?: unknown): boolean;
	$close(result?: unknown): boolean;
}

export interface IModalSettings<T = undefined> {
	templateUrl?: string | (() => string) | undefined;
	template?: string | (() => string) | undefined;
	scope?: IScope | IModalScope | undefined;
	controller?: string | (() => void) | Array<string | (()=>void)> | undefined;
	controllerAs?: string | undefined;
	/**
	 * @default false
	 */
	bindToController?: boolean | undefined;
	resolve?: T; //{ [key: string]: string | (()=>void) | Array<string | (()=>void)> | object } | undefined;
	/**
	 * @default true
	 */
	animation?: boolean | undefined;
	/**
	 * @default true
	 */
	backdrop?: boolean | string | undefined;
	/**
	 * @default true
	 */
	keyboard?: boolean | undefined;
	backdropClass?: string | undefined;
	windowClass?: string | undefined;
	size?: string | undefined;
	windowTemplateUrl?: string | undefined;
	/**
	 * @default "model-open"
	 */
	openedClass?: string | undefined;
	windowTopClass?: string | undefined;
	/**
	 * @default "body"
	 */
	appendTo?: IAugmentedJQuery | undefined;
	component?: string | undefined;
	ariaDescribedBy?: string | undefined;
	ariaLabelledBy?: string | undefined;
}


export interface IModalService {
	getPromiseChain(): IPromise<unknown>;
	open<T, U>(options: IModalSettings<T>): IModalInstanceService<U>;
}

export interface IModalInstanceService<T = undefined> {
	close(result?: T): void;
	dismiss(reason?: T): void;
	result: IPromise<T>;
	opened: IPromise<unknown>;
	rendered: IPromise<unknown>;
	closed: IPromise<unknown>;
}
