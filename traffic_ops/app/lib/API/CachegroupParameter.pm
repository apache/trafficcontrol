package API::CachegroupParameter;
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

sub index {
	my $self = shift;
	my @data;
	my $orderby = $self->param('orderby') || "cachegroup";
	my $rs_data = $self->db->resultset("CachegroupParameter")->search( undef, { prefetch => [ 'cachegroup', 'parameter' ], order_by => $orderby } );
	while ( my $row = $rs_data->next ) {
		push(
			@data, {
				"cachegroup"   => $row->cachegroup->name,
				"parameter"    => $row->parameter->id,
				"last_updated" => $row->last_updated,
			}
		);
	}
	$self->success( { "cachegroupParameters" => \@data } );
}

sub create {
	my $self 	= shift;
	my $params 	= $self->req->json;

	if ( !&is_oper($self) ) {
		return $self->forbidden();
	}

	if ( !defined($params) ) {
		return $self->alert("parameters must be in JSON format.");
	}
	if ( ref($params) ne 'ARRAY' ) {
		my @temparry;
		push(@temparry, $params);
		$params = \@temparry;
	}
	if ( scalar(@{ $params }) == 0 ) {
		return $self->alert("parameters array length is 0.");
	}

	$self->db->txn_begin();
	foreach my $param (@{ $params }) {
		my $cg = $self->db->resultset('Cachegroup')->find( { id => $param->{cacheGroupId} } );
		if ( !defined($cg) ) {
			$self->db->txn_rollback();
			return $self->alert("Cache Group with id: " . $param->{cacheGroupId} . " doesn't exist");
		}
		my $parameter = $self->db->resultset('Parameter')->find( { id => $param->{parameterId} } );
		if ( !defined($parameter) ) {
			$self->db->txn_rollback();
			return $self->alert("Parameter with id: " . $param->{parameterId} . " doesn't exist");
		}
		my $cg_param = $self->db->resultset('CachegroupParameter')->find( { parameter => $parameter->id, cachegroup => $cg->id } );
		if ( defined($cg_param) ) {
			$self->db->txn_rollback();
			return $self->alert("parameter: " . $param->{parameterId} . " already associated with cachegroup: " . $param->{cacheGroupId});
		}
		$self->db->resultset('CachegroupParameter')->create( { parameter => $parameter->id, cachegroup => $cg->id } )->insert();
	}
	$self->db->txn_commit();

	&log( $self, "New cache group parameter associations were created.", "APICHANGE" );

	my $response = $params;
	return $self->success($response, "Cachegroup parameter associations were created.");
}

sub delete {
	my $self 			= shift;
	my $cg_id 			= $self->param('cachegroup_id');
	my $parameter_id 	= $self->param('parameter_id');

	if ( !&is_oper($self) ) {
		return $self->forbidden();
	}

	my $cg = $self->db->resultset('Cachegroup')->find( { id => $cg_id } );
	if ( !defined($cg) ) {
		return $self->not_found();
	}

	my $parameter = $self->db->resultset('Parameter')->find( { id => $parameter_id } );
	if ( !defined($parameter) ) {
		return $self->not_found();
	}

	my $delete = $self->db->resultset('CachegroupParameter')->find( { parameter => $parameter->id, cachegroup => $cg->id } );
	if ( !defined($delete) ) {
		return $self->alert("parameter: $parameter_id isn't associated with cachegroup: $cg_id.");
	}

	$delete->delete();

	&log( $self, "Deleted cache group parameter " . $cg->name . " <-> " . $parameter->name, "APICHANGE" );

	return $self->success_message("Cachegroup parameter association was deleted.");
}



1;
