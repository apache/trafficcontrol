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

	&navbarpage($self);

	# select * from steering_target where deliveryservice = ds_id;
	my $steering = { ds_id => $ds_id, ds_name => $self->get_ds_name($ds_id) };
	my $st_rs = $self->db->resultset('SteeringTarget')->search( { deliveryservice => $ds_id } );
	if ( $st_rs > 0 ) {
		my %steering_targets;
		while ( my $row = $st_rs->next ) {
			$steering_targets{ $row->target } = $row->weight;
		}
		my @keys = sort keys %steering_targets;
		$steering->{'target_id_1'}     = $keys[0];
		$steering->{'target_id_2'}     = $keys[1];
		$steering->{'target_name_1'}   = $self->get_ds_name( $keys[0] );
		$steering->{'target_name_2'}   = $self->get_ds_name( $keys[1] );
		$steering->{'target_id_1_weight'}   = $self->get_target_weight( $ds_id, $keys[0] );
		$steering->{'target_id_2_weight'}   = $self->get_target_weight( $ds_id, $keys[1] );
	}
	if (!defined($steering->{'target_id_1_weight'})) { $steering->{'target_id_1_weight'} = 0; }
	if (!defined($steering->{'target_id_2_weight'})) { $steering->{'target_id_2_weight'} = 0; }
	
	$self->stash(
		steering       => $steering,
		ds_data        => $self->get_deliveryservices(),
		fbox_layout    => 1
	);
}

sub get_target_weight{
	my $self = shift;
	my $ds_id = shift;
	my $target_id = shift;
	my $weight = $self->db->resultset('SteeringTarget')->search( { -and => [target => $target_id, deliveryservice => $ds_id] } )->get_column('weight')->single();
	return $weight;
}

sub get_ds_name {
	my $self  = shift;
	my $ds_id = shift;
	return $self->db->resultset('Deliveryservice')->search( { id => $ds_id } )->get_column('xml_id')->single();
}

sub get_deliveryservices {
	my $self = shift;
	my %ds_data;
	my $rs = $self->db->resultset('Deliveryservice')->search(undef, { prefetch => [ 'type' ] });
	while ( my $row = $rs->next ) {
		if ( $row->type->name =~ m/^HTTP/ ) {
			$ds_data{ $row->id } = $row->xml_id;
		}
	}

	return \%ds_data;
}

sub update {
	my $self  = shift;
	my $ds_id = $self->param('id');
	my $tid1  = $self->param('steering.target_id_1');
	my $tid2  = $self->param('steering.target_id_2');
	my $tid1_weight = $self->param('steering.target_id_1_weight');
	my $tid2_weight = $self->param('steering.target_id_2_weight');
	if ( $tid1_weight eq "" ) { $tid1_weight = 0; }
	if ( $tid2_weight eq "" ) { $tid2_weight = 0; }
	if ( $self->is_valid() ) {
		my $targets;
		$targets->{$tid1} = $tid1_weight;
		$targets->{$tid2} = $tid2_weight;
		

		#delete current entries
		my $delete = $self->db->resultset('SteeringTarget')
			->search( { deliveryservice => $ds_id } );
		if ( defined($delete) ) {
			$delete->delete();
		}

		#add new entries
		foreach my $target ( keys %$targets ) {
			my $insert = $self->db->resultset('SteeringTarget')->create(
				{   deliveryservice => $ds_id,
					target          => $target,
					weight          => $targets->{$target},
				}
			);

			$insert->insert();
		}

		$self->flash(
			      message => "Successfully saved steering assignments for "
				. $self->get_ds_name($ds_id)
				. "!" );

		$self->redirect_to("/ds/$ds_id/steering");
	}
	else {
		&stash_role($self);
		my $target_name_1;
		my $target_name_2;
		my $target_id_1_weight;
		my $target_id_2_weight;
		if ($tid1 ) {
			$target_name_1 = $self->get_ds_name($tid1);
			$target_id_1_weight = $self->get_target_weight( $ds_id, $tid1 );
		}
		if ($tid2 ) {
			$target_name_2 = $self->get_ds_name($tid2);
			$target_id_2_weight = $self->get_target_weight( $ds_id, $tid2 );
		}
		$self->stash(
			steering => {
				ds_id           => $ds_id,
				ds_name         => $self->get_ds_name($ds_id),
				target_id_1     => $tid1,
				target_id_2     => $tid2,
				target_name_1   => $target_name_1,
				target_name_2   => $target_name_2,
				target_id_1_weight => $target_id_1_weight,
				target_id_2_weight => $target_id_2_weight
			},
			ds_data        => $self->get_deliveryservices(),
			fbox_layout    => 1
		);
		$self->render("steering/index");
	}
}

sub is_valid {
	my $self  = shift;

	#validate DSs are in the same CDN (same profile...)
	my $t1 = $self->param('steering.target_id_1');
	my $t2 = $self->param('steering.target_id_2');
	my $t1_profile;
	my $t2_profile;
	my $tid1_weight = $self->param('steering.target_id_1_weight');
	my $t1_name = $self->param('steering.target_name_1');
	my $tid2_weight = $self->param('steering.target_id_2_weight');
	my $t2_name = $self->param('steering.target_name_2');

	unless ( $t1 eq '' ) {
		$t1_profile = $self->get_ds_profile( $self->param('steering.target_id_1') );
	}
	unless ( $t2 eq '' ) {
		$t2_profile = $self->get_ds_profile( $self->param('steering.target_id_2') );
	}
	
	unless ( $t1 ) {
		$self->field('steering.target_id_1')->is_equal( "",  "Steering targets cannot be blank!" );
	}
	unless ( $t2 ) {
		$self->field('steering.target_id_2')->is_equal( "",  "Steering targets cannot be blank!" );
	}

	unless ( $t1_profile eq $t2_profile ) {
		$self->field('steering.target_id_1')->is_equal( "",  "Target Deliveryservices must be in the same CDN!" );
	}
	unless ( $tid1_weight eq int($tid1_weight) && $tid1_weight >= 0 ) {
		$self->field('steering.target_id_1_weight')->is_equal( "", "Error: \"$tid1_weight\" is not a valid integer of 0 or greater." );
	}
	unless ( $tid2_weight eq int($tid2_weight) && $tid2_weight >= 0 ) {
		$self->field('steering.target_id_2_weight')->is_equal( "", "Error: \"$tid2_weight\" is not a valid integer of 0 or greater." );
	}

	return $self->valid();
}

sub get_ds_profile {
	my $self  = shift;
	my $ds_id = shift;
	my $ds = $self->db->resultset('Deliveryservice')->search( { 'me.id' => $ds_id }, { prefetch => ['profile'] } )->single();
	return $ds->profile->name;
}

1;
