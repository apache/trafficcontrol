import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import { MatDialog } from '@angular/material/dialog';
import { ActivatedRoute, Router } from '@angular/router';
import { StatusesModel, StatusesService } from 'src/app/api/statuses.service';
import { DecisionDialogComponent, DecisionDialogData } from 'src/app/shared/dialogs/decision-dialog/decision-dialog.component';

@Component({
  selector: 'tp-status-details',
  templateUrl: './status-details.component.html',
  styleUrls: ['./status-details.component.scss']
})
export class StatusDetailsComponent implements OnInit {

  id: string | null = null;
  statusDetails: StatusesModel | null = null;
  statusDetailsForm!: FormGroup;
  loading = false;
  submitting = false;
  submitted = false;

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private fb: FormBuilder,
    private readonly dialog: MatDialog,
    private statusesService: StatusesService) { }

  ngOnInit(): void {
    // Form is built here
    this.statusDetailsForm = this.fb.group({
      name: ['', Validators.required],
      description: ['', Validators.required],
    });

    // Getting id from the route
    this.id = this.route.snapshot.paramMap.get('id');

    // we check whether params is a number if not we shall assume user wants to add a new status.
    if (!this.isNew) {
      this.loading = true;
      this.statusDetailsForm.addControl('id', new FormControl(''));
      this.statusDetailsForm.addControl('lastUpdated', new FormControl(''));
      this.getStatusDetails();
    }
  }

  /*
   * Reloads the servers table data. 
   * @param id is the id passed in route for this page if this is a edit view.
  */
  async getStatusDetails(): Promise<void> {
    const id = this.id as string
    this.statusDetails = await this.statusesService.getStatuses(id);
    const data = {
      name: this.statusDetails.name,
      description: this.statusDetails.description,
      lastUpdated: new Date(),
      id: this.statusDetails.id
    } as any;
    this.statusDetailsForm.patchValue(data);
    this.loading = false;
  }

  // On submitting the form we check for whether we are performing Create or Edit
  onSubmit() {
    if (this.isNew) {
      this.createStatus()

    } else {
      this.updateStatus();
    }
  }

  // For Creating a new status
  createStatus() {
    this.statusesService.createStatus(this.statusDetailsForm.value).then((res: any) => {
      if (res) {
        this.id = res?.id
        this.router.navigate([`/core/statuses/${this.id}`]);
      }
    })
  }

  // For updating the Status
  updateStatus() {
    this.statusesService.updateStatus(this.statusDetailsForm.value, Number(this.id));
  }

  // Deleteting status
  async deleteStatus() {
    const ref = this.dialog.open<DecisionDialogComponent, DecisionDialogData, boolean>(DecisionDialogComponent, {
      data: {
        message: `This action CANNOT be undone. This will permanently delete '${this.statusDetails?.name}'.`,
        title: `Delete Status: ${this.statusDetails?.name}`
      }
    });

    if (await ref.afterClosed().toPromise()) {
      const id = Number(this.id);
      this.statusesService.deleteStatus(id).then(() => {
        this.router.navigate([`/core/statuses`]);
      })
    }

  }

  // Title for the page
  get title(): string {
    return this.isNew ? 'Add New Status' : 'Edit Status';
  }

  // Checking for params to ensure given id is a number
  get isNew() {
    return this.id === 'new' && isNaN(Number(this.id));
  }
}
