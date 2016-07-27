# Traffic Portal Installation

### 1. Build

* Download Traffic Control repo

    ```
    $ mkdir foo
    $ cd foo
    $ git clone http://github.com/Comcast/traffic_control.git
    ```

* Bring up Vagrant environment (http://www.vagrantup.com)

    ```
    $ cp traffic_control/traffic_portal/build/Vagrantfile ./
    $ vagrant up
	```

* ssh into vagrant environment

    ```
    $ vagrant ssh
    ```  

* [OPTIONAL] Set RPM variables

  * BRANCH (enter version (i.e. 1.6.0) or leave default (master))
  * BUILD_NUMBER (defaults to number of git commits)

* Build the RPM

    ```
    $ cd /vagrant/traffic_control/traffic_portal/build
    $ ./build_rpm.sh
    ```

### 2. Install

* Install the RPM

    ```
    $ cd /vagrant/traffic_control/traffic_portal/build
    $ sudo yum install -y traffic_portal-$BRANCH-$BUILD_NUMBER.x86_64.rpm
    ```

### 3. Configure

* Configure Traffic Portal

    ```
    $ cd /etc/traffic_portal/conf
    $ sudo cp config-template.js config.js
    $ sudo vi config.js (read the inline comments)
    ```

### 4. Run

* Start Traffic Portal

    ```
    $ sudo service traffic_portal start
    ```

#### Notes

    - This is known to work with CentOS 6.7 as the Vagrant environment
