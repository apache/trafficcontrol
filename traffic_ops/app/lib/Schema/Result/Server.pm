use utf8;
package Schema::Result::Server;

# Created by DBIx::Class::Schema::Loader
# DO NOT MODIFY THE FIRST PART OF THIS FILE

=head1 NAME

Schema::Result::Server

=cut

use strict;
use warnings;

use base 'DBIx::Class::Core';

=head1 TABLE: C<server>

=cut

__PACKAGE__->table("server");

=head1 ACCESSORS

=head2 id

  data_type: 'bigint'
  is_auto_increment: 1
  is_nullable: 0
  sequence: 'server_id_seq'

=head2 host_name

  data_type: 'varchar'
  is_nullable: 0
  size: 63

=head2 domain_name

  data_type: 'varchar'
  is_nullable: 0
  size: 63

=head2 tcp_port

  data_type: 'bigint'
  is_nullable: 1

=head2 xmpp_id

  data_type: 'varchar'
  is_nullable: 1
  size: 256

=head2 xmpp_passwd

  data_type: 'varchar'
  is_nullable: 1
  size: 45

=head2 interface_name

  data_type: 'varchar'
  is_nullable: 0
  size: 45

=head2 ip_address

  data_type: 'varchar'
  is_nullable: 0
  size: 45

=head2 ip_netmask

  data_type: 'varchar'
  is_nullable: 0
  size: 45

=head2 ip_gateway

  data_type: 'varchar'
  is_nullable: 0
  size: 45

=head2 ip6_address

  data_type: 'varchar'
  is_nullable: 1
  size: 50

=head2 ip6_gateway

  data_type: 'varchar'
  is_nullable: 1
  size: 50

=head2 interface_mtu

  data_type: 'bigint'
  default_value: 9000
  is_nullable: 0

=head2 phys_location

  data_type: 'bigint'
  is_foreign_key: 1
  is_nullable: 0

=head2 rack

  data_type: 'varchar'
  is_nullable: 1
  size: 64

=head2 cachegroup

  data_type: 'bigint'
  default_value: 0
  is_foreign_key: 1
  is_nullable: 0

=head2 type

  data_type: 'bigint'
  is_foreign_key: 1
  is_nullable: 0

=head2 status

  data_type: 'bigint'
  is_foreign_key: 1
  is_nullable: 0

=head2 offline_reason

  data_type: 'varchar'
  is_nullable: 1
  size: 256

=head2 upd_pending

  data_type: 'boolean'
  default_value: false
  is_nullable: 0

=head2 profile

  data_type: 'bigint'
  is_foreign_key: 1
  is_nullable: 0

=head2 cdn_id

  data_type: 'bigint'
  is_foreign_key: 1
  is_nullable: 0

=head2 mgmt_ip_address

  data_type: 'varchar'
  is_nullable: 1
  size: 45

=head2 mgmt_ip_netmask

  data_type: 'varchar'
  is_nullable: 1
  size: 45

=head2 mgmt_ip_gateway

  data_type: 'varchar'
  is_nullable: 1
  size: 45

=head2 ilo_ip_address

  data_type: 'varchar'
  is_nullable: 1
  size: 45

=head2 ilo_ip_netmask

  data_type: 'varchar'
  is_nullable: 1
  size: 45

=head2 ilo_ip_gateway

  data_type: 'varchar'
  is_nullable: 1
  size: 45

=head2 ilo_username

  data_type: 'varchar'
  is_nullable: 1
  size: 45

=head2 ilo_password

  data_type: 'varchar'
  is_nullable: 1
  size: 45

=head2 router_host_name

  data_type: 'varchar'
  is_nullable: 1
  size: 256

=head2 router_port_name

  data_type: 'varchar'
  is_nullable: 1
  size: 256

=head2 guid

  data_type: 'varchar'
  is_nullable: 1
  size: 45

=head2 last_updated

  data_type: 'timestamp with time zone'
  default_value: current_timestamp
  is_nullable: 1
  original: {default_value => \"now()"}

=head2 https_port

  data_type: 'bigint'
  is_nullable: 1

=cut

__PACKAGE__->add_columns(
  "id",
  {
    data_type         => "bigint",
    is_auto_increment => 1,
    is_nullable       => 0,
    sequence          => "server_id_seq",
  },
  "host_name",
  { data_type => "varchar", is_nullable => 0, size => 63 },
  "domain_name",
  { data_type => "varchar", is_nullable => 0, size => 63 },
  "tcp_port",
  { data_type => "bigint", is_nullable => 1 },
  "xmpp_id",
  { data_type => "varchar", is_nullable => 1, size => 256 },
  "xmpp_passwd",
  { data_type => "varchar", is_nullable => 1, size => 45 },
  "interface_name",
  { data_type => "varchar", is_nullable => 0, size => 45 },
  "ip_address",
  { data_type => "varchar", is_nullable => 0, size => 45 },
  "ip_netmask",
  { data_type => "varchar", is_nullable => 0, size => 45 },
  "ip_gateway",
  { data_type => "varchar", is_nullable => 0, size => 45 },
  "ip6_address",
  { data_type => "varchar", is_nullable => 1, size => 50 },
  "ip6_gateway",
  { data_type => "varchar", is_nullable => 1, size => 50 },
  "interface_mtu",
  { data_type => "bigint", default_value => 9000, is_nullable => 0 },
  "phys_location",
  { data_type => "bigint", is_foreign_key => 1, is_nullable => 0 },
  "rack",
  { data_type => "varchar", is_nullable => 1, size => 64 },
  "cachegroup",
  {
    data_type      => "bigint",
    default_value  => 0,
    is_foreign_key => 1,
    is_nullable    => 0,
  },
  "type",
  { data_type => "bigint", is_foreign_key => 1, is_nullable => 0 },
  "status",
  { data_type => "bigint", is_foreign_key => 1, is_nullable => 0 },
  "offline_reason",
  { data_type => "varchar", is_nullable => 1, size => 256 },
  "upd_pending",
  { data_type => "boolean", default_value => \"false", is_nullable => 0 },
  "profile",
  { data_type => "bigint", is_foreign_key => 1, is_nullable => 0 },
  "cdn_id",
  { data_type => "bigint", is_foreign_key => 1, is_nullable => 0 },
  "mgmt_ip_address",
  { data_type => "varchar", is_nullable => 1, size => 45 },
  "mgmt_ip_netmask",
  { data_type => "varchar", is_nullable => 1, size => 45 },
  "mgmt_ip_gateway",
  { data_type => "varchar", is_nullable => 1, size => 45 },
  "ilo_ip_address",
  { data_type => "varchar", is_nullable => 1, size => 45 },
  "ilo_ip_netmask",
  { data_type => "varchar", is_nullable => 1, size => 45 },
  "ilo_ip_gateway",
  { data_type => "varchar", is_nullable => 1, size => 45 },
  "ilo_username",
  { data_type => "varchar", is_nullable => 1, size => 45 },
  "ilo_password",
  { data_type => "varchar", is_nullable => 1, size => 45 },
  "router_host_name",
  { data_type => "varchar", is_nullable => 1, size => 256 },
  "router_port_name",
  { data_type => "varchar", is_nullable => 1, size => 256 },
  "guid",
  { data_type => "varchar", is_nullable => 1, size => 45 },
  "last_updated",
  {
    data_type     => "timestamp with time zone",
    default_value => \"current_timestamp",
    is_nullable   => 1,
    original      => { default_value => \"now()" },
  },
  "https_port",
  { data_type => "bigint", is_nullable => 1 },
);

=head1 PRIMARY KEY

=over 4

=item * L</id>

=item * L</cachegroup>

=item * L</type>

=item * L</status>

=item * L</profile>

=back

=cut

__PACKAGE__->set_primary_key("id", "cachegroup", "type", "status", "profile");

=head1 UNIQUE CONSTRAINTS

=head2 C<idx_62588_ip6_profile>

=over 4

=item * L</ip6_address>

=item * L</profile>

=back

=cut

__PACKAGE__->add_unique_constraint("idx_62588_ip6_profile", ["ip6_address", "profile"]);

=head2 C<idx_62588_ip_profile>

=over 4

=item * L</ip_address>

=item * L</profile>

=back

=cut

__PACKAGE__->add_unique_constraint("idx_62588_ip_profile", ["ip_address", "profile"]);

=head2 C<idx_62588_se_id_unique>

=over 4

=item * L</id>

=back

=cut

__PACKAGE__->add_unique_constraint("idx_62588_se_id_unique", ["id"]);

=head1 RELATIONS

=head2 cachegroup

Type: belongs_to

Related object: L<Schema::Result::Cachegroup>

=cut

__PACKAGE__->belongs_to(
  "cachegroup",
  "Schema::Result::Cachegroup",
  { id => "cachegroup" },
  { is_deferrable => 0, on_delete => "CASCADE", on_update => "RESTRICT" },
);

=head2 cdn

Type: belongs_to

Related object: L<Schema::Result::Cdn>

=cut

__PACKAGE__->belongs_to(
  "cdn",
  "Schema::Result::Cdn",
  { id => "cdn_id" },
  { is_deferrable => 0, on_delete => "RESTRICT", on_update => "RESTRICT" },
);

=head2 deliveryservice_servers

Type: has_many

Related object: L<Schema::Result::DeliveryserviceServer>

=cut

__PACKAGE__->has_many(
  "deliveryservice_servers",
  "Schema::Result::DeliveryserviceServer",
  { "foreign.server" => "self.id" },
  { cascade_copy => 0, cascade_delete => 0 },
);

=head2 hwinfos

Type: has_many

Related object: L<Schema::Result::Hwinfo>

=cut

__PACKAGE__->has_many(
  "hwinfos",
  "Schema::Result::Hwinfo",
  { "foreign.serverid" => "self.id" },
  { cascade_copy => 0, cascade_delete => 0 },
);

=head2 phys_location

Type: belongs_to

Related object: L<Schema::Result::PhysLocation>

=cut

__PACKAGE__->belongs_to(
  "phys_location",
  "Schema::Result::PhysLocation",
  { id => "phys_location" },
  { is_deferrable => 0, on_delete => "NO ACTION", on_update => "NO ACTION" },
);

=head2 profile

Type: belongs_to

Related object: L<Schema::Result::Profile>

=cut

__PACKAGE__->belongs_to(
  "profile",
  "Schema::Result::Profile",
  { id => "profile" },
  { is_deferrable => 0, on_delete => "NO ACTION", on_update => "NO ACTION" },
);

=head2 servercheck

Type: might_have

Related object: L<Schema::Result::Servercheck>

=cut

__PACKAGE__->might_have(
  "servercheck",
  "Schema::Result::Servercheck",
  { "foreign.server" => "self.id" },
  { cascade_copy => 0, cascade_delete => 0 },
);

=head2 status

Type: belongs_to

Related object: L<Schema::Result::Status>

=cut

__PACKAGE__->belongs_to(
  "status",
  "Schema::Result::Status",
  { id => "status" },
  { is_deferrable => 0, on_delete => "NO ACTION", on_update => "NO ACTION" },
);

=head2 type

Type: belongs_to

Related object: L<Schema::Result::Type>

=cut

__PACKAGE__->belongs_to(
  "type",
  "Schema::Result::Type",
  { id => "type" },
  { is_deferrable => 0, on_delete => "NO ACTION", on_update => "NO ACTION" },
);


# Created by DBIx::Class::Schema::Loader v0.07045 @ 2016-11-15 08:31:12
# DO NOT MODIFY THIS OR ANYTHING ABOVE! md5sum:6bpdIJoG8OpFcQLdaRdvIA
# These lines were loaded from '/Users/drichard/projects/github.com/traffic_control/traffic_ops/app/lib/Schema/Result/Server.pm' found in @INC.
# They are now part of the custom portion of this file
# for you to hand-edit.  If you do not either delete
# this section or remove that file from @INC, this section
# will be repeated redundantly when you re-create this
# file again via Loader!  See skip_load_external to disable
# this feature.

use utf8;
package Schema::Result::Server;

# Created by DBIx::Class::Schema::Loader
# DO NOT MODIFY THE FIRST PART OF THIS FILE

=head1 NAME

Schema::Result::Server

=cut

use strict;
use warnings;

use base 'DBIx::Class::Core';

=head1 TABLE: C<server>

=cut

__PACKAGE__->table("server");

=head1 ACCESSORS

=head2 id

  data_type: 'integer'
  is_auto_increment: 1
  is_nullable: 0

=head2 host_name

  data_type: 'varchar'
  is_nullable: 0
  size: 63

=head2 domain_name

  data_type: 'varchar'
  is_nullable: 0
  size: 63

=head2 tcp_port

  data_type: 'integer'
  extra: {unsigned => 1}
  is_nullable: 1

=head2 xmpp_id

  data_type: 'varchar'
  is_nullable: 1
  size: 256

=head2 xmpp_passwd

  data_type: 'varchar'
  is_nullable: 1
  size: 45

=head2 interface_name

  data_type: 'varchar'
  is_nullable: 0
  size: 45

=head2 ip_address

  data_type: 'varchar'
  is_nullable: 0
  size: 45

=head2 ip_netmask

  data_type: 'varchar'
  is_nullable: 0
  size: 45

=head2 ip_gateway

  data_type: 'varchar'
  is_nullable: 0
  size: 45

=head2 ip6_address

  data_type: 'varchar'
  is_nullable: 1
  size: 50

=head2 ip6_gateway

  data_type: 'varchar'
  is_nullable: 1
  size: 50

=head2 interface_mtu

  data_type: 'integer'
  default_value: 9000
  is_nullable: 0

=head2 phys_location

  data_type: 'integer'
  is_foreign_key: 1
  is_nullable: 0

=head2 rack

  data_type: 'varchar'
  is_nullable: 1
  size: 64

=head2 cachegroup

  data_type: 'integer'
  default_value: 0
  is_foreign_key: 1
  is_nullable: 0

=head2 type

  data_type: 'integer'
  is_foreign_key: 1
  is_nullable: 0

=head2 status

  data_type: 'integer'
  is_foreign_key: 1
  is_nullable: 0

=head2 upd_pending

  data_type: 'tinyint'
  default_value: 0
  is_nullable: 0

=head2 profile

  data_type: 'integer'
  is_foreign_key: 1
  is_nullable: 0

=head2 mgmt_ip_address

  data_type: 'varchar'
  is_nullable: 1
  size: 45

=head2 mgmt_ip_netmask

  data_type: 'varchar'
  is_nullable: 1
  size: 45

=head2 mgmt_ip_gateway

  data_type: 'varchar'
  is_nullable: 1
  size: 45

=head2 ilo_ip_address

  data_type: 'varchar'
  is_nullable: 1
  size: 45

=head2 ilo_ip_netmask

  data_type: 'varchar'
  is_nullable: 1
  size: 45

=head2 ilo_ip_gateway

  data_type: 'varchar'
  is_nullable: 1
  size: 45

=head2 ilo_username

  data_type: 'varchar'
  is_nullable: 1
  size: 45

=head2 ilo_password

  data_type: 'varchar'
  is_nullable: 1
  size: 45

=head2 router_host_name

  data_type: 'varchar'
  is_nullable: 1
  size: 256

=head2 router_port_name

  data_type: 'varchar'
  is_nullable: 1
  size: 256

=head2 last_updated

  data_type: 'timestamp'
  datetime_undef_if_invalid: 1
  default_value: current_timestamp
  is_nullable: 1

=cut

__PACKAGE__->add_columns(
  "id",
  { data_type => "integer", is_auto_increment => 1, is_nullable => 0 },
  "host_name",
  { data_type => "varchar", is_nullable => 0, size => 63 },
  "domain_name",
  { data_type => "varchar", is_nullable => 0, size => 63 },
  "tcp_port",
  { data_type => "integer", extra => { unsigned => 1 }, is_nullable => 1 },
  "xmpp_id",
  { data_type => "varchar", is_nullable => 1, size => 256 },
  "xmpp_passwd",
  { data_type => "varchar", is_nullable => 1, size => 45 },
  "interface_name",
  { data_type => "varchar", is_nullable => 0, size => 45 },
  "ip_address",
  { data_type => "varchar", is_nullable => 0, size => 45 },
  "ip_netmask",
  { data_type => "varchar", is_nullable => 0, size => 45 },
  "ip_gateway",
  { data_type => "varchar", is_nullable => 0, size => 45 },
  "ip6_address",
  { data_type => "varchar", is_nullable => 1, size => 50 },
  "ip6_gateway",
  { data_type => "varchar", is_nullable => 1, size => 50 },
  "interface_mtu",
  { data_type => "integer", default_value => 9000, is_nullable => 0 },
  "phys_location",
  { data_type => "integer", is_foreign_key => 1, is_nullable => 0 },
  "rack",
  { data_type => "varchar", is_nullable => 1, size => 64 },
  "cachegroup",
  {
    data_type      => "integer",
    default_value  => 0,
    is_foreign_key => 1,
    is_nullable    => 0,
  },
  "type",
  { data_type => "integer", is_foreign_key => 1, is_nullable => 0 },
  "status",
  { data_type => "integer", is_foreign_key => 1, is_nullable => 0 },
  "upd_pending",
  { data_type => "tinyint", default_value => 0, is_nullable => 0 },
  "profile",
  { data_type => "integer", is_foreign_key => 1, is_nullable => 0 },
  "mgmt_ip_address",
  { data_type => "varchar", is_nullable => 1, size => 45 },
  "mgmt_ip_netmask",
  { data_type => "varchar", is_nullable => 1, size => 45 },
  "mgmt_ip_gateway",
  { data_type => "varchar", is_nullable => 1, size => 45 },
  "ilo_ip_address",
  { data_type => "varchar", is_nullable => 1, size => 45 },
  "ilo_ip_netmask",
  { data_type => "varchar", is_nullable => 1, size => 45 },
  "ilo_ip_gateway",
  { data_type => "varchar", is_nullable => 1, size => 45 },
  "ilo_username",
  { data_type => "varchar", is_nullable => 1, size => 45 },
  "ilo_password",
  { data_type => "varchar", is_nullable => 1, size => 45 },
  "router_host_name",
  { data_type => "varchar", is_nullable => 1, size => 256 },
  "router_port_name",
  { data_type => "varchar", is_nullable => 1, size => 256 },
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

=item * L</cachegroup>

=item * L</type>

=item * L</status>

=item * L</profile>

=back

=cut

__PACKAGE__->set_primary_key("id", "cachegroup", "type", "status", "profile");

=head2 C<se_id_UNIQUE>

=over 4

=item * L</id>

=back

=cut

__PACKAGE__->add_unique_constraint("se_id_UNIQUE", ["id"]);

=head1 RELATIONS

=head2 cachegroup

Type: belongs_to

Related object: L<Schema::Result::Cachegroup>

=cut

__PACKAGE__->belongs_to(
  "cachegroup",
  "Schema::Result::Cachegroup",
  { id => "cachegroup" },
  { is_deferrable => 1, on_delete => "CASCADE", on_update => "RESTRICT" },
);

=head2 deliveryservice_servers

Type: has_many

Related object: L<Schema::Result::DeliveryserviceServer>

=cut

__PACKAGE__->has_many(
  "deliveryservice_servers",
  "Schema::Result::DeliveryserviceServer",
  { "foreign.server" => "self.id" },
  { cascade_copy => 0, cascade_delete => 0 },
);

=head2 hwinfos

Type: has_many

Related object: L<Schema::Result::Hwinfo>

=cut

__PACKAGE__->has_many(
  "hwinfos",
  "Schema::Result::Hwinfo",
  { "foreign.serverid" => "self.id" },
  { cascade_copy => 0, cascade_delete => 0 },
);

=head2 phys_location

Type: belongs_to

Related object: L<Schema::Result::PhysLocation>

=cut

__PACKAGE__->belongs_to(
  "phys_location",
  "Schema::Result::PhysLocation",
  { id => "phys_location" },
  { is_deferrable => 1, on_delete => "NO ACTION", on_update => "NO ACTION" },
);

=head2 profile

Type: belongs_to

Related object: L<Schema::Result::Profile>

=cut

__PACKAGE__->belongs_to(
  "profile",
  "Schema::Result::Profile",
  { id => "profile" },
  { is_deferrable => 1, on_delete => "NO ACTION", on_update => "NO ACTION" },
);

=head2 servercheck

Type: might_have

Related object: L<Schema::Result::Servercheck>

=cut

__PACKAGE__->might_have(
  "servercheck",
  "Schema::Result::Servercheck",
  { "foreign.server" => "self.id" },
  { cascade_copy => 0, cascade_delete => 0 },
);

=head2 status

Type: belongs_to

Related object: L<Schema::Result::Status>

=cut

__PACKAGE__->belongs_to(
  "status",
  "Schema::Result::Status",
  { id => "status" },
  { is_deferrable => 1, on_delete => "NO ACTION", on_update => "NO ACTION" },
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


# Created by DBIx::Class::Schema::Loader v0.07043 @ 2015-05-21 13:27:11
# DO NOT MODIFY THIS OR ANYTHING ABOVE! md5sum:mecAHLqlmqBRoRpMOHaiOQ


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
# End of lines loaded from '/Users/drichard/projects/github.com/traffic_control/traffic_ops/app/lib/Schema/Result/Server.pm'


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
