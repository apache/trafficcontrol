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

# Traffic Portal Angular7
This project was generated with [Angular CLI](https://github.com/angular/angular-cli) version 7.1.4.

Traffic Portal Angular7 is an attempt at making a self-service-first UI for Apache Traffic Control CDN architectures. The eventual goal is to slowly eat Traffic Portal by importing all of its functionality once the self-service UI is complete.

## Setup
First, clone the [Apache Traffic Control repository](git://github.com/apache/trafficcontrol). This needs to be done, because this server imports the Apache Traffic Control logo from the repository.

### Prerequisites
Traffic Portal Angular7 runs on [NodeJS](https://nodejs.org/) version 10 (8 theoretically works also - untested) and uses its built-in NPM package manager to manage dependencies.

### Building and Running
To set up the Angular7 project, first install all dependencies, then build and run the server-side and client-side modules.

E.g. using the [Angular CLI](https://github.com/angular/angular-cli)

```bash
# If you don't want the development dependencies, add `--production`
# (only needs to be done once unless dependencies change)
npm install

# If you installed @angular/cli globally (`npm install -g @angular/cli`), then you can use these
# commands to build the modules

# Add `--prod` for production deployment
ng build

# For production deployment, run ng run traffic-portal-angular7:server:production
ng run traffic-portal-angular7:server

# Necessary to pull in out-of-scope asset
cp -f ../../misc/logos/ATC-SVG.svg dist/browser/

# Runs the server locally at http://localhost:4000
node dist/server.js
```
Note that calling `node` manually like this will allow you to pass command line parameters to the server. Currently, the only supported parametor is an optional hostname that refers to a Traffic Ops instance. This may or may not include a schema and port, e.g. `https://trafficops.infra.ciab.test`, `trafficops.infra.ciab.test` and `trafficops.infra.ciab.test:443` are all equivalent (because the default schema is `https://` and the default port is 443). If this value is not provided on the command line, the value of the `TO_URL` environment variable will be used. If the Traffic Ops server host is not defined either in the environment or on the command line (e.g. when using the below `npm` commands without `TO_URL` set), an error will be issued and the server will refuse to start.

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

#### Debug Mode
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

By default this will set up an Angular debug-mode server to run, listening on port 4200 (as opposed to the production-mode default of 4000 [presumably so that you could run both at the same time]). Note that regardless of production-mode SSL configuration (TODO), this will **only serve unencrypted HTTP responses by default**. Also, unlike production mode which compiles things ahead of time, this will compile resources on-the-fly so that making changes to a file is immediately "live" without requiring the developer to restart the debug server. Pretty neat, desu.

## Contributing
This project uses `tslint` and an `.editorconfig` file for maintaining code style. If your editor doesn't support `.editorconfig` (VS Code does out-of-the-box as I understand, but Vim and Sublime Text need plugins. Atom ~~doesn't~~ shouldn't exist) then you'll want to manually configure it so as to avoid linting errors. There's quite a bit going on, but the big ones are:

These apply to all files:

- No trailing whitespace before line-endings
- Unix line-endings
- Ensure line ending at end-of-file (POSIX-compliance)
- Tabs not spaces for indentation (spaces may be used for alignment with multi-line statements AFTER indenting appropriately)

These apply to Typescript, specifically

- Don't use `var` - only `const` and `let` are allowed
- Prefer single quotes to double quotes for string literals
- *Document your code - we use JSDoc here*

## Supporting old Traffic Ops versions
This UI is built to work with an API at version 1.5. All endpoints will use this version by default, so when pointing it at a server that only supports e.g. a max of 1.4, you'll need to do something heinous: edit a source file. In the [`src/app/services/api.service.ts`](./src/app/services/api.service.ts) file, change the line `public API_VERSION = '1.5';` to the appropriate version, e.g. `public API_VERSION = '1.4';`. This will be easier in the future<sup>Citation needed</sup>.

## Code scaffolding
Run `ng generate component component-name` to generate a new component. You can also use `ng generate directive|pipe|service|class|guard|interface|enum|module`.

## Running unit tests
Run `ng test` to execute the unit tests via [Karma](https://karma-runner.github.io).

## Running end-to-end tests
Run `ng e2e` to execute the end-to-end tests via [Protractor](http://www.protractortest.org/).

## Further help
To get more help on the Angular CLI use `ng help` or go check out the [Angular CLI README](https://github.com/angular/angular-cli/blob/master/README.md).
