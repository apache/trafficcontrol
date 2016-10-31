use utf8;
package Schema::Result::PhysLocation;

# Created by DBIx::Class::Schema::Loader
# DO NOT MODIFY THE FIRST PART OF THIS FILE

=head1 NAME

Schema::Result::PhysLocation

=cut

use strict;
use warnings;

use base 'DBIx::Class::Core';

=head1 TABLE: C<phys_location>

=cut

__PACKAGE__->table("phys_location");

=head1 ACCESSORS

=head2 id

  data_type: 'integer'
  is_auto_increment: 1
  is_nullable: 0

=head2 name

  data_type: 'varchar'
  is_nullable: 0
  size: 45

=head2 short_name

  data_type: 'varchar'
  is_nullable: 0
  size: 12

=head2 address

  data_type: 'varchar'
  is_nullable: 0
  size: 128

=head2 city

  data_type: 'varchar'
  is_nullable: 0
  size: 128

=head2 state

  data_type: 'varchar'
  is_nullable: 0
  size: 2

=head2 zip

  data_type: 'varchar'
  is_nullable: 0
  size: 5

=head2 poc

  data_type: 'varchar'
  is_nullable: 1
  size: 128

=head2 phone

  data_type: 'varchar'
  is_nullable: 1
  size: 45

=head2 email

  data_type: 'varchar'
  is_nullable: 1
  size: 128

=head2 comments

  data_type: 'varchar'
  is_nullable: 1
  size: 256

=head2 region

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
  "id",
  { data_type => "integer", is_auto_increment => 1, is_nullable => 0 },
  "name",
  { data_type => "varchar", is_nullable => 0, size => 45 },
  "short_name",
  { data_type => "varchar", is_nullable => 0, size => 12 },
  "address",
  { data_type => "varchar", is_nullable => 0, size => 128 },
  "city",
  { data_type => "varchar", is_nullable => 0, size => 128 },
  "state",
  { data_type => "varchar", is_nullable => 0, size => 2 },
  "zip",
  { data_type => "varchar", is_nullable => 0, size => 5 },
  "poc",
  { data_type => "varchar", is_nullable => 1, size => 128 },
  "phone",
  { data_type => "varchar", is_nullable => 1, size => 45 },
  "email",
  { data_type => "varchar", is_nullable => 1, size => 128 },
  "comments",
  { data_type => "varchar", is_nullable => 1, size => 256 },
  "region",
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

=item * L</id>

=back

=cut

__PACKAGE__->set_primary_key("id");

=head1 UNIQUE CONSTRAINTS

=head2 C<name_UNIQUE>

=over 4

=item * L</name>

=back

=cut

__PACKAGE__->add_unique_constraint("name_UNIQUE", ["name"]);

=head2 C<short_name_UNIQUE>

=over 4

=item * L</short_name>

=back

=cut

__PACKAGE__->add_unique_constraint("short_name_UNIQUE", ["short_name"]);

=head1 RELATIONS

=head2 region

Type: belongs_to

Related object: L<Schema::Result::Region>

=cut

__PACKAGE__->belongs_to(
  "region",
  "Schema::Result::Region",
  { id => "region" },
  { is_deferrable => 1, on_delete => "NO ACTION", on_update => "NO ACTION" },
);

=head2 servers

Type: has_many

Related object: L<Schema::Result::Server>

=cut

__PACKAGE__->has_many(
  "servers",
  "Schema::Result::Server",
  { "foreign.phys_location" => "self.id" },
  { cascade_copy => 0, cascade_delete => 0 },
);


# Created by DBIx::Class::Schema::Loader v0.07043 @ 2015-05-21 13:27:11
# DO NOT MODIFY THIS OR ANYTHING ABOVE! md5sum:P49VpYuUIcx6aOOUORrQyA


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
