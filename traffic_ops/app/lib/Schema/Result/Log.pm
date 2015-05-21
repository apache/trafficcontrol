use utf8;
package Schema::Result::Log;

# Created by DBIx::Class::Schema::Loader
# DO NOT MODIFY THE FIRST PART OF THIS FILE

=head1 NAME

Schema::Result::Log

=cut

use strict;
use warnings;

use base 'DBIx::Class::Core';

=head1 TABLE: C<log>

=cut

__PACKAGE__->table("log");

=head1 ACCESSORS

=head2 id

  data_type: 'integer'
  is_auto_increment: 1
  is_nullable: 0

=head2 level

  data_type: 'varchar'
  is_nullable: 1
  size: 45

=head2 message

  data_type: 'varchar'
  is_nullable: 0
  size: 1024

=head2 tm_user

  data_type: 'integer'
  is_foreign_key: 1
  is_nullable: 0

=head2 ticketnum

  data_type: 'varchar'
  is_nullable: 1
  size: 64

=head2 last_updated

  data_type: 'timestamp'
  datetime_undef_if_invalid: 1
  default_value: current_timestamp
  is_nullable: 1

=cut

__PACKAGE__->add_columns(
  "id",
  { data_type => "integer", is_auto_increment => 1, is_nullable => 0 },
  "level",
  { data_type => "varchar", is_nullable => 1, size => 45 },
  "message",
  { data_type => "varchar", is_nullable => 0, size => 1024 },
  "tm_user",
  { data_type => "integer", is_foreign_key => 1, is_nullable => 0 },
  "ticketnum",
  { data_type => "varchar", is_nullable => 1, size => 64 },
  "last_updated",
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

=item * L</tm_user>

=back

=cut

__PACKAGE__->set_primary_key("id", "tm_user");

=head1 RELATIONS

=head2 tm_user

Type: belongs_to

Related object: L<Schema::Result::TmUser>

=cut

__PACKAGE__->belongs_to(
  "tm_user",
  "Schema::Result::TmUser",
  { id => "tm_user" },
  { is_deferrable => 1, on_delete => "NO ACTION", on_update => "NO ACTION" },
);


# Created by DBIx::Class::Schema::Loader v0.07043 @ 2015-05-21 13:27:11
# DO NOT MODIFY THIS OR ANYTHING ABOVE! md5sum:XkJt7i4956YEhggYnKRF3A


# You can replace this text with custom code or comments, and it will be preserved on regeneration
1;
