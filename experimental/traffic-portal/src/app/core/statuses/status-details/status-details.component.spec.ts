import { ComponentFixture, TestBed } from '@angular/core/testing';

import { StatusDetailsComponent } from './status-details.component';

describe('StatusDetailsComponent', () => {
  let component: StatusDetailsComponent;
  let fixture: ComponentFixture<StatusDetailsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ StatusDetailsComponent ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(StatusDetailsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
