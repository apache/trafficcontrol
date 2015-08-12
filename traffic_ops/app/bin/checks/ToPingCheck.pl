#!/usr/bin/perl
#
# Copyright 2015 Comcast Cable Communications Management, LLC
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
#

# Plugin for the "ping" and "MTU" check.
#
# example cron entry
# 0 * * * * root /opt/traffic_ops/app/bin/checks/ToPingCheck.pl -c "{\"base_url\": \"https://localhost\", \"check_name\": \"10G\", \"name\": \"10G_PING\", \"select\": \"ipAddress\"}"
# example cron entry with select array
# 0 * * * * root /opt/traffic_ops/app/bin/checks/ToPingCheck.pl -c "{\"base_url\": \"https://localhost\", \"check_name\": \"10G\", \"name\": \"10G_PING\", \"select\": [\"hostName\",\"domainName\"]}"
# example cron entry with syslog
# 0 * * * * root /opt/traffic_ops/app/bin/checks/ToPingCheck.pl -c "{\"base_url\": \"https://localhost\", \"check_name\": \"10G\", \"name\": \"10G_PING\", \"select\": \"ipAddress\", \"syslog_facility\": \"local0\"}"
# example cron entry for MTU
# 0 0 * * * root /opt/traffic_ops/app/bin/checks/ToPingCheck.pl -c "{\"base_url\": \"https://localhost\", \"check_name\": \"MTU\", \"name\": \"Max Trans Unit\", \"select\": \"ipAddress\", \"syslog_facility\": \"local0\"}"
# 0 0 * * * root /opt/traffic_ops/app/bin/checks/ToPingCheck.pl -c "{\"base_url\": \"https://localhost\", \"check_name\": \"MTU\", \"name\": \"Max Trans Unit\", \"select\": \"ip6Address\", \"syslog_facility\": \"local0\"}"

use strict;
use warnings;

$|++;

use Data::Dumper;
use Getopt::Std;
use Log::Log4perl qw(:easy);
use JSON;
use Extensions::Helper;
use Sys::Syslog qw(:standard :macros);

my $VERSION = "0.02";

my %args = ();
getopts( "l:c:", \%args );

Log::Log4perl->easy_init($ERROR);
if ( defined( $args{l} ) ) {
	if    ( $args{l} == 1 ) { Log::Log4perl->easy_init($INFO); }
	elsif ( $args{l} == 2 ) { Log::Log4perl->easy_init($DEBUG); }
	elsif ( $args{l} == 3 ) { Log::Log4perl->easy_init($TRACE); }
	elsif ( $args{l} > 3 )  { Log::Log4perl->easy_init($TRACE); }
	else                    { Log::Log4perl->easy_init($INFO); }
}

# For syslog messages
setlogmask(LOG_UPTO(LOG_INFO));

DEBUG( "Including DEBUG messages in output. Config is \'" . $args{c} . "\'" );
TRACE( "Including TRACE messages in output. Config is \'" . $args{c} . "\'" );

if ( !defined( $args{c} ) ) {
	&help();
	exit(1);
}

my $jconf = undef;
eval { $jconf = decode_json( $args{c} ) };
if ($@) {
	ERROR("Bad json config: $@");
	exit(1);
}

my $sslg = undef;
if (defined($jconf->{syslog_facility})) {
   openlog ('ToChecks', '', $jconf->{syslog_facility});
   $sslg = 1;
}

TRACE Dumper($jconf);
my $b_url = $jconf->{base_url};
Extensions::Helper->import();
my $ext = Extensions::Helper->new( { base_url => $b_url, token => '91504CE6-8E4A-46B2-9F9F-FE7C15228498' } );

my $jdataserver    = $ext->get(Extensions::Helper::SERVERLIST_PATH);
my $select         = $jconf->{select};
my $check_name     = &trim($jconf->{check_name});
my $check_lng_name = &trim($jconf->{name});

foreach my $server ( @{$jdataserver} ) {
	if ( $server->{type} eq 'EDGE' || $server->{type} eq 'MID' ) {
      my $srv_nm = $server->{hostName}.".".$server->{domainName};
      my $srv_status = &trim($server->{status});
		my $ip = undef;
      my $pingable = undef;
      my $size = &trim($server->{interfaceMtu});
      $size = $size - 28;

      # select in the jconf is mandatory. TODO should probably error if not there
      if ( ref($select) eq 'ARRAY' ) {
         DEBUG "select is an array";
         $select->[0] = &trim($select->[0]);
         $select->[1] = &trim($select->[1]);
         $ip = $server->{ $select->[0] } . "." . $server->{ $select->[1] };
      }
      else {
         DEBUG "select is not an array";
         $select = &trim($select);
         $ip = &trim($server->{$select});
         DEBUG "ip: ".$ip;
      }
      if (!defined($ip) || ($ip eq '')) {
         next;
      }

      if ($check_name =~ m/^MTU$/) {
         $pingable = &ping_check($ip, $size);
      } else {
         $pingable = &ping_check( $ip, 30 );
      }
      if ($pingable && $sslg) {
         $ip =~ s/\/\d+$// if ( $ip =~ /:/ );
         my @tmp = ($srv_nm,$check_name,$check_lng_name,'OK',$srv_status,$ip);
         syslog(LOG_INFO, "hostname=%s check=%s name=\"%s\" result=%s status=%s target=%s", @tmp);
      } elsif ($sslg) {
         $ip =~ s/\/\d+$// if ( $ip =~ /:/ );
         my @tmp = ($srv_nm,$check_name,$check_lng_name,'FAIL',$srv_status,$ip);
         syslog(LOG_ERR, "hostname=%s check=%s name=\"%s\" result=%s status=%s target=%s", @tmp);
      }
		DEBUG $check_name . " >> " . $server->{hostName} . ": " . $select . " = " . $ip . " ---> " . $pingable . "\n";
		$ext->post_result( $server->{id}, $check_name, $pingable );
	}
}

closelog();

sub help {
	print "The -c argument is mandatory\n";
}

sub ping_check {
	my $ping_target = shift;    # use address to bypass DNS and FQDN to check DNS
	my $size        = shift;

	if ( !defined($ping_target) ) {
		print "Nothing to ping!\n";
		return 0;
	}

	if ( !defined($size) ) {
		$size = 30;
	}

	TRACE "Ping checking " . $ping_target." with: ".$size;

	my $cmd;
	if ( $ping_target =~ /:/ ) {
		$ping_target =~ s/\/\d+$//;
		$cmd = '/bin/ping6 -M do -s ' . $size . ' -c 2 ' . $ping_target . ' 2>&1 > /dev/null';
	}
	else {
		$cmd = '/bin/ping -M do -s ' . $size . ' -c 2 ' . $ping_target . ' 2>&1 > /dev/null';
	}

	system($cmd);
	if ( $? != 0 ) {
		ERROR $ping_target . " is NOT Pingable (with " . $size . " packet size)";
		return 0;
	}
	return 1;
}

sub ltrim { my $s = shift; $s =~ s/^\s+//;       return $s };
sub rtrim { my $s = shift; $s =~ s/\s+$//;       return $s };
sub  trim { my $s = shift; $s =~ s/^\s+|\s+$//g; return $s };
