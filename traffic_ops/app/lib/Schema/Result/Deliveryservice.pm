use utf8;
package Schema::Result::Deliveryservice;

# Created by DBIx::Class::Schema::Loader
# DO NOT MODIFY THE FIRST PART OF THIS FILE

=head1 NAME

Schema::Result::Deliveryservice

=cut

use strict;
use warnings;

use base 'DBIx::Class::Core';

=head1 TABLE: C<deliveryservice>

=cut

__PACKAGE__->table("deliveryservice");

=head1 ACCESSORS

=head2 id

  data_type: 'bigint'
  is_auto_increment: 1
  is_nullable: 0
  sequence: 'deliveryservice_id_seq'

=head2 xml_id

  data_type: 'text'
  is_nullable: 0

=head2 active

  data_type: 'boolean'
  default_value: false
  is_nullable: 0

=head2 anonymous_blocking_enabled
  data_type: 'boolean'
  default_value: false
  is_nullable: 0

=head2 dscp

  data_type: 'bigint'
  is_nullable: 0

=head2 signing_algorithm

  data_type: 'deliveryservice_signature_type'
  is_nullable: 1

=head2 qstring_ignore

  data_type: 'smallint'
  is_nullable: 1

=head2 geo_limit

  data_type: 'smallint'
  default_value: 0
  is_nullable: 1

=head2 http_bypass_fqdn

  data_type: 'text'
  is_nullable: 1

=head2 dns_bypass_ip

  data_type: 'text'
  is_nullable: 1

=head2 dns_bypass_ip6

  data_type: 'text'
  is_nullable: 1

=head2 dns_bypass_ttl

  data_type: 'bigint'
  is_nullable: 1

=head2 type

  data_type: 'bigint'
  is_foreign_key: 1
  is_nullable: 0

=head2 profile

  data_type: 'bigint'
  is_foreign_key: 1
  is_nullable: 1

=head2 cdn_id

  data_type: 'bigint'
  is_foreign_key: 1
  is_nullable: 0

=head2 ccr_dns_ttl

  data_type: 'bigint'
  is_nullable: 1

=head2 global_max_mbps

  data_type: 'bigint'
  is_nullable: 1

=head2 global_max_tps

  data_type: 'bigint'
  is_nullable: 1

=head2 long_desc

  data_type: 'text'
  is_nullable: 1

=head2 long_desc_1

  data_type: 'text'
  is_nullable: 1

=head2 long_desc_2

  data_type: 'text'
  is_nullable: 1

=head2 max_dns_answers

  data_type: 'bigint'
  default_value: 5
  is_nullable: 1

=head2 info_url

  data_type: 'text'
  is_nullable: 1

=head2 miss_lat

  data_type: 'numeric'
  is_nullable: 1

=head2 miss_long

  data_type: 'numeric'
  is_nullable: 1

=head2 check_path

  data_type: 'text'
  is_nullable: 1

=head2 last_updated

  data_type: 'timestamp with time zone'
  default_value: current_timestamp
  is_nullable: 0
  original: {default_value => \"now()"}

=head2 protocol

  data_type: 'smallint'
  default_value: 0
  is_nullable: 1

=head2 ssl_key_version

  data_type: 'bigint'
  default_value: 0
  is_nullable: 1

=head2 ipv6_routing_enabled

  data_type: 'boolean'
  default_value: false
  is_nullable: 1

=head2 range_request_handling

  data_type: 'smallint'
  default_value: 0
  is_nullable: 1

=head2 edge_header_rewrite

  data_type: 'text'
  is_nullable: 1

=head2 origin_shield

  data_type: 'text'
  is_nullable: 1

=head2 mid_header_rewrite

  data_type: 'text'
  is_nullable: 1

=head2 regex_remap

  data_type: 'text'
  is_nullable: 1

=head2 cacheurl

  data_type: 'text'
  is_nullable: 1

=head2 remap_text

  data_type: 'text'
  is_nullable: 1

=head2 multi_site_origin

  data_type: 'boolean'
  default_value: false
  is_nullable: 1

=head2 display_name

  data_type: 'text'
  is_nullable: 0

=head2 tr_response_headers

  data_type: 'text'
  is_nullable: 1

=head2 initial_dispersion

  data_type: 'bigint'
  default_value: 1
  is_nullable: 1

=head2 dns_bypass_cname

  data_type: 'text'
  is_nullable: 1

=head2 tr_request_headers

  data_type: 'text'
  is_nullable: 1

=head2 regional_geo_blocking

  data_type: 'boolean'
  default_value: false
  is_nullable: 0

=head2 geo_provider

  data_type: 'smallint'
  default_value: 0
  is_nullable: 1

=head2 geo_limit_countries

  data_type: 'text'
  is_nullable: 1

=head2 logs_enabled

  data_type: 'boolean'
  default_value: false
  is_nullable: 1

=head2 multi_site_origin_algorithm

  data_type: 'smallint'
  is_nullable: 1

=head2 geolimit_redirect_url

  data_type: 'text'
  is_nullable: 1

=head2 tenant_id

  data_type: 'bigint'
  is_foreign_key: 1
  is_nullable: 1

=head2 routing_name

  data_type: 'text'
  default_value: 'cdn'
  is_nullable: 0

=head2 deep_caching_type

  data_type: 'enum'
  default_value: 'NEVER'
  extra: {custom_type_name => "deep_caching_type",list => ["NEVER","ALWAYS"]}
  is_nullable: 1

=head2 fq_pacing_rate

  data_type: 'bigint'
  default_value: 0
  is_nullable: 1

=cut

__PACKAGE__->add_columns(
  "id",
  {
    data_type         => "bigint",
    is_auto_increment => 1,
    is_nullable       => 0,
    sequence          => "deliveryservice_id_seq",
  },
  "xml_id",
  { data_type => "text", is_nullable => 0 },
  "active",
  { data_type => "boolean", default_value => \"false", is_nullable => 0 },
  "anonymous_blocking_enabled",
  { data_type => "boolean", default_value => \"false", is_nullable => 0 },
  "dscp",
  { data_type => "bigint", is_nullable => 0 },
  "signing_algorithm",
  { data_type => "deliveryservice_signature_type", is_nullable => 1 },
  "qstring_ignore",
  { data_type => "smallint", is_nullable => 1 },
  "geo_limit",
  { data_type => "smallint", default_value => 0, is_nullable => 1 },
  "http_bypass_fqdn",
  { data_type => "text", is_nullable => 1 },
  "dns_bypass_ip",
  { data_type => "text", is_nullable => 1 },
  "dns_bypass_ip6",
  { data_type => "text", is_nullable => 1 },
  "dns_bypass_ttl",
  { data_type => "bigint", is_nullable => 1 },
  "type",
  { data_type => "bigint", is_foreign_key => 1, is_nullable => 0 },
  "profile",
  { data_type => "bigint", is_foreign_key => 1, is_nullable => 1 },
  "cdn_id",
  { data_type => "bigint", is_foreign_key => 1, is_nullable => 0 },
  "ccr_dns_ttl",
  { data_type => "bigint", is_nullable => 1 },
  "global_max_mbps",
  { data_type => "bigint", is_nullable => 1 },
  "global_max_tps",
  { data_type => "bigint", is_nullable => 1 },
  "long_desc",
  { data_type => "text", is_nullable => 1 },
  "long_desc_1",
  { data_type => "text", is_nullable => 1 },
  "long_desc_2",
  { data_type => "text", is_nullable => 1 },
  "max_dns_answers",
  { data_type => "bigint", default_value => 5, is_nullable => 1 },
  "info_url",
  { data_type => "text", is_nullable => 1 },
  "miss_lat",
  { data_type => "numeric", is_nullable => 1 },
  "miss_long",
  { data_type => "numeric", is_nullable => 1 },
  "check_path",
  { data_type => "text", is_nullable => 1 },
  "last_updated",
  {
    data_type     => "timestamp with time zone",
    default_value => \"current_timestamp",
    is_nullable   => 0,
    original      => { default_value => \"now()" },
  },
  "protocol",
  { data_type => "smallint", default_value => 0, is_nullable => 1 },
  "ssl_key_version",
  { data_type => "bigint", default_value => 0, is_nullable => 1 },
  "ipv6_routing_enabled",
  { data_type => "boolean", default_value => \"false", is_nullable => 1 },
  "range_request_handling",
  { data_type => "smallint", default_value => 0, is_nullable => 1 },
  "edge_header_rewrite",
  { data_type => "text", is_nullable => 1 },
  "origin_shield",
  { data_type => "text", is_nullable => 1 },
  "mid_header_rewrite",
  { data_type => "text", is_nullable => 1 },
  "regex_remap",
  { data_type => "text", is_nullable => 1 },
  "cacheurl",
  { data_type => "text", is_nullable => 1 },
  "remap_text",
  { data_type => "text", is_nullable => 1 },
  "multi_site_origin",
  { data_type => "boolean", default_value => \"false", is_nullable => 1 },
  "display_name",
  { data_type => "text", is_nullable => 0 },
  "tr_response_headers",
  { data_type => "text", is_nullable => 1 },
  "initial_dispersion",
  { data_type => "bigint", default_value => 1, is_nullable => 1 },
  "dns_bypass_cname",
  { data_type => "text", is_nullable => 1 },
  "tr_request_headers",
  { data_type => "text", is_nullable => 1 },
  "regional_geo_blocking",
  { data_type => "boolean", default_value => \"false", is_nullable => 0 },
  "geo_provider",
  { data_type => "smallint", default_value => 0, is_nullable => 1 },
  "geo_limit_countries",
  { data_type => "text", is_nullable => 1 },
  "logs_enabled",
  { data_type => "boolean", default_value => \"false", is_nullable => 1 },
  "multi_site_origin_algorithm",
  { data_type => "smallint", is_nullable => 1 },
  "geolimit_redirect_url",
  { data_type => "text", is_nullable => 1 },
  "tenant_id",
  { data_type => "bigint", is_foreign_key => 1, is_nullable => 1 },
  "routing_name",
  { data_type => "text", default_value => "cdn", is_nullable => 0 },
  "deep_caching_type",
  {
    data_type => "enum",
    default_value => "NEVER",
    extra => {
      custom_type_name => "deep_caching_type",
      list => ["NEVER", "ALWAYS"],
    },
    is_nullable => 1,
  },
  "fq_pacing_rate",
  { data_type => "bigint", default_value => 0, is_nullable => 1 },
);

=head1 PRIMARY KEY

=over 4

=item * L</id>

=item * L</type>

=back

=cut

__PACKAGE__->set_primary_key("id", "type");

=head1 UNIQUE CONSTRAINTS

=head2 C<idx_140234_ds_id_unique>

=over 4

=item * L</id>

=back

=cut

__PACKAGE__->add_unique_constraint("idx_140234_ds_id_unique", ["id"]);

=head2 C<idx_140234_ds_name_unique>

=over 4

=item * L</xml_id>

=back

=cut

__PACKAGE__->add_unique_constraint("idx_140234_ds_name_unique", ["xml_id"]);

=head1 RELATIONS

=head2 cdn

Type: belongs_to

Related object: L<Schema::Result::Cdn>

=cut

__PACKAGE__->belongs_to(
  "cdn",
  "Schema::Result::Cdn",
  { id => "cdn_id" },
  { is_deferrable => 0, on_delete => "RESTRICT", on_update => "RESTRICT" },
);

=head2 deliveryservice_regexes

Type: has_many

Related object: L<Schema::Result::DeliveryserviceRegex>

=cut

__PACKAGE__->has_many(
  "deliveryservice_regexes",
  "Schema::Result::DeliveryserviceRegex",
  { "foreign.deliveryservice" => "self.id" },
  { cascade_copy => 0, cascade_delete => 0 },
);

=head2 deliveryservice_servers

Type: has_many

Related object: L<Schema::Result::DeliveryserviceServer>

=cut

__PACKAGE__->has_many(
  "deliveryservice_servers",
  "Schema::Result::DeliveryserviceServer",
  { "foreign.deliveryservice" => "self.id" },
  { cascade_copy => 0, cascade_delete => 0 },
);

=head2 deliveryservice_tmusers

Type: has_many

Related object: L<Schema::Result::DeliveryserviceTmuser>

=cut

__PACKAGE__->has_many(
  "deliveryservice_tmusers",
  "Schema::Result::DeliveryserviceTmuser",
  { "foreign.deliveryservice" => "self.id" },
  { cascade_copy => 0, cascade_delete => 0 },
);

=head2 federation_deliveryservices

Type: has_many

Related object: L<Schema::Result::FederationDeliveryservice>

=cut

__PACKAGE__->has_many(
  "federation_deliveryservices",
  "Schema::Result::FederationDeliveryservice",
  { "foreign.deliveryservice" => "self.id" },
  { cascade_copy => 0, cascade_delete => 0 },
);

=head2 jobs

Type: has_many

Related object: L<Schema::Result::Job>

=cut

__PACKAGE__->has_many(
  "jobs",
  "Schema::Result::Job",
  { "foreign.job_deliveryservice" => "self.id" },
  { cascade_copy => 0, cascade_delete => 0 },
);

=head2 origins

Type: has_many

Related object: L<Schema::Result::Origin>

=cut

__PACKAGE__->has_many(
  "origins",
  "Schema::Result::Origin",
  { "foreign.deliveryservice" => "self.id" },
  { cascade_copy => 0, cascade_delete => 0 },
);

=head2 profile

Type: belongs_to

Related object: L<Schema::Result::Profile>

=cut

__PACKAGE__->belongs_to(
  "profile",
  "Schema::Result::Profile",
  { id => "profile" },
  {
    is_deferrable => 0,
    join_type     => "LEFT",
    on_delete     => "NO ACTION",
    on_update     => "NO ACTION",
  },
);

=head2 staticdnsentries

Type: has_many

Related object: L<Schema::Result::Staticdnsentry>

=cut

__PACKAGE__->has_many(
  "staticdnsentries",
  "Schema::Result::Staticdnsentry",
  { "foreign.deliveryservice" => "self.id" },
  { cascade_copy => 0, cascade_delete => 0 },
);

=head2 steering_target_deliveryservices

Type: has_many

Related object: L<Schema::Result::SteeringTarget>

=cut

__PACKAGE__->has_many(
  "steering_target_deliveryservices",
  "Schema::Result::SteeringTarget",
  { "foreign.deliveryservice" => "self.id" },
  { cascade_copy => 0, cascade_delete => 0 },
);

=head2 steering_target_targets

Type: has_many

Related object: L<Schema::Result::SteeringTarget>

=cut

__PACKAGE__->has_many(
  "steering_target_targets",
  "Schema::Result::SteeringTarget",
  { "foreign.target" => "self.id" },
  { cascade_copy => 0, cascade_delete => 0 },
);

=head2 tenant

Type: belongs_to

Related object: L<Schema::Result::Tenant>

=cut

__PACKAGE__->belongs_to(
  "tenant",
  "Schema::Result::Tenant",
  { id => "tenant_id" },
  {
    is_deferrable => 0,
    join_type     => "LEFT",
    on_delete     => "NO ACTION",
    on_update     => "NO ACTION",
  },
);

=head2 type

Type: belongs_to

Related object: L<Schema::Result::Type>

=cut

__PACKAGE__->belongs_to(
  "type",
  "Schema::Result::Type",
  { id => "type" },
  { is_deferrable => 0, on_delete => "NO ACTION", on_update => "NO ACTION" },
);


# Created by DBIx::Class::Schema::Loader v0.07042 @ 2018-05-17 16:24:12
# DO NOT MODIFY THIS OR ANYTHING ABOVE! md5sum:Yjz2V+duaN88QPILxLqoHg

# You can replace this text with custom code or comments, and it will be preserved on regeneration
#
# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.
1;
