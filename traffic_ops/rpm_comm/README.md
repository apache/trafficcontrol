# Create a TO RPM with Dependencies

1. Download repo Traffic Control (or from your favorite repo)
```
$ git clone http://github.com/Comcast/traffic_control.git
```
2. Bring up Vagrant environment (http://www.vagrantup.com)
```
$ cd <repo dir>
$ cp traffic_control/traffic_ops/rpm_comm/Vagrantfile ./
$ vagrant up
```
3. ssh into vagrant environment
```
$ vagrant ssh
```
4. Build the RPM

```
$ cd /vagrant/traffic_control/rpm_comm/build
$ ./build.sh -h
./build.sh [-b <branch>] [-c | -cc] [-g <gitrepo> | -r <repodir>]
[-w <working directory]

Don't run this script ever.
   -b  | --branch         Git branch
                          default: master
   -c  | --clean          Make a fresh start but, leave carton
   -cc | --clean_carton   Make a fresh start
   -g  | --gitrepo        Git repository.
                          default:
                             https://github.com/Comcast/traffic_control.git
   -h  | --help           Print this message
   -r  | --repodir        Location of an already downloaded traffic_control
                          repo (e.g. /vagrant/traffic_control). This overrides
                          the -g option. Note that the script expects this to
                          be the root directory of the project and in the
                          correct branch.
   -w  | --workspace      Working directory
                          default: home directory
===============================================================================
$ ./build.sh -b my-feature-branch -g https://github.com/myuser/traffic_control.git -c
```

Notes:  
This is known to work with CentOS 6.6 as the Vagrant environment.
