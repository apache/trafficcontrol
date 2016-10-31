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
# RTR check extension. Checks the status of the caches as seen by the Traffic Router
#
# example cron entry
# 20 * * * * root /opt/traffic_ops/app/bin/checks/ToRTRCheck.pl -c "{\"base_url\": \"https://localhost\", \"check_name\": \"RTR\", \"name\": \"Content Router Check\"}" >> /var/log/traffic_ops/extensionCheck.log 2>&1
#
# example cron entry with syslog
# 20 * * * * root /opt/traffic_ops/app/bin/checks/ToRTRCheck.pl -c "{\"base_url\": \"https://localhost\", \"check_name\": \"RTR\", \"name\": \"Content Router Check\", \"syslog_facility\": \"local0\"}" > /dev/null 2>&1

use strict;
use warnings;

use LWP::UserAgent;
use Data::Dumper;
use Getopt::Std;
use Log::Log4perl qw(:easy);
use JSON;
use Extensions::Helper;
use Sys::Syslog qw(:standard :macros);
use IO::Handle;

my $VERSION = "0.04";

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
	elsif ( $args{l} == 3 ) { Log::Log4perl->easy_init($TRACE); }
	elsif ( $args{l} > 3 )  { Log::Log4perl->easy_init($TRACE); }
	else                    { Log::Log4perl->easy_init($INFO); }
}

DEBUG( "Including DEBUG messages in output. Config is \'" . $args{c} . "\'" );
TRACE( "Including TRACE messages in output. Config is \'" . $args{c} . "\'" );

if ( !defined( $args{c} ) ) {
   ERROR "-c not defined";
   print "\n\n";
	&help();
	exit(1);
}

# check the command line args
my $jconf = undef;
eval { $jconf = decode_json( $args{c} ) };
if ($@) {
	ERROR("Bad json config: $@");
   print "\n\n";
   &help();
	exit(1);
}

# Setup Syslogging if requested
my $sslg = undef;
my $chck_lng_nm;
if (defined($jconf->{syslog_facility})) {
   $chck_lng_nm = $jconf->{name};
   TRACE "syslog is defined";
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

my $check_name = $jconf->{check_name};
if ( $check_name ne "RTR" ) {
	ERROR "This Check Extension is exclusively for the RTR (Router) check.";
   print "\n\n";
   &help();
	exit(4);
}

TRACE Dumper($jconf);
my $b_url = $jconf->{base_url};
Extensions::Helper->import();
my $ext = Extensions::Helper->new( { base_url => $b_url, token => '91504CE6-8E4A-46B2-9F9F-FE7C15228498' } );

my $jdataserver = $ext->get(Extensions::Helper::SERVERLIST_PATH);

my $ua = LWP::UserAgent->new;
$ua->timeout(3);

my %cdn_name;    # cdn_name by server _and_ profile... Don't name your server and profile the same.
my %ccr_assoc;
my %server_assoc;
my %api_port_assoc;
foreach my $server ( @{$jdataserver} ) {
	$server_assoc{ $server->{hostName} . "." . $server->{domainName} }->{id} = $server->{id};
	$server_assoc{ $server->{hostName} . "." . $server->{domainName} }->{cmsStatus} = $server->{status};
	if ( !defined( $cdn_name{ $server->{profile} } ) ) {
		TRACE "Getting info for server profile " . $server->{profile} . " for server: " . $server->{hostName};
		my $plist = $ext->get( Extensions::Helper::PARAMETER_PATH . '/profile/' . $server->{profile} . '.json' );
      foreach my $param ( @{$plist} ) {
			if ( $param->{name} eq 'api.port' && $param->{configFile} eq 'server.xml' ) {
				$api_port_assoc{ $server->{profile} } = $param->{value};
			}
		}
	}
	$cdn_name{ $server->{hostName} . "." . $server->{domainName} } = $server->{cdnName};
	next unless ( $server->{type} eq 'CCR' );
   DEBUG Dumper($server);

	# next unless ( $server->{host_name} eq 'odol-ccr-chi-08');
	my $new_ccr;
	$new_ccr->{status}    = "OPERATIONAL";           # we are optimistic!
   $new_ccr->{cmsStatus} = $server->{status};
	$new_ccr->{ipAddress} = $server->{ipAddress};
	my $name              = $server->{hostName} . "." . $server->{domainName};
	$new_ccr->{name}      = $name;
	$new_ccr->{profile}   = $server->{profile};
	$new_ccr->{apiPort}   = $server->{apiPort};
	$new_ccr->{apiPort}   = $api_port_assoc{ $server->{profile} };
	$new_ccr->{cdnName}   = $server->{cdnName};
	$ccr_assoc{$name}     = $new_ccr;
	DEBUG "Adding CCR " . $new_ccr->{name};
}

INFO "Starting CCR /crs page checks";
my @ccr_err = ();
my %healthy;
my %unhealthy;
my %ccr_healthy_caches;
my %ccr_unhealthy_caches;

my $total_ccrs;
my $bad_ccrs;
foreach my $content_router ( keys %ccr_assoc ) {
   TRACE "content_router: ".$content_router;
	my $ccr_status = &get_crs_stat( $ccr_assoc{$content_router}->{ipAddress},
                                   $ccr_assoc{$content_router}->{apiPort},
                                   $ccr_assoc{$content_router}->{name}
                                 );

   if (ref($ccr_status) eq 'HASH') {
	   foreach my $loc ( keys %{$ccr_status} ) {
	   	foreach my $cache ( @{ $ccr_status->{$loc}->{caches} } ) {
            if (($force == 3) || ($force == 4) || ($force == 5)) {
               # Force healthy
               $cache->{cacheOnline} = 1;
            }
            elsif (($force == 6) || ($force == 7)) {
               # Force unhealthy
               $cache->{cacheOnline} = 0;
            }
	   		if ( $cache->{cacheOnline} ) {
	   			$ccr_healthy_caches{$content_router}++;
               TRACE "healthy_cache: ".$cache->{fqdn};
	   			$healthy{ $cache->{fqdn} }++;
	   		}
	   		else {
	   			$ccr_unhealthy_caches{$content_router}++;
               TRACE "unhealthy_cache: ".$cache->{fqdn};
	   			$unhealthy{ $cache->{fqdn} }++;
	   		}
	   	}
	   }
   } else {
      TRACE "incrementing bad_ccrs";
		$bad_ccrs->{ $ccr_assoc{$content_router}->{cdnName} }++;
   }
	$total_ccrs->{ $ccr_assoc{$content_router}->{cdnName} }++;
}

my $tmp = JSON->new->pretty(1)->encode($total_ccrs);
TRACE "total_ccrs: ".$tmp;

$tmp = JSON->new->pretty(1)->encode(\%ccr_healthy_caches);
TRACE "ccr_healthy_caches: ".$tmp;

# Check the Content Routers to see if they are aware of ONLINE Cache servers
if (($force == 0) || ($force == 1) || ($force == 2)) {
   foreach my $content_router ( keys %ccr_unhealthy_caches, keys %ccr_healthy_caches ) {

      #if ( !defined( $bad_ccrs->{ $ccr_assoc{$content_router}->{cdn_name} } ) ) {
      #	$bad_ccrs->{ $ccr_assoc{content_router}->{cdn_name} } = 0;
      #}
      TRACE "ccr_assoc{content_router}->{cdnName}: ".$ccr_assoc{$content_router}->{cdnName};
      TRACE "content_router: ".$content_router;
      if ($force > 0) {
         $ccr_healthy_caches{$content_router} = undef;
      }
      if (!defined( $ccr_healthy_caches{$content_router} )) {
         my $msg = $content_router . " has NO caches marked ONLINE";
         ERROR $msg;
         $ccr_assoc{$content_router}->{status} = "DOWN";
         $bad_ccrs->{ $ccr_assoc{$content_router}->{cdnName} }++;
         if ($sslg) {
            my $result;
            # Only FAIL if online or reported
            if (($ccr_assoc{$content_router}->{cmsStatus} =~ m/ONLINE/) || ($ccr_assoc{$content_router}->{cmsStatus} =~ m/REPORTED/)) {
               $result = "FAIL";
            }
            else {
               $result = "OK";
            }
            if ($force == 1) {
               $result = "FAIL";
               $msg = "Force: FAIL ".$msg;
            }
            elsif ($force == 2) {
               $result = "OK";
               $msg = "Force: OK ".$msg;
            }
            my @tmp = ($content_router,$check_name,$chck_lng_nm,$result,$ccr_assoc{$content_router}->{cmsStatus},$msg);
            syslog(LOG_ERR, "hostname=%s check=%s name=\"%s\" result=%s status=%s msg=\"%s\"", @tmp);
         }
      }
   }
}

$tmp = JSON->new->pretty(1)->encode(\%ccr_assoc);
TRACE "ccr_assoc: ".$tmp;

$tmp = JSON->new->pretty(1)->encode(\%healthy);
TRACE "healthy: ".$tmp;

# Check to see if the cache is being reported healthy by all the content routers
if (($force == 0) || ($force == 3) || ($force == 4) || ($force == 5)) {
   foreach my $cache ( sort keys %healthy ) {

      DEBUG "total_ccrs: ".$total_ccrs->{ $cdn_name{$cache} } . " bad_ccrs: " .(defined( $bad_ccrs->{ $cdn_name{$cache} } ) ? $bad_ccrs->{ $cdn_name{$cache} } : 0 ). "\n";
      DEBUG "healthy{cache}: ".$healthy{$cache};
      DEBUG "cache: ".$cache;
      if ($force == 3) {
         $healthy{$cache} = 0;
         $server_assoc{$cache}->{cmsStatus} = "FORCE_FAIL";
      }
      elsif ($force == 4) {
         $healthy{$cache} = 0;
         $server_assoc{$cache}->{cmsStatus} = "FORCE_OK";
      }
      elsif ($force == 5) {
         $healthy{$cache} = $total_ccrs->{ $cdn_name{$cache} }
      }
      if ($healthy{$cache} < ( $total_ccrs->{ $cdn_name{$cache} } - ( defined( $bad_ccrs->{ $cdn_name{$cache} } ) ? $bad_ccrs->{ $cdn_name{$cache} } : 0 ) ))  {
         my $msg = $cache . " => " . $healthy{$cache} . " out of " . ( $total_ccrs->{ $cdn_name{$cache} } - ( defined( $bad_ccrs->{ $cdn_name{$cache} } ) ? $bad_ccrs->{ $cdn_name{$cache} } : 0 )) . " healthy Content Routers think it is OK.";
         ERROR $msg;
         if ($sslg) {
            my $result;
            # Only FAIL if cache online or reported
            if (($server_assoc{$cache}->{cmsStatus} =~ m/ONLINE/) || ($server_assoc{$cache}->{cmsStatus} =~ m/REPORTED/)) {
               $result = "FAIL";
            }
            else {
               $result = "OK";
            }
            if ($force == 3) {
               $result = "FAIL";
               $msg = "Force: FAIL ".$msg;
            }
            elsif ($force == 4) {
               $result = "OK";
               $msg = "Force: OK ".$msg;
            }
            my @tmp = ($cache,$check_name,$chck_lng_nm,$result,$server_assoc{$cache}->{cmsStatus},$msg);
            syslog(LOG_ERR, "hostname=%s check=%s name=\"%s\" result=%s status=%s msg=\"%s\"", @tmp);
         }
         $ext->post_result( $server_assoc{$cache}->{id}, $check_name, 0 ) if (!$quiet);
      }
      elsif ($healthy{$cache} == $total_ccrs->{ $cdn_name{$cache} }) {
         my $msg = $cache . " is marked ONLINE by all Content Routers.";
         if ($force == 5) {
            $msg = "Force: OK ".$msg;
         }
         INFO $msg;
         if ($sslg) {
            my @tmp = ($cache,$check_name,$chck_lng_nm, 'OK',$server_assoc{$cache}->{cmsStatus},$msg);
            syslog(LOG_ERR, "hostname=%s check=%s name=\"%s\" result=%s status=%s msg=\"%s\"", @tmp);
         }

         $ext->post_result( $server_assoc{$cache}->{id}, $check_name, 1 ) if (!$quiet);
      }
      else {
         # Should never get here
         ERROR $healthy{$cache}."Didn't match for any health checks";
      }
   }
}

# Report on caches that are OFFLINE
if (($force == 0) || ($force == 6) || ($force == 7)) {
   foreach my $cache ( sort keys %unhealthy ) {
      if ($force == 6) {
         $server_assoc{$cache}->{cmsStatus} = "FORCE_FAIL";
         $unhealthy{$cache} = $total_ccrs->{ $cdn_name{$cache} };
      } elsif ($force == 7) {
         $server_assoc{$cache}->{cmsStatus} = "FORCE_OK";
         $unhealthy{$cache} = $total_ccrs->{ $cdn_name{$cache} };
      }
      if ($unhealthy{$cache} == $total_ccrs->{ $cdn_name{$cache} }) {
         my $msg = $cache . " is marked OFFLINE by all Content Routers.";
         ERROR $msg;
         if ($sslg) {
            my $result;
            # Only FAIL if cache online or reported
            if (($server_assoc{$cache}->{cmsStatus} =~ m/ONLINE/) || ($server_assoc{$cache}->{cmsStatus} =~ m/REPORTED/)) {
               $result = "FAIL";
            }
            else {
               $result = "OK";
            }
            if ($force == 6) {
               $msg = "Force: FAIL ".$msg;
               $result = "FAIL";
            }
            elsif ($force == 7) {
               $msg = "Force: OK ".$msg;
               $result = "OK";
            }
            my @tmp = ($cache,$check_name,$chck_lng_nm,$result,$server_assoc{$cache}->{cmsStatus},$msg);
            syslog(LOG_ERR, "hostname=%s check=%s name=\"%s\" result=%s status=%s msg=\"%s\"", @tmp);
         }
         $ext->post_result( $server_assoc{$cache}->{id}, $check_name, 0 ) if (!$quiet);
      }
   }
}

# Report on healthy content routers.
if (($force == 0) || ($force == 8)) {
   foreach my $ccr (keys %ccr_assoc) {
      # By this time we should have already logged any 'FAIL' messages
      if (($ccr_assoc{$ccr}->{status} eq "OPERATIONAL") || ($force == 8)) {
         my $msg = $ccr . " is OK.";
         if ($force == 8) {
            $msg = "Force: OK ".$msg;
         }
         INFO $msg;
         if ($sslg) {
            my @tmp = ($ccr,$check_name,$chck_lng_nm,'OK',$ccr_assoc{$ccr}->{cmsStatus},$msg);
            syslog(LOG_ERR, "hostname=%s check=%s name=\"%s\" result=%s status=%s msg=\"%s\"", @tmp);
         }
      }
   }
}

closelog();

sub get_crs_stat() {
	my $ccr  = shift;
	my $port = shift || 80;
   my $ccr_name = shift;
	my $url  = 'http://' . $ccr . ":" . $port . '/crs/locations';

	my $ua = LWP::UserAgent->new;
	TRACE "getting locations: " . $url;
   TRACE "Force: ".$force;
	my $response = $ua->get( $url );

	if ( !$response->is_success || ($force == 9) || ($force == 10)) {
		ERROR $ccr . " Not responding! - " . $response->status_line . ".";
		$ccr_assoc{$ccr_name}->{status} = "DOWN";
      # log a FAIL message for this content router.
      if ($sslg) {
         my $result;
         my $msg = "Unable to connect to content router API. response: ".$response->status_line;
         # Only FAIL if online or reported
         if (($ccr_assoc{$ccr_name}->{cmsStatus} =~ m/ONLINE/) || ($ccr_assoc{$ccr_name}->{cmsStatus} =~ m/REPORTED/)) {
            $result = "FAIL";
         }
         else {
            $result = "OK";
         }
         if ($force == 9) {
            $result = "FAIL";
            $msg = "Force: FAIL ".$msg;
         } elsif ($force == 10) {
            $result = "OK";
            $msg = "Force: OK ".$msg;
         }
         my @tmp = ($ccr_name,$check_name,$chck_lng_nm,$result,
                    $ccr_assoc{$ccr_name}->{cmsStatus},$msg);
         syslog(LOG_ERR, "hostname=%s check=%s name=\"%s\" result=%s status=%s msg=\"%s\"", @tmp);
      }
		return;
	}
	my $loc_var = JSON->new->utf8->decode( $response->content );

   # get the status of each location from the content router
	my $status;
	foreach my $location ( sort @{ $loc_var->{locations} } ) {
		TRACE "getting " . $url ."/".$location . "/caches";
		my $response = $ua->get( $url."/".$location."/caches" );
		if ( !$response->is_success  || ($force == 11) || ($force == 12)) {
			ERROR $ccr . " Not responding! - " . $response->status_line . ".";
         if ($sslg) {
            my $result;
            my $msg = "Unable to get cache info for $location. response: ".$response->status_line;
            # Only FAIL if online or reported
            if (($ccr_assoc{$ccr_name}->{cmsStatus} =~ m/ONLINE/) || ($ccr_assoc{$ccr_name}->{cmsStatus} =~ m/REPORTED/)) {
               $result = "FAIL";
            }
            else {
               $result = "OK";
            }
            if ($force == 11) {
               $result = "FAIL";
               $msg = "Force: FAIL ".$msg
            } elsif ($force == 12) {
               $result = "OK";
               $msg = "Force: OK ".$msg
            }
            # Log a FAIL message for this content router
            my @tmp = ($ccr_name,$check_name,$chck_lng_nm,$result,
                       $ccr_assoc{$ccr_name}->{cmsStatus},$msg);
            syslog(LOG_ERR, "hostname=%s check=%s name=\"%s\" result=%s status=%s msg=\"%s\"", @tmp);
         }
			next;
		}
		my $loc_stat = JSON->new->utf8->decode( $response->content );
		$status->{$location} = $loc_stat;
	}
	return $status;
}

sub help() {
   print "ToRTRCheck.pl -c \"{\\\"base_url\\\": \\\"https://localhost\\\", \\\"check_name\\\": \\\"RTR\\\"[, \\\"name\\\": \\\"Content Router Checks\\\", \\\"syslog_facility\\\": \\\"local0\\\"}]\" [-f <1-8>] [-l <1-3>]\n";
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
   print "        1: FAIL - Content Router\n";
   print "        2: OK - Content Router\n";
   print "        3: FAIL - # of CRs reporting on cache\n";
   print "        4: OK - # of CRs reporting on cache\n";
   print "        5: OK - Cache healthy by all CRs\n";
   print "        6: FAIL - Cache offline but ONLINE or REPORTED in Traffic Ops.\n";
   print "        7: OK - Cache offline but ONLINE or REPORTED in Traffic Ops.\n";
   print "        8: OK - For all Content Routers.\n";
   print "        9: FAIL - Could not connect to Content Router API. Content Router is\n";
   print "             ONLINE or REPORTED in Traffic Ops.\n";
   print "        10: OK - Could not connect to Content Router API. Content Router is not\n";
   print "             ONLINE or REPORTED in Traffic Ops.\n";
   print "        11: FAIL - Could not get cache info for a location from Content Router\n";
   print "             API. Content Router is ONLINE or REPORTED in Traffic Ops.\n";
   print "        12: OK - Could not get cache info for a location Content Router API.\n";
   print "             Content Router is not ONLINE or REPORTED in Traffic Ops.\n";
   print "-h   Print this message\n";
   print "-l   Debug level\n";
   print "-q   Don't post results to Traffic Ops.\n";
   print "================================================================================\n";
   # the above line of equal signs is 80 columns
   print "\n";
}
