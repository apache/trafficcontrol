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
use Data::Dumper;

sub index {
	my $self  = shift;
	my $ds_id = $self->param('id');
	#print STDERR Dumper($ds_id);

	&navbarpage($self);

	# select * from steering_target where deliveryservice = ds_id;
	#my $steering = { ds_id => $ds_id, ds_name => $self->get_ds_name($ds_id) };
	my $steering_obj;
	my @steering;
	my $st_rs = $self->db->resultset('SteeringTarget')->search( { deliveryservice => $ds_id, type => "weight" }, { order_by => 'value DESC' } );
	my $i = 0;
	if ( $st_rs > 0 ) {
		while ( my $row = $st_rs->next ) {
			$steering_obj->{"target_$i"}->{'target_id'} = $row->target;
			$steering_obj->{"target_$i"}->{'target_name'}   = $self->get_ds_name( $row->target );
			$steering_obj->{"target_$i"}->{'target_value'}   = $row->value;
			if (!defined($steering_obj->{"target_$i"}->{'target_value'})) { $steering_obj->{"target_$i"}->{'target_value'} = 0; }
			$steering_obj->{"target_$i"}->{'target_type'}   = $row->type;
			#print STDERR Dumper($steering_obj->{"target_$i"}->{'target_type'});
			push ( @steering, $steering_obj->{"target_$i"} );
			$i++;
		}
	}
	$st_rs = $self->db->resultset('SteeringTarget')->search( { deliveryservice => $ds_id, type => "order" }, { order_by => 'value ASC' } );
	if ( $st_rs > 0 ) {
		while ( my $row = $st_rs->next ) {
			$steering_obj->{"target_$i"}->{'target_id'} = $row->target;
			$steering_obj->{"target_$i"}->{'target_name'}   = $self->get_ds_name( $row->target );
			$steering_obj->{"target_$i"}->{'target_value'}   = $row->value;
			if (!defined($steering_obj->{"target_$i"}->{'target_value'})) { $steering_obj->{"target_$i"}->{'target_value'} = 0; }
			$steering_obj->{"target_$i"}->{'target_type'}   = $row->type;
			#print STDERR Dumper($steering_obj->{"target_$i"}->{'target_type'});
			push ( @steering, $steering_obj->{"target_$i"} );
			$i++;
		}
	}

	$self->stash(
		ds_id          => $ds_id,
		ds_name        => $self->get_ds_name($ds_id),
		steering       => \@steering,
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

sub testupdate {
	my $self = shift;
	my $ds_id = $self->param('id');
	my $st = $self->param('st');
	my @target_id = $self->param('st.target_id');
	my @target_value = $self->param('st.target_value');
	my @target_type = $self->param('st.target_type');
	my @target_delete = $self->param('st.target_delete');
	my @all = $self->req->params();
	print STDERR Dumper(\@all);
	print STDERR Dumper(\@target_id);
	print STDERR Dumper(\@target_value);
	print STDERR Dumper(\@target_type);
	foreach my $i (0 .. $#target_id) {
	#foreach my $id, $weight (@{$st->{'target_id'}}, @{$st->{'target_weight'}}) {
		print STDERR Dumper($target_id[$i]);
		print STDERR Dumper($target_value[$i]);
		print STDERR Dumper($target_type[$i]);
	}
	$self->redirect_to("/ds/$ds_id/steering");
}

sub update {
	my $self = shift;
	my $ds_id = $self->param('id');
	my @target_id = $self->param('st.target_id');
	my @target_value = $self->param('st.target_value');
	my @target_type = $self->param('st.target_type');
	#my $st = $self->param('st');
	my @targets;
	my $steering_obj;
	foreach my $i (0 .. $#target_id) {
	#foreach my $id (@{$st->{'target_id'}}) {
		#print STDERR Dumper($i);
		#print STDERR Dumper(@target_id[$i]);
		#print STDERR Dumper(@target_weight[$i]);
		if ( $target_id[$i] eq '' ) {
			#print STDERR Dumper("This one is blank");
			next;
		}
		if ( $target_value[$i] eq "" ) { $target_value[$i] = 0 };
		$steering_obj->{"target_$i"}->{'target_id'} = $target_id[$i];
		$steering_obj->{"target_$i"}->{'target_value'} = $target_value[$i];
		$steering_obj->{"target_$i"}->{'target_type'} = $target_type[$i];
		push ( @targets, $steering_obj->{"target_$i"} );
	}
	print STDERR Dumper(\@targets);
	#if ( 1 ==1 ) {
	if ( $self->is_valid(\@targets) ) {
		#delete current entries
		my $delete = $self->db->resultset('SteeringTarget')
			->search( { deliveryservice => $ds_id } );
		if ( defined($delete) ) {
			$delete->delete();
		}
		
		#add new entries
		#my $i = 0;
		foreach my $i ( keys @targets ) {
			print STDERR Dumper($targets[$i]->{'target_id'});
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
		print STDERR Dumper("at else somehow");
		my $steering_obj;
		my @steering;
		my $st_rs = $self->db->resultset('SteeringTarget')->search( { deliveryservice => $ds_id, type => "weight" }, { order_by => 'value DESC' } );
		my $i = 0;
		if ( $st_rs > 0 ) {
			while ( my $row = $st_rs->next ) {
				$steering_obj->{"target_$i"}->{'target_id'} = $row->target;
				$steering_obj->{"target_$i"}->{'target_name'}   = $self->get_ds_name( $row->target );
				$steering_obj->{"target_$i"}->{'target_value'}   = $row->value;
				if (!defined($steering_obj->{"target_$i"}->{'target_value'})) { $steering_obj->{"target_$i"}->{'target_value'} = 0; }
				$steering_obj->{"target_$i"}->{'target_type'}   = $row->type;
				#print STDERR Dumper($steering_obj->{"target_$i"}->{'target_type'});
				push ( @steering, $steering_obj->{"target_$i"} );
				$i++;
			}
		}
		$st_rs = $self->db->resultset('SteeringTarget')->search( { deliveryservice => $ds_id, type => "order" }, { order_by => 'value ASC' } );
		if ( $st_rs > 0 ) {
			while ( my $row = $st_rs->next ) {
				$steering_obj->{"target_$i"}->{'target_id'} = $row->target;
				$steering_obj->{"target_$i"}->{'target_name'}   = $self->get_ds_name( $row->target );
				$steering_obj->{"target_$i"}->{'target_value'}   = $row->value;
				if (!defined($steering_obj->{"target_$i"}->{'target_value'})) { $steering_obj->{"target_$i"}->{'target_value'} = 0; }
				$steering_obj->{"target_$i"}->{'target_type'}   = $row->type;
				#print STDERR Dumper($steering_obj->{"target_$i"}->{'target_type'});
				push ( @steering, $steering_obj->{"target_$i"} );
				$i++;
			}
		}
		&stash_role($self);
		$self->stash(
			ds_id          => $ds_id,
			ds_name        => $self->get_ds_name($ds_id),
			steering       => \@steering,
			ds_data        => $self->get_deliveryservices(),
			fbox_layout    => 1
		);
		$self->render("steering/index");
	}

	$self->redirect_to("/ds/$ds_id/steering");
}

sub old_update {
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
	my @targets = @{$_[0]};
	my $last_cdn;
	
	foreach my $i ( keys @targets ) {
		print STDERR Dumper($targets[$i]->{'target_id'});
		my $cdn = $self->get_ds_cdn( $targets[$i]->{'target_id'} );
		if ( defined($last_cdn) ) {
			if ( $cdn == $last_cdn ) {
				next;
			}
			else { 
				$self->flash(message => "Target Deliveryservices must be in the same CDN!" );
				return;
			}
		}
		else {
			$last_cdn = $cdn;
			next;
		}
		
	}

	return $self->valid;

}


sub get_ds_cdn {
	my $self  = shift;
	my $ds_id = shift;
	my $ds = $self->db->resultset('Deliveryservice')->search( { 'me.id' => $ds_id } )->single();
	return $ds->cdn_id;
}

1;
