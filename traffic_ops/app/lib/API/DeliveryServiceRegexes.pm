package API::DeliveryServiceRegexes;
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
use UI::DeliveryService;
use Mojo::Base 'Mojolicious::Controller';
use Data::Dumper;
use Common::ReturnCodes qw(SUCCESS ERROR);

sub index {
	my $self = shift;

	my $rs;
	if ( &is_privileged($self) ) {
		$rs = $self->db->resultset('Deliveryservice')->search( undef, { prefetch => [ 'cdn', 'deliveryservice_regexes' ], order_by => 'xml_id' } );

		my @regexes;
		while ( my $row = $rs->next ) {
			my $cdn_name = defined( $row->cdn_id ) ? $row->cdn->name : "";
			my $xml_id   = defined( $row->xml_id ) ? $row->xml_id    : "";

			my $re_rs = $row->deliveryservice_regexes;

			my @matchlist;
			while ( my $re_row = $re_rs->next ) {
				push(
					@matchlist, {
						type      => $re_row->regex->type->name,
						pattern   => $re_row->regex->pattern,
						setNumber => $re_row->set_number,
					}
				);
			}
			my $delivery_service->{dsName} = $xml_id;
			$delivery_service->{regexes} = \@matchlist;
			push( @regexes, $delivery_service );
		}

		return $self->success( \@regexes );
	}
	else {
		return $self->forbidden("Forbidden. Insufficent privileges.");
	}

}

1;
