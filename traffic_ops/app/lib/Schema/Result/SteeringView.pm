package Schema::Result::SteeringView;

# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
# 
#   http://www.apache.org/licenses/LICENSE-2.0
# 
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.


use strict;
use warnings;

use base 'DBIx::Class::Core';

__PACKAGE__->table_class('DBIx::Class::ResultSource::View');

__PACKAGE__->table("SteeringView");

__PACKAGE__->result_source_instance->is_virtual(1);

__PACKAGE__->result_source_instance->view_definition(
    "select s.xml_id as steering_xml_id, s.id as steering_id, t.xml_id as target_xml_id, t.id as target_id, value, tp.name as type from steering_target
    join deliveryservice s on s.id = steering_target.deliveryservice
    join deliveryservice t on t.id = steering_target.target
    join type tp on tp.id = steering_target.type"
);

__PACKAGE__->add_columns(
    "steering_xml_id",
    { data_type => "varchar", is_nullable => 0, size => 50 },
    "steering_id",
    { data_type => "integer", is_nullable => 0, size => 11 },
    "target_xml_id",
    { data_type => "varchar", is_nullable => 0, size => 50 },
    "target_id",
    { data_type => "integer", is_nullable => 0, size => 11 },
    "value",
    { data_type => "integer", extra => { unsigned => 1 }, is_nullable => 0 },
    "type",
    { data_type => "integer", is_nullable => 0 }
);

1;
