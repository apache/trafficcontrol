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
my $dbh    = Schema->database_handle;
my $t      = Test::Mojo->new('TrafficOps');
no warnings 'once';

ok $t->post_ok( '/login', => form => { u => Test::TestHelper::ADMIN_USER, p => Test::TestHelper::ADMIN_USER_PASSWORD } )->status_is(302)
	->or( sub { diag $t->tx->res->content->asset->{content}; } );

#unload data for a clean test
Test::TestHelper->unload_core_data($schema);

#load core test data
Test::TestHelper->load_core_data($schema);

my $q      = 'select * from parameter limit 1';
my $get_ds = $dbh->prepare($q);
$get_ds->execute();
my $p = $get_ds->fetchall_arrayref( {} );
$get_ds->finish();

# jsons
$t->get_ok('/dataparameter')->status_is(200)->json_has('/3/last_updated')->json_has('/3/value')->json_has('/3/name')->json_has('/3/id');
$t->get_ok('/dataparameter/orderby/last_updated')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )->json_has('/3/last_updated')
	->json_has('/3/value')->json_has('/3/name')->json_has('/3/id');

# create a new param
$t->post_ok( '/parameter/create' => form => { name => 'auto_tstinsertparam', config_file => 'auto_tstfile', value => 'auto_tstvalue', profile => '13' } )
	->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );

# modify and delete it
$q = 'select id from parameter where name = \'auto_tstinsertparam\'';
my $get_param = $dbh->prepare($q);
$get_param->execute();
$p = $get_param->fetchall_arrayref( {} );
$get_param->finish();
my $i = 0;
while ( defined( $p->[$i] ) ) {
	my $id = $p->[$i]->{id};
	$t->post_ok( '/parameter/update/'
			. $id => form => { name => 'auto_tstinsertparam', config_file => 'auto_tstfile', value => 'auto_tst_UPDATED_value', profile => '13,14' } )
		->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );
	$t->get_ok( '/parameter/delete/' . $id )->status_is(302);    # no need to delete the profile_parameter entry - the cascade does that.
	$i++;
}

# test secure parameter: create a non-secure parameter, modify to secure parameter, delete it
$t->post_ok( '/parameter/create' => form => { name => 'auto_tstinsertparam', config_file => 'auto_tstfile', value => 'auto_tstvalue', profile => '13' } )
	->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );
$q = 'select id from parameter where name = \'auto_tstinsertparam\'';
$get_param = $dbh->prepare($q);
$get_param->execute();
$p = $get_param->fetchall_arrayref( {} );
$get_param->finish();
$i = 0;
while ( defined( $p->[$i] ) ) {
	my $id = $p->[$i]->{id};
	$t->post_ok( '/parameter/update/'
			. $id => form => { name => 'auto_tstinsertparam', config_file => 'auto_tstfile', value => 'auto_tst_UPDATED_value', secure => '1', profile => '13' } )
		->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );
	$t->get_ok( '/parameter/delete/' . $id )->status_is(302);    # no need to delete the profile_parameter entry - the cascade does that.
	$i++;
}

ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );
ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );
$dbh->disconnect();
done_testing();
