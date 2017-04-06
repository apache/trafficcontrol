use utf8;
package Schema::Result::Capability;

# Created by DBIx::Class::Schema::Loader
# DO NOT MODIFY THE FIRST PART OF THIS FILE

=head1 NAME

Schema::Result::Capability

=cut

use strict;
use warnings;

use base 'DBIx::Class::Core';

=head1 TABLE: C<capability>

=cut

__PACKAGE__->table("capability");

=head1 ACCESSORS

=head2 name

  data_type: 'text'
  is_nullable: 0

=head2 description

  data_type: 'text'
  is_nullable: 1

=head2 last_updated

  data_type: 'timestamp with time zone'
  default_value: current_timestamp
  is_nullable: 1
  original: {default_value => \"now()"}

=cut

__PACKAGE__->add_columns(
  "name",
  { data_type => "text", is_nullable => 0 },
  "description",
  { data_type => "text", is_nullable => 1 },
  "last_updated",
  {
    data_type     => "timestamp with time zone",
    default_value => \"current_timestamp",
    is_nullable   => 1,
    original      => { default_value => \"now()" },
  },
);

=head1 PRIMARY KEY

=over 4

=item * L</name>

=back

=cut

__PACKAGE__->set_primary_key("name");

=head1 RELATIONS

=head2 api_capabilities

Type: has_many

Related object: L<Schema::Result::ApiCapability>

=cut

__PACKAGE__->has_many(
  "api_capabilities",
  "Schema::Result::ApiCapability",
  { "foreign.capability" => "self.name" },
  { cascade_copy => 0, cascade_delete => 0 },
);

=head2 role_capabilities

Type: has_many

Related object: L<Schema::Result::RoleCapability>

=cut

__PACKAGE__->has_many(
  "role_capabilities",
  "Schema::Result::RoleCapability",
  { "foreign.cap_name" => "self.name" },
  { cascade_copy => 0, cascade_delete => 0 },
);


# Created by DBIx::Class::Schema::Loader v0.07046 @ 2017-04-01 22:22:35
# DO NOT MODIFY THIS OR ANYTHING ABOVE! md5sum:34+RZwrrOVdouhv+bD2V/Q


# You can replace this text with custom code or comments, and it will be preserved on regeneration
1;
