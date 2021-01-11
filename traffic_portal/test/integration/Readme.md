# Traffic Portal Test Automation 

### Test Development Environment Setup
---
* Install [Node](http://nodejs.org) (v6.x.x or later)
* Follow setup steps described [here](http://www.protractortest.org/#/tutorial#setup)
* Now install packages: `npm install`
* Now install protractor `npm i protractor`
* Now install typescript `npm install typescript@3.6.4 -g`
* Now install selenium standalone- `sudo webdriver-manager update`
* In a separate command line window, run `sudo webdriver-manager start` and keep it running.

### Run Tests
---
start [CDN in a Box](https://traffic-control-cdn.readthedocs.io/en/latest/admin/quick_howto/ciab.html)
run `tsc` from traffic_portal/integration directory. This command will compile and convert the typescript files into javascript files. The generated js files are available in traffic_portal/integration/GeneratedCode directory
run `protractor ./GeneratedCode/config.js` from traffic_portal/integration directory. This command will run the protractor test from the CDN in a box environment.

