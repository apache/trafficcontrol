use utf8;
package Schema::Result::Origin;

# Created by DBIx::Class::Schema::Loader
# DO NOT MODIFY THE FIRST PART OF THIS FILE

=head1 NAME

Schema::Result::Origin

=cut

use strict;
use warnings;

use base 'DBIx::Class::Core';

=head1 TABLE: C<origin>

=cut

__PACKAGE__->table("origin");

=head1 ACCESSORS

=head2 id

  data_type: 'bigint'
  is_auto_increment: 1
  is_nullable: 0
  sequence: 'origin_id_seq'

=head2 name

  data_type: 'text'
  is_nullable: 0

=head2 fqdn

  data_type: 'text'
  is_nullable: 0

=head2 protocol

  data_type: 'enum'
  default_value: 'http'
  extra: {custom_type_name => "origin_protocol",list => ["http","https"]}
  is_nullable: 0

=head2 is_primary

  data_type: 'boolean'
  default_value: false
  is_nullable: 0

=head2 port

  data_type: 'bigint'
  is_nullable: 1

=head2 ip_address

  data_type: 'text'
  is_nullable: 1

=head2 ip6_address

  data_type: 'text'
  is_nullable: 1

=head2 deliveryservice

  data_type: 'bigint'
  is_foreign_key: 1
  is_nullable: 0

=head2 coordinate

  data_type: 'bigint'
  is_foreign_key: 1
  is_nullable: 1

=head2 profile

  data_type: 'bigint'
  is_foreign_key: 1
  is_nullable: 1

=head2 cachegroup

  data_type: 'bigint'
  is_foreign_key: 1
  is_nullable: 1

=head2 tenant

  data_type: 'bigint'
  is_foreign_key: 1
  is_nullable: 1

=head2 last_updated

  data_type: 'timestamp with time zone'
  default_value: current_timestamp
  is_nullable: 0
  original: {default_value => \"now()"}

=cut

__PACKAGE__->add_columns(
  "id",
  {
    data_type         => "bigint",
    is_auto_increment => 1,
    is_nullable       => 0,
    sequence          => "origin_id_seq",
  },
  "name",
  { data_type => "text", is_nullable => 0 },
  "fqdn",
  { data_type => "text", is_nullable => 0 },
  "protocol",
  {
    data_type => "enum",
    default_value => "http",
    extra => { custom_type_name => "origin_protocol", list => ["http", "https"] },
    is_nullable => 0,
  },
  "is_primary",
  { data_type => "boolean", default_value => \"false", is_nullable => 0 },
  "port",
  { data_type => "bigint", is_nullable => 1 },
  "ip_address",
  { data_type => "text", is_nullable => 1 },
  "ip6_address",
  { data_type => "text", is_nullable => 1 },
  "deliveryservice",
  { data_type => "bigint", is_foreign_key => 1, is_nullable => 0 },
  "coordinate",
  { data_type => "bigint", is_foreign_key => 1, is_nullable => 1 },
  "profile",
  { data_type => "bigint", is_foreign_key => 1, is_nullable => 1 },
  "cachegroup",
  { data_type => "bigint", is_foreign_key => 1, is_nullable => 1 },
  "tenant",
  { data_type => "bigint", is_foreign_key => 1, is_nullable => 1 },
  "last_updated",
  {
    data_type     => "timestamp with time zone",
    default_value => \"current_timestamp",
    is_nullable   => 0,
    original      => { default_value => \"now()" },
  },
);

=head1 PRIMARY KEY

=over 4

=item * L</id>

=back

=cut

__PACKAGE__->set_primary_key("id");

=head1 UNIQUE CONSTRAINTS

=head2 C<origin_name_key>

=over 4

=item * L</name>

=back

=cut

__PACKAGE__->add_unique_constraint("origin_name_key", ["name"]);

=head1 RELATIONS

=head2 cachegroup

Type: belongs_to

Related object: L<Schema::Result::Cachegroup>

=cut

__PACKAGE__->belongs_to(
  "cachegroup",
  "Schema::Result::Cachegroup",
  { id => "cachegroup" },
  {
    is_deferrable => 0,
    join_type     => "LEFT",
    on_delete     => "RESTRICT",
    on_update     => "NO ACTION",
  },
);

=head2 coordinate

Type: belongs_to

Related object: L<Schema::Result::Coordinate>

=cut

__PACKAGE__->belongs_to(
  "coordinate",
  "Schema::Result::Coordinate",
  { id => "coordinate" },
  {
    is_deferrable => 0,
    join_type     => "LEFT",
    on_delete     => "RESTRICT",
    on_update     => "NO ACTION",
  },
);

=head2 deliveryservice

Type: belongs_to

Related object: L<Schema::Result::Deliveryservice>

=cut

__PACKAGE__->belongs_to(
  "deliveryservice",
  "Schema::Result::Deliveryservice",
  { id => "deliveryservice" },
  { is_deferrable => 0, on_delete => "CASCADE", on_update => "NO ACTION" },
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
    on_delete     => "RESTRICT",
    on_update     => "NO ACTION",
  },
);

=head2 tenant

Type: belongs_to

Related object: L<Schema::Result::Tenant>

=cut

__PACKAGE__->belongs_to(
  "tenant",
  "Schema::Result::Tenant",
  { id => "tenant" },
  {
    is_deferrable => 0,
    join_type     => "LEFT",
    on_delete     => "RESTRICT",
    on_update     => "NO ACTION",
  },
);


# Created by DBIx::Class::Schema::Loader v0.07042 @ 2018-05-15 16:06:00
# DO NOT MODIFY THIS OR ANYTHING ABOVE! md5sum:EFdWsJg/ANV/vUHBHfK0iA

#
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
