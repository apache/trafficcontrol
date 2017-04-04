use utf8;
package Schema::Result::Cdn;

# Created by DBIx::Class::Schema::Loader
# DO NOT MODIFY THE FIRST PART OF THIS FILE

=head1 NAME

Schema::Result::Cdn

=cut

use strict;
use warnings;

use base 'DBIx::Class::Core';

=head1 TABLE: C<cdn>

=cut

__PACKAGE__->table("cdn");

=head1 ACCESSORS

=head2 id

  data_type: 'bigint'
  is_auto_increment: 1
  is_nullable: 0
  sequence: 'cdn_id_seq'

=head2 name

  data_type: 'text'
  is_nullable: 0

=head2 last_updated

  data_type: 'timestamp with time zone'
  default_value: current_timestamp
  is_nullable: 0
  original: {default_value => \"now()"}

=head2 dnssec_enabled

  data_type: 'boolean'
  default_value: false
  is_nullable: 0

=head2 domain_name

  data_type: 'text'
  is_nullable: 0

=cut

__PACKAGE__->add_columns(
  "id",
  {
    data_type         => "bigint",
    is_auto_increment => 1,
    is_nullable       => 0,
    sequence          => "cdn_id_seq",
  },
  "name",
  { data_type => "text", is_nullable => 0 },
  "last_updated",
  {
    data_type     => "timestamp with time zone",
    default_value => \"current_timestamp",
    is_nullable   => 0,
    original      => { default_value => \"now()" },
  },
  "dnssec_enabled",
  { data_type => "boolean", default_value => \"false", is_nullable => 0 },
  "domain_name",
  { data_type => "text", is_nullable => 0 },
);

=head1 PRIMARY KEY

=over 4

=item * L</id>

=back

=cut

__PACKAGE__->set_primary_key("id");

=head1 UNIQUE CONSTRAINTS

=head2 C<idx_18210_cdn_cdn_unique>

=over 4

=item * L</name>

=back

=cut

__PACKAGE__->add_unique_constraint("idx_18210_cdn_cdn_unique", ["name"]);

=head1 RELATIONS

=head2 deliveryservices

Type: has_many

Related object: L<Schema::Result::Deliveryservice>

=cut

__PACKAGE__->has_many(
  "deliveryservices",
  "Schema::Result::Deliveryservice",
  { "foreign.cdn_id" => "self.id" },
  { cascade_copy => 0, cascade_delete => 0 },
);

=head2 profiles

Type: has_many

Related object: L<Schema::Result::Profile>

=cut

__PACKAGE__->has_many(
  "profiles",
  "Schema::Result::Profile",
  { "foreign.cdn" => "self.id" },
  { cascade_copy => 0, cascade_delete => 0 },
);

=head2 servers

Type: has_many

Related object: L<Schema::Result::Server>

=cut

__PACKAGE__->has_many(
  "servers",
  "Schema::Result::Server",
  { "foreign.cdn_id" => "self.id" },
  { cascade_copy => 0, cascade_delete => 0 },
);

=head2 snapshot

Type: might_have

Related object: L<Schema::Result::Snapshot>

=cut

__PACKAGE__->might_have(
  "snapshot",
  "Schema::Result::Snapshot",
  { "foreign.cdn" => "self.name" },
  { cascade_copy => 0, cascade_delete => 0 },
);


# Created by DBIx::Class::Schema::Loader v0.07046 @ 2016-12-26 14:46:31
# DO NOT MODIFY THIS OR ANYTHING ABOVE! md5sum:QlfY6K6vVxW5C9vsHc4YqA


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
