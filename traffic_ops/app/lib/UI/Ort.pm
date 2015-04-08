package UI::Ort;
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

use UI::Utils;
use Mojo::Base 'Mojolicious::Controller';
use Data::Dumper;

sub ort1 {
    my $self = shift;
	my $data_obj;
    my $host_name = $self->param('hostname');

	my %condition = ( 'servers.host_name' => $host_name );
	my $rs_profile = $self->db->resultset('Profile')->search( \%condition, { join => 'servers', columns => [qw/name id/] } );
  	my $row = $rs_profile->next;
    $data_obj->{'profile'}->{'name'} = $row->name;
    $data_obj->{'profile'}->{'id'} = $row->id;

#  $rs_route_type = $self->db->resultset('Type')->search( { -or => [ name => 'HTTP', name => 'HTTP_NO_CACHE', name => 'HTTP_LIVE' ] } );
    %condition = ( 'profile_parameters.profile' => $data_obj->{'profile'}->{'id'}, -or => ['name' => 'location', 'name' => 'CDN_name'] );
    my $rs_config = $self->db->resultset('Parameter')->search( \%condition, { join => 'profile_parameters' } );
    while ( my $row = $rs_config->next ) {
		if ( $row->name eq 'location' ) {
        	$data_obj->{'config_files'}->{$row->config_file}->{'location'} = $row->value;
    	}
		else {
        	$data_obj->{'other'}->{$row->name}= $row->value;
		}
	}
	$self->render( json => $data_obj );	
}

sub __get_json_parameter_list_by_host {
	my $self = shift;
	my $host = shift;
	my $value = shift;
	my $key_name = shift || "name";
	my $key_value = shift || "value";
	my $data_obj = [];
	
	my %condition = ( 'servers.host_name' => $host );
	my $rs_profile = $self->db->resultset('Profile')->search( \%condition, { join => 'servers', columns => [qw/name id/] } );
	my $row = $rs_profile->next;
	
	if (defined($row) && defined($row->id)) {
		my $id = $row->id;
		    
		%condition = ( 'profile_parameters.profile' => $id, 'config_file' => $value );
		my $rs_config = $self->db->resultset('Parameter')->search( \%condition, { join => 'profile_parameters' } );
		while ( my $row = $rs_config->next ) {
			# name = package name, value = package version
			push(@{$data_obj}, { $key_name => $row->name, $key_value => $row->value });
		}	
	}
	
	return($data_obj);
}

sub __get_json_parameter_by_host {
	my $self = shift;
	my $host = shift;
	my $parameter = shift;
	my $value = shift;
	my $key_name = shift || "name";
	my $key_value = shift || "value";
	my $data_obj;
	
	my %condition = ( 'servers.host_name' => $host );
	my $rs_profile = $self->db->resultset('Profile')->search( \%condition, { join => 'servers', columns => [qw/name id/] } );
	my $row = $rs_profile->next;
	my $id = $row->id;
	
	%condition = ( 'profile_parameters.profile' => $id, 'config_file' => $value, name => $parameter );
	my $rs_config = $self->db->resultset('Parameter')->search( \%condition, { join => 'profile_parameters' } );
	$row = $rs_config->next;
	
	if (defined($row) && defined($row->name) && defined($row->value)) {
		$data_obj->{$key_name} = $row->name;
		$data_obj->{$key_value} = $row->value;
	} else {
		# this is to ensure that we send an empty json response
		$data_obj->{""};
	}
	
	return($data_obj);
}

sub get_package_versions {
	my $self = shift;
	my $host_name = $self->param("hostname");	
	my $data_obj = __get_json_parameter_list_by_host($self, $host_name, "package", "name", "version");
	
	$self->render( json => $data_obj );
}

sub get_package_version {
	my $self = shift;
	my $host_name = $self->param("hostname");
	my $package = $self->param("package");
	my $data_obj = __get_json_parameter_by_host($self, $host_name, $package, "package", "name", "version");
	    
	$self->render( json => $data_obj );
}

sub get_chkconfig {
	my $self = shift;
	my $host_name = $self->param("hostname");	
	my $data_obj = __get_json_parameter_list_by_host($self, $host_name, "chkconfig");
	
	$self->render( json => $data_obj );
}

sub get_package_chkconfig {
	my $self = shift;
	my $host_name = $self->param("hostname");
	my $package = $self->param("package");
	my $data_obj = __get_json_parameter_by_host($self, $host_name, $package, "chkconfig");
	    
	$self->render( json => $data_obj );
}

1;
