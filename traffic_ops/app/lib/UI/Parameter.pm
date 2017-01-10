package UI::Parameter;
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

#Table View
sub index {
	my $self = shift;

	my $filter = $self->param('filter');
	my $value  = $self->param('byvalue');

	my $filter_title = "Cachegroup/Profile name";
	if ( defined($filter) ) {
		$self->stash(
			filter => $filter,
			value  => $value
		);
		$filter_title = $filter . " name";
		$filter_title =~ s/(^\w)/\U$1/x;
	}
	$self->stash( filter_title => $filter_title );
	&navbarpage($self);
}

sub view {
	my $self = shift;
	my $mode = $self->param('mode');
	my $id   = $self->param('id');

	my $rs_param = $self->db->resultset('Parameter')->search( { id => $id } );
	my $data = $rs_param->single;
	$self->conceal_secure_parameter_value( $data->{'_column_data'}{'secure'}, \$data->{'_column_data'}{'value'} );
	$self->stash( parameter => $data );

	&stash_role($self);

	my %assigned_profiles    = $self->get_assigned_profiles();
	my %assigned_cachegroups = $self->get_assigned_cachegroups();
	$self->stash(
		assigned_profiles    => \%assigned_profiles,
		assigned_cachegroups => \%assigned_cachegroups,
		fbox_layout          => 1,
	);
}

sub get_assigned_profiles {
	my $self        = shift;
	my $id          = $self->param('id');
	my @profile_ids = $self->db->resultset('ProfileParameter')->search( { parameter => $id } )->get_column('profile')->all;
	my %assigned_profiles;
	foreach my $p_id (@profile_ids) {
		my $profile = $self->db->resultset('Profile')->search( { id => $p_id } )->single;
		$assigned_profiles{$p_id} = {
			desc => $profile->description,
			name => $profile->name
		};
	}
	return %assigned_profiles;
}

sub get_assigned_cachegroups {
	my $self   = shift;
	my $id     = $self->param('id');
	my @cg_ids = $self->db->resultset('CachegroupParameter')->search( { parameter => $id } )->get_column('cachegroup')->all;
	my %assigned_cachegroups;
	foreach my $l_id (@cg_ids) {
		my $cachegroup = $self->db->resultset('Cachegroup')->search( { id => $l_id } )->single;
		$assigned_cachegroups{$l_id} = {
			short_name => $cachegroup->short_name,
			name       => $cachegroup->name
		};
	}
	return %assigned_cachegroups;
}

# Read
sub readparameter {
	my $self = shift;
	my @data;
	my $orderby = "name";
	$orderby = $self->param('orderby') if ( defined $self->param('orderby') );
	my $rs_data = $self->db->resultset("Parameter")->search( undef, { order_by => $orderby } );
	while ( my $row = $rs_data->next ) {
		my $value = $row->value;
		$self->conceal_secure_parameter_value( $row->secure, \$value );
		push(
			@data, {
				"id"           => $row->id,
				"name"         => $row->name,
				"config_file"  => $row->config_file,
				"value"        => $value,
				"secure"       => $row->secure,
				"last_updated" => $row->last_updated,
			}
		);
	}
	$self->render( json => \@data );
}

sub readparameter_for_profile {
	my $self         = shift;
	my $profile_name = $self->param('profile_name');

	my $rs_data = $self->db->resultset("ProfileParameter")->search( { 'profile.name' => $profile_name }, { prefetch => [ 'parameter', 'profile' ] } );
	my @data = ();
	while ( my $row = $rs_data->next ) {
		my $value = $row->parameter->value;
		$self->conceal_secure_parameter_value( $row->parameter->secure, \$value );
		push(
			@data, {
				"name"         => $row->parameter->name,
				"config_file"  => $row->parameter->config_file,
				"value"        => $value,
				"secure"       => $row->parameter->secure,
				"last_updated" => $row->parameter->last_updated,
			}
		);
	}
	$self->render( json => \@data );
}

# Delete
sub delete {
	my $self = shift;
	my $id   = $self->param('id');

	if ( !&is_oper($self) ) {
		$self->flash( alertmsg => "No can do. Get more privs." );
	}
	else {
		my $secure = $self->db->resultset('Parameter')->search( { id => $id } )->get_column('secure')->single();
		if ( (1==$secure) && !&is_admin($self) ) {
			$self->flash( alertmsg => "Forbidden. Admin role required to delete a secure parameter." );
		}
		else {
			my $p_name = $self->db->resultset('Parameter')->search( { id => $id } )->get_column('name')->single();
			my $delete = $self->db->resultset('Parameter')->search( { id => $id } );
			$delete->delete();
			&log( $self, "Delete parameter " . $p_name, "UICHANGE" );
		}
	}
	return $self->redirect_to('/close_fancybox.html');
}

# Update
sub update {
	my $self = shift;
	my $id   = $self->param('id');

	if ( $self->is_valid() ) {
		my $update = $self->db->resultset('Parameter')->find( { id => $self->param('id') } );
		$update->name( $self->param('parameter.name') );
		$update->config_file( $self->param('parameter.config_file') );
		$update->value( $self->param('parameter.value') );
		if ( &is_admin($self) ) {
			my $secure = defined( $self->param('parameter.secure') ) ? $self->param('parameter.secure') : 0;
			$update->secure($secure);
		}
		$update->update();

		# if the update has failed, we don't even get here, we go to the exception page.
		&log( $self, "Update parameter with name: " . $self->param('parameter.name'), "UICHANGE" );
		$self->flash( message => "Parameter updated successfully." );
		return $self->redirect_to( '/parameter/' . $self->param('id') );
	}
	else {
		&stash_role($self);
		my $rs_param             = $self->db->resultset('Parameter')->search( { id => $id } );
		my $data                 = $rs_param->single;
		$self->conceal_secure_parameter_value( $data->{'_column_data'}{'secure'}, \$data->{'_column_data'}{'value'} );
		my %assigned_profiles    = &get_assigned_profiles($self);
		my %assigned_cachegroups = &get_assigned_cachegroups($self);
		$self->stash(
			assigned_profiles    => \%assigned_profiles,
			assigned_cachegroups => \%assigned_cachegroups,
			parameter            => $data,
			fbox_layout          => 1
		);
		$self->render('/parameter/view');
	}
}

sub is_valid {
	my $self        = shift;
	my $mode        = shift;
	my $name        = $self->param('parameter.name');
	my $config_file = $self->param('parameter.config_file');
	my $value       = $self->param('parameter.value');
	my $secure      = defined( $self->param('parameter.secure') ) ? $self->param('parameter.secure') : 0;

	#Check permissions
	if ( !&is_oper($self) ) {
		$self->field('parameter.name')->is_equal( "", "You do not have the permissions to perform this operation!" );
	}
	if ( !&is_admin($self) ) {
		my $id = $self->param('id');
		if ( defined($id) ) {
			my $rs_param = $self->db->resultset('Parameter')->search( { id => $id } );
			my $data = $rs_param->single;
			my $original_secure = $data->{'_column_data'}{'secure'};
			if ( $original_secure == 1 ) {
				$self->field('parameter.name')->is_equal( "", "You must be an ADMIN to modify a secure parameter!" );
			}
		} else {
			if ( $secure == 1 ) {
				$self->field('parameter.name')->is_equal( "", "You must be an ADMIN to create a secure parameter!" );
			}
		}
	}

	#Check required fields
	$self->field('parameter.name')->is_required;
	$self->field('parameter.config_file')->is_required;
	$self->field('parameter.value')->is_required;

	#Make sure the same Parameter doesn't already exist
	my $existing_param =
		$self->db->resultset('Parameter')->search( { name => $name, value => $value, config_file => $config_file, secure => $secure } )->get_column('id')->single();
	if ($existing_param) {
		$self->field('parameter.name')
			->is_equal( "", "A parameter with the name \"$name\", config_file \"$config_file\", and value \"$value\" already exists." );
	}
	return $self->valid;
}

# Create
sub create {
	my $self = shift;

	my $new_id = -1;

	if ( $self->is_valid() ) {
		my $secure = defined( $self->param('parameter.secure') ) ? $self->param('parameter.secure') : 0;
		my $insert = $self->db->resultset('Parameter')->create(
			{
				name        => $self->param('parameter.name'),
				config_file => $self->param('parameter.config_file'),
				value       => $self->param('parameter.value'),
				secure      => $secure,
			}
		);
		$insert->insert();

		# if the insert has failed, we don't even get here, we go to the exception page.
		&log( $self, "Create parameter with name " . $self->param('parameter.name') . " and value " . $self->param('parameter.value'), "UICHANGE" );

		$new_id = $insert->id();

		if ( $new_id == -1 ) {
			my $referer = $self->req->headers->header('referer');
			return $self->redirect_to($referer);
		}
		else {
			$self->flash( message => "Parameter added successfully!  Please associate to profiles or cache groups now." );
			return $self->redirect_to( '/parameter/' . $new_id );
		}
	}
	else {
		&stash_role($self);
		$self->stash( parameter => { secure => 0 }, fbox_layout => 1 );
		$self->render('parameter/add');
	}
}

# add parameter view
sub add {
	my $self = shift;
	&stash_role($self);
	$self->stash( fbox_layout => 1, parameter => { secure => 0 } );
}

# conceal secure parameter value
sub conceal_secure_parameter_value {
	my $self = shift;
	my $secure = shift;
	my $value = shift;
	if ( $secure == 1 && !&is_admin($self) ) {
		$$value = '*********';
	}
}

1;
