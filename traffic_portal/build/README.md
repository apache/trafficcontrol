# Traffic Portal Installation

1. Download Traffic Control repo
    ```
    $ mkdir foo
    $ cd foo
    $ git clone http://github.com/Comcast/traffic_control.git
    ```

2. Bring up Vagrant environment (http://www.vagrantup.com)
    ```
    $ cp traffic_control/traffic_portal/build/Vagrantfile ./
    $ vagrant up
```

3. ssh into vagrant environment
    ```
    $ vagrant ssh
    ```  

4. [OPTIONAL] Set RPM variables

  * BRANCH (enter version (i.e. 1.6.0) or leave default (master))
  * BUILD_NUMBER (defaults to number of git commits)

5. Build the RPM
    ```
    $ cd /vagrant/traffic_control/traffic_portal/build
    $ bash -x ./build_rpm.sh
    ```

6. Install the RPM
    ```
    $ cd /vagrant/traffic_control/traffic_portal/build
    $ sudo yum install -y traffic_portal-$BRANCH-$BUILD_NUMBER.x86_64.rpm
    ```

6. Configure Traffic Portal
    ```
    $ cd /etc/traffic_portal/conf
    $ sudo cp config-template.js config.js
    $ sudo vi config.js
    ```

7. Start Traffic Portal
    ```
    $ sudo service traffic_portal start
    ```

#### Notes

    - This is known to work with CentOS 6.7 as the Vagrant environment
    - /etc/traffic_portal/conf/config.js is consumed by /opt/traffic_portal/server/server.js on server startup
