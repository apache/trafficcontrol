package API::Parameter;
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
use UI::Parameter;
use Mojo::Base 'Mojolicious::Controller';
use Data::Dumper;
use POSIX qw(strftime);
use Time::Local;
use LWP;
use MojoPlugins::Response;
use MojoPlugins::Job;
use Utils::Helper::ResponseHelper;

sub index {
	my $self    = shift;
	my $rs_data = $self->db->resultset("ProfileParameter")->search( undef, { prefetch => [ 'parameter', 'profile' ] } );
	my @data    = ();
	while ( my $row = $rs_data->next ) {
		my $value = $row->parameter->value;
		&UI::Parameter::conceal_secure_parameter_value( $self, $row->parameter->secure, \$value );
		push(
			@data, {
				"name"        => $row->parameter->name,
				"id"          => $row->parameter->id,
				"configFile"  => $row->parameter->config_file,
				"value"       => $value,
				"secure"      => $row->parameter->secure,
				"lastUpdated" => $row->parameter->last_updated,
			}
		);
	}
	$self->success( \@data );
}

sub profile {
	my $self         = shift;
	my $profile_name = $self->param('name');

	my $rs_data = $self->db->resultset("ProfileParameter")->search( { 'profile.name' => $profile_name }, { prefetch => [ 'parameter', 'profile' ] } );
	my @data = ();
	while ( my $row = $rs_data->next ) {
		my $value = $row->parameter->value;
		&UI::Parameter::conceal_secure_parameter_value( $self, $row->parameter->secure, \$value );
		push(
			@data, {
				"name"        => $row->parameter->name,
				"id"          => $row->parameter->id,
				"configFile"  => $row->parameter->config_file,
				"value"       => $value,
				"secure"      => $row->parameter->secure,
				"lastUpdated" => $row->parameter->last_updated,
			}
		);
	}
	$self->success( \@data );
}

1;
