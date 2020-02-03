#!/usr/bin/python
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


'''
    This reads a configuration file and checks for functioning
    network links in /sys/class/net/*, then emits a ks.cfg network line.
    '''

from __future__ import print_function
import os
import re


global TO_LOG
# This "logs" to stdout which is captured during kickstart
TO_LOG = True

# This is the standard interface we install to. It is set to a speed value of
# 5 (vice 100,1000 or 10000) later on and any other interface will override it
# if you've got something faster installed.
standard_interface=['p4p1']

## These are configuration settings:
# Name of Configuration file:
cfg_file = "network.cfg"

# Where linux is putting the interface stuff:
iface_dir = '/sys/class/net/'

ignore_interfaces = ['lo','bond']

# Where we kickstart mounts the ISO, and our config directory:
base_cfg_dir = '/mnt/stage2/ks_scripts/'

# Remember the ? makes the match non-greedy. This is important.
cfg_line = re.compile("\s*(?P<key>.*?)=(?P<value>.*)\s*$")

# Pick the interface speed for bonding, or "Auto".
# Auto assumes you want the fastest connections with more than 1 interface,
# Or if there's not 2 interfaces at the same speed you want the fastest.
# auto is expected to be a string, otherwise use integers:
# Speed is in megs. 1000 is 1 gig, 10000 is 10G.

iface_speed = 'auto'

restring = iface_dir + "(?P<iface>.*)/speed"
iface_search = re.compile(restring)


def read_config(config_file):
    ''' Reads our network config file and hands back a dict of key:value
    pairs '''

    net_cfg = {}
    with open(config_file,'r') as cfg:
        network_lines = cfg.readlines()
        for line in network_lines:
            if cfg_line.match(line):
                key = cfg_line.match(line).group('key')
                value = cfg_line.match(line).group('value')
                net_cfg[key] = value
    return net_cfg


def find_usable_net_devs(location):
    ''' Search through iface_dir looking for /speed files.
    Build a dict keyed on speed (in otherwords the speed is the key with a list of
    interfaces as teh value). '''
    # We "pre-seed" the dictionary with the standard interface names at a
    # speed of 5 so that if there's nothing else we set that up. This
    # makes it easier  to reconfigure later.
    ifaces = {5:standard_interface}
    bad_ifaces={}
    devs = os.listdir(location)
    for dev in devs:
        dev_path = os.path.join(location,dev,'speed')
        add=True
        if os.path.isfile(dev_path):
            with open(dev_path,'r') as iface:
                try:
                    speed = iface.readlines()
                    # speed should only have one line:
                    speed = int(speed[0])
                # if there is no link some drivers/cards/whatever will
                # throw an IOError when you try to read the speed file.
                except IOError:
                    speed = 0
        # Other cards will return a -1, which is fine, but *some* of them
        # return a 65535. Those we set to 0 as well.
        if speed == 65535:
            speed = 0
        for i_face in ignore_interfaces:
            if i_face in dev:
                add = False
        if speed  <= 0:
            add = False
        if TO_LOG:
            print(add, dev)
        if add:
            if speed in ifaces:
                this_speed = ifaces[speed]
                this_speed.append(dev)
                ifaces[speed]=this_speed
            else:
                ifaces[speed]=[dev]
        else:
            bad_ifaces[dev] = speed
    print("We find these interfaces have link and might be useful:", ifaces)
    if TO_LOG:
        print("And these aren't useful:", bad_ifaces)
    return ifaces


def useable_interfaces(net_devs, nc, iface_speed):
    ''' This takes a go at figuring out which interfaces to use.'''
    iface_list = False
    notes = False
    if TO_LOG:
        print("in usable interfaces")

    if "bond" not in nc['BOND_DEVICE'].lower():
        if TO_LOG:
            print("useable interfaces if not", nc['BOND_DEVICE'])
        #  Not doing a bond, so  we check to make sure the requested device,
        #  nc['BOND_DEVICE'], is in the list of devices with carrier:
        if nc['BOND_DEVICE'] == '""':
            # In this case we have no network interface in the configuration but we
            # network settings.
            # First we check how many net_devs we have:
            if TO_LOG:
                print("nc['BOND_DEVICE']=''", len(net_devs), net_devs)
            if len(net_devs) == 1: # This is a dict of speed: devices
                speeds = net_devs.keys()
                speeds.sort(reverse=True)
                speed = speeds[0]
                possibles = net_devs[speed]
                if TO_LOG:
                    print(possibles)
                # At this point we have options, but no information, so:
                notes = "No device in the configuration file and multiple devices found. Picking the first"
                iface_list = [possibles[0]]
        else:
            if TO_LOG:
                print("inner else")
            for speed in net_devs:
                if nc['BOND_DEVICE'] in net_devs[speed]:
                    iface_list = [nc['BOND_DEVICE']]
                else:
                    iface_list = [nc['BOND_DEVICE']]
                    notes = "{0} did not have carrier at install time, and may not work".format(nc['BOND_DEVICE'])
    elif iface_speed != 'auto':
        if len(net_devs[iface_speed]) > 0:
            iface_list = net_devs[iface_speed]
        else:
            notes = "no devices set to {0}".format(iface_speed)
    else: # This SHOULD be iface_speed == auto, and nc['BOND_DEVCE'] containing bond.
        #  if not it is anyway.
        # Thus we are doing a bond of some sort.
        # This gives us the fastest interfaces first:
        speeds = [k for k in sorted(net_devs.keys(), reverse=True)]
        fastest = speeds[0]
        # Walk through "speeds" and take the first one that has more than one
        # interface. This will only set iface_list if there are 2 or more interfaces:
        # previous_speed = 0
        for i in speeds:
            if len(net_devs[i]) > 1:
                iface_list = net_devs[i]
                break
        if TO_LOG:
            print("iface list:", iface_list)
        # if iface_list is still false, and we are requesting a bond, we will
        # want the fastest interface with link:
        if (iface_list == False) and ("bond" in nc['BOND_DEVICE'].lower()):
                if TO_LOG:
                    print(len(net_devs), net_devs, i)
                if len(net_devs) == 0:
                    iface_list = net_devs
                    notes = "no devices found for the bond. Will not have network after reboot"
                else:
                    iface_list = net_devs[fastest] # This is assuming that we'll want to bond the fastest interaface.
                    if TO_LOG:
                        print("dev:", net_devs[fastest])
    if TO_LOG:
        print(iface_list, notes)
    return iface_list, notes


# Find our network configuration file:
if os.path.isfile(os.path.join(base_cfg_dir,cfg_file)):
    cfg_path = os.path.join(base_cfg_dir,cfg_file)
elif os.path.isfile(cfg_file):
    cfg_path = cfg_file
else:
    cfg_path = ''

if cfg_path:
    nc = read_config(cfg_path)
else:
    # if we don't have a working config file we use this
    # The IPs and hostnames are bad.
    nc = { IPADDR:"10.0.0.2",
        NETMASK:"255.255.255.252",
        GATEWAY:"10.0.0.1",
        BOND_DEVICE:"bond0",
        MTU:"9000",
        NAMESERVER:"192.168.0.1",
        HOSTNAME:"bad.example.com",
        NETWORKING_IPV6:"yes",
        IPV6ADDR:" 2001:0db8:0a0b:12f0:0000:0000:0000:0002/64",
        IPV6_DEFAULTGW:" 2001:0db8:0a0b:12f0:0000:0000:0000:0001",
        BONDING_OPTS:"miimon=100 mode=4 lacp_rate=fast xmit_hash_policy=layer3+4",
        DHCP:"no" }
# This should be set to no in the config file, but that could change:
if "DHCP" not in nc:
    nc['DHCP']='no'

net_devs = find_usable_net_devs(iface_dir)
bondable_iface, iface_problems = useable_interfaces(net_devs, nc, iface_speed)

# turn bondable_iface into a string for the network line:
if bondable_iface and len(bondable_iface) > 1:
    dev_list = bondable_iface
    dev_str = dev_list.pop()
    for d in dev_list:
        dev_str = dev_str + "," + d
else:
    dev_str = bondable_iface[0]


if ('y' in nc['NETWORKING_IPV6'].lower()) and re.search(":",nc['IPV6ADDR']):
    IPV6 = "--ipv6=" + nc["IPV6ADDR"]
else:
    if 'y' in nc['NETWORKING_IPV6'].lower():
        if iface_problems is False:
            iface_problems = "IPv6 enabled but no address provided"
        else:
            iface_problems = "{0} and IPv6 enabled but no address provided".format(iface_problems)
    if re.search(":",nc['IPV6ADDR']):
        if iface_problems is False:
            iface_problems = "IPv6 is disabled, but IPV6ADDR was set to {0}".format(nc['IPV6ADDR'])
        else:
            iface_problems = "{0} and IPv6 is disabled, but IPV6ADDR was set to {1}".format(iface_problems, nc['IPV6ADDR'])
    IPV6 = "--noipv6"

if "bond" in nc['BOND_DEVICE'].lower():
    bond_stuff = "--device={BOND_DEVICE} --bondslaves={0} --bondopts={BONDING_OPTS}".format(dev_str, **nc)
elif nc['BOND_DEVICE'] in dev_str:
    bond_stuff = "--device={0}".format(nc["BOND_DEVICE"])
elif bondable_iface and nc['BOND_DEVICE'] == '""' :
    print("**")
    print("No device (BOND_DEVICE) specified it he config, found", bondable_iface, "with link, using it.")
    print("**")
    bond_stuff = "--device={0}".format(bondable_iface[0])
else:
    print("**")
    print(nc["BOND_DEVICE"], "not found within $usable_devices, setting anyway, this probably won't work")
    print("**")
    bond_stuff = "--device={0}".format(nc["BOND_DEVICE"])

if 'yes' in nc['DHCP'].lower() or not bondable_iface:
    network_line = "network --bootproto=dhcp --device={BOND_DEVICE} --hostname={HOSTNAME}".format(**nc)
else:
    network_line = "network --bootproto=static {0} --activate {1} --ip={IPADDR} --netmask={NETMASK} --gateway={GATEWAY} --nameserver={NAMESERVER} --mtu={MTU} --hostname={HOSTNAME} \n".format(
            bond_stuff, IPV6, **nc)

if iface_problems:
    network_line = "# Problems found: {0}\n{1}".format(iface_problems,network_line)

with open('/tmp/network_line','w') as OUT:
    OUT.write(network_line)
