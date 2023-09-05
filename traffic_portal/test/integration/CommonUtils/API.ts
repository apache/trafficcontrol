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

// API Utility
import { Agent } from "https";

import axios from 'axios';
import type {AxiosResponse, AxiosError} from "axios";
import randomIpv6 from "random-ipv6";

import { hasProperty } from "./utils";
import { randomize } from '../config';
import { AlertLevel, isAlert, logAlert, TestingConfig } from "../config.model";

interface GetRequest {
    queryKey: string;
    queryValue: string | number | boolean;
    replace: string | number;
    route: string;
}

export interface IDData extends Record<string | number, unknown> {
    getRequest?: Array<GetRequest>;
    route?: string;
}

interface APIDataData extends Record<PropertyKey, unknown>, IDData {
    id?: unknown;
}

export interface APIData {
    action: string;
    data: Array<APIDataData>;
    method: string;
    route: string;
}

/**
 * Checks if an object is an AxiosError, usually useful in `try`/`catch` blocks
 * around axios calls.
 *
 * @param e The object to check.
 * @returns Whether or not `e` is an AxiosError.
 */
function isAxiosError(e: unknown): e is AxiosError {
    if (typeof(e) !== "object" || e === null) {
        return false;
    }
    if (!hasProperty(e, "isAxiosError", "boolean")) {
        return false;
    }
    return e.isAxiosError;
}

export class API {

    private cookie = "";

    /**
     * This controls the alert levels that get logged - levels not in this set
     * are not logged
     */
    private readonly alertLevels = new Set<AlertLevel>(["warning", "error", "info"]);
    /**
     * Stores login information for the admin-level user.
     */
    private readonly loginInfo: {
        password: string;
        username: string;
    };
    /**
     * The URL base used for the Traffic Ops API.
     *
     * Trailing `/` is guaranteed.
     *
     * @example
     * "https://localhost:6443/api/4.0/"
     */
    private readonly apiURL: string;

    /**
     * @param cfg The testing configuration.
     */
    constructor(cfg: TestingConfig) {
        axios.defaults.headers.common['Accept'] = 'application/json'
        axios.defaults.headers.common['Authorization'] = 'No-Auth'
        axios.defaults.headers.common['Content-Type'] = 'application/json'
        axios.defaults.httpsAgent = new Agent({ rejectUnauthorized: false })
        if (cfg.alertLevels) {
            this.alertLevels = new Set(cfg.alertLevels);
        }
        this.loginInfo = cfg.login;
        this.apiURL = cfg.apiUrl.endsWith("/") ? cfg.apiUrl : `${cfg.apiUrl}/`;
    }

    /**
     * Logs the API client into Traffic Ops.
     *
     * @returns The API response from logging in.
     * @throws {Error} when login fails, or when Traffic Ops doesn't return a cookie.
     */
    public async Login(): Promise<AxiosResponse<unknown>> {
        const data = {
            p: this.loginInfo.password,
            u: this.loginInfo.username,
        }
        const response = await this.getResponse("post", "/user/login", data);
        const h = response.headers as object;
        if (!hasProperty(h, "set-cookie", "Array") || h["set-cookie"].length < 1) {
            throw new Error("Traffic Ops response did not set a cookie");
        }
        const cookie = await h["set-cookie"][0];
        if (typeof(cookie) !== "string") {
            throw new Error(`non-string cookie: ${cookie}`);
        }
        this.cookie = cookie;
        return response
    }

    /**
     * Retrieves a response from the API.
     *
     * Alerts will be logged if they are found - even if an error occurs and is
     * thrown.
     *
     * @param method The request method to use.
     * @param path The path to request, relative to the configured TO API URL.
     * @returns The server's response.
     * @throws {unknown} when the request fails for any reason. If an error
     * response was returned from the API, it was logged, so there's no need to
     * dig into the properties of these errors, really.
     */
    private async getResponse(method: "get" | "delete", path: string): Promise<AxiosResponse>;
    /**
     * Retrieves a response from the API.
     *
     * Alerts will be logged if they are found - even if an error occurs and is
     * thrown.
     *
     * @param method The request method to use.
     * @param path The path to request, relative to the configured TO API URL.
     * @param data Data to send in the body of the POST request.
     * @returns The server's response.
     * @throws {unknown} when the request fails for any reason. If an error
     * response was returned from the API, it was logged, so there's no need to
     * dig into the properties of these errors, really.
     */
    private async getResponse(method: "post", path: string, data: unknown): Promise<AxiosResponse>;
    private async getResponse(method: "post" | "get" | "delete", path: string, data?: unknown): Promise<AxiosResponse> {
        if (method === "post" && data === undefined) {
            throw new TypeError("request body must be given for POST requests");
        }

        const url = `${this.apiURL}${path.replace(/^\/+/g, "")}`;
        const conf = {
            method,
            url,
            headers: { Cookie: this.cookie },
            data
        }

        let throwable;
        let resp: AxiosResponse<unknown>;
        try {
            resp = await axios(conf);
        } catch(e) {
            if (!isAxiosError(e) || !e.response) {
                console.debug("non-axios error or axios error with no response thrown");
                throw e;
            }
            resp = e.response;
            throwable = e;
        }
        if (typeof(resp.data) === "object" && resp.data !== null && hasProperty(resp.data, "alerts", "Array")) {
            for (const a of resp.data.alerts) {
                if (isAlert(a) && this.alertLevels.has(a.level)) {
                    logAlert(a, `${method.toUpperCase()} ${url} (${resp.status} ${resp.statusText}):`);
                }
            }
        }
        if (throwable) {
            throw throwable;
        }
        return resp;
    }

    public async SendRequest<T extends IDData>(route: string, method: string, data: T): Promise<void> {
        let response
        this.Randomize(data)

        if(data.hasOwnProperty('getRequest')){
            let response;
            try {
                response = await this.GetId(data);
            } catch (e) {
                let msg = e instanceof Error ? e.message : String(e);
                if (response) {
                    msg = `response status: ${response.statusText}, response data: ${response.data} - ${msg}`;
                }
                throw new Error(`Failed to get id: ${msg}`);
            }
        }

        switch (method) {
            case "post":
                response = await this.getResponse("post", route, data);
                break;
            case "get":
                response = await this.getResponse("get", route);
                break;
            case "delete":
                if (!data.route) {
                    throw new Error("DELETE requests must include a 'route' data property")
                }
                if ((data.route).includes('?name')){
                    data.route = data.route + randomize
                }
                if ((data.route).includes('?id')){
                    if (!hasProperty(data, "id")) {
                        throw new Error("route specified an 'id' query parameter, but data had no 'id' property");
                    }
                    data.route = data.route + data.id;
                }
                if((data.route).includes('/service_categories/')){
                    data.route = data.route + randomize
                }
                response = await this.getResponse("delete", data.route);
                break;
            default:
                throw new Error(`unrecognized request method: '${method}'`);
        }
        if (response.status == 200 || response.status == 201) {
            return;
        } else {
            console.log("Reponse Data: " , response.data);
            console.log("Response: " , response);
            throw new Error(`request failed: response status: '${response.statusText}' response data: '${response.data}'`);
        }
    }

    public async GetId(data: IDData): Promise<null | AxiosResponse<unknown>> {
        if (!data.getRequest) {
            return null;
        }
        for (const request of data.getRequest) {
            let query = `?${encodeURIComponent(request.queryKey)}=`;
            if (request.queryValue === 'admin' || request.queryValue === 'operations' || request.queryValue === 'read-only'){
                query += encodeURIComponent(request.queryValue);
            }else{
                query += encodeURIComponent(request.queryValue+randomize);
            }
            const response = await this.getResponse("get", request.route + query)
            if (response.status == 200) {
                if(request.hasOwnProperty('isArray')){
                    data[request.replace] = [await response.data.response[0].id];
                } else if (request.replace === "route") {
                    data.route = data.route + response.data.response[0].id;
                } else {
                    data[request.replace] = await response.data.response[0].id;
                }
            } else {
                // todo: should this be getting cut short like this?
                return response
            }
        }
        return null
    }

    public Randomize(data: object): void {
       if (hasProperty(data, "fullName")) {
           if (hasProperty(data, "email")) {
               data.email = data.fullName + randomize + data.email;
           }
           data.fullName = data.fullName + randomize;
       }
        if (hasProperty(data, "hostName")) {
            data.hostName = data.hostName + randomize;
        }
        if (hasProperty(data, "ipAddress")) {
            const rand = () => Math.floor(Math.random()*255)+1;
            data.ipAddress = `${rand()}.${rand()}.${rand()}.${rand()}`;
        }
        if(hasProperty(data, 'name') && !(hasProperty(data, "noRandomize") && data.noRandomize === true)) {
            data.name = data.name + randomize;
        }
        if(hasProperty(data, 'requiredCapability')) {
            data.requiredCapability = data.requiredCapability + randomize;
        }
        if(hasProperty(data, 'requiredCapabilities', 'Array')) {
            data.requiredCapabilities.forEach((_cap, i) => {
                data.requiredCapabilities[i] += randomize;
            })
        }
        if(hasProperty(data, 'serverCapability')) {
            data.serverCapability = data.serverCapability + randomize;
        }
        if(hasProperty(data, 'username')) {
            data.username = data.username + randomize;
        }
        if(hasProperty(data, 'xmlId')) {
            data.xmlId = data.xmlId + randomize;
        }
        if(hasProperty(data, 'shortName')) {
            data.shortName = data.shortName + randomize;
        }
        if(hasProperty(data, 'divisionName')) {
            data.divisionName = data.divisionName + randomize;
        }
        if(hasProperty(data, 'domainName')) {
            data.domainName = data.domainName + randomize;
        }
        if(hasProperty(data, 'nodes', "Array")){
            data.nodes.map(i => {
                if (typeof(i) === "object" && i !== null && hasProperty(i, "cachegroup")) {
                    i.cachegroup = i.cachegroup + randomize;
                }
            });
        }
        if(hasProperty(data, 'interfaces', "Array")){
            const ipv6 = randomIpv6();
            for (const inf of data.interfaces) {
                if (typeof(inf) === "object" && inf !== null && hasProperty(inf, "ipAddresses", "Array")) {
                    for (const ip of inf.ipAddresses) {
                        (ip as Record<"address", string>).address = ipv6.toString();
                    }
                }
            }
        }
        if(hasProperty(data, 'profiles', "Array")){
            for (const index in data.profiles) {
                data.profiles[index] = data.profiles[index]+randomize
            }
        }
    }

    public async UseAPI(data: Array<APIData>): Promise<void> {
        const response = await this.Login();
        if (response.status === 200) {
            for(let i = 0; i < data.length; i++){
                for(let j = 0; j < data[i].data.length; j++){
                    const route = data[i].data[j].route ?? data[i].route;
                    try {
                        await this.SendRequest(route, data[i].method, data[i].data[j]);
                    } catch (output) {
                        if (output instanceof Error) {
                            output = output.message;
                        }
                        console.debug(`${data[i].method} ${route}`);
                        console.debug("DATA:", data[i].data[j]);
                        throw new Error(`UseAPI failed on Action ${data[i].action} with index ${i}, and Data index ${j}: ${output}`);
                    }
                }
            }
        } else if (response.status == undefined) {
            throw new Error(`Error requesting ${this.apiURL}: ${response}`);
        } else {
            throw new Error(`Login failed: Response Status: '${response.statusText}'' Response Data: '${response.data}'`);
        }
    }
}
