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
# RPM spec file for Traffic Ops (tm).
#

%define TRAFFIC_OPS_USER trafops
%define TRAFFIC_OPS_GROUP trafops
%define TRAFFIC_OPS_LOG_DIR /var/log/traffic_ops

Summary:          Traffic Ops UI
Name:             traffic_ops
Version:          %{traffic_control_version}
Release:          %{build_number}
License:          Apache License, Version 2.0
Group:            Base System/System Tools
Prefix:           /opt/traffic_ops
Source:           %{_sourcedir}/traffic_ops-%{version}.tgz
URL:	          https://github.com/Comcast/traffic_control/
Vendor:	          Comcast
Packager:         daniel_kirkwood at Cable dot Comcast dot com
AutoReqProv:      no
Requires:         cpanminus, expat-devel, gcc-c++, libcurl, libpcap-devel, mkisofs, tar
Requires:         openssl-devel, perl, perl-core, perl-DBD-Pg, perl-DBI, perl-Digest-SHA1
Requires:	  libidn-devel, libcurl-devel
Requires:         perl-JSON, perl-libwww-perl, perl-Test-CPAN-Meta, perl-WWW-Curl
Requires(pre):    /usr/sbin/useradd, /usr/bin/getent
Requires(postun): /usr/sbin/userdel

%define PACKAGEDIR %{prefix}

%description
Installs Traffic Ops.

Built: %(date) by %{getenv: USER}

%prep

%setup

%build
    # update version referenced in the source
    perl -pi.bak -e 's/__VERSION__/%{version}-%{release}/' app/lib/UI/Utils.pm
    # compile go executables used during postinstall
    # suppress strip of go execs
    %define debug_package %{nil}

    export GOPATH="$(pwd)/install/go"
    export GOBIN="$(pwd)/install/bin"

    echo "Compiling go executables"
    for d in install/go/src/comcast.com/*; do
	(cd "$d" && go get -ldflags "-B 0x%{commit}" -v ) || \
	    { echo "Could not compile $d"; exit 1; }
    done

%install

    if [ -d $RPM_BUILD_ROOT ]; then
		%__rm -rf $RPM_BUILD_ROOT
    fi

    if [ ! -d $RPM_BUILD_ROOT/%{PACKAGEDIR} ]; then
		%__mkdir -p $RPM_BUILD_ROOT/%{PACKAGEDIR}
    fi

    %__cp -R $RPM_BUILD_DIR/traffic_ops-%{version}/* $RPM_BUILD_ROOT/%{PACKAGEDIR}

    if [ ! -d $RPM_BUILD_ROOT/%{PACKAGEDIR}/app/public/CRConfig-Snapshots ]; then
        %__mkdir -p $RPM_BUILD_ROOT/%{PACKAGEDIR}/app/public/CRConfig-Snapshots
    fi
    if [ ! -d $RPM_BUILD_ROOT/%{PACKAGEDIR}/app/public/routing ]; then
        %__mkdir -p $RPM_BUILD_ROOT/%{PACKAGEDIR}/app/public/routing
    fi

%pre
    /usr/bin/getent group %{TRAFFIC_OPS_GROUP} || /usr/sbin/groupadd -r %{TRAFFIC_OPS_GROUP}
    /usr/bin/getent passwd %{TRAFFIC_OPS_USER} || /usr/sbin/useradd -r -d %{PACKAGEDIR} -s /sbin/nologin %{TRAFFIC_OPS_USER} -g %{TRAFFIC_OPS_GROUP}
    if [ -d %{PACKAGEDIR}/app/conf ]; then
	  echo -e "\nBacking up config files.\n"
	  if [ -f /var/tmp/traffic_ops-backup.tar ]; then
		  %__rm /var/tmp/traffic_ops-backup.tar
	  fi
	  cd %{PACKAGEDIR} && tar cf /var/tmp/traffic_ops-backup.tar app/public/*Snapshots app/public/routing  app/conf app/db/dbconf.yml app/local app/cpanfile.snapshot
    fi

    # upgrade
    if [ "$1" == "2" ]; then
	service traffic_ops stop
    fi

%post

    %__cp %{PACKAGEDIR}/etc/init.d/traffic_ops /etc/init.d/traffic_ops
    %__cp %{PACKAGEDIR}/etc/cron.d/trafops_dnssec_refresh /etc/cron.d/trafops_dnssec_refresh
     %__cp %{PACKAGEDIR}/etc/logrotate.d/traffic_ops /etc/logrotate.d/traffic_ops
     %__cp %{PACKAGEDIR}/etc/logrotate.d/traffic_ops_access /etc/logrotate.d/traffic_ops_access
    %__chown root:root /etc/init.d/traffic_ops
    %__chown root:root /etc/cron.d/trafops_dnssec_refresh
    %__chown root:root /etc/logrotate.d/traffic_ops
    %__chown root:root /etc/logrotate.d/traffic_ops_access
    %__chmod +x /etc/init.d/traffic_ops
    %__chmod +x %{PACKAGEDIR}/install/bin/*
    /sbin/chkconfig --add traffic_ops 
	
    %__mkdir -p %{TRAFFIC_OPS_LOG_DIR}

    if [ -f /var/tmp/traffic_ops-backup.tar ]; then
    	echo -e "\nRestoring config files.\n"
		cd %{PACKAGEDIR} && tar xf /var/tmp/traffic_ops-backup.tar
    fi

    # install
    if [ "$1" = "1" ]; then
      # see postinstall, the .reconfigure file triggers init().
      /bin/touch %{PACKAGEDIR}/.reconfigure
    	echo -e "\nRun /opt/traffic_ops/install/bin/postinstall from the root home directory to complete the install.\n"
    fi

    # upgrade
    if [ "$1" == "2" ]; then
		    /opt/traffic_ops/install/bin/migratedb
        echo -e "\nUpgrade complete.\n\n"
    	 echo -e "\nRun /opt/traffic_ops/install/bin/postinstall from the root home directory to complete the update.\n"
        echo -e "To start Traffic Ops:  service traffic_ops start\n";
        echo -e "To stop Traffic Ops:   service traffic_ops stop\n\n";
    fi
    /bin/chown -R %{TRAFFIC_OPS_USER}:%{TRAFFIC_OPS_GROUP} %{PACKAGEDIR}
    /bin/chown -R %{TRAFFIC_OPS_USER}:%{TRAFFIC_OPS_GROUP} %{TRAFFIC_OPS_LOG_DIR}

%preun

if [ "$1" = "0" ]; then
    # stop service before starting the uninstall
    service traffic_ops stop
fi

%postun

if [ "$1" = "0" ]; then
	# this is an uninstall
	%__rm -rf %{PACKAGEDIR}
	%__rm /etc/init.d/traffic_ops
    /usr/bin/getent passwd %{TRAFFIC_OPS_USER} || /usr/sbin/userdel %{TRAFFIC_OPS_USER} 
    /usr/bin/getent group %{TRAFFIC_OPS_GROUP} || /usr/sbin/groupdel %{TRAFFIC_OPS_GROUP}
fi

%files
%defattr(644,root,root,755)
%attr(755,root,root) %{PACKAGEDIR}/app/bin/*
%attr(755,root,root) %{PACKAGEDIR}/app/script/*
%attr(755,root,root) %{PACKAGEDIR}/app/db/*.pl
%attr(755,root,root) %{PACKAGEDIR}/app/db/*.sh
%config(noreplace)/opt/traffic_ops/app/conf/*
%{PACKAGEDIR}/app/cpanfile
%{PACKAGEDIR}/app/db
%{PACKAGEDIR}/app/lib
%{PACKAGEDIR}/app/public
%{PACKAGEDIR}/app/templates
%{PACKAGEDIR}/install
%exclude %{PACKAGEDIR}/install/go
%{PACKAGEDIR}/etc
%doc %{PACKAGEDIR}/doc
