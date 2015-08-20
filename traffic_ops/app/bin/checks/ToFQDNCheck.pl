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

# Plugin for the "FQDN" check.
#
# example cron entry
# 0 * * * * /opt/traffic_ops/app/bin/checks/ToFQDNCheck.pl -c "{\"base_url\": \"https://localhost\", \"check_name\": \"FQDN\", \"name\": \"FQDN\"}"
# example cron entry with syslog
# 0 * * * * /opt/traffic_ops/app/bin/checks/ToFQDNCheck.pl -c "{\"base_url\": \"https://localhost\", \"check_name\": \"FQDN\", \"name\": \"FQDN\", \"syslog_facility\": \"local0\"}"

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

my $VERSION = "0.01";

STDOUT->autoflush(1);

my %args = ();
getopts( "l:c:", \%args );

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
if (defined($jconf->{syslog_facility})) {
   setlogmask(LOG_UPTO(LOG_INFO));
   openlog ('ToChecks', '', $jconf->{syslog_facility});
   $sslg = 1;
}

TRACE Dumper($jconf);
my $b_url = $jconf->{base_url};
Extensions::Helper->import();
my $ext = Extensions::Helper->new( { base_url => $b_url, token => '91504CE6-8E4A-46B2-9F9F-FE7C15228498' } );

my $jdataserver    = $ext->get(Extensions::Helper::SERVERLIST_PATH);
my $chck_nm     = $jconf->{check_name};
my $chck_lng_nm    = $jconf->{name};
foreach my $server ( @{$jdataserver} ) {
	if ( $server->{type} eq 'EDGE' || $server->{type} eq 'MID' ) {
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

      my @rslt = &fqdn_check( $srv_nm, "A" );
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
            syslog(LOG_INFO, "hostname=%s check=%s name=\"%s\" result=%s target=A", @tmp);
         }
      }

      # Check IPv6
      if ($srv_ip6 !~ m/not defined/) {
         @rslt = &fqdn_check($srv_nm, "AAAA");
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
               syslog(LOG_INFO, "hostname=%s check=%s name=\"%s\" result=%s target=AAAA", @tmp);
            }
         }
      }

		DEBUG $chck_nm . " >> ".$srv_nm." result: ".$rslt[1]." status: ".$status;
      # assuming that if someone is asking for output they are debugging script
      # and don't post to DB. This allows for testing on prod server because
      # it is hard to test everything in the lab. Is it not? And who's going to
      # know anyway besides you and me. ;)
		$ext->post_result( $server->{id}, $chck_nm, $status )
         if (!defined($args{l}));
	}
}

closelog();

sub help {
	print "The -c argument is mandatory\n";
}

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
