package API::DeliveryService::KeysUrlSig;
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

use Mojo::Base 'Mojolicious::Controller';
use Data::Dumper;
use API::Keys;
use Utils::Helper;
use JSON;
use UI::Utils;
use constant URL_SIG_KEYS_BUCKET => "url_sig_keys";
use Exporter qw(import);
our @EXPORT_OK = qw(URL_SIG_KEYS_BUCKET);

sub view_by_xmlid {
	my $self                = shift;
	my $xml_id              = $self->param('xmlId');
	my $response_container  = $self->riak_get( URL_SIG_KEYS_BUCKET, $xml_id );
	my $url_sig_values_json = $response_container->{"response"};

	#$self->app->log->debug( "url_sig_values_json #-> " . Dumper($url_sig_values_json) );
	return $self->success($url_sig_values_json);
}

sub generate {
	my $self        = shift;
	my $xml_id      = $self->param('xmlId');
	my $config_file = $self->url_sig_config_file_name($xml_id);
	$self->app->log->info( "Generating New keys for config_file:  " . $config_file );

	my $current_user = $self->current_user()->{username};
	&log( $self, "Generated new url_sig_keys for " . $xml_id, "APICHANGE" );

	my $rs = $self->db->resultset("Deliveryservice")->find( { xml_id => $xml_id } );
	my $ds_id;
	if ( defined($rs) ) {
		$ds_id = $rs->id;
	}

	my $helper = new Utils::Helper( { mojo => $self } );

	# Admins can always do this, otherwise verify the user
	if ( ( defined($rs) && $helper->is_valid_delivery_service($ds_id) ) ) {
		if ( &is_admin($self) || $helper->is_delivery_service_assigned($ds_id) ) {
			my $url_sig_key_values_json = $self->generate_random_sigs_for_ds();
			if ( defined($rs) ) {

				#				$self->app->log->debug( "URL_SIG_KEYS_BUCKET, #-> " . URL_SIG_KEYS_BUCKET, );
				#				$self->app->log->debug( "config_file #-> " . $config_file );
				#				$self->app->log->debug( "url_sig_key_values_json #-> " . $url_sig_key_values_json );
				my $response_container = $self->riak_put( URL_SIG_KEYS_BUCKET, $config_file, $url_sig_key_values_json );
				my $response           = $response_container->{"response"};
				my $rc                 = $response->{_rc};
				if ( $rc eq '204' ) {
					return $self->success_message("Successfully generated and stored keys");
				}
				else {
					return $self->alert( $response->{_content} );
				}
			}
		}
		else {
			return $self->forbidden("Forbidden. Delivery service not assigned to user.");
		}
	}
	else {
		return $self->alert("Delivery Service '$xml_id' does not exist.");
	}
}

sub generate_random_sigs_for_ds {
	my $self  = shift;
	my $len   = 32;
	my @chars = ( 'a' .. 'z', 'A' .. 'Z', '0' .. '9', '_' );
	my $url_sig_keys;
	foreach my $i ( 0 .. 15 ) {
		my $v;
		foreach ( 1 .. $len ) {
			$v .= $chars[ rand @chars ];
		}
		my $k = "key$i";
		$self->app->log->info( "Generating..." . $k );
		$url_sig_keys->{$k} = $v;
		$i++;
	}
	return encode_json($url_sig_keys);
}

sub url_sig_config_file_name {
	my $self   = shift;
	my $xml_id = shift;

	return "url_sig_$xml_id.config";
}

1;
