import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { SshCellRendererComponent } from './ssh-cell-renderer.component';

describe('SshCellRendererComponent', () => {
  let component: SshCellRendererComponent;
  let fixture: ComponentFixture<SshCellRendererComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ SshCellRendererComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(SshCellRendererComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
