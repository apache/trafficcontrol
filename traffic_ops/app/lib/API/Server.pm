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

sub index {
	my $self   = shift;
	my $ds_id  = $self->param('dsId');
	my $helper = new Utils::Helper( { mojo => $self } );
	if ( defined $ds_id ) {
		if ( !$helper->is_valid_delivery_service($ds_id) ) {
			return $self->alert("Delivery Service does not exist.");
		}
		if ( $self->is_delivery_service_assigned($ds_id) || &is_admin($self) || &is_oper($self) ) {
			my $data = getserverdata( $self, $ds_id );
			$self->success($data);
		}
		else {
			$self->forbidden();
		}
	}
	else {
		if ( &is_admin($self) || &is_oper($self) ) {
			my $data = getserverdata($self);
			$self->success($data);
		}
		else {
			$self->forbidden();
		}
	}
}

sub getserverdata {
	my $self  = shift;
	my $ds_id = shift;
	my @data;
	my $isadmin = &is_admin($self);
	my $orderby = $self->param('orderby') || "host_name";
	my $servers;
	if ( defined $ds_id ) {

		# we want the edge cache servers and mid cache servers (but only mids if the delivery service uses mids)
		my @deliveryservice_servers_edge = $self->db->resultset('DeliveryserviceServer')->search(
			{
				deliveryservice => $ds_id,
			}
		)->get_column('server')->all();

		my $ds = $self->db->resultset('Deliveryservice')->search( { 'me.id' => $ds_id }, { prefetch => ['type'] } )->single();
		my @criteria     = [ { 'me.id' => { -in => \@deliveryservice_servers_edge } } ];
		my $subsel       = '(SELECT id FROM type where name = "MID")';
		my @types_no_mid = qw( HTTP_NO_CACHE HTTP_LIVE DNS_LIVE );                         # currently these are the ds types that bypass the mids
		if ( !grep { $_ eq $ds->type->name } @types_no_mid ) {
			push( @criteria, { 'me.type' => { -in => \$subsel }, 'me.cdn_id' => $ds->cdn_id } );
		}

		$servers = $self->db->resultset('Server')->search(
			[@criteria], {
				prefetch => [ 'cdn', 'cachegroup', 'type', 'profile', 'status', 'phys_location' ],
				order_by => 'me.' . $orderby,
			}
		);
	}
	else {
		# get all servers
		$servers = $self->db->resultset('Server')->search(
			undef, {
				prefetch => [ 'cdn', 'cachegroup', 'type', 'profile', 'status', 'phys_location' ],
				order_by => 'me.' . $orderby,
			}
		);
	}

	while ( my $row = $servers->next ) {
		my $cdn_name = defined( $row->cdn_id ) ? $row->cdn->name : "";

		push(
			@data, {
				"id"             => $row->id,
				"hostName"       => $row->host_name,
				"domainName"     => $row->domain_name,
				"tcpPort"        => $row->tcp_port,
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
				"physLocation"   => $row->phys_location->name,
				"rack"           => $row->rack,
				"type"           => $row->type->name,
				"status"         => $row->status->name,
				"profile"        => $row->profile->name,
				"cdnName"        => $cdn_name,
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

			}
		);
	}
	return ( \@data );
}

sub summary {
	my $self = shift;

	# TODO: drichardson - loop through this select to make this more dynamic.
	# Based on this: SELECT * FROM TYPE WHERE ID IN (SELECT TYPE FROM SERVER);
	my $edges   = $self->get_count_by_type('EDGE');
	my $mids    = $self->get_count_by_type('MID');
	my $rascals = $self->get_count_by_type('RASCAL');
	my $ccrs    = $self->get_count_by_type('CCR');
	my $redis   = $self->get_count_by_type('REDIS');

	my $response_body = [
		{ type => 'CCR',    count => $ccrs },
		{ type => 'EDGE',   count => $edges },
		{ type => 'MID',    count => $mids },
		{ type => 'REDIS',  count => $redis },
		{ type => 'RASCAL', count => $rascals }
	];
	return $self->success($response_body);
}

sub get_count_by_type {
	my $self      = shift;
	my $type_name = shift;
	return $self->db->resultset('Server')->search( { 'type.name' => $type_name }, { join => 'type' } )->count();
}

sub details {
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
			"xmppId"         => $row->xmpp_id,
			"xmppPasswd"     => $row->xmpp_passwd,
			"interfaceName"  => $row->interface_name,
			"ipAddress"      => $row->ip_address,
			"ipNetmask"      => $row->ip_netmask,
			"ipGateway"      => $row->ip_gateway,
			"ip6Address"     => $row->ip6_address,
			"ip6Gateway"     => $row->ip6_gateway,
			"interfaceMtu"   => $row->interface_mtu,
			"cachegroup"     => $row->cachegroup->name,
			"physLocation"   => $row->phys_location->name,
			"rack"           => $row->rack,
			"type"           => $row->type->name,
			"status"         => $row->status->name,
			"profile"        => $row->profile->name,
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

1;
