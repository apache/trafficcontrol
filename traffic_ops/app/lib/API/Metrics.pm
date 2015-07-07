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
use Common::ReturnCodes qw(SUCCESS ERROR);

use Mojo::Base 'Mojolicious::Controller';

sub index {
	my $self        = shift;
	my $metric      = $self->param("metric");
	my $server_type = $self->param("server_type");
	my $m           = new Extensions::Delegate::Metrics($self);
	my ( $rc, $result ) = $m->get_etl_metrics();
	$self->app->log->debug( "result #-> " . Dumper($result) );
	if ( $rc == SUCCESS ) {
		return ( $self->success($result) );
	}
	else {
		return ( $self->alert($result) );
	}
}

1;
