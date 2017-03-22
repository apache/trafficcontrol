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

package Schema::Result::ServersForCdn;

use strict;
use warnings;

use base 'DBIx::Class::Core';

__PACKAGE__->table_class('DBIx::Class::ResultSource::View');

__PACKAGE__->table("ServersForCdn");

__PACKAGE__->result_source_instance->is_virtual(1);

__PACKAGE__->result_source_instance->view_definition(
	"SELECT
		 me.host_name,
		 me.domain_name,
		 me.tcp_port,
		 me.interface_name,
		 me.ip_address,
		 me.ip6_address,
		 me.xmpp_id,
		 type.name as type,
		 status.name as status,
		 cachegroup.name as cachegroup,
		 profile.name as profile
	FROM server me
	 JOIN type type ON type.id = me.type
	 JOIN status status ON status.id = me.status
	 JOIN cachegroup cachegroup ON cachegroup.id = me.cachegroup
	 JOIN profile profile ON profile.id = me.profile
	 JOIN cdn cdn ON cdn.id = me.cdn_id
	WHERE cdn.name = ?"
);

__PACKAGE__->add_columns(
	"host_name",
	{ data_type => "varchar", is_nullable => 0 },
	"domain_name",
	{ data_type => "varchar", is_nullable => 0 },
	"tcp_port",
	{ data_type => "integer", extra => { unsigned => 1 }, is_nullable => 1 },
	"interface_name",
	{ data_type => "varchar", is_nullable => 0 },
	"ip_address",
	{ data_type => "varchar", is_nullable => 0 },
	"ip6_address",
	{ data_type => "varchar", is_nullable => 1 },
	"xmpp_id",
	{ data_type => "varchar", is_nullable => 1 },
	"type",
	{ data_type => "varchar", is_nullable => 0 },
	"status",
	{ data_type => "varchar", is_nullable => 0 },
	"cachegroup",
	{ data_type => "varchar", is_nullable => 0 },
	"profile",
	{ data_type => "varchar", is_nullable => 0 },
);

1;
