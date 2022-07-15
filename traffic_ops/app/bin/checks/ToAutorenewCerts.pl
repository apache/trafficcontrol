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
# DNSSEC refresh, checks to see if DNSSEC keys need to be re-generated.
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

my $to_user= $jconf->{user};
my $to_pass= $jconf->{pass};

if ( !defined($to_user) || $to_user eq '' ) {
	ERROR "Config missing \"user\" key, this script now requires \"user\" and \"pass\" keys, as the endpoint used by this script now requires admin-level authentication.";
	exit(1);
}

if ( !defined($to_pass) || $to_pass eq '' ) {
	ERROR "Config missing \"pass\" key, this script now requires \"user\" and \"pass\" keys, as the endpoint used by this script now requires admin-level authentication.";
	exit(1);
}

my $ua = LWP::UserAgent->new;
$ua->timeout(30);
$ua->ssl_opts(verify_hostname => 0);
$ua->cookie_jar( {} );

my $login_url = "$b_url/api/4.0/user/login";
TRACE "posting $login_url";
my $req = HTTP::Request->new( 'POST', $login_url );
$req->header( 'Content-Type' => 'application/json' );

$req->content( "{\"u\": \"$to_user\",\"p\": \"$to_pass\"}" );
my $login_response = $ua->request($req);
if ( ! $login_response->is_success ) {
	ERROR "Error trying to update keys, login failed, response was " . $login_response->status_line;
	exit(1);
}

my $url       = "$b_url/api/4.0/letsencrypt/autorenew/";
TRACE "getting $url";
my $response = $ua->post($url);
if ( $response->is_success ) {
	DEBUG "Successfully refreshed keys response was " . $response->decoded_content;
}
else {
 ERROR "Error trying to update keys, response was " . $response->status_line;
}

sub help() {
	print
		"ToAutorenewCerts.pl -c \"{\\\"base_url\\\": \\\"https://localhost\\\", \\\"user\\\": \\\"user\\\", \\\"pass\\\": \\\"password\\\"}\"\n";
	print "\n";
	print "-c   json formatted list of variables\n";
	print "     base_url: required\n";
	print "        URL of the Traffic Ops server.\n";
	print "     user: required\n";
	print "        The Traffic Ops user.\n";
	print "     pass: required\n";
	print "        The password for the user.\n";
	print "-l   Logging level. 1 - 6. 1 being least (FATAL). 6 being most (TRACE). Default\n";
    print "     is 1.\n";
	print "================================================================================\n";

	# the above line of equal signs is 80 columns
	print "\n";
}
