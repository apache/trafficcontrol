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
# RTR check extension. Checks the status of the caches as seen by the Traffic Router
#
# example cron entry
# 20 * * * * root /opt/traffic_ops/app/bin/checks/ToRTRCheck.pl -c "{\"base_url\": \"https://localhost\", \"check_name\": \"RTR\", \"name\": \"Content Router\"}"
# example cron entry with syslog
# 20 * * * * root /opt/traffic_ops/app/bin/checks/ToRTRCheck.pl -c "{\"base_url\": \"https://localhost\", \"check_name\": \"RTR\", \"name\": \"Content Router\", \"syslog_facility\": \"local0\"}"

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

my $VERSION = "0.02";

STDOUT->autoflush(1);

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

DEBUG( "Including DEBUG messages in output. Config is \'" . $args{c} . "\'" );
TRACE( "Including TRACE messages in output. Config is \'" . $args{c} . "\'" );

if ( !defined( $args{c} ) ) {
	&help();
	exit(1);
}

# check the command line args
my $jconf = undef;
eval { $jconf = decode_json( $args{c} ) };
if ($@) {
	ERROR("Bad json config: $@");
	exit(1);
}

# Setup Syslogging if requested
my $sslg = undef;
if (defined($jconf->{syslog_facility})) {
   TRACE "syslog is defined";
   setlogmask(LOG_UPTO(LOG_INFO));
   openlog ('ToChecks', '', $jconf->{syslog_facility});
   $sslg = 1;
}

TRACE Dumper($jconf);
my $b_url = $jconf->{base_url};
Extensions::Helper->import();
my $ext = Extensions::Helper->new( { base_url => $b_url, token => '91504CE6-8E4A-46B2-9F9F-FE7C15228498' } );

my $jdataserver = $ext->get(Extensions::Helper::SERVERLIST_PATH);
my $check_name = $jconf->{check_name};
my $chck_lng_nm = $jconf->{name};

if ( $check_name ne "RTR" ) {
	ERROR "This Check Extension is exclusively for the RTR (Router) check.";
	exit(4);
}

my $ua = LWP::UserAgent->new;
$ua->timeout(3);

my %cdn_name;    # cdn_name by server _and_ profile... Don't name your server and profile the same.
my %ccr_assoc;
my %server_assoc;
my %api_port_assoc;
foreach my $server ( @{$jdataserver} ) {
	$server_assoc{ $server->{hostName} . "." . $server->{domainName} } = $server->{id};
	if ( !defined( $cdn_name{ $server->{profile} } ) ) {
		TRACE "Getting info for profile " . $server->{profile};
		my $plist = $ext->get( Extensions::Helper::PARAMETER_PATH . '/profile/' . $server->{profile} . '.json' );
		foreach my $param ( @{$plist} ) {
			if ( $param->{name} eq 'CDN_name' ) {
				$cdn_name{ $server->{profile} } = $param->{value};
			}
			elsif ( $param->{name} eq 'api.port' && $param->{configFile} eq 'server.xml' ) {
				$api_port_assoc{ $server->{profile} } = $param->{value};
			}
		}
	}
	$cdn_name{ $server->{hostName} . "." . $server->{domainName} } = $cdn_name{ $server->{profile} };
	next unless ( $server->{type} eq 'CCR' );

	# next unless ( $server->{host_name} eq 'odol-ccr-chi-08');
	my $new_ccr;
	$new_ccr->{status}     = "OPERATIONAL";           # we are optimistic!
	$new_ccr->{ipAddress} = $server->{ipAddress};
	my $name = $server->{hostName} . "." . $server->{domainName};
	$new_ccr->{name}     = $name;
	$new_ccr->{profile}  = $server->{profile};
	$new_ccr->{apiPort} = $server->{apiPort};
	$new_ccr->{apiPort} = $api_port_assoc{ $server->{profile} };
	$new_ccr->{cdnName} = $cdn_name{ $server->{profile} };
	$ccr_assoc{$name}    = $new_ccr;
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
      TRACE "bad_ccr++";
		$bad_ccrs->{ $ccr_assoc{$content_router}->{cdnName} }++;
   }
	$total_ccrs->{ $ccr_assoc{$content_router}->{cdnName} }++;
}

my $tmp = JSON->new->pretty(1)->encode($total_ccrs);
TRACE "total_ccrs: ".$tmp;

$tmp = JSON->new->pretty(1)->encode(\%ccr_healthy_caches);
TRACE "ccr_healthy_caches: ".$tmp;

# Check the Content Routers to see if they are aware of ONLINE Cache servers
#my $good_ccr = "";
foreach my $content_router ( keys %ccr_unhealthy_caches, keys %ccr_healthy_caches ) {

	#if ( !defined( $bad_ccrs->{ $ccr_assoc{$content_router}->{cdn_name} } ) ) {
	#	$bad_ccrs->{ $ccr_assoc{content_router}->{cdn_name} } = 0;
	#}
   TRACE "ccr_assoc{content_router}->{cdnName}: ".$ccr_assoc{$content_router}->{cdnName};
   TRACE "content_router: ".$content_router;
	if ( !defined( $ccr_healthy_caches{$content_router} ) ) {
      my $msg = $content_router . " has NO caches marked ONLINE";
		ERROR $msg;
		$ccr_assoc{$content_router}->{status} = "DOWN";
		$bad_ccrs->{ $ccr_assoc{$content_router}->{cdnName} }++;
      if ($sslg) {
         my @tmp = ($content_router, "RTR", "RTR", 'FAIL',$msg);
         syslog(LOG_ERR, "hostname=%s check=%s name=\"%s\" result=%s msg=\"%s\"", @tmp);
      }
	}
	#else {
	#	$good_ccr = $content_router;
   #   TRACE "good_ccr: ".$good_ccr;
	#}
}

$tmp = JSON->new->pretty(1)->encode(\%ccr_assoc);
TRACE "ccr_assoc: ".$tmp;

$tmp = JSON->new->pretty(1)->encode(\%healthy);
TRACE "healthy: ".$tmp;

# Check to see if the cache is being reported healthy by all the content routers
foreach my $cache ( sort keys %healthy ) {

   DEBUG "total_ccrs: ".$total_ccrs->{ $cdn_name{$cache} } . " bad_ccrs: " .(defined( $bad_ccrs->{ $cdn_name{$cache} } ) ? $bad_ccrs->{ $cdn_name{$cache} } : 0 ). "\n";
   DEBUG "healthy{cache}: ".$healthy{$cache};
   DEBUG "cache: ".$cache;
	if ( $healthy{$cache} < ( $total_ccrs->{ $cdn_name{$cache} } - ( defined( $bad_ccrs->{ $cdn_name{$cache} } ) ? $bad_ccrs->{ $cdn_name{$cache} } : 0 ) ) ) {
      my $msg = $cache . " => " . $healthy{$cache} . " out of " . ( $total_ccrs - $bad_ccrs ) . " healthy CCRs think it is OK.";
		ERROR $msg;
      if ($sslg) {
         my @tmp = ($cache, "RTR", "RTR", 'FAIL',$msg);
         syslog(LOG_ERR, "hostname=%s check=%s name=\"%s\" result=%s msg=\"%s\"", @tmp);
      }
		$ext->post_result( $server_assoc{$cache}, $check_name, 0 );
	}
}

# Report on caches that are OFFLINE
foreach my $cache ( sort keys %unhealthy ) {
	if ( $unhealthy{$cache} == $total_ccrs->{ $cdn_name{$cache} } ) {
      my $msg = $cache . " is marked OFFLINE by all CCRs.";
		ERROR $msg;
      if ($sslg) {
         my @tmp = ($cache, "RTR", "RTR", 'FAIL',$msg);
         syslog(LOG_ERR, "hostname=%s check=%s name=\"%s\" result=%s msg=\"%s\"", @tmp);
      }
		$ext->post_result( $server_assoc{$cache}, $check_name, 0 );
	}
}
# Report on healthy caches
foreach my $cache ( sort keys %healthy ) {
	if ( $healthy{$cache} == $total_ccrs->{ $cdn_name{$cache} } ) {
      my $msg = $cache . " is marked ONLINE by all CCRs.";
		INFO $msg;
      if ($sslg) {
         my @tmp = ($cache, "RTR", "RTR", 'OK');
         syslog(LOG_ERR, "hostname=%s check=%s name=\"%s\" result=%s", @tmp);
      }
		$ext->post_result( $server_assoc{$cache}, $check_name, 1 );
	}
}
# Report on healthy content routers.
foreach my $ccr (keys %ccr_assoc) {
   # By this time we should have already logged any 'FAIL' messages
   if ($ccr_assoc{$ccr}->{status} eq "OPERATIONAL") {
      my $msg = $ccr . " is OK.";
		INFO $msg;
      if ($sslg) {
         my @tmp = ($ccr, "RTR", "RTR", 'OK');
         syslog(LOG_ERR, "hostname=%s check=%s name=\"%s\" result=%s", @tmp);
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
	my $response = $ua->get( $url );

	if ( !$response->is_success ) {
		ERROR $ccr . " Not responding! - " . $response->status_line . ".";
		$ccr_assoc{$ccr_name}->{status} = "DOWN";
      # log a FAIL message for this content router.
      if ($sslg) {
         my @tmp = ($ccr_name, "RTR", "RTR", 'FAIL',
                    "Unable to connect to content router API. response: ".$response->status_line);
         syslog(LOG_ERR, "hostname=%s check=%s name=\"%s\" result=%s msg=\"%s\"", @tmp);
      }
		return;
	}
	my $loc_var = JSON->new->utf8->decode( $response->content );

   # get the status of each location from the content router
	my $status;
	foreach my $location ( sort @{ $loc_var->{locations} } ) {
		TRACE "getting " . $url ."/".$location . "/caches";
		my $response = $ua->get( $url."/".$location."/caches" );
		if ( !$response->is_success ) {
			ERROR $ccr . " Not responding! - " . $response->status_line . ".";
         if ($sslg) {
            # Log a FAIL message for this content router
            my @tmp = ($ccr_name, "RTR", "RTR", 'FAIL',
                       "Unable to get cache info for $location. response: ".$response->status_line);
            syslog(LOG_ERR, "hostname=%s check=%s name=\"%s\" result=%s msg=\"%s\"", @tmp);
         }
			next;
		}
		my $loc_stat = JSON->new->utf8->decode( $response->content );
		$status->{$location} = $loc_stat;
	}
	return $status;
}

# TODO
sub help {
   print "-c is a required option\n";
}
