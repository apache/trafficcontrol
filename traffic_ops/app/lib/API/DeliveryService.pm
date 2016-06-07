package API::DeliveryService;
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
use Mojolicious::Validator;
use Mojolicious::Validator::Validation;
use Email::Valid;
use Validate::Tiny ':all';
use Data::Dumper;
use Common::ReturnCodes qw(SUCCESS ERROR);
use JSON;
use MojoPlugins::Response;
use UI::DeliveryService;

sub delivery_services {
	my $self         = shift;
	my $id           = $self->param('id');
	my $logs_enabled = $self->param('logsEnabled');
	my $current_user = $self->current_user()->{username};

	my $rs;
	my $tm_user_id;
	my $forbidden;
	if ( defined($id) || defined($logs_enabled) ) {
		( $forbidden, $rs, $tm_user_id ) = $self->get_delivery_service_params( $current_user, $id, $logs_enabled );
	}
	else {
		( $rs, $tm_user_id ) = $self->get_delivery_services_by_user($current_user);
	}

	my @data;
	if ( defined($rs) ) {
		while ( my $row = $rs->next ) {
			my $cdn_name  = defined( $row->cdn_id ) ? $row->cdn->name : "";
			my $re_rs     = $row->deliveryservice_regexes;
			my @matchlist = ();
			while ( my $re_row = $re_rs->next ) {
				push(
					@matchlist, {
						type      => $re_row->regex->type->name,
						pattern   => $re_row->regex->pattern,
						setNumber => $re_row->set_number,
					}
				);
			}
			my $cdn_domain = &UI::DeliveryService::get_cdn_domain( $self, $row->id );
			my $regexp_set = &UI::DeliveryService::get_regexp_set( $self, $row->id );
			my @example_urls = &UI::DeliveryService::get_example_urls( $self, $row->id, $regexp_set, $row, $cdn_domain, $row->protocol );
			push(
				@data, {
					"id"                   => $row->id,
					"xmlId"                => $row->xml_id,
					"displayName"          => $row->display_name,
					"dscp"                 => $row->dscp,
					"signed"               => \$row->signed,
					"qstringIgnore"        => $row->qstring_ignore,
					"geoLimit"             => $row->geo_limit,
					"geoLimitCountries"    => $row->geo_limit_countries,
					"geoProvider"          => $row->geo_provider,
					"httpBypassFqdn"       => $row->http_bypass_fqdn,
					"dnsBypassIp"          => $row->dns_bypass_ip,
					"dnsBypassIp6"         => $row->dns_bypass_ip6,
					"dnsBypassCname"       => $row->dns_bypass_cname,
					"dnsBypassTtl"         => $row->dns_bypass_ttl,
					"orgServerFqdn"        => $row->org_server_fqdn,
					"multiSiteOrigin"      => $row->multi_site_origin,
					"ccrDnsTtl"            => $row->ccr_dns_ttl,
					"type"                 => $row->type->name,
					"profileName"          => $row->profile->name,
					"profileDescription"   => $row->profile->description,
					"cdnName"              => $cdn_name,
					"globalMaxMbps"        => $row->global_max_mbps,
					"globalMaxTps"         => $row->global_max_tps,
					"headerRewrite"        => $row->edge_header_rewrite,
					"edgeHeaderRewrite"    => $row->edge_header_rewrite,
					"midHeaderRewrite"     => $row->mid_header_rewrite,
					"trResponseHeaders"    => $row->tr_response_headers,
					"regexRemap"           => $row->regex_remap,
					"longDesc"             => $row->long_desc,
					"longDesc1"            => $row->long_desc_1,
					"longDesc2"            => $row->long_desc_2,
					"maxDnsAnswers"        => $row->max_dns_answers,
					"infoUrl"              => $row->info_url,
					"missLat"              => $row->miss_lat,
					"missLong"             => $row->miss_long,
					"checkPath"            => $row->check_path,
					"matchList"            => \@matchlist,
					"active"               => \$row->active,
					"protocol"             => $row->protocol,
					"ipv6RoutingEnabled"   => \$row->ipv6_routing_enabled,
					"rangeRequestHandling" => $row->range_request_handling,
					"cacheurl"             => $row->cacheurl,
					"remapText"            => $row->remap_text,
					"initialDispersion"    => $row->initial_dispersion,
					"exampleURLs"          => \@example_urls,
					"logsEnabled"          => \$row->logs_enabled,
				}
			);
		}
	}

	return defined($forbidden) ? $self->forbidden() : $self->success( \@data );
}

sub get_delivery_services_by_user {
	my $self         = shift;
	my $current_user = shift;

	my $tm_user_id;
	my $rs;
	if ( &is_privileged($self) ) {
		$rs = $self->db->resultset('Deliveryservice')->search( undef, { prefetch => [ 'cdn', 'deliveryservice_regexes' ], order_by => 'xml_id' } );
	}
	else {
		my $tm_user = $self->db->resultset('TmUser')->search( { username => $current_user } )->single();
		$tm_user_id = $tm_user->id;

		my @ds_ids = $self->db->resultset('DeliveryserviceTmuser')->search( { tm_user_id => $tm_user_id } )->get_column('deliveryservice')->all();
		$rs = $self->db->resultset('Deliveryservice')
			->search( { 'me.id' => { -in => \@ds_ids } }, { prefetch => [ 'cdn', 'deliveryservice_regexes' ], order_by => 'xml_id' } );
	}

	return ( $rs, $tm_user_id );
}

sub get_delivery_service_params {
	my $self         = shift;
	my $current_user = shift;
	my $id           = shift;
	my $logs_enabled = shift;

	# Convert to 1 or 0
	$logs_enabled = $logs_enabled ? 1 : 0;

	my $tm_user_id;
	my $rs;
	my $forbidden;
	my $condition;
	if ( &is_privileged($self) ) {
		if ( defined($id) ) {
			$condition = ( { 'me.id' => $id } );
		}
		else {
			$condition = ( { 'me.logs_enabled' => $logs_enabled } );
		}
		my @ds_ids = $rs =
			$self->db->resultset('Deliveryservice')->search( $condition, { prefetch => [ 'cdn', 'deliveryservice_regexes' ], order_by => 'xml_id' } );
	}
	elsif ( $self->is_delivery_service_assigned($id) ) {
		my $tm_user = $self->db->resultset('TmUser')->search( { username => $current_user } )->single();
		$tm_user_id = $tm_user->id;

		my @ds_ids =
			$self->db->resultset('DeliveryserviceTmuser')->search( { tm_user_id => $tm_user_id, deliveryservice => $id } )->get_column('deliveryservice')
			->all();
		$rs =
			$self->db->resultset('Deliveryservice')
			->search( { 'me.id' => { -in => \@ds_ids } }, { prefetch => [ 'cdn', 'deliveryservice_regexes' ], order_by => 'xml_id' } );
	}
	elsif ( !$self->is_delivery_service_assigned($id) ) {
		$forbidden = "true";
	}

	return ( $forbidden, $rs, $tm_user_id );
}

sub routing {
	my $self = shift;

	# get and pass { cdn_name => $foo } into get_routing_stats
	my $id = $self->param('id');

	if ( $self->is_valid_delivery_service($id) ) {
		if ( $self->is_delivery_service_assigned($id) || &is_admin($self) || &is_oper($self) ) {
			my $result = $self->db->resultset("Deliveryservice")->search( { 'me.id' => $id }, { prefetch => ['cdn'] } )->single();
			my $cdn_name = $result->cdn->name;

			# we expect type to be a dns or http type, but strip off any trailing bit
			my $stat_key = lc( $result->type->name );
			$stat_key =~ s/^(dns|http).*/$1/;
			$stat_key .= "Map";
			my $re_rs = $result->deliveryservice_regexes;
			my @patterns;
			while ( my $re_row = $re_rs->next ) {
				push( @patterns, $re_row->regex->pattern );
			}

			my $e = $self->get_routing_stats( { stat_key => $stat_key, patterns => \@patterns, cdn_name => $cdn_name } );
			if ( defined($e) ) {
				$self->alert($e);
			}
		}
		else {
			$self->forbidden();
		}
	}
	else {
		$self->not_found();
	}
}

sub capacity {
	my $self = shift;

	# get and pass { cdn_name => $foo } into get_cache_capacity
	my $id = $self->param('id');

	if ( $self->is_valid_delivery_service($id) ) {
		if ( $self->is_delivery_service_assigned($id) || &is_admin($self) || &is_oper($self) ) {
			my $result = $self->db->resultset("Deliveryservice")->search( { 'me.id' => $id }, { prefetch => ['cdn'] } )->single();
			my $cdn_name = $result->cdn->name;

			$self->get_cache_capacity( { delivery_service => $result->xml_id, cdn_name => $cdn_name } );
		}
		else {
			$self->forbidden();
		}
	}
	else {
		$self->not_found();
	}
}

sub health {
	my $self = shift;
	my $id   = $self->param('id');

	if ( $self->is_valid_delivery_service($id) ) {
		if ( $self->is_delivery_service_assigned($id) || &is_admin($self) || &is_oper($self) ) {
			my $result = $self->db->resultset("Deliveryservice")->search( { 'me.id' => $id }, { prefetch => ['cdn'] } )->single();
			my $cdn_name = $result->cdn->name;

			return ( $self->get_cache_health( { server_type => "caches", delivery_service => $result->xml_id, cdn_name => $cdn_name } ) );
		}
		else {
			$self->forbidden();
		}
	}
	else {
		$self->not_found();
	}
}

sub state {

	my $self = shift;
	my $id   = $self->param('id');

	if ( $self->is_valid_delivery_service($id) ) {
		if ( $self->is_delivery_service_assigned($id) || &is_admin($self) || &is_oper($self) ) {
			my $result      = $self->db->resultset("Deliveryservice")->search( { 'me.id' => $id }, { prefetch => ['cdn'] } )->single();
			my $cdn_name    = $result->cdn->name;
			my $ds_name     = $result->xml_id;
			my $rascal_data = $self->get_rascal_state_data( { type => "RASCAL", state_type => "deliveryServices", cdn_name => $cdn_name } );

			# scalar refs get converted into json booleans
			my $data = {
				enabled  => \0,
				failover => {
					enabled     => \0,
					configured  => \0,
					destination => undef,
					locations   => []
				}
			};

			if ( exists( $rascal_data->{$cdn_name} ) && exists( $rascal_data->{$cdn_name}->{state}->{$ds_name} ) ) {
				my $health_config = $self->get_health_config($cdn_name);
				my $c             = $rascal_data->{$cdn_name}->{config}->{deliveryServices}->{$ds_name};
				my $r             = $rascal_data->{$cdn_name}->{state}->{$ds_name};

				if ( exists( $health_config->{deliveryServices}->{$ds_name} ) ) {
					my $h = $health_config->{deliveryServices}->{$ds_name};

					if ( $h->{status} eq "REPORTED" ) {
						$data->{enabled} = \1;
					}

					if ( !$r->{isAvailable} ) {
						$data->{failover}->{enabled}   = \1;
						$data->{failover}->{locations} = $r->{disabledLocations};
					}

					if ( exists( $h->{"health.threshold.total.kbps"} ) ) {

						# get current kbps, calculate percent used
						$data->{failover}->{configured} = \1;
						push( @{ $data->{failover}->{limits} }, { metric => "total_kbps", limit => $h->{"health.threshold.total.kbps"} } );
					}

					if ( exists( $h->{"health.threshold.total.tps_total"} ) ) {

						# get current tps, calculate percent used
						$data->{failover}->{configured} = \1;
						push( @{ $data->{failover}->{limits} }, { metric => "total_tps", limit => $h->{"health.threshold.total.tps_total"} } );
					}

					if ( exists( $c->{bypassDestination} ) ) {
						my @k        = keys( %{ $c->{bypassDestination} } );
						my $type     = shift(@k);
						my $location = undef;

						if ( $type eq "DNS" ) {
							$location = $c->{bypassDestination}->{$type}->{ip};
						}
						elsif ( $type eq "HTTP" ) {
							my $port = ( exists( $c->{bypassDestination}->{$type}->{port} ) ) ? ":" . $c->{bypassDestination}->{$type}->{port} : "";
							$location = sprintf( "http://%s%s", $c->{bypassDestination}->{$type}->{fqdn}, $port );
						}

						$data->{failover}->{destination} = {
							type     => $type,
							location => $location
						};
					}
				}
			}

			$self->success($data);
		}
		else {
			$self->forbidden();
		}
	}
	else {
		$self->not_found();
	}
}

sub request {
	my $self     = shift;
	my $email_to = $self->req->json->{emailTo};
	my $details  = $self->req->json->{details};

	my $is_email_valid = Email::Valid->address($email_to);

	if ( !$is_email_valid ) {
		return $self->alert("Please provide a valid email address to send the delivery service request to.");
	}

	my ( $is_valid, $result ) = $self->is_deliveryservice_request_valid($details);

	if ($is_valid) {
		if ( $self->send_deliveryservice_request( $email_to, $details ) ) {
			return $self->success_message( "Delivery Service request sent to " . $email_to );
		}
	}
	else {
		return $self->alert($result);
	}
}

sub is_deliveryservice_request_valid {
	my $self    = shift;
	my $details = shift;

	my $rules = {
		fields => [
			qw/customer contentType deliveryProtocol routingType serviceDesc peakBPSEstimate peakTPSEstimate maxLibrarySizeEstimate originURL hasOriginDynamicRemap originTestFile hasOriginACLWhitelist originHeaders otherOriginSecurity queryStringHandling rangeRequestHandling hasSignedURLs hasNegativeCachingCustomization negativeCachingCustomizationNote serviceAliases rateLimitingGBPS rateLimitingTPS overflowService headerRewriteEdge headerRewriteMid headerRewriteRedirectRouter notes/
		],

		# Validation checks to perform
		checks => [

			# required deliveryservice request fields
			[
				qw/customer contentType deliveryProtocol routingType serviceDesc peakBPSEstimate peakTPSEstimate maxLibrarySizeEstimate originURL hasOriginDynamicRemap originTestFile hasOriginACLWhitelist queryStringHandling rangeRequestHandling hasSignedURLs hasNegativeCachingCustomization rateLimitingGBPS rateLimitingTPS/
			] => is_required("is required")

		]
	};

	# Validate the input against the rules
	my $result = validate( $details, $rules );

	if ( $result->{success} ) {
		return ( 1, $result->{data} );
	}
	else {
		return ( 0, $result->{error} );
	}
}

sub assign_servers {
	my $self      = shift;
	my $ds_xml_Id = $self->param('xml_id');
	my $params    = $self->req->json;

	if ( !defined($params) ) {
		return $self->alert("parameters are JSON format, please check!");
	}
	if ( !&is_oper($self) ) {
		return $self->alert("You must be an ADMIN or OPER to perform this operation!");
	}

	if ( !exists( $params->{server_names} ) ) {
		return $self->alert("Parameter 'server_names' is required.");
	}

	my $dsid = $self->db->resultset('Deliveryservice')->search( { xml_id => $ds_xml_Id } )->get_column('id')->single();
	if ( !defined($dsid) ) {
		return $self->alert( "DeliveryService[" . $ds_xml_Id . "] is not found." );
	}

	my @server_ids;
	my $svrs = $params->{server_names};
	foreach my $svr (@$svrs) {
		my $svr_id = $self->db->resultset('Server')->search( { host_name => $svr } )->get_column('id')->single();
		if ( !defined($svr_id) ) {
			return $self->alert( "Server[" . $svr . "] is not found in database." );
		}
		push( @server_ids, $svr_id );
	}

	# clean up
	my $delete = $self->db->resultset('DeliveryserviceServer')->search( { deliveryservice => $dsid } );
	$delete->delete();

	# assign servers
	foreach my $s_id (@server_ids) {
		my $insert = $self->db->resultset('DeliveryserviceServer')->create(
			{
				deliveryservice => $dsid,
				server          => $s_id,
			}
		);
		$insert->insert();
	}

	my $ds = $self->db->resultset('Deliveryservice')->search( { id => $dsid } )->single();
	&UI::DeliveryService::header_rewrite( $self, $ds->id, $ds->profile, $ds->xml_id, $ds->edge_header_rewrite, "edge" );

	my $response;
	$response->{xml_id} = $ds->xml_id;
	$response->{'server_names'} = \@$svrs;

	return $self->success($response);
}

1;
