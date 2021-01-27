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
run `tsc` from trafficportal directory. This command will compile and convert the typescript files into javascript files. The generated js files are available in trafficPortal/GeneratedCode directory
run `protractor ./GeneratedCode/config.js --params.baseUrl='https://localhost:443'` from trafficportal directory. This command will run the protractor test from the environment user input.


### Test Automation Framework
![TP Automation Framework](/trafficportal/tp-framework.png)
