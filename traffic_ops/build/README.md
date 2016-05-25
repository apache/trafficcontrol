# Create a TO RPM with Dependencies

1. Download repo Traffic Control (or from your favorite repo)
```
$ git clone http://github.com/Comcast/traffic_control.git
```
2. Bring up Vagrant environment (http://www.vagrantup.com)
```
$ cd <repo dir>
$ cp traffic_control/traffic_ops/build/Vagrantfile ./
$ vagrant up
```
3. ssh into vagrant environment
```
$ vagrant ssh
```
4. **OPTIONAL** Set environment variables to control build.  All are
   automatically set to reasonable defaults (in parentheses) and it is
   recommended to leave them unset.   They can be overridden if necessary:
   - *BRANCH* (master)
   - *HOTFIX\_BRANCH* (none)
   - *WORKSPACE* (top level of local repo tree -- workspace must be a clone of the repository)
   - *BUILD\_NUMBER* (# of commits in branch + last commit identifier)
5. Build the RPM
```
$ cd /vagrant/traffic_control/traffic_ops/build
$ ./build_rpm.sh
```
Notes:  
This is known to work with CentOS 6.7 as the Vagrant environment.
