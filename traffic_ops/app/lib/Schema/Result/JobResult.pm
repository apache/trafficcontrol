use utf8;
package Schema::Result::JobResult;

# Created by DBIx::Class::Schema::Loader
# DO NOT MODIFY THE FIRST PART OF THIS FILE

=head1 NAME

Schema::Result::JobResult

=cut

use strict;
use warnings;

use base 'DBIx::Class::Core';

=head1 TABLE: C<job_result>

=cut

__PACKAGE__->table("job_result");

=head1 ACCESSORS

=head2 id

  data_type: 'integer'
  is_auto_increment: 1
  is_nullable: 0

=head2 job

  data_type: 'integer'
  is_foreign_key: 1
  is_nullable: 0

=head2 agent

  data_type: 'integer'
  is_foreign_key: 1
  is_nullable: 0

=head2 result

  data_type: 'varchar'
  is_nullable: 0
  size: 48

=head2 description

  data_type: 'varchar'
  is_nullable: 1
  size: 512

=head2 last_updated

  data_type: 'timestamp'
  datetime_undef_if_invalid: 1
  default_value: current_timestamp
  is_nullable: 1

=cut

__PACKAGE__->add_columns(
  "id",
  { data_type => "integer", is_auto_increment => 1, is_nullable => 0 },
  "job",
  { data_type => "integer", is_foreign_key => 1, is_nullable => 0 },
  "agent",
  { data_type => "integer", is_foreign_key => 1, is_nullable => 0 },
  "result",
  { data_type => "varchar", is_nullable => 0, size => 48 },
  "description",
  { data_type => "varchar", is_nullable => 1, size => 512 },
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

=head1 RELATIONS

=head2 agent

Type: belongs_to

Related object: L<Schema::Result::JobAgent>

=cut

__PACKAGE__->belongs_to(
  "agent",
  "Schema::Result::JobAgent",
  { id => "agent" },
  { is_deferrable => 1, on_delete => "CASCADE", on_update => "NO ACTION" },
);

=head2 job

Type: belongs_to

Related object: L<Schema::Result::Job>

=cut

__PACKAGE__->belongs_to(
  "job",
  "Schema::Result::Job",
  { id => "job" },
  { is_deferrable => 1, on_delete => "CASCADE", on_update => "NO ACTION" },
);


# Created by DBIx::Class::Schema::Loader v0.07043 @ 2015-05-21 13:27:11
# DO NOT MODIFY THIS OR ANYTHING ABOVE! md5sum:6+gwHGyMRYmILsJvuVcKyQ


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
