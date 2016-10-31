use utf8;
package Schema::Result::ToExtension;

# Created by DBIx::Class::Schema::Loader
# DO NOT MODIFY THE FIRST PART OF THIS FILE

=head1 NAME

Schema::Result::ToExtension

=cut

use strict;
use warnings;

use base 'DBIx::Class::Core';

=head1 TABLE: C<to_extension>

=cut

__PACKAGE__->table("to_extension");

=head1 ACCESSORS

=head2 id

  data_type: 'integer'
  is_auto_increment: 1
  is_nullable: 0

=head2 name

  data_type: 'varchar'
  is_nullable: 0
  size: 45

=head2 version

  data_type: 'varchar'
  is_nullable: 0
  size: 45

=head2 info_url

  data_type: 'varchar'
  is_nullable: 0
  size: 45

=head2 script_file

  data_type: 'varchar'
  is_nullable: 0
  size: 45

=head2 isactive

  data_type: 'tinyint'
  is_nullable: 0

=head2 additional_config_json

  data_type: 'varchar'
  is_nullable: 1
  size: 4096

=head2 description

  data_type: 'varchar'
  is_nullable: 1
  size: 4096

=head2 servercheck_short_name

  data_type: 'varchar'
  is_nullable: 1
  size: 8

=head2 servercheck_column_name

  data_type: 'varchar'
  is_nullable: 1
  size: 10

=head2 type

  data_type: 'integer'
  is_foreign_key: 1
  is_nullable: 0

=head2 last_updated

  data_type: 'timestamp'
  datetime_undef_if_invalid: 1
  default_value: current_timestamp
  is_nullable: 0

=cut

__PACKAGE__->add_columns(
  "id",
  { data_type => "integer", is_auto_increment => 1, is_nullable => 0 },
  "name",
  { data_type => "varchar", is_nullable => 0, size => 45 },
  "version",
  { data_type => "varchar", is_nullable => 0, size => 45 },
  "info_url",
  { data_type => "varchar", is_nullable => 0, size => 45 },
  "script_file",
  { data_type => "varchar", is_nullable => 0, size => 45 },
  "isactive",
  { data_type => "tinyint", is_nullable => 0 },
  "additional_config_json",
  { data_type => "varchar", is_nullable => 1, size => 4096 },
  "description",
  { data_type => "varchar", is_nullable => 1, size => 4096 },
  "servercheck_short_name",
  { data_type => "varchar", is_nullable => 1, size => 8 },
  "servercheck_column_name",
  { data_type => "varchar", is_nullable => 1, size => 10 },
  "type",
  { data_type => "integer", is_foreign_key => 1, is_nullable => 0 },
  "last_updated",
  {
    data_type => "timestamp",
    datetime_undef_if_invalid => 1,
    default_value => \"current_timestamp",
    is_nullable => 0,
  },
);

=head1 PRIMARY KEY

=over 4

=item * L</id>

=back

=cut

__PACKAGE__->set_primary_key("id");

=head1 RELATIONS

=head2 type

Type: belongs_to

Related object: L<Schema::Result::Type>

=cut

__PACKAGE__->belongs_to(
  "type",
  "Schema::Result::Type",
  { id => "type" },
  { is_deferrable => 1, on_delete => "NO ACTION", on_update => "NO ACTION" },
);


# Created by DBIx::Class::Schema::Loader v0.07043 @ 2015-05-21 13:27:11
# DO NOT MODIFY THIS OR ANYTHING ABOVE! md5sum:/YZDsslpM0Bp0vcpV6WEMw


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
