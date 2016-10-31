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

package Schema::Result::ServersParentCachegroupList;

# this view returns the parent cachegroups for a list of servers.
# to use, do
#
# $rs = $self->db->resultset('ServersParentCachegroupList')->search({});
#

use strict;
use warnings;

use base 'DBIx::Class::Core';

__PACKAGE__->table_class('DBIx::Class::ResultSource::View');

__PACKAGE__->table("ServersParentCachegroupList");

__PACKAGE__->result_source_instance->is_virtual(1);

__PACKAGE__->result_source_instance->view_definition( "
SELECT
	server.id as server_id,
    cachegroup.parent_cachegroup_id AS parent_cachegroup_id
FROM
    cachegroup
        JOIN server ON cachegroup.id = server.cachegroup
WHERE cachegroup.parent_cachegroup_id IS NOT NULL
"
);

__PACKAGE__->add_columns(
	"server_id",      { data_type => "integer", is_nullable => 0 },
	"parent_cachegroup_id",      { data_type => "integer", is_nullable => 0 },
);

1;