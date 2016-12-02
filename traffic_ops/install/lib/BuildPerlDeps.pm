#!/usr/bin/perl
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

use lib qw(/opt/traffic_ops/install/lib /opt/traffic_ops/lib/perl5 /opt/traffic_ops/app/lib);

package BuildPerlDeps;

use InstallUtils qw{ :all };

use base qw{ Exporter };
our @EXPORT_OK = qw{ build };
our %EXPORT_TAGS = ( all => \@EXPORT_OK );

sub build {
    my $opt_i = shift;

    my @dependencies = ( "expat-devel", "mod_ssl", "mkisofs", "libpcap", "libpcap-devel", "libcurl", "libcurl-devel", "mysql-server", "mysql-devel", "openssl", "openssl-devel", "cpan", "gcc", "make", "pkgconfig", "automake", "autoconf", "libtool", "gettext", "libidn-devel" );

    my $msg = << 'EOF';

This script will build and package the required Traffic Ops perl modules.
In order to complete this operation, Development tools such as the gcc
compiler will be installed on this machine.

EOF

    $ENV{PERL_MM_USE_DEFAULT}    = 1;
    $ENV{PERL_MM_NONINTERACTIVE} = 1;
    $ENV{AUTOMATED_TESTING}      = 1;

    my $result;

    if ( $ENV{USER} ne "root" ) {
        errorOut("You must run this script as the root user");
    }

    logger( $msg, "info" );

    chdir("/opt/traffic_ops/app");

    if ( defined $opt_i && $opt_i == 1 ) {
        if ( !-x "/usr/bin/yum" ) {
            errorOut("You must install 'yum'");
        }

        logger( "Installing dependencies", "info" );
        $result = execCommand( "/usr/bin/yum", "install", @dependencies );
        if ( $result != 0 ) {
            errorOut("Dependency installation failed, look through the output and correct the problem");
        }
        logger( "Building perl modules", "info" );

        $result = execCommand( "/opt/traffic_ops/install/bin/cpan.sh", "pi_custom_log=" . $::cpanLogFile, "/opt/traffic_ops/install/bin/yaml.txt" );
        if ( $result != 0 ) {
            errorOut("Failed to install YAML, look through the output and correct the problem");
        }

        $result = execCommand( "/opt/traffic_ops/install/bin/cpan.sh", "pi_custom_log=" . $::cpanLogFile, "/opt/traffic_ops/install/bin/carton.txt" );
        if ( $result != 0 ) {
            errorOut("Failed to install Carton, look through the output and correct the problem");
        }
    }

    $result = execCommand( "/usr/local/bin/carton", "install", "--deployment", "--cached" );
    if ( $result != 0 ) {
        errorOut("Failure to build required perl modules, check the output and correct the problem");
    }

    if ( !-s "/opt/traffic_ops/lib/perl5" ) {
        logger( "Linking perl libraries...", "info" );
        if ( !-d "/opt/traffic_ops/lib" ) {
            mkdir("/opt/traffic_ops/lib");
        }
        symlink( "/opt/traffic_ops/app/local/lib/perl5", "/opt/traffic_ops/lib/perl5" );
        execCommand( "/bin/chown", "-R", "trafops:trafops", "/opt/traffic_ops/lib" );
    }
    logger( "Installing perl scripts", "info" );
    chdir("/opt/traffic_ops/app/local/bin");
    my $rc = execCommand( "/bin/cp", "-R", ".", "/opt/traffic_ops/app/bin" );
    if ( $rc != 0 ) {
        logger( "Failed to copy perl scripts to /opt/traffic_ops/app/bin", "error" );
    }

    return 0;
}
