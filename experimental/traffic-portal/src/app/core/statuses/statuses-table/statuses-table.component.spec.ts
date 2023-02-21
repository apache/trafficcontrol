import { ComponentFixture, TestBed } from '@angular/core/testing';

import { StatusesTableComponent } from './statuses-table.component';

describe('StatusesTableComponent', () => {
  let component: StatusesTableComponent;
  let fixture: ComponentFixture<StatusesTableComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ StatusesTableComponent ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(StatusesTableComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
