package Extensions::Delegate::Statistics;
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

# Stub class that provides the "hook" for implementing custom statistics gathering.
# See TrafficOps.pm for the route that points here.
# TODO: drichardson - add documentation here and in the public doc (.rst).
sub new {
	my $self  = {};
	my $class = shift;
	return ( bless( $self, $class ) );
}

sub get_stats {
	return ( 1, "No Traffic Ops Extension implemented for 'Statistics->get_stats()'" );
}

sub get_usage_overview {
	return ( 1, "No Traffic Ops Extension implemented for 'Statistics->get_usage_overview()'" );
}

1;
