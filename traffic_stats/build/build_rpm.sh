GIT_SHORT_REVISION=`git rev-parse --short HEAD`; export GIT_SHORT_REVISION
VERSION="1.2.1"; export VERSION
SRC_BASE_DIR="/tmp/.traffic_stats.$$"; export SRC_BASE_DIR
SRC_NAME="traffic_stats-${VERSION}"; export SRC_NAME
SRC_DIR="${SRC_BASE_DIR}/${SRC_NAME}"; export SRC_DIR
export WORKSPACE=/vol1/jenkins/jobs
GOPATH=$HOME/go; export GOPATH

#. ~/.bash_profile

##build traffic_ops client
/usr/bin/rsync -av --delete $WORKSPACE/traffic_stats/workspace/traffic_ops/client/ $GOPATH/src/github.com/comcast/traffic_control/traffic_ops/client/
cd $GOPATH/src/github.com/comcast/traffic_control/traffic_ops/client/
/usr/local/bin/go install

##build influxdb client
#cd $GOPATH/src/github.com/influxdb/influxdb/client/
#git pull
#/usr/local/bin/go install


#traffic_stats
/usr/bin/rsync -av --delete $WORKSPACE/traffic_stats/workspace/traffic_stats/ $GOPATH/src/github.com/comcast/traffic_control/traffic_stats/
cd $GOPATH/src/github.com/comcast/traffic_control/traffic_stats
/usr/local/bin/go get

sed -i -e "s/@VERSION@/$VERSION/g" traffic_stats.spec
sed -i -e "s/@RELEASE@/$GIT_SHORT_REVISION/g" traffic_stats.spec
#go get all
/usr/local/bin/go build traffic_stats.go


rm -rf $WORKSPACE/traffic_stats/SOURCES/traffic_stats*
rm -rf $WORKSPACE/traffic_stats/RPMS/x86_64/traffic_stats*

mkdir -p ${SRC_DIR}/opt/traffic_stats
mkdir -p ${SRC_DIR}/opt/traffic_stats/bin
mkdir -p ${SRC_DIR}/opt/traffic_stats/conf
mkdir -p ${SRC_DIR}/opt/traffic_stats/var/log
mkdir -p ${SRC_DIR}/etc/init.d
mkdir -p ${SRC_DIR}/etc/logrotate.d

cp $GOPATH/src/github.com/comcast/traffic_control/traffic_stats/traffic_stats ${SRC_DIR}/opt/traffic_stats/bin
cp $WORKSPACE/traffic_stats/workspace/traffic_stats/traffic_stats.cfg ${SRC_DIR}/opt/traffic_stats/conf/traffic_stats.cfg
cp $WORKSPACE/traffic_stats/workspace/traffic_stats/traffic_stats_seelog.xml ${SRC_DIR}/opt/traffic_stats/conf
cp $WORKSPACE/traffic_stats/workspace/traffic_stats/traffic_stats.init ${SRC_DIR}/etc/init.d/traffic_stats
cp $WORKSPACE/traffic_stats/workspace/traffic_stats/traffic_stats.logrotate ${SRC_DIR}/etc/logrotate.d/traffic_stats

cd $SRC_BASE_DIR
tar -zcvf ~/rpmbuild/SOURCES/${SRC_NAME}.tar.gz *
cd -
rm -rf $SRC_BASE_DIR


rpmbuild -bb traffic_stats.spec

