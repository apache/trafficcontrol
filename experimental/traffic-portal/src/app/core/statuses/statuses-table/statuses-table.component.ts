import { Component, OnInit } from '@angular/core';
import { BehaviorSubject } from 'rxjs';
import { StatusesModel, StatusesService } from 'src/app/api/statuses.service';
import { ContextMenuItem } from 'src/app/shared/generic-table/generic-table.component';

@Component({
  selector: 'tp-statuses-table',
  templateUrl: './statuses-table.component.html',
  styleUrls: ['./statuses-table.component.scss']
})
export class StatusesTableComponent implements OnInit {

  statuses: any | null = null;
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
  public contextMenuItems: Array<ContextMenuItem<StatusesModel>> = [
    {
      href: (u: StatusesModel): string => `${u.id}`,
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
    private statusesService: StatusesService
  ) { }

  ngOnInit(): void {
    this.getStatuses();
  }

  /** Reloads the servers table data. */
  async getStatuses(): Promise<void> {
    this.statuses = await this.statusesService.getStatuses();
  }

  /**
 * Updates the "search" query parameter in the URL every time the search
 * text input changes.
 */
  public updateURL(): void {
    this.fuzzySubject.next(this.searchText);
  }
}
