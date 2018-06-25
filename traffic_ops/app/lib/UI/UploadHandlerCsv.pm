package UI::UploadHandlerCsv;
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
use Data::Dumper;
use UI::Utils;
use UI::Server;
use Mojo::Base 'Mojolicious::Controller';
use Mojo::Base 'Mojolicious::Plugin';
use JSON;
use Mojo::JSON;
use Mojo::Upload;

use Mojo::Upload;

# use Mojo::Asset::File;
use Mojo::Log;
use Text::ParseWords;

sub getCdnCheckData {
	my $self = shift;
	my @data;
	my $cdns       = '';
	my $cdnHashRef = {};
	my $orderby        = "name";
	$orderby = $self->param('orderby') if ( defined $self->param('orderby') );
	my $rs_data = $self->db->resultset("Cdn")->search( undef, { order_by => $orderby } );
	while ( my $row = $rs_data->next ) {
		$cdns .= $row->name . ',';
		$cdnHashRef->{ $row->name } = $row->id;
		push(
			@data, {
				"id"           => $row->id,
				"name"         => $row->name,
			}
		);
	}
	return $cdnHashRef;
}

sub getCachegroupCheckData {    # renamed to 'CacheGroup'
	my $self = shift;
	my @data;
	my $cachegroups       = '';
	my $cachegroupHashRef = {};
	my %idnames;
	my $orderby = "name";

	#$orderby = $self->param('orderby') if ( defined $self->param('orderby') );
	# Can't figure out how to do the join on the same table
	my $rs_idnames = $self->db->resultset("Cachegroup")->search( undef, { columns => [qw/id name/] } );
	while ( my $row = $rs_idnames->next ) {
		$idnames{ $row->id } = $row->name;
	}
	my $rs_data = $self->db->resultset("Cachegroup")->search( undef, { prefetch => [ { 'type' => undef, } ], order_by => 'me.' . $orderby } );
	while ( my $row = $rs_data->next ) {
		$cachegroups .= $row->name . ',';
		$cachegroupHashRef->{ $row->name } = $row->id;
		if ( defined $row->parent_cachegroup_id ) {
			push(
				@data, {
					"id"                     => $row->id,
					"name"                   => $row->name,
					"short_name"             => $row->short_name,
					"last_updated"           => $row->last_updated,
					"parent_cachegroup_id"   => $row->parent_cachegroup_id,
					"parent_cachegroup_name" => $idnames{ $row->parent_cachegroup_id },
					"type_id"                => $row->type->id,
					"type_name"              => $row->type->name,
				}
			);
		}
		else {
			push(
				@data, {
					"id"                     => $row->id,
					"name"                   => $row->name,
					"short_name"             => $row->short_name,
					"last_updated"           => $row->last_updated,
					"parent_cachegroup_id"   => $row->parent_cachegroup_id,
					"parent_cachegroup_name" => undef,
					"type_id"                => $row->type->id,
					"type_name"              => $row->type->name,
				}
			);
		}
	}
	return $cachegroupHashRef;
}

sub getTypeCheckData {
	my $self = shift;
	my @data;
	my $types       = '';
	my $typeHashRef = {};
	my $orderby     = "name";
	$orderby = $self->param('orderby') if ( defined $self->param('orderby') );
	my $rs_data = $self->db->resultset("Type")->search( undef, { order_by => $orderby } );
	while ( my $row = $rs_data->next ) {
		$types .= $row->name . ',';
		$typeHashRef->{ $row->name } = $row->id;
		push(
			@data, {
				"id"           => $row->id,
				"name"         => $row->name,
				"description"  => $row->description,
				"use_in_table" => $row->use_in_table,
				"last_updated" => $row->last_updated,
			}
		);
	}
	return $typeHashRef;
}

sub getProfileCheckData {
	my $self = shift;
	my @data;
	my $profiles       = '';
	my $profileHashRef = {};
	my $orderby        = "name";
	$orderby = $self->param('orderby') if ( defined $self->param('orderby') );
	my $rs_data = $self->db->resultset("Profile")->search( undef, { order_by => $orderby } );
	while ( my $row = $rs_data->next ) {
		$profiles .= $row->name . ',';
		$profileHashRef->{ $row->name } = $row->id;
		push(
			@data, {
				"id"           => $row->id,
				"name"         => $row->name,
				"description"  => $row->description,
				"last_updated" => $row->last_updated,
			}
		);
	}
	return $profileHashRef;
}

sub getPhysLocationCheckData {
	my $self = shift;
	my @data;
	my $physLocations       = '';
	my $physLocationHashRef = {};
	my $rs_data             = $self->db->resultset("PhysLocation")->search( undef, { prefetch => ['region'], order_by => 'me.name' } );
	while ( my $row = $rs_data->next ) {
		next if $row->short_name eq 'UNDEF';
		$physLocations .= $row->name . ',';
		$physLocationHashRef->{ $row->name } = $row->id;
		push(
			@data, {
				"id"         => $row->id,
				"name"       => $row->name,
				"short_name" => $row->short_name,
				"address"    => $row->address,
				"city"       => $row->city,
				"state"      => $row->state,
				"zip"        => $row->zip,
				"poc"        => $row->poc,
				"phone"      => $row->phone,
				"email"      => $row->email,
				"comments"   => $row->comments,
				"region"     => $row->region->name,
			}
		);
	}
	return $physLocationHashRef;
}

sub getParamHashRef {
	my $p            = $_[0];
	my $lineNumber   = $_[1];
	my $paramHashRef = {};
	$paramHashRef->{'host_name'}        = $p->[0];
	$paramHashRef->{'domain_name'}      = $p->[1];
	$paramHashRef->{'interface_name'}   = $p->[2];
	$paramHashRef->{'ip_address'}       = $p->[3];
	$paramHashRef->{'ip_netmask'}       = $p->[4];
	$paramHashRef->{'ip_gateway'}       = $p->[5];
	$paramHashRef->{'ip6_address'}      = $p->[6];
	$paramHashRef->{'ip6_gateway'}      = $p->[7];
	$paramHashRef->{'interface_mtu'}    = $p->[8];
	$paramHashRef->{'cdn'}              = $p->[9];
	$paramHashRef->{'cachegroup'}       = $p->[10];
	$paramHashRef->{'phys_location'}    = $p->[11];
	$paramHashRef->{'rack'}             = $p->[12];
	$paramHashRef->{'type'}             = $p->[13];
	$paramHashRef->{'profile'}          = $p->[14];
	$paramHashRef->{'tcp_port'}         = $p->[15];
	$paramHashRef->{'mgmt_ip_address'}  = $p->[16];
	$paramHashRef->{'mgmt_ip_netmask'}  = $p->[17];
	$paramHashRef->{'mgmt_ip_gateway'}  = $p->[18];
	$paramHashRef->{'ilo_ip_address'}   = $p->[19];
	$paramHashRef->{'ilo_ip_netmask'}   = $p->[20];
	$paramHashRef->{'ilo_ip_gateway'}   = $p->[21];
	$paramHashRef->{'ilo_username'}     = $p->[22];
	$paramHashRef->{'ilo_password'}     = $p->[23];
	$paramHashRef->{'router_host_name'} = $p->[24];
	$paramHashRef->{'router_port_name'} = $p->[25];
	$paramHashRef->{'https_port'}       = $p->[26];
	$paramHashRef->{'offline_reason'}   = $p->[27];
	$paramHashRef->{'status'}           = '';
	$paramHashRef->{'csv_line_number'}  = $lineNumber;
	return $paramHashRef;
}

sub checkNamedValues {
	my $lineNumber          = shift;
	my $errorLineDelim      = shift;
	my $enteredCdn          = shift;
	my $cdnHashRef          = shift;
	my $enteredCachegroup   = shift;
	my $cachegroupHashRef   = shift;
	my $enteredType         = shift;
	my $typeHashRef         = shift;
	my $enteredProfile      = shift;
	my $profileHashRef      = shift;
	my $enteredPhysLocation = shift;
	my $physLocationHashRef = shift;
	my $processCSVErrors    = '';

	# allow integers for backward compatability but if non-integer then validate as well
	if ( !exists $cdnHashRef->{$enteredCdn} ) {
		$processCSVErrors
			.= $errorLineDelim
			. "[LINE #:"
			. $lineNumber
			. "]<span style='color:blue;'>CDN NOT VALID["
			. $enteredCdn
			. "] CASE SENSITIVE.</span>";
	}

	if ( !exists $cachegroupHashRef->{$enteredCachegroup} ) {
		$processCSVErrors
			.= $errorLineDelim
			. "[LINE #:"
			. $lineNumber
			. "]<span style='color:blue;'>CACHE_GROUP NOT VALID["
			. $enteredCachegroup
			. "] CASE SENSITIVE.</span>";
	}

	if ( !exists $physLocationHashRef->{$enteredPhysLocation} ) {
		$processCSVErrors
			.= $errorLineDelim
			. "[LINE #:"
			. $lineNumber
			. "]<span style='color:blue;'>PHYS LOCATION NOT VALID["
			. $enteredPhysLocation
			. "] CASE SENSITIVE.</span>";
	}

	if ( !exists $typeHashRef->{$enteredType} ) {
		$processCSVErrors
			.= $errorLineDelim . "[LINE #:" . $lineNumber . "]<span style='color:blue;'>TYPE NOT VALID[" . $enteredType . "] CASE SENSITIVE.</span>";
	}

	if ( !exists $profileHashRef->{$enteredProfile} && $enteredProfile !~ /\d+/) {
		$processCSVErrors
			.= $errorLineDelim . "[LINE #:" . $lineNumber . "]<span style='color:blue;'>PROFILE NOT VALID[" . $enteredProfile . "] CASE SENSITIVE.</span>";
	}

	print $processCSVErrors;
	return $processCSVErrors;
}

sub processCSV {
	my $self                = shift;
	my $fileNameAndPath     = shift;
	my $cdnHashRef          = shift;
	my $cachegroupHashRef   = shift;
	my $typeHashRef         = shift;
	my $profileHashRef      = shift;
	my $physLocationHashRef = shift;
	my $processCSVErrors    = '';
	my $lineNumber          = 1;
	my $errorLineDelim      = '</li><li>';
	open( INPUTFILE, "<$fileNameAndPath" );

	while (<INPUTFILE>) {
		my $line = $_;
		chomp($line);
		if ( $line !~ /\r\n/ ) {
			my $new = $line;
			$new =~ tr/\r\n/\n/d;    # transcode any carriage return characters <ctrl-m> to line feed characters
			my @processedLines = split( '\n', $new );
			foreach my $pLine (@processedLines) {
				if ( $pLine !~ /HEADER:/ && $pLine !~ /^host/ && $pLine !~ /^#/ ) {    # ignore the header and comments
					my $delim = ',';
					my $keep  = 0;

					#print Dumper($pLine);
					my @p = parse_line( $delim, $keep, $pLine );

					#print Dumper(@p);
					my $itemCount    = scalar @p;
					my $correctCount = 28;
					if ( $itemCount > $correctCount || $itemCount < $correctCount ) {
						$processCSVErrors
							.= $errorLineDelim
							. "[LINE #:"
							. $lineNumber
							. "] [ITEM COUNT = "
							. $itemCount . "/"
							. $correctCount
							. "] PLEASE FIX EACH LINE AND ENSURE AN ITEM COUNT = " . $correctCount . ".";
						$processCSVErrors
							.= "</li><ul><li style='color:blue;'>"
							. "[host,domain,int,ip4,subnet,gw,ip6,gw6,mtu,cdn,cachegroup,phys_loc,rack,type,prof,port,1g_ip,1g_subnet,1g_gw,ilo_ip,ilo_subnet,ilo_gw,ilo_user,ilo_pwd,r_host,r_port,https_port,offline_reason]"
							. "</li></ul>";
					}
					else {
						my $paramHashRef = &getParamHashRef( \@p, $lineNumber );
						$paramHashRef = &replaceNamedLookupValues( $paramHashRef, $cdnHashRef, $cachegroupHashRef, $typeHashRef, $profileHashRef, $physLocationHashRef );

						# print Dumper($paramHashRef);
						$processCSVErrors .= &UI::Server::check_server_input( $self, $paramHashRef );
						my $enteredCdn          = $p[9];
						my $enteredCachegroup   = $p[10];
						my $enteredPhysLocation = $p[11];
						my $enteredType         = $p[13];
						my $enteredProfile      = $p[14];
						$processCSVErrors .= &checkNamedValues(
							$lineNumber, $errorLineDelim, $enteredCdn, $cdnHashRef, $enteredCachegroup, $cachegroupHashRef, $enteredType,
							$typeHashRef, $enteredProfile, $profileHashRef, $enteredPhysLocation, $physLocationHashRef
						);
					}
				}
				$lineNumber++;
			}
		}
		else {
			if ( $line !~ /HEADER:/ && $line !~ /^host/ && $line !~ /^#/ ) {    # ignore the header and comments
				my $delim        = ',';
				my $keep         = 0;
				my @p            = parse_line( $delim, $keep, $line );
				my $itemCount    = scalar @p;
				my $correctCount = 28;
				if ( $itemCount > $correctCount || $itemCount < $correctCount ) {
					$processCSVErrors
						.= $errorLineDelim
						. "[LINE #:"
						. $lineNumber
						. "] [ITEM COUNT = "
						. $itemCount . "/"
						. $correctCount
						. "] PLEASE FIX EACH LINE AND ENSURE AN ITEM COUNT = " . $correctCount . ".";
					$processCSVErrors
						.= "</li><ul><li style='color:blue;'>"
						. "[host,domain,int,ip4,subnet,gw,ip6,gw6,mtu,cdn,cachegroup,phys_loc,rack,type,prof,port,1g_ip,1g_subnet,1g_gw,ilo_ip,ilo_subnet,ilo_gw,ilo_user,ilo_pwd,r_host,r_port,https_port,offline_reason]"
						. "</li></ul>";
				}
				else {
					my $paramHashRef = &getParamHashRef( \@p, $lineNumber );
					$paramHashRef = &replaceNamedLookupValues( $paramHashRef, $cdnHashRef, $cachegroupHashRef, $typeHashRef, $profileHashRef, $physLocationHashRef );
					$processCSVErrors .= &UI::Server::check_server_input( $self, $paramHashRef );
					my $enteredCdn          = $p[9];
					my $enteredCachegroup   = $p[10];
					my $enteredPhysLocation = $p[11];
					my $enteredType         = $p[13];
					my $enteredProfile      = $p[14];
					$processCSVErrors .= &checkNamedValues(
						$lineNumber, $errorLineDelim, $enteredCdn, $cdnHashRef, $enteredCachegroup, $cachegroupHashRef, $enteredType,
						$typeHashRef, $enteredProfile, $profileHashRef, $enteredPhysLocation, $physLocationHashRef
					);
				}
			}
			$lineNumber++;
		}
	}
	close(INPUTFILE);
	print $processCSVErrors;
	return $processCSVErrors;
}

sub replaceNamedLookupValues {
	my $paramHashRef        = shift;
	my $cdnHashRef          = shift;
	my $cachegroupHashRef   = shift;
	my $typeHashRef         = shift;
	my $profileHashRef      = shift;
	my $physLocationHashRef = shift;

	$paramHashRef->{'cdn'} = $cdnHashRef->{ $paramHashRef->{'cdn'} };

	#  if ($paramHashRef->{'cachegroup'} !~ /^[+-]?\d+$/) {  # if not an integer
	$paramHashRef->{'cachegroup'} = $cachegroupHashRef->{ $paramHashRef->{'cachegroup'} };

	#  }
	#  if ($paramHashRef->{'type'} !~ /^[+-]?\d+$/) {
	$paramHashRef->{'type'} = $typeHashRef->{ $paramHashRef->{'type'} };

	#  }
	if ($paramHashRef->{'profile'} !~ /^[+-]?\d+$/) {
		$paramHashRef->{'profile'} = $profileHashRef->{ $paramHashRef->{'profile'} };

	}
	#  if ($paramHashRef->{'phys_location'} !~ /^[+-]?\d+$/) {
	$paramHashRef->{'phys_location'} = $physLocationHashRef->{ $paramHashRef->{'phys_location'} };

	#  }
	return $paramHashRef;
}

sub processSynchronizeCSV {
	my $self                = shift;
	my $fileNameAndPath     = shift;
	my $cdnHashRef          = shift;
	my $cachegroupHashRef   = shift;
	my $typeHashRef         = shift;
	my $profileHashRef      = shift;
	my $physLocationHashRef = shift;
	my $processCSVErrors    = '';
	my $lineNumber          = 1;
	my $errorLineDelim      = '</li><li>';
	open( INPUTFILE, "<$fileNameAndPath" );

	while (<INPUTFILE>) {
		my $line = $_;
		chomp($line);
		if ( $line !~ /\r\n/ ) {
			my $new = $line;
			$new =~ tr/\r\n/\n/d;    # transcode any carriage return characters <ctrl-m> to line feed characters
			my @processedLines = split( '\n', $new );
			foreach my $pLine (@processedLines) {
				if ( $pLine !~ /HEADER:/ && $pLine !~ /^host/ && $pLine !~ /^#/ ) {    # ignore the header and comments
					my $delim        = ',';
					my $keep         = 0;
					my @p            = parse_line( $delim, $keep, $pLine );
					my $itemCount    = scalar @p;
					my $correctCount = 28;
					if ( $itemCount > $correctCount || $itemCount < $correctCount ) {
						$processCSVErrors
							.= $errorLineDelim
							. "[LINE #:"
							. $lineNumber
							. "] [ITEM COUNT = "
							. $itemCount . "/"
							. $correctCount
							. "] PLEASE FIX EACH LINE AND ENSURE AN ITEM COUNT = " . $correctCount . ".";
						$processCSVErrors
							.= "</li><ul><li style='color:blue;'>"
							. "[host,domain,int,ip4,subnet,gw,ip6,gw6,mtu,cdn,cachegroup,phys_loc,rack,type,prof,port,1g_ip,1g_subnet,1g_gw,ilo_ip,ilo_subnet,ilo_gw,ilo_user,ilo_pwd,r_host,r_port,https_port,offline_reason]"
							. "</li></ul>";
					}
					else {
						my $paramHashRef = &getParamHashRef( \@p, $lineNumber );
						$paramHashRef = &replaceNamedLookupValues( $paramHashRef, $cdnHashRef, $cachegroupHashRef, $typeHashRef, $profileHashRef, $physLocationHashRef );

						# insert/create new record
						eval { &UI::Server::create( $self, $paramHashRef ) };
						if ($@) {
							my $exceptionError = $@;
							$exceptionError =~ s/\://g;

							#$exceptionError =~ s/\(//g;
							#$exceptionError =~ s/\)//g;
							#$exceptionError =~ s/\[//g;
							#$exceptionError =~ s/\]//g;
							$exceptionError =~ s/\"//g;
							$exceptionError =~ s/\'//g;

							#$exceptionError =~ s/\.//g;
							#$exceptionError =~ s/\?//g;
							#$exceptionError =~ s/\,//g;
							#$exceptionError =~ s/\=//g;
							$exceptionError =~ s/\\//g;

							#$exceptionError =~ s/\///g;
							$exceptionError =~ s{\n}{ }g;
							$processCSVErrors .= '<li style=\"display:none;\">[LINE #' . $lineNumber . '] [EXCEPTION_ERROR] - ' . $exceptionError . '</li>';
							$processCSVErrors =~ s/\[LINE/\<\/li\>\<li\>\[LINE/g;
						}
					}
				}
				$lineNumber++;
			}
		}
		else {
			if ( $line !~ /HEADER:/ && $line !~ /^host/ && $line !~ /^#/ ) {    # ignore the header and comments
				my $delim        = ',';
				my $keep         = 0;
				my @p            = parse_line( $delim, $keep, $line );
				my $itemCount    = scalar @p;
				my $correctCount = 28;
				if ( $itemCount > $correctCount || $itemCount < $correctCount ) {
					$processCSVErrors
						.= $errorLineDelim
						. "[LINE #:"
						. $lineNumber
						. "] [ITEM COUNT = "
						. $itemCount . "/"
						. $correctCount
						. "] PLEASE FIX EACH LINE AND ENSURE AN ITEM COUNT = " . $correctCount . ".";
					$processCSVErrors
						.= "</li><ul><li style='color:blue;'>"
						. "[host,domain,int,ip4,subnet,gw,ip6,gw6,mtu,cdn,cachegroup,phys_loc,rack,type,prof,port,1g_ip,1g_subnet,1g_gw,ilo_ip,ilo_subnet,ilo_gw,ilo_user,ilo_pwd,r_host,r_port,https_port,offline_reason]"
						. "</li></ul>";
				}
				else {
					my $paramHashRef = &getParamHashRef( \@p, $lineNumber );
					$paramHashRef = &replaceNamedLookupValues( $paramHashRef, $cdnHashRef, $cachegroupHashRef, $typeHashRef, $profileHashRef, $physLocationHashRef );

					# insert/create new record
					eval { &UI::Server::createserver( $self, $paramHashRef ) };
					if ($@) {
						my $exceptionError = $@;
						$exceptionError =~ s/\://g;

						#$exceptionError =~ s/\(//g;
						#$exceptionError =~ s/\)//g;
						#$exceptionError =~ s/\[//g;
						#$exceptionError =~ s/\]//g;
						$exceptionError =~ s/\"//g;
						$exceptionError =~ s/\'//g;

						#$exceptionError =~ s/\.//g;
						#$exceptionError =~ s/\?//g;
						#$exceptionError =~ s/\,//g;
						#$exceptionError =~ s/\=//g;
						$exceptionError =~ s/\\//g;

						#$exceptionError =~ s/\///g;
						$exceptionError =~ s{\n}{ }g;
						$processCSVErrors .= '<li style=\"display:none;\">[LINE #' . $lineNumber . '] [EXCEPTION_ERROR] - ' . $exceptionError . '</li>';
						$processCSVErrors =~ s/\[LINE/\<\/li\>\<li\>\[LINE/g;
					}
				}
			}
			$lineNumber++;
		}
	}
	close(INPUTFILE);
	return $processCSVErrors;
}

sub upload {
	my $self = shift;

	my $serverPath = '/tmp/';
	my $url        = $self->req->url->to_abs;
	my $userinfo   = $self->req->url->to_abs->userinfo;
	my $host       = $self->req->url->to_abs->host;

	my $upload = $self->param('file-0');

	return $self->render_exception( status => 400, '[.csv files only] Upload size exceeded the 1GByte(1073741824 bytes) limit!' )
		if $upload->size > 1073741824;    # 1Gbyte limit
	                                      #if $upload->size > 5242880;  # 5Mbyte limit
	my $fileNameAndPath = $serverPath . $upload->filename;
	$upload->move_to($fileNameAndPath);

	my $cdnHashRef          = &getCdnCheckData($self);
	my $cachegroupHashRef   = &getCachegroupCheckData($self);
	my $typeHashRef         = &getTypeCheckData($self);
	my $profileHashRef      = &getProfileCheckData($self);
	my $physLocationHashRef = &getPhysLocationCheckData($self);
	my $processCSVErrors    = &processCSV( $self, $fileNameAndPath, $cdnHashRef, $cachegroupHashRef, $typeHashRef, $profileHashRef, $physLocationHashRef );

	if ( length($processCSVErrors) <= 0 ) {
		$processCSVErrors = &processSynchronizeCSV( $self, $fileNameAndPath, $cdnHashRef, $cachegroupHashRef, $typeHashRef, $profileHashRef, $physLocationHashRef );
	}
	return $self->render( json => "{\"success\":true,\"serverpath\":\""
			. $serverPath
			. "\",\"filename\":\""
			. $upload->filename
			. "\",\"size\":"
			. $upload->size
			. ",\"processcsverrors\":\""
			. $processCSVErrors
			. "\"}" );

}

1;
