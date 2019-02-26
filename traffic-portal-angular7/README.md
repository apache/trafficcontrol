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

Traffic Portal Angular7 is a re-write of Traffic Portal in the latest, stable release of Angular. The goal is for rendering to take place mainly server-side, and while many pages will have Javascript-only features (e.g. graphs, sorting tables etc.) attemtps will be made to limit client-side scripting as much as possible. Ideally, most pages will work without Javascript at all, but this early in the project it's unclear whether that is realistic.

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
cp -f ../misc/logos/ATC-SVG.svg dist/browser/

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

## Code scaffolding
Run `ng generate component component-name` to generate a new component. You can also use `ng generate directive|pipe|service|class|guard|interface|enum|module`.

## Running unit tests
Run `ng test` to execute the unit tests via [Karma](https://karma-runner.github.io).

## Running end-to-end tests
Run `ng e2e` to execute the end-to-end tests via [Protractor](http://www.protractortest.org/).

## Further help
To get more help on the Angular CLI use `ng help` or go check out the [Angular CLI README](https://github.com/angular/angular-cli/blob/master/README.md).
