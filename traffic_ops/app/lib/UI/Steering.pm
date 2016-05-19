package UI::Steering;
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

# JvD Note: you always want to put Utils as the first use. Sh*t don't work if it's after the Mojo lines.
use UI::Utils;
use Mojo::Base 'Mojolicious::Controller';
use Data::Dumper;
use UI::DeliveryService;
use API::Cdn;
use Scalar::Util qw(looks_like_number);
use JSON;
use POSIX qw(strftime);
use Date::Parse;

sub index {
	my $self = shift;
	my $ds_id = $self->param('id');

	&navbarpage($self);
	# select * from steering_target where deliveryservice = ds_id;
	my $steering = {ds_id => $ds_id, ds_name => $self->get_ds_name($ds_id)};
	my $st_rs = $self->db->resultset('SteeringTarget')->search( {deliveryservice => $ds_id} );
	my %target_filters;
	if ($st_rs > 0) {
		my %steering_targets;
		while (my $row = $st_rs->next) {
			$steering_targets{$row->target} = $row->weight;
		}
		my @keys = sort keys %steering_targets;
		 $steering->{'target_id_1'} = $keys[0];
		 $steering->{'target_id_2'} = $keys[1];
		 $steering->{'target_name_1'} = $self->get_ds_name($keys[0]);
		 $steering->{'target_name_2'} = $self->get_ds_name($keys[1]);
		 $steering->{'target_weight_1'} = $steering_targets{$keys[0]};
		 $steering->{'target_weight_2'} = $steering_targets{$keys[1]};
		 foreach my $targetid (@keys) {
		 	my $regex_rs = $self->db->resultset('DeliveryserviceRegex')->search( 
		 		{
		 			deliveryservice => $targetid, 
		 			'type.name' => 'STEERING_REGEXP' 
		 		}, 
		 		{
		 			join => {regex => {type => undef} } } 
		 	);
			my @regexes;
			if ($regex_rs > 0) {
				while (my $row = $regex_rs->next) {
					push @regexes, $row->regex->pattern;
				}
			$target_filters{$targetid} = \@regexes;
			}
		 }
	}	

	$self->stash( 
		steering => $steering,
		ds_data => $self->get_deliveryservices(),
		target_filters => \%target_filters,	
		fbox_layout => 1
	);
}

sub get_ds_name {
	my $self = shift;
	my $ds_id = shift;
	return $self->db->resultset('Deliveryservice')->search( { id => $ds_id } )->get_column('xml_id')->single();
}

sub get_deliveryservices {
	my $self = shift;
	my %ds_data;
	my $rs = $self->db->resultset("Deliveryservice")->search( undef, { order_by => "xml_id" } );
	while ( my $row = $rs->next ) {
		$ds_data{$row->id} = $row->xml_id;
	}
	return \%ds_data
}

sub update {
	my $self = shift;
	my $ds_id = $self->param('id');
	my $tid1 = $self->param('steering.target_id_1');
	my $tid2 = $self->param('steering.target_id_2');
	my $tw1 = $self->param('steering.target_weight_1');
	my $tw2 = $self->param('steering.target_weight_2');
	if ( $self->is_valid() ) {
		my $targets;
		$targets->{$tid1} = $tw1;
		$targets->{$tid2} = $tw2;

		#delete current entries
		my $delete = $self->db->resultset('SteeringTarget')->search( { deliveryservice => $ds_id } );
		if (defined($delete)) {
			$delete->delete();
		}

		#add new entries
		foreach my $target (keys $targets) {
			my $insert = $self->db->resultset('SteeringTarget')->create( { 
				deliveryservice => $ds_id,  
				target => $target,
				weight => $targets->{$target},
				}  );
			$insert->insert();
		}
			
		$self->flash( message => "Successfully saved steering assignments for " . $self->get_ds_name($ds_id) . "!" );

		$self->redirect_to("/ds/$ds_id/steering");
	} else {
		&stash_role($self);
		$self->stash( 
			steering => {
				ds_id => $ds_id, 
				ds_name => $self->get_ds_name($ds_id),
				target_id_1 => $tid1,
				target_id_2 => $tid2,
				target_name_1 => $self->get_ds_name($tid1),
				target_name_2 => $self->get_ds_name($tid2),
				target_weight_1 => $tw1,
				target_weight_2 => $tw2
				}, 
			ds_data => $self->get_deliveryservices(),
			fbox_layout => 1 
		);
		$self->render("steering/index");
	}
}

sub is_valid {
	my $self           = shift;
	my $ds_id = $self->param('id');

	# $self->field('steering.target_name_1')->is_equal( "", "This is an error!" );
	##TODO:  ADD VALIDATION!!!

	return $self->valid;
	# return 0;
}

1;
