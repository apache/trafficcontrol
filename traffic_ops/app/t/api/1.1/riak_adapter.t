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
use DBI;
use strict;
use warnings;
no warnings 'once';
use warnings 'all';
use Test::TestHelper;
use Test::MockModule;
use Connection::RiakAdapter;
use Data::Dumper;

use constant BUCKET => "mybucket";
use constant KEY    => "mybucket";

#no_transactions=>1 ==> keep fixtures after every execution, beware of duplicate data!
#no_transactions=>0 ==> delete fixtures after every execution

#GET
my $riak_util = Connection::RiakAdapter->new( Test::TestHelper::ADMIN_USER, Test::TestHelper::ADMIN_USER_PASSWORD );

my $fake_answer = "OK";
my $fake_lwp    = new Test::MockModule( 'LWP::UserAgent', no_auto => 1 );
my $fake_header = HTTP::Headers->new;
$fake_header->header( 'Content-Type' => 'application/json' );    # set

my $fake_response = HTTP::Response->new( 200, undef, $fake_header, $fake_answer );
$fake_lwp->mock( 'delete', sub { return $fake_response } );
$fake_lwp->mock( 'get',    sub { return $fake_response } );
$fake_lwp->mock( 'put',    sub { return $fake_response } );

my $key_uri = $riak_util->get_key_uri( BUCKET, KEY );
is( $key_uri, "/riak/" . BUCKET . "/" . KEY );

#GET
my $response = $riak_util->get( BUCKET, KEY );
is( $response->{_rc},                        200 );
is( $response->{_headers}->{'content-type'}, 'application/json' );
is( $fake_answer, $response->{_content} );

#PUT
$response = $riak_util->put( BUCKET, KEY, "value1" );
is( $response->{_rc}, 200 );
is( $fake_answer, $response->{_content} );

#PUT - With Content Type
$response = $riak_util->put( BUCKET, KEY, "value1", "application/json" );
is( $response->{_rc}, 200 );
is( $fake_answer,     $response->{_content} );

#DELETE
$response = $riak_util->delete( BUCKET, KEY );
is( $response->{_rc}, 200 );
is( $fake_answer, $response->{_content} );

#PING
$response = $riak_util->ping();
is( $response->{_rc}, 200 );
is( $fake_answer,     $response->{_content} );

#STATS
$response = $riak_util->stats();
is( $response->{_rc}, 200 );
is( $fake_answer,     $response->{_content} );

done_testing();
