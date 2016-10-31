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

package Schema::Result::ServerTypes;

use strict;
use warnings;

use base 'DBIx::Class::Core';

__PACKAGE__->table_class('DBIx::Class::ResultSource::View');

__PACKAGE__->table("ServerTypes");

__PACKAGE__->result_source_instance->is_virtual(1);

__PACKAGE__->result_source_instance->view_definition("SELECT * FROM TYPE WHERE ID IN( SELECT TYPE FROM SERVER )");

__PACKAGE__->add_columns(
	"id",
	{ data_type => "integer", is_auto_increment => 1, is_nullable => 0 },
	"name",
	{ data_type => "varchar", is_nullable => 0, size => 45 },
	"description",
	{ data_type => "varchar", is_nullable => 1, size => 256 },
	"use_in_table",
	{ data_type => "varchar", is_nullable => 1, size => 45 },
	"last_updated", {
		data_type                 => "timestamp",
		datetime_undef_if_invalid => 1,
		default_value             => \"current_timestamp",
		is_nullable               => 1,
	},
);
1;
