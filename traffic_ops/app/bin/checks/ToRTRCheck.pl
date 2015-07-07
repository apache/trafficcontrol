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
# TRTR check extension. Checks the status of the caches as seen by the Traffic Router
#

use strict;
use warnings;

$|++;

use LWP::UserAgent;
use Data::Dumper;
use Getopt::Std;
use Log::Log4perl qw(:easy);
use JSON;
use Extensions::Helper;

my $VERSION = "0.01";
my $hostn   = `hostname`;
chomp($hostn);

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

my $jconf = undef;
eval { $jconf = decode_json( $args{c} ) };
if ($@) {
	ERROR("Bad json config: $@");
	exit(1);
}

# TRACE Dumper($jconf);
my $b_url = $jconf->{base_url};
Extensions::Helper->import();
my $ext = Extensions::Helper->new( { base_url => $b_url, token => '91504CE6-8E4A-46B2-9F9F-FE7C15228498' } );

my $jdataserver = $ext->get(Extensions::Helper::SERVERLIST_PATH);
my $match       = $jconf->{match};

my $check_name = $jconf->{check_name};
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
foreach my $content_router ( keys %ccr_assoc ) {
	my $ccr_status = &get_crs_stat( $ccr_assoc{$content_router}->{ipAddress}, $ccr_assoc{$content_router}->{apiPort} );
	foreach my $loc ( keys %{$ccr_status} ) {
		foreach my $cache ( @{ $ccr_status->{$loc}->{caches} } ) {
			if ( $cache->{cacheOnline} ) {
				$ccr_healthy_caches{$content_router}++;
				$healthy{ $cache->{fqdn} }++;
			}
			else {
				$ccr_unhealthy_caches{$content_router}++;
				$unhealthy{ $cache->{fqdn} }++;
			}
		}
	}
	$total_ccrs->{ $ccr_assoc{$content_router}->{cdnName} }++;
}

my $bad_ccrs;
my $good_ccr = "";
foreach my $content_router ( keys %ccr_unhealthy_caches, keys %ccr_healthy_caches ) {

	# if ( !defined( $bad_ccrs->{ $ccr_assoc{$content_router}->{cdn_name} } ) ) {
	# 	$bad_ccrs->{ $ccr_assoc{content_router}->{cdn_name} } = 0;
	# }
	if ( !defined( $ccr_healthy_caches{$content_router} ) ) {
		ERROR $content_router . " has NO caches marked ONLINE - disconnected from XMPP!?!?!?";
		$ccr_assoc{$content_router}->{status} = "DOWN";
		$bad_ccrs->{ $ccr_assoc{$content_router}->{cdnName} }++;
	}
	else {
		$good_ccr = $content_router;
	}
}

foreach my $cache ( sort keys %healthy ) {

	# print $total_ccrs->{ $cdn_name{$cache} } . " - " . $bad_ccrs->{ $cdn_name{$cache} } . "\n";
	if ( $healthy{$cache} < ( $total_ccrs->{ $cdn_name{$cache} } - ( defined( $bad_ccrs->{ $cdn_name{$cache} } ) ? $bad_ccrs->{ $cdn_name{$cache} } : 0 ) ) ) {
		ERROR "MINORITY REPORT for  " . $cache . " => " . $healthy{$cache} . " out of " . ( $total_ccrs - $bad_ccrs ) . " healthy CCRs think it is OK.";
		$ext->post_result( $server_assoc{$cache}, $check_name, 0 );
	}
}
foreach my $cache ( sort keys %unhealthy ) {
	if ( $unhealthy{$cache} == $total_ccrs->{ $cdn_name{$cache} } ) {
		ERROR $cache . " is marked OFFLINE by all CCRs.";
		$ext->post_result( $server_assoc{$cache}, $check_name, 0 );
	}
}
foreach my $cache ( sort keys %healthy ) {
	if ( $healthy{$cache} == $total_ccrs->{ $cdn_name{$cache} } ) {
		$ext->post_result( $server_assoc{$cache}, $check_name, 1 );
	}
}

sub get_crs_stat() {
	my $ccr  = shift;
	my $port = shift || 80;
	my $url  = 'http://' . $ccr . ":" . $port . '/crs';

	my $ua = LWP::UserAgent->new;
	TRACE "getting " . $url;
	my $response = $ua->get( $url . '/locations' );

	if ( !$response->is_success ) {
		ERROR $ccr . " Not responding! - " . $response->status_line . ".";
		$ccr_assoc{$ccr}->{status} = "DOWN";
		return;
	}
	my $loc_var = JSON->new->utf8->decode( $response->content );

	my $status;
	foreach my $location ( sort @{ $loc_var->{locations} } ) {
		TRACE "getting " . $url . "/locations/" . $location . "/caches";
		my $response = $ua->get( $url . "/locations/" . $location . "/caches" );
		if ( !$response->is_success ) {
			ERROR $ccr . " Not responding! - " . $response->status_line . ".";
			next;
		}
		my $loc_stat = JSON->new->utf8->decode( $response->content );
		$status->{$location} = $loc_stat;
	}
	return $status;
}

