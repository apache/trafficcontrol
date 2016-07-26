cd /
git clone ${GITREPO}
cd /traffic_control/traffic_monitor
git checkout ${BRANCH}
./build/build_rpm.sh
