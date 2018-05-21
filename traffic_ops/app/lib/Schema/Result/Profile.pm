use utf8;
package Schema::Result::Profile;

# Created by DBIx::Class::Schema::Loader
# DO NOT MODIFY THE FIRST PART OF THIS FILE

=head1 NAME

Schema::Result::Profile

=cut

use strict;
use warnings;

use base 'DBIx::Class::Core';

=head1 TABLE: C<profile>

=cut

__PACKAGE__->table("profile");

=head1 ACCESSORS

=head2 id

  data_type: 'bigint'
  is_auto_increment: 1
  is_nullable: 0
  sequence: 'profile_id_seq'

=head2 name

  data_type: 'text'
  is_nullable: 0

=head2 description

  data_type: 'text'
  is_nullable: 1

=head2 last_updated

  data_type: 'timestamp with time zone'
  default_value: current_timestamp
  is_nullable: 0
  original: {default_value => \"now()"}

=head2 type

  data_type: 'enum'
  extra: {custom_type_name => "profile_type",list => ["ATS_PROFILE","TR_PROFILE","TM_PROFILE","TS_PROFILE","TP_PROFILE","INFLUXDB_PROFILE","RIAK_PROFILE","SPLUNK_PROFILE","DS_PROFILE","ORG_PROFILE","KAFKA_PROFILE","LOGSTASH_PROFILE","ES_PROFILE","UNK_PROFILE"]}
  is_nullable: 0

=head2 cdn

  data_type: 'bigint'
  is_foreign_key: 1
  is_nullable: 1

=head2 routing_disabled

  data_type: 'boolean'
  default_value: false
  is_nullable: 0

=cut

__PACKAGE__->add_columns(
  "id",
  {
    data_type         => "bigint",
    is_auto_increment => 1,
    is_nullable       => 0,
    sequence          => "profile_id_seq",
  },
  "name",
  { data_type => "text", is_nullable => 0 },
  "description",
  { data_type => "text", is_nullable => 1 },
  "last_updated",
  {
    data_type     => "timestamp with time zone",
    default_value => \"current_timestamp",
    is_nullable   => 0,
    original      => { default_value => \"now()" },
  },
  "type",
  {
    data_type => "enum",
    extra => {
      custom_type_name => "profile_type",
      list => [
        "ATS_PROFILE",
        "TR_PROFILE",
        "TM_PROFILE",
        "TS_PROFILE",
        "TP_PROFILE",
        "INFLUXDB_PROFILE",
        "RIAK_PROFILE",
        "SPLUNK_PROFILE",
        "DS_PROFILE",
        "ORG_PROFILE",
        "KAFKA_PROFILE",
        "LOGSTASH_PROFILE",
        "ES_PROFILE",
        "UNK_PROFILE",
      ],
    },
    is_nullable => 0,
  },
  "cdn",
  { data_type => "bigint", is_foreign_key => 1, is_nullable => 1 },
  "routing_disabled",
  { data_type => "boolean", default_value => \"false", is_nullable => 0 },
);

=head1 PRIMARY KEY

=over 4

=item * L</id>

=back

=cut

__PACKAGE__->set_primary_key("id");

=head1 UNIQUE CONSTRAINTS

=head2 C<idx_140397_name_unique>

=over 4

=item * L</name>

=back

=cut

__PACKAGE__->add_unique_constraint("idx_140397_name_unique", ["name"]);

=head1 RELATIONS

=head2 cdn

Type: belongs_to

Related object: L<Schema::Result::Cdn>

=cut

__PACKAGE__->belongs_to(
  "cdn",
  "Schema::Result::Cdn",
  { id => "cdn" },
  {
    is_deferrable => 0,
    join_type     => "LEFT",
    on_delete     => "RESTRICT",
    on_update     => "RESTRICT",
  },
);

=head2 deliveryservices

Type: has_many

Related object: L<Schema::Result::Deliveryservice>

=cut

__PACKAGE__->has_many(
  "deliveryservices",
  "Schema::Result::Deliveryservice",
  { "foreign.profile" => "self.id" },
  { cascade_copy => 0, cascade_delete => 0 },
);

=head2 origins

Type: has_many

Related object: L<Schema::Result::Origin>

=cut

__PACKAGE__->has_many(
  "origins",
  "Schema::Result::Origin",
  { "foreign.profile" => "self.id" },
  { cascade_copy => 0, cascade_delete => 0 },
);

=head2 profile_parameters

Type: has_many

Related object: L<Schema::Result::ProfileParameter>

=cut

__PACKAGE__->has_many(
  "profile_parameters",
  "Schema::Result::ProfileParameter",
  { "foreign.profile" => "self.id" },
  { cascade_copy => 0, cascade_delete => 0 },
);

=head2 servers

Type: has_many

Related object: L<Schema::Result::Server>

=cut

__PACKAGE__->has_many(
  "servers",
  "Schema::Result::Server",
  { "foreign.profile" => "self.id" },
  { cascade_copy => 0, cascade_delete => 0 },
);


# Created by DBIx::Class::Schema::Loader v0.07042 @ 2018-05-15 16:06:00
# DO NOT MODIFY THIS OR ANYTHING ABOVE! md5sum:eX4sXLElEMpaA0xfHTv5lw


# You can replace this text with custom code or comments, and it will be preserved on regeneration
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
#
1;
