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
# CDU check extension. Populates the 'CDU' (Cache Disk Usage) column.
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

TRACE Dumper($jconf);
my $b_url = $jconf->{base_url};
Extensions::Helper->import();
my $ext = Extensions::Helper->new( { base_url => $b_url, token => '91504CE6-8E4A-46B2-9F9F-FE7C15228498' } );

my $jdataserver = $ext->get( Extensions::Helper::SERVERLIST_PATH );

my $check_name = $jconf->{check_name};
if ( $check_name ne "CDU" ) {
	ERROR "This Check Extension is exclusively for CDU (Cache Disk Usage).";
	exit(4);
}

my $ua = LWP::UserAgent->new;
$ua->timeout(3);

foreach my $server ( @{$jdataserver} ) {
	if ( $server->{type} =~ m/^EDGE/ || $server->{type} =~ m/^MID/ ) {    # We know this is "CHR, so we know what we want
		my $ip        = $server->{ipAddress};
		my $host_name = $server->{hostName};
		my $interface = $server->{interfaceName};
		my $port      = $server->{tcpPort};
		my $url       = 'http://' . $ip . ':' . $port . '/_astats?application=bytes_used;bytes_total&inf.name=' . $interface;
		TRACE "getting $url";
		my $response = $ua->get($url);
		if ( $response->is_success ) {
			my $stats_var = JSON->new->utf8->decode( $response->content );

			# TODO: add ability to display volume 1 and volume2 stats - it's all there.
			# For now, just replace the old CDU check.
			# "proxy.process.cache.bytes_used": 21591297110528,
			# "proxy.process.cache.bytes_total": 21655715577856,
			# "proxy.process.cache.ram_cache.bytes_used": 33849232896,
			# "proxy.process.cache.volume_1.bytes_used": 21586800279552,
			# "proxy.process.cache.volume_1.bytes_total": 21587001606144,
			# "proxy.process.cache.volume_1.ram_cache.bytes_used": 33746226688,
			# "proxy.process.cache.volume_2.bytes_used": 4496830976,
			# "proxy.process.cache.volume_2.bytes_total": 68713971712,
			# "proxy.process.cache.volume_2.ram_cache.bytes_used": 103006208,

			my $size                  = $stats_var->{ats}{'proxy.process.cache.bytes_total'};
			my $used                  = $stats_var->{ats}{'proxy.process.cache.bytes_used'};
			if ( $size == 0 ) {
				ERROR "$host_name: cache size is 0!";
				next;
			}
			my $percentage_cache_used = sprintf( "%3d", ( $used / $size ) * 100 );
			TRACE "$host_name: percentage cache used == " . $percentage_cache_used;
			$ext->post_result( $server->{id}, $check_name, $percentage_cache_used );
		}
		else {
			ERROR "Can't get _astats for " . $ip;
			$ext->post_result( $server->{id}, $check_name, -1 );
		}

	}
}
