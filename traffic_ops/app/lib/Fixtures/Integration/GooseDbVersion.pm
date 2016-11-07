package Fixtures::Integration::GooseDbVersion;

# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
# 
#   http://www.apache.org/licenses/LICENSE-2.0
# 
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.


# Do not edit! Generated code.
# See https://github.com/Comcast/traffic_control/wiki/The%20Kabletown%20example

use Moose;
extends 'DBIx::Class::EasyFixture';
use namespace::autoclean;

my %definition_for = (
'0' => { new => 'GooseDbVersion', => using => { id => '1', is_applied => '1', tstamp => '2015-12-04 07:46:20', version_id => '0', }, }, 
'1' => { new => 'GooseDbVersion', => using => { tstamp => '2015-12-04 07:46:21', version_id => '20141222103718', id => '2', is_applied => '1', }, }, 
'2' => { new => 'GooseDbVersion', => using => { version_id => '20150108100000', id => '3', is_applied => '1', tstamp => '2015-12-04 07:46:21', }, }, 
'3' => { new => 'GooseDbVersion', => using => { id => '4', is_applied => '1', tstamp => '2015-12-04 07:46:21', version_id => '20150205100000', }, }, 
'4' => { new => 'GooseDbVersion', => using => { version_id => '20150209100000', id => '5', is_applied => '1', tstamp => '2015-12-04 07:46:21', }, }, 
'5' => { new => 'GooseDbVersion', => using => { is_applied => '1', tstamp => '2015-12-04 07:46:21', version_id => '20150210100000', id => '6', }, }, 
'6' => { new => 'GooseDbVersion', => using => { tstamp => '2015-12-04 07:46:21', version_id => '20150304100000', id => '7', is_applied => '1', }, }, 
'7' => { new => 'GooseDbVersion', => using => { version_id => '20150310100000', id => '8', is_applied => '1', tstamp => '2015-12-04 07:46:21', }, }, 
'8' => { new => 'GooseDbVersion', => using => { id => '9', is_applied => '1', tstamp => '2015-12-04 07:46:21', version_id => '20150316100000', }, }, 
'9' => { new => 'GooseDbVersion', => using => { id => '10', is_applied => '1', tstamp => '2015-12-04 07:46:22', version_id => '20150331105256', }, }, 
'10' => { new => 'GooseDbVersion', => using => { id => '11', is_applied => '1', tstamp => '2015-12-04 07:46:22', version_id => '20150501100000', }, }, 
'11' => { new => 'GooseDbVersion', => using => { id => '12', is_applied => '1', tstamp => '2015-12-04 07:46:22', version_id => '20150503100001', }, }, 
'12' => { new => 'GooseDbVersion', => using => { is_applied => '1', tstamp => '2015-12-04 07:46:22', version_id => '20150504100000', id => '13', }, }, 
'13' => { new => 'GooseDbVersion', => using => { id => '14', is_applied => '1', tstamp => '2015-12-04 07:46:22', version_id => '20150504100001', }, }, 
'14' => { new => 'GooseDbVersion', => using => { id => '15', is_applied => '1', tstamp => '2015-12-04 07:46:22', version_id => '20150521100000', }, }, 
'15' => { new => 'GooseDbVersion', => using => { id => '16', is_applied => '1', tstamp => '2015-12-04 07:46:22', version_id => '20150530100000', }, }, 
'16' => { new => 'GooseDbVersion', => using => { tstamp => '2015-12-04 07:46:22', version_id => '20150618100000', id => '17', is_applied => '1', }, }, 
'17' => { new => 'GooseDbVersion', => using => { version_id => '20150626100000', id => '18', is_applied => '1', tstamp => '2015-12-04 07:46:22', }, }, 
'18' => { new => 'GooseDbVersion', => using => { is_applied => '1', tstamp => '2015-12-04 07:46:22', version_id => '20150706084134', id => '19', }, }, 
'19' => { new => 'GooseDbVersion', => using => { id => '20', is_applied => '1', tstamp => '2015-12-04 07:46:22', version_id => '20150721000000', }, }, 
'20' => { new => 'GooseDbVersion', => using => { id => '21', is_applied => '1', tstamp => '2015-12-04 07:46:22', version_id => '20150722100000', }, }, 
'21' => { new => 'GooseDbVersion', => using => { tstamp => '2015-12-04 07:46:22', version_id => '20150728000000', id => '22', is_applied => '1', }, }, 
'22' => { new => 'GooseDbVersion', => using => { id => '23', is_applied => '1', tstamp => '2015-12-04 07:46:22', version_id => '20150804000000', }, }, 
'23' => { new => 'GooseDbVersion', => using => { id => '24', is_applied => '1', tstamp => '2015-12-04 07:46:22', version_id => '20150807000000', }, }, 
'24' => { new => 'GooseDbVersion', => using => { is_applied => '1', tstamp => '2015-12-04 07:46:22', version_id => '20150825175644', id => '25', }, }, 
'25' => { new => 'GooseDbVersion', => using => { tstamp => '2015-12-04 07:46:22', version_id => '20150922092122', id => '26', is_applied => '1', }, }, 
'26' => { new => 'GooseDbVersion', => using => { id => '27', is_applied => '1', tstamp => '2015-12-04 07:46:22', version_id => '20150925020500', }, }, 
'27' => { new => 'GooseDbVersion', => using => { id => '28', is_applied => '1', tstamp => '2015-12-04 07:46:22', version_id => '20151020143912', }, }, 
'28' => { new => 'GooseDbVersion', => using => { is_applied => '1', tstamp => '2015-12-04 07:46:22', version_id => '20151021000000', id => '29', }, }, 
'29' => { new => 'GooseDbVersion', => using => { id => '30', is_applied => '1', tstamp => '2015-12-04 07:46:22', version_id => '20151027152323', }, }, 
'30' => { new => 'GooseDbVersion', => using => { tstamp => '2015-12-04 07:46:22', version_id => '20151107000000', id => '31', is_applied => '1', }, }, 
); 

sub name {
		return "GooseDbVersion";
}

sub get_definition { 
		my ( $self, $name ) = @_;
		return $definition_for{$name};
}

sub all_fixture_names {
		return keys %definition_for;
}

__PACKAGE__->meta->make_immutable;
1;
