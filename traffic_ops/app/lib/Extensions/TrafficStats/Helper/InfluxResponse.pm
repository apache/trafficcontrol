package Extensions::TrafficStats::Helper::InfluxResponse;
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
use UI::Utils;
use constant FIVE_MINUTES => 5;

use Data::Dumper;
use JSON;
use POSIX qw(strftime);
use POSIX qw(localtime);
use HTTP::Date;
use Common::ReturnCodes qw(SUCCESS ERROR);

my $args = shift;

sub new {
	my $self  = {};
	my $class = shift;
	$args = shift;

	return ( bless( $self, $class ) );
}

sub parse_retention_period_in_seconds {
	my $self             = shift;
	my $retention_period = shift;

	local $/;

	my ( $hour, $minutes, $seconds ) = $retention_period =~ /(\d*)h(\d*)m(\d*)s/ms;

	my $hour_in_seconds    = $hour * 60 * 60;
	my $minutes_in_seconds = $minutes * 60;
	my $total_seconds      = $hour_in_seconds + $minutes_in_seconds + $seconds;

	return ($total_seconds);
}

1;
