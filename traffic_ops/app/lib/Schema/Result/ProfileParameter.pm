use utf8;
package Schema::Result::ProfileParameter;

# Created by DBIx::Class::Schema::Loader
# DO NOT MODIFY THE FIRST PART OF THIS FILE

=head1 NAME

Schema::Result::ProfileParameter

=cut

use strict;
use warnings;

use base 'DBIx::Class::Core';

=head1 TABLE: C<profile_parameter>

=cut

__PACKAGE__->table("profile_parameter");

=head1 ACCESSORS

=head2 profile

  data_type: 'bigint'
  is_foreign_key: 1
  is_nullable: 0

=head2 parameter

  data_type: 'bigint'
  is_foreign_key: 1
  is_nullable: 0

=head2 last_updated

  data_type: 'timestamp with time zone'
  default_value: current_timestamp
  is_nullable: 1
  original: {default_value => \"now()"}

=cut

__PACKAGE__->add_columns(
  "profile",
  { data_type => "bigint", is_foreign_key => 1, is_nullable => 0 },
  "parameter",
  { data_type => "bigint", is_foreign_key => 1, is_nullable => 0 },
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

=item * L</profile>

=item * L</parameter>

=back

=cut

__PACKAGE__->set_primary_key("profile", "parameter");

=head1 RELATIONS

=head2 parameter

Type: belongs_to

Related object: L<Schema::Result::Parameter>

=cut

__PACKAGE__->belongs_to(
  "parameter",
  "Schema::Result::Parameter",
  { id => "parameter" },
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
  { is_deferrable => 0, on_delete => "CASCADE", on_update => "NO ACTION" },
);


# Created by DBIx::Class::Schema::Loader v0.07045 @ 2016-11-15 09:35:47
# DO NOT MODIFY THIS OR ANYTHING ABOVE! md5sum:bg3IAJfz8GX6uS5UXMPblA


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
