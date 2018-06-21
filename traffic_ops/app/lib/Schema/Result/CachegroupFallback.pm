use utf8;
package Schema::Result::CachegroupFallback;

# Created by DBIx::Class::Schema::Loader
# DO NOT MODIFY THE FIRST PART OF THIS FILE

=head1 NAME

Schema::Result::CachegroupFallback

=cut

use strict;
use warnings;

use base 'DBIx::Class::Core';

=head1 TABLE: C<cachegroup_fallbacks>

=cut

__PACKAGE__->table("cachegroup_fallbacks");

=head1 ACCESSORS

=head2 primary_cg

  data_type: 'bigint'
  is_foreign_key: 1
  is_nullable: 1

=head2 backup_cg

  data_type: 'bigint'
  is_foreign_key: 1
  is_nullable: 1

=head2 set_order

  data_type: 'bigint'
  is_nullable: 0

=cut

__PACKAGE__->add_columns(
  "primary_cg",
  { data_type => "bigint", is_foreign_key => 1, is_nullable => 1 },
  "backup_cg",
  { data_type => "bigint", is_foreign_key => 1, is_nullable => 1 },
  "set_order",
  { data_type => "bigint", is_nullable => 0 },
);

=head1 UNIQUE CONSTRAINTS

=head2 C<cachegroup_fallbacks_primary_cg_backup_cg_key>

=over 4

=item * L</primary_cg>

=item * L</backup_cg>

=back

=cut

__PACKAGE__->add_unique_constraint(
  "cachegroup_fallbacks_primary_cg_backup_cg_key",
  ["primary_cg", "backup_cg"],
);

=head2 C<cachegroup_fallbacks_primary_cg_set_order_key>

=over 4

=item * L</primary_cg>

=item * L</set_order>

=back

=cut

__PACKAGE__->add_unique_constraint(
  "cachegroup_fallbacks_primary_cg_set_order_key",
  ["primary_cg", "set_order"],
);

=head1 RELATIONS

=head2 backup_cg

Type: belongs_to

Related object: L<Schema::Result::Cachegroup>

=cut

__PACKAGE__->belongs_to(
  "backup_cg",
  "Schema::Result::Cachegroup",
  { id => "backup_cg" },
  {
    is_deferrable => 0,
    join_type     => "LEFT",
    on_delete     => "CASCADE",
    on_update     => "NO ACTION",
  },
);

=head2 primary_cg

Type: belongs_to

Related object: L<Schema::Result::Cachegroup>

=cut

__PACKAGE__->belongs_to(
  "primary_cg",
  "Schema::Result::Cachegroup",
  { id => "primary_cg" },
  {
    is_deferrable => 0,
    join_type     => "LEFT",
    on_delete     => "CASCADE",
    on_update     => "NO ACTION",
  },
);


# Created by DBIx::Class::Schema::Loader v0.07048 @ 2018-03-20 04:15:32
# DO NOT MODIFY THIS OR ANYTHING ABOVE! md5sum:9bJ/JA5FNpy0LYu1KRdQqA

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
