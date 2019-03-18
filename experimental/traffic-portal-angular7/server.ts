/*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*/

import 'zone.js/dist/zone-node';
import {enableProdMode} from '@angular/core';
// Express Engine
import {ngExpressEngine} from '@nguniversal/express-engine';
// Import module map for lazy loading
import {provideModuleMap} from '@nguniversal/module-map-ngfactory-loader';

import * as express from 'express';
import {parse} from 'url';
import {request} from 'https';
import {join} from 'path';
import * as zlib from 'zlib';

import {environment} from './src/environments/environment';

// Faster server renders w/ Prod mode (dev mode never needed)
enableProdMode();

// Get a URL for a Traffic Ops instance
let to_url_raw = '';
if (process.argv.length >= 3) {
	to_url_raw = process.argv[2];
}
else if (process.env.hasOwnProperty("TO_URL")) {
	to_url_raw = process.env.TO_URL;
}
else {
	console.error("Must define a Traffic Ops URL, either on the command line or TO_URL environment variable");
	process.exit(1);
}

let to_port;
let to_host;
const to_url_split = to_url_raw.split(':', 2);
if (to_url_split.length === 1) {
	to_host = to_url_split[0];
	to_port = 443;
}
else if (to_url_split.length === 2) {
	if (to_url_split[0].toLowerCase() === 'https') {
		to_host = to_url_split[1];
		if (to_host.length < 3) {
			console.error("Malformed Traffic Ops URL:", to_url_raw);
			process.exit(1);
		}
		to_host = to_host.slice(2);
		to_port = 443;
	}
	else {
		to_host = to_url_split[0];
		try {
			to_port = Number(to_url_split[1]);
		}
		catch (e) {
			console.error("Malformed Traffic Ops URL:", to_url_raw);
			console.debug("Exception:", e);
			process.exit(1);
		}
	}
}
else {
	to_host = to_url_split[1];
	if (to_host.length < 3) {
		console.error("Malformed Traffic Ops URL:", to_url_raw);
		process.exit(1);
	}
	to_host = to_host.slice(2);

	try {
		to_port = Number(to_url_split[2]);
	}
	catch (e) {
		console.error("Malformed Traffic Ops URL:", to_url_raw);
		console.debug("Exception:", e);
		process.exit(1);
	}
}

console.debug("TO_HOST:", to_host, "TO_PORT:", to_port);

// Ignore untrusted certificate signers (TODO: this should be an option)
process.env["NODE_TLS_REJECT_UNAUTHORIZED"] = '0';

// Express server
const app = express();

const PORT = process.env.PORT || 4000;
const DIST_FOLDER = join(process.cwd(), 'dist/browser');

// * NOTE :: leave this as require() since this file is built Dynamically from webpack
const {AppServerModuleNgFactory, LAZY_MODULE_MAP} = require('./dist/server/main');

// Our Universal express-engine (found @ https://github.com/angular/universal/tree/master/modules/express-engine)
app.engine('html', ngExpressEngine({
	bootstrap: AppServerModuleNgFactory,
	providers: [
		provideModuleMap(LAZY_MODULE_MAP)
	]
}));

app.set('view engine', 'html');
app.set('views', DIST_FOLDER);

// When in a dev environment, serve changes quickly
let m = '1s';
if (environment.production) {
	m = '1y';
}

// Static files
app.get('*.*', express.static(DIST_FOLDER, {
	maxAge: m
}));

// Forward API requests to Traffic Ops
// Note that this doesn't handle compression/encoding, just transparently
// proxies arbitrary data
app.use('/api/**', (req, res) => {
	console.debug(`Making TO API request to \`${req.originalUrl}\``);

	let fwdRequest = {
		host: to_host,
		port: to_port,
		path: parse(req.originalUrl).path,
		method: req.method,
		headers: req.headers
	};


	const proxiedRequest = request(fwdRequest, (r) => {
		res.writeHead(r.statusCode, r.headers);
		r.pipe(res);
	});
	req.pipe(proxiedRequest);
});

// Default route shows the dash
app.get('*', (req, res) => {
	res.render('index', { req });
});

app.enable('trust proxy');

// Start up the Node server
app.listen(PORT, () => {
	console.log(`Node Express server listening on http://localhost:${PORT}`);
});
