#!/usr/bin/perl

use LWP::Simple;
use JSON qw( decode_json );
use Data::Dumper;
use strict;
use warnings;

my $requrl = 'http://localhost:8080/request';

my $json = get( $requrl );
die "Could not get $requrl!" unless defined $json;
my $tables = decode_json( $json );
print Dumper($tables);
foreach my $table ( @{ $tables } ) {
	print $table . " => ";
	my $url = 'http://localhost:8080/api/' . $table . '?format=moosefixture&join=no';
	my $t = get($url);
	my $filename;
	if ($t =~ /package Fixtures::Integration::(\S+);/) {
		$filename = $1 . '.pm';
	}
	print $filename;
	open(FH, ">$filename");
	print FH $t;
	close(FH);
	print "\n";
}


