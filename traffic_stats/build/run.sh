cd /
git clone ${GITREPO}
cd /traffic_control/traffic_stats
git checkout ${BRANCH}
./build/build_rpm.sh
