import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { StatusesTableComponent } from './statuses-table.component';
import { RouterModule } from '@angular/router';
import { SharedModule } from 'src/app/shared/shared.module';
import { MatCardModule } from '@angular/material/card';
import { StatusesService } from 'src/app/api/statuses.service';
import { FormsModule } from '@angular/forms';

const StatusesTableRouting = RouterModule.forChild([
  {
    path: '',
    component: StatusesTableComponent
  }
]);

@NgModule({
  declarations: [
    StatusesTableComponent
  ],
  imports: [
    CommonModule,
    StatusesTableRouting,
    FormsModule,
    MatCardModule,
    SharedModule
  ],
  providers:[
    StatusesService
  ]
})
export class StatusesTableModule { }
