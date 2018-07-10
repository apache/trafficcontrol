use utf8;
package Schema::Result::ProfileTypeValue;

# Created by DBIx::Class::Schema::Loader
# DO NOT MODIFY THE FIRST PART OF THIS FILE

=head1 NAME

Schema::Result::ProfileTypeValue

=cut

use strict;
use warnings;

use base 'DBIx::Class::Core';
__PACKAGE__->table_class("DBIx::Class::ResultSource::View");

=head1 TABLE: C<profile_type_values>

=cut

__PACKAGE__->table("profile_type_values");
__PACKAGE__->result_source_instance->view_definition(" SELECT unnest(enum_range(NULL::profile_type)) AS value\n  ORDER BY (unnest(enum_range(NULL::profile_type)))");

=head1 ACCESSORS

=head2 value

  data_type: 'enum'
  extra: {custom_type_name => "profile_type",list => ["ATS_PROFILE","TR_PROFILE","TM_PROFILE","TS_PROFILE","TP_PROFILE","INFLUXDB_PROFILE","RIAK_PROFILE","SPLUNK_PROFILE","DS_PROFILE","ORG_PROFILE","KAFKA_PROFILE","LOGSTASH_PROFILE","ES_PROFILE","UNK_PROFILE","GROVE_PROFILE"]}
  is_nullable: 1

=cut

__PACKAGE__->add_columns(
  "value",
  {
    data_type => "enum",
    extra => {
      custom_type_name => "profile_type",
      list => [
        "ATS_PROFILE",
        "TR_PROFILE",
        "TM_PROFILE",
        "TS_PROFILE",
        "TP_PROFILE",
        "INFLUXDB_PROFILE",
        "RIAK_PROFILE",
        "SPLUNK_PROFILE",
        "DS_PROFILE",
        "ORG_PROFILE",
        "KAFKA_PROFILE",
        "LOGSTASH_PROFILE",
        "ES_PROFILE",
        "UNK_PROFILE",
        "GROVE_PROFILE",
      ],
    },
    is_nullable => 1,
  },
);


# Created by DBIx::Class::Schema::Loader v0.07049 @ 2018-07-09 22:32:39
# DO NOT MODIFY THIS OR ANYTHING ABOVE! md5sum:YGAQ8Wqg7vVCtrtfiBIcAw


# You can replace this text with custom code or comments, and it will be preserved on regeneration
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
#
1;
