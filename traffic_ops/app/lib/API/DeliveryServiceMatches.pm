package API::DeliveryServiceMatches;
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
use Utils::Tenant;
use UI::DeliveryService;
use Mojo::Base 'Mojolicious::Controller';
use Data::Dumper;
use Common::ReturnCodes qw(SUCCESS ERROR);

sub index {
	my $self = shift;
	my $format = $self->param("format") || "";

	my $rs;
    $rs = $self->db->resultset('Deliveryservice')->search( undef, { prefetch => [ 'cdn', 'type', { 'deliveryservice_regexes' => 'regex' }  ], order_by => 'xml_id' } );

    my @matches;
    my $tenant_utils = Utils::Tenant->new($self);
    my $tenants_data = $tenant_utils->create_tenants_data_from_db();
    while ( my $row = $rs->next ) {
        if (!$tenant_utils->is_ds_resource_accessible($tenants_data, $row->tenant_id)) {
            next;
        }
        my @match_patterns;
        my $cdn_name = defined( $row->cdn_id ) ? $row->cdn->name : "";
        my $xml_id   = defined( $row->xml_id ) ? $row->xml_id    : "";
        my $active   = defined( $row->active ) ? $row->active    : "";

        if ($active) {

            # Attach the remap_text host for Teak
            my $remap_text = $row->remap_text;

            #print "remap_text #-> (" . $remap_text . ")\n";
            if ( defined($remap_text) ) {

                #$self->app->log->debug( "remap_text #-> " . Dumper($remap_text) );
                my ($remap_text_match) = $remap_text =~ /regex_map http:\/\/(.*?[-])(.*)/;
                push( @match_patterns, $remap_text_match );
            }
            else {
                my $type_name = defined( $row->type ) ? $row->type->name : "";

                my $regexes = $row->deliveryservice_regexes;

                # if delivery service of type ANY_MAP
                # pull from regex_map_text, and only match for (http://)(quick-)(.*)
                my $regex_pattern;
                while ( my $r = $regexes->next ) {
                    $regex_pattern = $r->regex->pattern;
                    my $match_pattern = $self->convert_regex_to_match_pattern($regex_pattern);
                    push( @match_patterns, $match_pattern );
                }

            }

            #print "match_patterns #-> (" . Dumper(@match_patterns) . ")\n";

            # Check if there are any match_patterns or that there is a Delivery Service regex has '\.*' in it.
            if (@match_patterns) {
                my $xml_id_underscores = $xml_id =~ s/[-]/_/g;
                my $delivery_service->{dsName} = $xml_id;
                $delivery_service->{patterns} = \@match_patterns;
                push( @matches, $delivery_service );
            }

        }
    }

    if ( $format =~ /file/ ) {
        return $self->render( json => \@matches );
    }
    else {
        return $self->success( \@matches );
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
