package Helper::DeliveryServiceStats;
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

use utf8;
use Data::Dumper;
use JSON;
use File::Slurp;
use Math::Round;

our @ISA = ("Helper::Stats");

sub series_name {
	my $self            = shift;
	my $cdn_name        = shift;
	my $ds_name         = shift;
	my $cachegroup_name = shift;
	my $metric_type     = shift;

	# 'series' section
	my $delim = ":";

	# Example: <cdn_name>:<deliveryservice_name>:<cache_group_name>:<metric_type>
	return sprintf( "%s$delim%s$delim%s$delim%s", $cdn_name, $ds_name, $cachegroup_name, $metric_type );
}

1;
