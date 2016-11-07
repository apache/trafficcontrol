package UI::Health;
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
use Utils::Rascal;
use Utils::CCR;
use Utils::Helper;
use Data::Dumper;
use Carp qw(cluck confess);

sub healthfull {
	my $self = shift;
	my $data_obj;
	my %condition = ( 'profile.name' => [ { like => 'MID%' }, { like => 'EDGE%' } ] );
	my $rs = $self->db->resultset('Server')->search( \%condition, { prefetch => [ { 'profile' => undef } ] } );
	while ( my $row = $rs->next ) {
		push( @{ $data_obj->{ $row->profile->name }->{'servers'} }, $row->host_name . "." . $row->domain_name );
	}
	%condition = ( 'parameter.config_file' => 'rascal.properties' );
	$rs = $self->db->resultset('ProfileParameter')->search( \%condition, { prefetch => [ { 'parameter' => undef }, { 'profile' => undef } ] } );
	while ( my $row = $rs->next ) {
		push(
			@{ $data_obj->{ $row->profile->name }->{'parameters'} },
			( { 'name' => $row->parameter->name, 'value' => $row->parameter->value, 'last_updated' => $row->parameter->last_updated } )
		);
	}
	$self->render( json => $data_obj );
}

sub healthprofile {
	my $self = shift;
	my $data_obj;
	my %condition = ( 'parameter.config_file' => 'rascal.properties' );
	my $rs = $self->db->resultset('ProfileParameter')->search( \%condition, { prefetch => [ { 'parameter' => undef }, { 'profile' => undef } ] } );
	while ( my $row = $rs->next ) {
		push(
			@{ $data_obj->{ $row->profile->name } },
			( { 'name' => $row->parameter->name, 'value' => $row->parameter->value, 'last_updated' => $row->parameter->last_updated } )
		);
	}
	$self->render( json => $data_obj );
}

sub rascal_config {
	my $self     = shift;
	my $data_obj = $self->get_health_config( $self->param('cdnname') );
	$self->render( json => $data_obj );
}

1;
