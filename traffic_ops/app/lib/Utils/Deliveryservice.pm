package Utils::DeliveryService;
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

use Mojo::UserAgent;
use utf8;
use Carp qw(cluck confess);
use Data::Dumper;

sub new {
	my $self  = {};
	my $class = shift;
	my $args  = shift;

	return ( bless( $self, $class ) );
}

sub build_match {
	my $self            = shift;
	my $cdn_name        = shift;
	my $ds_name         = shift;
	my $cachegroup_name = shift;
	my $peak_usage_type = shift;
	return $cdn_name . ":" . $ds_name . ":" . $cachegroup_name . ":all:" . $peak_usage_type;
}
1;
