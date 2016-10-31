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

package Schema::Result::RegexByDeliveryServiceList;

# this view returns the regexp set for a delivery services, ordered by type, set_number.
# to use, do 
#
# $rs = $self->db->resultset('RegexByDeliveryServiceList')->search();
#
# where $id is the deliveryservice id.

use strict;
use warnings;

use base 'DBIx::Class::Core';

__PACKAGE__->table_class('DBIx::Class::ResultSource::View');

__PACKAGE__->table("RegexByDeliveryServiceList");

__PACKAGE__->result_source_instance->is_virtual(1);

__PACKAGE__->result_source_instance->view_definition( "
SELECT 
    deliveryservice.xml_id AS shortname,
	deliveryservice.id AS ds_id,
    regex.pattern AS pattern,
    type.name AS type,
    deliveryservice_regex.set_number AS set_number
FROM
    deliveryservice
        JOIN
    deliveryservice_regex ON deliveryservice_regex.deliveryservice = deliveryservice.id
        JOIN
    regex ON deliveryservice_regex.regex = regex.id
        JOIN
    type ON regex.type = type.id
ORDER BY type , deliveryservice_regex.set_number
"
);

__PACKAGE__->add_columns(
	"shortname",  { data_type => "varchar", is_nullable => 0, size => 45 }, 
	"ds_id",      { data_type => "integer", is_nullable => 0 },
	"pattern",    { data_type => "varchar", is_nullable => 0, size => 45 },
	"type",       { data_type => "varchar", is_nullable => 0, size => 45 }, 
	"set_number", { data_type => "integer", is_nullable => 0 },
);

1;