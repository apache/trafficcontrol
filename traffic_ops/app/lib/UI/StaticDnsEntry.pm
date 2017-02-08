package UI::StaticDnsEntry;
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

sub edit {
	my $self = shift;
	my $mode = $self->param('mode');
	my $dsid = $self->param('id');
	my $ds;
	&stash_role($self);
	my $static_dns;
	my $i = 0;
	my $rs = $self->db->resultset('Staticdnsentry')->search( { deliveryservice => $dsid }, { prefetch => [ 'cachegroup', 'deliveryservice' ] } );
	while ( my $row = $rs->next ) {
		$static_dns->[$i] = $row;
		$i++;
		$ds->{xml_id} = $row->deliveryservice->xml_id;
	}
	$ds->{id} = $dsid;

	$self->stash( ds          => $ds );
	$self->stash( static_dns  => $static_dns );
	$self->stash( fbox_layout => 1 );
}

sub check_update_input {
	my $self = shift;
	my $err  = "";

	if ( !&is_oper($self) ) {
		$err .= "You do not have enough privileges to modify this.";
		return $err;
	}

	foreach my $param ( $self->param ) {
		if ( $self->param($param) eq "" ) {
			$err .= $param . " cannot be empty.";
			last;
		}
		my $atype_id     = &type_id( $self, "A_RECORD" );
		my $aaaatype_id  = &type_id( $self, "AAAA_RECORD" );
		my $cnametype_id = &type_id( $self, "CNAME_RECORD" );
		if ( $param =~ /^address_(.*\d+)/ ) {
			if ( $self->param( 'type_' . $1 ) == $atype_id && !&is_ipaddress( $self->param($param) ) ) {
				$err .= $self->param($param) . " is not a valid IPv4 address.";
				last;
			}
			elsif ( $self->param( 'type_' . $1 ) == $aaaatype_id && !&is_ip6address( $self->param($param) ) ) {
				$err .= $self->param($param) . " is not a valid IPv6 address.";
				last;
			}
			elsif ( $self->param( 'type_' . $1 ) == $cnametype_id && $self->param($param) !~ /.*\.$/ ) {
				$err .= $self->param($param) . " is not a valid cname.";
				last;
			}
		}
		if ( $param =~ /^ttl_/ && $self->param($param) =~ m/[a-zA-Z]/ ) {
			$err .= $self->param($param) . " is not a valid ttl (NaN).";
		}
	}
	return $err;
}

sub update_assignments {
	my $self = shift;
	my $dsid = $self->param('dsid');

	my $err = &check_update_input($self);
	if ( defined($err) && $err ne "" ) {
		$self->flash( alertmsg => $err );
		my $referer = $self->req->headers->header('referer');
		if ( defined($referer) ) {
			return $self->redirect_to($referer);
		}
		else {
			return $self->render( text => "ERR = " . $err, layout => undef );    # for testing - $referer is not defined there.
		}
	}

	# my @active_entries = ();
	foreach my $param ( $self->param ) {
		if ( $param =~ /host_new_(\d+)/ ) {
			my $host       = $self->param($param);
			my $ttl        = $self->param( 'ttl_new_' . $1 );
			my $type       = $self->param( 'type_new_' . $1 );
			my $address    = $self->param( 'address_new_' . $1 );
			my $cachegroup = $self->param( 'cg_new_' . $1 );
			my $insert_dns = $self->db->resultset('Staticdnsentry')->create(
				{
					host            => $host,
					address         => $address,
					type            => $type,
					ttl             => $ttl,
					deliveryservice => $dsid,
					cachegroup      => $cachegroup,
				}
			);
			$insert_dns->insert();
			my $new_id = $insert_dns->id;
			&log( $self, "Create static dns entry " . $host . "->" . $address . " for DS " . $dsid, "UICHANGE" );
		}
		elsif ( $param =~ /host_(\d+)/ ) {
			my $sdns_id = $1;

			my $host       = $self->param($param);
			my $ttl        = $self->param( 'ttl_' . $sdns_id );
			my $type       = $self->param( 'type_' . $sdns_id );
			my $address    = $self->param( 'address_' . $sdns_id );
			my $cachegroup = $self->param( 'cg_' . $1 );
			my $update     = $self->db->resultset('Staticdnsentry')->find( { id => $sdns_id } );
			my %hash       = (
				host            => $host,
				ttl             => $ttl,
				type            => $type,
				address         => $address,
				deliveryservice => $dsid,
				cachegroup      => $cachegroup,
			);
			$update->update( \%hash );
			&log( $self, "Update static dns entry " . $host . "->" . $address . " for DS " . $dsid, "UICHANGE" );
		}
	}
	my $referer = $self->req->headers->header('referer');
	return $self->redirect_to($referer);
}

sub delete {
	my $self = shift;
	my $id   = $self->param('id');

	if ( !&is_oper($self) ) {
		$self->flash( alertmsg => "No can do. Get more privs." );
	}
	else {
		my $deleted = $self->db->resultset('Staticdnsentry')->search( { 'me.id' => $id }, { prefetch => [ 'cachegroup', 'deliveryservice' ] } )->single();
		my $delete = $self->db->resultset('Staticdnsentry')->search( { id => $id } );
		$delete->delete();
		&log( $self, "Delete static dns entry " . $deleted->host . "->" . $deleted->address . " from " . $deleted->deliveryservice->xml_id, "UICHANGE" );
	}
	return $self->redirect_to('/close_fancybox.html');
}

sub read {
	my $self = shift;

	my @data;
	my $orderby = "deliveryservice";
	$orderby = $self->param('orderby') if ( defined $self->param('orderby') );
	my $rs_data = $self->db->resultset("Staticdnsentry")->search( undef, { prefetch => [ 'deliveryservice', 'type', 'cachegroup' ], order_by => $orderby } );
	while ( my $row = $rs_data->next ) {
		push(
			@data, {
				"deliveryservice" => $row->deliveryservice->xml_id,
				"host"            => $row->host,
				"ttl"             => $row->ttl,
				"address"         => $row->address,
				"type"            => $row->type->name,
				"cachegroup"      => $row->cachegroup->name,
			}
		);
	}
	$self->render( json => \@data );
}

# Create not needed - only created through update_assignments

1;
