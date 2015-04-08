package main;
use Mojo::Base -strict;
use Test::More;
use Test::Mojo;
use Schema;
use Test::IntegrationTestHelper;
use strict;
use warnings;

BEGIN { $ENV{MOJO_MODE} = "integration" }
my $schema = Schema->connect_to_database;
my $t      = Test::Mojo->new('TrafficOps');


# Initialize the Kabletown CDN database. This will take a while
# diag "\n\n\n ***** PLEASE BE CAREFUL!!!! THIS TEST WILL BLOW AWAY YOUR DATABASE! DO YOU WANT TO CONTINUE?? (y/n):";
# my $ans = <STDIN>;
# chomp($ans);
# if ( $ans ne "y" ) {
# 	done_testing();
# 	exit();
# }
Test::IntegrationTestHelper->unload_core_data($schema);
Test::IntegrationTestHelper->load_core_data($schema);

$t->post_ok( '/login', => form => { u => 'admin', p => 'password' } )->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );

done_testing();
