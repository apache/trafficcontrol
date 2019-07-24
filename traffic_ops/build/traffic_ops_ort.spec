#
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
#
# RPM spec file for Traffic Stats (tm).
#
%define debug_package %{nil}
Name:		traffic_ops_ort
Summary:	Installs ORT script for Traffic Control caches
Version:	%{traffic_control_version}
Release:	%{build_number}
License:	Apache License, Version 2.0
Group:		Applications/Communications
Source0:	traffic_ops_ort-%{version}.tgz
URL:		https://github.com/apache/trafficcontrol/
Vendor:		Apache Software Foundation
Packager:	daniel_kirkwood at Cable dot Comcast dot com
%{?el6:Requires: perl-JSON, perl-libwww-perl, perl-Crypt-SSLeay, perl-Digest-SHA}
%{?el7:Requires: perl-JSON, perl-libwww-perl, perl-Crypt-SSLeay, perl-LWP-Protocol-https, perl-Digest-SHA}


%description
Installs ORT script for Traffic Ops caches

%prep
tar xvf %{SOURCE0} -C $RPM_SOURCE_DIR


%build
export GOPATH=$(pwd)
# Create build area with proper gopath structure
mkdir -p src pkg bin || { echo "Could not create directories in $(pwd): $!"; exit 1; }

go_get_version() {
  local src=$1
  local version=$2
  (
   cd $src && \
   git checkout $version && \
   go get -v \
  )
}

# build all internal go dependencies (expects package being built as argument)
build_dependencies () {
    IFS=$'\n'
    array=($(go list -f '{{ join .Deps "\n" }}' | grep trafficcontrol | grep -v $1))
    echo "array: AA${array}AA";

    prefix=github.com/apache/trafficcontrol
    for (( i=0; i<${#array[@]}; i++ )); do
        curPkg=${array[i]};
        curPkgShort=${curPkg#$prefix};
        echo "checking $curPkg";
        godir=$GOPATH/src/$curPkg;
        if [ ! -d "$godir" ]; then
          ( echo "building $curPkg" && \
            mkdir -p "$godir" && \
            cd "$godir" && \
            cp -r "$TC_DIR$curPkgShort"/* . && \
            build_dependencies "$curPkgShort" && \
            go get -v && \
            echo "go building $curPkgShort at $(pwd)" && \
            go build \
          ) || { echo "Could not build go $curPkgShort at $(pwd): $!"; exit 1; };
        fi
     done
}

#build atstccfg binary
godir=src/github.com/apache/trafficcontrol/traffic_ops/ort/atstccfg
oldpwd=$(pwd)
( mkdir -p "$godir" && \
  cd "$godir" && \
  cp -r "$TC_DIR"/traffic_ops/ort/atstccfg/* . && \
  build_dependencies atstccfg  && \
  #with proper vendoring go get would be  unneeded.
  go get -d -v && \
  go build -ldflags "-X main.GitRevision=`git rev-parse HEAD` -X main.BuildTimestamp=`date +'%Y-%M-%dT%H:%M:%s'` -X main.Version=%{traffic_control_version}"
) || { echo "Could not build go program at $(pwd): $!"; exit 1; }


%install
mkdir -p ${RPM_BUILD_ROOT}/opt/ort
cp -p ${RPM_SOURCE_DIR}/traffic_ops_ort-%{version}/traffic_ops_ort.pl ${RPM_BUILD_ROOT}/opt/ort
cp -p ${RPM_SOURCE_DIR}/traffic_ops_ort-%{version}/supermicro_udev_mapper.pl ${RPM_BUILD_ROOT}/opt/ort

src=src/github.com/apache/trafficcontrol/traffic_ops/ort/atstccfg
cp -p "$src"/atstccfg ${RPM_BUILD_ROOT}/opt/ort

%clean
rm -rf ${RPM_BUILD_ROOT}

%post

%files
%attr(755, root, root)
/opt/ort/traffic_ops_ort.pl
/opt/ort/supermicro_udev_mapper.pl
/opt/ort/atstccfg

%changelog
