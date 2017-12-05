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

#work in legancy - no tenant  mode
my $useTenancyParamId = &get_param_id('use_tenancy');
ok $t->put_ok('/api/1.2/parameters/' . $useTenancyParamId => {Accept => 'application/json'} => json => {
            'value'      => '0',
        })->status_is(200)
        ->or( sub { diag $t->tx->res->content->asset->{content}; } )
        ->json_is( "/response/name" => "use_tenancy" )
        ->json_is( "/response/configFile" => "global" )
        ->json_is( "/response/value" => "0" )
    , 'Was the disabling paramter set?';

# there is currently 1 delivery service assigned to user with id=200
$t->get_ok('/api/1.2/users/200/deliveryservices')->status_is(200)->$count_response(1)
    ->or( sub { diag $t->tx->res->content->asset->{content}; } );

# there are currently 13 delivery services NOT assigned to user with id=200,
# with tenancy that is accesssible by the user
$t->get_ok('/api/1.2/user/200/deliveryservices/available')->status_is(200)->$count_response(14)
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
# with tenancy that is accesssible by the user
$t->get_ok('/api/1.2/user/200/deliveryservices/available')->status_is(200)->$count_response(13)
    ->or( sub { diag $t->tx->res->content->asset->{content}; } );

# now remove ds=300 from user=200
ok $t->delete_ok('/api/1.2/deliveryservice_user/300/200')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );

# now remove ds=300 from user=200 again, should fail as it is no longer there
ok $t->delete_ok('/api/1.2/deliveryservice_user/200/200')->status_is(404)->or( sub { diag $t->tx->res->content->asset->{content}; } );

# there is now 1 delivery service assigned to user with id=200
$t->get_ok('/api/1.2/users/200/deliveryservices')->status_is(200)->$count_response(1)
    ->or( sub { diag $t->tx->res->content->asset->{content}; } );

# there are now 13 delivery services NOT assigned to user with id=200
# with tenancy that is accesssible by the user
$t->get_ok('/api/1.2/user/200/deliveryservices/available')->status_is(200)->$count_response(14)
    ->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );

################## Tenancy testing - user tenancy point of view
ok $t->post_ok( '/login', => form => { u => Test::TestHelper::ADMIN_ROOT_USER, p => Test::TestHelper::ADMIN_ROOT_USER_PASSWORD } )->status_is(302)
        ->or( sub { diag $t->tx->res->content->asset->{content}; } ), 'Should login?';

#re-enable tenancy
ok $t->put_ok('/api/1.2/parameters/' . $useTenancyParamId => {Accept => 'application/json'} => json => {
            'value'      => '1',
        })->status_is(200)
        ->or( sub { diag $t->tx->res->content->asset->{content}; } )
        ->json_is( "/response/name" => "use_tenancy" )
        ->json_is( "/response/configFile" => "global" )
        ->json_is( "/response/value" => "1" )
    , 'Was the disabling paramter unset?';

my $portal_user_id = $schema->resultset('TmUser')->find( { username => Test::TestHelper::PORTAL_ROOT_USER } )->id;
# there is currently 0 delivery service assigned to PORTAL_ROOT_USER, but the feature is disabled
$t->get_ok('/api/1.2/users/'.$portal_user_id.'/deliveryservices')->status_is(200)->$count_response(15)
    ->or( sub { diag $t->tx->res->content->asset->{content}; } );
# there are currently 15 delivery services NOT assigned to PORTAL_ROOT_USER
$t->get_ok('/api/1.2/user/'.$portal_user_id.'/deliveryservices/available')->status_is(200)->$count_response(15)
    ->or( sub { diag $t->tx->res->content->asset->{content}; } );
ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );


#undef tenant user cannot read the table for the PORTAL_ROOT_USER
ok $t->post_ok( '/login', => form => { u => Test::TestHelper::ADMIN_USER, p => Test::TestHelper::ADMIN_USER_PASSWORD } )->status_is(302)
        ->or( sub { diag $t->tx->res->content->asset->{content}; } ), 'Should login?';
$t->get_ok('/api/1.2/users/'.$portal_user_id.'/deliveryservices')->status_is(403)
    ->or( sub { diag $t->tx->res->content->asset->{content}; } );
$t->get_ok('/api/1.2/user/'.$portal_user_id.'/deliveryservices/available')->status_is(403)
    ->or( sub { diag $t->tx->res->content->asset->{content}; } );
ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );

#root tenant user can read the table for the PORTAL_ROOT_USER
ok $t->post_ok( '/login', => form => { u => Test::TestHelper::ADMIN_ROOT_USER, p => Test::TestHelper::ADMIN_ROOT_USER_PASSWORD } )->status_is(302)
        ->or( sub { diag $t->tx->res->content->asset->{content}; } ), 'Should login?';
# there is currently 0 delivery service assigned to PORTAL_ROOT_USER, but the feature is disabled
$t->get_ok('/api/1.2/users/'.$portal_user_id.'/deliveryservices')->status_is(200)->$count_response(15)
    ->or( sub { diag $t->tx->res->content->asset->{content}; } );
# there are currently 15 delivery services NOT assigned to PORTAL_ROOT_USER
$t->get_ok('/api/1.2/user/'.$portal_user_id.'/deliveryservices/available')->status_is(200)->$count_response(15)
    ->or( sub { diag $t->tx->res->content->asset->{content}; } );
ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );


#undef tenant user cannot assign the ds for the PORTAL_ROOT_USER
ok $t->post_ok( '/login', => form => { u => Test::TestHelper::ADMIN_USER, p => Test::TestHelper::ADMIN_USER_PASSWORD } )->status_is(302)
        ->or( sub { diag $t->tx->res->content->asset->{content}; } ), 'Should login?';
ok $t->post_ok('/api/1.2/deliveryservice_user' => {Accept => 'application/json'} => json => {
            "userId" => $portal_user_id,
            "deliveryServices" => [ 300 ]
        })
        ->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content}; } )
        ->json_is( "/alerts/0/text" => "Invalid user. This user is not available to you for assignment.")
    , 'Does the delivery services assignment blocked?';
ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );

#root tenant user can assign the ds for the PORTAL_ROOT_USER
ok $t->post_ok( '/login', => form => { u => Test::TestHelper::ADMIN_ROOT_USER, p => Test::TestHelper::ADMIN_ROOT_USER_PASSWORD } )->status_is(302)
        ->or( sub { diag $t->tx->res->content->asset->{content}; } ), 'Should login?';
# verifying no change in non-root tenant call, but the feature is disabled
$t->get_ok('/api/1.2/users/'.$portal_user_id.'/deliveryservices')->status_is(200)->$count_response(15)
    ->or( sub { diag $t->tx->res->content->asset->{content}; } );
ok $t->post_ok('/api/1.2/deliveryservice_user' => {Accept => 'application/json'} => json => {
            "userId" => $portal_user_id,
            "deliveryServices" => [ 300 ]
        })
        ->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
        ->json_is( "/response/userId" => $portal_user_id )
        ->json_is( "/response/deliveryServices/0" => 300 )
        ->json_is( "/alerts/0/level" => "success" )
        ->json_is( "/alerts/0/text" => "Delivery service assignments complete." )
    , 'Does the delivery services assign details return?';
ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );


#verify the change
#undef tenant user cannot read the table for the PORTAL_ROOT_USER
ok $t->post_ok( '/login', => form => { u => Test::TestHelper::ADMIN_USER, p => Test::TestHelper::ADMIN_USER_PASSWORD } )->status_is(302)
        ->or( sub { diag $t->tx->res->content->asset->{content}; } ), 'Should login?';
$t->get_ok('/api/1.2/users/'.$portal_user_id.'/deliveryservices')->status_is(403)
    ->or( sub { diag $t->tx->res->content->asset->{content}; } );
$t->get_ok('/api/1.2/user/'.$portal_user_id.'/deliveryservices/available')->status_is(403)
    ->or( sub { diag $t->tx->res->content->asset->{content}; } );
ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );

#root tenant user can read the table for the PORTAL_ROOT_USER
ok $t->post_ok( '/login', => form => { u => Test::TestHelper::ADMIN_ROOT_USER, p => Test::TestHelper::ADMIN_ROOT_USER_PASSWORD } )->status_is(302)
        ->or( sub { diag $t->tx->res->content->asset->{content}; } ), 'Should login?';
# there is currently 1 delivery service assigned to PORTAL_ROOT_USER
$t->get_ok('/api/1.2/users/'.$portal_user_id.'/deliveryservices')->status_is(200)->$count_response(15)
    ->or( sub { diag $t->tx->res->content->asset->{content}; } );
# there are currently 14 delivery services NOT assigned to PORTAL_ROOT_USER
$t->get_ok('/api/1.2/user/'.$portal_user_id.'/deliveryservices/available')->status_is(200)->$count_response(14)
    ->or( sub { diag $t->tx->res->content->asset->{content}; } );
ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );


# now remove ds=300 from PORTAL_ROOT_USER
#undef tenant user cannot delete the ds for the PORTAL_ROOT_USER
ok $t->post_ok( '/login', => form => { u => Test::TestHelper::ADMIN_USER, p => Test::TestHelper::ADMIN_USER_PASSWORD } )->status_is(302)
        ->or( sub { diag $t->tx->res->content->asset->{content}; } ), 'Should login?';
ok $t->delete_ok('/api/1.2/deliveryservice_user/300/'.$portal_user_id)->status_is(403)->or( sub { diag $t->tx->res->content->asset->{content}; } );
ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );

# root tenant user can assign the ds for the PORTAL_ROOT_USER
ok $t->post_ok( '/login', => form => { u => Test::TestHelper::ADMIN_ROOT_USER, p => Test::TestHelper::ADMIN_ROOT_USER_PASSWORD } )->status_is(302)
        ->or( sub { diag $t->tx->res->content->asset->{content}; } ), 'Should login?';
# verifying no change in non-root tenant call, but the feature is disabled
$t->get_ok('/api/1.2/users/'.$portal_user_id.'/deliveryservices')->status_is(200)->$count_response(15)
    ->or( sub { diag $t->tx->res->content->asset->{content}; } );
ok $t->delete_ok('/api/1.2/deliveryservice_user/300/'.$portal_user_id)->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );
$t->get_ok('/api/1.2/users/'.$portal_user_id.'/deliveryservices')->status_is(200)->$count_response(15)
    ->or( sub { diag $t->tx->res->content->asset->{content}; } );
ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );


$dbh->disconnect();
done_testing();

sub get_param_id {
    my $name = shift;
    my $q      = "select id from parameter where name = \'$name\'";
    my $get_svr = $dbh->prepare($q);
    $get_svr->execute();
    my $p = $get_svr->fetchall_arrayref( {} );
    $get_svr->finish();
    my $id = $p->[0]->{id};
    return $id;
}
