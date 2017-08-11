package API::Asn;
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

use Mojo::Base 'Mojolicious::Controller';
use Data::Dumper;

# Index
sub index {
	my $self    = shift;
	my $cg_id   = $self->param('cachegroup');

	my %criteria;
	if ( defined $cg_id ) {
		$criteria{'cachegroup'} = $cg_id;
	}

	my @data;
	my $orderby = $self->param('orderby') || "asn";
	my $rs_data = $self->db->resultset("Asn")->search( \%criteria, { prefetch => [ { 'cachegroup' => undef } ], order_by => "me." . $orderby } );
	while ( my $row = $rs_data->next ) {
		push(
			@data, {
				"id"           => $row->id,
				"asn"          => $row->asn,
				"cachegroupId" => $row->cachegroup->id,
				"cachegroup"   => $row->cachegroup->name,
				"lastUpdated"  => $row->last_updated
			}
		);
	}
	$self->success( \@data );
}

sub index_v11 {
	my $self = shift;
	my @data;
	my $orderby = $self->param('orderby') || "asn";
	my $rs_data = $self->db->resultset("Asn")->search( undef, { prefetch => [ { 'cachegroup' => undef } ], order_by => "me." . $orderby } );
	while ( my $row = $rs_data->next ) {
		push(
			@data, {
				"id"          => $row->id,
				"asn"         => $row->asn,
				"cachegroup"  => $row->cachegroup->name,
				"lastUpdated" => $row->last_updated,
			}
		);
	}
	$self->success( { "asns" => \@data } );
}

# Show
sub show {
	my $self = shift;
	my $id   = $self->param('id');

	my $rs_data = $self->db->resultset("Asn")->search( { 'me.id' => $id }, { prefetch => ['cachegroup'] } );
	my @data = ();
	while ( my $row = $rs_data->next ) {
		push(
			@data, {
				"id"           => $row->id,
				"asn"          => $row->asn,
				"cachegroupId" => $row->cachegroup->id,
				"cachegroup"   => $row->cachegroup->name,
				"lastUpdated"  => $row->last_updated
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

	my $asn = $self->db->resultset('Asn')->find( { id => $id } );
	if ( !defined($asn) ) {
		return $self->not_found();
	}

	if ( !defined($params) ) {
		return $self->alert("parameters must be in JSON format.");
	}

	if ( !defined( $params->{asn} ) ) {
		return $self->alert("ASN is required.");
	}

	if ( !defined( $params->{cachegroupId} ) ) {
		return $self->alert("Cachegroup is required.");
	}

	my $values = {
		asn        => $params->{asn},
		cachegroup => $params->{cachegroupId}
	};

	my $rs = $asn->update($values);
	if ( $rs ) {
		my $response;
		$response->{id}           = $rs->id;
		$response->{asn}          = $rs->asn;
		$response->{cachegroupId} = $rs->cachegroup->id;
		$response->{cachegroup}   = $rs->cachegroup->name;
		$response->{lastUpdated}  = $rs->last_updated;
		&log( $self, "Updated ASN name '" . $rs->asn . "' for id: " . $rs->id, "APICHANGE" );
		return $self->success( $response, "ASN update was successful." );
	}
	else {
		return $self->alert("ASN update failed.");
	}

}

sub create {
	my $self   = shift;
	my $params = $self->req->json;

	if ( !&is_oper($self) ) {
		return $self->forbidden();
	}

	if ( !defined($params->{asn}) ) {
		return $self->alert("ASN is required.");
	}

	if ( !defined($params->{cachegroupId}) ) {
		return $self->alert("Cachegroup Id is required.");
	}

	my $values = {
		asn 		=> $params->{asn} ,
		cachegroup 	=> $params->{cachegroupId}
	};

	my $insert = $self->db->resultset('Asn')->create($values);
	my $rs = $insert->insert();
	if ($rs) {
		my $response;
		$response->{id}          	= $rs->id;
		$response->{asn}        	= $rs->asn;
		$response->{cachegroupId}   = $rs->cachegroup->id;
		$response->{cachegroup}   	= $rs->cachegroup->name;
		$response->{lastUpdated} 	= $rs->last_updated;

		&log( $self, "Created ASN name '" . $rs->asn . "' for id: " . $rs->id, "APICHANGE" );

		return $self->success( $response, "ASN create was successful." );
	}
	else {
		return $self->alert("ASN create failed.");
	}

}

sub delete {
	my $self = shift;
	my $id     = $self->param('id');

	if ( !&is_oper($self) ) {
		return $self->forbidden();
	}

	my $asn = $self->db->resultset('Asn')->find( { id => $id } );
	if ( !defined($asn) ) {
		return $self->not_found();
	}

	my $rs = $asn->delete();
	if ($rs) {
		return $self->success_message("ASN deleted.");
	} else {
		return $self->alert( "ASN delete failed." );
	}
}
1;
