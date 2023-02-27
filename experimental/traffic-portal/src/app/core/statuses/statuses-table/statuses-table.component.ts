/*
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
import { Component, OnInit } from '@angular/core';
import { BehaviorSubject } from 'rxjs';
import { ServerService } from 'src/app/api';
import { ContextMenuItem } from 'src/app/shared/generic-table/generic-table.component';
import { ResponseStatus } from 'trafficops-types';

@Component({
  selector: 'tp-statuses-table',
  templateUrl: './statuses-table.component.html',
  styleUrls: ['./statuses-table.component.scss']
})
export class StatusesTableComponent implements OnInit {

  statuses: ResponseStatus[] | null = null;
  columnDefs = [
    {
      field: "name",
      headerName: "Name",
      hide: false
    },
    {
      field: "description",
      headerName: "Description",
      hide: false
    }];

  /** The current search text. */
  public searchText = "";

  /** Definitions for the context menu items (which act on user data). */
  public contextMenuItems: Array<ContextMenuItem<ResponseStatus>> = [
    {
      href: (u: ResponseStatus): string => `${u.id}`,
      name: "View Status Details"
    },
    {
      href: (): string => `new`,
      name: "Create New Status"
    }
  ];

  /** Emits changes to the fuzzy search text. */
  public fuzzySubject = new BehaviorSubject("");
  constructor(
    private serverService: ServerService,
  ) { }

  ngOnInit(): void {
    this.getStatuses();
  }

  /** Reloads the servers table data. */
  async getStatuses(): Promise<void> {
    this.statuses = await this.serverService.getStatuses();
  }

  /**
 * Updates the "search" query parameter in the URL every time the search
 * text input changes.
 */
  public updateURL(): void {
    this.fuzzySubject.next(this.searchText);
  }
}
