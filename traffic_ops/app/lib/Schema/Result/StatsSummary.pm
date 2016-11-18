use utf8;
package Schema::Result::StatsSummary;

# Created by DBIx::Class::Schema::Loader
# DO NOT MODIFY THE FIRST PART OF THIS FILE

=head1 NAME

Schema::Result::StatsSummary

=cut

use strict;
use warnings;

use base 'DBIx::Class::Core';

=head1 TABLE: C<stats_summary>

=cut

__PACKAGE__->table("stats_summary");

=head1 ACCESSORS

=head2 id

  data_type: 'bigint'
  is_auto_increment: 1
  is_nullable: 0
  sequence: 'stats_summary_id_seq'

=head2 cdn_name

  data_type: 'text'
  default_value: 'all'
  is_nullable: 0

=head2 deliveryservice_name

  data_type: 'text'
  is_nullable: 0

=head2 stat_name

  data_type: 'text'
  is_nullable: 0

=head2 stat_value

  data_type: 'double precision'
  is_nullable: 0

=head2 summary_time

  data_type: 'timestamp with time zone'
  default_value: current_timestamp
  is_nullable: 0
  original: {default_value => \"now()"}

=head2 stat_date

  data_type: 'date'
  is_nullable: 1

=cut

__PACKAGE__->add_columns(
  "id",
  {
    data_type         => "bigint",
    is_auto_increment => 1,
    is_nullable       => 0,
    sequence          => "stats_summary_id_seq",
  },
  "cdn_name",
  { data_type => "text", default_value => "all", is_nullable => 0 },
  "deliveryservice_name",
  { data_type => "text", is_nullable => 0 },
  "stat_name",
  { data_type => "text", is_nullable => 0 },
  "stat_value",
  { data_type => "double precision", is_nullable => 0 },
  "summary_time",
  {
    data_type     => "timestamp with time zone",
    default_value => \"current_timestamp",
    is_nullable   => 0,
    original      => { default_value => \"now()" },
  },
  "stat_date",
  { data_type => "date", is_nullable => 1 },
);

=head1 PRIMARY KEY

=over 4

=item * L</id>

=back

=cut

__PACKAGE__->set_primary_key("id");


# Created by DBIx::Class::Schema::Loader v0.07046 @ 2016-11-18 22:45:19
# DO NOT MODIFY THIS OR ANYTHING ABOVE! md5sum:1osLwbS/Nzx/0LXJcCmZcA


# You can replace this text with custom code or comments, and it will be preserved on regeneration
1;
