
FROM centos:7
RUN printf "\n\
[google-chrome] \n\
name=google-chrome \n\
baseurl=https://dl.google.com/linux/chrome/rpm/stable/x86_64 \n\
enabled=1 \n\
gpgcheck=1 \n\
gpgkey=https://dl-ssl.google.com/linux/linux_signing_key.pub \n\
" > /etc/yum.repos.d/google-chrome.repo

RUN yum update -y

RUN yum install -y  unzip java-1.8.0-openjdk git google-chrome-stable bind-utils net-tools
RUN yum install -y google-chrome-stable

# RUN curl -LO https://nodejs.org/dist/v10.15.0/node-v10.15.0-linux-x64.tar.xz

# # Use bsdtar, because regular tar doesn't work with Docker
# RUN yum install -y bsdtar
# RUN bsdtar --strip-components 1 -xzvf node-v* -C /usr/local

RUN version=$(curl -s https://chromedriver.storage.googleapis.com/LATEST_RELEASE) && \
wget -N https://chromedriver.storage.googleapis.com/${version}/chromedriver_linux64.zip -P ~/ && \ 
 unzip ~/chromedriver_linux64.zip -d ~/ && \
 rm ~/chromedriver_linux64.zip && \
  mv -f ~/chromedriver /usr/local/bin/chromedriver && \
 chown root:root /usr/local/bin/chromedriver && \
 chmod 0755 /usr/local/bin/chromedriver && \
 curl -sL https://rpm.nodesource.com/setup_10.x | bash - && \
  yum update && \
  yum install -y nodejs && \
  rm -rf /var/lib/apt/lists/* && \
  npm install got@latest -g
# Install Protractor and initialized Webdriver
RUN npm install -g protractor@^5.4 && \
  npm install typescript@3.6.4 -g && \
  webdriver-manager update 
  
COPY ./trafficportal /trafficportal
