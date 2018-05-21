use utf8;
package Schema::Result::Coordinate;

# Created by DBIx::Class::Schema::Loader
# DO NOT MODIFY THE FIRST PART OF THIS FILE

=head1 NAME

Schema::Result::Coordinate

=cut

use strict;
use warnings;

use base 'DBIx::Class::Core';

=head1 TABLE: C<coordinate>

=cut

__PACKAGE__->table("coordinate");

=head1 ACCESSORS

=head2 id

  data_type: 'bigint'
  is_auto_increment: 1
  is_nullable: 0
  sequence: 'coordinate_id_seq'

=head2 name

  data_type: 'text'
  is_nullable: 0

=head2 latitude

  data_type: 'numeric'
  default_value: 0.0
  is_nullable: 0

=head2 longitude

  data_type: 'numeric'
  default_value: 0.0
  is_nullable: 0

=head2 last_updated

  data_type: 'timestamp with time zone'
  default_value: current_timestamp
  is_nullable: 0
  original: {default_value => \"now()"}

=cut

__PACKAGE__->add_columns(
  "id",
  {
    data_type         => "bigint",
    is_auto_increment => 1,
    is_nullable       => 0,
    sequence          => "coordinate_id_seq",
  },
  "name",
  { data_type => "text", is_nullable => 0 },
  "latitude",
  { data_type => "numeric", default_value => "0.0", is_nullable => 0 },
  "longitude",
  { data_type => "numeric", default_value => "0.0", is_nullable => 0 },
  "last_updated",
  {
    data_type     => "timestamp with time zone",
    default_value => \"current_timestamp",
    is_nullable   => 0,
    original      => { default_value => \"now()" },
  },
);

=head1 PRIMARY KEY

=over 4

=item * L</id>

=back

=cut

__PACKAGE__->set_primary_key("id");

=head1 UNIQUE CONSTRAINTS

=head2 C<coordinate_name_key>

=over 4

=item * L</name>

=back

=cut

__PACKAGE__->add_unique_constraint("coordinate_name_key", ["name"]);

=head1 RELATIONS

=head2 origins

Type: has_many

Related object: L<Schema::Result::Origin>

=cut

__PACKAGE__->has_many(
  "origins",
  "Schema::Result::Origin",
  { "foreign.coordinate" => "self.id" },
  { cascade_copy => 0, cascade_delete => 0 },
);


# Created by DBIx::Class::Schema::Loader v0.07042 @ 2018-05-15 16:06:00
# DO NOT MODIFY THIS OR ANYTHING ABOVE! md5sum:FZ64Zkbh+B6CECd1k/h66w


# You can replace this text with custom code or comments, and it will be preserved on regeneration
1;
