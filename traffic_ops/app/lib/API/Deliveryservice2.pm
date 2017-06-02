package API::Deliveryservice2;
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
use Mojolicious::Validator;
use Mojolicious::Validator::Validation;
use Email::Valid;
use Validate::Tiny ':all';
use Data::Dumper;
use Common::ReturnCodes qw(SUCCESS ERROR);
use JSON;
use MojoPlugins::Response;
use UI::DeliveryService;
use Scalar::Util qw(looks_like_number);

my $valid_server_types = {
	edge => "EDGE",
	mid  => "MID",
};

# this structure maps the above types to the allowed metrics below
my $valid_metric_types = {
	origin_tps => "mid",
	ooff       => "mid",
};

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

			my $cdn_domain = $self->get_cdn_domain_by_ds_id($row->id);
			my $regexp_set = &UI::DeliveryService::get_regexp_set( $self, $row->id );
			my @example_urls = &UI::DeliveryService::get_example_urls( $self, $row->id, $regexp_set, $row, $cdn_domain, $row->protocol );

			push(
				@data, {
					"active"                   => \$row->active,
					"cacheurl"                 => $row->cacheurl,
					"ccrDnsTtl"                => $row->ccr_dns_ttl,
					"cdnId"                    => $row->cdn->id,
					"cdnName"                  => $row->cdn->name,
					"checkPath"                => $row->check_path,
					"displayName"              => $row->display_name,
					"dnsBypassCname"           => $row->dns_bypass_cname,
					"dnsBypassIp"              => $row->dns_bypass_ip,
					"dnsBypassIp6"             => $row->dns_bypass_ip6,
					"dnsBypassTtl"             => $row->dns_bypass_ttl,
					"dscp"                     => $row->dscp,
					"edgeHeaderRewrite"        => $row->edge_header_rewrite,
					"exampleURLs"              => \@example_urls,
					"geoLimitRedirectURL"      => $row->geolimit_redirect_url,
					"geoLimit"                 => $row->geo_limit,
					"geoLimitCountries"        => $row->geo_limit_countries,
					"geoProvider"              => $row->geo_provider,
					"globalMaxMbps"            => $row->global_max_mbps,
					"globalMaxTps"             => $row->global_max_tps,
					"httpBypassFqdn"           => $row->http_bypass_fqdn,
					"id"                       => $row->id,
					"infoUrl"                  => $row->info_url,
					"initialDispersion"        => $row->initial_dispersion,
					"ipv6RoutingEnabled"       => \$row->ipv6_routing_enabled,
					"lastUpdated"              => $row->last_updated,
					"logsEnabled"              => \$row->logs_enabled,
					"longDesc"                 => $row->long_desc,
					"longDesc1"                => $row->long_desc_1,
					"longDesc2"                => $row->long_desc_2,
					"matchList"                => \@matchlist,
					"maxDnsAnswers"            => $row->max_dns_answers,
					"midHeaderRewrite"         => $row->mid_header_rewrite,
					"missLat"                  => $row->miss_lat,
					"missLong"                 => $row->miss_long,
					"multiSiteOrigin"          => \$row->multi_site_origin,
					# "multiSiteOriginAlgorithm" => $row->multi_site_origin_algorithm,
					"orgServerFqdn"            => $row->org_server_fqdn,
					"originShield"             => $row->origin_shield,
					"profileId"                => $row->profile->id,
					"profileName"              => $row->profile->name,
					"profileDescription"       => $row->profile->description,
					"protocol"                 => $row->protocol,
					"qstringIgnore"            => $row->qstring_ignore,
					"rangeRequestHandling"     => $row->range_request_handling,
					"regexRemap"               => $row->regex_remap,
					"regionalGeoBlocking"      => \$row->regional_geo_blocking,
					"remapText"                => $row->remap_text,
					"signed"                   => \$row->signed,
					"sslKeyVersion"            => $row->ssl_key_version,
					"trRequestHeaders"         => $row->tr_request_headers,
					"trResponseHeaders"        => $row->tr_response_headers,
					"type"                     => $row->type->name,
					"typeId"                   => $row->type->id,
					"xmlId"                    => $row->xml_id
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
		$rs =
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

sub update_profileparameter {
	my $self   = shift;
	my $ds_id = shift;
	my $profile_id = shift;
	my $params = shift;

	&UI::DeliveryService::header_rewrite( $self, $ds_id, $profile_id, $params->{xmlId}, $params->{edgeHeaderRewrite}, "edge" );
	&UI::DeliveryService::header_rewrite( $self, $ds_id, $profile_id, $params->{xmlId}, $params->{midHeaderRewrite},  "mid" );
	&UI::DeliveryService::regex_remap( $self, $ds_id, $profile_id, $params->{xmlId}, $params->{regexRemap} );
	&UI::DeliveryService::cacheurl( $self, $ds_id, $profile_id, $params->{xmlId}, $params->{cacheurl} );
}

sub create {
	my $self   = shift;
	my $params = $self->req->json;

	if ( !&is_oper($self) ) {
		return $self->forbidden();
	}

	my ($transformed_params, $err) = (undef, undef);
	($transformed_params, $err) = $self->_check_params($params);
	if ( defined($err) ) {
		return $self->alert($err);
	}

	my $existing = $self->db->resultset('Deliveryservice')->search( { xml_id => $params->{xmlId} } )->get_column('xml_id')->single();
	if ( $existing ) {
			$self->alert("a delivery service with xmlId " . $params->{xmlId} . " already exists." );
	}

	my $value=$self->new_value($params, $transformed_params);
	my $insert = $self->db->resultset('Deliveryservice')->create($value);
	$insert->insert();
	my $new_id = $insert->id;

	if ( $new_id > 0 ) {
		my $patterns = $params->{matchList};
		foreach my $re (@$patterns) {
			my $type = $self->db->resultset('Type')->search( { name => $re->{type} } )->get_column('id')->single();
			my $regexp = $re->{pattern};

			my $insert = $self->db->resultset('Regex')->create(
				{
					pattern => $regexp,
					type    => $type,
				}
			);
			$insert->insert();
			my $new_re_id = $insert->id;

			my $de_re_insert = $self->db->resultset('DeliveryserviceRegex')->create(
				{
					regex           => $new_re_id,
					deliveryservice => $new_id,
					set_number      => defined($re->{setNumber}) ? $re->{setNumber} : 0,
				}
			);
			$de_re_insert->insert();
		}

		my $profile_id=$transformed_params->{ profile_id };
		$self->update_profileparameter($new_id, $profile_id, $params);

		my $cdn_rs = $self->db->resultset('Cdn')->search( { id => $transformed_params->{cdn_id} } )->single();
		my $dnssec_enabled = $cdn_rs->dnssec_enabled;
		if ( $dnssec_enabled == 1 ) {
			$self->app->log->debug("dnssec is enabled, creating dnssec keys");
			&UI::DeliveryService::create_dnssec_keys( $self, $cdn_rs->name, $params->{xmlId}, $new_id );
		}

		&log( $self, "Create deliveryservice with xml_id: " . $params->{xmlId}, " APICHANGE" );

		my $response = $self->get_response($new_id);
		return $self->success($response, "Delivery service was created: " . $new_id);
	}

	my $r = "Create Delivery Service fail, insert to database failed.";
	return $self->alert($r);
}

sub nodef_to_default {
	my $self    = shift;
	my $v       = shift;
	my $default = shift;

    return $v || $default;
}

sub get_types {
	my $self         = shift;
	my $use_in_table = shift;
	my $types;
	my $rs = $self->db->resultset('Type')->search( { use_in_table => $use_in_table } );
	while ( my $row = $rs->next ) {
		$types->{ $row->name } = $row->id;
	}
	return $types;
}

sub _check_params {
	my $self = shift;
	my $params = shift;
	my $ds_id = shift;
	my $transformed_params = undef;

	if ( !defined($params) ) {
		return (undef, "parameters should be in json format, please check!");
	}

	if ( !defined($params->{xmlId}) ) {
		return (undef, "parameter xmlId is must." );
	}

	if (!defined($params->{active})) {
		return (undef, "parameter active is must." );
	}

	if (looks_like_number($params->{active})) {
		if ($params->{active} == 1) {
			$transformed_params->{active} = 1;
		} elsif ($params->{active} == 0) {
			$transformed_params->{active} = 0;
		} else {
			return (undef, "active must be 1|0");
		}
	} else {
		if ($params->{active} eq "true") {
			$transformed_params->{active} = 1;
		} elsif ($params->{active} eq "false") {
			$transformed_params->{active} = 0;
		} else {
			return (undef, "active must be true|false");
		}
	}

	if ( defined($params->{type}) ) {
		my $rs = $self->get_types("deliveryservice");
		if ( !exists $rs->{ $params->{type} } ) {
			return (undef, "type (" . $params->{type} . ") must be deliveryservice type." );
		}
		else {
			$transformed_params->{type} = $rs->{ $params->{type} };
		}
	} else {
		return (undef, "parameter type is must." );
	}

	if (!defined($params->{protocol})) {
		return (undef, "parameter protocol is must." );
	}

	my $proto_num = $params->{protocol};

	if (!looks_like_number($proto_num) || $proto_num < 0 || $proto_num > 3) {
		return (undef, "protocol must be 0|1|2|3." );
	}


	if ( defined($params->{profileName}) ) {
		my $pname = $params->{profileName};
		my $profile =  $self->db->resultset('Profile')->search( { 'me.name' => $pname }, { prefetch => ['cdn'] } )->single();
		if ( !defined($profile) || $profile->cdn->name ne  $params->{cdnName} ) {
			return (undef, "profileName (" . $params->{profileName} . ") does not exist, or is not on the same CDN as " . $params->{cdnName} );
		}
		else {
			$transformed_params->{ profile_id } = $profile->id;
		}
	} else {
		return (undef, "parameter profileName is must." );
	}


	my $cdn_id = undef;
	if ( defined($params->{cdnName}) ) {
		$cdn_id = $self->db->resultset('Cdn')->search( { name => $params->{cdnName} } )->get_column('id')->single();
		if ( !defined $cdn_id ) {
			return (undef, "cdnName (" . $params->{cdnName} . ") does not exists." );
		} else {
			$transformed_params->{ cdn_id } = $cdn_id;
		}
	} else {
		return (undef, "parameter cdnName is must." );
	}

	if ( defined($params->{matchList}) ) {
		my $match_list = $params->{matchList};

		if ((scalar $match_list) == 0) {
			return (undef, "At least have 1 pattern in matchList.");
		}

		my $cdn_domain = undef;

		if (defined($ds_id)) {
			$cdn_domain = $self->get_cdn_domain_by_ds_id($ds_id);
		} else {
			my $profile_id = $self->get_profile_id_for_name($params->{profileName});
			$cdn_domain = $self->get_cdn_domain_by_profile_id($profile_id);
		}

		foreach my $match_item (@$match_list) {
			my $conflicting_regex = $self->find_existing_host_regex($match_item->{'type'}, $match_item->{'pattern'}, $cdn_domain, $cdn_id, $ds_id);
			if (defined($conflicting_regex)) {
				return(undef, "Another delivery service is already using host regex $conflicting_regex");
			}
		}
	} else {
		return (undef, "parameter matchList is must." );
	}

	if ( defined($params->{multiSiteOrigin}) ) {
		if ( !( ( $params->{multiSiteOrigin} eq "0" ) || ( $params->{multiSiteOrigin} eq "1" ) ) ) {
			return (undef, "multiSiteOrigin must be 0|1." );
		}
	} else {
		return (undef, "parameter multiSiteOrigin is must." );
	}

	if ( !defined($params->{displayName}) ) {
		return (undef, "parameter displayName is must." );
	}

	if ( defined($params->{orgServerFqdn}) ) {
		if ( $params->{orgServerFqdn} !~ /^https?:\/\// ) {
			return (undef, "orgServerFqdn must start with http(s)://" );
		}
	} else {
		return (undef, "parameter orgServerFqdn is must." );
	}

	if ( defined($params->{logsEnabled}) ) {
		if ( $params->{logsEnabled} eq "true" || $params->{logsEnabled} == 1 ) {
			$transformed_params->{logsEnabled} = 1;
		} elsif ( $params->{logsEnabled} eq "false" || $params->{logsEnabled} == 0 ) {
			$transformed_params->{logsEnabled} = 0;
		} else {
			return (undef, "logsEnabled must be true|false." );
		}
	} else {
		$transformed_params->{logsEnabled} = 0;
	}

	return ($transformed_params, undef);
}

sub new_value {
	my $self = shift;
	my $params = shift;
	my $transformed_params = shift;

	my $value = {
			xml_id                 => $params->{xmlId},
			display_name           => $params->{displayName},
			dscp                   => $self->nodef_to_default( $params->{dscp}, 0 ),
			signed                 => $self->nodef_to_default( $params->{signed}, 0 ),
			qstring_ignore         => $params->{qstringIgnore},
			geo_limit              => $params->{geoLimit},
			geo_limit_countries    => $params->{geoLimitCountries},
			geolimit_redirect_url  => $params->{geoLimitRedirectURL},
			geo_provider           => $params->{geoProvider},
			http_bypass_fqdn       => $params->{httpBypassFqdn},
			dns_bypass_ip          => $params->{dnsBypassIp},
			dns_bypass_ip6         => $params->{dnsBypassIp6},
			dns_bypass_cname       => $params->{dnsBypassCname},
			dns_bypass_ttl         => $params->{dnsBypassTtl},
			org_server_fqdn        => $params->{orgServerFqdn},
			multi_site_origin      => $params->{multiSiteOrigin},
			ccr_dns_ttl            => $params->{ccrDnsTtl},
			type                   => $transformed_params->{type},
			profile                => $transformed_params->{profile_id},
			cdn_id                 => $transformed_params->{cdn_id},
			global_max_mbps        => $self->nodef_to_default( $params->{globalMaxMbps}, 0 ),
			global_max_tps         => $self->nodef_to_default( $params->{globalMaxTps}, 0 ),
			miss_lat               => $params->{missLat},
			miss_long              => $params->{missLong},
			long_desc              => $params->{longDesc},
			long_desc_1            => $params->{longDesc1},
			long_desc_2            => $params->{longDesc2},
			max_dns_answers        => $self->nodef_to_default( $params->{maxDnsAnswers}, 0 ),
			info_url               => $params->{infoUrl},
			check_path             => $params->{checkPath},
			active                 => $transformed_params->{active},
			protocol               => $params->{protocol},
			ipv6_routing_enabled   => $params->{ipv6RoutingEnabled},
			range_request_handling => $params->{rangeRequestHandling},
			edge_header_rewrite    => $params->{edgeHeaderRewrite},
			mid_header_rewrite     => $params->{midHeaderRewrite},
			regex_remap            => $params->{regexRemap},
			origin_shield          => $params->{originShield},
			cacheurl               => $params->{cacheurl},
			remap_text             => $params->{remapText},
			initial_dispersion     => $params->{initialDispersion},
			regional_geo_blocking  => $self->nodef_to_default($params->{regionalGeoBlocking}, 0),
			ssl_key_version        => $params->{sslKeyVersion},
			tr_request_headers     => $params->{trRequestHeaders},
			tr_response_headers    => $params->{trResponseHeaders},
			logs_enabled           => $transformed_params->{logsEnabled},
		};

	return $value;
}

sub get_response {
	my $self   = shift;
	my $ds_id  = shift;

	my $response;
	my $rs = $self->db->resultset('Deliveryservice')->find( { id => $ds_id } );
	if ( defined($rs) ) {
		my $cdn_name = $self->db->resultset('Cdn')->search( { id => $rs->cdn_id } )->get_column('name')->single();

		$response->{id}                     = $rs->id;
		$response->{xmlId}                  = $rs->xml_id;
		$response->{active}                 = $rs->active==1 ? "true" : "false";
		$response->{dscp}                   = $rs->dscp;
		$response->{signed}                 = $rs->signed;
		$response->{qstringIgnore}          = $rs->qstring_ignore;
		$response->{geoLimit}               = $rs->geo_limit;
		$response->{geoLimitCountries}      = $rs->geo_limit_countries;
		$response->{geoLimitRedirectURL}    = $rs->geolimit_redirect_url;
		$response->{geoProvider}            = $rs->geo_provider;
		$response->{httpBypassFqdn}         = $rs->http_bypass_fqdn;
		$response->{dnsBypassIp}            = $rs->dns_bypass_ip;
		$response->{dnsBypassIp6}           = $rs->dns_bypass_ip6;
		$response->{dnsBypassTtl}           = $rs->dns_bypass_ttl;
		$response->{orgServerFqdn}          = $rs->org_server_fqdn;
		$response->{type}                   = $rs->type->name;
		$response->{profileName}            = $rs->profile->name;
		$response->{cdnName}                = $cdn_name;
		$response->{ccrDnsTtl}              = $rs->ccr_dns_ttl;
		$response->{globalMaxMbps}          = $rs->global_max_mbps;
		$response->{globalMaxTps}           = $rs->global_max_tps;
		$response->{longDesc}               = $rs->long_desc;
		$response->{longDesc1}              = $rs->long_desc_1;
		$response->{longDesc2}              = $rs->long_desc_2;
		$response->{maxDnsAnswers}          = $rs->max_dns_answers;
		$response->{infoUrl}                = $rs->info_url;
		$response->{missLat}                = $rs->miss_lat;
		$response->{missLong}               = $rs->miss_long;
		$response->{checkPath}              = $rs->check_path;
		$response->{protocol}               = $rs->protocol;
		$response->{sslKeyVersion}          = $rs->ssl_key_version;
		$response->{ipv6RoutingEnabled}     = $rs->ipv6_routing_enabled;
		$response->{rangeRequestHandling}   = $rs->range_request_handling;
		$response->{edgeHeaderRewrite}      = $rs->edge_header_rewrite;
		$response->{originShield}           = $rs->origin_shield;
		$response->{midHeaderRewrite}       = $rs->mid_header_rewrite;
		$response->{regexRemap}             = $rs->regex_remap;
		$response->{cacheurl}               = $rs->cacheurl;
		$response->{remapText}              = $rs->remap_text;
		$response->{multiSiteOrigin}        = $rs->multi_site_origin;
		$response->{displayName}            = $rs->display_name;
		$response->{trResponseHeaders}      = $rs->tr_response_headers;
		$response->{initialDispersion}      = $rs->initial_dispersion;
		$response->{dnsBypassCname}         = $rs->dns_bypass_cname;
		$response->{regionalGeoBlocking}    = $rs->regional_geo_blocking;
		$response->{trRequestHeaders}       = $rs->tr_request_headers;
		$response->{logsEnabled}            = $rs->logs_enabled==1 ? "true" : "false";
	}

	my @pats = ();
	$rs = $self->db->resultset('DeliveryserviceRegex')->search( { deliveryservice => $ds_id } );
	while ( my $row = $rs->next ) {
		push(
			@pats, {
				'pattern'   => $row->regex->pattern,
				'type'      => $row->regex->type->name,
				'setNumber' => $row->set_number,
			}
		);
	}
	$response->{matchList} = \@pats;

	return $response;
}

sub update {
	my $self   = shift;
	my $id     = $self->param('id');
	my $params = $self->req->json;

	if ( !&is_oper($self) ) {
		return $self->forbidden();
	}

	my $ds = $self->db->resultset('Deliveryservice')->find( { id => $id } );
	if ( !defined($ds) ) {
		return $self->not_found();
	}

	my ($transformed_params, $err) = (undef, undef);
	($transformed_params, $err) = $self->_check_params($params, $id);
	if ( defined($err) ) {
		return $self->alert($err);
	}

	my $existing = $self->db->resultset('Deliveryservice')->search( { xml_id => $params->{xmlId} } )->get_column('xml_id')->single();
	if ( $existing && $existing ne $ds->xml_id ) {
		$self->alert("a delivery service with xmlId " . $params->{xmlId} . " already exists." );
	}
	if ( $transformed_params->{ type } != $ds->type->id ) {
		return $self->alert("delivery service type can't be changed");
	}

	my $value=$self->new_value($params, $transformed_params);
	$ds->update($value);

	if ( defined($params->{matchList}) ) {
		my $patterns     = $params->{matchList};
		my $patterns_len = @$patterns;

		my $rs = $self->db->resultset('RegexesForDeliveryService')->search( {}, { bind => [$id] } );
		my $last_number = $rs->count;

		my $row = $rs->next;
		my $update_number;
		my $re;
		for ( $update_number=0; $update_number < $last_number && $update_number < $patterns_len; $update_number++ ) {
			$re = @$patterns[$update_number];
			my $type = $self->db->resultset('Type')->search( { name => $re->{type} } )->get_column('id')->single();
			my $update = $self->db->resultset('Regex')->find( { id => $row->id } );
			$update->update(
				{
					pattern => $re->{pattern},
					type    => $type,
				}
			);
			$update = $self->db->resultset('DeliveryserviceRegex')->find( { deliveryservice => $id, regex => $row->id } );
			$update->update( { set_number => defined($re->{setNumber}) ? $re->{setNumber} : 0 } );
			$row = $rs->next;
		}

		if ( $patterns_len > $last_number ) {
			for ( ; $update_number < $patterns_len; $update_number++ ) {
				$re = @$patterns[$update_number];
				my $type = $self->db->resultset('Type')->search( { name => $re->{type} } )->get_column('id')->single();
				my $insert = $self->db->resultset('Regex')->create(
					{
						pattern => $re->{pattern},
						type    => $type,
					}
				);
				$insert->insert();
				my $new_re_id = $insert->id;
				my $de_re_insert = $self->db->resultset('DeliveryserviceRegex')->create(
					{
						regex           => $new_re_id,
						deliveryservice => $id,
						set_number      => defined($re->{setNumber}) ? $re->{setNumber} : 0,
					}
				);
				$de_re_insert->insert();
			}
		}

		while ( $row ) {
			my $delete_re = $self->db->resultset('Regex')->search( { id => $row->id } );
			$delete_re->delete();
			$row = $rs->next;
		}
	}

	my $profile_id=$transformed_params->{ profile_id };
	$self->update_profileparameter($id, $profile_id, $params);

	&log( $self, "Update deliveryservice with xml_id: " . $params->{xmlId}, " APICHANGE" );

	my $response = $self->get_response($id);
	return $self->success($response, "Delivery service was updated: " . $id);
}

1;
