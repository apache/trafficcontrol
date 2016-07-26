cd /
git clone ${GITREPO}
cd /traffic_control/traffic_router
git checkout ${BRANCH}
./build/build_rpm.sh
