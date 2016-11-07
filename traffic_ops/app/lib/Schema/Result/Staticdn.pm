use utf8;
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
#
package Schema::Result::Staticdn;

# Created by DBIx::Class::Schema::Loader
# DO NOT MODIFY THE FIRST PART OF THIS FILE

=head1 NAME

Schema::Result::Staticdn

=cut

use strict;
use warnings;

use base 'DBIx::Class::Core';

=head1 TABLE: C<staticdns>

=cut

__PACKAGE__->table("staticdns");

=head1 ACCESSORS

=head2 id

  data_type: 'integer'
  is_auto_increment: 1
  is_nullable: 0

=head2 host

  data_type: 'varchar'
  is_nullable: 0
  size: 45

=head2 address

  data_type: 'varchar'
  is_nullable: 0
  size: 45

=head2 type

  data_type: 'integer'
  is_foreign_key: 1
  is_nullable: 0

=head2 ttl

  data_type: 'integer'
  default_value: 3600
  is_nullable: 0

=cut

__PACKAGE__->add_columns(
  "id",
  { data_type => "integer", is_auto_increment => 1, is_nullable => 0 },
  "host",
  { data_type => "varchar", is_nullable => 0, size => 45 },
  "address",
  { data_type => "varchar", is_nullable => 0, size => 45 },
  "type",
  { data_type => "integer", is_foreign_key => 1, is_nullable => 0 },
  "ttl",
  { data_type => "integer", default_value => 3600, is_nullable => 0 },
);

=head1 PRIMARY KEY

=over 4

=item * L</id>

=back

=cut

__PACKAGE__->set_primary_key("id");

=head1 UNIQUE CONSTRAINTS

=head2 C<host_name_UNIQUE>

=over 4

=item * L</host>

=back

=cut

__PACKAGE__->add_unique_constraint("host_name_UNIQUE", ["host"]);

=head1 RELATIONS

=head2 staticdns_deliveryservices

Type: has_many

Related object: L<Schema::Result::StaticdnsDeliveryservice>

=cut

__PACKAGE__->has_many(
  "staticdns_deliveryservices",
  "Schema::Result::StaticdnsDeliveryservice",
  { "foreign.staticdns" => "self.id" },
  { cascade_copy => 0, cascade_delete => 0 },
);

=head2 staticdns_locations

Type: has_many

Related object: L<Schema::Result::StaticdnsLocation>

=cut

__PACKAGE__->has_many(
  "staticdns_locations",
  "Schema::Result::StaticdnsLocation",
  { "foreign.staticdns" => "self.id" },
  { cascade_copy => 0, cascade_delete => 0 },
);

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

=head2 deliveryservices

Type: many_to_many

Composing rels: L</staticdns_deliveryservices> -> deliveryservice

=cut

__PACKAGE__->many_to_many(
  "deliveryservices",
  "staticdns_deliveryservices",
  "deliveryservice",
);

=head2 locations

Type: many_to_many

Composing rels: L</staticdns_locations> -> location

=cut

__PACKAGE__->many_to_many("locations", "staticdns_locations", "location");


# Created by DBIx::Class::Schema::Loader v0.07036 @ 2013-10-18 16:59:25
# DO NOT MODIFY THIS OR ANYTHING ABOVE! md5sum:C3Jv7K0u4um1s/Y1PH0ljA


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
