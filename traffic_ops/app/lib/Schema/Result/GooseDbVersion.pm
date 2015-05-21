use utf8;
package Schema::Result::GooseDbVersion;

# Created by DBIx::Class::Schema::Loader
# DO NOT MODIFY THE FIRST PART OF THIS FILE

=head1 NAME

Schema::Result::GooseDbVersion

=cut

use strict;
use warnings;

use base 'DBIx::Class::Core';

=head1 TABLE: C<goose_db_version>

=cut

__PACKAGE__->table("goose_db_version");

=head1 ACCESSORS

=head2 id

  data_type: 'bigint'
  extra: {unsigned => 1}
  is_auto_increment: 1
  is_nullable: 0

=head2 version_id

  data_type: 'bigint'
  is_nullable: 0

=head2 is_applied

  data_type: 'tinyint'
  is_nullable: 0

=head2 tstamp

  data_type: 'timestamp'
  datetime_undef_if_invalid: 1
  default_value: current_timestamp
  is_nullable: 1

=cut

__PACKAGE__->add_columns(
  "id",
  {
    data_type => "bigint",
    extra => { unsigned => 1 },
    is_auto_increment => 1,
    is_nullable => 0,
  },
  "version_id",
  { data_type => "bigint", is_nullable => 0 },
  "is_applied",
  { data_type => "tinyint", is_nullable => 0 },
  "tstamp",
  {
    data_type => "timestamp",
    datetime_undef_if_invalid => 1,
    default_value => \"current_timestamp",
    is_nullable => 1,
  },
);

=head1 PRIMARY KEY

=over 4

=item * L</id>

=back

=cut

__PACKAGE__->set_primary_key("id");


# Created by DBIx::Class::Schema::Loader v0.07043 @ 2015-05-21 13:27:11
# DO NOT MODIFY THIS OR ANYTHING ABOVE! md5sum:utMNoRDzBw9a2hNaMxsP3A


# You can replace this text with custom code or comments, and it will be preserved on regeneration
1;
