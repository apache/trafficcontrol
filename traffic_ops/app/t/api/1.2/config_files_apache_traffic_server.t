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

# Check ATS config metadata
ok $t->get_ok("/api/1.2/servers/100/configfiles/ats")->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content} } );

# Check server routes
ok $t->get_ok("/api/1.2/profiles/100/configfiles/ats/12M_facts")->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content} } );
ok $t->get_ok("/api/1.2/servers/100/configfiles/ats/ip_allow.config")->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content} } );
ok $t->get_ok("/api/1.2/servers/100/configfiles/ats/remap.config")->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content} } );
ok $t->get_ok("/api/1.2/servers/100/configfiles/ats/parent.config")->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content} } );
ok $t->get_ok("/api/1.2/profiles/100/configfiles/ats/cache.config")->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content} } );
ok $t->get_ok("/api/1.2/servers/100/configfiles/ats/hosting.config")->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content} } );

# Check CDN routes
ok $t->get_ok("/api/1.2/cdns/cdn1/configfiles/ats/cacheurl.config")->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content} } );
ok $t->get_ok("/api/1.2/cdns/cdn1/configfiles/ats/regex_revalidate.config")->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content} } );

# Check profile routes
ok $t->get_ok("/api/1.2/profiles/100/configfiles/ats/50-ats.rules")->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content} } );
ok $t->get_ok("/api/1.2/profiles/100/configfiles/ats/drop_qstring.config")->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content} } );
ok $t->get_ok("/api/1.2/profiles/100/configfiles/ats/volume.config")->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content} } );
ok $t->get_ok("/api/1.2/profiles/100/configfiles/ats/astats.config")->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content} } );
ok $t->get_ok("/api/1.2/profiles/100/configfiles/ats/storage.config")->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content} } );


# Count the 'response number'
my $count_response = sub {
	my ( $t, $count ) = @_;
	my $json = decode_json( $t->tx->res->content->asset->slurp );
	my $r    = $json->{response};
	return $t->success( is( scalar(@$r), $count ) );
};

$t->get_ok('/api/1.2/deliveryservices?logsEnabled=true')->status_is(200)->$count_response(3)
	->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->put_ok('/api/1.2/snapshot/cdn1')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );
$dbh->disconnect();
done_testing();


