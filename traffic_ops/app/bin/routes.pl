#!/usr/bin/env perl
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
use Mojo::Base -strict;
use Schema;
use Test::IntegrationTestHelper;
use strict;
use warnings;
use Data::Dumper;

my $t = Test::Mojo->new('TrafficOps');

#print "t: " . Dumper( $t->ua->server->app->routes->children->[0]->{pattern} );
foreach my $i ( $t->ua->server->app->routes->children ) {
	foreach my $j (@$i) {
		my $method     = $j->{via}->[0];                          #GET/POST
		my $path       = $j->{pattern}{pattern};                  #/url
		my $package    = $j->{pattern}{defaults}{namespace};      # UI/API
		my $format     = $j->{pattern}{constraints}{format}[0];
		my $controller = $j->{pattern}{defaults}{controller};
		my $action     = $j->{pattern}{defaults}{action};
		if ( defined($package) && defined($method) && defined($action) && defined($path) && defined($controller) ) {

			#print "$method\t$path \t\t\t\t{:action =>$action, :package =>$package, :controller=>$controller} \n";
			my $max_length = 80;
			my $method_and_path = sprintf( "%-6s %s", $method, $path );
			if ( defined($format) ) {
				$method_and_path = $method_and_path . "." . $format;
			}

			my $method_and_path_length  = length($method_and_path);
			my $spacing                 = ' ' x ( $max_length - $method_and_path_length );
			my $fully_qualified_package = $package . "::" . $controller . "->" . $action;
			my $line                    = sprintf( "%s %s %s\n", $method_and_path, $spacing, $fully_qualified_package );
			print($line);

			#printf( "%s\n", '-' x length($line) );

			#printf( "%-5s %-40s {:action => %s, :package=> %s, :controller=> %s}\n", $method, $path, $action, $package, $controller );
		}

		#print "j: " . Dumper( $j->{pattern}{pattern} );
	}
}

#print "t: " . Dumper( $t->ua->server->app->routes->children->[0]->pattern );
