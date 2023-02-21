import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { APIService } from './base-api.service';

export interface StatusesModel {
	description?: string;
	id?: number;
	lastUpdated?: Date;
	name?: string;
}

@Injectable({
	providedIn: 'root'
})
export class StatusesService extends APIService {

	/**
	 * Injects the Angular HTTP client service into the parent constructor.
	 * @param http The Angular HTTP client service.
	 */
	constructor(http: HttpClient) {
		super(http);
	}

	public async getStatuses(idOrName: number | string): Promise<StatusesModel>;
	public async getStatuses(): Promise<Array<StatusesModel>>;
	/**
	 * @param id Specify either the integral, unique identifier (number.
	 * @returns The requested status(s).
	 */
	public async getStatuses(id?: number | string): Promise<Array<StatusesModel> | StatusesModel> {
		const path = "statuses";
		if (id !== undefined) {
			let statuses;
			statuses = await this.get<[StatusesModel]>(path, undefined, { id: String(id) }).toPromise();
			if (statuses.length < 1) {
				throw new Error(`no such statuses '${id}'`);
			}
			return statuses[0];
		}
		return this.get<Array<StatusesModel>>(path).toPromise();
	}

	/**
	 * Creating new Status.
	 * @param data containes name and description for the status.
	 * @returns The 'response' property of the TO status response. See TO API docs.
	 */
	public async createStatus(data: StatusesModel) {
		const path = "statuses";
		return this.post<StatusesModel>(path, data).toPromise();
	}

	/**
	 * Updates status.
	 * @param data containes name and description for the status., unique identifier thereof.
	 * @param id The Status ID
	 */
	public async updateStatus(data: StatusesModel, id: number): Promise<StatusesModel | undefined> {
		const path = `statuses/${id}`;
		return this.put(path, data).toPromise();
	}

	/**
	 * Deletes an existing Status.
	 * @param id The Status ID
	 */
	public async deleteStatus(id: number): Promise<void> {
		return this.delete(`statuses/${id}`).toPromise();
	}
}
