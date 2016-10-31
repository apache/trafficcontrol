#!/bin/bash

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


# This script looks at a configuration file (disk.cfg) and at the type
# of machine it's running on to build a drive configuration, which it writes in
# /tmp/disk_config
#
# ks.cfg will then pick it up with an %include /tmp/harddrive

hardware=`dmidecode -s system-product-name | grep -v '^#' | sed -e 's/\r//g' | sed -e 's/\n//g'`

# This is where Traffic Ops GenIso.pm writes disk.cfg.
# Well, more specifically it's where kickstart mounts 
# what GenIso writes. 
config_location='/mnt/stage2/ks_scripts/disk.cfg'

# This is the file that ks.cfg looks for (e.g what this script writes):
disk_config_loc='/tmp/drive_config'

boot_drives=''

if [ -e $config_location ]
then
    source $config_location
elif [ -e 'disk.cfg' ]
then
    source 'disk.cfg' # This is for testing outside of kickstart envs. 
fi

boot_drives_old=$boot_drives
case $hardware in
    SSG-6047R-*)
    if [ "$boot_drives" == '' ]
    then 
        boot_drives='sdb,sdc'
    fi 
    first_drive=`echo $boot_drives | awk -F, ' { print $1 } '`
    second_drive=`echo $boot_drives | awk -F, ' { print $2 } '`

    cat <<EOF >> $disk_config_loc
# $hardware and $boot_drives_old
# Disk config
clearpart --all --initlabel --drives=$boot_drives
bootloader --location=mbr --driveorder=$boot_drives --append="crashkernel=auto"

ignoredisk --only-use=$boot_drives

part raid.boot.b --size=500 --ondisk=$first_drive
part raid.boot.c --size=500 --ondisk=$second_drive
part swap --size=2048 --ondisk=$first_drive
part swap --size=2048 --ondisk=$second_drive
part raid.root.b --size=1 --grow --ondisk=$first_drive
part raid.root.c --size=1 --grow --ondisk=$second_drive
raid /boot --fstype=ext4 --device=md0 --level=1 raid.boot.b raid.boot.c
raid / --fstype=ext4 --device=md1 --level=1 raid.root.b raid.root.c
EOF
;;
    *)
    if [ "$boot_drives" == "" ]
    then
        boot_drives='sda'
    fi
    cat <<EOF1 >> $disk_config_loc
# $hardware and $boot_drives_old
# Disk config
bootloader --location=mbr
clearpart --all --initlabel
zerombr yes
ignoredisk --only-use=$boot_drives
part /boot --fstype=ext4 --label=boot --size=500
part swap --size=4000
part / --fstype=ext4 --label=root --size=1 --grow
EOF1
;;
esac
