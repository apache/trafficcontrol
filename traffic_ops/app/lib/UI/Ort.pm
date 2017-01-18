package UI::Ort;
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

use UI::Utils;
use Mojo::Base 'Mojolicious::Controller';
use Data::Dumper;

sub ort1 {
	my $self = shift;
	my $data_obj;
	my $host_name = $self->param('hostname');

	my %condition = ( 'me.host_name' => $host_name );
	my $rs_profile = $self->db->resultset('Server')
		->search( \%condition, { prefectch => [ 'cdn', 'profile' ] } );

	my $row = $rs_profile->next;
	if ($row) {
		my $cdn_name = defined( $row->cdn_id ) ? $row->cdn->name : "";

		$data_obj->{'profile'}->{'name'}   = $row->profile->name;
		$data_obj->{'profile'}->{'id'}     = $row->profile->id;
		$data_obj->{'other'}->{'CDN_name'} = $cdn_name;

		%condition = (
			'profile_parameters.profile' => $data_obj->{'profile'}->{'id'},
			-or                          => [ 'name' => 'location' ]
		);
		my $rs_config = $self->db->resultset('Parameter')
			->search( \%condition, { join => 'profile_parameters' } );
		while ( my $row = $rs_config->next ) {
			if ( $row->name eq 'location' ) {
				$data_obj->{'config_files'}->{ $row->config_file }
					->{'location'} = $row->value;
			}
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

	my $profile_id = $self->db->resultset('Server')->search( { host_name => $host } )->get_column('profile')->single();
	
	my %condition = ( 'profile_parameters.profile' => $profile_id, config_file => $value );
	my $rs_config = $self->db->resultset('Parameter')->search( \%condition, { join => 'profile_parameters' } );

	while ( my $row = $rs_config->next ) {
		push(@{$data_obj}, { $key_name => $row->name, $key_value => $row->value });
	}

	return($data_obj);
}

sub get_package_versions {
	my $self = shift;
	my $host_name = $self->param("hostname");	
	my $data_obj = __get_json_parameter_list_by_host($self, $host_name, "package", "name", "version");
	
	$self->render( json => $data_obj );
}

sub get_chkconfig {
	my $self = shift;
	my $host_name = $self->param("hostname");	
	my $data_obj = __get_json_parameter_list_by_host($self, $host_name, "chkconfig");
	
	$self->render( json => $data_obj );
}

1;
