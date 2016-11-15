use utf8;
package Schema::Result::JobAgent;

# Created by DBIx::Class::Schema::Loader
# DO NOT MODIFY THE FIRST PART OF THIS FILE

=head1 NAME

Schema::Result::JobAgent

=cut

use strict;
use warnings;

use base 'DBIx::Class::Core';

=head1 TABLE: C<job_agent>

=cut

__PACKAGE__->table("job_agent");

=head1 ACCESSORS

=head2 id

  data_type: 'bigint'
  is_auto_increment: 1
  is_nullable: 0
  sequence: 'job_agent_id_seq'

=head2 name

  data_type: 'varchar'
  is_nullable: 1
  size: 128

=head2 description

  data_type: 'varchar'
  is_nullable: 1
  size: 512

=head2 active

  data_type: 'integer'
  default_value: 0
  is_nullable: 0

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
    sequence          => "job_agent_id_seq",
  },
  "name",
  { data_type => "varchar", is_nullable => 1, size => 128 },
  "description",
  { data_type => "varchar", is_nullable => 1, size => 512 },
  "active",
  { data_type => "integer", default_value => 0, is_nullable => 0 },
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

=back

=cut

__PACKAGE__->set_primary_key("id");

=head1 RELATIONS

=head2 job_results

Type: has_many

Related object: L<Schema::Result::JobResult>

=cut

__PACKAGE__->has_many(
  "job_results",
  "Schema::Result::JobResult",
  { "foreign.agent" => "self.id" },
  { cascade_copy => 0, cascade_delete => 0 },
);

=head2 jobs

Type: has_many

Related object: L<Schema::Result::Job>

=cut

__PACKAGE__->has_many(
  "jobs",
  "Schema::Result::Job",
  { "foreign.agent" => "self.id" },
  { cascade_copy => 0, cascade_delete => 0 },
);


# Created by DBIx::Class::Schema::Loader v0.07045 @ 2016-11-15 08:31:12
# DO NOT MODIFY THIS OR ANYTHING ABOVE! md5sum:1IIrmQijG0mO9uYkgiU55w


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
