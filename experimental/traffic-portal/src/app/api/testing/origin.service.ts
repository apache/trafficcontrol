/*
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
import { Injectable } from "@angular/core";
import type { RequestOrigin, RequestOriginResponse } from "trafficops-types";

import {
	CacheGroupService,
	DeliveryServiceService,
	ProfileService,
	UserService,
} from "..";

/**
 * OriginService exposes API functionality relating to Origins.
 */
@Injectable()
export class OriginService {
	private readonly origins: Array<RequestOriginResponse> = [
		{
			cachegroup: null,
			cachegroupId: null,
			coordinate: null,
			coordinateId: null,
			deliveryService: "test",
			deliveryServiceId: 1,
			fqdn: "origin.infra.ciab.test",
			id: 1,
			ip6Address: null,
			ipAddress: null,
			isPrimary: false,
			lastUpdated: new Date(),
			name: "test",
			port: null,
			profile: null,
			profileId: null,
			protocol: "http",
			tenant: "root",
			tenantId: 1,
		},
	];

	constructor(
		private readonly userService: UserService,
		private readonly cacheGroupService: CacheGroupService,
		private readonly profileService: ProfileService,
		private readonly dsService: DeliveryServiceService
	) {}

	/**
	 * Gets a specific Origin.
	 *
	 * @param nameOrID Either the integral, unique identifier (number) or name
	 * (string) of the Origin to be returned.
	 * @returns The requested Origin.
	 */
	public async getOrigins(nameOrID: number | string): Promise<RequestOriginResponse>;
	/**
	 * Gets all Origins.
	 *
	 * @returns All stored Origins.
	 */
	public async getOrigins(): Promise<Array<RequestOriginResponse>>;
	/**
	 * Gets one or all Origins.
	 *
	 * @param nameOrID Optionally the integral, unique identifier (number) or
	 * name (string) of a single Origin to be returned.
	 * @returns The requested Origin(s).
	 */
	public async getOrigins(
		nameOrID?: string | number
	): Promise<Array<RequestOriginResponse> | RequestOriginResponse> {
		if (nameOrID) {
			let origin;
			switch (typeof nameOrID) {
				case "string":
					origin = this.origins.find((d) => d.name === nameOrID);
					break;
				case "number":
					origin = this.origins.find((d) => d.id === nameOrID);
			}
			if (!origin) {
				throw new Error(`no such Origin: ${nameOrID}`);
			}
			return origin;
		}
		return this.origins;
	}

	/**
	 * Replaces the current definition of a Origin with the one given.
	 *
	 * @param origin The new Origin.
	 * @returns The updated Origin.
	 */
	public async updateOrigin(
		origin: RequestOriginResponse
	): Promise<RequestOriginResponse> {
		const id = this.origins.findIndex((d) => d.id === origin.id);
		if (id === -1) {
			throw new Error(`no such Origin: ${origin.id}`);
		}
		this.origins[id] = origin;
		return origin;
	}

	/**
	 * Creates a new Origin.
	 *
	 * @param origin The Origin to create.
	 * @returns The created Origin.
	 */
	public async createOrigin(
		origin: RequestOrigin
	): Promise<RequestOriginResponse> {
		const tenant = await this.userService.getTenants(origin.tenantID);
		const ds = await this.dsService.getDeliveryServices(
			origin.deliveryServiceId
		);
		let profile = null;
		if (!!origin?.profileId) {
			profile = await this.profileService.getProfiles(origin.profileId);
		}
		let coordinate = null;
		if (!!origin?.coordinateId) {
			coordinate = await this.cacheGroupService.getCoordinates(
				origin.coordinateId
			);
		}
		let cacheGroup = null;
		if (!!origin?.cachegroupId) {
			cacheGroup = await this.cacheGroupService.getCacheGroups(
				origin.cachegroupId
			);
		}

		const created = {
			cachegroup: cacheGroup?.name ?? null,
			cachegroupId: cacheGroup?.id ?? null,
			coordinate: coordinate?.name ?? null,
			coordinateId: coordinate?.id ?? null,
			deliveryService: ds.displayName ?? null,
			deliveryServiceId: ds.id ?? null,
			fqdn: "",
			id: 1,
			ip6Address: null,
			ipAddress: null,
			isPrimary: null,
			lastUpdated: new Date(),
			name: "",
			port: null,
			profile: profile?.name ?? null,
			profileId: profile?.id ?? null,
			protocol: "https" as never,
			tenant: tenant.name ?? null,
			tenantId: tenant.id ?? null,
		};
		this.origins.push(created);
		return created;
	}

	/**
	 * Deletes an existing Origin.
	 *
	 * @param origin The Origin to be deleted (or its ID)
	 */
	public async deleteOrigin(
		origin: RequestOriginResponse | number
	): Promise<void> {
		const id = typeof origin === "number" ? origin : origin.id;
		const index = this.origins.findIndex((d) => d.id === id);
		if (index === -1) {
			throw new Error(`no such Origin: ${id}`);
		}
		this.origins.splice(index, 1);
	}
}
