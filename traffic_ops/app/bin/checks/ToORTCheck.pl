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
# ORT check extension. Checks how many errors there are running the "ort" script on the cache
#
# example cron entry
# 40 * * * * ssh_key_user /opt/traffic_ops/app/bin/checks/ToORTCheck.pl -c "{\"base_url\": \"https://localhost\", \"check_name\": \"ORT\", \"name\": \"Operational Readiness Test\", \"ssh_user\": \"<ssh_key_user>\", \"to_user\": \"<some_user>\", \"to_pass\": \"<some_pass>\"}"
# example cron entry with syslog
# 40 * * * * ssh_key_user /opt/traffic_ops/app/bin/checks/ToORTCheck.pl -c "{\"base_url\": \"https://localhost\", \"check_name\": \"ORT\", \"name\": \"Operational Readiness Test\", \"ssh_user\": \"<ssh_key_user>\", \"to_user\": \"<some_user>\", \"to_pass\": \"<some_pass>\", \"syslog_facility\": \"local0\"}"

# TODO: use tokens instead of username and password.

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
getopts( "hl:c:", \%args );

if ($args{h}) {
   &help();
   exit();
}

if ( !defined( $args{c} ) ) {
   ERROR "-c not defined";
	&help();
	exit(1);
}

my $jconf = undef;
eval { $jconf = decode_json( $args{c} ) };
if ($@) {
	ERROR("Bad json config: $@");
   &help();
	exit(1);
}

my $check_name  = $jconf->{check_name};
my $chck_lng_nm = $jconf->{name};
my $to_user     = $jconf->{to_user};
my $to_pass     = $jconf->{to_pass};

if ( $check_name ne "ORT" ) {
	ERROR "This Check Extension is exclusively for the ORT (Operational Readiness Test) check.";
   &help();
	exit(4);
}

my $sslg = undef;
if (defined($jconf->{syslog_facility})) {
   openlog ('ToChecks', '', $jconf->{syslog_facility});
   $sslg = 1;
}

Log::Log4perl->easy_init($ERROR);
if ( defined( $args{l} ) ) {
	if    ( $args{l} == 1 ) { Log::Log4perl->easy_init($INFO); }
	elsif ( $args{l} == 2 ) { Log::Log4perl->easy_init($DEBUG); }
	elsif ( $args{l} >= 3 ) { Log::Log4perl->easy_init($TRACE); }
	else                    { Log::Log4perl->easy_init($INFO); }
}

DEBUG( "Including DEBUG messages in output. Config is \'" . $args{c} . "\'" );
TRACE( "Including TRACE messages in output. Config is \'" . $args{c} . "\'" );

my $b_url = $jconf->{base_url};
Extensions::Helper->import();
my $ext = Extensions::Helper->new( { base_url => $b_url, token => '91504CE6-8E4A-46B2-9F9F-FE7C15228498' } );

my $jdataserver = $ext->get(Extensions::Helper::SERVERLIST_PATH);

my $glbl_prms = $ext->get( '/api/1.1/parameters/profile/global.json' );
my $to_url;
foreach my $p (@{$glbl_prms}) {
   if ($p->{name} eq 'tm.url') {
      $to_url = $p->{value};
   }
}

foreach my $server ( @{$jdataserver} ) {
	if ( $server->{type} eq 'EDGE' || $server->{type} eq 'MID' ) {
		&ort_check( $server->{ipAddress}, $server->{hostName}, $server->{id} );
	}
}

sub ort_check() {
	my $ipaddr    = shift;
	my $host_name = shift;
	my $host_id   = shift;

	my $cmd = "/usr/bin/sudo /opt/ort/traffic_ops_ort.pl report WARN ".$to_url." '".$to_user.":".$to_pass."'";
	#$cmd = "ssh -t -o \"StrictHostKeyChecking no\" -i ~$username/.ssh/."$key." -l " . $username . " " . $ipaddr . " " . $cmd;
	$cmd = "ssh -t -o \"StrictHostKeyChecking no\" ".$ipaddr." ".$cmd." 2>&1";
	TRACE $host_name . " running " . $cmd;
	my $out = `$cmd`;

   my $results = "/var/www/html/ort/" . $host_id;

	if ( !-d $results ) {
		mkdir($results);
	}

   # TODO integrate this file into Traffic Ops
   my $filename = $results."/results";
   open(my $fh, '>', $filename);
   if (!$fh) {
      ERROR "Could not open file '$filename' $!";
   } else {
      print $fh $out;
   }
   close $fh;

	my @lines = split( /\n/, $out );
	my $errors = grep /^ERROR/, @lines;
   my $fatals = grep /^FATAL/, @lines;
   my $warns = grep /^WARN/, @lines;
   if ($out =~ m/Connection timed out/) {
      ERROR "Could not connect to server.";
      $errors = -1;
   }
   if ($sslg) {
      if ($errors == -1) {
         my @tmp = ($host_name,$check_name,$chck_lng_nm,'FAIL');
         syslog(LOG_ERR, "hostname=%s check=%s name=\"%s\" result=%s msg=\"Could not SSH to server\"", @tmp);
      } else {
         my @tmp = ($host_name,$check_name,$chck_lng_nm,'OK',$errors,$fatals,$warns);
         syslog(LOG_ERR, "hostname=%s check=%s name=\"%s\" result=%s errors=%s fatals=%s warnings=%s", @tmp);
      }
   }
	TRACE $host_name . " score: " . $errors;
	$ext->post_result( $host_id, $check_name, $errors );
}

sub help() {
   print "ToORTCheck.pl -c \"{\\\"base_url\\\": \\\"https://localhost\\\", \\\"check_name\\\": \\\"ORT\\\", \\\"name\\\": \\\"Operational Readiness Test\\\", \\\"syslog_facility\\\": \\\"local0\\\"}\"\n";
   print "\n";
   print "-c   json formatted list of variables\n";
   print "     base_url: required\n";
   print "        URL of the Traffic Ops server.\n";
   print "     check_name: required\n";
   print "        The name of this check. Don't ask.\n";
   print "     to_user: required\n";
   print "        Traffic Ops user.\n";
   print "     to_pass: required\n";
   print "        Password for the Traffic Ops user\n";
   print "     name: optional\n";
   print "        The long name of this check. used in conjuction with syslog_facility.\n";
   print "     syslog_facility: optional\n";
   print "        The syslog facility to send messages. Requires the \"name\" option to\n";
   print "        be set.\n";
   print "-h   Print this message\n";
   print "-l   Debug level\n";
   print "================================================================================\n";
   # the above line of equal signs is 80 columns
   print "Note: The user running this script must have an authorized public key on the\n";
   print "      cache servers and must also have sudoers access.\n";
   print "\n";
}
