#### Traffic Ops UI 2.0 Experimental

##### Prerequisites

###### Ruby and Compass installed

-- ruby -v (brew install ruby)

-- compass -v (gem install compass)

###### Node, Grunt and Bower installed

-- node -v (brew install node)

-- grunt --version (npm install -g grunt-cli)

-- bower -v (npm install -g bower)

##### Install

You'll need to install the node modules defined in ./package.json for the build process so:

cd ./

npm install

You'll need to install the components required for the app to run (i.e. angular.js) defined in bower.json so:

cd ./

bower install

Finally, create the build artifacts and start up an Express server for traffic ops to run on localhost:8080:

cd ./

grunt

go to localhost:8080

Once you do all this you'll end up with these directories (none of which are stored in the code repository):

./node_modules

./app/bower_components

./app/dist

API destination is found in:

./conf/config.js

##### Source files

Source files are found in:

./app/src