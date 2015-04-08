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

use strict;
use warnings;

$|++;

my $username = "jvando001";

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

my $b_url = $jconf->{base_url};
Extensions::Helper->import();
my $ext = Extensions::Helper->new( { base_url => $b_url, token => '91504CE6-8E4A-46B2-9F9F-FE7C15228498' } );

my $jdataserver = $ext->get(Extensions::Helper::SERVERLIST_PATH);
my $match       = $jconf->{match};

my $select     = $jconf->{select};
my $check_name = $jconf->{check_name};
if ( $check_name ne "ORT" ) {
	ERROR "This Check Extension is exclusively for the ORT (Operational Readiness Test) check.";
	exit(4);
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

	my $cmd = "/usr/bin/sudo /opt/ort/ipcdn_install_ort.pl report WARN";
	$cmd = "ssh -o \"StrictHostKeyChecking no\" -i ~$username/.ssh/id_dsa -l " . $username . " " . $ipaddr . " " . $cmd;
	TRACE $host_name . " running " . $cmd;
	my $out = `$cmd`;

	my $results = "/var/www/html/ort/" . $host_id;

	if ( !-d $results ) {
		mkdir($results);
	}

	my @lines = split( /\n/, $out );
	my $count = grep /^ERROR/, @lines;
	TRACE $host_name . " score: " . $count;
	$ext->post_result( $host_id, $check_name, $count );
}

