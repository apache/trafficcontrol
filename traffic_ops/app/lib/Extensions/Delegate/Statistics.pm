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

sub new {
	my $self  = {};
	my $class = shift;
	my $args  = shift;

	return ( bless( $self, $class ) );
}

sub long_term {
	return ( 1, "No Traffic Ops Extension implemented for 'Statistics->long_term()'" );
}

sub short_term_redis {
	return ( 1, "No Traffic Ops Extension implemented for 'Statistics->short_term()'" );
}

sub short_term {
	return ( 1, "No Traffic Ops Extension implemented for 'Statistics->short_term()'" );
}
1;
