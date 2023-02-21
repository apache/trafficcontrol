import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { StatusDetailsComponent } from './status-details.component';
import { RouterModule } from '@angular/router';
import { ReactiveFormsModule } from '@angular/forms';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import {MatGridListModule} from '@angular/material/grid-list';
import { MatCardModule } from '@angular/material/card';
import { SharedModule } from 'src/app/shared/shared.module';
import { MatButtonModule } from '@angular/material/button';

const StatusDetailRouting = RouterModule.forChild([
  {
    path: '',
    component: StatusDetailsComponent
  }
]);

@NgModule({
  declarations: [
    StatusDetailsComponent
  ],
  imports: [
    CommonModule,
    StatusDetailRouting,
    ReactiveFormsModule,
    MatFormFieldModule,
    MatInputModule,
    MatGridListModule,
    MatCardModule,
    MatButtonModule,
    SharedModule
  ]
})
export class StatusDetailsModule { }
