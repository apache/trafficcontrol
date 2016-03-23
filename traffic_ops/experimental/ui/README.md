#### Traffic Ops UI 2.0

An AngularJS client served from a lightweight Node.js web server. TO UI 2.0 was designed to consume the TO 2.0 RESTful API.

##### Installation

1. Install prerequisite software (OS X)

    - brew install node
    - npm install -g compass
    - npm install -g bower
    - npm install -g grunt-cli
    
2. Navigate to UI root

3. Load application dependencies

    - npm install

4. Load client-side dependencies

    - bower install
    
##### Configuration

1. Configure the Node.js web server to proxy api requests to the API URL

    - vim ./conf/config.js
    - set api.base_url to http://api-domain.com or leave default value
    
##### Run

1. Package, deploy and start Node.js server

    - grunt

2. Head over to http://localhost:8080

##### Notes

    - Node.js server configuration is found in ./server/server.js
    - Source files are found in ./app/src
    - Build artifacts are found in ./app/dist
