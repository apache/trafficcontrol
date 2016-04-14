package API::DeliveryServiceMatches;
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
use UI::DeliveryService;
use Mojo::Base 'Mojolicious::Controller';
use Data::Dumper;
use Common::ReturnCodes qw(SUCCESS ERROR);

sub index {
	my $self   = shift;
	my $format = $self->param("format");

	my $rs;
	if ( &is_privileged($self) ) {
		$rs = $self->db->resultset('Deliveryservice')->search( undef, { prefetch => [ 'cdn', 'deliveryservice_regexes' ], order_by => 'xml_id' } );

		my @matches;
		while ( my $row = $rs->next ) {
			my $cdn_name = defined( $row->cdn_id ) ? $row->cdn->name : "";
			my $xml_id   = defined( $row->xml_id ) ? $row->xml_id    : "";

			my $regexes = $row->deliveryservice_regexes;

			my @match_patterns;
			while ( my $r = $regexes->next ) {
				my $match_pattern = $self->convert_regex_to_match_pattern( $r->regex->pattern );
				push( @match_patterns, $match_pattern );
			}
			my $delivery_service->{dsName} = $xml_id;
			$delivery_service->{patterns} = \@match_patterns;
			push( @matches, $delivery_service );
		}

		if ( $format =~ /file/ ) {
			return $self->render( json => \@matches );
		}
		else {
			return $self->success( \@matches );
		}
	}
	else {
		return $self->forbidden();
	}

}

sub convert_regex_to_match_pattern {
	my $self  = shift;
	my $regex = shift;

	my $match_pattern = $regex;
	$match_pattern =~ s/\\//g;
	$match_pattern =~ s/\.\*//g;
	return $match_pattern;
}

1;
