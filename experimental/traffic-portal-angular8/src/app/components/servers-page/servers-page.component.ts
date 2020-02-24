import { Component, OnInit } from '@angular/core';
import { FormControl } from '@angular/forms';
import { Router, ActivatedRoute, ParamMap } from '@angular/router';

import { first } from 'rxjs/operators';

import { ServerService } from '../../services/api';
import { Server } from '../../models';
import { orderBy, fuzzyScore } from '../../utils';

@Component({
	selector: 'servers-page',
	templateUrl: './servers-page.component.html',
	styleUrls: ['./servers-page.component.scss']
})
export class ServersPageComponent implements OnInit {

	fuzzControl: FormControl;
	servers: Server[];
	filteredServers: Server[];

	constructor (private readonly router: Router, private readonly route: ActivatedRoute, private readonly api: ServerService) { }

	ngOnInit (): void {
		const searchParam = this.route.snapshot.queryParamMap.get('search');
		this.fuzzControl = new FormControl(searchParam || "");
		this.api.getServers().pipe(first()).subscribe(
			(r: Server[]) => {
				this.servers = orderBy(r, 'hostName') as Server[];
				this.filteredServers = Array.from(this.servers);
				this.sort();
			}
		);
	}

	updateURL (e: Event) {
		e.preventDefault();
		this.sort();
		if (this.fuzzControl.value === '') {
			this.router.navigate([], {replaceUrl: true, queryParams: null});
		} else if (this.fuzzControl.value) {
			this.router.navigate([], {replaceUrl: true, queryParams: {search: this.fuzzControl.value}});
		}
	}

	tracker = (unused_item, s: Server) => {return s.id;};

	sort () {
		this.filteredServers = this.servers.map(
			x => [x, fuzzyScore(x.hostName, this.fuzzControl.value)]
		).filter(x=>x[1]!==Infinity).sort(
			(a, b) => {
				if (a[1] > b[1]) {
					return 1;
				}
				if (a[1] < b[1]) {
					return -1;
				}
				return 0;
			}
		).map(x=>x[0]) as Server[];
	}
}
