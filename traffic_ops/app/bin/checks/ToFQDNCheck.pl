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

# Plugin for the "FQDN" check.
#
# example cron entry
# 0 * * * * root /opt/traffic_ops/app/bin/checks/ToFQDNCheck.pl -c "{\"base_url\": \"https://localhost\", \"check_name\": \"FQDN\"}" >> /var/log/traffic_ops/extensionCheck.log 2>&1
# example cron entry with syslog
# 0 * * * * root /opt/traffic_ops/app/bin/checks/ToFQDNCheck.pl -c "{\"base_url\": \"https://localhost\", \"check_name\": \"FQDN\", \"name\": \"DNS Lookup\", \"syslog_facility\": \"local0\"}" > /dev/null 2>&1

use strict;
use warnings;

use IO::Handle;
use Log::Log4perl qw(:easy);
use Data::Dumper;
use Getopt::Std;
use JSON;
use Extensions::Helper;
use Sys::Syslog qw(:standard :macros);
use Net::DNS;
use NetAddr::IP;

my $VERSION = "0.02";

STDOUT->autoflush(1);

my %args = ();
getopts( "c:f:hl:q", \%args );

if ($args{h}) {
   &help();
   exit();
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
my $chck_lng_nm;
if (defined($jconf->{syslog_facility})) {
   $chck_lng_nm = $jconf->{name};
   setlogmask(LOG_UPTO(LOG_INFO));
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

TRACE Dumper($jconf);
my $b_url = $jconf->{base_url};
Extensions::Helper->import();
my $ext = Extensions::Helper->new( { base_url => $b_url, token => '91504CE6-8E4A-46B2-9F9F-FE7C15228498' } );

my $jdataserver    = $ext->get(Extensions::Helper::SERVERLIST_PATH);
my $chck_nm     = $jconf->{check_name};
foreach my $server ( @{$jdataserver} ) {
	if ( $server->{type} =~ m/^EDGE/ || $server->{type} =~ m/^MID/ ) {
      my $status = 1;
      my $srv_nm = $server->{hostName}.".".$server->{domainName};
      my $srv_ip = $server->{ipAddress};
      my $srv_ip6;

      if (defined($server->{ip6Address})) {
         $srv_ip6 = $server->{ip6Address};
         $srv_ip6 =~ s/\/\d+$//;
      } else {
         $srv_ip6 = 'not defined';
      }

      my @rslt;
      if ($force == 0) {
         @rslt = &fqdn_check( $srv_nm, "A" );
      } elsif ($force == 1) {
         @rslt = (0, "match");
      } elsif ($force == 2) {
         @rslt = (0, "NXDOMAIN");
      } elsif ($force == 3) {
         @rslt = (1, "127.0.0.1");
      } elsif ($force == 4) {
         @rslt = (1, $srv_ip);
      }
      if (!$rslt[0]) {  # IPv4 query failed
         # IPv4 DNS lookup failed
         if ($rslt[1] =~ m/match/) {
            if ($sslg) {
               my @tmp = ($srv_nm, $chck_nm, $chck_lng_nm, 'FAIL',
                          "DNS A record '' does not match IPv4 in Traffic Ops $srv_ip");
               syslog(LOG_ERR, "hostname=%s check=%s name=\"%s\" result=%s target=A msg=\"%s\"", @tmp);
            }
            ERROR "hostname: ".$srv_nm." check: ".$chck_nm." "
         } else {
            if ($sslg) {
               my @tmp = ($srv_nm, $chck_nm, $chck_lng_nm, 'FAIL',$rslt[1]);
               syslog(LOG_ERR, "hostname=%s check=%s name=\"%s\" result=%s target=A msg=\"%s\"", @tmp);
            }
         }
         my $tmp = "hostname: ".$srv_nm;
         $tmp .= " result: FAIL";
         $tmp .= " dns return: ".$rslt[1]." db: ".$srv_ip;
         ERROR $tmp;
         $status = 0;
      } elsif ($rslt[1] !~ m/$srv_ip/) { # IPv4 query success
         # check to see if the IPv4 from DNS matches DB
         if ($sslg) {
            my @tmp = ($srv_nm, $chck_nm, $chck_lng_nm, 'FAIL',
                       "DNS A record $rslt[1] does not match IPv4 in Traffic Ops $srv_ip");
            syslog(LOG_ERR, "hostname=%s check=%s name=\"%s\" result=%s target=A msg=\"%s\"", @tmp);
         }
         my $tmp = "hostname: ".$srv_nm;
         $tmp .= " result: FAIL";
         $tmp .= " dns return: ".$rslt[1]." db: ".$srv_ip;
         ERROR $tmp;
         $status = 0;
      } else {
         if ($sslg) {
            my @tmp = ($srv_nm, $chck_nm, $chck_lng_nm, 'OK');
            syslog(LOG_INFO, "hostname=%s check=%s name=\"%s\" result=%s target=A msg=\"\"", @tmp);
         }
      }

      # Check IPv6
      if ($srv_ip6 !~ m/not defined/) {
         if ($force == 0) {
            @rslt = &fqdn_check($srv_nm, "AAAA");
         } elsif ($force == 1) {
            @rslt = (0, "match");
         } elsif ($force == 2) {
            @rslt = (0, "NXDOMAIN");
         } elsif ($force == 3) {
            @rslt = (1, "::1");
         } elsif ($force == 4) {
            @rslt = (1, $srv_ip6);
         }
         if (!$rslt[0]) {
            # IPv6 DNS lookup failed
            if ($rslt[1] =~ m/match/) {
               if ($sslg) {
                  my @tmp = ($srv_nm, $chck_nm, $chck_lng_nm, 'FAIL',
                             "DNS AAAA record '' does not match IPv6 in Traffic Ops $srv_ip6");
                  syslog(LOG_ERR, "hostname=%s check=%s name=\"%s\" result=%s target=AAAA msg=\"%s\"", @tmp);
               }
            } else {
               if ($sslg) {
                  my @tmp = ($srv_nm, $chck_nm, $chck_lng_nm, 'FAIL',$rslt[1]);
                  syslog(LOG_ERR, "hostname=%s check=%s name=\"%s\" result=%s target=AAAA msg=\"%s\"", @tmp);
               }
            }
            my $tmp = "hostname: ".$srv_nm;
            $tmp .= " result: FAIL";
            $tmp .= " dns return: ".$rslt[1]." db: ".$srv_ip6;
            ERROR $tmp;
            $status = 0;
         } elsif ($rslt[1] !~ m/$srv_ip6/i) {
            # check to see if the IPv6 from DNS matches DB
            if ($sslg) {
               my @tmp = ($srv_nm, $chck_nm, $chck_lng_nm, 'FAIL',
                          'DNS AAAA record $rslt[1] does not match IPv6 in Traffic Ops $srv_ip6');
               syslog(LOG_ERR, "hostname=%s check=%s name=\"%s\" result=%s target=AAAA msg=\"%s\"", @tmp);
            }
            my $tmp = "hostname: ".$srv_nm;
            $tmp .= " result: FAIL";
            $tmp .= " dns return: ".$rslt[1]." db: ".$srv_ip6;
            ERROR $tmp;
            $status = 0;
         } else {
            if ($sslg) {
               my @tmp = ($srv_nm, $chck_nm, $chck_lng_nm, 'OK');
               syslog(LOG_INFO, "hostname=%s check=%s name=\"%s\" result=%s target=AAAA msg=\"\"", @tmp);
            }
         }
      }

		DEBUG $chck_nm . " >> ".$srv_nm." result: ".$rslt[1]." status: ".$status;
		$ext->post_result( $server->{id}, $chck_nm, $status )
         if (!$quiet);
	}
}

closelog();

sub fqdn_check {
   my ($hostname,$type) = @_;

   my ($resolver,$reply,$matched,$ip);
   my (@result);

   $resolver = Net::DNS::Resolver->new;
   # I tried query($hostname,$type,''); but, it didn't seem to work right.
   if ($type =~ m/AAAA/) {
      $reply = $resolver->query($hostname,'AAAA','');
   } else {
      $reply = $resolver->query($hostname);
   }

   #DEBUG "Net::DNS version: ".Net::DNS->version;
   DEBUG "hostname: ".$hostname." type: ".$type."\n";
   
   if ($reply) {
      foreach my $rr ($reply->answer) {
         DEBUG "hostname: ".$hostname." answer: ".$rr->address;
         next unless $rr->type eq $type;
         $matched = 1;
         if ($type =~ m/AAAA/) {
            $ip = new NetAddr::IP ($rr->address);
            $ip = $ip->short;
            TRACE "New IP: ".$ip;
         } else {
            $ip = $rr->address;
         }

         @result = (1,$ip);
         my @tmp = ($hostname, $ip);
         DEBUG "hostname: ".$hostname." address: ".$ip;
      }
   } else {
      @result = (0,$resolver->errorstring);
      #my @tmp = ($hostname, $resolver->errorstring);
      WARN "query failed: ".$resolver->errorstring;
   }
   if (!$matched) { # makes sure type matched
      my $msg = "match";
      @result = (0,$msg);
   }

   return @result;
}

sub help {
   print "ToFQDNCheck.pl -c \"{\\\"base_url\\\": \\\"https://localhost\\\", \\\"check_name\\\": \\\"FQDN\\\"[, \\\"name\\\": \\\"DNS Lookup\\\", \\\"syslog_facility\\\": \\\"local0\\\"]}\" [-f <1-4>] [-l <1-3>]\n";
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
   print "        1: FAIL Blank A record in DNS\n";
   print "        2: FAIL DNS failure\n";
   print "        3: FAIL mis-match between DNS and Traffic Ops.\n";
   print "        4: OK\n";
   print "-h   Print this message\n";
   print "-l   Debug level\n";
   print "-q   Don't post results to Traffic Ops.\n";
   print "================================================================================\n";
   # the above line of equal signs is 80 columns
   print "\n";
}

