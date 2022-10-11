/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
import { Component, Inject } from "@angular/core";
import { MAT_DIALOG_DATA } from "@angular/material/dialog";

/**
 * Contains the structure of the data the TextDialogComponent expects
 */
export interface TextDialogData {
	title: string;
	message: string;
}

/**
 * TextDialogComponent contains code for displaying a simple text mat-dialog.
 */
@Component({
	selector: "tp-text-dialog",
	styleUrls: ["./text-dialog.component.scss"],
	templateUrl: "./text-dialog.component.html"
})
export class TextDialogComponent {

	constructor(@Inject(MAT_DIALOG_DATA) public readonly dialogData: TextDialogData) {
	}
}
