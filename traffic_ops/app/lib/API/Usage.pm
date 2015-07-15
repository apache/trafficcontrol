package API::Usage;
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
use JSON;
use Common::ReturnCodes qw(SUCCESS ERROR);
use Extensions::Delegate::Statistics;
use Utils::Helper::Extensions;
Utils::Helper::Extensions->use;

sub deliveryservice {
	my $self  = shift;
	my $ds_id = $self->param("ds_id");

	if ( $self->is_valid_delivery_service($ds_id) ) {
		if ( $self->is_delivery_service_assigned($ds_id) || &is_admin($self) || &is_oper($self) ) {

			my $stats = new Extensions::Delegate::Statistics($self);
			my ( $rc, $result ) = $stats->get_deliveryservice_usage();

			if ( $rc == SUCCESS ) {
				return $self->success($result);
			}
			else {
				return $self->alert($result);
			}
		}
		else {
			return $self->forbidden();
		}
	}
	else {
		$self->success( {} );
	}
}

1;
