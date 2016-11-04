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

package Schema::Result::RascalHostsByCdnAll;

use strict;
use warnings;

use base 'DBIx::Class::Core';

__PACKAGE__->table_class('DBIx::Class::ResultSource::View');

__PACKAGE__->table("RascalHostsByCdnAll");

__PACKAGE__->result_source_instance->is_virtual(1);

__PACKAGE__->result_source_instance->view_definition(
	"SELECT s.host_name,s.domain_name,s.tcp_port,c.name as cdn_name,st.name as status
  FROM server s
  JOIN status st ON st.id = s.status
  JOIN cdn c ON c.id = s.cdn_id 
  WHERE  s.type = (SELECT type.id FROM type WHERE name='RASCAL') 
  GROUP BY c.name, s.host_name, s.domain_name, s.tcp_port, st.name"
);

__PACKAGE__->add_columns(
	"host_name",
	{ data_type => "varchar", is_nullable => 0, size => 45 },
	"domain_name",
	{ data_type => "varchar", is_nullable => 0, size => 45 },
	"tcp_port",
	{ data_type => "integer", extra => { unsigned => 1 }, is_nullable => 1 },
	"cdn_name",
	{ data_type => "varchar", is_nullable => 0, size => 45 },
	"status",
	{ data_type => "varchar", is_nullable => 0, size => 45 },
);

1;
