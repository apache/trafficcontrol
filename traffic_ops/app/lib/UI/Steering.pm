package UI::Steering;
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

# JvD Note: you always want to put Utils as the first use. Sh*t don't work if it's after the Mojo lines.
use UI::Utils;
use Mojo::Base 'Mojolicious::Controller';
use Mojo::Parameters;
use Data::Dumper;
use UI::DeliveryService;
use API::Cdn;
use Scalar::Util qw(looks_like_number);
use JSON;
use POSIX qw(strftime);
use Date::Parse;

sub index {
	my $self  = shift;
	my $ds_id = $self->param('id');

	my ($type_names, $type_ids) = $self->get_types();

	my @steering = $self->get_target_data($ds_id, $type_ids);

	my @targets;
	foreach my $i ( keys @steering ) {
		push ( @targets, $steering[$i]->{'target_id'} );
	}

	my %ds_data = $self->get_deliveryservices($ds_id, \@targets);

	&navbarpage($self);

	$self->stash(
		ds_id          => $ds_id,
		ds_name        => $self->get_ds_name($ds_id),
		steering       => \@steering,
		ds_data        => \%ds_data,
		types          => $type_names,
		fbox_layout    => 1
	);
}

sub get_types {
	my $self = shift;
	my $t_rs = $self->db->resultset('Type')->search( { use_in_table => 'steering_target'} );
	my $type_names;
	my $type_ids;

	if ( $t_rs > 0 ) {
		while ( my $row = $t_rs->next ) {
			$type_names->{$row->id} = $row->name;
			$type_ids->{$row->name} = $row->id;
		}
	}

	return ($type_names, $type_ids);
}

sub get_target_data {
	my $self = shift;
	my $ds_id = shift;
	my $type_ids = shift;
	my $steering_obj;
	my @steering;
	my @positive_order_steering;

	my $neg_order_rs = $self->db->resultset('SteeringTarget')->search( { deliveryservice => $ds_id, type => $type_ids->{'STEERING_ORDER'}, value => { '<', 0 } }, { order_by => 'value ASC' } );

	if ( $neg_order_rs > 0 ) {
		my $i = 0;
		while ( my $row = $neg_order_rs->next ) {
			my $t = $steering_obj->{"target_$i"};
			$t->{'target_id'} = $row->target;
			$t->{'target_name'}   = $self->get_ds_name( $row->target );
			$t->{'target_value'}   = $row->value;
			if (!defined($t->{'target_value'})) { $t->{'target_value'} = 0; }
			$t->{'target_type'}   = $row->type->id;
			push ( @steering, $t );
			$i++;
		}	
	}

	my $weight_rs = $self->db->resultset('SteeringTarget')->search( { deliveryservice => $ds_id, type => $type_ids->{'STEERING_WEIGHT'} }, { order_by => 'value DESC' } );

	if ( $weight_rs > 0 ) {
		my $i = 0;
		while ( my $row = $weight_rs->next ) {
			my $t = $steering_obj->{"target_$i"};
			$t->{'target_id'} = $row->target;
			$t->{'target_name'}   = $self->get_ds_name( $row->target );
			$t->{'target_value'}   = $row->value;
			if (!defined($t->{'target_value'})) { $t->{'target_value'} = 0; }
			$t->{'target_type'}   = $row->type->id;
			push ( @steering, $t );
			$i++;
		}
	}

	my $pos_order_rs = $self->db->resultset('SteeringTarget')->search( { deliveryservice => $ds_id, type => $type_ids->{'STEERING_ORDER'}, value => { '>=', 0 } }, { order_by => 'value ASC' } );

	if ( $pos_order_rs > 0 ) {
		my $i = 0;
		while ( my $row = $pos_order_rs->next ) {
			my $t = $steering_obj->{"target_$i"};
			$t->{'target_id'} = $row->target;
			$t->{'target_name'}   = $self->get_ds_name( $row->target );
			$t->{'target_value'}   = $row->value;
			if (!defined($t->{'target_value'})) { $t->{'target_value'} = 0; }
			$t->{'target_type'}   = $row->type->id;
			push ( @steering, $t );
			$i++;
		}	
	}
	return @steering;
}

sub get_ds_name {
	my $self  = shift;
	my $ds_id = shift;
	return $self->db->resultset('Deliveryservice')->search( { id => $ds_id } )->get_column('xml_id')->single();
}

sub get_cdn {
	my $self = shift;
	my $ds_id = shift;
	return $self->db->resultset('Deliveryservice')->search( { id => $ds_id } )->get_column('cdn_id')->single();

}

sub get_deliveryservices {
	my $self = shift;
	my $ds_id = shift;
	my @targets = @{$_[0]};

	my $cdn_id = $self->get_cdn($ds_id);
	my %ds_data;
	#search for only the delivery services that match the CDN ID of the supplied delivery service.
	my $rs = $self->db->resultset('Deliveryservice')->search({ cdn_id => $cdn_id } , { prefetch => [ 'type' ] });
	while ( my $row = $rs->next ) {
		my $ds = $row->id;
		if ( $row->type->name =~ m/^HTTP/ ) {
			if (!grep( /$ds/, @targets )) {
				$ds_data{ $row->id } = $row->xml_id;
			}
		}
	}

	return %ds_data;
}

sub update {
	my $self = shift;
	my $ds_id = $self->param('id');
	my @target_id = $self->param('st.target_id');
	my @target_value = $self->param('st.target_value');
	my @target_type = $self->param('st.target_type');
	my @targets;
	my $steering_obj;
	foreach my $i (0 .. $#target_id) {
		#look for and remove the blank entries - this filters out the deleted entries and the unused new target entry.
		if ( $target_id[$i] eq '' ) {
			next;
		}
		if ( $target_value[$i] eq "" ) { $target_value[$i] = 0 };
		$steering_obj->{"target_$i"}->{'target_id'} = $target_id[$i];
		$steering_obj->{"target_$i"}->{'target_value'} = $target_value[$i];
		$steering_obj->{"target_$i"}->{'target_type'} = $target_type[$i];
		push ( @targets, $steering_obj->{"target_$i"} );
	}
	if ( $self->is_valid(\@targets) ) {
		#delete current entries
		my $delete = $self->db->resultset('SteeringTarget')
			->search( { deliveryservice => $ds_id } );
		if ( defined($delete) ) {
			$delete->delete();
		}
		
		#add new entries
		foreach my $i ( keys @targets ) {
			my $insert = $self->db->resultset('SteeringTarget')->create(
				{   deliveryservice => $ds_id,
					target          => $targets[$i]->{'target_id'},
					value           => $targets[$i]->{'target_value'},
					type            => $targets[$i]->{'target_type'}
				}
			);
			$insert->insert();
		}
		
		$self->flash(
			      message => "Successfully saved steering assignments for "
				. $self->get_ds_name($ds_id)
				. "!" );
	}
	else {
		
		
		my ($type_names, $type_ids) = $self->get_types();
	
		my @steering = $self->get_target_data($ds_id, $type_ids);

		my @targets;
		foreach my $i ( keys @steering ) {
			push ( @targets, $steering[$i]->{'target_id'} );
		}

		my %ds_data = $self->get_deliveryservices($ds_id, \@targets);
		
		&stash_role($self);
		$self->stash(
			ds_id          => $ds_id,
			ds_name        => $self->get_ds_name($ds_id),
			steering       => \@steering,
			ds_data        => \%ds_data,
			types          => $type_names,
			fbox_layout    => 1
		);
		$self->render("steering/index");
	}

	$self->redirect_to("/ds/$ds_id/steering");
}


sub is_valid {
	my $self  = shift;
	my @targets = @{$_[0]};
	my %tracker;

	foreach my $i ( keys @targets ) {
		my $t = $targets[$i];
		my $t_name = $self->db->resultset('Type')->search( { id => "$t->{'target_type'}" } )->get_column('name')->single();
		if ( $t_name eq "STEERING_ORDER" && $t->{'target_value'} ne int($t->{'target_value'})) {
			$self->flash(message => "STEERING_ORDER values must be integers." );
			return;
		}
		elsif ( $t_name eq "STEERING_WEIGHT" && ( $t->{'target_value'} ne int($t->{'target_value'}) || ($t->{'target_value'} < 0 ) ) )  {
			$self->flash(message => "STEERING_WEIGHT values must be integers greater than 0." );
			return;
		}
		if (exists($t->{'target_id'})) {
			$tracker{$t->{'target_id'}}++;
			if ( $tracker{$t->{'target_id'}} > 1 ) {
				$self->flash(message => "Target delivery services must be unique." );
				return;
			}
		}
	}

	return $self->valid;
}

1;
