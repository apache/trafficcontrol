#!/usr/bin/perl
#
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

use strict;
use warnings;
use Getopt::Long;
use feature qw(switch);

my $dispersion = undef;
my $retries = 5;
my $wait_for_parents = 1;
my $login_dispersion = 0;
my $reval_wait_time = 60;
my $reval_in_use = 0;
my $rev_proxy_disable = 0;
my $skip_os_check = 0;
my $override_hostname_short = '';
my $to_timeout_ms = 30000;
my $syncds_updates_ipallow = 0;
my $traffic_ops_insecure = 0;
my $via_string_release = 0;
my $dns_local_bind = 0;
my $disable_parent_config_comments = 0;

print "ERROR traffic_ops_ort.pl is deprecated and will be removed in the next major version! Please upgrade to t3c\n";

GetOptions( "dispersion=i"       => \$dispersion, # dispersion (in seconds)
            "retries=i"          => \$retries,
            "wait_for_parents=i" => \$wait_for_parents,
            "login_dispersion=i" => \$login_dispersion,
            "rev_proxy_disable=i" => \$rev_proxy_disable,
            "skip_os_check=i" => \$skip_os_check,
            "override_hostname_short=s" => \$override_hostname_short,
            "to_timeout_ms=i" => \$to_timeout_ms,
            "syncds_updates_ipallow=i" => \$syncds_updates_ipallow,
            "traffic_ops_insecure=i" => \$traffic_ops_insecure,
            "via_string_release=i" => \$via_string_release,
            "dns_local_bind=i" => \$dns_local_bind,
            "disable_parent_config_comments=i" => \$disable_parent_config_comments,
          );

my $cmd = 't3c apply';

if ( defined $dispersion ) {
	my $sleeptime = rand($dispersion);
	print "ERROR t3c no longer has a dispersion feature, Please upgrade to t3c, and use shell commands to randomly sleep if necessary. Sleeping for rand($dispersion)=$sleeptime\n";
	sleep($sleeptime);
}
if ( defined $retries ) {
	$cmd .= ' --num-retries=' . $retries;
}
if ( defined  $wait_for_parents && $wait_for_parents == 0 ) {
	$cmd .= ' --wait-for-parents=false';
}
if ( defined $login_dispersion ) {
	my $sleeptime = rand($login_dispersion);
	print "ERROR t3c no longer has any dispersion feature, Please upgrade to t3c, and use shell commands to randomly sleep if necessary. Sleeping for rand($login_dispersion)=$sleeptime\n";
	sleep($sleeptime);
}
if ( defined $rev_proxy_disable && $rev_proxy_disable == 1 ) {
	$cmd .= ' --rev-proxy-disable=true';
}
if ( defined $skip_os_check ) {
	$cmd .= ' --skip-os-check=' . $skip_os_check;
}
if ( defined $override_hostname_short ) {
	$cmd .= ' --cache-host-name=' . $override_hostname_short;
} else {
	$cmd .= ' --cache-host-name=' . `hostname -s`;
}
if ( defined $to_timeout_ms ) {
	$cmd .= ' --traffic-ops-timeout-milliseconds=' . $to_timeout_ms;
}
if ( $syncds_updates_ipallow == 1 ) {
	$cmd .= ' --syncds-updates-ipallow=true';
}
if ( defined $traffic_ops_insecure ) {
	$cmd .= ' --traffic-ops-insecure=' . $traffic_ops_insecure;
}
if ( $via_string_release != 1 ) {
	$cmd .= ' --omit-via-string-release=true';
}
if ( $dns_local_bind == 1 ) {
	$cmd .= ' --dns-local-bind=true';
}
if ( $disable_parent_config_comments == 1 ) {
	$cmd .= ' --disable-parent-config-comments=true';
}

my $mode = $ARGV[0];
if ( defined $mode ) {
	$cmd .= ' --run-mode=' . $mode;
}

if ( defined( $ARGV[1] ) ) {
	my $log_level = "";
        my $ort_log_level = uc($ARGV[1]);
	given ( $ort_log_level ) {
                when ("ALL")   { $log_level = "-vv"; }
                when ("TRACE") { $log_level = "-vv"; }
                when ("DEBUG") { $log_level = "-vv"; }
                when ("INFO")  { $log_level = "-v"; }
                when ("WARN")  { $log_level = ""; }
                when ("ERROR") { $log_level = ""; }
                when ("FATAL") { $log_level = ""; }
                when ("NONE")  { $log_level = ""; }
	}
	$cmd .= ' ' . $log_level;
}

if ( defined( $ARGV[2] ) ) {
	my $to_url = $ARGV[2];
	$to_url =~ s/\/*$//g;
	$cmd .= ' --traffic-ops-url=' . $to_url;
}
if ( defined( $ARGV[3] ) ) {
	my ( $to_user, $to_pass ) = split( /:/, $ARGV[3] );
	$cmd .= ' --traffic-ops-user=' . $to_user;
	$cmd .= ' --traffic-ops-password=' . $to_pass;
}
else {
	&usage();
	exit 1;
}

sub usage {
	print "====-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-====\n";
	print "Usage: ./traffic_ops_ort.pl <Mode> <Log_Level> <Traffic_Ops_URL> <Traffic_Ops_Login> [optional flags]\n";
	print "Usage: ./traffic_ops_ort.pl <Mode> <Log_Level> <Traffic_Ops_URL> <Traffic_Ops_Login> [optional flags]\n";
	print "\t<Mode> = interactive - asks questions during config process.\n";
	print "\t<Mode> = report - prints config differences and exits.\n";
	print "\t<Mode> = badass - attempts to fix all config differences that it can.\n";
	print "\t<Mode> = syncds - syncs delivery services with what is configured in Traffic Ops.\n";
	print "\t<Mode> = revalidate - checks for updated revalidations in Traffic Ops and applies them.  Requires Traffic Ops 2.1.\n";
	print "\n";
	print "\t<Log_Level> => ALL, TRACE, DEBUG, INFO, WARN, ERROR, FATAL, NONE\n";
	print "\n";
	print "\t<Traffic_Ops_URL> = URL to Traffic Ops host. Example: https://trafficops.company.net\n";
	print "\n";
	print "\t<Traffic_Ops_Login> => Example: 'username:password' \n";
	print "\n\t[optional flags]:\n";
	print "\t   dispersion=<time>              => wait a random number between 0 and <time> before starting. Default = 300.\n";
	print "\t   login_dispersion=<time>        => wait a random number between 0 and <time> before login. Default = 0.\n";
	print "\t   retries=<number>               => retry connection to Traffic Ops URL <number> times. Default = 3.\n";
	print "\t   wait_for_parents=<0|1>         => do not update if parent_pending = 1 in the update json. Default = 1, wait for parents.\n";
	print "\t   rev_proxy_disable=<0|1>        => bypass the reverse proxy even if one has been configured Default = 0.\n";
	print "\t   skip_os_check=<0|1>            => bypass the check for a supported CentOS version. Default = 0.\n";
	print "\t   override_hostname_short=<text> => override the short hostname of the OS for config generation. Default = ''.\n";
	print "\t   to_timeout_ms=<time>           => the Traffic Ops request timeout in milliseconds. Default = 30000 (30 seconds).\n";
	print "\t   syncds_updates_ipallow=<0|1>   => Update ip_allow.config in syncds mode, which may trigger an ATS bug blocking random addresses on load! Default = 0, only update on badass and restart.\n";
	print "\t   traffic_ops_insecure=<0|1>     => Turns off certificate checking when connecting to Traffic Ops.\n";
	print "\t   via_string_release=<0|1>       => change the ATS via string to be the rpm release instead of the actual ATS version number\n";
	print "\t   dns_local_bind=<0|1>           => set the server service addresses to the ATS config dns local bind address\n";
	print "\t   disable_parent_config_comments=<0|1>     => do not write line comments to the parent.config file\n";
	print "====-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-====\n";
	exit 1;
}

exec ("$cmd");
