# Traffic Portal Installation

### 1. Build w/ Docker

* Download Traffic Control repo

    ```
    $ git clone https://github.com/apache/incubator-trafficcontrol.git
    ```

* Build the RPM

    ```
    $ cd traffic_control/build
    $ ./docker-build.sh -r https://github.com/apache/incubator-trafficcontrol.git -b master traffic_portal
    ```

### 2. Install

* Install the Node.js JavaScript runtime

    ```
    $ curl --silent --location https://rpm.nodesource.com/setup_6.x | sudo bash -
    $ sudo yum install -y nodejs
    ```

* Install the Traffic Portal RPM

    ```
    $ sudo yum install -y traffic_portal-[version]-[commits].[sha].x86_64.rpm
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

* Navigate to Traffic Portal

    ```
    $ http://localhost:8080
    ```

#### Notes

    - Traffic Portal consumes the Traffic Ops API, therefore, an instance of Traffic Ops must be running.
    - This is known to work with CentOS 6.7 and Centos 7 as the host environment.
