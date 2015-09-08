# Create a TO RPM with Dependencies

1. Download repo Traffic Control (or from your favorite repo)
```
$ git clone http://github.com/Comcast/traffic_control.git
```
2. Bring up Vagrant environment (http://www.vagrantup.com)
```
$ cd <repo dir>
$ cp traffic_control/traffic_ops/rpm/Vagrantfile ./
$ vagrant up
```
3. ssh into vagrant environment
```
$ vagrant ssh
```
4. set environment variables to control build: BRANCH, HOTFIX\_BRANCH,
   WORKSPACE are automatically set to reasonable defaults.  BUILD\_NUMBER is
   created from the build.number file in the rpm directory and is incremented
   automatically.
```
export BRANCH=master
export HOTFIX_BRANCH=hotfix
export WORKSPACE=$HOME/workspace
export BUILD_NUMBER=100
```
5. Build the RPM
```
$ cd /vagrant/traffic_control/rpm/build
$ ./build_rpm.sh
```
Notes:  
This is known to work with CentOS 6.6 as the Vagrant environment.
