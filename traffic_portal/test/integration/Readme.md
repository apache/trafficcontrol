# Traffic Portal Test Automation 

### Test Development Environment Setup
---
* Install [Node](http://nodejs.org) (v6.x.x or later) `brew install node`
* Follow setup steps described [here](http://www.protractortest.org/#/tutorial#setup)
* Now install packages manager: `npm install`
* Now install protractor `npm i protractor`
* Now install typescript `npm install typescript@3.6.4 -g`
* Now install selenium standalone- `sudo webdriver-manager update`
* In a separate command line window, run `sudo webdriver-manager start` and keep it running.
* Run CDN-in-a-Box and make sure all the components and features display.

### Run Tests
---
run `tsc` from integration directory. This command will compile and convert the typescript files into javascript files. The generated js files are available in integration/GeneratedCode directory
run `protractor ./GeneratedCode/config.js --params.baseUrl='https://localhost:443'` from integration directory. This command will run the protractor test from the environment user input.

### Test Automation Framework
![TP Automation Framework](/trafficportal/tp-framework.png)
