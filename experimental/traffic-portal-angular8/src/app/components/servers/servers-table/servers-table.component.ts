import { Component, OnInit } from '@angular/core';

import {first } from 'rxjs/operators';

import { ServerService } from '../../../services/api';
import { Server } from '../../../models/server';

@Component({
	selector: 'servers-table',
	templateUrl: './servers-table.component.html',
	styleUrls: ['./servers-table.component.scss']
})
export class ServersTableComponent implements OnInit {

	public servers: Array<Server>;

	columnDefs = [
		{headerName: 'ID', field: 'id' },
		{headerName: 'Hostname', field: 'hostName' },
		{headerName: 'Profile', field: 'profile'}
	];

	constructor(private readonly api: ServerService) {
		this.servers = [];
	}

	ngOnInit() {
		this.api.getServers().pipe(first()).subscribe(
			(r: Array<Server>) => {
				this.servers = r;
			}
		);
	}

}
