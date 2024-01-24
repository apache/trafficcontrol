<!--
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
-->

# Traffic Portal
This project was generated with
[Angular CLI](https://github.com/angular/angular-cli) version 10.2.0.

This is an attempt at making a self-service-first UI for Apache Traffic Control
CDN architectures. The eventual goal is to slowly eat Traffic Portal by
importing all of its functionality once the self-service UI is complete.

## Prerequisites
Traffic Portal runs on [NodeJS](https://nodejs.org/) version 16 (or later) and
uses its built-in NPM package manager to manage dependencies.

## Building and Running
To set up the Angular project for a production or testing environment

1. install all dependencies with `npm install` (`pnpm` is **not** supported with
`ngcc` so it won't work)

	```bash
	# If you don't want the development dependencies, add `--production`
	# (only needs to be done once unless dependencies change)
	npm install
	```

1. build and run the server-side (optional; production or production-like) and
client-side modules.

	E.g. running in server-side-rendering mode using the
	[Angular CLI](https://github.com/angular/angular-cli)

	```bash
	# Add `--prod` for production deployment
	ng build

	# For production deployment, run ng run traffic-portal:server:production
	ng run traffic-portal:server

	# Runs the server locally at http://localhost:4000
	node dist/server.js
	```

	E.g. running in server-side-rendering mode using NPM scripts

	```bash
	# These commands don't require a globally-available `@angular/cli` install,
	# but shadow intermediate steps (mostly those are also available as NPM
	# scripts, check out package.json to see them)

	# This builds for production deployment by default
	npm run build:ssr

	# Runs the server locally at http://localhost:4000
	npm run serve:ssr
	```

### Server Command Line Arguments
By default, the Traffic Portal server will attempt to connect to a Traffic Ops server running at the location specified by the `TO_URL` environment variable. If one is not set or you wish to override it, this (and other behavior) can be configured by passing arguments on the command line. To see the available command line options, pass the `-h`/`--help` flag to the server, e.g.

```console
$ CMD="node dist/server.js" # The arguments can be passed to either command, use whichever you like
$ CMD="npm run serve:ssr --" # (`--` signals the end of options to `npm`)
$ $CMD --help # output is verbatim, but omits some things npm might echo e.g. the actual command
usage: server.js [-h] [-v] [-t TRAFFIC_OPS] [-k] [-p PORT] [-c CERT_PATH -K KEY_PATH] [-C CONFIG_FILE]

Traffic Portal re-written in modern Angular.

Optional arguments:
  -h, --help            Show this help message and exit.
  -v, --version         Show program's version number and exit.
  -t TRAFFIC_OPS, --traffic-ops TRAFFIC_OPS
                        Specify the Traffic Ops host/URL, including port.
                        (Default: uses the `TO_URL` environment variable)
  -k, --insecure        Skip Traffic Ops server certificate validation. This
                        affects requests from Traffic Portal to Traffic Ops
                        AND signature verification of any passed SSL
  -p PORT, --port PORT  Specify the port on which Traffic Portal will listen
                        (Default: 4200)
  -d DIST_PATH, --browser-folder DIST_PATH
                        Specifiy locaiton for the folder that holds the
                        browser files
  -c CERT_PATH, --cert-path CERT_PATH
                        Specify a location for an SSL certificate to be used
                        by Traffic Portal. (Requires `-K`/`--key-path`. If
                        both are omitted, will serve using HTTP)
  -K KEY_PATH, --key-path KEY_PATH
                        Specify a location for an SSL certificate to be used
                        by Traffic Portal. (Requires `-c`/`--cert-path`. If
                        both are omitted, will serve using HTTP)
  -C CONFIG_FILE, --config-file CONFIG_FILE
                        Specify a path to a configuration file - options are
                        overridden by command-line flags.
$
```

### Debug Mode
Because we need to proxy API requests back to Traffic Ops, running in debug mode
is a bit more involved than it normally would be. Specifically, unless your
Traffic Ops instance is listening at `https://localhost:6443/` it'll require
making a new file somewhere with the following information in JSON format:

```json
{
	"/api": {
		"target": "Traffic Ops server URL here - e.g. https://trafficops.apache.test",
		"secure":
		 "This should be one of the literal boolean values 'true' or 'false' to indicate if certificate authenticity should be checked."
	}
}
```
More documentation on the configuration options available in this file can be
found in
[the relevant section of the angular-cli documentation](https://github.com/angular/angular-cli/blob/master/docs/documentation/stories/proxy.md).
This step isn't necessary in a production deployment because the proxy is built
into the server-side rendering server.

By default the file `proxy.json` in the project directory (i.e. the same one as
this README.md file) will be used to define proxy settings, which causes it to
expect a TO instance at `https://localhost:6443` (CDN-in-a-Box's default Traffic
Ops port).  To change where this looks for a TO instance, make your own proxy
configuration file and pass it in with `--proxy-config /path/to/proxy.json`.

```bash
# Using globally installed Angular CLI
## Using default proxy
ng serve
## Using custom proxy settings
ng serve --proxy-config /path/to/custom-proxy.json
# using NPM scripts:
## Using default proxy
npm start
## Using custom proxy settings
npm start -- --proxy-config /path/to/custom-proxy.json
```

By default this will set up an Angular debug-mode server to run, listening on
port 4200. Note that regardless of production-mode SSL configuration, this will
**only serve unencrypted HTTP responses by default**. Also, unlike production
mode which compiles things ahead of time, this will compile resources on-the-fly
so that making changes to a file is immediately "live" without requiring the
developer to restart the debug server.

**This debugging mode server is NOT safe for production environments - not only
does it not server HTTPS by default but the server itself _is not audited for
security flaws_ - use this for development and testing ONLY.**

## Running the Tests
Coverage is pretty abysmal at the moment, but unit tests can be run using the
[Angular CLI](https://github.com/angular/angular-cli). To run the unit tests,
use the command `ng test` (dependencies must first be installed). This will
attempt to open Chrome, Firefox and Opera, so ideally you would have those
installed prior to running the tests.

End-to-end testing uses [Cypress](https://www.cypress.io/) and can be run by
using `ng e2e`. More detailed instructions can be found in the README in the
`cypress/` directory.

## Extending Traffic Portal
Traffic Portal supports extending functionality through the use of Angular modules.
The `Custom` module (located at `src/app/custom/`) contains the code to do so and any additional
functionality should be added here as you would to any other Angular module. By default,
this module is not built or included in the main bundle, to enable this modify the environment
(`src/environments`) variable `customModule` to true.

## Contributing
This project uses `eslint` and an `.editorconfig` file for maintaining code
style. If your editor doesn't support `.editorconfig` (VS Code does
out-of-the-box as I understand, but Vim and Sublime Text need plugins) then
you'll want to manually configure it so as to avoid linting errors. There's
quite a bit going on, but the big ones are:

These apply to all files:

- No trailing whitespace before line-endings
- Unix line-endings
- Ensure line ending at end-of-file (POSIX-compliance)
- Tabs not spaces for indentation

These apply to Typescript, specifically

- Don't use `var` - only `const` and `let` are allowed
- Prefer double quotes to single quotes for string literals
- *Document your code - we use JSDoc here*

Code _must_ pass linting to be accepted. To run the linter:
```bash
# Using Angular CLI
ng lint
# Using NPM scripts
npm run lint
```

## Supporting old Traffic Ops versions
This UI is built to work with an API at version 3.0 in development mode
(configured in `./src/environments/environment.ts`) and 2.0 in production mode
(configured in `./src/environments/environment.prod.ts`). All endpoints will use
this version by default, so when pointing it at a server that only supports e.g.
a max of 1.5, you'll need to do edit the source file for your environment.

## Browser Support
This UI obviously requires Javascript, but beyond that the hope is that any
HTML5/CSS3/DOM3-compliant browser should work. Specifically, the
`./.browserslistrc` file defines our browser support, but in terms of browsers
supported without regard for versioning, our aim is:

- Mozilla Firefox
- Google Chrome
- Chromium
- Opera
- Vivaldi
- Microsoft Edge

Internet Explorer and Safari are notably absent from this list. These browsers
are standard-defying nightmares, and we refuse to support browsers that are
end-of-life and/or do not recevie fixes for critical bugs.

The goal is to continuously support these browser in their latest and
penultimate major release versions.

## Code scaffolding
Run `ng generate component component-name` to generate a new component. You can
also use `ng generate directive|pipe|service|class|guard|interface|enum|module`.

It may ask you to then specify a module, because it doesn't know if the thing
you're generating should be imported into the client-side code or the
server-side code - nearly always the thing you're trying to do should be done on
the client-side, so point it to the absolute or relative location of the
`src/app` directory.

However, the files generated via this scaffolding **will** fail linting.
Generally most of those errors can be fixed automatically with `ng lint --fix`
(or equivalently `npm run lint -- --fix`), though.

## Further help
To get more help on the Angular CLI use `ng help` or go check out the
[Angular CLI README](https://github.com/angular/angular-cli/blob/master/README.md).
