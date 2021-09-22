import {Component, EventEmitter, Input, Output, OnInit} from "@angular/core";

/**
 *
 */
@Component({
	selector: "tp-demo",
	styleUrls: ["./demo.component.scss"],
	templateUrl: "./demo.component.html"
})
export class DemoComponent implements OnInit {
	@Input() // Parent -> Child
	public input: string = "";

	@Output() // Child -> Parent
	public output = new EventEmitter<string>();

	// Two way binding
	@Input()
	public twoWay: string = "";

	@Output()
	public twoWayChange = new EventEmitter<string>();

	constructor() { }

	/**
	 *
	 */
	ngOnInit(): void {
		this.output.emit("output");
	}

}
