use utf8;
package Schema::Result::Cachegroup;

# Created by DBIx::Class::Schema::Loader
# DO NOT MODIFY THE FIRST PART OF THIS FILE

=head1 NAME

Schema::Result::Cachegroup

=cut

use strict;
use warnings;

use base 'DBIx::Class::Core';

=head1 TABLE: C<cachegroup>

=cut

__PACKAGE__->table("cachegroup");

=head1 ACCESSORS

=head2 id

  data_type: 'bigint'
  is_auto_increment: 1
  is_nullable: 0
  sequence: 'cachegroup_id_seq'

=head2 name

  data_type: 'text'
  is_nullable: 0

=head2 short_name

  data_type: 'text'
  is_nullable: 0

=head2 parent_cachegroup_id

  data_type: 'bigint'
  is_foreign_key: 1
  is_nullable: 1

=head2 secondary_parent_cachegroup_id

  data_type: 'bigint'
  is_foreign_key: 1
  is_nullable: 1

=head2 type

  data_type: 'bigint'
  is_foreign_key: 1
  is_nullable: 0

=head2 last_updated

  data_type: 'timestamp with time zone'
  default_value: current_timestamp
  is_nullable: 0
  original: {default_value => \"now()"}

=head2 fallback_to_closest

  data_type: 'boolean'
  default_value: true
  is_nullable: 1

=head2 coordinate

  data_type: 'bigint'
  is_foreign_key: 1
  is_nullable: 1

=cut

__PACKAGE__->add_columns(
  "id",
  {
    data_type         => "bigint",
    is_auto_increment => 1,
    is_nullable       => 0,
    sequence          => "cachegroup_id_seq",
  },
  "name",
  { data_type => "text", is_nullable => 0 },
  "short_name",
  { data_type => "text", is_nullable => 0 },
  "parent_cachegroup_id",
  { data_type => "bigint", is_foreign_key => 1, is_nullable => 1 },
  "secondary_parent_cachegroup_id",
  { data_type => "bigint", is_foreign_key => 1, is_nullable => 1 },
  "type",
  { data_type => "bigint", is_foreign_key => 1, is_nullable => 0 },
  "last_updated",
  {
    data_type     => "timestamp with time zone",
    default_value => \"current_timestamp",
    is_nullable   => 0,
    original      => { default_value => \"now()" },
  },
  "fallback_to_closest",
  { data_type => "boolean", default_value => \"true", is_nullable => 1 },
  "coordinate",
  { data_type => "bigint", is_foreign_key => 1, is_nullable => 1 },
);

=head1 PRIMARY KEY

=over 4

=item * L</id>

=item * L</type>

=back

=cut

__PACKAGE__->set_primary_key("id", "type");

=head1 UNIQUE CONSTRAINTS

=head2 C<idx_140208_cg_name_unique>

=over 4

=item * L</name>

=back

=cut

__PACKAGE__->add_unique_constraint("idx_140208_cg_name_unique", ["name"]);

=head2 C<idx_140208_cg_short_unique>

=over 4

=item * L</short_name>

=back

=cut

__PACKAGE__->add_unique_constraint("idx_140208_cg_short_unique", ["short_name"]);

=head2 C<idx_140208_lo_id_unique>

=over 4

=item * L</id>

=back

=cut

__PACKAGE__->add_unique_constraint("idx_140208_lo_id_unique", ["id"]);

=head1 RELATIONS

=head2 asns

Type: has_many

Related object: L<Schema::Result::Asn>

=cut

__PACKAGE__->has_many(
  "asns",
  "Schema::Result::Asn",
  { "foreign.cachegroup" => "self.id" },
  { cascade_copy => 0, cascade_delete => 0 },
);

=head2 cachegroup_fallbacks_backup_cgs

Type: has_many

Related object: L<Schema::Result::CachegroupFallback>

=cut

__PACKAGE__->has_many(
  "cachegroup_fallbacks_backup_cgs",
  "Schema::Result::CachegroupFallback",
  { "foreign.backup_cg" => "self.id" },
  { cascade_copy => 0, cascade_delete => 0 },
);

=head2 cachegroup_fallbacks_primary_cgs

Type: has_many

Related object: L<Schema::Result::CachegroupFallback>

=cut

__PACKAGE__->has_many(
  "cachegroup_fallbacks_primary_cgs",
  "Schema::Result::CachegroupFallback",
  { "foreign.primary_cg" => "self.id" },
  { cascade_copy => 0, cascade_delete => 0 },
);

=head2 cachegroup_parameters

Type: has_many

Related object: L<Schema::Result::CachegroupParameter>

=cut

__PACKAGE__->has_many(
  "cachegroup_parameters",
  "Schema::Result::CachegroupParameter",
  { "foreign.cachegroup" => "self.id" },
  { cascade_copy => 0, cascade_delete => 0 },
);

=head2 cachegroup_secondary_parent_cachegroups

Type: has_many

Related object: L<Schema::Result::Cachegroup>

=cut

__PACKAGE__->has_many(
  "cachegroup_secondary_parent_cachegroups",
  "Schema::Result::Cachegroup",
  { "foreign.secondary_parent_cachegroup_id" => "self.id" },
  { cascade_copy => 0, cascade_delete => 0 },
);

=head2 cachegroups

Type: has_many

Related object: L<Schema::Result::Cachegroup>

=cut

__PACKAGE__->has_many(
  "cachegroups",
  "Schema::Result::Cachegroup",
  { "foreign.parent_cachegroup_id" => "self.id" },
  { cascade_copy => 0, cascade_delete => 0 },
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
    on_delete     => "NO ACTION",
    on_update     => "NO ACTION",
  },
);

=head2 origins

Type: has_many

Related object: L<Schema::Result::Origin>

=cut

__PACKAGE__->has_many(
  "origins",
  "Schema::Result::Origin",
  { "foreign.cachegroup" => "self.id" },
  { cascade_copy => 0, cascade_delete => 0 },
);

=head2 parent_cachegroup

Type: belongs_to

Related object: L<Schema::Result::Cachegroup>

=cut

__PACKAGE__->belongs_to(
  "parent_cachegroup",
  "Schema::Result::Cachegroup",
  { id => "parent_cachegroup_id" },
  {
    is_deferrable => 0,
    join_type     => "LEFT",
    on_delete     => "NO ACTION",
    on_update     => "NO ACTION",
  },
);

=head2 secondary_parent_cachegroup

Type: belongs_to

Related object: L<Schema::Result::Cachegroup>

=cut

__PACKAGE__->belongs_to(
  "secondary_parent_cachegroup",
  "Schema::Result::Cachegroup",
  { id => "secondary_parent_cachegroup_id" },
  {
    is_deferrable => 0,
    join_type     => "LEFT",
    on_delete     => "NO ACTION",
    on_update     => "NO ACTION",
  },
);

=head2 servers

Type: has_many

Related object: L<Schema::Result::Server>

=cut

__PACKAGE__->has_many(
  "servers",
  "Schema::Result::Server",
  { "foreign.cachegroup" => "self.id" },
  { cascade_copy => 0, cascade_delete => 0 },
);

=head2 staticdnsentries

Type: has_many

Related object: L<Schema::Result::Staticdnsentry>

=cut

__PACKAGE__->has_many(
  "staticdnsentries",
  "Schema::Result::Staticdnsentry",
  { "foreign.cachegroup" => "self.id" },
  { cascade_copy => 0, cascade_delete => 0 },
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


# Created by DBIx::Class::Schema::Loader v0.07042 @ 2018-06-27 16:34:28
# DO NOT MODIFY THIS OR ANYTHING ABOVE! md5sum:Xv7TVScj3TvTuzk/Gd9Mug

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
