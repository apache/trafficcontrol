use utf8;
package Schema::Result::Coordinate;

# Created by DBIx::Class::Schema::Loader
# DO NOT MODIFY THE FIRST PART OF THIS FILE

=head1 NAME

Schema::Result::Coordinate

=cut

use strict;
use warnings;

use base 'DBIx::Class::Core';

=head1 TABLE: C<coordinate>

=cut

__PACKAGE__->table("coordinate");

=head1 ACCESSORS

=head2 id

  data_type: 'bigint'
  is_auto_increment: 1
  is_nullable: 0
  sequence: 'coordinate_id_seq'

=head2 name

  data_type: 'text'
  is_nullable: 0

=head2 latitude

  data_type: 'numeric'
  default_value: 0.0
  is_nullable: 0

=head2 longitude

  data_type: 'numeric'
  default_value: 0.0
  is_nullable: 0

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
    sequence          => "coordinate_id_seq",
  },
  "name",
  { data_type => "text", is_nullable => 0 },
  "latitude",
  { data_type => "numeric", default_value => "0.0", is_nullable => 0 },
  "longitude",
  { data_type => "numeric", default_value => "0.0", is_nullable => 0 },
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

=head2 C<coordinate_name_key>

=over 4

=item * L</name>

=back

=cut

__PACKAGE__->add_unique_constraint("coordinate_name_key", ["name"]);

=head1 RELATIONS

=head2 cachegroups

Type: has_many

Related object: L<Schema::Result::Cachegroup>

=cut

__PACKAGE__->has_many(
  "cachegroups",
  "Schema::Result::Cachegroup",
  { "foreign.coordinate" => "self.id" },
  { cascade_copy => 0, cascade_delete => 0 },
);

=head2 origins

Type: has_many

Related object: L<Schema::Result::Origin>

=cut

__PACKAGE__->has_many(
  "origins",
  "Schema::Result::Origin",
  { "foreign.coordinate" => "self.id" },
  { cascade_copy => 0, cascade_delete => 0 },
);


# Created by DBIx::Class::Schema::Loader v0.07042 @ 2018-06-27 16:34:28
# DO NOT MODIFY THIS OR ANYTHING ABOVE! md5sum:TdhJL1P7uk/07Oz2Y73Plw

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
