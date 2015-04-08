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
# DSCP check extension. Populates the 'DSCP' column.
#

use strict;
use warnings;

$|++;

use Data::Dumper;
use Getopt::Std;
use Log::Log4perl qw(:easy);
use Net::PcapUtils;
use NetPacket::Ethernet qw(:strip);
use NetPacket::IP qw(:strip);
use NetPacket::TCP;
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

TRACE Dumper($jconf);
my $b_url = $jconf->{base_url};
Extensions::Helper->import();
my $ext = Extensions::Helper->new( { base_url => $b_url, token => '91504CE6-8E4A-46B2-9F9F-FE7C15228498' } );

my $match = $jconf->{match};

my $check_name = $jconf->{check_name};
if ( $check_name ne "DSCP" ) {
	ERROR "This Check Extension is exclusively for DSCP.";
	exit(4);
}

my %ds_info           = ();
my $jdeliveryservices = $ext->get( Extensions::Helper::DSLIST_PATH );

foreach my $ds ( @{$jdeliveryservices} ) {
	$ds_info{ $ds->{id} } = $ds;
}

my %domain_name_for_profile = ();
my $jdataserver             = $ext->get( Extensions::Helper::SERVERLIST_PATH );
foreach my $server ( @{$jdataserver} ) {
	next unless $server->{type} eq 'EDGE';    # We know this is DSCP, so we know we want edges only
	my $ip        = $server->{ipAddress};
	my $host_name = $server->{hostName};
	my $details   = $ext->get( '/api/1.1/servers/hostname/' . $host_name . '/details.json' );
	foreach my $dsid ( @{ $details->{deliveryservices} } ) {
		my $ds = $ds_info{$dsid};
		if ( $ds->{active} && defined( $ds->{checkPath} ) && $ds->{checkPath} ne "" && $ds->{protocol} == 0 ) {
			my $prefix = $host_name;
			if ( $ds->{type} =~ /^DNS/ ) {
				$prefix = 'edge';
			}
			my $url = 'http://' . $prefix;
			foreach my $match ( @{ $ds->{matchList} } ) {
				if ( $match->{type} eq 'HOST_REGEXP' ) {
					$url .= $match->{pattern};
					$url =~ s/\\//g;
					$url =~ s/\.\*//g;
				}
			}
			if ( !defined( $domain_name_for_profile{ $ds->{profileName} } ) ) {
				my $param_list = $ext->get( '/api/1.1/parameters/profile/' . $ds->{profileName} . '.json' );    ## TODO: create /api, use that
				foreach my $p ( @{$param_list} ) {
					if ( $p->{name} eq 'domain_name' ) {
						$domain_name_for_profile{ $ds->{profileName} } = $p->{value};
					}
				}
			}
			$url .= $domain_name_for_profile{ $ds->{profileName} };
			$url .= $ds->{checkPath};

			my $dscp_found = &get_dscp( $url, $server->{ipAddress}, "p3p1" );
			if ( $dscp_found == $ds->{dscp} ) {
				TRACE "Success";
				$ext->post_result( $server->{id}, $check_name, 1 );
			} else {
				TRACE "Fail";
				$ext->post_result( $server->{id}, $check_name, 0 );
			}
			last;
		}
	}
}

sub get_dscp() {
	my $url = shift;
	my $ip  = shift;
	my $dev = shift;

	my $tos     = 0;
	my $max_len = 0;

	my $src_port = int( rand( 65535 - 1024 ) ) + 1024;
	TRACE "get_dscp ip:" . $ip . " url:" . $url . " dev:" . $dev . " port:" . $src_port . "\n";

	# Use curl to get some traffic from the URL, but send the command to the background, so the capture that follows
	# is while traffic is being returned
	system( "(sleep 1; curl --local-port " . $src_port . " --ipv4 -s $url 2>&1 > /dev/null || ping -c 10 $ip 2>&1 > /dev/null)  &" );

	Net::PcapUtils::loop(
		sub {
			my ( $user, $hdr, $pkt ) = @_;
			my $ip_obj = NetPacket::IP->decode( eth_strip($pkt) );

			TRACE " <=> $ip_obj->{src_ip} -> $ip_obj->{dest_ip} $ip_obj->{proto} tos $ip_obj->{tos} len $ip_obj->{len}\n";
			my $tcp_obj = NetPacket::TCP->decode( $ip_obj->{data} );
			TRACE " TCP1 $ip_obj->{src_ip}:$tcp_obj->{src_port} -> $ip_obj->{dest_ip}:$tcp_obj->{dest_port} $ip_obj->{proto} tos $ip_obj->{tos} len $ip_obj->{len}\n";
			if ( $ip_obj->{src_ip} eq $ip && $ip_obj->{len} > $max_len && $ip_obj->{proto} == 6 ) {
				my $tcp_obj = NetPacket::TCP->decode( $ip_obj->{data} );
				TRACE " TCP2 $ip_obj->{src_ip}:$tcp_obj->{src_port} -> $ip_obj->{dest_ip}:$tcp_obj->{dest_port} $ip_obj->{proto} tos $ip_obj->{tos} len $ip_obj->{len}\n";
				if ( $tcp_obj->{src_port} == 80 && $tcp_obj->{dest_port} == $src_port ) {
					TRACE " TCP3 $ip_obj->{src_ip}:$tcp_obj->{src_port} -> $ip_obj->{dest_ip}:$tcp_obj->{dest_port} $ip_obj->{proto} tos $ip_obj->{tos} len $ip_obj->{len}\n";
					$max_len = $ip_obj->{len};
					$tos     = $ip_obj->{tos};
				}
			}
		},
		FILTER     => 'host ' . $ip,
		DEV        => $dev,
		NUMPACKETS => 7,
		TIMEOUT    => 10
	);

	my $dscp = $tos >> 2;
	TRACE "returning " . $dscp;
	return $dscp;
}

