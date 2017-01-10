use utf8;
package Schema::Result::Servercheck;

# Created by DBIx::Class::Schema::Loader
# DO NOT MODIFY THE FIRST PART OF THIS FILE

=head1 NAME

Schema::Result::Servercheck

=cut

use strict;
use warnings;

use base 'DBIx::Class::Core';

=head1 TABLE: C<servercheck>

=cut

__PACKAGE__->table("servercheck");

=head1 ACCESSORS

=head2 id

  data_type: 'bigint'
  is_auto_increment: 1
  is_nullable: 0
  sequence: 'servercheck_id_seq'

=head2 server

  data_type: 'bigint'
  is_foreign_key: 1
  is_nullable: 0

=head2 aa

  data_type: 'bigint'
  is_nullable: 1

=head2 ab

  data_type: 'bigint'
  is_nullable: 1

=head2 ac

  data_type: 'bigint'
  is_nullable: 1

=head2 ad

  data_type: 'bigint'
  is_nullable: 1

=head2 ae

  data_type: 'bigint'
  is_nullable: 1

=head2 af

  data_type: 'bigint'
  is_nullable: 1

=head2 ag

  data_type: 'bigint'
  is_nullable: 1

=head2 ah

  data_type: 'bigint'
  is_nullable: 1

=head2 ai

  data_type: 'bigint'
  is_nullable: 1

=head2 aj

  data_type: 'bigint'
  is_nullable: 1

=head2 ak

  data_type: 'bigint'
  is_nullable: 1

=head2 al

  data_type: 'bigint'
  is_nullable: 1

=head2 am

  data_type: 'bigint'
  is_nullable: 1

=head2 an

  data_type: 'bigint'
  is_nullable: 1

=head2 ao

  data_type: 'bigint'
  is_nullable: 1

=head2 ap

  data_type: 'bigint'
  is_nullable: 1

=head2 aq

  data_type: 'bigint'
  is_nullable: 1

=head2 ar

  data_type: 'bigint'
  is_nullable: 1

=head2 bf

  data_type: 'bigint'
  is_nullable: 1

=head2 at

  data_type: 'bigint'
  is_nullable: 1

=head2 au

  data_type: 'bigint'
  is_nullable: 1

=head2 av

  data_type: 'bigint'
  is_nullable: 1

=head2 aw

  data_type: 'bigint'
  is_nullable: 1

=head2 ax

  data_type: 'bigint'
  is_nullable: 1

=head2 ay

  data_type: 'bigint'
  is_nullable: 1

=head2 az

  data_type: 'bigint'
  is_nullable: 1

=head2 ba

  data_type: 'bigint'
  is_nullable: 1

=head2 bb

  data_type: 'bigint'
  is_nullable: 1

=head2 bc

  data_type: 'bigint'
  is_nullable: 1

=head2 bd

  data_type: 'bigint'
  is_nullable: 1

=head2 be

  data_type: 'bigint'
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
    sequence          => "servercheck_id_seq",
  },
  "server",
  { data_type => "bigint", is_foreign_key => 1, is_nullable => 0 },
  "aa",
  { data_type => "bigint", is_nullable => 1 },
  "ab",
  { data_type => "bigint", is_nullable => 1 },
  "ac",
  { data_type => "bigint", is_nullable => 1 },
  "ad",
  { data_type => "bigint", is_nullable => 1 },
  "ae",
  { data_type => "bigint", is_nullable => 1 },
  "af",
  { data_type => "bigint", is_nullable => 1 },
  "ag",
  { data_type => "bigint", is_nullable => 1 },
  "ah",
  { data_type => "bigint", is_nullable => 1 },
  "ai",
  { data_type => "bigint", is_nullable => 1 },
  "aj",
  { data_type => "bigint", is_nullable => 1 },
  "ak",
  { data_type => "bigint", is_nullable => 1 },
  "al",
  { data_type => "bigint", is_nullable => 1 },
  "am",
  { data_type => "bigint", is_nullable => 1 },
  "an",
  { data_type => "bigint", is_nullable => 1 },
  "ao",
  { data_type => "bigint", is_nullable => 1 },
  "ap",
  { data_type => "bigint", is_nullable => 1 },
  "aq",
  { data_type => "bigint", is_nullable => 1 },
  "ar",
  { data_type => "bigint", is_nullable => 1 },
  "bf",
  { data_type => "bigint", is_nullable => 1 },
  "at",
  { data_type => "bigint", is_nullable => 1 },
  "au",
  { data_type => "bigint", is_nullable => 1 },
  "av",
  { data_type => "bigint", is_nullable => 1 },
  "aw",
  { data_type => "bigint", is_nullable => 1 },
  "ax",
  { data_type => "bigint", is_nullable => 1 },
  "ay",
  { data_type => "bigint", is_nullable => 1 },
  "az",
  { data_type => "bigint", is_nullable => 1 },
  "ba",
  { data_type => "bigint", is_nullable => 1 },
  "bb",
  { data_type => "bigint", is_nullable => 1 },
  "bc",
  { data_type => "bigint", is_nullable => 1 },
  "bd",
  { data_type => "bigint", is_nullable => 1 },
  "be",
  { data_type => "bigint", is_nullable => 1 },
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

=item * L</server>

=back

=cut

__PACKAGE__->set_primary_key("id", "server");

=head1 UNIQUE CONSTRAINTS

=head2 C<idx_54498_server>

=over 4

=item * L</server>

=back

=cut

__PACKAGE__->add_unique_constraint("idx_54498_server", ["server"]);

=head2 C<idx_54498_ses_id_unique>

=over 4

=item * L</id>

=back

=cut

__PACKAGE__->add_unique_constraint("idx_54498_ses_id_unique", ["id"]);

=head1 RELATIONS

=head2 server

Type: belongs_to

Related object: L<Schema::Result::Server>

=cut

__PACKAGE__->belongs_to(
  "server",
  "Schema::Result::Server",
  { id => "server" },
  { is_deferrable => 0, on_delete => "CASCADE", on_update => "NO ACTION" },
);


# Created by DBIx::Class::Schema::Loader v0.07046 @ 2016-11-18 22:45:19
# DO NOT MODIFY THIS OR ANYTHING ABOVE! md5sum:GgkEbC3GilmdyYWOguEY7A


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
