package API::Metrics;
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
use Math::Round;
use Data::Dumper;
use POSIX qw(strftime);
use Carp qw(cluck confess);

use Mojo::Base 'Mojolicious::Controller';

my $valid_server_types = {
	edge => "EDGE",
	mid  => "MID",
};

# this structure maps the above types to the allowed metrics below
my $valid_metric_types = {
	origin_tps => "mid",
	ooff       => "mid",
};

sub index {
	my $self        = shift;
	my $server_type = $self->param("server_type");
	my $metric      = $self->param("metric");
	if ( exists( $valid_metric_types->{$metric} ) ) {
		$self->param( type => $valid_server_types->{$server_type} );
		return ( $self->etl_metrics() );
	}
	else {
		my @valid_types;
		foreach my $type ( keys %$valid_metric_types ) {
			push( @valid_types, "'" . $type . "'" );
		}
		return $self->alert( { error => "Invalid metric_type passed '" . $metric . "' valid types are " . join( ", ", @valid_types ) } );
	}
}

1;
