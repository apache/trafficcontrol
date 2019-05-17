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
import {existsSync, readFileSync, statSync} from 'fs';
import {request as HTTPRequest} from 'http';
import {request as HTTPSRequest, createServer} from 'https';
import {join} from 'path';
import {parse} from 'url';
import * as zlib from 'zlib';

import {ArgumentParser} from 'argparse';

import {environment} from './src/environments/environment';

// Faster server renders w/ Prod mode (dev mode never needed)
enableProdMode();

const VERSION = '4.0.0';

const parser = new ArgumentParser({
	version: VERSION,
	addHelp: true,
	description: 'A re-imagining of Traffic Portal with server-side rendering in Angular7.'
});
parser.addArgument(['-t', '--traffic-ops'], {
	help: 'Specify the Traffic Ops host/URL, including port. (Default: uses the `TO_URL` environment variable)',
	type: (arg: string) => {
		try {
			return new URL(arg);
		} catch (e) {
			if (e instanceof TypeError) {
				return new URL('https://' + arg);
			}
			throw e;
		}
	}
});
parser.addArgument(['-k', '--insecure'], {
	help: 'Skip Traffic Ops server certificate validation.' +
	      'This affects requests from Traffic Portal to Traffic Ops AND signature verification of any passed SSL keys/certificates',
	action: 'storeTrue'
});
parser.addArgument(['-p', '--port'], {
	help: 'Specify the port on which Traffic Portal will listen (Default: 4200)',
	type: Number,
	defaultValue: 4200
});
parser.addArgument(['-c', '--cert-path'], {
	help: 'Specify a location for an SSL certificate to be used by Traffic Portal. (Requires `-K`/`--key-path`.' +
	      ' If both are omitted, will serve using HTTP)',
	type: String
});
parser.addArgument(['-K', '--key-path'], {
	help: 'Specify a location for an SSL certificate to be used by Traffic Portal. (Requires `-c`/`--cert-path`.' +
	      ' If both are omitted, will serve using HTTP)',
	type: String
});

const args = parser.parseArgs();

if (isNaN(args.port) || args.port <= 0 || args.port > 65535) {
	console.error('Invalid listen port:', args.port);
	process.exit(1);
}

let to_url: URL;
if (args.traffic_ops) {
	to_url = args.traffic_ops;
} else if (process.env.hasOwnProperty('TO_URL')) {
	try {
		to_url = new URL((process.env as any).TO_URL);
	} catch (e) {
		console.error('Invalid Traffic Ops URL set in environment variable:', (process.env as any).TO_URL);
		process.exit(1);
	}
} else {
	console.error('Must define a Traffic Ops URL, either on the command line or TO_URL environment variable');
	process.exit(1);
}

let to_host: string;
let to_port: number;
let to_use_SSL: boolean;

if (!to_url.hostname || to_url.hostname.length <= 0) {
	console.error("'%s' is not a valid Traffic Ops URL! (hint: try -h/--help)", to_url.href);
	process.exit(1);
}
to_host = to_url.hostname;

if (to_url.protocol) {
	switch (to_url.protocol) {
		case 'http:':
			to_use_SSL = false;
			break;
		case 'https:':
			to_use_SSL = true;
			break;
		default:
			console.error("Unknown/unsupported protocol: '%s'", to_url.protocol);
			process.exit(1);
	}
} else {
	to_use_SSL = true;
}

if (to_url.port) {
	to_port = Number(to_url.port);
	if (isNaN(to_port) || to_port > 65535 || to_port <= 0) {
		console.error('Invalid port: ', to_port);
		process.exit(1);
	}
} else if (to_use_SSL) {
	to_port = 443;
} else {
	to_port = 80;
}

const TO_URL = 'http' + (to_use_SSL ? 's' : '') + '://' + to_host + ':' + String(to_port);

if ((args.cert_path && !args.key_path) || (!args.cert_path && args.key_path)) {
	console.error('Either both `-c`/`--cert-path` and `-K`/`--key-path` must be given, or neither.');
	process.exit(1);
}
const serveSSL = args.cert_path && args.key_path;

console.debug('Traffic Ops server at:', TO_URL);

// Ignore untrusted certificate signers
(process.env as any).NODE_TLS_REJECT_UNAUTHORIZED = args.insecure ? '0' : '1';

const request = to_use_SSL ? HTTPSRequest : HTTPRequest;

console.debug('Pinging Traffic Ops server...');
const pingRequest = request({
		host:   to_host,
		port:   to_port,
		path:   '/api/1.4/ping',
		method: 'GET'
	},
	response => {
		if ((response as any).aborted || (response as any).statusCode !== 200) {
			console.error("Failed to ping Traffic Ops server! Is '%s' correct?", TO_URL);
			if (response.hasOwnProperty('statusCode') && response.hasOwnProperty('statusMessage')) {
				console.debug('Response status code was', (response as any).statusCode, (response as any).statusMessage);
			}
			response.pipe(process.stderr);
			process.exit(2);
		}
		console.debug('Ping succeeded.');
	}
);
pingRequest.on('error', e => {
	console.error('Failed to contact Traffic Ops server!');
	console.error(e);
	process.exit(2);
});
pingRequest.end();

// Read in SSL key/cert if present.
let key: string;
let cert: string;
if (serveSSL) {
	if (!existsSync(args.key_path)) {
		console.error('%s: no such file or directory', args.key_path);
		process.exit(1);
	}
	if (statSync(args.key_path).isDirectory()) {
		console.error('%s: is a directory', args.key_path);
		process.exit(1);
	}
	if (!existsSync(args.cert_path)) {
		console.error('%s: no such file or directory', args.cert_path);
		process.exit(1);
	}
	if (statSync(args.cert_path).isDirectory()) {
		console.error('%s: is a directory', args.cert_path);
		process.exit(1);
	}
	try {
		key = readFileSync(args.key_path, 'utf8');
		cert = readFileSync(args.cert_path, 'utf8');
	} catch (e) {
		console.error('An error occurred reading SSL certificate/key files:', e);
		process.exit(1);
	}
}


// Express server
const app = express();
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

	const fwdRequest = {
		host:    to_host,
		port:    to_port,
		path:    parse(req.originalUrl).path,
		method:  req.method,
		headers: req.headers
	};

	try {
		const proxiedRequest = request(fwdRequest, (r) => {
			res.writeHead(r.statusCode, r.headers);
			r.pipe(res);
		});
		req.pipe(proxiedRequest);
	} catch (e) {
		console.error(e);
		res.end();
		req.end();
	}
});

// Default route shows the dash
app.get('*', (req, res) => {
	try {
		res.render('index', { req });
	} catch (e) {
		console.error(e);
		res.end();
	}
});

app.enable('trust proxy');

// Start up the Node server
function logMsg () {
	if (serveSSL) {
		console.log(`Node Express server listening on https://localhost:${args.port}`);
	} else {
		console.log(`Node Express server listening on http://localhost:${args.port}`);
	}
}
if (serveSSL) {
	createServer({
		key: key,
		cert: cert,
		rejectUnauthorized: !args.insecure,
	}, app).listen(args.port, logMsg);
} else {
	app.listen(args.port, logMsg);
}
