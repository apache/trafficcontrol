package UI::DataAll;
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

sub availableds {
	my $self = shift;
	my @data;
	my $id = $self->param('id');
	my %dsids;
	my %takendsids;

	my $rs_takendsids = undef;
	$rs_takendsids = $self->db->resultset("DeliveryserviceTmuser")->search( { 'tm_user_id' => $id } );

	while ( my $row = $rs_takendsids->next ) {
		$takendsids{ $row->deliveryservice->id } = undef;
	}

	my $rs_links = $self->db->resultset("Deliveryservice")->search( undef, { order_by => "xml_id" } );
	while ( my $row = $rs_links->next ) {
		if ( !exists( $takendsids{ $row->id } ) ) {
			push( @data, { "id" => $row->id, "xml_id" => $row->xml_id } );
		}
	}

	$self->render( json => \@data );
}

# deprecated @see API/DeliveryServiceServer#index
sub data_links {
	my $self = shift;
	my @data;
	my $orderby = "deliveryservice";
	$orderby = $self->param('orderby') if ( defined $self->param('orderby') );
	my $rs_data = $self->db->resultset("DeliveryserviceServer")->search( undef, { order_by => $orderby } );
	while ( my $row = $rs_data->next ) {
		push(
			@data, {
				"deliveryservice" => $row->deliveryservice->id,
				"server"          => $row->server->id,
				"last_updated"    => $row->last_updated,
			}
		);
	}
	$self->render( json => \@data );
}

sub data_server {
	my $self = shift;
	my @data;
	my $orderby = "host_name";
	$orderby = $self->param('orderby') if ( defined $self->param('orderby') );
	my $rs_data = $self->db->resultset("Server")->search( undef, { order_by => $orderby } );
	while ( my $row = $rs_data->next ) {
		push(
			@data, {
				"id"               => $row->id,
				"host_name"        => $row->host_name,
				"domain_name"      => $row->domain_name,
				"tcp_port"         => $row->tcp_port,
				"xmpp_id"          => $row->xmpp_id,
				"xmpp_passwd"      => $row->xmpp_passwd,
				"interface_name"   => $row->interface_name,
				"ip_address"       => $row->ip_address,
				"ip_netmask"       => $row->ip_netmask,
				"ip_gateway"       => $row->ip_gateway,
				"interface_mtu"    => $row->interface_mtu,
				"location"         => $row->location->id,
				"type"             => $row->type->id,
				"status"           => $row->status->id,
				"profile"          => $row->profile->id,
				"ilo_ip_address"   => $row->ilo_ip_address,
				"ilo_ip_netmask"   => $row->ilo_ip_netmask,
				"ilo_ip_gateway"   => $row->ilo_ip_gateway,
				"ilo_username"     => $row->ilo_username,
				"ilo_password"     => $row->ilo_password,
				"router_host_name" => $row->router_host_name,
				"router_port_name" => $row->router_port_name,
			}
		);
	}
	$self->render( json => \@data );
}

# deprecated @see API/ProfileParameters#domains
# TODO JvD - this is the 3rd copy of the exact same function!
sub data_domains {
	my $self = shift;
	my @data;

	my $rs = $self->db->resultset('Profile')->search( { 'me.name' => { -like => 'CCR%' } }, { prefetch => ['cdn'] } );
	while ( my $row = $rs->next ) {
		push(
			@data, {
				"domainName"         => $row->cdn->domain_name,
				"parameterId"        => -1,  # it's not a parameter anymore
				"profileId"          => $row->id,
				"profileName"        => $row->name,
				"profileDescription" => $row->description,
			}
		);

	}
	$self->render( json => \@data );
}
1;
