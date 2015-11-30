package WebDep;

#
# Copyright 2015 Comcast Cable Communications Management, LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

use strict;
use warnings;
use WWW::Curl::Easy;
use Digest::MD5 qw{md5 md5_hex md5_base64};
use File::Path;
use File::Basename qw{dirname};
use IO::Uncompress::Unzip qw{$UnzipError};
use Data::Dumper;
use Cwd qw{getcwd abs_path};

# columns from web_deps.txt in order
my @columns = qw{name version cdn_location compression filename source_dir final_dir};

# class method returns list of WebDep objects loaded from file
sub getDeps {
	my $class = shift;

	my $webdeps_file = shift;
	my @deps;
	open my $fh, '<', $webdeps_file or die "Can't open $webdeps_file\n";

	# final_dir within webdeps_file are absolute or relative to directory that file is in
	my $oldcwd      = getcwd();
	my $webdeps_dir = dirname($webdeps_file);
	chdir($webdeps_dir);
	while (<$fh>) {

		# comments only if line starts with #
		next if /^#/;
		chomp;
		my @parts = split( /\s*,\s*/, $_ );
		my $obj = bless {}, $class;
		for my $attrib (@columns) {
			$obj->{$attrib} = shift @parts;
		}
		$obj->{filename} =~ s{^/}{};
		$obj->{source_dir} =~ s{/$}{};
		$obj->{final_dir} = abs_path( $obj->{final_dir} );
		push( @deps, $obj );
	}
	chdir($oldcwd);
	return @deps;
}

sub getSrcFileName {
	my $self = shift;
	return join( '/', $self->{source_dir}, $self->{filename} );
}

sub getDestFileName {
	my $self = shift;
	return join( '/', $self->{final_dir}, $self->{filename} );
}

sub getCdnLoc {
	my $self = shift;
	my ( $response_body, $err ) = _getCdnLoc( $self->{cdn_location} );
	return ( $response_body, $err );
}

sub getContent {
	my $self = shift;
	my ( $cdnloc, $err ) = $self->getCdnLoc();
	if ( defined $err ) {
		die "$err\n";
	}
	my $srcfn = $self->getSrcFileName();
	if ( !exists $self->{content} ) {
		if ( $self->{compression} eq 'zip' ) {

			#my $u = IO::Uncompress::Unzip->new( $cdnloc, Name => $srcfn ) or die "IO::Uncompress::Unzip failed: $UnzipError\n";
			my $u = new IO::Uncompress::Unzip($cdnloc) or die "IO::Uncompress::Unzip failed: $UnzipError\n";
			my $buffer;
			my $content = "";
			while ( $u->read($buffer) > 0 ) {
				$content .= $buffer;
			}
			$self->{content} = $content;
		}
		else {
			# no compression
			$self->{content} = $cdnloc;
		}
	}
	return $self->{content};
}

sub needsUpdating {
	my $self = shift;
	my $fn   = $self->getDestFileName();

	if ( -f $fn ) {

		# checksum and compare
		open my $fh, '<', $fn or die "Can't open existing file: $fn\n";
		my $md5_existing = md5_hex(<$fh>);
		close $fh;

		my $md5_new = md5_hex( $self->getContent() );
		return ( $md5_new ne $md5_existing );
	}
	return 1;
}

sub update {
	my $self = shift;
	my $err;
	my $srcfn  = $self->getSrcFileName();
	my $destfn = $self->getDestFileName();
	my $action = "";

	# download archive
	if ( $self->needsUpdating() ) {
		if ( -f $destfn ) {
			$action = "Replaced";
		}
		else {
			$action = "Created";
			my $final_dir = $self->{final_dir};

			if ( !-d $final_dir ) {
				print "Making dir: $final_dir\n";
				mkpath($final_dir);
			}
			open my $ofn, '>', $destfn or die "Can't write to $destfn\n";
			print $ofn, $self->getContent();
			close $ofn;
		}
	}
	else {
		$action = "Kept";
	}
	return ( $action, $err );
}

####################################################################
# Utilities
sub execCommand {
	my ( $command, @args ) = @_;
	my $pid    = fork();
	my $result = 0;

	if ( $pid == 0 ) {
		exec( $command, @args );
		exit 0;
	}
	else {
		wait;
		$result = $?;
		if ( $result != 0 ) {
			print "ERROR executing: $command,  args: " . join( ' ', @args ) . "\n";
		}
	}
	return $result;
}

sub curlMe {
	my $url  = shift;
	my $curl = WWW::Curl::Easy->new;
	my $response_body;
	my $err;    # undef if no error

	$curl->setopt( CURLOPT_VERBOSE, 0 );
	if ( $url =~ m/https/ ) {
		$curl->setopt( CURLOPT_SSL_VERIFYHOST, 0 );
		$curl->setopt( CURLOPT_SSL_VERIFYPEER, 0 );
	}
	$curl->setopt( CURLOPT_IPRESOLVE,      1 );
	$curl->setopt( CURLOPT_FOLLOWLOCATION, 1 );
	$curl->setopt( CURLOPT_CONNECTTIMEOUT, 5 );
	$curl->setopt( CURLOPT_TIMEOUT,        15 );
	$curl->setopt( CURLOPT_HEADER,         0 );
	$curl->setopt( CURLOPT_URL,            $url );
	$curl->setopt( CURLOPT_WRITEDATA,      \$response_body );
	my $retcode       = $curl->perform;
	my $response_code = $curl->getinfo(CURLINFO_HTTP_CODE);

	if ( $response_code != 200 ) {
		$err = "Got HTTP $response_code response for '$url'";
	}
	elsif ( length($response_body) == 0 ) {
		$err = "URL: $url returned empty!!";
	}
	return ( $response_body, $err );
}

{
	# cache cdn locs -- some files extracted from same downloaded archive
	my %cdn_loc_for;

	sub _getCdnLoc {
		my $cdnloc = shift;
		my $err;
		if ( !exists $cdn_loc_for{$cdnloc} ) {
			my $response_body;
			( $response_body, $err ) = curlMe($cdnloc);
			$cdn_loc_for{$cdnloc} = $response_body;
		}

		# could be undef indicating previous error
		return ( $cdn_loc_for{$cdnloc}, $err );
	}
}

1;
