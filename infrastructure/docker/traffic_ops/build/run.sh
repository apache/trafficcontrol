cd /
git clone ${GITREPO}
cd /traffic_control
git checkout ${BRANCH}
cd /traffic_control/traffic_ops
./build/build_rpm.sh
