FROM ubuntu:latest
MAINTAINER ramiro@rarenivar.com

RUN locale-gen en_US.UTF-8
RUN mkdir -p src
RUN cd src
RUN apt-get install -y curl
RUN curl https://install.meteor.com/ | sh
RUN meteor update
RUN apt-get clean
COPY . /src
WORKDIR src/

CMD ["meteor"]
