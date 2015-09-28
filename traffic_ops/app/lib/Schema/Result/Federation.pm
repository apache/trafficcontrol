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

=head2 federation_resolver_id

  data_type: 'integer'
  is_foreign_key: 1
  is_nullable: 0

=head2 name

  data_type: 'varchar'
  is_nullable: 0
  size: 1024

=head2 description

  data_type: 'varchar'
  is_nullable: 1
  size: 1024

=head2 cname

  data_type: 'varchar'
  is_nullable: 0
  size: 1024

=head2 ttl

  data_type: 'integer'
  is_nullable: 0

=head2 type

  data_type: 'integer'
  is_foreign_key: 1
  is_nullable: 0

=head2 last_updated

  data_type: 'timestamp'
  datetime_undef_if_invalid: 1
  default_value: current_timestamp
  is_nullable: 0

=cut

__PACKAGE__->add_columns(
  "id",
  { data_type => "integer", is_auto_increment => 1, is_nullable => 0 },
  "federation_resolver_id",
  { data_type => "integer", is_foreign_key => 1, is_nullable => 0 },
  "name",
  { data_type => "varchar", is_nullable => 0, size => 1024 },
  "description",
  { data_type => "varchar", is_nullable => 1, size => 1024 },
  "cname",
  { data_type => "varchar", is_nullable => 0, size => 1024 },
  "ttl",
  { data_type => "integer", is_nullable => 0 },
  "type",
  { data_type => "integer", is_foreign_key => 1, is_nullable => 0 },
  "last_updated",
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

=item * L</type>

=back

=cut

__PACKAGE__->set_primary_key("id", "type");

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

=head2 federation_resolver

Type: belongs_to

Related object: L<Schema::Result::FederationResolver>

=cut

__PACKAGE__->belongs_to(
  "federation_resolver",
  "Schema::Result::FederationResolver",
  { id => "federation_resolver_id" },
  { is_deferrable => 1, on_delete => "CASCADE", on_update => "CASCADE" },
);

=head2 type

Type: belongs_to

Related object: L<Schema::Result::Type>

=cut

__PACKAGE__->belongs_to(
  "type",
  "Schema::Result::Type",
  { id => "type" },
  { is_deferrable => 1, on_delete => "NO ACTION", on_update => "NO ACTION" },
);


# Created by DBIx::Class::Schema::Loader v0.07042 @ 2015-09-28 13:05:02
# DO NOT MODIFY THIS OR ANYTHING ABOVE! md5sum:NilKXY+sxtZI83Bfg/RdXw


# You can replace this text with custom code or comments, and it will be preserved on regeneration
1;
