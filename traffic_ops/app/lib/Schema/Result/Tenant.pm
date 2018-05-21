use utf8;
package Schema::Result::Tenant;

# Created by DBIx::Class::Schema::Loader
# DO NOT MODIFY THE FIRST PART OF THIS FILE

=head1 NAME

Schema::Result::Tenant

=cut

use strict;
use warnings;

use base 'DBIx::Class::Core';

=head1 TABLE: C<tenant>

=cut

__PACKAGE__->table("tenant");

=head1 ACCESSORS

=head2 id

  data_type: 'bigint'
  is_auto_increment: 1
  is_nullable: 0
  sequence: 'tenant_id_seq'

=head2 name

  data_type: 'text'
  is_nullable: 0

=head2 active

  data_type: 'boolean'
  default_value: false
  is_nullable: 0

=head2 parent_id

  data_type: 'bigint'
  default_value: 1
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
    sequence          => "tenant_id_seq",
  },
  "name",
  { data_type => "text", is_nullable => 0 },
  "active",
  { data_type => "boolean", default_value => \"false", is_nullable => 0 },
  "parent_id",
  {
    data_type      => "bigint",
    default_value  => 1,
    is_foreign_key => 1,
    is_nullable    => 1,
  },
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

=head2 C<tenant_name_key>

=over 4

=item * L</name>

=back

=cut

__PACKAGE__->add_unique_constraint("tenant_name_key", ["name"]);

=head1 RELATIONS

=head2 deliveryservices

Type: has_many

Related object: L<Schema::Result::Deliveryservice>

=cut

__PACKAGE__->has_many(
  "deliveryservices",
  "Schema::Result::Deliveryservice",
  { "foreign.tenant_id" => "self.id" },
  { cascade_copy => 0, cascade_delete => 0 },
);

=head2 origins

Type: has_many

Related object: L<Schema::Result::Origin>

=cut

__PACKAGE__->has_many(
  "origins",
  "Schema::Result::Origin",
  { "foreign.tenant" => "self.id" },
  { cascade_copy => 0, cascade_delete => 0 },
);

=head2 parent

Type: belongs_to

Related object: L<Schema::Result::Tenant>

=cut

__PACKAGE__->belongs_to(
  "parent",
  "Schema::Result::Tenant",
  { id => "parent_id" },
  {
    is_deferrable => 0,
    join_type     => "LEFT",
    on_delete     => "NO ACTION",
    on_update     => "NO ACTION",
  },
);

=head2 tenants

Type: has_many

Related object: L<Schema::Result::Tenant>

=cut

__PACKAGE__->has_many(
  "tenants",
  "Schema::Result::Tenant",
  { "foreign.parent_id" => "self.id" },
  { cascade_copy => 0, cascade_delete => 0 },
);

=head2 tm_users

Type: has_many

Related object: L<Schema::Result::TmUser>

=cut

__PACKAGE__->has_many(
  "tm_users",
  "Schema::Result::TmUser",
  { "foreign.tenant_id" => "self.id" },
  { cascade_copy => 0, cascade_delete => 0 },
);


# Created by DBIx::Class::Schema::Loader v0.07042 @ 2018-05-15 16:06:00
# DO NOT MODIFY THIS OR ANYTHING ABOVE! md5sum:fLrBvjW6JLyIRv59Qkrr+Q

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
