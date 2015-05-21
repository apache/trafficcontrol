use utf8;
package Schema::Result::Profile;

# Created by DBIx::Class::Schema::Loader
# DO NOT MODIFY THE FIRST PART OF THIS FILE

=head1 NAME

Schema::Result::Profile

=cut

use strict;
use warnings;

use base 'DBIx::Class::Core';

=head1 TABLE: C<profile>

=cut

__PACKAGE__->table("profile");

=head1 ACCESSORS

=head2 id

  data_type: 'integer'
  is_auto_increment: 1
  is_nullable: 0

=head2 name

  data_type: 'varchar'
  is_nullable: 0
  size: 45

=head2 description

  data_type: 'varchar'
  is_nullable: 1
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
  "name",
  { data_type => "varchar", is_nullable => 0, size => 45 },
  "description",
  { data_type => "varchar", is_nullable => 1, size => 256 },
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

=head2 C<name_UNIQUE>

=over 4

=item * L</name>

=back

=cut

__PACKAGE__->add_unique_constraint("name_UNIQUE", ["name"]);

=head1 RELATIONS

=head2 deliveryservices

Type: has_many

Related object: L<Schema::Result::Deliveryservice>

=cut

__PACKAGE__->has_many(
  "deliveryservices",
  "Schema::Result::Deliveryservice",
  { "foreign.profile" => "self.id" },
  { cascade_copy => 0, cascade_delete => 0 },
);

=head2 profile_parameters

Type: has_many

Related object: L<Schema::Result::ProfileParameter>

=cut

__PACKAGE__->has_many(
  "profile_parameters",
  "Schema::Result::ProfileParameter",
  { "foreign.profile" => "self.id" },
  { cascade_copy => 0, cascade_delete => 0 },
);

=head2 servers

Type: has_many

Related object: L<Schema::Result::Server>

=cut

__PACKAGE__->has_many(
  "servers",
  "Schema::Result::Server",
  { "foreign.profile" => "self.id" },
  { cascade_copy => 0, cascade_delete => 0 },
);


# Created by DBIx::Class::Schema::Loader v0.07043 @ 2015-05-21 13:27:11
# DO NOT MODIFY THIS OR ANYTHING ABOVE! md5sum:KX3bLkMCHjiZKm7x44Lw1A


# You can replace this text with custom code or comments, and it will be preserved on regeneration
1;
