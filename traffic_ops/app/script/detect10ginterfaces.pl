#!/usr/bin/perl
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
use strict;
use File::Basename;
my @ifaces = </sys/class/net/*>;
for my $iface(@ifaces) {
	open(IN, "< $iface/speed");
	my $line = <IN>;
	close IN;
	if ($line =~ m/^10000$/) {
		my $name = basename($iface);
		open (MAC, "< $iface/address");
		my $address = <MAC>;
		close MAC;
		chomp $address;
		chomp $name;
		open (OUT, "> /etc/sysconfig/network-scripts/ifcfg-$name");
		print OUT <<EOF;
DEVICE="$name"
HWADDR="$address"
ONBOOT="no"
SLAVE="yes"
MASTER="bond0"
EOF
		close OUT;
	}
}
