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
use Schema;
use Test::TestHelper;
use strict;
use warnings;
use Schema;
use Fixtures::TmUser;
use Test::TestHelper;

BEGIN { $ENV{MOJO_MODE} = "test" }

my $dbh    = Schema->database_handle;
my $schema = Schema->connect_to_database;
my $t      = Test::Mojo->new('TrafficOps');
my $t2_id;

#unload data for a clean test
Test::TestHelper->unload_core_data($schema);

#load core test data
Test::TestHelper->load_core_data($schema);

#login
ok $t->post_ok(
	'/login',
	=> form => {
		u => Test::TestHelper::ADMIN_USER,
		p => Test::TestHelper::ADMIN_USER_PASSWORD
	}
)->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );

# ####################### RW testing - careful with these! #####################################################

# create
#HTTP DS
ok $t->post_ok(
	'/ds/create' => form => {
		'ds.active'                      => '1',
		'ds.ccr_dns_ttl'                 => '3600',
		'ds.check_path'                  => '/clientaccesspolicy.xml',
		'ds.dns_bypass_ip'               => '',
		'ds.dns_bypass_ip6'              => '',
		'ds.dns_bypass_cname'            => '',
		'ds.dns_bypass_ttl'              => '30',
		'ds.dscp'                        => '40',
		'ds.routing_name'                => 'foo',
		'ds.deep_caching_type'           => 'NEVER',
		'ds.geo_limit'                   => '0',
		'ds.geo_limit_countries'         => '',
		'ds.geo_provider'                => '1',
		'ds.global_max_mbps'             => '',
		'ds.global_max_tps'              => '',
		'ds.http_bypass_fqdn'            => '',
		'ds.info_url'                    => '',
		'ds.long_desc'                   => 'description',
		'ds.long_desc_1'                 => 'ccp',
		'ds.long_desc_2'                 => 'Columbus',
		'ds.max_dns_answers'             => '0',
		'ds.miss_lat'                    => '41.881944',
		'ds.miss_long'                   => '-87.627778',
		'origin.org_server_fqdn'         => 'http://jvd.knutsel.com',
		'ds.multi_site_origin'           => '0',
		'ds.multi_site_origin_algorithm' => '0',
		'ds.profile'                     => '100',
		'ds.cdn_id'                      => '100',
		'ds.qstring_ignore'              => '0',
		're_order_0'                     => '0',
		're_re_0'                        => '.*\.jvdtest\..*',
		're_type_0'                      => 'HOST_REGEXP',
		'ds.signed'                      => '0',
		'ds.type'                        => '8',
		'ds.xml_id'                      => 'tst_xml_id_1',
		'ds.protocol'                    => '3',
		'ds.edge_header_rewrite'         => '',
		'ds.mid_header_rewrite'          => '',
		'ds.regex_remap'                 => '',
		'ds.origin_shield'               => '',
		'ds.range_request_handling'      => '0',
		'ds.ipv6_routing_enabled'        => '1',
		'ds.display_name'                => 'display name 1',
		'ds.regional_geo_blocking'       => '1',
		'ds.geolimit_redirect_url'       => '',
		'ds.anonymous_blocking_enabled'  => '1',
	}
)->status_is(302), "create HTTP delivery service";
my $t1_id = &get_ds_id('tst_xml_id_1');
ok defined($t1_id), "validated http ds was added";

# DNS DS
ok $t->post_ok(
	'/ds/create' => form => {
		'ds.active'                      => '0',
		'ds.ccr_dns_ttl'                 => '30',
		'ds.check_path'                  => '/clientaccesspolicy.xml',
		'ds.dns_bypass_ip'               => '',
		'ds.dns_bypass_ip6'              => '',
		'ds.dns_bypass_cname'            => '',
		'ds.dns_bypass_ttl'              => '30',
		'ds.dscp'                        => '42',
		'ds.routing_name'                => '',
		'ds.deep_caching_type'           => 'NEVER',
		'ds.geo_limit'                   => '0',
		'ds.geo_limit_countries'         => '',
		'ds.global_max_mbps'             => '',
		'ds.global_max_tps'              => '',
		'ds.http_bypass_fqdn'            => '',
		'ds.info_url'                    => '',
		'ds.long_desc'                   => '',
		'ds.long_desc_1'                 => 'ccp',
		'ds.long_desc_2'                 => 'Columbus',
		'ds.max_dns_answers'             => '0',
		'ds.miss_lat'                    => '41.881944',
		'ds.miss_long'                   => '-87.627778',
		'origin.org_server_fqdn'         => 'http://jvd-1.knutsel.com',
		'ds.multi_site_origin'           => '0',
		'ds.multi_site_origin_algorithm' => '0',
		'ds.profile'                     => '100',
		'ds.cdn_id'                      => '100',
		'ds.qstring_ignore'              => '0',
		'ds.signed'                      => '0',
		'ds.type'                        => '9',
		'ds.xml_id'                      => 'tst_xml_id_2',
		'ds.protocol'                    => '0',
		'ds.edge_header_rewrite'         => '',
		'ds.mid_header_rewrite'          => '',
		'ds.regex_remap'                 => '',
		'ds.range_request_handling'      => '0',
		'ds.origin_shield'               => '',
		're_order_0'                     => '0',
		're_re_0'                        => '.*\.jvdtest-1\..*',
		're_type_0'                      => 'HOST_REGEXP',
		'ds.ipv6_routing_enabled'        => '0',
		'ds.display_name'                => 'display name 2',
		'ds.regional_geo_blocking'       => '0',
		'ds.geolimit_redirect_url'       => '',
		'ds.anonymous_blocking_enabled'  => '0',
	}
)->status_is(302), "create DNS DeliveryService";
$t2_id = &get_ds_id('tst_xml_id_2');
ok defined($t2_id), "validated dns ds was added";

#create DS ALL FIELDS
ok $t->post_ok(
	'/ds/create' => form => {
		'ds.active'                      => '1',
		'ds.ccr_dns_ttl'                 => '3600',
		'ds.check_path'                  => '/clientaccesspolicy.xml',
		'ds.dns_bypass_ip'               => '10.10.10.10',
		'ds.dns_bypass_ip6'              => '2001:558:fee8:180::2/64',
		'ds.dns_bypass_cname'            => 'bypass.knutsel.com',
		'ds.dns_bypass_ttl'              => '30',
		'ds.dscp'                        => '40',
		'ds.routing_name'                => 'foo',
		'ds.deep_caching_type'           => 'NEVER',
		'ds.geo_limit'                   => '1',
		'ds.geo_limit_countries'         => '',
		'ds.global_max_mbps'             => '30G',
		'ds.global_max_tps'              => '10000',
		'ds.http_bypass_fqdn'            => 'overflow.knutsel.com',
		'ds.info_url'                    => 'http://knutsel.com',
		'ds.long_desc'                   => 'long',
		'ds.long_desc_1'                 => 'cust',
		'ds.long_desc_2'                 => 'service',
		'ds.max_dns_answers'             => '0',
		'ds.miss_lat'                    => '41.881944',
		'ds.miss_long'                   => '-87.627778',
		'origin.org_server_fqdn'         => 'http://jvd.knutsel.com',
		'ds.multi_site_origin'           => '0',
		'ds.multi_site_origin_algorithm' => '0',
		'ds.profile'                     => '100',
		'ds.cdn_id'                      => '100',
		'ds.qstring_ignore'              => '1',
		'ds.signed'                      => '1',
		'ds.type'                        => '9',
		'ds.xml_id'                      => 'tst_xml_id_3',
		'ds.protocol'                    => '0',
		'ds.edge_header_rewrite'         => '',
		'ds.mid_header_rewrite'          => '',
		'ds.regex_remap'                 => '',
		'ds.range_request_handling'      => '0',
		'ds.origin_shield'               => '',
		're_order_0'                     => '0',
		're_re_0'                        => '.*\.jvdtest-3\..*',
		're_type_0'                      => 'HOST_REGEXP',
		're_order_1'                     => '0',
		're_re_1'                        => '/path/to/goodies/.*',
		're_type_1'                      => 'PATH_REGEXP',
		'ds.ipv6_routing_enabled'        => '1',
		'ds.display_name'                => 'display name 3',
		'ds.regional_geo_blocking'       => '0',
		'ds.geolimit_redirect_url'       => 'http://knutsel3.com',
		'ds.anonymous_blocking_enabled'  => '0',
	}
)->status_is(302), "create HTTP_NO_CACHE deliveryservice";

#add - validate 200 response (data is actually added to DB when create is called)
ok $t->get_ok('/ds/add')->status_is(200), "validate add screen";

# validate existing delivery service
ok $t->get_ok('/ds/1')->status_is(200), "validate existing delivery service";

# validate existing delivery service
ok $t->get_ok('/ds/1')->status_is(200), "validate existing delivery service";

$t2_id = &get_ds_id('tst_xml_id_2');
ok defined($t2_id), "validated delivery service with all fields was added";

# update DS
#post update
ok $t->post_ok(
	"/ds/$t2_id/update" => form => {
		'ds.active'                      => '0',
		'ds.ccr_dns_ttl'                 => '3601',
		'ds.check_path'                  => '/clientaccesspolicy.xml_update',
		'ds.dns_bypass_ip'               => '10.10.10.11',
		'ds.dns_bypass_ip6'              => '2001:558:fee8:180::1/64',
		'ds.dns_bypass_cname'            => 'updateby.knutsel.com',
		'ds.dns_bypass_ttl'              => '31',
		'ds.dscp'                        => '41',
		'ds.routing_name'                => 'bar',
		'ds.deep_caching_type'           => 'NEVER',
		'ds.geo_limit'                   => '2',
		'ds.geo_limit_countries'         => '',
		'ds.geo_provider'                => '1',
		'ds.global_max_mbps'             => '4T',
		'ds.http_bypass_fqdn'            => '',
		'ds.global_max_tps'              => '10001',
		'ds.fq_pacing_rate'              => '500k',
		'ds.info_url'                    => 'http://knutsel-update.com',
		'ds.long_desc'                   => 'long_update',
		'ds.long_desc_1'                 => 'cust_update',
		'ds.long_desc_2'                 => 'service_update',
		'ds.max_dns_answers'             => '1',
		'ds.miss_lat'                    => '0',
		'ds.miss_long'                   => '0',
		'origin.org_server_fqdn'         => 'http://update.knutsel.com',
		'ds.multi_site_origin'           => '0',
		'ds.multi_site_origin_algorithm' => '0',
		'ds.profile'                     => '200',
		'ds.cdn_id'                      => '200',
		'ds.qstring_ignore'              => '0',
		'ds.signed'                      => '0',
		'ds.type'                        => '7',
		'ds.xml_id'                      => 'tst_xml_id_2',
		'ds.protocol'                    => '1',
		'ds.edge_header_rewrite'         => '',
		'ds.mid_header_rewrite'          => '',
		'ds.regex_remap'                 => '',
		'ds.origin_shield'               => '',
		'ds.range_request_handling'      => '0',
		're_order_0'                     => '0',
		're_re_0'                        => '.*\.jvdtest-3_update\..*',
		're_type_0'                      => 'HOST_REGEXP',
		'ds.ipv6_routing_enabled'        => '1',
		'ds.display_name'                => 'Testing Delivery Service',
		'ds.tr_response_headers'         => '',
		'ds.regional_geo_blocking'       => '1',
		'ds.geolimit_redirect_url'       => 'http://update.redirect.url.com',
		'ds.anonymous_blocking_enabled'  => '1',
	}
)->status_is(302), "update deliveryservice";


#delete delivery service
ok $t->get_ok("/ds/$t2_id/delete")->status_is(302), "delete ds";

#validate it was deleted
$t2_id = &get_ds_id('tst_xml_id_2');
ok !defined($t2_id), "validated delivery service was deleted";

sub get_ds_id {
	my $xml_id = shift;
	my $q      = "select id from deliveryservice where xml_id = \'$xml_id\'";
	my $get_ds = $dbh->prepare($q);
	$get_ds->execute();
	my $p = $get_ds->fetchall_arrayref( {} );
	$get_ds->finish();
	my $id = $p->[0]->{id};
	return $id;
}
ok $t->get_ok('/logout')->status_is(302)
  ->or( sub { diag $t->tx->res->content->asset->{content}; } );
$dbh->disconnect();
done_testing();
