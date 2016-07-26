cd /
git clone ${GITREPO}
cd /traffic_control/traffic_ops
git checkout ${BRANCH}
./build/build_rpm.sh
