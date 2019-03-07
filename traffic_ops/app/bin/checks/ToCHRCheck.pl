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
# CHR check extension. Populates the 'CHR' (Cache Hit Ratio) column.
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
my $tmpdir = '/tmp/gggstats';

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

# JvD TESTING
#local $/;
#open( my $fh, '<', '/tmp/servers.json' );
#my $json_text   = <$fh>;
#$jdataserver = decode_json( $json_text );
# JvD TESTING


my $match       = $jconf->{match};

my $select     = $jconf->{select};
my $check_name = $jconf->{check_name};
if ( $check_name ne "CHR" ) {
	ERROR "This Check Extension is exclusively for CHR (Cache Hit Ratio).";
	exit(4);
}

my $ua = LWP::UserAgent->new;
$ua->timeout(3);
my $i                  = 0;
my $check_against_prev = 1;
if ( !-d $tmpdir ) {
	mkdir($tmpdir);
	$check_against_prev = 0;
}

foreach my $server ( @{$jdataserver} ) {
	if ( $server->{type} =~ m/^EDGE/ || $server->{type} =~ m/^MID/ ) {    # We know this is "CHR, so we know what we want
		my $ip        = $server->{ipAddress};
		my $host_name = $server->{hostName};
		my $interface = $server->{interfaceName};
		my $port      = $server->{tcpPort};
		my $url       = 'http://' . $ip . ':' . $port . '/_astats?application=proxy.process.http.transaction_counts&inf.name=' . $interface;
		TRACE "getting $url";
		my $response = $ua->get($url);
		if ( $response->is_success ) {
			my $stats_var = JSON->new->utf8->decode( $response->content );
			my $hits =
				  $stats_var->{'ats'}{'proxy.process.http.transaction_counts.hit_fresh'}
				+ $stats_var->{'ats'}{'proxy.process.http.transaction_counts.hit_fresh.process'}
				+ $stats_var->{'ats'}{'proxy.process.http.transaction_counts.hit_revalidated'};
			my $miss =
				  $stats_var->{'ats'}{'proxy.process.http.transaction_counts.miss_cold'}
				+ $stats_var->{'ats'}{'proxy.process.http.transaction_counts.miss_not_cacheable'}
				+ $stats_var->{'ats'}{'proxy.process.http.transaction_counts.miss_changed'}
				+ $stats_var->{'ats'}{'proxy.process.http.transaction_counts.miss_client_no_cache'};
			my $errors =
				  $stats_var->{'ats'}{'proxy.process.http.transaction_counts.errors.aborts'}
				+ $stats_var->{'ats'}{'proxy.process.http.transaction_counts.errors.possible_aborts'}
				+ $stats_var->{'ats'}{'proxy.process.http.transaction_counts.errors.connect_failed'}
				+ $stats_var->{'ats'}{'proxy.process.http.transaction_counts.errors.other'}
				+ $stats_var->{'ats'}{'proxy.process.http.transaction_counts.other.unclassified'};

			my $filename = $tmpdir . "/" . $host_name . ".stats";
			if ($check_against_prev) {
				my $ftime = ( stat $filename )[9];
				my $secs  = time() - $ftime;
				if ( -f $filename ) {
					open( FILE, "<$filename" );
					my $jstring;
					sysread( FILE, $jstring, -s $filename );
					close(FILE);
					my $prev_var = JSON->new->utf8->decode($jstring);
					my $prev_hits =
						  $prev_var->{'ats'}{'proxy.process.http.transaction_counts.hit_fresh'}
						+ $prev_var->{'ats'}{'proxy.process.http.transaction_counts.hit_fresh.process'}
						+ $prev_var->{'ats'}{'proxy.process.http.transaction_counts.hit_revalidated'};
					my $prev_miss =
						  $prev_var->{'ats'}{'proxy.process.http.transaction_counts.miss_cold'}
						+ $prev_var->{'ats'}{'proxy.process.http.transaction_counts.miss_not_cacheable'}
						+ $prev_var->{'ats'}{'proxy.process.http.transaction_counts.miss_changed'}
						+ $prev_var->{'ats'}{'proxy.process.http.transaction_counts.miss_client_no_cache'};
					my $prev_errors =
						  $prev_var->{'ats'}{'proxy.process.http.transaction_counts.errors.aborts'}
						+ $prev_var->{'ats'}{'proxy.process.http.transaction_counts.errors.possible_aborts'}
						+ $prev_var->{'ats'}{'proxy.process.http.transaction_counts.errors.connect_failed'}
						+ $prev_var->{'ats'}{'proxy.process.http.transaction_counts.errors.other'}
						+ $prev_var->{'ats'}{'proxy.process.http.transaction_counts.other.unclassified'};

					TRACE " prev_hits == $prev_hits misses == $prev_miss errors == $prev_errors\n";
					my $h  = $hits - $prev_hits;
					my $m  = $miss - $prev_miss;
					my $e  = $errors - $prev_errors;
					my $t  = ( $h + $m + $e );
					if ( $t != 0 ) {
						my $hr = sprintf( "%d", ($h / $t) * 100 );
						TRACE "$host_name: hitratio: $hr\%, errors: $e, period: $secs\n";
						$ext->post_result( $server->{id}, $check_name, $hr );
				  } else {
						TRACE "$host_name: No transaction_counts! Service enabled?"
					}
				}
			}
			else {
				$ext->post_result( $server->{id}, $check_name, 0 );
			}
			open( FILE, ">$filename" ) || die $!;
			print FILE $response->content;
			close(FILE);
		}
		else {
			ERROR "Can't get _astats for " . $ip;
			$ext->post_result( $server->{id}, $check_name, -1 );
		}

	}
}
