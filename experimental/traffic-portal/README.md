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
This project was generated with [Angular CLI](https://github.com/angular/angular-cli) version 10.2.0.

This is an attempt at making a self-service-first UI for Apache Traffic Control CDN architectures. The eventual goal is to slowly eat Traffic Portal by importing all of its functionality once the self-service UI is complete.

## Prerequisites
Traffic Portal runs on [NodeJS](https://nodejs.org/) version 13  and uses its built-in NPM package manager to manage dependencies.

## Building and Running
To set up the Angular project, first install all dependencies, then build and run the server-side and client-side modules.

E.g. using the [Angular CLI](https://github.com/angular/angular-cli)

```bash
# If you don't want the development dependencies, add `--production`
# (only needs to be done once unless dependencies change)
npm install

# If you installed @angular/cli globally (`npm install -g @angular/cli`), then you can use these
# commands to build the modules

# Add `--prod` for production deployment
ng build

# For production deployment, run ng run traffic-portal-angular8:server:production
ng run traffic-portal-angular8:server

# Runs the server locally at http://localhost:4000
node dist/server.js
```

E.g. using NPM scripts

```bash
# If you don't want the development dependencies, add `--production`
# (only needs to be done once unless dependencies change)
npm install

# These commands are less verbose and don't require a globally-available `@angular/cli` install,
# but shadow intermediate steps (mostly those are also available as NPM scripts,
# check out package.json to see them)

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

Note that only certificates for `localhost` are accepted at the time of this writing.

### Debug Mode
Because we need to proxy API requests back to Traffic Ops, running in debug mode is a bit more involved than it normally would be. Specifically, it'll require making a new file somewhere with the following information in JSON format:

```json
{
	"/api": {
		"target": "Traffic Ops server URL here - e.g. https://trafficops.apache.test",
		"secure":
		 "This should be one of the literal boolean values 'true' or 'false' to indicate if certificate authenticity should be checked."
	}
}
```
More documentation on the configuration options available in this file can be found in [the relevant section of the angular-cli documentation](https://github.com/angular/angular-cli/blob/master/docs/documentation/stories/proxy.md). This step isn't necessary in a production deployment because the proxy is built into the server-side rendering server.

Now, assuming this in the project directory (i.e. the same one as this README.md file) and named e.g. `proxy.json`, a debug-mode server can be started like so:

```bash
# Note that this pre-supposes you've globally installed the angular-cli
# e.g. via `npm install -g @angular/cli`
ng serve --proxy-config ./proxy.json
```

By default this will set up an Angular debug-mode server to run, listening on port 4200. Note that regardless of production-mode SSL configuration, this will **only serve unencrypted HTTP responses by default**. Also, unlike production mode which compiles things ahead of time, this will compile resources on-the-fly so that making changes to a file is immediately "live" without requiring the developer to restart the debug server. Pretty neat, desu.

## Running the Tests
Coverage is pretty abysmal at the moment, but unit tests can be run using the [Angular CLI](https://github.com/angular/angular-cli). To run the unit tests, use the command `ng test` (dependencies must first be installed). This will attempt to open Chrome, Firefox and Opera, so ideally you would have those installed prior to running the tests.

End-to-end testing is broken at the time of this writing, but to run it anyway use `ng e2e`.

## Contributing
This project uses `tslint` and an `.editorconfig` file for maintaining code style. If your editor doesn't support `.editorconfig` (VS Code does out-of-the-box as I understand, but Vim and Sublime Text need plugins. Atom ~~doesn't~~ shouldn't exist) then you'll want to manually configure it so as to avoid linting errors. There's quite a bit going on, but the big ones are:

These apply to all files:

- No trailing whitespace before line-endings
- Unix line-endings
- Ensure line ending at end-of-file (POSIX-compliance)
- Tabs not spaces for indentation (spaces may be used for alignment with multi-line statements AFTER indenting appropriately)

These apply to Typescript, specifically

- Don't use `var` - only `const` and `let` are allowed
- Prefer double quotes to single quotes for string literals
- *Document your code - we use JSDoc here*

## Supporting old Traffic Ops versions
This UI is built to work with an API at version 1.5. All endpoints will use this version by default, so when pointing it at a server that only supports e.g. a max of 1.4, you'll need to do something heinous: edit a source file. In the [`src/app/services/api.service.ts`](./src/app/services/api.service.ts) file, change the line `public API_VERSION = '1.5';` to the appropriate version, e.g. `public API_VERSION = '1.4';`. This will be easier in the future<sup>Citation needed</sup>.

## Browser Support
This UI obviously requires Javascript, but beyond that the hope is that any HTML5/CSS3/DOM3-compliant browser should work. Specifically, testing is being done using the latest versions of:

- Google Chrome (not Chromium atm, maybe in the future)
- Opera
- Vivaldi (support tenuous)
- Mozilla Firefox
- Microsoft Edge (Once the rendering engine turns into Chromium - Chakra is not HTML5/CSS3-compliant)

... and the goal is to continuously support these browser in their latest and penultimate major release versions. Safari isn't tested because I don't have a Mac, but support for that would be nice given the high usage by Traffic Control users/admins/devs.

## Code scaffolding
Run `ng generate component component-name` to generate a new component. You can also use `ng generate directive|pipe|service|class|guard|interface|enum|module`.

However, the files generated via this scaffolding **will** fail linting. Generally most of those errors can be fixed automatically with `ng lint --fix`, though.

## Further help
To get more help on the Angular CLI use `ng help` or go check out the [Angular CLI README](https://github.com/angular/angular-cli/blob/master/README.md).
