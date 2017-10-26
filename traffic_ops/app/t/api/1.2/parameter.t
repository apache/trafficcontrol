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

#no_transactions=>1 ==> keep fixtures after every execution, beware of duplicate data!
#no_transactions=>0 ==> delete fixtures after every execution

BEGIN { $ENV{MOJO_MODE} = "test" }

my $schema = Schema->connect_to_database;
my $schema_values = { schema => $schema, no_transactions => 1 };
my $dbh    = Schema->database_handle;
my $t      = Test::Mojo->new('TrafficOps');

Test::TestHelper->unload_core_data($schema);

# Load the test data up until 'cachegroup', because this test case creates
# them.
Test::TestHelper->load_all_fixtures( Fixtures::Tenant->new($schema_values) );
Test::TestHelper->load_all_fixtures( Fixtures::Cdn->new($schema_values) );
Test::TestHelper->load_all_fixtures( Fixtures::Role->new($schema_values) );
Test::TestHelper->load_all_fixtures( Fixtures::TmUser->new($schema_values) );
Test::TestHelper->load_all_fixtures( Fixtures::Status->new($schema_values) );
Test::TestHelper->load_all_fixtures( Fixtures::Type->new($schema_values) );
Test::TestHelper->load_all_fixtures( Fixtures::Profile->new($schema_values) );

ok $t->post_ok( '/login', => form => { u => Test::TestHelper::ADMIN_USER, p => Test::TestHelper::ADMIN_USER_PASSWORD } )->status_is(302)
	->or( sub { diag $t->tx->res->content->asset->{content}; } ), 'Should login?';

ok $t->post_ok('/api/1.2/parameters' => {Accept => 'application/json'} => json => 
        {
            'name'  => 'param10',
            'configFile' => 'configFile10',
            'value'      => 'value10',
            'secure'     => '0'
        }
    )->status_is(200)
	->or( sub { diag $t->tx->res->content->asset->{content}; } )
    ->json_is( "/response/0/name" => "param10" )
    ->json_is( "/response/0/configFile" => "configFile10" )
    ->json_is( "/response/0/value" => "value10" )
    ->json_is( "/response/0/secure" => "0" )
		, 'Does the paramters created return?';

ok $t->post_ok('/api/1.2/parameters' => {Accept => 'application/json'} => json => 
	[
        {
            'name'  => 'param1',
            'configFile' => 'configFile1',
            'value'      => 'value1',
            'secure'     => '0'
        },
        {
            'name'  => 'param2',
            'configFile' => 'configFile2',
            'value'      => 'value2',
            'secure'     => '1'
        }
    ])->status_is(200)
	->or( sub { diag $t->tx->res->content->asset->{content}; } )
    ->json_is( "/response/0/name" => "param1" )
    ->json_is( "/response/0/configFile" => "configFile1" )
    ->json_is( "/response/0/value" => "value1" )
    ->json_is( "/response/0/secure" => "0" )
    ->json_is( "/response/1/name" => "param2" )
    ->json_is( "/response/1/configFile" => "configFile2" )
    ->json_is( "/response/1/value" => "value2" )
    ->json_is( "/response/1/secure" => "1" )
		, 'Does the paramters created return?';

ok $t->post_ok('/api/1.2/parameters' => {Accept => 'application/json'} => json => [
        {
            'name'  => 'param3',
            'configFile' => 'configFile3',
            'value'      => 'value3',
            'secure'     => '0'
        },
        {
             name        => 'domain_name',
             value       => 'foo.com',
             configFile  => 'CRConfig.json',
            'secure'     => '0'
        }
    ])->status_is(200)
	->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->post_ok('/api/1.2/parameters' => {Accept => 'application/json'} => json => [
        {
            'name'  => 'param3',
            'configFile' => 'configFile3',
            'value'      => 'value3',
            'secure'     => '0'
        },
        {
             name        => 'domain_name',
             value       => 'foo.com',
             configFile  => 'CRConfig.json',
            'secure'     => '0'
        }
    ])->status_is(400)
	->or( sub { diag $t->tx->res->content->asset->{content}; } )
	->json_is( "/alerts/0/text" => "parameter [name:param3 , configFile:configFile3 , value:value3] already exists." )
		, 'Does the paramters created return?';

ok $t->post_ok('/api/1.2/parameters' => {Accept => 'application/json'} => json => [
        {
            'name'  => 'param3',
            'configFile' => 'configFile3',
            'value'      => 'value3',
            'secure'     => '0'
        },
        {
            'name'  => 'param3',
             configFile  => 'CRConfig.json',
            'secure'     => '0'
        }
    ])->status_is(400)
	->or( sub { diag $t->tx->res->content->asset->{content}; } )
	->json_is( "/alerts/0/text" => 'parameter [name:param3 , configFile:configFile3 , value:value3] already exists.' )
		, 'Does the paramters create return?';

ok $t->post_ok('/api/1.2/parameters' => {Accept => 'application/json'} => json => [
        {
            'name'  => 'param3',
            'configFile' => 'configFile3',
            'value'      => 'value3',
            'secure'     => '0'
        },
        {
            'secure'     => '0'
        }
    ])->status_is(400)
	->or( sub { diag $t->tx->res->content->asset->{content}; } )
		, 'Does the paramters created return?';

my $para_id = &get_param_id('param2');

ok $t->put_ok('/api/1.2/parameters/' . $para_id => {Accept => 'application/json'} => json => {
            'value'      => 'value2.1',
            'secure'     => '0'
    })->status_is(200)
	->or( sub { diag $t->tx->res->content->asset->{content}; } )
	->json_is( "/response/name" => "param2" )
	->json_is( "/response/configFile" => "configFile2" )
	->json_is( "/response/value" => "value2.1" )
	->json_is( "/response/secure" => "0" )
		, 'Does the paramters modified return?';

ok $t->put_ok('/api/1.2/parameters/' . $para_id => {Accept => 'application/json'} => json => {
            'name'  => 'param2.1',
            'configFile' => 'configFile2.1',
            'secure'     => '1'
    })->status_is(200)
	->or( sub { diag $t->tx->res->content->asset->{content}; } )
	->json_is( "/response/name" => "param2.1" )
	->json_is( "/response/configFile" => "configFile2.1" )
	->json_is( "/response/value" => "value2.1" )
	->json_is( "/response/secure" => "1" )
		, 'Does the paramters modified return?';

ok $t->put_ok('/api/1.2/parameters/0' => {Accept => 'application/json'} => json => {
    })->status_is(404)
	->or( sub { diag $t->tx->res->content->asset->{content}; } )
		, 'Does the paramters modified return?';

ok $t->delete_ok('/api/1.2/parameters/' . $para_id )->status_is(200)
	->or( sub { diag $t->tx->res->content->asset->{content}; } )
		, 'Does the parameter delete return?';

$para_id = &get_param_id('param10');
ok $t->post_ok('/api/1.2/profileparameters' => {Accept => 'application/json'} => json => {
	"profileId" => 300, "parameterId" => $para_id })->status_is(200)
	->or( sub { diag $t->tx->res->content->asset->{content}; } )
	->json_is( "/response/0/profileId" => "300" )
	->json_is( "/response/0/parameterId" => $para_id )
		, 'Does the profile parameter details return?';

ok $t->delete_ok('/api/1.2/parameters/' . $para_id )->status_is(400)
	->or( sub { diag $t->tx->res->content->asset->{content}; } )
	->json_like( "/alerts/0/text" => qr/has profile associated/ )
		, 'Does the parameter delete return?';

ok $t->post_ok('/api/1.2/parameters/validate' => {Accept => 'application/json'} => json => {
            'name'  => 'param1',
            'configFile' => 'configFile1',
            'value'      => 'value1'
    })->status_is(200)
	->or( sub { diag $t->tx->res->content->asset->{content}; } )
    ->json_like( "/response/id" => qr/^\d+$/ )
	->json_is( "/response/name" => "param1" )
	->json_is( "/response/configFile" => "configFile1" )
	->json_is( "/response/value" => "value1" )
	->json_is( "/response/secure" => "0" )
		, 'Does the paramters validate return?';

ok $t->post_ok('/api/1.2/parameters/validate' => {Accept => 'application/json'} => json => {
            'configFile' => 'configFile1',
            'value'      => 'value1'
    })->status_is(400)
	->or( sub { diag $t->tx->res->content->asset->{content}; } )
	->json_like( "/alerts/0/text" => qr/is required.$/ )
	->or( sub { diag $t->tx->res->content->asset->{content}; } )
		, 'Does the paramters validate return?';
ok $t->post_ok('/api/1.2/parameters/validate' => {Accept => 'application/json'} => json => {
            'name'  => 'param1',
            'value'      => 'value1'
    })->status_is(400)
	->or( sub { diag $t->tx->res->content->asset->{content}; } )
	->json_like( "/alerts/0/text" => qr/is required.$/ )
	->or( sub { diag $t->tx->res->content->asset->{content}; } )
		, 'Does the paramters validate return?';
ok $t->post_ok('/api/1.2/parameters/validate' => {Accept => 'application/json'} => json => {
            'name'  => 'param1',
            'configFile' => 'configFile1',
    })->status_is(400)
	->or( sub { diag $t->tx->res->content->asset->{content}; } )
	->json_like( "/alerts/0/text" => qr/is required.$/ )
	->or( sub { diag $t->tx->res->content->asset->{content}; } )
		, 'Does the paramters validate return?';
ok $t->post_ok('/api/1.2/parameters/validate' => {Accept => 'application/json'} => json => {
            'name'  => 'noexist',
            'configFile' => 'noexist',
            'value'      => 'noexist'
    })->status_is(400)
	->or( sub { diag $t->tx->res->content->asset->{content}; } )
	->json_like( "/alerts/0/text" => qr/does not exist.$/ )
	->or( sub { diag $t->tx->res->content->asset->{content}; } )
		, 'Does the paramters validate return?';


#checking if a parameter vaule can be changed to "0"
ok $t->post_ok('/api/1.2/parameters' => {Accept => 'application/json'} => json => [
			{
				'name'  => 'default1',
				'configFile' => 'configFile3',
				'value'      => '1',
				'secure'     => '0'
			}]
	)->status_is(200)
	, 'Adding the parameter with default 1';

$para_id = &get_param_id('default1');
ok $t->get_ok('/api/1.2/parameters/'. $para_id)->status_is(200)
		->or( sub { diag $t->tx->res->content->asset->{content}; } )
		->json_is( "/response/0/name" => "default1" )
		->json_is( "/response/0/value" => "1" )
		->json_is( "/response/0/configFile" => "configFile3" )
	, 'Does the paramter get return?';

ok $t->put_ok('/api/1.2/parameters/' . $para_id => {Accept => 'application/json'} => json => {
			'value'      => '0',
		})->status_is(200)
		->or( sub { diag $t->tx->res->content->asset->{content}; } )
		->json_is( "/response/name" => "default1" )
		->json_is( "/response/configFile" => "configFile3" )
		->json_is( "/response/value" => "0" )
	, 'Was the paramters modification return?';

ok $t->get_ok('/api/1.2/parameters/'. $para_id)->status_is(200)
		->or( sub { diag $t->tx->res->content->asset->{content}; } )
		->json_is( "/response/0/name" => "default1" )
		->json_is( "/response/0/value" => "0" )
		->json_is( "/response/0/configFile" => "configFile3" )
	, 'Was the parameter really changed?';

ok $t->delete_ok('/api/1.2/parameters/' . $para_id )->status_is(200)
	, 'Does the paramter deleted?';




ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->post_ok( '/login', => form => { u =>Test::TestHelper::FEDERATION_USER , p => Test::TestHelper::FEDERATION_USER_PASSWORD } )->status_is(302)
	->or( sub { diag $t->tx->res->content->asset->{content}; } ), 'Should login?';

ok $t->post_ok('/api/1.2/parameters' => {Accept => 'application/json'} => json => [
        {
            'name'  => 'param3',
            'configFile' => 'configFile3',
            'value'      => 'value3',
            'secure'     => '0'
        }]
    )->status_is(403)
	->or( sub { diag $t->tx->res->content->asset->{content}; } )
	->json_is( "/alerts/0/text" => "You must be an admin or oper to perform this operation!" )
		, 'Does the paramters created return?';

ok $t->put_ok('/api/1.2/parameters/' . $para_id => {Accept => 'application/json'} => json => {
    })->status_is(403)
	->or( sub { diag $t->tx->res->content->asset->{content}; } )
	->json_is( "/alerts/0/text" => "You must be an admin or oper to perform this operation!" )
		, 'Does the paramters modified return?';

$para_id = &get_param_id('param1');
ok $t->delete_ok('/api/1.2/parameters/' . $para_id )->status_is(403)
	->or( sub { diag $t->tx->res->content->asset->{content}; } )
	->json_is( "/alerts/0/text" => "You must be an admin or oper to perform this operation!" )
		, 'Does the paramter delete return?';

$para_id = &get_param_id('domain_name');
ok $t->get_ok('/api/1.2/parameters/'. $para_id)->status_is(200)
	->or( sub { diag $t->tx->res->content->asset->{content}; } )
	->json_is( "/response/0/name" => "domain_name" )
	->json_is( "/response/0/value" => "foo.com" )
	->json_is( "/response/0/configFile" => "CRConfig.json" )
	->json_is( "/response/0/secure" => "0" )
		, 'Does the paramter get return?';

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



