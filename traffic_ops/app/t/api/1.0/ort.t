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
use Data::Dumper;
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

#$t->app->log->level('trace');

Test::TestHelper->unload_core_data($schema);
Test::TestHelper->load_core_data($schema);

ok $t->post_ok( '/login', => form => { u => Test::TestHelper::ADMIN_USER, p => Test::TestHelper::ADMIN_USER_PASSWORD } )->status_is(302)
	->or( sub { diag $t->tx->res->content->asset->{content}; } );

# Count the 'config_files'
my $config_files_count = sub {
	my ( $t, $count ) = @_;
	my $json               = decode_json( $t->tx->res->content->asset->slurp );
	my $config_files       = $json->{config_files};
	my $config_files_count = keys( %{$config_files} );
	return $t->success( is( $config_files_count, $count ) );
};

$t->get_ok('/ort/atlanta-edge-01/ort1')->status_is(200)->$config_files_count(16)->or( sub { diag $t->tx->res->content->asset->{content}; } );
$t->get_ok('/ort/atlanta-mid-01/ort1')->status_is(200)->$config_files_count(16)->or( sub  { diag $t->tx->res->content->asset->{content}; } );

&ort_check('atlanta-edge-01');
&genfiles_check('atlanta-edge-01');

&ort_check('atlanta-mid-01');
&genfiles_check('atlanta-mid-01');

sub ort_check {
	my $host_name = shift;
	$t->get_ok( '/ort/' . $host_name . '/ort1' )->status_is(200)->or( sub      { diag $t->tx->res->content->asset->{content}; } );
	$t->get_ok( '/ort/' . $host_name . '/packages' )->status_is(200)->or( sub  { diag $t->tx->res->content->asset->{content}; } );
	$t->get_ok( '/ort/' . $host_name . '/chkconfig' )->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );
}

sub genfiles_check {
	my $host_name = shift;
	$t->get_ok( '/genfiles/view/' . $host_name . '/50-ats.rules' )->status_is(200)->or( sub        { diag $t->tx->res->content->asset->{content}; } );
	$t->get_ok( '/genfiles/view/' . $host_name . '/astats.config' )->status_is(200)->or( sub       { diag $t->tx->res->content->asset->{content}; } );
	$t->get_ok( '/genfiles/view/' . $host_name . '/cache.config' )->status_is(200)->or( sub        { diag $t->tx->res->content->asset->{content}; } );
	$t->get_ok( '/genfiles/view/' . $host_name . '/cacheurl.config' )->status_is(200)->or( sub     { diag $t->tx->res->content->asset->{content}; } );
	$t->get_ok( '/genfiles/view/' . $host_name . '/drop_qstring.config' )->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );
	$t->get_ok( '/genfiles/view/' . $host_name . '/hosting.config' )->status_is(200)->or( sub      { diag $t->tx->res->content->asset->{content}; } );
	$t->get_ok( '/genfiles/view/' . $host_name . '/ip_allow.config' )->status_is(200)->or( sub     { diag $t->tx->res->content->asset->{content}; } );
	$t->get_ok( '/genfiles/view/' . $host_name . '/logs_xml.config' )->status_is(200)->or( sub     { diag $t->tx->res->content->asset->{content}; } );
	$t->get_ok( '/genfiles/view/' . $host_name . '/parent.config' )->status_is(200)->or( sub       { diag $t->tx->res->content->asset->{content}; } )
		->content_like( qr/^# DO NOT EDIT/,                   'parent.config: has at least lead comment' )
		->content_like( qr/Generated for \w+(-mid|-edge.* parent=".*" secondary_parent=".*")/s, 'parent.config: if edge, see parent= and secondary_parent=' );

	$t->get_ok( '/genfiles/view/' . $host_name . '/plugin.config' )->status_is(200)->or( sub  { diag $t->tx->res->content->asset->{content}; } );
	$t->get_ok( '/genfiles/view/' . $host_name . '/records.config' )->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );

	#$t->get_ok( '/genfiles/view/' . $host_name . '/regex_revalidate.config' )->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );
	$t->get_ok( '/genfiles/view/' . $host_name . '/remap.config' )->status_is(200)->or( sub   { diag $t->tx->res->content->asset->{content}; } );
	$t->get_ok( '/genfiles/view/' . $host_name . '/crontab_root' )->status_is(200)->or( sub   { diag $t->tx->res->content->asset->{content}; } );
	$t->get_ok( '/genfiles/view/' . $host_name . '/storage.config' )->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );
	$t->get_ok( '/genfiles/view/' . $host_name . '/volume.config' )->status_is(200)->or( sub  { diag $t->tx->res->content->asset->{content}; } );
}

ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );
$dbh->disconnect();
done_testing();
