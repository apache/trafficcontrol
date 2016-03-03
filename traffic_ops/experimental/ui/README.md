#### Traffic Ops UI 2.0

An AngularJS client served from a lightweight Node.js web server. TO UI 2.0 was designed to consume the TO 2.0 RESTful API.

##### Installation

1. Install prerequisite software

    - brew install ruby
    - gem install compass
    - brew install node
    - npm install -g bower
    - npm install -g grunt-cli

2. Navigate to UI root

    - cd ./

3. Load app dependencies

    - npm install

4. Load client-side dependencies

    - bower install

5. Package, deploy and start Node.js server

    - grunt (for dev mode)
    - grunt dist (for prod mode)

6. Head over to localhost:8080

##### Notes

    - API destination is defined in ./conf/config.js (make sure this points to the TO 2.0 API)
    - Node.js server configuration is found in ./server/server.js
    - Source files are found in ./app/src
    - Build artifacts are found in ./app/dist
