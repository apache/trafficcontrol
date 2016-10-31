use utf8;
package Schema::Result::CachegroupParameter;

# Created by DBIx::Class::Schema::Loader
# DO NOT MODIFY THE FIRST PART OF THIS FILE

=head1 NAME

Schema::Result::CachegroupParameter

=cut

use strict;
use warnings;

use base 'DBIx::Class::Core';

=head1 TABLE: C<cachegroup_parameter>

=cut

__PACKAGE__->table("cachegroup_parameter");

=head1 ACCESSORS

=head2 cachegroup

  data_type: 'integer'
  default_value: 0
  is_foreign_key: 1
  is_nullable: 0

=head2 parameter

  data_type: 'integer'
  is_foreign_key: 1
  is_nullable: 0

=head2 last_updated

  data_type: 'timestamp'
  datetime_undef_if_invalid: 1
  default_value: current_timestamp
  is_nullable: 1

=cut

__PACKAGE__->add_columns(
  "cachegroup",
  {
    data_type      => "integer",
    default_value  => 0,
    is_foreign_key => 1,
    is_nullable    => 0,
  },
  "parameter",
  { data_type => "integer", is_foreign_key => 1, is_nullable => 0 },
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

=item * L</cachegroup>

=item * L</parameter>

=back

=cut

__PACKAGE__->set_primary_key("cachegroup", "parameter");

=head1 RELATIONS

=head2 cachegroup

Type: belongs_to

Related object: L<Schema::Result::Cachegroup>

=cut

__PACKAGE__->belongs_to(
  "cachegroup",
  "Schema::Result::Cachegroup",
  { id => "cachegroup" },
  { is_deferrable => 1, on_delete => "CASCADE", on_update => "NO ACTION" },
);

=head2 parameter

Type: belongs_to

Related object: L<Schema::Result::Parameter>

=cut

__PACKAGE__->belongs_to(
  "parameter",
  "Schema::Result::Parameter",
  { id => "parameter" },
  { is_deferrable => 1, on_delete => "CASCADE", on_update => "NO ACTION" },
);


# Created by DBIx::Class::Schema::Loader v0.07043 @ 2015-05-21 13:27:11
# DO NOT MODIFY THIS OR ANYTHING ABOVE! md5sum:BRGhe7It8jc9lvi87bvjHw


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
