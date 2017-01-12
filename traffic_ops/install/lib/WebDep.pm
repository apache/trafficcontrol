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
        $obj->{source_dir} =~ s{/$}{};

        my $final_dir = $obj->{final_dir};
        if ( !-d $final_dir ) {
            print "Making dir: $final_dir\n";
            mkpath($final_dir);
        }
        $obj->{final_dir} = abs_path($final_dir);
        push( @deps, $obj );
    }
    chdir($oldcwd);
    return @deps;
}

sub getSrcFileName {
    my $self = shift;
    my @parts;
    push @parts, $self->{source_dir} if $self->{source_dir} ne '';
    push @parts, $self->{filename};
    return join( '/', @parts );
}

sub getDestFileName {
    my $self = shift;
    return join( '/', $self->{final_dir}, $self->{filename} );
}

sub getDownloadContent {
    my $self = shift;
    my ( $response_body, $err ) = _getDownloadContent( $self->{cdn_location} );
    return ( $response_body, $err );
}

sub getContent {
    my $self = shift;
    my ( $content, $err ) = $self->getDownloadContent();
    if ( defined $err ) {
        die "$err\n";
    }
    my $srcfn = $self->getSrcFileName();
    if ( !exists $self->{content} ) {
        if ( $self->{compression} eq 'zip' ) {

            #print "Unzipping $srcfn\n";
            my $u = IO::Uncompress::Unzip->new( \$content ) or die "IO::Uncompress::Unzip failed: $UnzipError\n";
            my $found;
            while ( $u->nextStream() > 0 && !$u->eof() ) {
                my $name = $u->getHeaderInfo()->{Name};
                if ( $name eq $srcfn ) {
                    $found = $name;
                    last;
                }
            }
            if ( !defined $found ) {
                die "$srcfn not found in " . $self->{cdn_location} . "\n";
            }

            undef $/;    # slurp mode
            $content = <$u>;
            $u->close();
        }

        $self->{content} = $content;
    }
    return $self->{content};
}

sub needsUpdating {
    my $self = shift;
    my $fn   = $self->getDestFileName();

    # checksum and compare
    open my $fh, '<', $fn or die "Can't open existing file: $fn\n";
    my $md5_existing = md5_hex(<$fh>);
    close $fh;

    my $md5_new       = md5_hex( $self->getContent() );
    my $needsUpdating = ( $md5_new ne $md5_existing );
    return $needsUpdating;
}

sub update {
    my $self = shift;
    my $err;
    my $srcfn  = $self->getSrcFileName();
    my $destfn = $self->getDestFileName();
    my $action = "";

    # download archive
    if ( -f $destfn ) {
        if ( !$self->needsUpdating() ) {
            $action = "Kept";
            return ( $action, $err );
        }

        # exists but needs to be replaced
        $action = "Replaced";
    }
    else {
        $action = "Created";
    }
    my $content = $self->getContent();
    open my $ofh, '>', $destfn or $err = "Can't write to $destfn";
    if ( !defined $err ) {
        print $ofh $content;
        close $ofh;
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
    my %content_for;

    sub _getDownloadContent {
        my $cdnloc = shift;
        my $err;
        if ( !exists $content_for{$cdnloc} ) {
            my $response_body;
            ( $response_body, $err ) = curlMe($cdnloc);
            $content_for{$cdnloc} = $response_body;
        }

        # could be undef indicating previous error
        return ( $content_for{$cdnloc}, $err );
    }
}

1;
