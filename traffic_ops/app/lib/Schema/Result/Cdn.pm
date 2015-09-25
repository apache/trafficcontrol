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

  data_type: 'integer'
  is_auto_increment: 1
  is_nullable: 0

=head2 name

  data_type: 'varchar'
  is_nullable: 0
  size: 1024

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
  { data_type => "varchar", is_nullable => 0, size => 1024 },
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
  { "foreign.cdn_id" => "self.id" },
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


# Created by DBIx::Class::Schema::Loader v0.07043 @ 2015-09-25 10:08:25
# DO NOT MODIFY THIS OR ANYTHING ABOVE! md5sum:SS/GRjmssf2s5XhUu2rEvw


# You can replace this text with custom code or comments, and it will be preserved on regeneration
#
# Copyright 2015 Comcast Cable Communications Management, LLC
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
