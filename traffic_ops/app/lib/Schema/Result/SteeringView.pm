package SteeringView;
use strict;
use warnings FATAL => 'all';

package Schema::Result::SteeringView;

use strict;
use warnings;

use base 'DBIx::Class::Core';

__PACKAGE__->table_class('DBIx::Class::ResultSource::View');

__PACKAGE__->table("SteeringView");

__PACKAGE__->result_source_instance->is_virtual(1);

__PACKAGE__->result_source_instance->view_definition(
    "select s.xml_id as steering_xml_id, s.id as steering_id, t.xml_id as target_xml_id, t.id as target_id, weight from steering_target
    join deliveryservice s on s.id = steering_target.deliveryservice
    join deliveryservice t on t.id = steering_target.target"
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
    "weight",
    { data_type => "integer", extra => { unsigned => 1 }, is_nullable => 0 },
);

1;