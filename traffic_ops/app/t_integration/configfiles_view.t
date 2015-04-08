package main;
use Mojo::Base -strict;
use Test::More;
use Test::Mojo;
use Data::Dumper;
use DBI;
use strict;
use warnings;

BEGIN { $ENV{MOJO_MODE} = "integration" }
my $t = Test::Mojo->new('TrafficOps');
no warnings 'once';
use warnings 'all';
my $api_version = '1.1';

$t->post_ok( '/login', => form => { u => 'admin', p => 'password' } )->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );

my $json = JSON->new->allow_nonref;
$t->get_ok('/dataprofile')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );
my $profile_arr  = $json->decode( $t->tx->res->content->asset->{content} );
my %profile_name = ();
my $i            = 0;
foreach my $p ( @{$profile_arr} ) {
	$profile_name{ $p->{id} } = $p->{name};
	$i++;
}

$t->get_ok( '/api/' . $api_version . '/servers.json' )->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content} } );
my $servers = $json->decode( $t->tx->res->content->asset->{content} );

$i = 0;
foreach my $server ( @{ $servers->{response} } ) {
	$t->get_ok( '/ort/' . $server->{hostName} . '/ort1' )->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );
	my $files = $json->decode( $t->tx->res->content->asset->{content} );
	if ( $server->{type} eq 'EDGE' || $server->{type} eq 'MID' ) {
		diag "... server: " . $server->{hostName} . ' ' . $files->{other}->{CDN_name} . ' ' . $files->{profile}->{name};
		my $path = '/genfiles/view/' . $server->{hostName} . '/all';
		$t->get_ok($path)->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );
		foreach my $file ( keys %{ $files->{config_files} } ) {
			$t->content_like(qr/$file/);

			# diag $file;
		}
	}
}

done_testing();
