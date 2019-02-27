import { Component, OnInit } from '@angular/core';
import { FormControl } from '@angular/forms';
import { Router, ActivatedRoute } from '@angular/router';
import { first }  from 'rxjs/operators';

import { AuthenticationService } from '../../services';

@Component({
	selector: 'login',
	templateUrl: './login.component.html',
	styleUrls: ['./login.component.scss']
})
export class LoginComponent implements OnInit {
	returnURL: string;

	u = new FormControl('');
	p = new FormControl('');

	constructor(private route: ActivatedRoute, private router: Router, private auth: AuthenticationService) { }

	ngOnInit() {
		this.returnURL = this.route.snapshot.queryParams['returnUrl'] || '/';
		console.log(this);
	}

	submitLogin(): void {
		this.auth.login(this.u.value, this.p.value).subscribe(
			(response) => {
				if (response) {
					console.log("LoginComponent: response:", response);
					this.router.navigate([this.returnURL]);
				}
			},
			(erro) => {
				console.error("LoginComponent: Error:", erro);
			}
		);
	}

}
