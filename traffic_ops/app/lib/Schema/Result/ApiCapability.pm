use utf8;
package Schema::Result::ApiCapability;

# Created by DBIx::Class::Schema::Loader
# DO NOT MODIFY THE FIRST PART OF THIS FILE

=head1 NAME

Schema::Result::ApiCapability

=cut

use strict;
use warnings;

use base 'DBIx::Class::Core';

=head1 TABLE: C<api_capability>

=cut

__PACKAGE__->table("api_capability");

=head1 ACCESSORS

=head2 id

  data_type: 'bigint'
  is_auto_increment: 1
  is_nullable: 0
  sequence: 'api_capability_id_seq'

=head2 http_method

  data_type: 'enum'
  extra: {custom_type_name => "http_method_t",list => ["GET","POST","PUT","PATCH","DELETE"]}
  is_nullable: 0

=head2 route

  data_type: 'text'
  is_nullable: 0

=head2 capability

  data_type: 'text'
  is_foreign_key: 1
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
    sequence          => "api_capability_id_seq",
  },
  "http_method",
  {
    data_type => "enum",
    extra => {
      custom_type_name => "http_method_t",
      list => ["GET", "POST", "PUT", "PATCH", "DELETE"],
    },
    is_nullable => 0,
  },
  "route",
  { data_type => "text", is_nullable => 0 },
  "capability",
  { data_type => "text", is_foreign_key => 1, is_nullable => 0 },
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

=head2 C<api_capability_http_method_route_capability_key>

=over 4

=item * L</http_method>

=item * L</route>

=item * L</capability>

=back

=cut

__PACKAGE__->add_unique_constraint(
  "api_capability_http_method_route_capability_key",
  ["http_method", "route", "capability"],
);

=head1 RELATIONS

=head2 capability

Type: belongs_to

Related object: L<Schema::Result::Capability>

=cut

__PACKAGE__->belongs_to(
  "capability",
  "Schema::Result::Capability",
  { name => "capability" },
  { is_deferrable => 0, on_delete => "RESTRICT", on_update => "NO ACTION" },
);


# Created by DBIx::Class::Schema::Loader v0.07046 @ 2017-05-21 10:15:00
# DO NOT MODIFY THIS OR ANYTHING ABOVE! md5sum:b1CNpOv08i47l8nNcqxLoA


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
