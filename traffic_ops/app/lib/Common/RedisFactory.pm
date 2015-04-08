package Common::RedisFactory;
#
# Copyright 2015 Comcast Cable Communications Management, LLC
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
use Carp qw(cluck confess);
use Data::Dumper;
my $singleton = undef;
my $redis;

sub new {
	my $class = shift;

	# redis_connection_string of undef will short circuit for Unit Testing
	my ( $c, $redis_connection_string ) = @_;
	my $self = bless { c => $c, redis_connection_string => $redis_connection_string }, $class;

	my $deps        = $c->app->get_dependencies();
	my $redis_class = $deps->{'redis_class'};

	#print "RedisFactory.redis_class #-> (" . $redis_class . ")\n";
	unless ( eval "use $redis_class;" ) {
		$redis = $redis_class->new(
			server       => $redis_connection_string,
			read_timeout => 5,
			debug        => 0,
			reconnect    => 1
		);
	}
	else {
		die( "Could not load redis class '" . $redis_class . "'" );
	}

	return $singleton if defined $singleton;

	$singleton = bless \$self, $class;

	#print "New RedisFactory\n";
	return $singleton;
}

sub connection {
	my $self = shift;
	return $redis;
}

1;
