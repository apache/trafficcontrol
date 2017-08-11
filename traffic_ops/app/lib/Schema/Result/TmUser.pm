use utf8;
package Schema::Result::TmUser;

# Created by DBIx::Class::Schema::Loader
# DO NOT MODIFY THE FIRST PART OF THIS FILE

=head1 NAME

Schema::Result::TmUser

=cut

use strict;
use warnings;

use base 'DBIx::Class::Core';

=head1 TABLE: C<tm_user>

=cut

__PACKAGE__->table("tm_user");

=head1 ACCESSORS

=head2 id

  data_type: 'bigint'
  is_auto_increment: 1
  is_nullable: 0
  sequence: 'tm_user_id_seq'

=head2 username

  data_type: 'text'
  is_nullable: 1

=head2 public_ssh_key

  data_type: 'text'
  is_nullable: 1

=head2 role

  data_type: 'bigint'
  is_foreign_key: 1
  is_nullable: 1

=head2 uid

  data_type: 'bigint'
  is_nullable: 1

=head2 gid

  data_type: 'bigint'
  is_nullable: 1

=head2 local_passwd

  data_type: 'text'
  is_nullable: 1

=head2 confirm_local_passwd

  data_type: 'text'
  is_nullable: 1

=head2 last_updated

  data_type: 'timestamp with time zone'
  default_value: current_timestamp
  is_nullable: 1
  original: {default_value => \"now()"}

=head2 company

  data_type: 'text'
  is_nullable: 1

=head2 email

  data_type: 'text'
  is_nullable: 1

=head2 full_name

  data_type: 'text'
  is_nullable: 1

=head2 new_user

  data_type: 'boolean'
  default_value: false
  is_nullable: 0

=head2 address_line1

  data_type: 'text'
  is_nullable: 1

=head2 address_line2

  data_type: 'text'
  is_nullable: 1

=head2 city

  data_type: 'text'
  is_nullable: 1

=head2 state_or_province

  data_type: 'text'
  is_nullable: 1

=head2 phone_number

  data_type: 'text'
  is_nullable: 1

=head2 postal_code

  data_type: 'text'
  is_nullable: 1

=head2 country

  data_type: 'text'
  is_nullable: 1

=head2 token

  data_type: 'text'
  is_nullable: 1

=head2 registration_sent

  data_type: 'timestamp with time zone'
  is_nullable: 1

=head2 tenant_id

  data_type: 'bigint'
  is_foreign_key: 1
  is_nullable: 1

=cut

__PACKAGE__->add_columns(
  "id",
  {
    data_type         => "bigint",
    is_auto_increment => 1,
    is_nullable       => 0,
    sequence          => "tm_user_id_seq",
  },
  "username",
  { data_type => "text", is_nullable => 1 },
  "public_ssh_key",
  { data_type => "text", is_nullable => 1 },
  "role",
  { data_type => "bigint", is_foreign_key => 1, is_nullable => 1 },
  "uid",
  { data_type => "bigint", is_nullable => 1 },
  "gid",
  { data_type => "bigint", is_nullable => 1 },
  "local_passwd",
  { data_type => "text", is_nullable => 1 },
  "confirm_local_passwd",
  { data_type => "text", is_nullable => 1 },
  "last_updated",
  {
    data_type     => "timestamp with time zone",
    default_value => \"current_timestamp",
    is_nullable   => 1,
    original      => { default_value => \"now()" },
  },
  "company",
  { data_type => "text", is_nullable => 1 },
  "email",
  { data_type => "text", is_nullable => 1 },
  "full_name",
  { data_type => "text", is_nullable => 1 },
  "new_user",
  { data_type => "boolean", default_value => \"false", is_nullable => 0 },
  "address_line1",
  { data_type => "text", is_nullable => 1 },
  "address_line2",
  { data_type => "text", is_nullable => 1 },
  "city",
  { data_type => "text", is_nullable => 1 },
  "state_or_province",
  { data_type => "text", is_nullable => 1 },
  "phone_number",
  { data_type => "text", is_nullable => 1 },
  "postal_code",
  { data_type => "text", is_nullable => 1 },
  "country",
  { data_type => "text", is_nullable => 1 },
  "token",
  { data_type => "text", is_nullable => 1 },
  "registration_sent",
  { data_type => "timestamp with time zone", is_nullable => 1 },
  "tenant_id",
  { data_type => "bigint", is_foreign_key => 1, is_nullable => 1 },
);

=head1 PRIMARY KEY

=over 4

=item * L</id>

=back

=cut

__PACKAGE__->set_primary_key("id");

=head1 UNIQUE CONSTRAINTS

=head2 C<idx_89765_tmuser_email_unique>

=over 4

=item * L</email>

=back

=cut

__PACKAGE__->add_unique_constraint("idx_89765_tmuser_email_unique", ["email"]);

=head2 C<idx_89765_username_unique>

=over 4

=item * L</username>

=back

=cut

__PACKAGE__->add_unique_constraint("idx_89765_username_unique", ["username"]);

=head1 RELATIONS

=head2 deliveryservice_tmusers

Type: has_many

Related object: L<Schema::Result::DeliveryserviceTmuser>

=cut

__PACKAGE__->has_many(
  "deliveryservice_tmusers",
  "Schema::Result::DeliveryserviceTmuser",
  { "foreign.tm_user_id" => "self.id" },
  { cascade_copy => 0, cascade_delete => 0 },
);

=head2 federation_tmusers

Type: has_many

Related object: L<Schema::Result::FederationTmuser>

=cut

__PACKAGE__->has_many(
  "federation_tmusers",
  "Schema::Result::FederationTmuser",
  { "foreign.tm_user" => "self.id" },
  { cascade_copy => 0, cascade_delete => 0 },
);

=head2 jobs

Type: has_many

Related object: L<Schema::Result::Job>

=cut

__PACKAGE__->has_many(
  "jobs",
  "Schema::Result::Job",
  { "foreign.job_user" => "self.id" },
  { cascade_copy => 0, cascade_delete => 0 },
);

=head2 logs

Type: has_many

Related object: L<Schema::Result::Log>

=cut

__PACKAGE__->has_many(
  "logs",
  "Schema::Result::Log",
  { "foreign.tm_user" => "self.id" },
  { cascade_copy => 0, cascade_delete => 0 },
);

=head2 role

Type: belongs_to

Related object: L<Schema::Result::Role>

=cut

__PACKAGE__->belongs_to(
  "role",
  "Schema::Result::Role",
  { id => "role" },
  {
    is_deferrable => 0,
    join_type     => "LEFT",
    on_delete     => "SET NULL",
    on_update     => "NO ACTION",
  },
);

=head2 tenant

Type: belongs_to

Related object: L<Schema::Result::Tenant>

=cut

__PACKAGE__->belongs_to(
  "tenant",
  "Schema::Result::Tenant",
  { id => "tenant_id" },
  {
    is_deferrable => 0,
    join_type     => "LEFT",
    on_delete     => "NO ACTION",
    on_update     => "NO ACTION",
  },
);


# Created by DBIx::Class::Schema::Loader v0.07046 @ 2017-02-19 10:20:47
# DO NOT MODIFY THIS OR ANYTHING ABOVE! md5sum:2lI3iG0t7INKH+xQq+lo9g


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
