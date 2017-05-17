FROM node:4-onbuild
#FROM buildpack-deps:jessie

RUN apt-get update -y && apt-get install libffi-dev ruby-dev rubygems vim -y

# replace this with your application's default port
RUN gem update --system && gem install --no-rdoc --no-ri compass && gem install --no-rdoc --no-ri sass -v 3.4.22
RUN npm install -g bower grunt-cli
RUN echo '{ "allow_root": true }' > /root/.bowerrc
RUN cd /usr/src/app && /usr/local/bin/bower install
RUN cd /usr/src/app && /usr/local/bin/grunt dist

EXPOSE 8888
