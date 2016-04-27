FROM node:5.10.1

# Need to run this to see what's in the various sources
RUN apt-get -qq update

# Get avconv
RUN apt-get install -y libav-tools

# Create app directory
RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app

COPY package.json /usr/src/app/
RUN npm install
COPY . /usr/src/app
COPY ./mongo.json /usr/src/app

CMD ["npm", "start"]
