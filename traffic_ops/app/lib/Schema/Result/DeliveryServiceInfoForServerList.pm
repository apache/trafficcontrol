use utf8;

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
#

package Schema::Result::DeliveryServiceInfoForServerList;

# this view returns the regexp set for a delivery services, ordered by type, set_number.
# to use, do
#
# $rs = $self->db->resultset('DeliveryServiceInfoForServerList')->search({}, { bind => [ \$server ]});
#
# where $id is the deliveryservice id.

use strict;
use warnings;

use base 'DBIx::Class::Core';

__PACKAGE__->table_class('DBIx::Class::ResultSource::View');

__PACKAGE__->table("DeliveryServiceInfoForServerList:");

__PACKAGE__->result_source_instance->is_virtual(1);

__PACKAGE__->result_source_instance->view_definition( "
SELECT DISTINCT
    deliveryservice.xml_id AS xml_id,
    deliveryservice.id AS ds_id,
    deliveryservice.dscp AS dscp,
    deliveryservice.signed AS signed,
    deliveryservice.qstring_ignore AS qstring_ignore,
    deliveryservice.org_server_fqdn as org_server_fqdn,
    deliveryservice.multi_site_origin as multi_site_origin,
    deliveryservice.range_request_handling as range_request_handling,
    deliveryservice.origin_shield as origin_shield,
    regex.pattern AS pattern,
    retype.name AS re_type,
    dstype.name AS ds_type,
    parameter.value AS domain_name,
    deliveryservice_regex.set_number AS set_number,
    deliveryservice.edge_header_rewrite as edge_header_rewrite,
    deliveryservice.regex_remap as regex_remap,
    deliveryservice.cacheurl as cacheurl,
    deliveryservice.remap_text as remap_text,
    mid_header_rewrite as mid_header_rewrite,
    deliveryservice.protocol as protocol
FROM
    deliveryservice
        JOIN deliveryservice_regex ON deliveryservice_regex.deliveryservice = deliveryservice.id
        JOIN regex ON deliveryservice_regex.regex = regex.id
        JOIN type as retype ON regex.type = retype.id
        JOIN type as dstype ON deliveryservice.type = dstype.id
        JOIN profile_parameter ON deliveryservice.profile = profile_parameter.profile
        JOIN parameter ON parameter.id = profile_parameter.parameter
        JOIN deliveryservice_server ON deliveryservice_server.deliveryservice = deliveryservice.id
        JOIN server ON deliveryservice_server.server = server.id
WHERE parameter.name = 'domain_name' AND server.id IN (?)
ORDER BY ds_id, re_type , deliveryservice_regex.set_number
"
);

__PACKAGE__->add_columns(
	"xml_id",                 { data_type => "varchar", is_nullable => 0, size => 45 },
	"org_server_fqdn",        { data_type => "varchar", is_nullable => 0, size => 45 },
	"multi_site_origin",      { data_type => "integer", is_nullable => 0 },
	"ds_id",                  { data_type => "integer", is_nullable => 0 },
	"dscp",                   { data_type => "integer", is_nullable => 0 },
	"signed",                 { data_type => "integer", is_nullable => 0 },
	"qstring_ignore",         { data_type => "integer", is_nullable => 0 },
	"pattern",                { data_type => "varchar", is_nullable => 0, size => 45 },
	"re_type",                { data_type => "varchar", is_nullable => 0, size => 45 },
	"ds_type",                { data_type => "varchar", is_nullable => 0, size => 45 },
	"set_number",             { data_type => "integer", is_nullable => 0 },
	"domain_name",            { data_type => "varchar", is_nullable => 0, size => 45 },
	"edge_header_rewrite",    { data_type => "varchar", is_nullable => 0, size => 1024 },
	"mid_header_rewrite",     { data_type => "varchar", is_nullable => 0, size => 1024 },
	"regex_remap",            { data_type => "varchar", is_nullable => 0, size => 1024 },
	"cacheurl",               { data_type => "varchar", is_nullable => 0, size => 1024 },
	"remap_text",             { data_type => "varchar", is_nullable => 0, size => 2048 },
	"protocol",               { data_type => "tinyint", is_nullable => 0, size => 4 },
	"range_request_handling", { data_type => "tinyint", is_nullable => 0, size => 4 },
	"origin_shield",          { data_type => "varchar", is_nullable => 0, size => 1024 },
);

1;
