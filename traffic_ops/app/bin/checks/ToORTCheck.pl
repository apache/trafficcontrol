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
#
# ORT check extension. Checks how many errors there are running the "ort" script on the cache
#
# example cron entry
# 40 * * * * ssh_key_edge_user /opt/traffic_ops/app/bin/checks/ToORTCheck.pl -c "{\"base_url\": \"https://localhost\", \"check_name\": \"ORT\"}" >> /var/log/traffic_ops/extensionCheck.log 2>&1
#
# example cron entry with syslog
# 40 * * * * ssh_key_edge_user /opt/traffic_ops/app/bin/checks/ToORTCheck.pl -c "{\"base_url\": \"https://localhost\", \"check_name\": \"ORT\", \"name\": \"Operational Readiness Test\", \"syslog_facility\": \"local0\"}" > /dev/null 2>&1

use strict;
use warnings;

$|++;

use Data::Dumper;
use Getopt::Std;
use Log::Log4perl qw(:easy);
use JSON;
use Extensions::Helper;
use Sys::Syslog qw(:standard :macros);

my $VERSION = "0.03";

my %args = ();
getopts( "c:f:hl:q", \%args );

if ($args{h}) {
   &help();
   exit();
}

if ( !defined( $args{c} ) ) {
   ERROR "-c not defined";
   print "\n\n";
	&help();
	exit(1);
}

my $jconf = undef;
eval { $jconf = decode_json( $args{c} ) };
if ($@) {
	ERROR("Bad json config: $@");
   print "\n\n";
   &help();
	exit(1);
}

my $check_name  = $jconf->{check_name};
my $to_user     = $jconf->{to_user};
my $to_pass     = $jconf->{to_pass};

if ( $check_name ne "ORT" ) {
	ERROR "This Check Extension is exclusively for the ORT (Operational Readiness Test) check.";
   print "\n\n";
   &help();
	exit(4);
}

my $sslg = undef;
my $chck_lng_nm;
if (defined($jconf->{syslog_facility})) {
   $chck_lng_nm = $jconf->{name};
   openlog ('ToChecks', '', $jconf->{syslog_facility});
   $sslg = 1;
}

my $force = 0;
if (defined($args{f})) {
   $force = $args{f};
}

my $quiet;
if ($args{q}) {
   $quiet = 1;
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

my $glbl_prms = $ext->get( '/api/4.0/profiles/name/GLOBAL/parameters' );
my $to_url;
foreach my $p (@{$glbl_prms}) {
   if ($p->{name} eq 'tm.url') {
      $to_url = $p->{value};
   }
}

foreach my $server ( @{$jdataserver} ) {
	if ( $server->{type} =~ m/^EDGE/ || $server->{type} =~ m/^MID/ ) {
		&ort_check( $server->{ipAddress}, $server->{hostName}, $server->{id},
                  $server->{domainName}, $server->{status} );
	}
}

sub ort_check() {
	my $ipaddr    = shift;
	my $host_name = shift;
	my $host_id   = shift;
   my $domain    = shift;
   my $status    = shift;

   my $ort_version = "unknown";
   my $host_name_fqdn = $host_name.".".$domain;

	my $cmd = "/usr/bin/sudo /opt/ort/traffic_ops_ort.pl report WARN ".$to_url." '".$to_user.":".$to_pass."'";
	#$cmd = "ssh -t -o \"StrictHostKeyChecking no\" -i ~$username/.ssh/."$key." -l " . $username . " " . $ipaddr . " " . $cmd;
	$cmd = "ssh -t -o ConnectTimeout=5 -o \"UserKnownHostsFile /dev/null\" -o \"StrictHostKeyChecking no\" -o \"BatchMode yes\" ".$ipaddr." ".$cmd." 2>&1";
	TRACE $host_name . " running " . $cmd;
   my $out = '';
   my $msg = '';
   my ($fatals, $errors, $warns);
   if ($force == 0) {
      $out = `$cmd`;
      my $results = "/var/www/html/ort/" . $host_id;

      if ( !-d $results ) {
         mkdir($results);
      }

      # TODO integrate this file into Traffic Ops GUI
      my $filename = $results."/results";
      open(my $fh, '>', $filename);
      if (!$fh) {
         ERROR "Could not open file '$filename' $!";
      } else {
         print $fh $out;
      }
      close $fh;

      my @lines = split( /\n/, $out );
      $errors = grep /^ERROR/, @lines;
      $fatals = grep /^FATAL/, @lines;
      $warns = grep /^WARN/, @lines;
      foreach my $line (@lines) {
         if ($line =~ m/^Version/) {
            chomp($line);
            my @tmp = split(" ", $line);
            $ort_version = $tmp[-1];
         }
      }
      DEBUG "ORT version: ".$ort_version;
      if ($out =~ m/Connection timed out/) {
         ERROR "Could not connect to server.";
         $msg = "Could not connect to server.";
         $errors = -1;
      } elsif ($out =~ m/Permission denied/) {
         $msg = "Permission Denied";
         ERROR "Permission Denied";
         $errors = -1;
      }
   } elsif ($force == 1) {
      $msg = "Force: FAIL";
      $errors = -1;
   } elsif ($force == 2) {
      $msg = "Force: OK";
      $errors = 1;
      $fatals = 1;
      $warns = 1;
   }


   if ($sslg) {
      if ($errors < 0) {
         my @tmp = ($host_name_fqdn,$check_name,$chck_lng_nm,'FAIL',$status,$ort_version,$msg);
         syslog(LOG_ERR, "hostname=%s check=%s name=\"%s\" result=%s status=%s ort_version=%s msg=\"%s\"", @tmp);
      } elsif (($errors >= 0) || ($force == 2)) {
         my @tmp = ($host_name_fqdn,$check_name,$chck_lng_nm,'OK',$status,$ort_version,$errors,$fatals,$warns,$msg);
         syslog(LOG_ERR, "hostname=%s check=%s name=\"%s\" result=%s status=%s ort_version=%s errors=%s fatals=%s warnings=%s msg=\"%s\"", @tmp);
      }
   }
	TRACE $host_name . " score: " . $errors;
	$ext->post_result( $host_id, $check_name, $errors ) if (!$quiet);
}

sub ltrim { my $s = shift; $s =~ s/^\s+//;       return $s };
sub rtrim { my $s = shift; $s =~ s/\s+$//;       return $s };
sub  trim { my $s = shift; $s =~ s/^\s+|\s+$//g; return $s };

sub help() {
   print "ToORTCheck.pl -c \"{\\\"base_url\\\": \\\"https://localhost\\\", \\\"check_name\\\": \\\"ORT\\\"[, \\\"name\\\": \\\"Operational Readiness Test\\\", \\\"syslog_facility\\\": \\\"local0\\\"]}\" [-f <1-2>] [-l <1-3>]\n";
   print "\n";
   print "-c   json formatted list of variables\n";
   print "     base_url: required\n";
   print "        URL of the Traffic Ops server.\n";
   print "     check_name: required\n";
   print "        The name of this check.\n";
   print "     name: optional\n";
   print "        The long name of this check. used in conjuction with syslog_facility.\n";
   print "     syslog_facility: optional\n";
   print "        The syslog facility to send messages. Requires the \"name\" option to\n";
   print "        be set.\n";
   print "-f   Force a FAIL or OK message\n";
   print "        1: FAIL\n";
   print "        2: OK\n";
   print "-h   Print this message\n";
   print "-l   Debug level\n";
   print "-q   Don't post results to Traffic Ops.\n";
   print "================================================================================\n";
   # the above line of equal signs is 80 columns
   print "Note: The user running this script must have an authorized public key on the\n";
   print "      cache servers and must also have sudoers access.\n";
   print "\n";
}
