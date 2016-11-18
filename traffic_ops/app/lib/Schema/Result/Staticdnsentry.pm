use utf8;
package Schema::Result::Staticdnsentry;

# Created by DBIx::Class::Schema::Loader
# DO NOT MODIFY THE FIRST PART OF THIS FILE

=head1 NAME

Schema::Result::Staticdnsentry

=cut

use strict;
use warnings;

use base 'DBIx::Class::Core';

=head1 TABLE: C<staticdnsentry>

=cut

__PACKAGE__->table("staticdnsentry");

=head1 ACCESSORS

=head2 id

  data_type: 'bigint'
  is_auto_increment: 1
  is_nullable: 0
  sequence: 'staticdnsentry_id_seq'

=head2 host

  data_type: 'text'
  is_nullable: 0

=head2 address

  data_type: 'text'
  is_nullable: 0

=head2 type

  data_type: 'bigint'
  is_foreign_key: 1
  is_nullable: 0

=head2 ttl

  data_type: 'bigint'
  default_value: 3600
  is_nullable: 0

=head2 deliveryservice

  data_type: 'bigint'
  is_foreign_key: 1
  is_nullable: 0

=head2 cachegroup

  data_type: 'bigint'
  is_foreign_key: 1
  is_nullable: 1

=head2 last_updated

  data_type: 'timestamp with time zone'
  default_value: current_timestamp
  is_nullable: 1
  original: {default_value => \"now()"}

=cut

__PACKAGE__->add_columns(
  "id",
  {
    data_type         => "bigint",
    is_auto_increment => 1,
    is_nullable       => 0,
    sequence          => "staticdnsentry_id_seq",
  },
  "host",
  { data_type => "text", is_nullable => 0 },
  "address",
  { data_type => "text", is_nullable => 0 },
  "type",
  { data_type => "bigint", is_foreign_key => 1, is_nullable => 0 },
  "ttl",
  { data_type => "bigint", default_value => 3600, is_nullable => 0 },
  "deliveryservice",
  { data_type => "bigint", is_foreign_key => 1, is_nullable => 0 },
  "cachegroup",
  { data_type => "bigint", is_foreign_key => 1, is_nullable => 1 },
  "last_updated",
  {
    data_type     => "timestamp with time zone",
    default_value => \"current_timestamp",
    is_nullable   => 1,
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

=head2 C<idx_54505_combi_unique>

=over 4

=item * L</host>

=item * L</address>

=item * L</deliveryservice>

=item * L</cachegroup>

=back

=cut

__PACKAGE__->add_unique_constraint(
  "idx_54505_combi_unique",
  ["host", "address", "deliveryservice", "cachegroup"],
);

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
    on_delete     => "NO ACTION",
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
  { is_deferrable => 0, on_delete => "NO ACTION", on_update => "NO ACTION" },
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


# Created by DBIx::Class::Schema::Loader v0.07046 @ 2016-11-18 22:45:19
# DO NOT MODIFY THIS OR ANYTHING ABOVE! md5sum:g7WdTA+fuHr6rlpFQR6R7w


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
