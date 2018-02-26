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
use Utils::Tenant;
use JSON;
use UI::Utils;
use Utils::Tenant;
use constant URL_SIG_KEYS_BUCKET => "url_sig_keys";
use Exporter qw(import);
our @EXPORT_OK = qw(URL_SIG_KEYS_BUCKET);


sub view_by_id {
	my $self                = shift;
	my $id              = $self->param('id');

	my $tenant_utils = Utils::Tenant->new($self);
	my $tenants_data = $tenant_utils->create_tenants_data_from_db();

	if ( !&is_privileged($self) and !$tenant_utils->use_tenancy()) {

		# check to see if deliveryservice is assigned to user, if not return forbidden
		return $self->forbidden("Forbidden. Delivery service not assigned to user.") if ( !$self->is_delivery_service_assigned($id) );
	}
	
	my $rs = $self->db->resultset("Deliveryservice")->find( { id => $id } ); 
	
	if ($tenant_utils->use_tenancy() and !$tenant_utils->is_ds_resource_accessible($tenants_data, $rs->tenant_id)) {
		return $self->forbidden("Forbidden. Delivery-service tenant is not available to the user.");
	}							
	
	my $xml_id;
	if ( defined($rs) ) {
		$xml_id = $rs->xml_id;
	}
	else {
		return $self->not_found("Delivery Service '$id' does not exist.");
	}

	my $config_file = $self->url_sig_config_file_name($xml_id);
	my $response_container  = $self->riak_get( URL_SIG_KEYS_BUCKET, $config_file );
	my $rc                  = $response_container->{"response"}->{_rc};
	if ( $rc eq '200' ) {
		my $url_sig_values_json = decode_json( $response_container->{"response"}->{_content} );
		return $self->success($url_sig_values_json);
	} else {
		my $error_msg = $response_container->{"response"}->{_content};
		$self->app->log->debug("received error code '$rc' from riak: '$error_msg'");
		return $self->success({}, "No url sig keys found");
	} 
}

sub view_by_xmlid {
	my $self                = shift;
	my $xml_id              = $self->param('xmlId');

	my $rs = $self->db->resultset("Deliveryservice")->find( { xml_id => $xml_id } );
	if ( !defined($rs) ) {
		return $self->not_found("Delivery Service '$xml_id' does not exist.");
	}
	my $tenant_utils = Utils::Tenant->new($self);
	my $tenants_data = $tenant_utils->create_tenants_data_from_db();
	if (!$tenant_utils->is_ds_resource_accessible($tenants_data, $rs->tenant_id)) {
		return $self->forbidden("Forbidden. Delivery-service tenant is not available to the user.");
	}

	my $config_file = $self->url_sig_config_file_name($xml_id);
	my $response_container  = $self->riak_get( URL_SIG_KEYS_BUCKET, $config_file );
	my $rc                  = $response_container->{"response"}->{_rc};
	if ( $rc eq '200' ) {
		my $url_sig_values_json = decode_json( $response_container->{"response"}->{_content} );
		return $self->success($url_sig_values_json);
	} else {
		my $error_msg = $response_container->{"response"}->{_content};
		$self->app->log->debug("received error code '$rc' from riak: '$error_msg'");
		return $self->success({}, "No url sig keys found");
	} 
}

sub copy_url_sig_keys {
	my $self                = shift;
	my $xml_id              = $self->param('xmlId'); #copying to this service
	my $copy_from_xml_id    = $self->param('copyFromXmlId'); # copying from this service

	if ( !&is_oper($self) ) {
		return $self->forbidden();
	}

	my $current_user = $self->current_user()->{username};
	my $is_oper = &is_oper($self);
	#check ds and generate config file name
	my $rs = $self->db->resultset("Deliveryservice")->find( { xml_id => $xml_id } ); 
	my $ds_id;
	if ( defined($rs) ) {
		$ds_id = $rs->id;
	}
	else {
		return $self->alert("Delivery Service '$xml_id' does not exist.");
	}
	my $tenant_utils = Utils::Tenant->new($self);
	my $tenants_data = $tenant_utils->create_tenants_data_from_db();
	if (!$tenant_utils->is_ds_resource_accessible($tenants_data, $rs->tenant_id)) {
		return $self->forbidden("Forbidden. Delivery-service tenant is not available to the user.");
	}
	my $config_file = $self->url_sig_config_file_name($xml_id);

	#check ds to copy from and generate config file name
	my $copy_rs = $self->db->resultset("Deliveryservice")->find( { xml_id => $copy_from_xml_id } );
	my $copy_ds_id;
	if ( defined($copy_rs) ) {
		$copy_ds_id = $copy_rs->id;
	}
	else {
		return $self->alert("Delivery Service to copy from '$copy_from_xml_id' does not exist.");
	}
	if (!$tenant_utils->is_ds_resource_accessible($tenants_data, $copy_rs->tenant_id)) {
		return $self->forbidden("Forbidden. Source delivery-service tenant is not available to the user.");
	}
	my $copy_config_file = $self->url_sig_config_file_name($copy_from_xml_id);

	my $helper = new Utils::Helper( { mojo => $self } );
	my $url_sig_key_values_json;

	#verify we can copy keys out
	if ( $helper->is_valid_delivery_service($copy_ds_id) ) {
		if ( $is_oper || $helper->is_delivery_service_assigned($copy_ds_id) || $tenant_utils->use_tenancy()) {
			my $response_container = $self->riak_get( URL_SIG_KEYS_BUCKET, $copy_config_file ); # verify this
			my $rc                 = $response_container->{"response"}->{_rc};
			if ( $rc eq '200' ) {
				$url_sig_key_values_json = $response_container->{"response"}->{_content};
			}
			else {
				my $error_msg = $response_container->{"response"}->{_content};
				$self->app->log->warn("received error code '$rc' from riak: '$error_msg'");
			}
		}
		else {
			return $self->forbidden("Forbidden. Delivery service to copy from not assigned to user.");
		}
	}
	else {
		return $self->alert("Delivery Service to copy from '$copy_from_xml_id' is not a valid delivery service.");
	}

	if ( defined($url_sig_key_values_json) ) { # verify we got keys copied
		# Ops can always do this, otherwise verify the user
		if ( $helper->is_valid_delivery_service($ds_id) ) {
			if ( $is_oper || $helper->is_delivery_service_assigned($ds_id) || $tenant_utils->use_tenancy()) {
				$self->app->log->debug( "url_sig_key_values_json #-> " . $url_sig_key_values_json );
				my $response_container = $self->riak_put( URL_SIG_KEYS_BUCKET, $config_file, $url_sig_key_values_json );
				my $response           = $response_container->{"response"};
				my $rc                 = $response->{_rc};
				if ( $rc eq '204' ) {
					&log( $self, "copied url_sig_keys from " . $copy_from_xml_id . " to " . $xml_id, "APICHANGE" );
					return $self->success_message("Successfully copied and stored keys");
				}
				else {
					my $error_msg = $response->{_content};
					$self->app->log->warn("received error code '$rc' from riak: '$error_msg'");
					return $self->alert( $response->{_content} );
				}
			}
			else {
				return $self->forbidden("Forbidden. Delivery service not assigned to user.");
			}
		}
		else {
			return $self->alert("Delivery Service '$xml_id' is not a valid delivery service.");
		}
	}
	else {
		return $self->alert("Unable to retrieve keys from Delivery Service '$copy_from_xml_id'");
	}
}

sub generate {
	my $self        = shift;
	my $xml_id      = $self->param('xmlId');
	my $config_file = $self->url_sig_config_file_name($xml_id);
	$self->app->log->info( "Generating New keys for config_file:  " . $config_file );

	if ( !&is_oper($self) ) {
		return $self->forbidden();
	}

	my $current_user = $self->current_user()->{username};

	my $tenant_utils = Utils::Tenant->new($self);
	my $tenants_data = $tenant_utils->create_tenants_data_from_db();
	my $rs = $self->db->resultset("Deliveryservice")->find( { xml_id => $xml_id } );
	my $ds_id;
	if ( defined($rs) ) {
		$ds_id = $rs->id;
		if (!$tenant_utils->is_ds_resource_accessible($tenants_data, $rs->tenant_id)) {
			return $self->forbidden("Forbidden. Delivery-service tenant is not available to the user.");
		}
	}

	my $helper = new Utils::Helper( { mojo => $self } );

	# Ops can always do this, otherwise verify the user
	if ( ( defined($rs) && $helper->is_valid_delivery_service($ds_id) ) ) {
		if ( &is_oper($self) || $helper->is_delivery_service_assigned($ds_id) || $tenant_utils->use_tenancy()) {
			my $url_sig_key_values_json = $self->generate_random_sigs_for_ds();
			if ( defined($rs) ) {

				#				$self->app->log->debug( "URL_SIG_KEYS_BUCKET, #-> " . URL_SIG_KEYS_BUCKET, );
				#				$self->app->log->debug( "config_file #-> " . $config_file );
				#				$self->app->log->debug( "url_sig_key_values_json #-> " . $url_sig_key_values_json );
				my $response_container = $self->riak_put( URL_SIG_KEYS_BUCKET, $config_file, $url_sig_key_values_json );
				my $response           = $response_container->{"response"};
				my $rc                 = $response->{_rc};
				if ( $rc eq '204' ) {
					&log( $self, "Generated new url_sig_keys for " . $xml_id, "APICHANGE" );
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
