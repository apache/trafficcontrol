package Fixtures::Integration::ToExtension;
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

my %definition_for = (
  'servercheck_aa' => {
    new   => 'ToExtension',
    using => {
      id                      => 1,
      name                    => 'ILO_PING',
      version                 => '1.0.0',
      info_url                => 'http://foo.com/bar.html',
      script_file             => 'ping',
      isactive                => '1',
      additional_config_json  => '{ "select": "ilo_ip_address", "cron": "9 * * * *" }',
      servercheck_column_name => 'aa',
      servercheck_short_name  => 'ILO',
      type                    => '31',
    },
  },
  'servercheck_ab' => {
    new   => 'ToExtension',
    using => {
      id                      => 2,
      name                    => '10G_PING',
      version                 => '1.0.0',
      info_url                => 'http://foo.com/bar.html',
      script_file             => 'ping',
      isactive                => '1',
      additional_config_json  => '{ "select": "ip_address", "cron": "18 * * * *" }',
      servercheck_column_name => 'ab',
      servercheck_short_name  => '10G',
      type                    => '31',
    },
  },
  'servercheck_ac' => {
    new   => 'ToExtension',
    using => {
      id                      => 3,
      name                    => 'FQDN_PING',
      version                 => '1.0.0',
      info_url                => 'http://foo.com/bar.html',
      script_file             => 'ping',
      isactive                => '1',
      additional_config_json  => '{ "select": "host_name", "cron": "27 * * * *" }',
      servercheck_column_name => 'ac',
      servercheck_short_name  => 'FQDN',
      type                    => '31',
    },
  },
  'servercheck_ad' => {
    new   => 'ToExtension',
    using => {
      id                      => 4,
      name                    => 'CHECK_DSCP',
      version                 => '1.0.0',
      info_url                => 'http://foo.com/bar.html',
      script_file             => 'dscp',
      isactive                => '1',
      additional_config_json  => '{ "select": "ilo_ip_address", "cron": "36 * * * *" }',
      servercheck_column_name => 'ad',
      servercheck_short_name  => 'DSCP',
      type                    => '31',
    },
  },
  'servercheck_ae' => {
    new   => 'ToExtension',
    using => {
      id                      => 5,
      name                    => 'OPEN',
      version                 => '1.0.0',
      info_url                => 'http://foo.com/bar.html',
      script_file             => 'dscp',
      isactive                => '0',
      additional_config_json  => '',
      servercheck_column_name => 'ae',
      servercheck_short_name  => '',
      type                    => '33',
    },
  },
  'servercheck_af' => {
    new   => 'ToExtension',
    using => {
      id                      => 6,
      name                    => 'OPEN',
      version                 => '1.0.0',
      info_url                => 'http://foo.com/bar.html',
      script_file             => 'dscp',
      isactive                => '0',
      additional_config_json  => '',
      servercheck_column_name => 'af',
      servercheck_short_name  => '',
      type                    => '33',
    },
  },
  'servercheck_ag' => {
    new   => 'ToExtension',
    using => {
      id                      => 7,
      name                    => 'IPV6_PING',
      version                 => '1.0.0',
      info_url                => 'http://foo.com/bar.html',
      script_file             => 'ping',
      isactive                => '1',
      additional_config_json  => '{ "select": "ip6_address", "cron": "0 * * * *" }',
      servercheck_column_name => 'ag',
      servercheck_short_name  => '10G6',
      type                    => '31',
    },
  },
  'servercheck_ah' => {
    new   => 'ToExtension',
    using => {
      id                      => 8,
      name                    => 'OPEN',
      version                 => '1.0.0',
      info_url                => 'http://foo.com/bar.html',
      script_file             => 'dscp',
      isactive                => '0',
      additional_config_json  => '',
      servercheck_column_name => 'ah',
      servercheck_short_name  => '',
      type                    => '33',
    },
  },
  'servercheck_ai' => {
    new   => 'ToExtension',
    using => {
      id                      => 9,
      name                    => 'CHECK_STATS',
      version                 => '1.0.0',
      info_url                => 'http://foo.com/bar.html',
      script_file             => 'ping',
      isactive                => '1',
      additional_config_json  => '{ "select": "ilo_ip_address", "cron": "54 * * * *" }',
      servercheck_column_name => 'ai',
      servercheck_short_name  => 'STAT',
      type                    => '31',
    },
  },
  'servercheck_aj' => {
    new   => 'ToExtension',
    using => {
      id                      => 10,
      name                    => 'OPEN',
      version                 => '1.0.0',
      info_url                => 'http://foo.com/bar.html',
      script_file             => 'dscp',
      isactive                => '0',
      additional_config_json  => '',
      servercheck_column_name => 'aj',
      servercheck_short_name  => '',
      type                    => '33',
    },
  },
  'servercheck_ak' => {
    new   => 'ToExtension',
    using => {
      id                      => 11,
      name                    => 'CHECK_MTU',
      version                 => '1.0.0',
      info_url                => 'http://foo.com/bar.html',
      script_file             => 'ping',
      isactive                => '1',
      additional_config_json  => '{ "select": "ip_address", "cron": "45 * * * *" }',
      servercheck_column_name => 'ak',
      servercheck_short_name  => 'MTU',
      type                    => '31',
    },
  },
  'servercheck_al' => {
    new   => 'ToExtension',
    using => {
      id                      => 12,
      name                    => 'CHECK_TRAFFIC_ROUTER_STATUS',
      version                 => '1.0.0',
      info_url                => 'http://foo.com/bar.html',
      script_file             => 'ping',
      isactive                => '1',
      additional_config_json  => '{ "select": "ilo_ip_address", "cron": "10 * * * *" }',
      servercheck_column_name => 'al',
      servercheck_short_name  => 'TRTR',
      type                    => '31',
    },
  },
  'servercheck_am' => {
    new   => 'ToExtension',
    using => {
      id                      => 13,
      name                    => 'CHECK_TRAFFIC_MONITOR_STATUS',
      version                 => '1.0.0',
      info_url                => 'http://foo.com/bar.html',
      script_file             => 'ping',
      isactive                => '1',
      additional_config_json  => '{ "select": "ip_address", "cron": "10 * * * *" }',
      servercheck_column_name => 'am',
      servercheck_short_name  => 'TRMO',
      type                    => '31',
    },
  },
  'servercheck_an' => {
    new   => 'ToExtension',
    using => {
      id                      => 14,
      name                    => 'CACHE_HIT_RATIO_LAST_15',
      version                 => '1.0.0',
      info_url                => 'http://foo.com/bar.html',
      script_file             => 'ping',
      isactive                => '1',
      additional_config_json  => '{ "select": "ilo_ip_address", "cron": "0,15,30,45 * * * *" }',
      servercheck_column_name => 'an',
      servercheck_short_name  => 'CHR',
      type                    => '32',
    },
  },
  'servercheck_ao' => {
    new   => 'ToExtension',
    using => {
      id                      => 15,
      name                    => 'DISK_UTILIZATION',
      version                 => '1.0.0',
      info_url                => 'http://foo.com/bar.html',
      script_file             => 'ping',
      isactive                => '1',
      additional_config_json  => '{ "select": "ilo_ip_address", "cron": "20 * * * *" }',
      servercheck_column_name => 'ao',
      servercheck_short_name  => 'CDU',
      type                    => '32',
    },
  },
  'servercheck_ap' => {
    new   => 'ToExtension',
    using => {
      id                      => 16,
      name                    => 'ORT_ERROR_COUNT',
      version                 => '1.0.0',
      info_url                => 'http://foo.com/bar.html',
      script_file             => 'ping',
      isactive                => '1',
      additional_config_json  => '{ "select": "ilo_ip_address", "cron": "40 * * * *" }',
      servercheck_column_name => 'ap',
      servercheck_short_name  => 'ORT',
      type                    => '32',
    },
  },
  'servercheck_aq' => {
    new   => 'ToExtension',
    using => {
      id                      => 17,
      name                    => 'OPEN',
      version                 => '1.0.0',
      info_url                => 'http://foo.com/bar.html',
      script_file             => 'dscp',
      isactive                => '0',
      additional_config_json  => '',
      servercheck_column_name => 'aq',
      servercheck_short_name  => '',
      type                    => '33',
    },
  },
  'servercheck_ar' => {
    new   => 'ToExtension',
    using => {
      id                      => 18,
      name                    => 'OPEN',
      version                 => '1.0.0',
      info_url                => 'http://foo.com/bar.html',
      script_file             => 'dscp',
      isactive                => '0',
      additional_config_json  => '',
      servercheck_column_name => 'ar',
      servercheck_short_name  => '',
      type                    => '33',
    },
  },
  'servercheck_as' => {
    new   => 'ToExtension',
    using => {
      id                      => 19,
      name                    => 'OPEN',
      version                 => '1.0.0',
      info_url                => 'http://foo.com/bar.html',
      script_file             => 'dscp',
      isactive                => '0',
      additional_config_json  => '',
      servercheck_column_name => 'as',
      servercheck_short_name  => '',
      type                    => '33',
    },
  },
  'servercheck_at' => {
    new   => 'ToExtension',
    using => {
      id                      => 20,
      name                    => 'OPEN',
      version                 => '1.0.0',
      info_url                => 'http://foo.com/bar.html',
      script_file             => 'dscp',
      isactive                => '0',
      additional_config_json  => '',
      servercheck_column_name => 'at',
      servercheck_short_name  => '',
      type                    => '33',
    },
  },
  'servercheck_au' => {
    new   => 'ToExtension',
    using => {
      id                      => 21,
      name                    => 'OPEN',
      version                 => '1.0.0',
      info_url                => 'http://foo.com/bar.html',
      script_file             => 'dscp',
      isactive                => '0',
      additional_config_json  => '',
      servercheck_column_name => 'au',
      servercheck_short_name  => '',
      type                    => '33',
    },
  },
  'servercheck_av' => {
    new   => 'ToExtension',
    using => {
      id                      => 22,
      name                    => 'OPEN',
      version                 => '1.0.0',
      info_url                => 'http://foo.com/bar.html',
      script_file             => 'dscp',
      isactive                => '0',
      additional_config_json  => '',
      servercheck_column_name => 'av',
      servercheck_short_name  => '',
      type                    => '33',
    },
  },
  'servercheck_aw' => {
    new   => 'ToExtension',
    using => {
      id                      => 23,
      name                    => 'OPEN',
      version                 => '1.0.0',
      info_url                => 'http://foo.com/bar.html',
      script_file             => 'dscp',
      isactive                => '0',
      additional_config_json  => '',
      servercheck_column_name => 'aw',
      servercheck_short_name  => '',
      type                    => '33',
    },
  },
  'servercheck_ax' => {
    new   => 'ToExtension',
    using => {
      id                      => 24,
      name                    => 'OPEN',
      version                 => '1.0.0',
      info_url                => 'http://foo.com/bar.html',
      script_file             => 'dscp',
      isactive                => '0',
      additional_config_json  => '',
      servercheck_column_name => 'ax',
      servercheck_short_name  => '',
      type                    => '33',
    },
  },
  'servercheck_ay' => {
    new   => 'ToExtension',
    using => {
      id                      => 25,
      name                    => 'OPEN',
      version                 => '1.0.0',
      info_url                => 'http://foo.com/bar.html',
      script_file             => 'dscp',
      isactive                => '0',
      additional_config_json  => '',
      servercheck_column_name => 'ay',
      servercheck_short_name  => '',
      type                    => '33',
    },
  },
  'servercheck_az' => {
    new   => 'ToExtension',
    using => {
      id                      => 26,
      name                    => 'OPEN',
      version                 => '1.0.0',
      info_url                => 'http://foo.com/bar.html',
      script_file             => 'dscp',
      isactive                => '0',
      additional_config_json  => '',
      servercheck_column_name => 'az',
      servercheck_short_name  => '',
      type                    => '33',
    },
  },
  'servercheck_ba' => {
    new   => 'ToExtension',
    using => {
      id                      => 27,
      name                    => 'OPEN',
      version                 => '1.0.0',
      info_url                => 'http://foo.com/bar.html',
      script_file             => 'dscp',
      isactive                => '0',
      additional_config_json  => '',
      servercheck_column_name => 'ba',
      servercheck_short_name  => '',
      type                    => '33',
    },
  },
  'servercheck_bb' => {
    new   => 'ToExtension',
    using => {
      id                      => 28,
      name                    => 'OPEN',
      version                 => '1.0.0',
      info_url                => 'http://foo.com/bar.html',
      script_file             => 'dscp',
      isactive                => '0',
      additional_config_json  => '',
      servercheck_column_name => 'bb',
      servercheck_short_name  => '',
      type                    => '33',
    },
  },
  'servercheck_bc' => {
    new   => 'ToExtension',
    using => {
      id                      => 29,
      name                    => 'OPEN',
      version                 => '1.0.0',
      info_url                => 'http://foo.com/bar.html',
      script_file             => 'dscp',
      isactive                => '0',
      additional_config_json  => '',
      servercheck_column_name => 'bc',
      servercheck_short_name  => '',
      type                    => '33',
    },
  },
  'servercheck_bd' => {
    new   => 'ToExtension',
    using => {
      id                      => 30,
      name                    => 'OPEN',
      version                 => '1.0.0',
      info_url                => 'http://foo.com/bar.html',
      script_file             => 'dscp',
      isactive                => '0',
      additional_config_json  => '',
      servercheck_column_name => 'bd',
      servercheck_short_name  => '',
      type                    => '33',
    },
  },
  'servercheck_be' => {
    new   => 'ToExtension',
    using => {
      id                      => 31,
      name                    => 'OPEN',
      version                 => '1.0.0',
      info_url                => 'http://foo.com/bar.html',
      script_file             => 'dscp',
      isactive                => '0',
      additional_config_json  => '',
      servercheck_column_name => 'be',
      servercheck_short_name  => '',
      type                    => '33',
    },
  },
);

sub name {
  return "ToExtension";
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
