import {
  Component,
  OnInit,
  OnChanges,
  EventEmitter,
  Output,
  Input,
  SimpleChanges,
} from '@angular/core';
import { EditorView } from '../../../models/content';

@Component({
  selector: 'app-apply-yaml',
  templateUrl: './apply-yaml.component.html',
  styleUrls: ['./apply-yaml.component.sass'],
})
export class ApplyYamlComponent implements OnInit, OnChanges {
  isOpenValue: boolean;
  view: EditorView;

  @Input()
  get isOpen() {
    return this.isOpenValue;
  }

  set isOpen(v: boolean) {
    this.isOpenValue = v;
    this.isOpenChange.emit(this.isOpenValue);
  }

  @Output()
  isOpenChange = new EventEmitter<boolean>();

  onCreate() {
    console.log('test');
    this.isOpenValue = true;
  }

  onCancel() {
    this.isOpen = false;
  }

  onSubmit() {
    this.isOpen = false;
  }

  constructor() {
    this.view = {
      config: {
        value: 'test: 123',
        language: 'yaml',
        readOnly: false,
        metadata: {},
      },
      metadata: { type: 'editor' },
    };
  }

  ngOnChanges(changes: SimpleChanges): void {}

  ngOnInit(): void {}
}
