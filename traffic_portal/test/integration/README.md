# Traffic Portal Test Automation 

**Test Development Environment Setup**

* Install [Node](http://nodejs.org) (v6.x.x or later): `brew install node`
* Follow setup steps described to install protractor: [here](http://www.protractortest.org/#/tutorial#setup)
* Now install packages manager: `npm install`
* Now install protractor: `npm i protractor`
* Now install typescript: `npm install typescript@3.6.4 -g`
* Now install selenium standalone: `sudo webdriver-manager update`
* In a separate command line window, run: `sudo webdriver-manager start` and keep it running.
* Run CDN-in-a-Box in a separate command line window and make sure all the components and features display.

**How To Run Tests**

Run this command below from integration directory. This command will compile and convert the typescript files into javascript files. The generated js files are available in integration/GeneratedCode directory.
```
tsc
```
After that, run this command below to run the protractor test from the environment user input.
```
protractor ./GeneratedCode/config.js --params.baseUrl='https://localhost:443'
```
**Command Line Parameters**

| Flag                            | Description                                                                                          |
| ------------------------------- | :--------------------------------------------------------------------------------------------------: |
| params.baseUrl                  | Environment test run on. Tests are written for cdn-in-a-box only. Do not run on other environment                                   |
| capabilities.shardTestFiles     | Input `true` or `false` to turn on or off parallelization. If the value is false, maxInstances will always count as 1. The default value in the config file = false                            |
| capabilities.maxInstances       | Input number of chrome instances that your machine can handle. Test will fail if local machine cannot handle a lot of chrome instances. The default value = 1    |
