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
use warnings;

$| = 1;

my $udevadm = "/sbin/udevadm";
my $kernel_device = shift(@ARGV);

if (!defined($kernel_device)) {
	die("Please specify a device");
}

if (! -x $udevadm) {
	die("$udevadm: $!");
}

# udevadm info --export-db |grep block |grep expander |egrep '^P'
# ...
# P: /devices/pci0000:80/0000:80:01.0/0000:81:00.0/host2/port-2:0/expander-2:0/port-2:0:7/end_device-2:0:7/target2:0:7/2:0:7:0/block/sdbd
# P: /devices/pci0000:80/0000:80:01.0/0000:81:00.0/host2/port-2:0/expander-2:0/port-2:0:8/end_device-2:0:8/target2:0:8/2:0:8:0/block/sdbe
# P: /devices/pci0000:80/0000:80:01.0/0000:81:00.0/host2/port-2:0/expander-2:0/port-2:0:9/end_device-2:0:9/target2:0:9/2:0:9:0/block/sdbf
# ...

my @out = `/sbin/udevadm info --export-db`;
my $devices = {};

for my $line (@out) {
	chomp($line);
	next unless ($line =~ /^P: .*$/ && $line =~ /block/ && $line =~ /expander/);

	my (undef, $device) = split(/^P: /, $line, 2);
	my @path_parts = split(/\//, $device);
	my $disk = pop(@path_parts);
	pop(@path_parts);
	my $path_ids = pop(@path_parts);
	my @target_parts = split(/:/, $path_ids);

	$devices->{$target_parts[0]}->{$target_parts[1]}->{$target_parts[2]}->{$target_parts[3]} = $disk;
}

my $matched = 0;
my $oa = 0;

for my $da (sort { $a <=> $b } (keys(%{$devices}))) {
	my $ob = 0;
	for my $db (sort { $a <=> $b } (keys(%{$devices->{$da}}))) {
		for my $dc (sort { $a <=> $b } (keys(%{$devices->{$da}->{$db}}))) {
			for my $dd (sort { $a <=> $b } (keys(%{$devices->{$da}->{$db}->{$dc}}))) {
				if ($kernel_device eq $devices->{$da}->{$db}->{$dc}->{$dd}) {
					printf("%d-%d-%02d\n", $oa, $ob, $dc);
					$matched = 1;
					last;
				}
			}

			last if ($matched);
		}

		last if ($matched);
		$ob++;
	}

	last if ($matched);
	$oa++;
}
