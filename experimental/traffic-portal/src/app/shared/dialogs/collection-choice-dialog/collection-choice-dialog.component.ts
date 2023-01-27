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

import type { DecisionDialogData } from "../decision-dialog/decision-dialog.component";

/**
 * Contains the structure of the data that CollectionChoiceDialogComponent
 * accepts.
 */
export interface CollectionChoiceDialogData<T = unknown> extends DecisionDialogData {
	/** The collection from which the user will choose. */
	collection: {
		/** A user-friendly name for the item. */
		label: string;
		/** The actual value of the item that you want back. */
		value: T;
	}[];
	/** If given, will be displayed as a hint to the user. */
	hint?: string | null | undefined;
	/** Used as a label for the input. Defaults to `message`. */
	label?: string | null | undefined;
	/** A prompt for the user so they know what they're choosing and why. */
	message: string;
}

/**
 * This dialog facilitates asking the user to choose an item from a list.
 */
@Component({
	selector: "tp-collection-choice-dialog",
	styleUrls: ["./collection-choice-dialog.component.scss"],
	templateUrl: "./collection-choice-dialog.component.html",
})
export class CollectionChoiceDialogComponent<T = unknown> {

	public selected: T | null = null;

	/** The label to use for the selection input. */
	public get label(): string {
		return this.dialogData.label ?? this.dialogData.message;
	}

	constructor(@Inject(MAT_DIALOG_DATA) public readonly dialogData: CollectionChoiceDialogData<T>) { }

}
