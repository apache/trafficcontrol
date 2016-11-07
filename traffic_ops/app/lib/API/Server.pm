package API::Server;
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
use UI::Utils;
use Mojo::Base 'Mojolicious::Controller';
use Data::Dumper;
use POSIX qw(strftime);
use Time::Local;
use LWP;
use UI::ConfigFiles;
use UI::Tools;
use MojoPlugins::Response;
use MojoPlugins::Job;
use Utils::Helper::ResponseHelper;
use String::CamelCase qw(decamelize);
use Validate::Tiny ':all';

sub index {
	my $self         = shift;
	my $current_user = $self->current_user()->{username};
	my $ds_id        = $self->param('dsId');
	my $type         = $self->param('type');
	my $status       = $self->param('status');
	my $profile_id   = $self->param('profileId');

	my $servers;
	my $forbidden;
	if ( defined $ds_id ) {
		( $forbidden, $servers ) = $self->get_servers_by_dsid( $current_user, $ds_id, $status );
	}
	elsif ( defined $type ) {
		$servers = $self->get_servers_by_type( $current_user, $type, $status );
	}
	elsif ( defined $profile_id ) {
		( $forbidden, $servers ) = $self->get_servers_by_profile_id( $profile_id );
	}
	else {
		$servers = $self->get_servers_by_status( $current_user, $status );
	}

	if ( defined($forbidden) ) {
		return $self->forbidden($forbidden);
	}

	my @data;
	if ( defined($servers) ) {
		my $is_admin = &is_admin($self);
		while ( my $row = $servers->next ) {
			push(
				@data, {
					"cachegroup"     => $row->cachegroup->name,
					"cachegroupId"   => $row->cachegroup->id,
					"cdnId"          => $row->cdn->id,
					"cdnName"        => $row->cdn->name,
					"domainName"     => $row->domain_name,
					"guid"           => $row->guid,
					"hostName"       => $row->host_name,
					"httpsPort"      => $row->https_port,
					"id"             => $row->id,
					"iloIpAddress"   => $row->ilo_ip_address,
					"iloIpNetmask"   => $row->ilo_ip_netmask,
					"iloIpGateway"   => $row->ilo_ip_gateway,
					"iloUsername"    => $row->ilo_username,
					"iloPassword"    => $is_admin ? $row->ilo_password : "",
					"interfaceMtu"   => $row->interface_mtu,
					"interfaceName"  => $row->interface_name,
					"ip6Address"     => $row->ip6_address,
					"ip6Gateway"     => $row->ip6_gateway,
					"ipAddress"      => $row->ip_address,
					"ipNetmask"      => $row->ip_netmask,
					"ipGateway"      => $row->ip_gateway,
					"lastUpdated"    => $row->last_updated,
					"mgmtIpAddress"  => $row->mgmt_ip_address,
					"mgmtIpNetmask"  => $row->mgmt_ip_netmask,
					"mgmtIpGateway"  => $row->mgmt_ip_gateway,
					"offlineReason" => $row->offline_reason,
					"physLocation"   => $row->phys_location->name,
					"physLocationId" => $row->phys_location->id,
					"profile"        => $row->profile->name,
					"profileId"      => $row->profile->id,
					"profileDesc"    => $row->profile->description,
					"rack"           => $row->rack,
					"routerHostName" => $row->router_host_name,
					"routerPortName" => $row->router_port_name,
					"status"         => $row->status->name,
					"statusId"       => $row->status->id,
					"tcpPort"        => $row->tcp_port,
					"type"           => $row->type->name,
					"typeId"         => $row->type->id,
					"updPending"     => \$row->upd_pending
				}
			);
		}
	}

	return $self->success( \@data );
}

sub show {
	my $self = shift;
	my $id   = $self->param('id');

	my $rs_data  = $self->db->resultset("Server")->search( { id => $id } );
	my @data     = ();
	my $is_admin = &is_admin($self);
	while ( my $row = $rs_data->next ) {
		push(
			@data, {
				"cachegroup"     => $row->cachegroup->name,
				"cachegroupId"   => $row->cachegroup->id,
				"cdnId"          => $row->cdn->id,
				"cdnName"        => $row->cdn->name,
				"domainName"     => $row->domain_name,
				"guid"           => $row->guid,
				"hostName"       => $row->host_name,
				"httpsPort"      => $row->https_port,
				"id"             => $row->id,
				"iloIpAddress"   => $row->ilo_ip_address,
				"iloIpNetmask"   => $row->ilo_ip_netmask,
				"iloIpGateway"   => $row->ilo_ip_gateway,
				"iloUsername"    => $row->ilo_username,
				"iloPassword"    => $is_admin ? $row->ilo_password : "********",
				"interfaceMtu"   => $row->interface_mtu,
				"interfaceName"  => $row->interface_name,
				"ip6Address"     => $row->ip6_address,
				"ip6Gateway"     => $row->ip6_gateway,
				"ipAddress"      => $row->ip_address,
				"ipNetmask"      => $row->ip_netmask,
				"ipGateway"      => $row->ip_gateway,
				"lastUpdated"    => $row->last_updated,
				"mgmtIpAddress"  => $row->mgmt_ip_address,
				"mgmtIpNetmask"  => $row->mgmt_ip_netmask,
				"mgmtIpGateway"  => $row->mgmt_ip_gateway,
				"offline_reason" => $row->offline_reason,
				"physLocation"   => $row->phys_location->name,
				"physLocationId" => $row->phys_location->id,
				"profile"        => $row->profile->name,
				"profileId"      => $row->profile->id,
				"profileDesc"    => $row->profile->description,
				"rack"           => $row->rack,
				"routerHostName" => $row->router_host_name,
				"routerPortName" => $row->router_port_name,
				"status"         => $row->status->name,
				"statusId"       => $row->status->id,
				"tcpPort"        => $row->tcp_port,
				"type"           => $row->type->name,
				"typeId"         => $row->type->id,
				"updPending"     => \$row->upd_pending
			}
		);
	}
	$self->success( \@data );
}

sub update {
	my $self   = shift;
	my $id     = $self->param('id');
	my $params = $self->req->json;

	if ( !&is_oper($self) ) {
		return $self->forbidden();
	}

	my ( $is_valid, $result ) = $self->is_server_valid($params);

	if ( !$is_valid ) {
		return $self->alert($result);
	}

	my $server = $self->db->resultset('Server')->find( { id => $id } );
	if ( !defined($server) ) {
		return $self->not_found();
	}

	my $values = {
		cachegroup               	=> $params->{cachegroupId},
		cdn_id                     	=> $params->{cdnId},
		domain_name               	=> $params->{domainName},
		host_name                   => $params->{hostName},
		https_port           		=> $params->{httpsPort},
		ilo_ip_address 				=> $params->{iloIpAddress},
		ilo_ip_netmask         		=> $params->{iloIpNetmask},
		ilo_ip_gateway            	=> $params->{iloIpGateway},
		ilo_username            	=> $params->{iloUsername},
		ilo_password          		=> $params->{iloPassword},
		interface_mtu           	=> $params->{interfaceMtu},
		interface_name            	=> $params->{interfaceName},
		ip6_address              	=> $params->{ip6Address},
		ip6_gateway              	=> $params->{ip6Gateway},
		ip_address             		=> $params->{ipAddress},
		ip_netmask             		=> $params->{ipNetmask},
		ip_gateway                	=> $params->{ipGateway},
		mgmt_ip_address           	=> $params->{mgmtIpAddress},
		mgmt_ip_netmask          	=> $params->{mgmtIpNetmask},
		mgmt_ip_gateway           	=> $params->{mgmtIpGateway},
		offline_reason            	=> $params->{offlineReason},
		phys_location            	=> $params->{physLocationId},
		profile             		=> $params->{profileId},
		rack                     	=> $params->{rack},
		router_host_name       		=> $params->{routerHostName},
		router_port_name          	=> $params->{routerPortName},
		status                   	=> $params->{statusId},
		tcp_port                 	=> $params->{tcpPort},
		type                     	=> $params->{typeId},
		upd_pending               	=> $params->{updPending}
	};

	my $rs = $server->update($values);
	if ($rs) {
		my @response;
		push(
			@response, {
				"cachegroup"     => $rs->cachegroup->name,
				"cachegroupId"   => $rs->cachegroup->id,
				"cdnId"          => $rs->cdn->id,
				"cdnName"        => $rs->cdn->name,
				"domainName"     => $rs->domain_name,
				"guid"           => $rs->guid,
				"hostName"       => $rs->host_name,
				"httpsPort"      => $rs->https_port,
				"id"             => $rs->id,
				"iloIpAddress"   => $rs->ilo_ip_address,
				"iloIpNetmask"   => $rs->ilo_ip_netmask,
				"iloIpGateway"   => $rs->ilo_ip_gateway,
				"iloUsername"    => $rs->ilo_username,
				"iloPassword"    => $rs->ilo_password,
				"interfaceMtu"   => $rs->interface_mtu,
				"interfaceName"  => $rs->interface_name,
				"ip6Address"     => $rs->ip6_address,
				"ip6Gateway"     => $rs->ip6_gateway,
				"ipAddress"      => $rs->ip_address,
				"ipNetmask"      => $rs->ip_netmask,
				"ipGateway"      => $rs->ip_gateway,
				"lastUpdated"    => $rs->last_updated,
				"mgmtIpAddress"  => $rs->mgmt_ip_address,
				"mgmtIpNetmask"  => $rs->mgmt_ip_netmask,
				"mgmtIpGateway"  => $rs->mgmt_ip_gateway,
				"offlineReason"  => $rs->offline_reason,
				"physLocation"   => $rs->phys_location->name,
				"physLocationId" => $rs->phys_location->id,
				"profile"        => $rs->profile->name,
				"profileId"      => $rs->profile->id,
				"profileDesc"    => $rs->profile->description,
				"rack"           => $rs->rack,
				"routerHostName" => $rs->router_host_name,
				"routerPortName" => $rs->router_port_name,
				"status"         => $rs->status->name,
				"statusId"       => $rs->status->id,
				"tcpPort"        => $rs->tcp_port,
				"type"           => $rs->type->name,
				"typeId"         => $rs->type->id,
				"updPending"     => \$rs->upd_pending
			}
		);

		&log( $self, "Updated server [ '" . $rs->host_name . "' ] with id: " . $rs->id, "APICHANGE" );

		return $self->success( \@response, "Cachegroup update was successful." );
	}
	else {
		return $self->alert("Cachegroup update failed.");
	}
}

sub create {
	my $self   = shift;
	my $params = $self->req->json;

	if ( !&is_oper($self) ) {
		return $self->forbidden();
	}

	my ( $is_valid, $result ) = $self->is_server_valid($params);

	if ( !$is_valid ) {
		return $self->alert($result);
	}

	my $values = {
		cachegroup               	=> $params->{cachegroupId},
		cdn_id                     	=> $params->{cdnId},
		domain_name               	=> $params->{domainName},
		host_name                   => $params->{hostName},
		https_port           		=> $params->{httpsPort},
		ilo_ip_address 				=> $params->{iloIpAddress},
		ilo_ip_netmask         		=> $params->{iloIpNetmask},
		ilo_ip_gateway            	=> $params->{iloIpGateway},
		ilo_username            	=> $params->{iloUsername},
		ilo_password          		=> $params->{iloPassword},
		interface_mtu           	=> $params->{interfaceMtu},
		interface_name            	=> $params->{interfaceName},
		ip6_address              	=> $params->{ip6Address},
		ip6_gateway              	=> $params->{ip6Gateway},
		ip_address             		=> $params->{ipAddress},
		ip_netmask             		=> $params->{ipNetmask},
		ip_gateway                	=> $params->{ipGateway},
		mgmt_ip_address           	=> $params->{mgmtIpAddress},
		mgmt_ip_netmask          	=> $params->{mgmtIpNetmask},
		mgmt_ip_gateway           	=> $params->{mgmtIpGateway},
		offline_reason            	=> $params->{offlineReason},
		phys_location            	=> $params->{physLocationId},
		profile             		=> $params->{profileId},
		rack                     	=> $params->{rack},
		router_host_name       		=> $params->{routerHostName},
		router_port_name          	=> $params->{routerPortName},
		status                   	=> $params->{statusId},
		tcp_port                 	=> $params->{tcpPort},
		type                     	=> $params->{typeId},
		upd_pending               	=> $params->{updPending}
	};

	my $insert = $self->db->resultset('Server')->create($values);
	my $rs = $insert->insert();
	if ($rs) {
		my @response;
		push(
			@response, {
				"cachegroup"     => $rs->cachegroup->name,
				"cachegroupId"   => $rs->cachegroup->id,
				"cdnId"          => $rs->cdn->id,
				"cdnName"        => $rs->cdn->name,
				"domainName"     => $rs->domain_name,
				"guid"           => $rs->guid,
				"hostName"       => $rs->host_name,
				"httpsPort"      => $rs->https_port,
				"id"             => $rs->id,
				"iloIpAddress"   => $rs->ilo_ip_address,
				"iloIpNetmask"   => $rs->ilo_ip_netmask,
				"iloIpGateway"   => $rs->ilo_ip_gateway,
				"iloUsername"    => $rs->ilo_username,
				"iloPassword"    => $rs->ilo_password,
				"interfaceMtu"   => $rs->interface_mtu,
				"interfaceName"  => $rs->interface_name,
				"ip6Address"     => $rs->ip6_address,
				"ip6Gateway"     => $rs->ip6_gateway,
				"ipAddress"      => $rs->ip_address,
				"ipNetmask"      => $rs->ip_netmask,
				"ipGateway"      => $rs->ip_gateway,
				"lastUpdated"    => $rs->last_updated,
				"mgmtIpAddress"  => $rs->mgmt_ip_address,
				"mgmtIpNetmask"  => $rs->mgmt_ip_netmask,
				"mgmtIpGateway"  => $rs->mgmt_ip_gateway,
				"offlineReason"  => $rs->offline_reason,
				"physLocation"   => $rs->phys_location->name,
				"physLocationId" => $rs->phys_location->id,
				"profile"        => $rs->profile->name,
				"profileId"      => $rs->profile->id,
				"profileDesc"    => $rs->profile->description,
				"rack"           => $rs->rack,
				"routerHostName" => $rs->router_host_name,
				"routerPortName" => $rs->router_port_name,
				"status"         => $rs->status->name,
				"statusId"       => $rs->status->id,
				"tcpPort"        => $rs->tcp_port,
				"type"           => $rs->type->name,
				"typeId"         => $rs->type->id,
				"updPending"     => \$rs->upd_pending
			}
		);

		&log( $self, "Created server [ '" . $rs->host_name . "' ] with id: " . $rs->id, "APICHANGE" );

		return $self->success( \@response, "Server creation was successful." );
	}
	else {
		return $self->alert("Server creation failed.");
	}
}



sub get_servers_by_status {
	my $self              = shift;
	my $current_user      = shift;
	my $status            = shift;
	my $orderby           = $self->param('orderby') || "hostName";
	my $orderby_snakecase = lcfirst( decamelize($orderby) );

	my $servers;
	if ( &is_privileged($self) ) {
		my %criteria;
		if ( defined $status ) {
			$criteria{'status.name'} = $status;
		}

		$servers = $self->db->resultset('Server')->search(
			\%criteria, {
				prefetch => [ 'cdn', 'cachegroup', 'type', 'profile', 'status', 'phys_location' ],
				order_by => 'me.' . $orderby_snakecase,
			}
		);
	}
	else {
		my $tm_user = $self->db->resultset('TmUser')->search( { username => $current_user } )->single();
		my @ds_ids = $self->db->resultset('DeliveryserviceTmuser')->search( { tm_user_id => $tm_user->id } )->get_column('deliveryservice')->all();

		my @ds_servers =
			$self->db->resultset('DeliveryserviceServer')->search( { deliveryservice => { -in => \@ds_ids } } )->get_column('server')->all();

		my %criteria = ( 'me.id' => { -in => \@ds_servers } );
		if ( defined $status ) {
			$criteria{'status.name'} = $status;
		}
		$servers = $self->db->resultset('Server')->search(
			\%criteria, {
				prefetch => [ 'cdn', 'cachegroup', 'type', 'profile', 'status', 'phys_location' ],
				order_by => 'me.' . $orderby_snakecase,
			}
		);
	}

	return $servers;
}

sub get_servers_by_dsid {
	my $self              = shift;
	my $current_user      = shift;
	my $ds_id             = shift;
	my $status            = shift;
	my $orderby           = $self->param('orderby') || "hostName";
	my $orderby_snakecase = lcfirst( decamelize($orderby) );
	my $helper            = new Utils::Helper( { mojo => $self } );

	my @ds_servers;
	my $forbidden;
	if ( &is_privileged($self) || $self->is_delivery_service_assigned($ds_id) ) {
		@ds_servers = $self->db->resultset('DeliveryserviceServer')->search( { deliveryservice => $ds_id } )->get_column('server')->all();
	}
	else {
		$forbidden = "Forbidden. Delivery service not assigned to user.";
	}

	my $servers;
	if ( scalar(@ds_servers) ) {
		my $ds = $self->db->resultset('Deliveryservice')->search( { 'me.id' => $ds_id }, { prefetch => ['type'] } )->single();
		my %criteria = ( -or => [ 'me.id' => { -in => \@ds_servers } ] );

		# currently these are the ds types that bypass the mids
		my @types_no_mid = qw( HTTP_NO_CACHE HTTP_LIVE DNS_LIVE );
		if ( !grep { $_ eq $ds->type->name } @types_no_mid ) {

# if the delivery service employs mids, we're gonna pull mid servers too by pulling the cachegroups of the edges and finding those cachegroups parent cachegroup...
# then we see which servers have cachegroup in parent cachegroup list...that's how we find mids for the ds :)
			my @parent_cachegroup_ids = $self->db->resultset('ServersParentCachegroupList')->search( { 'me.server_id' => { -in => \@ds_servers } } )
				->get_column('parent_cachegroup_id')->all();
			push @{ $criteria{-or} }, { 'me.cachegroup' => { -in => \@parent_cachegroup_ids } };
		}

		if ( defined $status ) {
			$criteria{'status.name'} = $status;
		}

		$servers = $self->db->resultset('Server')->search(
			\%criteria, {
				prefetch => [ 'cdn', 'cachegroup', 'type', 'profile', 'status', 'phys_location' ],
				order_by => 'me.' . $orderby_snakecase,
			}
		);
	}

	return ( $forbidden, $servers );
}

sub get_servers_by_type {
	my $self              = shift;
	my $current_user      = shift;
	my $type              = shift;
	my $status            = shift;
	my $orderby           = $self->param('orderby') || "hostName";
	my $orderby_snakecase = lcfirst( decamelize($orderby) );

	my $servers;
	if ( &is_privileged($self) ) {
		my %criteria = ( 'type.name' => $type );
		if ( defined $status ) {
			$criteria{'status.name'} = $status;
		}

		$servers = $self->db->resultset('Server')->search(
			\%criteria, {
				prefetch => [ 'cdn', 'cachegroup', 'type', 'profile', 'status', 'phys_location' ],
				order_by => 'me.' . $orderby_snakecase,
			}
		);
	}
	else {
		my $tm_user = $self->db->resultset('TmUser')->search( { username => $current_user } )->single();
		my @ds_ids = $self->db->resultset('DeliveryserviceTmuser')->search( { tm_user_id => $tm_user->id } )->get_column('deliveryservice')->all();

		my @ds_servers =
			$self->db->resultset('DeliveryserviceServer')->search( { deliveryservice => { -in => \@ds_ids } } )->get_column('server')->all();

		my %criteria = ( 'me.id' => { -in => \@ds_servers }, 'type.name' => $type );
		if ( defined $status ) {
			$criteria{'status.name'} = $status;
		}

		$servers = $self->db->resultset('Server')->search(
			\%criteria, {
				prefetch => [ 'cdn', 'cachegroup', 'type', 'profile', 'status', 'phys_location' ],
				order_by => 'me.' . $orderby_snakecase,
			}
		);
	}

	return $servers;
}

sub totals {
	my $self = shift;

	my @data;
	my @rs = $self->db->resultset('ServerTypes')->search();
	foreach my $rs (@rs) {
		my $type_name = $rs->name;
		my $count     = $self->get_count_by_type($type_name);
		push(
			@data, {
				"type"  => $rs->name,
				"count" => $count,
			}
		);
	}

	return $self->success( \@data );

}

sub get_count_by_type {
	my $self      = shift;
	my $type_name = shift;
	return $self->db->resultset('Server')->search( { 'type.name' => $type_name }, { join => 'type' } )->count();
}

sub details_v11 {
	my $self = shift;
	my @data;
	my $isadmin   = &is_admin($self);
	my $host_name = $self->param('name');
	my $rs_data   = $self->db->resultset('Server')->search( { host_name => $host_name },
		{ prefetch => [ 'cachegroup', 'type', 'profile', 'status', 'phys_location', 'hwinfos', 'deliveryservice_servers' ], } );
	while ( my $row = $rs_data->next ) {

		my $serv = {
			"id"             => $row->id,
			"hostName"       => $row->host_name,
			"domainName"     => $row->domain_name,
			"tcpPort"        => $row->tcp_port,
			"httpsPort"      => $row->https_port,
			"xmppId"         => $row->xmpp_id,
			"xmppPasswd"     => $isadmin ? $row->xmpp_passwd : "********",
			"interfaceName"  => $row->interface_name,
			"ipAddress"      => $row->ip_address,
			"ipNetmask"      => $row->ip_netmask,
			"ipGateway"      => $row->ip_gateway,
			"ip6Address"     => $row->ip6_address,
			"ip6Gateway"     => $row->ip6_gateway,
			"interfaceMtu"   => $row->interface_mtu,
			"cachegroup"     => $row->cachegroup->name,
			"physLocation"   => $row->phys_location->name,
			"guid"           => $row->guid,
			"rack"           => $row->rack,
			"type"           => $row->type->name,
			"status"         => $row->status->name,
			"offline_reason" => $row->offline_reason,
			"profile"        => $row->profile->name,
			"profileDesc"    => $row->profile->description,
			"mgmtIpAddress"  => $row->mgmt_ip_address,
			"mgmtIpNetmask"  => $row->mgmt_ip_netmask,
			"mgmtIpGateway"  => $row->mgmt_ip_gateway,
			"iloIpAddress"   => $row->ilo_ip_address,
			"iloIpNetmask"   => $row->ilo_ip_netmask,
			"iloIpGateway"   => $row->ilo_ip_gateway,
			"iloUsername"    => $row->ilo_username,
			"iloPassword"    => $isadmin ? $row->ilo_password : "********",
			"routerHostName" => $row->router_host_name,
			"routerPortName" => $row->router_port_name,
		};
		my $hw_rs = $row->hwinfos;
		while ( my $hwinfo_row = $hw_rs->next ) {
			$serv->{hardwareInfo}->{ $hwinfo_row->description } = $hwinfo_row->val;
		}

		my $rs_ds_data = $row->deliveryservice_servers;
		while ( my $dsrow = $rs_ds_data->next ) {
			push( @{ $serv->{deliveryservices} }, $dsrow->deliveryservice->id );
		}

		push( @data, $serv );
	}
	$self->success(@data);
}

sub details {
	my $self              = shift;
	my $orderby           = $self->param('orderby') || "hostName";
	my $orderby_snakecase = lcfirst( decamelize($orderby) );
	my $limit             = $self->param('limit') || 1000;
	my @data;
	my $isadmin          = &is_admin($self);
	my $phys_location_id = $self->param('physLocationID');
	my $host_name        = $self->param('hostName');

	if ( !defined($phys_location_id) && !defined($host_name) ) {
		return $self->alert("Missing required fields: 'hostName' or 'physLocationID'");
	}

	my $rs_data = $self->db->resultset('Server')->search(
		[ { host_name => $host_name }, { phys_location => $phys_location_id } ], {
			prefetch => [ 'cachegroup', 'type', 'profile', 'status', 'phys_location', 'hwinfos', 'deliveryservice_servers' ],
			order_by => 'me.' . $orderby_snakecase
		}
	);

	if ( $rs_data->count() > 0 ) {

		while ( my $row = $rs_data->next ) {

			my $serv = {
				"id"             => $row->id,
				"hostName"       => $row->host_name,
				"domainName"     => $row->domain_name,
				"tcpPort"        => $row->tcp_port,
				"httpsPort"      => $row->https_port,
				"xmppId"         => $row->xmpp_id,
				"xmppPasswd"     => $isadmin ? $row->xmpp_passwd : "********",
				"interfaceName"  => $row->interface_name,
				"ipAddress"      => $row->ip_address,
				"ipNetmask"      => $row->ip_netmask,
				"ipGateway"      => $row->ip_gateway,
				"ip6Address"     => $row->ip6_address,
				"ip6Gateway"     => $row->ip6_gateway,
				"interfaceMtu"   => $row->interface_mtu,
				"cachegroup"     => $row->cachegroup->name,
				"physLocation"   => $row->phys_location->name,
				"guid"           => $row->guid,
				"rack"           => $row->rack,
				"type"           => $row->type->name,
				"status"         => $row->status->name,
				"offline_reason" => $row->offline_reason,
				"profile"        => $row->profile->name,
				"profileDesc"    => $row->profile->description,
				"mgmtIpAddress"  => $row->mgmt_ip_address,
				"mgmtIpNetmask"  => $row->mgmt_ip_netmask,
				"mgmtIpGateway"  => $row->mgmt_ip_gateway,
				"iloIpAddress"   => $row->ilo_ip_address,
				"iloIpNetmask"   => $row->ilo_ip_netmask,
				"iloIpGateway"   => $row->ilo_ip_gateway,
				"iloUsername"    => $row->ilo_username,
				"routerHostName" => $row->router_host_name,
				"routerPortName" => $row->router_port_name,
			};
			my $hw_rs = $row->hwinfos;
			while ( my $hwinfo_row = $hw_rs->next ) {
				$serv->{hardwareInfo}->{ $hwinfo_row->description } = $hwinfo_row->val;
			}

			my $rs_ds_data = $row->deliveryservice_servers;
			while ( my $dsrow = $rs_ds_data->next ) {
				push( @{ $serv->{deliveryservices} }, $dsrow->deliveryservice->id );
			}

			push( @data, $serv );
		}
		my $size = @data;
		$self->success( \@data, undef, $orderby, $limit, $size );
	}
	else {
		$self->success( [] );
	}
}

sub delete {
	my ( $params, $data, $err ) = ( undef, undef, undef );
	my $self = shift;

	if ( !&is_oper($self) ) {
		return $self->forbidden("Forbidden. You must have the operations role to perform this operation.");
	}

	my $id = $self->param('id');
	my $server = $self->db->resultset('Server')->find( { id => $id } );
	if ( !defined($server) ) {
		return $self->not_found();
	}
	my $delete = $self->db->resultset('Server')->search( { id => $id } );
	my $host_name = $delete->get_column('host_name')->single();
	$delete->delete();

	&log( $self, "Delete server with id:" . $id . " named " . $host_name, "APICHANGE" );

	return $self->success_message( "Server was deleted: " . $host_name );
}

sub postupdatequeue {
	my $self   = shift;
	my $params = $self->req->json;
	my $id     = $self->param('id');
	if ( !&is_oper($self) ) {
		return $self->forbidden("Forbidden. You must have the operations role to perform this operation.");
	}

	my $update = $self->db->resultset('Server')->find( { id => $id } );
	if ( !defined($update) ) {
		return $self->alert("Failed to find server id = $id");
	}

	my $setqueue = $params->{action};
	if ( !defined($setqueue) ) {
		return $self->alert("action needed, should be queue or dequeue.");
	}
	if ( $setqueue eq "queue" ) {
		$setqueue = 1;
	}
	elsif ( $setqueue eq "dequeue" ) {
		$setqueue = 0;
	}
	else {
		return $self->alert("action should be queue or dequeue.");
	}
	$update->update( { upd_pending => $setqueue } );

	my $response;
	$response->{serverId} = $id;
	$response->{action} = ( $setqueue == 1 ) ? "queue" : "dequeue";
	return $self->success($response);
}

sub get_servers_by_profile_id {
	my $self              = shift;
	my $profile_id        = shift;

	my $forbidden;
	my $servers;
	if ( !&is_oper($self) ) {
		$forbidden = "Forbidden. You must have the operations role to perform this operation.";
		return ( $forbidden, $servers );
	}

	my $servers = $self->db->resultset('Server')->search( { profile => $profile_id } );
	return ( $forbidden, $servers );
}

sub is_server_valid {
	my $self   = shift;
	my $params = shift;

	if (!$self->is_valid_server_type($params->{typeId})) {
		return ( 0, "Invalid server type" );
	}

	my $rules = {
		fields => [ qw/cachegroupId cdnId domainName hostName httpsPort iloIpAddress iloIpNetmask iloIpGateway iloUsername iloPassword interfaceMtu interfaceName ip6Address ip6Gateway ipAddress ipNetmask ipGateway mgmtIpAddress mgmtIpNetmask mgmtIpGateway offlineReason physLocationId profileId rack routerHostName routerPortName statusId tcpPort typeId updPending/ ],

		# Validation checks to perform
		checks => [
			cachegroupId => [ is_required("is required") ],
			cdnId => [ is_required("is required") ],
			domainName => [ is_required("is required") ],
			hostName => [ is_required("is required") ],
			interfaceMtu => [ is_required("is required") ],
			interfaceName => [ is_required("is required") ],
			ipAddress => [ is_required("is required") ],
			ipNetmask => [ is_required("is required") ],
			ipGateway => [ is_required("is required") ],
			physLocationId => [ is_required("is required") ],
			profileId => [ is_required("is required") ],
			statusId => [ is_required("is required") ],
			typeId => [ is_required("is required") ],
			updPending => [ is_required("is required") ]
		]
	};

	# Validate the input against the rules
	my $result = validate( $params, $rules );

	if ( $result->{success} ) {
		return ( 1, $result->{data} );
	}
	else {
		return ( 0, $result->{error} );
	}
}

sub is_valid_server_type {
	my $self     = shift;
	my $type_id = shift;

	my $rs = $self->db->resultset("Type")->find( { id => $type_id } );
	if ( defined($rs) && ( $rs->use_in_table eq "server" ) ) {
		return 1;
	}
	return 0;
}

1;
