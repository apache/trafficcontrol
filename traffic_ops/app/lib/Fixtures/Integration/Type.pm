package Fixtures::Integration::Type;
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
#
use Moose;
extends 'DBIx::Class::EasyFixture';
use namespace::autoclean;
use Digest::SHA1 qw(sha1_hex);

my %definition_for = (
  EDGE => {
    new   => 'Type',
    using => {
      id           => 1,
      name         => 'EDGE',
      description  => 'Edge Cache',
      use_in_table => 'server',
    },
  },
  MID => {
    new   => 'Type',
    using => {
      id           => 2,
      name         => 'MID',
      description  => 'Mid Tier Cache',
      use_in_table => 'server',
    },
  },
  ORG => {
    new   => 'Type',
    using => {
      id           => 3,
      name         => 'ORG',
      description  => 'Origin',
      use_in_table => 'server',
    },
  },
  CCR => {
    new   => 'Type',
    using => {
      id           => 4,
      name         => 'CCR',
      description  => 'Comcast Content Router (aka Traffic Router)',
      use_in_table => 'server',
    },
  },
  EDGE_LOC => {
    new   => 'Type',
    using => {
      id           => 6,
      name         => 'EDGE_LOC',
      description  => 'Edge Cachegroup',
      use_in_table => 'cachegroup',
    },
  },
  MID_LOC => {
    new   => 'Type',
    using => {
      id           => 7,
      name         => 'MID_LOC',
      description  => 'Mid Cachegroup',
      use_in_table => 'cachegroup',
    },
  },
  HTTP => {
    new   => 'Type',
    using => {
      id           => 8,
      name         => 'HTTP',
      description  => 'HTTP Content Routing',
      use_in_table => 'deliveryservice',
    },
  },
  DNS => {
    new   => 'Type',
    using => {
      id           => 9,
      name         => 'DNS',
      description  => 'DNS Content Routing',
      use_in_table => 'deliveryservice',
    },
  },
  RIAK_SERVER => {
    new   => 'Type',
    using => {
      id           => 10,
      name         => 'RIAK',
      description  => 'Riak keystore',
      use_in_table => 'server',
    },
  },
  HTTP_NO_CACHE => {
    new   => 'Type',
    using => {
      id           => 11,
      name         => 'HTTP_NO_CACHE',
      description  => 'HTTP Content Routing, no caching',
      use_in_table => 'deliveryservice',
    },
  },
  HTTP_LIVE => {
    new   => 'Type',
    using => {
      id           => 13,
      name         => 'HTTP_LIVE',
      description  => 'HTTP Content routing cache in RAM',
      use_in_table => 'deliveryservice',
    },
  },
  RASCAL => {
    new   => 'Type',
    using => {
      id           => 15,
      name         => 'RASCAL',
      description  => 'Rascal (aka Traffic Monitor) server',
      use_in_table => 'server',
    },
  },
  HOST_REGEXP => {
    new   => 'Type',
    using => {
      id           => 18,
      name         => 'HOST_REGEXP',
      description  => 'Host header regular expression',
      use_in_table => 'regex',
    },
  },
  PATH_REGEXP => {
    new   => 'Type',
    using => {
      id           => 19,
      name         => 'PATH_REGEXP',
      description  => 'URL path regular expression',
      use_in_table => 'regex',
    },
  },
  HEADER_REGEXP => {
    new   => 'Type',
    using => {
      id           => 20,
      name         => 'HEADER_REGEXP',
      description  => 'HTTP header regular expression',
      use_in_table => 'regex',
    },
  },
  A_RECORD => {
    new   => 'Type',
    using => {
      id           => 21,
      name         => 'A_RECORD',
      description  => 'Static DNS A entry',
      use_in_table => 'staticdnsentry',
    },
  },
  AAAA_RECORD => {
    new   => 'Type',
    using => {
      id           => 22,
      name         => 'AAAA_RECORD',
      description  => 'Static DNS AAAA entry',
      use_in_table => 'staticdnsentry',
    },
  },
  CNAME_RECORD => {
    new   => 'Type',
    using => {
      id           => 23,
      name         => 'CNAME_RECORD',
      description  => 'Static DNS CNAME entry',
      use_in_table => 'staticdnsentry',
    },
  },
  HTTP_LIVE_NATNL => {
    new   => 'Type',
    using => {
      id           => 24,
      name         => 'HTTP_LIVE_NATNL',
      description  => 'HTTP Content routing, RAM cache, National',
      use_in_table => 'deliveryservice',
    },
  },
  REDIS => {
    new   => 'Type',
    using => {
      id           => 25,
      name         => 'REDIS',
      description  => 'Redis stats gateway',
      use_in_table => 'server',
    },
  },
  DNS_LIVE_NATNL => {
    new   => 'Type',
    using => {
      id           => 26,
      name         => 'DNS_LIVE_NATNL',
      description  => 'DNS Content routing, RAM cache, National',
      use_in_table => 'deliveryservice',
    },
  },
  DNS_LIVE => {
    new   => 'Type',
    using => {
      id          => 27,
      name        => 'DNS_LIVE',
      description => 'DNS Content routing, RAM cache, Lo
cal',
      use_in_table => 'deliveryservice',
    },
  },
  LOCAL => {
    new   => 'Type',
    using => {
      id           => 28,
      name         => 'LOCAL',
      description  => 'Local User',
      use_in_table => 'tm_user',
    },
  },
  ACTIVE_DIRECTORY => {
    new   => 'Type',
    using => {
      id           => 29,
      name         => 'ACTIVE_DIRECTORY',
      description  => 'Active Directory User',
      use_in_table => 'tm_user',
    },
  },
  TOOLS_SERVER => {
    new   => 'Type',
    using => {
      id           => 30,
      name         => 'TOOLS_SERVER',
      description  => 'Ops hosts for managment ',
      use_in_table => 'server',
    },
  },
  CHECK_PLUGIN_BOOL => {
    new   => 'Type',
    using => {
      id           => 31,
      name         => 'CHECK_EXTENSION_BOOL',
      description  => 'TO Extension for checkmark in Server Check',
      use_in_table => 'to_extension',
    },
  },
  CHECK_PLUGIN_INT => {
    new   => 'Type',
    using => {
      id           => 32,
      name         => 'CHECK_EXTENSION_NUM',
      description  => 'TO Extenstion for int value in Server Check',
      use_in_table => 'to_extension',
    },
  },
  CHECK_PLUGIN_OPEN_SLOT => {
    new   => 'Type',
    using => {
      id           => 33,
      name         => 'CHECK_EXTENSION_OPEN_SLOT',
      description  => 'Open slot for check in Server Status',
      use_in_table => 'to_extension',
    },
  },
  CONFIG_PLUGIN => {
    new   => 'Type',
    using => {
      id           => 34,
      name         => 'CONFIG_EXTENSION',
      description  => 'Extension for additional configuration file',
      use_in_table => 'to_extension',
    },
  },
  STATISTIC_PLUGIN => {
    new   => 'Type',
    using => {
      id           => 35,
      name         => 'STATISTIC_EXTENSION',
      description  => 'Extension source for 12M graphs',
      use_in_table => 'to_extension',
    },
  },
);

sub name {
  return "Type";
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
