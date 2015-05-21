use utf8;
package Schema::Result::Hwinfo;

# Created by DBIx::Class::Schema::Loader
# DO NOT MODIFY THE FIRST PART OF THIS FILE

=head1 NAME

Schema::Result::Hwinfo

=cut

use strict;
use warnings;

use base 'DBIx::Class::Core';

=head1 TABLE: C<hwinfo>

=cut

__PACKAGE__->table("hwinfo");

=head1 ACCESSORS

=head2 id

  data_type: 'integer'
  is_auto_increment: 1
  is_nullable: 0

=head2 serverid

  data_type: 'integer'
  is_foreign_key: 1
  is_nullable: 0

=head2 description

  data_type: 'varchar'
  is_nullable: 0
  size: 256

=head2 val

  data_type: 'varchar'
  is_nullable: 0
  size: 256

=head2 last_updated

  data_type: 'timestamp'
  datetime_undef_if_invalid: 1
  default_value: current_timestamp
  is_nullable: 1

=cut

__PACKAGE__->add_columns(
  "id",
  { data_type => "integer", is_auto_increment => 1, is_nullable => 0 },
  "serverid",
  { data_type => "integer", is_foreign_key => 1, is_nullable => 0 },
  "description",
  { data_type => "varchar", is_nullable => 0, size => 256 },
  "val",
  { data_type => "varchar", is_nullable => 0, size => 256 },
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

=back

=cut

__PACKAGE__->set_primary_key("id");

=head1 UNIQUE CONSTRAINTS

=head2 C<serverid>

=over 4

=item * L</serverid>

=item * L</description>

=back

=cut

__PACKAGE__->add_unique_constraint("serverid", ["serverid", "description"]);

=head1 RELATIONS

=head2 serverid

Type: belongs_to

Related object: L<Schema::Result::Server>

=cut

__PACKAGE__->belongs_to(
  "serverid",
  "Schema::Result::Server",
  { id => "serverid" },
  { is_deferrable => 1, on_delete => "CASCADE", on_update => "NO ACTION" },
);


# Created by DBIx::Class::Schema::Loader v0.07043 @ 2015-05-21 13:27:11
# DO NOT MODIFY THIS OR ANYTHING ABOVE! md5sum:Oq3GIhYgP8bbIBTlZvetMQ


# You can replace this text with custom code or comments, and it will be preserved on regeneration
1;
