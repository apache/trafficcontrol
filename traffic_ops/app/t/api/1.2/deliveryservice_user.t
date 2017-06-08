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
#    ->json_is( "/response/hostname" => "http://10.75.168.91")
#    ->json_is( "/response/ds_assigned/1" => "ds1")
#            , 'Does the deliveryservice details return?';
use Mojo::Base -strict;
use Test::More;
use Test::Mojo;
use DBI;
use JSON;
use strict;
use warnings;
no warnings 'once';
use warnings 'all';
use Test::TestHelper;

#no_transactions=>1 ==> keep fixtures after every execution, beware of duplicate data!
#no_transactions=>0 ==> delete fixtures after every execution

BEGIN { $ENV{MOJO_MODE} = "test" }

my $schema = Schema->connect_to_database;
my $dbh    = Schema->database_handle;
my $t      = Test::Mojo->new('TrafficOps');

Test::TestHelper->unload_core_data($schema);
Test::TestHelper->load_core_data($schema);

ok $t->post_ok( '/login', => form => { u => Test::TestHelper::ADMIN_USER, p => Test::TestHelper::ADMIN_USER_PASSWORD } )->status_is(302)
        ->or( sub { diag $t->tx->res->content->asset->{content}; } ), 'Should login?';

# Count the 'response number'
my $count_response = sub {
    my ( $t, $count ) = @_;
    my $json = decode_json( $t->tx->res->content->asset->slurp );
    my $r    = $json->{response};
    return $t->success( is( scalar(@$r), $count ) );
};

# there is currently 1 delivery service assigned to user with id=200
$t->get_ok('/api/1.2/users/200/deliveryservices')->status_is(200)->$count_response(1)
    ->or( sub { diag $t->tx->res->content->asset->{content}; } );

# there are currently 12 delivery services NOT assigned to user with id=200
$t->get_ok('/api/1.2/user/200/deliveryservices/available')->status_is(200)->$count_response(12)
    ->or( sub { diag $t->tx->res->content->asset->{content}; } );

# assign one ds to user with id=200
ok $t->post_ok('/api/1.2/deliveryservice_user' => {Accept => 'application/json'} => json => {
            "userId" => 200,
            "deliveryServices" => [ 300 ]
        })
        ->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
        ->json_is( "/response/userId" => 200 )
        ->json_is( "/response/deliveryServices/0" => 300 )
        ->json_is( "/alerts/0/level" => "success" )
        ->json_is( "/alerts/0/text" => "Delivery service assignments complete." )
    , 'Does the delivery services assign details return?';

# there are now 2 delivery services assigned to user with id=200
$t->get_ok('/api/1.2/users/200/deliveryservices')->status_is(200)->$count_response(2)
    ->or( sub { diag $t->tx->res->content->asset->{content}; } );

# there are now 11 delivery services NOT assigned to user with id=200
$t->get_ok('/api/1.2/user/200/deliveryservices/available')->status_is(200)->$count_response(11)
    ->or( sub { diag $t->tx->res->content->asset->{content}; } );

# now remove ds=300 from user=200
ok $t->delete_ok('/api/1.2/deliveryservice_user/300/200')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );

# now remove ds=300 from user=200 again, should fail as it is no longer there
ok $t->delete_ok('/api/1.2/deliveryservice_user/200/200')->status_is(404)->or( sub { diag $t->tx->res->content->asset->{content}; } );

# there is now 1 delivery service assigned to user with id=200
$t->get_ok('/api/1.2/users/200/deliveryservices')->status_is(200)->$count_response(1)
    ->or( sub { diag $t->tx->res->content->asset->{content}; } );

# there are now 12 delivery services NOT assigned to user with id=200
$t->get_ok('/api/1.2/user/200/deliveryservices/available')->status_is(200)->$count_response(12)
    ->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );
$dbh->disconnect();
done_testing();
