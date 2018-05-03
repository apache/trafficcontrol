package main;
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
use Mojo::Base -strict;
use Test::More;
use Test::Mojo;
use Data::Dumper;
use DBI;
use Test::TestHelper;
use Test::MockModule;
use Test::MockObject;
use strict;
use warnings;
use JSON;

BEGIN { $ENV{MOJO_MODE} = "integration" }
my $t = Test::Mojo->new('TrafficOps');
no warnings 'once';
use warnings 'all';
my $api_version = '1.1';

my $fake_lwp = new Test::MockModule( 'LWP::UserAgent', no_auto => 1 );
my $fake_get = HTTP::Response->new( 200, undef, HTTP::Headers->new, encode_json({response => "success"}) );
$fake_lwp->mock( 'get', sub { return $fake_get } );
my $fake_put = HTTP::Response->new( 204, undef, HTTP::Headers->new, undef );
$fake_lwp->mock( 'put', sub { return $fake_put } );
my $fake_delete = HTTP::Response->new( 204, undef, HTTP::Headers->new, undef );
$fake_lwp->mock( 'delete', sub { return $fake_delete } );


$t->post_ok( '/login', => form => { u => 'admin', p => 'password' } )->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );

my $json = JSON->new->allow_nonref;
$t->get_ok('/dataprofile')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );
my $profile_arr  = $json->decode( $t->tx->res->content->asset->{content} );
my %profile_name = ();
my $i            = 0;
foreach my $p ( @{$profile_arr} ) {
	$profile_name{ $p->{id} } = $p->{name};
	$i++;
}

$t->get_ok( '/api/' . $api_version . '/servers.json' )->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );
my $servers = $json->decode( $t->tx->res->content->asset->{content} );

$i = 0;
foreach my $server ( @{ $servers->{response} } ) {
	$t->get_ok( '/ort/' . $server->{hostName} . '/ort1' )->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );
	my $files = $json->decode( $t->tx->res->content->asset->{content} );
	diag "... server: " . $server->{hostName} . ' ' . $files->{other}->{CDN_name} . ' ' . $files->{profile}->{name};
	if ( $server->{type} eq 'EDGE' || $server->{type} eq 'MID' ) {
		foreach my $file ( keys %{ $files->{config_files} } ) {
			my $path = '/genfiles/view/' . $server->{hostName} . '/' . $file;
			$t->get_ok($path)->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );
		}
	}
}

done_testing();
