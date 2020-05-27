import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ApplyYamlComponent } from './apply-yaml.component';

describe('ApplyYamlComponent', () => {
  let component: ApplyYamlComponent;
  let fixture: ComponentFixture<ApplyYamlComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ApplyYamlComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ApplyYamlComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
