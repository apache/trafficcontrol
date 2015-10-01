use utf8;
package Schema::Result::Federation;

# Created by DBIx::Class::Schema::Loader
# DO NOT MODIFY THE FIRST PART OF THIS FILE

=head1 NAME

Schema::Result::Federation

=cut

use strict;
use warnings;

use base 'DBIx::Class::Core';

=head1 TABLE: C<federation>

=cut

__PACKAGE__->table("federation");

=head1 ACCESSORS

=head2 id

  data_type: 'integer'
  is_auto_increment: 1
  is_nullable: 0

=head2 cname

  data_type: 'varchar'
  is_nullable: 0
  size: 1024

=head2 description

  data_type: 'varchar'
  is_nullable: 1
  size: 1024

=head2 ttl

  data_type: 'integer'
  is_nullable: 0

=head2 role

  data_type: 'integer'
  is_foreign_key: 1
  is_nullable: 0

=head2 last_updated

  data_type: 'timestamp'
  datetime_undef_if_invalid: 1
  default_value: current_timestamp
  is_nullable: 1

=cut

__PACKAGE__->add_columns(
  "id",
  { data_type => "integer", is_auto_increment => 1, is_nullable => 0 },
  "cname",
  { data_type => "varchar", is_nullable => 0, size => 1024 },
  "description",
  { data_type => "varchar", is_nullable => 1, size => 1024 },
  "ttl",
  { data_type => "integer", is_nullable => 0 },
  "role",
  { data_type => "integer", is_foreign_key => 1, is_nullable => 0 },
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

=head1 RELATIONS

=head2 federation_deliveryservices

Type: has_many

Related object: L<Schema::Result::FederationDeliveryservice>

=cut

__PACKAGE__->has_many(
  "federation_deliveryservices",
  "Schema::Result::FederationDeliveryservice",
  { "foreign.federation" => "self.id" },
  { cascade_copy => 0, cascade_delete => 0 },
);

=head2 federation_federation_resolvers

Type: has_many

Related object: L<Schema::Result::FederationFederationResolver>

=cut

__PACKAGE__->has_many(
  "federation_federation_resolvers",
  "Schema::Result::FederationFederationResolver",
  { "foreign.federation" => "self.id" },
  { cascade_copy => 0, cascade_delete => 0 },
);

=head2 federation_tmusers

Type: has_many

Related object: L<Schema::Result::FederationTmuser>

=cut

__PACKAGE__->has_many(
  "federation_tmusers",
  "Schema::Result::FederationTmuser",
  { "foreign.federation" => "self.id" },
  { cascade_copy => 0, cascade_delete => 0 },
);

=head2 role

Type: belongs_to

Related object: L<Schema::Result::Role>

=cut

__PACKAGE__->belongs_to(
  "role",
  "Schema::Result::Role",
  { id => "role" },
  { is_deferrable => 1, on_delete => "NO ACTION", on_update => "NO ACTION" },
);


# Created by DBIx::Class::Schema::Loader v0.07042 @ 2015-10-01 13:37:14
# DO NOT MODIFY THIS OR ANYTHING ABOVE! md5sum:NeN+3mzLzD0MuaoCmkoLww


# You can replace this text with custom code or comments, and it will be preserved on regeneration
1;
