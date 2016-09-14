package API::Server;
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

sub index {
	my $self         = shift;
	my $current_user = $self->current_user()->{username};
	my $ds_id        = $self->param('dsId');
	my $type         = $self->param('type');
	my $status       = $self->param('status');

	my $servers;
	my $forbidden;
	if ( defined $ds_id ) {
		( $forbidden, $servers ) = $self->get_servers_by_dsid( $current_user, $ds_id, $status );
	}
	elsif ( defined $type ) {
		$servers = $self->get_servers_by_type( $current_user, $type, $status );
	}
	else {
		$servers = $self->get_servers( $current_user, $status );
	}

	my @data;
	if ( defined($servers) ) {
		my $is_admin = &is_admin($self);
		while ( my $row = $servers->next ) {
			my $cdn_name = defined( $row->cdn_id ) ? $row->cdn->name : "";

			push(
				@data, {
					"id"             => $row->id,
					"hostName"       => $row->host_name,
					"domainName"     => $row->domain_name,
					"tcpPort"        => $row->tcp_port,
					"httpsPort"      => $row->https_port,
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
					"cdnName"        => $cdn_name,
					"mgmtIpAddress"  => $row->mgmt_ip_address,
					"mgmtIpNetmask"  => $row->mgmt_ip_netmask,
					"mgmtIpGateway"  => $row->mgmt_ip_gateway,
					"iloIpAddress"   => $row->ilo_ip_address,
					"iloIpNetmask"   => $row->ilo_ip_netmask,
					"iloIpGateway"   => $row->ilo_ip_gateway,
					"iloUsername"    => $row->ilo_username,
					"iloPassword"    => $is_admin ? $row->ilo_password : "********",
					"routerHostName" => $row->router_host_name,
					"routerPortName" => $row->router_port_name,
					"lastUpdated"    => $row->last_updated,

				}
			);
		}
	}

	return defined($forbidden) ? $self->forbidden("Forbidden. Delivery service not assigned to user.") : $self->success( \@data );
}

sub get_servers {
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
	my $ds_id              = shift;
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
		$forbidden = "true";
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
			my @parent_cachegroup_ids = $self->db->resultset('ServersParentCachegroupList')->search( { 'me.server_id' => { -in => \@ds_servers } } )->get_column('parent_cachegroup_id')->all();
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
				"httpsPort"        => $row->https_port,
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

sub check_server_params {
	my $self        = shift;
	my $json        = shift;
	my $update_base = shift;
	my %params      = %{$json};
	my $err         = undef;
	my %errFields   = ();

	# Required field checks
	if ( !defined( $json->{'hostName'} ) ) {
		$errFields{'hostName'} = 'is required';
	}
	if ( !defined( $json->{'domainName'} ) ) {
		$errFields{'domainName'} = 'is required';
	}
	if ( !defined( $json->{'cachegroup'} ) ) {
		$errFields{'cachegroup'} = 'is required';
	}
	if ( !defined( $json->{'interfaceName'} ) ) {
		$errFields{'interfaceName'} = 'is required';
	}
	if ( !defined( $json->{'ipAddress'} ) ) {
		$errFields{'ipAddress'} = 'is required';
	}
	if ( !defined( $json->{'ipNetmask'} ) ) {
		$errFields{'ipNetmask'} = 'is required';
	}
	if ( !defined( $json->{'ipGateway'} ) ) {
		$errFields{'ipGateway'} = 'is required';
	}
	if ( !defined( $json->{'interfaceMtu'} ) ) {
		$errFields{'interfaceMtu'} = 'is required';
	}
	if ( !defined( $json->{'physLocation'} ) ) {
		$errFields{'physLocation'} = 'is required';
	}
	if ( !defined( $json->{'type'} ) ) {
		$errFields{'type'} = 'is required';
	}
	if ( !defined( $json->{'profile'} ) ) {
		$errFields{'profile'} = 'is required';
	}
	if ( !defined( $json->{'cdnName'} ) ) {
		$errFields{'cdnName'} = 'is required';
	}
	if ( %errFields ) {
		return (\%params, \%errFields);
	}

	# Valid value checks
	if ( defined( $json->{'interfaceMtu'} ) ) {
		if ( $json->{'interfaceMtu'} != '1500' && $json->{'interfaceMtu'} != '9000' ) {
			return ( \%params, "'interfaceMtu' '$json->{'interfaceMtu'}' not equal to 1500 or 9000!" );
		}
	}

	if ( defined( $json->{'tcpPort'} ) ) {
		$params{'tcpPort'} = int( $json->{'tcpPort'} );
	}
	elsif ( !defined($update_base) ) {
		$params{'tcpPort'} = 80;
	}
	if ( defined( $json->{'httpsPort'} ) ) {
		$params{'httpsPort'} = int( $json->{'httpsPort'} );
	}
	elsif ( !defined($update_base) ) {
		$params{'httpsPort'} = 443;
	}

	eval { $params{'cachegroup'} = $self->db->resultset('Cachegroup')->search( { name => $json->{'cachegroup'} } )->get_column('id')->single(); };
	if ( $@ || ( !defined( $params{'cachegroup'} ) ) ) { # $@ holds Perl errors
		return ( \%params, "'cachegroup' $json->{'cachegroup'} not found!" );
	}

	eval { $params{'cdnId'} = $self->db->resultset('Cdn')->search( { name => $json->{'cdnName'} } )->get_column('id')->single(); };

	eval { $params{'type'} = &type_id( $self, $json->{'type'} ); };
	if ( $@ || ( !defined( $params{'type'} ) ) ) { # $@ holds Perl errors
		return ( \%params, "'type' $json->{'type'} not found!" );
	}

	eval { $params{'profile'} = &profile_id( $self, $json->{'profile'} ); };
	if ( $@ || ( !defined( $params{'profile'} ) ) ) { # $@ holds Perl errors
		return ( \%params, "'profile' $json->{'profile'} not found!" );
	}

	eval {
		$params{'physLocation'} = $self->db->resultset('PhysLocation')->search( { name => $json->{'physLocation'} } )->get_column('id')->single();
	};
	if ( $@ || ( !defined( $params{'physLocation'} ) ) ) { # $@ holds Perl errors
		return ( \%params, "'physLocation' $json->{'physLocation'} not found!" );
	}

	# IP address checks
	foreach my $ipstr (
		$json->{'ipAddress'},     $json->{'ipNetmask'},      $json->{'ipGateway'},      $json->{'iloIpAddress'}, $json->{'iloIpNetmask'},
		$json->{'iloIpGateway'}, $json->{'mgmtIpAddress'}, $json->{'mgmtIpNetmask'}, $json->{'mgmtIpGateway'}
		)
	{
		if ( !defined($ipstr) || $ipstr eq "" ) {
			next;
		}    # already checked for mandatory.
		if ( !&is_ipaddress($ipstr) ) {
			return ( \%params, $ipstr . " is not a valid IPv4 address or netmask" );
		}
	}

	if (   defined( $json->{'ip6Address'} )
		&& $json->{'ip6Address'} ne ""
		&& !&is_ip6address( $json->{'ip6Address'} ) )
	{
		return ( \%params, "Address " . $json->{'ip6Address'} . " is not a valid IPv6 address " );
	}
	if (   defined( $json->{'ip6Gateway'} )
		&& $json->{'ip6Gateway'} ne ""
		&& !&is_ip6address( $json->{'ip6Gateway'} ) )
	{
		return ( \%params, "Address " . $json->{'ip6Address'} . " is not a valid IPv6 address " );
	}

	my $ip_used =
		$self->db->resultset('Server')
			->search(
				{ -and =>
					[
						'me.ip_address' => $json->{'ipAddress'},
						'profile.name' => $json->{'profile'},
						'me.id' => { '!=' => (defined($update_base)) ? $update_base->id : 0 }
					]
				},
				{
					join   => [ 'profile' ]
				}
		)->single();
	if ( $ip_used ) {
		return ( \%params, $json->{'ipAddress'} . " is already being used by a server with the same profile" );
	}

	if ( defined( $json->{'ip6Address'} ) && $json->{'ip6Address'} ne "" ) {
		my $ip6_used =
			$self->db->resultset('Server')
				->search(
				{ -and =>
					[
						'me.ip6_address' => $json->{'ip6Address'},
						'profile.name' => $json->{'profile'},
						'me.id' => { '!=' => (defined($update_base)) ? $update_base->id : 0 }
					]
				},
				{
					join   => [ 'profile' ]
				}
			)->single();
		if ( $ip6_used ) {
			return ( \%params, $json->{'ip6Address'} . " is already being used by a server with the same profile" );
		}
	}

	# Netmask checks
	if ( defined( $json->{'ipNetmask'} )
		&& $json->{'mgmtIpNetmask'} ne ""
		&& !&is_netmask( $json->{'ipNetmask'} ) ) {
		return ( \%params, $json->{'ipNetmask'} . " is not a valid netmask" );
	}
	if (   defined( $json->{'iloIpNetmask'} )
		&& $json->{'iloIpNetmask'} ne ""
		&& !&is_netmask( $json->{'iloIpNetmask'} ) )
	{
		return ( \%params, $json->{'iloIpNetmask'} . " is not a valid netmask" );
	}
	if (   defined( $json->{'mgmtIpNetmask'} )
		&& $json->{'mgmtIpNetmask'} ne ""
		&& !&is_netmask( $json->{'mgmtIpNetmask'} ) )
	{
		return ( \%params, $json->{'mgmtIpNetmask'} . " is not a valid netmask" );
	}

	if ( ( defined( $json->{'ip6Address'} ) && $json->{'ip6Address'} ne "" )
		|| ( defined( $json->{'ip6Gateway'} ) && $json->{'ip6Gateway'} ne "" ) )
	{
		if ( defined($update_base) ) {
			if ( !defined( $json->{'ip6Address'} ) ) {
				$json->{'ip6Address'} = $update_base->{'ip6_address'};
			}
			if ( !defined( $json->{'ip6Gateway'} ) ) {
				$json->{'ip6Gateway'} = $update_base->{'ip6_gateway'};
			}
		}
		if ( !&in_same_net( $json->{'ip6Address'}, $json->{'ip6Gateway'} ) ) {
			return ( \%params, $json->{'ip6Address'} . " and " . $json->{'ip6Gateway'} . " are not in same network" );
		}
	}

	my $ipstr1;
	my $ipstr2;
	if (   ( defined( $json->{'ipAddress'} ) && $json->{'ipAddress'} ne "" )
		|| ( defined( $json->{'ipNetmask'} ) && $json->{'ipNetmask'} ne "" )
		|| ( defined( $json->{'ipGateway'} ) && $json->{'ipGateway'} ne "" ) )
	{
		if ( !defined( $json->{'ipAddress'} ) ) {
			return ( \%params, "ipAddress is not found" );
		}
		$ipstr1 = $json->{'ipAddress'} . "/" . $json->{'ipNetmask'};
		$ipstr2 = $json->{'ipGateway'} . "/" . $json->{'ipNetmask'};
		if ( defined( $json->{'ipNetmask'} ) && $json->{'ipNetmask'} ne "" && !&in_same_net( $ipstr1, $ipstr2 ) ) {
			return ( \%params, $json->{'ipAddress'} . " and " . $json->{'ipGateway'} . " are not in same network" );
		}
	}

	if ( ( defined( $json->{'iloIpAddress'} ) && $json->{'iloIpAddress'} ne "" )
		|| ( defined( $json->{'iloIpNetmask'} ) && $json->{'iloIpNetmask'} ne "" )
		|| ( defined( $json->{'iloIpGateway'} ) && $json->{'iloIpGateway'} ne "" ) )
	{
		if ( defined($update_base) ) {
			if ( !defined( $json->{'iloIpAddress'} ) ) {
				$json->{'iloIpAddress'} = $update_base->ilo_ip_address;
			}
			if ( !defined( $json->{'iloIpNetmask'} ) ) {
				$json->{'iloIpNetmask'} = $update_base->ilo_ip_netmask;
			}
			if ( !defined( $json->{'iloIpGateway'} ) ) {
				$json->{'iloIpGateway'} = $update_base->ilo_ip_gateway;
			}
		}
		$ipstr1 = $json->{'iloIpAddress'} . "/" . $json->{'iloIpNetmask'};
		$ipstr2 = $json->{'iloIpGateway'} . "/" . $json->{'iloIpNetmask'};
		if ( $json->{'iloIpGateway'} ne ""
			&& !&in_same_net( $ipstr1, $ipstr2 ) )
		{
			return ( \%params, $json->{'iloIpAddress'} . " and " . $json->{'iloIpGateway'} . " are not in same network" );
		}
	}

	if (   ( defined( $json->{'mgmtIpAddress'} ) && $json->{'mgmtIpAddress'} ne "" )
		|| ( defined( $json->{'mgmtIpNetmask'} ) && $json->{'mgmtIpNetmask'} ne "" )
		|| ( defined( $json->{'mgmtIpGateway'} ) && $json->{'mgmtIpGateway'} ne "" ) )
	{
		if ( defined($update_base) ) {
			if ( !defined( $json->{'mgmtIpAddress'} ) ) {
				$json->{'mgmtIpAddress'} = $update_base->mgmt_ip_address;
			}
			if ( !defined( $json->{'mgmtIpNetmask'} ) ) {
				$json->{'mgmtIpNetmask'} = $update_base->mgmt_ip_netmask;
			}
			if ( !defined( $json->{'mgmtIpGateway'} ) ) {
				$json->{'mgmtIpGateway'} = $update_base->mgmt_ip_gateway;
			}
		}
		$ipstr1 = $json->{'mgmtIpAddress'} . "/" . $json->{'mgmtIpNetmask'};
		$ipstr2 = $json->{'mgmtIpGateway'} . "/" . $json->{'mgmtIpNetmask'};
		if ( $json->{'mgmtIpGateway'} ne ""
			&& !&in_same_net( $ipstr1, $ipstr2 ) )
		{
			return ( \%params, $json->{'mgmtIpAddress'} . " and " . $json->{'mgmtIpGateway'} . " are not in same network" );
		}
	}

	if ( defined( $json->{'tcpPort'} ) && $json->{'tcpPort'} !~ /\d+/ ) {
		return ( \%params, $json->{'tcpPort'} . " is not a valid tcp port" );
	}
	if ( defined( $json->{'httpsPort'} ) && $json->{'httpsPort'} !~ /\d+/ ) {
		return ( \%params, $json->{'httpsPort'} . " is not a valid https port" );
	}

	return ( \%params, $err );
}

sub get_server_by_id {
	my $self = shift;
	my $id   = shift;
	my $row;
	my $isadmin = &is_admin($self);
	eval { $row = $self->db->resultset('Server')->find( { id => $id } ); };
	if ($@) { # $@ holds Perl errors
		$self->app->log->error("Failed to get server id = $id: $@");
		return ( undef, "Failed to get server id = $id: $@" );
	}
	my $data = {
		"id"             => $row->id,
		"hostName"       => $row->host_name,
		"domainName"     => $row->domain_name,
		"tcpPort"        => $row->tcp_port,
		"httpsPort"      => $row->https_port,
		"xmppId"         => $row->xmpp_id,
		"xmppPasswd"     => "**********",
		"interfaceName"  => $row->interface_name,
		"ipAddress"      => $row->ip_address,
		"ipNetmask"      => $row->ip_netmask,
		"ipGateway"      => $row->ip_gateway,
		"ip6Address"     => $row->ip6_address,
		"ip6Gateway"     => $row->ip6_gateway,
		"interfaceMtu"   => $row->interface_mtu,
		"cachegroup"     => $row->cachegroup->name,
		"cdn_id"         => $row->cdn_id,
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
		"lastUpdated"    => $row->last_updated,

	};
	return ( $data, undef );
}

sub create {
	my ( $params, $data, $err ) = ( undef, undef, undef );
	my $self = shift;

	my $json = $self->req->json;
	if ( !&is_oper($self) ) {
		return $self->forbidden("Forbidden. Insufficent permissions.");
	}

	( $params, $err ) = $self->check_server_params( $json, undef );
	if ( defined($err) ) {
		return $self->alert( $err );
	}

	my $new_id      = -1;
	my $xmpp_passwd = "BOOGER";
	my $insert;
	if ( defined( $json->{'ip6Address'} )
		&& $json->{'ip6Address'} ne "" )
	{
		eval {
			$insert = $self->db->resultset('Server')->create(
				{
					host_name        => $json->{'hostName'},
					domain_name      => $json->{'domainName'},
					tcp_port         => $params->{'tcpPort'},
					https_port         => $params->{'httpsPort'},
					xmpp_id          => $json->{'hostName'},                                                           # TODO JvD remove me later.
					xmpp_passwd      => $xmpp_passwd,
					interface_name   => $json->{'interfaceName'},
					ip_address       => $json->{'ipAddress'},
					ip_netmask       => $json->{'ipNetmask'},
					ip_gateway       => $json->{'ipGateway'},
					ip6_address      => $json->{'ip6Address'},
					ip6_gateway      => $json->{'ip6Gateway'},
					interface_mtu    => $json->{'interfaceMtu'},
					cachegroup       => $params->{'cachegroup'},
					cdn_id           => $params->{'cdnId'},
					phys_location    => $params->{'physLocation'},
					guid             => $json->{'guid'},
					rack             => $json->{'rack'},
					type             => $params->{'type'},
					status           => &admin_status_id( $self, $json->{'type'} eq "EDGE" ? "REPORTED" : "ONLINE" ),
					offline_reason   => $json->{'offline_reason'},
					profile          => $params->{'profile'},
					mgmt_ip_address  => $json->{'mgmtIpAddress'},
					mgmt_ip_netmask  => $json->{'mgmtIpNetmask'},
					mgmt_ip_gateway  => $json->{'mgmtIpGateway'},
					ilo_ip_address   => $json->{'iloIpAddress'},
					ilo_ip_netmask   => $json->{'iloIpNetmask'},
					ilo_ip_gateway   => $json->{'iloIpGateway'},
					ilo_username     => $json->{'iloUsername'},
					ilo_password     => $json->{'iloPassword'},
					router_host_name => $json->{'routerHostName'},
					router_port_name => $json->{'routerPortName'},
				}
			);
		};
		if ($@) { # $@ holds Perl errors
			$self->app->log->error("Failed to create server: $@");
			return $self->alert( { Error => "Failed to create server: $@" } );
		}
	}
	else {
		eval {
			$insert = $self->db->resultset('Server')->create(
				{
					host_name        => $json->{'hostName'},
					domain_name      => $json->{'domainName'},
					tcp_port         => $params->{'tcpPort'},
					https_port         => $params->{'httpsPort'},
					xmpp_id          => $json->{'hostName'},                                                           # TODO JvD remove me later.
					xmpp_passwd      => $xmpp_passwd,
					interface_name   => $json->{'interfaceName'},
					ip_address       => $json->{'ipAddress'},
					ip_netmask       => $json->{'ipNetmask'},
					ip_gateway       => $json->{'ipGateway'},
					interface_mtu    => $json->{'interfaceMtu'},
					cachegroup       => $params->{'cachegroup'},
					cdn_id           => $params->{'cdnId'},
					phys_location    => $params->{'physLocation'},
					guid             => $json->{'guid'},
					rack             => $json->{'rack'},
					type             => $params->{'type'},
					status           => &admin_status_id( $self, $json->{'type'} eq "EDGE" ? "REPORTED" : "ONLINE" ),
					offline_reason   => $json->{'offline_reason'},
					profile          => $params->{'profile'},
					mgmt_ip_address  => $json->{'mgmtIpAddress'},
					mgmt_ip_netmask  => $json->{'mgmtIpNetmask'},
					mgmt_ip_gateway  => $json->{'mgmtIpGateway'},
					ilo_ip_address   => $json->{'iloIpAddress'},
					ilo_ip_netmask   => $json->{'iloIpNetmask'},
					ilo_ip_gateway   => $json->{'iloIpGateway'},
					ilo_username     => $json->{'iloUsername'},
					ilo_password     => $json->{'iloPassword'},
					router_host_name => $json->{'routerHostName'},
					router_port_name => $json->{'routerPortName'},
				}
			);
		};
		if ($@) { # $@ holds Perl errors
			$self->app->log->error("Failed to create server: $@");
			return $self->alert( { Error => "Failed to create server: $@" } );
		}
	}
	$insert->insert();
	$new_id = $insert->id;
	if (   $json->{'type'} eq "EDGE"
		|| $json->{'type'} eq "MID" )
	{
		$insert = $self->db->resultset('Servercheck')->create( { server => $new_id, } );
		$insert->insert();
	}

	# if the insert has failed, we don't even get here, we go to the exception page.
	&log( $self, "Create server with hostname:" . $json->{'hostName'}, "APICHANGE" );

	( $data, $err ) = $self->get_server_by_id($new_id);
	if ( defined($err) ) {
		return $self->alert( { Error => $err } );
	}
	$self->success($data, "Server successfully created: " . $json->{'hostName'});
}

sub update {
	my ( $params, $data, $err ) = ( undef, undef, undef );
	my $self = shift;
	my $json = $self->req->json;
	if ( !&is_oper($self) ) {
		return $self->forbidden("Forbidden. Insufficent permissions.");
	}

	my $id = $self->param('id');

	# get resultset for original and one to be updated.  Use to examine diffs to propagate the effects of the change.
	my $org_server = $self->db->resultset('Server')->find( { id => $id } );
	if ( !defined($org_server) ) {
		return $self->not_found();
	}
	( $params, $err ) = $self->check_server_params( $json, $org_server );
	if ( defined($err) ) {
		return $self->alert( $err );
	}

	my $update = $self->db->resultset('Server')->find( { id => $id } );
	eval {
		$update->update(
			{
				host_name      => defined( $params->{'hostName'} )      ? $params->{'hostName'}      : $update->host_name,
				domain_name    => defined( $params->{'domainName'} )    ? $params->{'domainName'}    : $update->domain_name,
				tcp_port       => defined( $params->{'tcpPort'} )       ? $params->{'tcpPort'}       : $update->tcp_port,
				https_port     => defined( $params->{'httpsPort'} )       ? $params->{'httpsPort'}   : $update->https_port,
				interface_name => defined( $params->{'interfaceName'} ) ? $params->{'interfaceName'} : $update->interface_name,
				ip_address     => defined( $params->{'ipAddress'} )     ? $params->{'ipAddress'}     : $update->ip_address,
				ip_netmask     => defined( $params->{'ipNetmask'} )     ? $params->{'ipNetmask'}     : $update->ip_netmask,
				ip_gateway     => defined( $params->{'ipGateway'} )     ? $params->{'ipGateway'}     : $update->ip_gateway,
				ip6_address => defined( $params->{'ip6Address'} ) && $params->{'ip6Address'} != "" ? $params->{'ip6Address'} : $update->ip6_address,
				ip6_gateway      => defined( $params->{'ip6Gateway'} )      ? $params->{'ip6Gateway'}      : $update->ip6_gateway,
				interface_mtu    => defined( $params->{'interfaceMtu'} )    ? $params->{'interfaceMtu'}    : $update->interface_mtu,
				cachegroup       => defined( $params->{'cachegroup'} )       ? $params->{'cachegroup'}       : $update->cachegroup->id,
				cdn_id           => defined( $params->{'cdnId'} )           ? $params->{'cdnId'}           : $update->cdn_id,
				phys_location    => defined( $params->{'physLocation'} )    ? $params->{'physLocation'}    : $update->phys_location->id,
				guid             => defined( $params->{'guid'} )             ? $params->{'guid'}             : $update->guid,
				rack             => defined( $params->{'rack'} )             ? $params->{'rack'}             : $update->rack,
				type             => defined( $params->{'type'} )             ? $params->{'type'}             : $update->type->id,
				status           => defined( $params->{'status'} )           ? $params->{'status'}           : $update->status->id,
				offline_reason   => defined( $params->{'offline_reason'} )    ? $params->{'offline_reason'}    : $update->offline_reason,
				profile          => defined( $params->{'profile'} )          ? $params->{'profile'}          : $update->profile->id,
				mgmt_ip_address  => defined( $params->{'mgmtIpAddress'} )  ? $params->{'mgmtIpAddress'}  : $update->mgmt_ip_address,
				mgmt_ip_netmask  => defined( $params->{'mgmtIpNetmask'} )  ? $params->{'mgmtIpNetmask'}  : $update->mgmt_ip_netmask,
				mgmt_ip_gateway  => defined( $params->{'mgmtIpGateway'} )  ? $params->{'mgmtIpGateway'}  : $update->mgmt_ip_gateway,
				ilo_ip_address   => defined( $params->{'iloIpAddress'} )   ? $params->{'iloIpAddress'}   : $update->ilo_ip_address,
				ilo_ip_netmask   => defined( $params->{'iloIpNetmask'} )   ? $params->{'iloIpNetmask'}   : $update->ilo_ip_netmask,
				ilo_ip_gateway   => defined( $params->{'iloIpGateway'} )   ? $params->{'iloIpGateway'}   : $update->ilo_ip_gateway,
				ilo_username     => defined( $params->{'iloUsername'} )     ? $params->{'iloUsername'}     : $update->ilo_username,
				ilo_password     => defined( $params->{'iloPassword'} )     ? $params->{'iloPassword'}     : $update->ilo_password,
				router_host_name => defined( $params->{'routerHostName'} ) ? $params->{'routerHostName'} : $update->router_host_name,
				router_port_name => defined( $params->{'routerPortName'} ) ? $params->{'routerPortName'} : $update->router_port_name,
			}
		);
	};
	if ($@) { # $@ holds Perl errors
		$self->app->log->error("Failed to update server id = $id: $@");
		return $self->alert( { Error => "Failed to update server: $@" } );
	}
	$update->update();

	if ( $org_server->profile->id != $update->profile->id ) {
		my $param =
			$self->db->resultset('ProfileParameter')
			->search( { -and => [ profile => $org_server->profile->id, 'parameter.config_file' => 'rascal-config.txt', 'parameter.name' => 'CDN_name' ] },
			{ prefetch => [ { parameter => undef }, { profile => undef } ] } )->single();
		my $org_cdn_name = "";
		if ( defined($param) ) {
			$org_cdn_name = $param->parameter->value;
		}

		$param =
			$self->db->resultset('ProfileParameter')
			->search( { -and => [ profile => $update->profile->id, 'parameter.config_file' => 'rascal-config.txt', 'parameter.name' => 'CDN_name' ] },
			{ prefetch => [ { parameter => undef }, { profile => undef } ] } )->single();
		my $upd_cdn_name = "";
		if ( defined($param) ) {
			$upd_cdn_name = $param->parameter->value;
		}

		if ( $upd_cdn_name ne $org_cdn_name ) {
			my $delete = $self->db->resultset('DeliveryserviceServer')->search( { server => $id } );
			$delete->delete();
			&log( $self, $update->host_name . " profile change assigns server to new CDN - deleting all DS assignments", "APICHANGE" );
		}
		if ( $org_server->type->id != $update->type->id ) {
			my $delete = $self->db->resultset('DeliveryserviceServer')->search( { server => $id } );
			$delete->delete();
			&log( $self, $update->host_name . " profile change changes cache type - deleting all DS assignments", "APICHANGE" );
		}
	}

	if ( $org_server->type->id != $update->type->id ) {

		# server type changed:  servercheck entry required for EDGE and MID, but not others. Add or remove servercheck entry accordingly
		my %need_servercheck = map { &type_id( $self, $_ ) => 1 } qw{ EDGE MID };
		my $newtype_id       = $update->type->id;
		my $servercheck      = $self->db->resultset('Servercheck')->search( { server => $id } );
		if ( $servercheck != 0 && !$need_servercheck{$newtype_id} ) {

			# servercheck entry found but not needed -- delete it
			$servercheck->delete();
			&log( $self, $update->host_name . " cache type change - deleting servercheck", "APICHANGE" );
		}
		elsif ( $servercheck == 0 && $need_servercheck{$newtype_id} ) {

			# servercheck entry not found but needed -- insert it
			$servercheck = $self->db->resultset('Servercheck')->create( { server => $id } );
			$servercheck->insert();
			&log( $self, $update->host_name . " cache type changed - adding servercheck", "APICHANGE" );
		}
	}

	# this just creates the log string for the log table / tab.
	my $lstring = "Update server " . $update->host_name . " ";
	foreach my $col ( keys %{ $org_server->{_column_data} } ) {
        my $colParam = $col;
        $colParam =~ s/_(\w)/\U$1/g;
		if ( defined( $params->{$colParam} )
			&& $params->{$colParam} ne ( $org_server->{_column_data}->{$col} // "" ) )
		{
			if ( $col eq 'ilo_password' || $col eq 'xmpp_passwd' ) {
				$lstring .= $col . "-> ***********";
			}
			else {
				$lstring .= $col . "->" . $params->{$colParam} . " ";
			}
		}
	}

	# if the update has failed, we don't even get here, we go to the exception page.
	&log( $self, $lstring, "APICHANGE" );

	( $data, $err ) = $self->get_server_by_id($id);
	if ( defined($err) ) {
		return $self->alert( { Error => $err } );
	}
	$self->success($data, "Server was updated: " . $update->host_name);
}

sub delete {
	my ( $params, $data, $err ) = ( undef, undef, undef );
	my $self = shift;

	if ( !&is_oper($self) ) {
		return $self->forbidden("Forbidden. Insufficent permissions.");
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

	return $self->success_message("Server was deleted: " . $host_name);
}

sub postupdatequeue {
	my $self   = shift;
	my $params = $self->req->json;
	my $id     = $self->param('id');
	if ( !&is_oper($self) ) {
		return $self->forbidden("Forbidden. Insufficent permissions.");
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

1;
