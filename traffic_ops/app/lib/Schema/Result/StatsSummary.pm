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

  data_type: 'integer'
  is_auto_increment: 1
  is_nullable: 0

=head2 cdn_name

  data_type: 'varchar'
  default_value: 'all'
  is_nullable: 0
  size: 255

=head2 deliveryservice_name

  data_type: 'varchar'
  is_nullable: 0
  size: 255

=head2 stat_name

  data_type: 'varchar'
  is_nullable: 0
  size: 255

=head2 stat_value

  data_type: 'float'
  is_nullable: 0

=head2 summary_timestamp

  data_type: 'timestamp'
  datetime_undef_if_invalid: 1
  default_value: current_timestamp
  is_nullable: 0

=cut

__PACKAGE__->add_columns(
  "id",
  { data_type => "integer", is_auto_increment => 1, is_nullable => 0 },
  "cdn_name",
  {
    data_type => "varchar",
    default_value => "all",
    is_nullable => 0,
    size => 255,
  },
  "deliveryservice_name",
  { data_type => "varchar", is_nullable => 0, size => 255 },
  "stat_name",
  { data_type => "varchar", is_nullable => 0, size => 255 },
  "stat_value",
  { data_type => "float", is_nullable => 0 },
  "summary_timestamp",
  {
    data_type => "timestamp",
    datetime_undef_if_invalid => 1,
    default_value => \"current_timestamp",
    is_nullable => 0,
  },
);

=head1 PRIMARY KEY

=over 4

=item * L</id>

=back

=cut

__PACKAGE__->set_primary_key("id");


# Created by DBIx::Class::Schema::Loader v0.07042 @ 2015-05-12 09:25:44
# DO NOT MODIFY THIS OR ANYTHING ABOVE! md5sum:c3qg7Ei5JXMWXjShC8fOKw


# You can replace this text with custom code or comments, and it will be preserved on regeneration
1;
