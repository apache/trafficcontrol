use utf8;
package Schema::Result::JobStatus;

# Created by DBIx::Class::Schema::Loader
# DO NOT MODIFY THE FIRST PART OF THIS FILE

=head1 NAME

Schema::Result::JobStatus

=cut

use strict;
use warnings;

use base 'DBIx::Class::Core';

=head1 TABLE: C<job_status>

=cut

__PACKAGE__->table("job_status");

=head1 ACCESSORS

=head2 id

  data_type: 'bigint'
  is_auto_increment: 1
  is_nullable: 0
  sequence: 'job_status_id_seq'

=head2 name

  data_type: 'text'
  is_nullable: 1

=head2 description

  data_type: 'text'
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
    sequence          => "job_status_id_seq",
  },
  "name",
  { data_type => "text", is_nullable => 1 },
  "description",
  { data_type => "text", is_nullable => 1 },
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

=head1 UNIQUE CONSTRAINTS

=head2 C<job_status_name_unique>

=over 4

=item * L</name>

=back

=cut

__PACKAGE__->add_unique_constraint("job_status_name_unique", ["name"]);

=head1 RELATIONS

=head2 jobs

Type: has_many

Related object: L<Schema::Result::Job>

=cut

__PACKAGE__->has_many(
  "jobs",
  "Schema::Result::Job",
  { "foreign.status" => "self.id" },
  { cascade_copy => 0, cascade_delete => 0 },
);


# Created by DBIx::Class::Schema::Loader v0.07042 @ 2018-05-01 14:12:13
# DO NOT MODIFY THIS OR ANYTHING ABOVE! md5sum:8K80djhLKC61MmOhbhyQYA


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
