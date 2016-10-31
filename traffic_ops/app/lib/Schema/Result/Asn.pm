use utf8;
package Schema::Result::Asn;

# Created by DBIx::Class::Schema::Loader
# DO NOT MODIFY THE FIRST PART OF THIS FILE

=head1 NAME

Schema::Result::Asn

=cut

use strict;
use warnings;

use base 'DBIx::Class::Core';

=head1 TABLE: C<asn>

=cut

__PACKAGE__->table("asn");

=head1 ACCESSORS

=head2 id

  data_type: 'integer'
  is_auto_increment: 1
  is_nullable: 0

=head2 asn

  data_type: 'integer'
  is_nullable: 0

=head2 cachegroup

  data_type: 'integer'
  default_value: 0
  is_foreign_key: 1
  is_nullable: 0

=head2 last_updated

  data_type: 'timestamp'
  datetime_undef_if_invalid: 1
  default_value: current_timestamp
  is_nullable: 1

=cut

__PACKAGE__->add_columns(
  "id",
  { data_type => "integer", is_auto_increment => 1, is_nullable => 0 },
  "asn",
  { data_type => "integer", is_nullable => 0 },
  "cachegroup",
  {
    data_type      => "integer",
    default_value  => 0,
    is_foreign_key => 1,
    is_nullable    => 0,
  },
  "last_updated",
  {
    data_type => "timestamp",
    datetime_undef_if_invalid => 1,
    default_value => \"current_timestamp",
    is_nullable => 1,
  },
);

=head1 PRIMARY KEY

=over 4

=item * L</id>

=item * L</cachegroup>

=back

=cut

__PACKAGE__->set_primary_key("id", "cachegroup");

=head1 UNIQUE CONSTRAINTS

=head2 C<cr_id_UNIQUE>

=over 4

=item * L</id>

=back

=cut

__PACKAGE__->add_unique_constraint("cr_id_UNIQUE", ["id"]);

=head1 RELATIONS

=head2 cachegroup

Type: belongs_to

Related object: L<Schema::Result::Cachegroup>

=cut

__PACKAGE__->belongs_to(
  "cachegroup",
  "Schema::Result::Cachegroup",
  { id => "cachegroup" },
  { is_deferrable => 1, on_delete => "NO ACTION", on_update => "NO ACTION" },
);


# Created by DBIx::Class::Schema::Loader v0.07043 @ 2015-05-21 13:27:11
# DO NOT MODIFY THIS OR ANYTHING ABOVE! md5sum:WGDVL2+i4jTulNrU+p3cHg


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
