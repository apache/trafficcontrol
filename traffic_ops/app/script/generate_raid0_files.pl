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

my $outputdir = "/var/www/files/disks";
my $fileprefix = "CreateRaid0-";
my $stripesize = 256; # 128=64k, 256=128k, etc
my $minbaynum = 0;
my $maxbaynum = 23;

for (my $i = $minbaynum; $i <= $maxbaynum ; $i++) {
	my $filename;
	if ($i < 10) {
		 $filename = "${outputdir}/${fileprefix}0${i}.xml";
	} else {
		 $filename = "${outputdir}/${fileprefix}${i}.xml";
	}
	open (OUT, "> $filename") || die("Could not open $filename for writing");
	print OUT <<EOF;
<p:CreateVirtualDisk_INPUT xmlns:p="http://schemas.dmtf.org/wbem/wscim/1/cim-schema/2/root/dcim/DCIM_RAIDService">
	<p:Target>RAID.Integrated.1-1</p:Target>
	<p:PDArray>Disk.Bay.$i:Enclosure.Internal.0-1:RAID.Integrated.1-1</p:PDArray>
	<p:VDPropNameArray>RAIDLevel</p:VDPropNameArray>
	<p:VDPropNameArray>SpanLength</p:VDPropNameArray>
	<p:VDPropNameArray>SpanDepth</p:VDPropNameArray>
	<p:VDPropNameArray>StripeSize</p:VDPropNameArray>
	<p:VDPropNameArray>VirtualDiskName</p:VDPropNameArray>
	<p:VDPropNameArray>ReadCachePolicy</p:VDPropNameArray>
	<p:VDPropNameArray>WriteCachePolicy</p:VDPropNameArray>
	<p:VDPropValueArray>2</p:VDPropValueArray>
	<p:VDPropValueArray>1</p:VDPropValueArray>
	<p:VDPropValueArray>1</p:VDPropValueArray>
	<p:VDPropValueArray>$stripesize</p:VDPropValueArray>
	<p:VDPropValueArray>Cachedisk_$i</p:VDPropValueArray>
	<p:VDPropValueArray>16</p:VDPropValueArray>
	<p:VDPropValueArray>2</p:VDPropValueArray>
</p:CreateVirtualDisk_INPUT>
EOF
	close OUT;
}
