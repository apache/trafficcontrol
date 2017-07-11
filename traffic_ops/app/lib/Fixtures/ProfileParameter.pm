package Fixtures::ProfileParameter;
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
use Moose;
extends 'DBIx::Class::EasyFixture';
use namespace::autoclean;
use Digest::SHA1 qw(sha1_hex);

my %definition_for = (
    rascal_properties1 => {
        new   => 'ProfileParameter',
        using => {
            profile   => 200,
            parameter => 4,
        },
    },
    rascal_properties2 => {
        new   => 'ProfileParameter',
        using => {
            profile   => 200,
            parameter => 5,
        },
    },
    rascal_properties3 => {
        new   => 'ProfileParameter',
        using => {
            profile   => 100,
            parameter => 6,
        },
    },
    rascal_properties4 => {
        new   => 'ProfileParameter',
        using => {
            profile   => 200,
            parameter => 6,
        },
    },
    edge1_key0 => {
        new   => 'ProfileParameter',
        using => {
            profile   => 100,
            parameter => 7,
        },
    },
    edge1_key1 => {
        new   => 'ProfileParameter',
        using => {
            profile   => 100,
            parameter => 8,
        },
    },
    edge1_key2 => {
        new   => 'ProfileParameter',
        using => {
            profile   => 100,
            parameter => 9,
        },
    },
    edge1_key3 => {
        new   => 'ProfileParameter',
        using => {
            profile   => 100,
            parameter => 10,
        },
    },
    edge1_key4 => {
        new   => 'ProfileParameter',
        using => {
            profile   => 100,
            parameter => 11,
        },
    },
    edge1_key5 => {
        new   => 'ProfileParameter',
        using => {
            profile   => 100,
            parameter => 12,
        },
    },
    edge1_key6 => {
        new   => 'ProfileParameter',
        using => {
            profile   => 100,
            parameter => 13,
        },
    },
    edge1_key7 => {
        new   => 'ProfileParameter',
        using => {
            profile   => 100,
            parameter => 14,
        },
    },
    edge1_key8 => {
        new   => 'ProfileParameter',
        using => {
            profile   => 100,
            parameter => 15,
        },
    },
    edge1_key9 => {
        new   => 'ProfileParameter',
        using => {
            profile   => 100,
            parameter => 16,
        },
    },
    edge1_key10 => {
        new   => 'ProfileParameter',
        using => {
            profile   => 100,
            parameter => 17,
        },
    },
    edge1_key11 => {
        new   => 'ProfileParameter',
        using => {
            profile   => 100,
            parameter => 18,
        },
    },
    edge1_key12 => {
        new   => 'ProfileParameter',
        using => {
            profile   => 100,
            parameter => 19,
        },
    },
    edge1_key13 => {
        new   => 'ProfileParameter',
        using => {
            profile   => 100,
            parameter => 20,
        },
    },
    edge1_key14 => {
        new   => 'ProfileParameter',
        using => {
            profile   => 100,
            parameter => 21,
        },
    },
    edge1_key15 => {
        new   => 'ProfileParameter',
        using => {
            profile   => 100,
            parameter => 22,
        },
    },
    'edge1_url_sig_cdl-c2.config_location' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 100,
            parameter => 23,
        },
    },
    'edge1_error_url' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 100,
            parameter => 24,
        },
    },
    'edge1_CONFIG-proxy.config.allocator.debug_filter' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 100,
            parameter => 25,
        },
    },
    'edge1_CONFIG-proxy.config.allocator.enable_reclaim' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 100,
            parameter => 26,
        },
    },
    'edge1_CONFIG-proxy.config.allocator.max_overage' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 100,
            parameter => 27,
        },
    },
    'edge1_CONFIG-proxy.config.diags.show_location' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 100,
            parameter => 28,
        },
    },
    'edge1_CONFIG-proxy.config.http.cache.allow_empty_doc' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 100,
            parameter => 29,
        },
    },
    'LOCAL-proxy.config.cache.interim.storage' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 100,
            parameter => 30,
        },
    },
    'edge1_CONFIG-proxy.config.http.parent_proxy.file' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 100,
            parameter => 31,
        },
    },
    'edge1_12M_location' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 100,
            parameter => 32,
        },
    },
    'edge1_cacheurl_location' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 100,
            parameter => 33,
        },
    },
    'edge1_ip_allow_location' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 100,
            parameter => 34,
        },
    },
    'edge1_astats_over_http.so' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 100,
            parameter => 35,
        },
    },
    'edge1_crontab_root_location' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 100,
            parameter => 36,
        },
    },
    'edge1_hdr_rw_cdl-c2.config_location' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 100,
            parameter => 37,
        },
    },
    'edge1_50-ats.rules_location' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 100,
            parameter => 38,
        },
    },
    'edge1_parent.config_location' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 100,
            parameter => 39,
        },
    },
    'edge1_remap.config_location' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 100,
            parameter => 40,
        },
    },
    'edge1_drop_qstring.config_location' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 100,
            parameter => 41,
        },
    },
    'edge1_LogFormat.Format' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 100,
            parameter => 42,
        },
    },
    'edge1_LogFormat.Name' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 100,
            parameter => 43,
        },
    },
    'edge1_LogObject.Format' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 100,
            parameter => 44,
        },
    },
    'edge1_LogObject.Filename' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 100,
            parameter => 45,
        },
    },
    'edge1_cache.config_location' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 100,
            parameter => 46,
        },
    },
    'edge1_CONFIG-proxy.config.cache.control.filename' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 100,
            parameter => 47,
        },
    },
    'edge1_regex_revalidate.so' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 100,
            parameter => 48,
        },
    },
    'edge1_regex_revalidate.config_location' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 100,
            parameter => 49,
        },
    },
    'edge1_hosting.config_location' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 100,
            parameter => 50,
        },
    },
    'edge1_volume.config_location' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 100,
            parameter => 51,
        },
    },
    'edge1_allow_ip' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 100,
            parameter => 52,
        },
    },
    'edge1_allow_ip6' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 100,
            parameter => 53,
        },
    },
    'edge1_record_types' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 100,
            parameter => 54,
        },
    },
    'edge1_astats.config_location' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 100,
            parameter => 55,
        },
    },
    'edge1_astats.config_path' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 100,
            parameter => 56,
        },
    },
    'edge1_storage.config_location' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 100,
            parameter => 57,
        },
    },
    'edge1_Drive_Prefix' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 100,
            parameter => 58,
        },
    },
    'edge1_Drive_Letters' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 100,
            parameter => 59,
        },
    },
    'edge1_Disk_Volume' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 100,
            parameter => 60,
        },
    },
    'edge1_CONFIG-proxy.config.hostdb.storage_size' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 100,
            parameter => 61,
        },
    },
    mid1_key0 => {
        new   => 'ProfileParameter',
        using => {
            profile   => 200,
            parameter => 7,
        },
    },
    mid1_key1 => {
        new   => 'ProfileParameter',
        using => {
            profile   => 200,
            parameter => 8,
        },
    },
    mid1_key2 => {
        new   => 'ProfileParameter',
        using => {
            profile   => 200,
            parameter => 9,
        },
    },
    mid1_key3 => {
        new   => 'ProfileParameter',
        using => {
            profile   => 200,
            parameter => 10,
        },
    },
    mid1_key4 => {
        new   => 'ProfileParameter',
        using => {
            profile   => 200,
            parameter => 11,
        },
    },
    mid1_key5 => {
        new   => 'ProfileParameter',
        using => {
            profile   => 200,
            parameter => 12,
        },
    },
    mid1_key6 => {
        new   => 'ProfileParameter',
        using => {
            profile   => 200,
            parameter => 13,
        },
    },
    mid1_key7 => {
        new   => 'ProfileParameter',
        using => {
            profile   => 200,
            parameter => 14,
        },
    },
    mid1_key8 => {
        new   => 'ProfileParameter',
        using => {
            profile   => 200,
            parameter => 15,
        },
    },
    mid1_key9 => {
        new   => 'ProfileParameter',
        using => {
            profile   => 200,
            parameter => 16,
        },
    },
    mid1_key10 => {
        new   => 'ProfileParameter',
        using => {
            profile   => 200,
            parameter => 17,
        },
    },
    mid1_key11 => {
        new   => 'ProfileParameter',
        using => {
            profile   => 200,
            parameter => 18,
        },
    },
    mid1_key12 => {
        new   => 'ProfileParameter',
        using => {
            profile   => 200,
            parameter => 19,
        },
    },
    mid1_key13 => {
        new   => 'ProfileParameter',
        using => {
            profile   => 200,
            parameter => 20,
        },
    },
    mid1_key14 => {
        new   => 'ProfileParameter',
        using => {
            profile   => 200,
            parameter => 21,
        },
    },
    mid1_key15 => {
        new   => 'ProfileParameter',
        using => {
            profile   => 200,
            parameter => 22,
        },
    },
    'mid1_url_sig_cdl-c2.config_location' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 200,
            parameter => 23,
        },
    },
    'mid1_error_url' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 200,
            parameter => 24,
        },
    },
    'mid1_CONFIG-proxy.config.allocator.debug_filter' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 200,
            parameter => 25,
        },
    },
    'mid1_CONFIG-proxy.config.allocator.enable_reclaim' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 200,
            parameter => 26,
        },
    },
    'mid1_CONFIG-proxy.config.allocator.max_overage' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 200,
            parameter => 27,
        },
    },
    'mid1_CONFIG-proxy.config.diags.show_location' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 200,
            parameter => 28,
        },
    },
    'mid1_CONFIG-proxy.config.http.cache.allow_empty_doc' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 200,
            parameter => 29,
        },
    },
    'LOCAL-proxy.config.cache.interim.storage' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 200,
            parameter => 30,
        },
    },
    'mid1_CONFIG-proxy.config.http.parent_proxy.file' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 200,
            parameter => 31,
        },
    },
    'mid1_12M_location' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 200,
            parameter => 32,
        },
    },
    'mid1_cacheurl_location' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 200,
            parameter => 33,
        },
    },
    'mid1_ip_allow_location' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 200,
            parameter => 34,
        },
    },
    'mid1_astats_over_http.so' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 200,
            parameter => 35,
        },
    },
    'mid1_crontab_root_location' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 200,
            parameter => 36,
        },
    },
    'mid1_hdr_rw_cdl-c2.config_location' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 200,
            parameter => 37,
        },
    },
    'mid1_50-ats.rules_location' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 200,
            parameter => 38,
        },
    },
    'mid1_parent.config_location' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 200,
            parameter => 39,
        },
    },
    'mid1_remap.config_location' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 200,
            parameter => 40,
        },
    },
    'mid1_drop_qstring.config_location' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 200,
            parameter => 41,
        },
    },
    'mid1_LogFormat.Format' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 200,
            parameter => 42,
        },
    },
    'mid1_LogFormat.Name' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 200,
            parameter => 43,
        },
    },
    'mid1_LogObject.Format' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 200,
            parameter => 44,
        },
    },
    'mid1_LogObject.Filename' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 200,
            parameter => 45,
        },
    },
    'mid1_cache.config_location' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 200,
            parameter => 46,
        },
    },
    'mid1_CONFIG-proxy.config.cache.control.filename' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 200,
            parameter => 47,
        },
    },
    'mid1_regex_revalidate.so' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 200,
            parameter => 48,
        },
    },
    'mid1_regex_revalidate.config_location' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 200,
            parameter => 49,
        },
    },
    'mid1_hosting.config_location' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 200,
            parameter => 50,
        },
    },
    'mid1_volume.config_location' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 200,
            parameter => 51,
        },
    },
    'mid1_allow_ip' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 200,
            parameter => 52,
        },
    },
    'mid1_allow_ip6' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 200,
            parameter => 53,
        },
    },
    'mid1_record_types' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 200,
            parameter => 54,
        },
    },
    'mid1_astats.config_location' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 200,
            parameter => 55,
        },
    },
    'mid1_astats.config_path' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 200,
            parameter => 56,
        },
    },
    'mid1_storage.config_location' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 200,
            parameter => 57,
        },
    },
    'mid1_Drive_Prefix' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 200,
            parameter => 58,
        },
    },
    'mid1_Drive_Letters' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 200,
            parameter => 59,
        },
    },
    'mid1_Disk_Volume' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 200,
            parameter => 60,
        },
    },
    'mid1_CONFIG-proxy.config.hostdb.storage_size' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 200,
            parameter => 61,
        },
    },
    'edge1_package_trafficserver' => {
        new   => 'ProfileParameter',
        using => {
            profile   => 100,
            parameter => 66,
        },
    },
);

sub get_definition {
    my ( $self, $name ) = @_;
    return $definition_for{$name};
}

sub all_fixture_names {
    return keys %definition_for;
}

__PACKAGE__->meta->make_immutable;

1;
