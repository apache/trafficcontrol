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

package BuildPerlDeps;

use strict;
use warnings;

use InstallUtils qw{ :all };

use base qw{ Exporter };
our @EXPORT_OK = qw{ build };
our %EXPORT_TAGS = ( all => \@EXPORT_OK );

sub build {
    my $opt_i       = shift;
    my $cpanLogFile = shift;

    my @dependencies = ( "expat-devel", "mod_ssl", "mkisofs", "libpcap", "libpcap-devel", "libcurl", "libcurl-devel", "openssl", "openssl-devel", "cpan", "gcc", "make", "pkgconfig", "automake", "autoconf", "libtool", "gettext", "libidn-devel" );

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

    InstallUtils::logger( $msg, "info" );

    chdir("/opt/traffic_ops/app");

    if ( defined $opt_i && $opt_i == 1 ) {
        if ( !-x "/usr/bin/yum" ) {
            errorOut("You must install 'yum'");
        }

        InstallUtils::logger( "Installing dependencies", "info" );
        $result = InstallUtils::execCommand( "/usr/bin/yum", "-y", "install", @dependencies );
        if ( $result != 0 ) {
            errorOut("Dependency installation failed, look through the output and correct the problem");
        }
        InstallUtils::logger( "Building perl modules", "info" );

        $result = InstallUtils::execCommand( "/usr/bin/cpan", "pi_custom_log=" . $cpanLogFile, "-if", "YAML" );
        if ( $result != 0 ) {
            errorOut("Failed to install YAML, look through the output and correct the problem");
        }

        $result = InstallUtils::execCommand( "/usr/bin/cpan", "pi_custom_log=" . $cpanLogFile, "-if", "MIYAGAWA/Carton-v1.0.15.tar.gz" );
        if ( $result != 0 ) {
            errorOut("Failed to install Carton, look through the output and correct the problem");
        }
    }

    $result = InstallUtils::execCommand( "/usr/local/bin/carton", "install" );
    if ( $result != 0 ) {
        errorOut("Failure to build required perl modules, check the output and correct the problem");
    }

    return 0;
}

1;
