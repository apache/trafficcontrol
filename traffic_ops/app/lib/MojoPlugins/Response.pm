package MojoPlugins::Response;
#
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
#
#

use Mojo::Base 'Mojolicious::Plugin';
use Carp qw(cluck confess);
use Data::Dumper;
use Hash::Merge qw( merge );

my $ERROR_LEVEL   = "error";
my $INFO_LEVEL    = "info";
my $WARNING_LEVEL = "warning";
my $SUCCESS_LEVEL = "success";

my $ALERTS_KEY   = "alerts";
my $LEVEL_KEY    = "level";
my $TEXT_KEY     = "text";
my $STATUS_KEY   = "status";
my $JSON_KEY     = "json";
my $RESPONSE_KEY = "response";
my $LIMIT_KEY    = "limit";
my $SIZE_KEY     = "size";
my $ORDERBY_KEY  = "orderby";
my $PAGE_KEY     = "page";
my $INFO_KEY     = "supplemental";

sub register {
	my ( $self, $app, $conf ) = @_;

	# Success (200) - With a JSON response and an optional success message
	$app->renderer->add_helper(
		success => sub {
			my $self    = shift || confess("Call on an instance of MojoPlugins::Response");
			my $body    = shift || confess("Please supply a response body hash.");

			# optional args
			my $message = shift;
			my $orderby = shift;
			my $limit   = shift;
			my $size    = shift;
			my $page    = shift;

			my $response_body = {
				$RESPONSE_KEY => $body
			};

			if ( defined($message) ) {
				$response_body = merge( $response_body, { $ALERTS_KEY   => [ { $LEVEL_KEY => $SUCCESS_LEVEL, $TEXT_KEY => $message } ] } );
			}
			if ( defined($orderby) ) {
				$response_body = merge( $response_body, { $ORDERBY_KEY => $orderby } );
			}
			if ( defined($limit) ) {
				$response_body = merge( $response_body, { $LIMIT_KEY => $limit } );
			}
			if ( defined($page) ) {
				$response_body = merge( $response_body, { $PAGE_KEY => $page } );
			}
			if ( defined($size) ) {
				$response_body = merge( $response_body, { $SIZE_KEY => $size } );
			}

			return $self->render( $STATUS_KEY => 200, $JSON_KEY => $response_body );
		}
	);

	# Success (200) - a JSON message response
	$app->renderer->add_helper(
		success_message => sub {
			my $self           = shift || confess("Call on an instance of MojoPlugins::Response");
			my $alert_messages = shift || confess("Please supply a response message text string.");
			my $info           = shift;

			my $response_body = { $ALERTS_KEY => [ { $LEVEL_KEY => $SUCCESS_LEVEL, $TEXT_KEY => $alert_messages } ] };
			if ( defined($info) ) {
				$response_body = {
					$ALERTS_KEY => [ { $LEVEL_KEY => $SUCCESS_LEVEL, $TEXT_KEY => $alert_messages } ],
					$INFO_KEY   => $info
				};
			}
			return $self->render( $STATUS_KEY => 200, $JSON_KEY => $response_body );
		}
	);

	# No Content (204)
	$app->renderer->add_helper(
		no_content => sub {
			my $self = shift || confess("Call on an instance of MojoPlugins::Response");

			my $response_body = { $ALERTS_KEY => [ { $LEVEL_KEY => $SUCCESS_LEVEL, $TEXT_KEY => "No Content" } ] };
			return $self->render( $STATUS_KEY => 204, $JSON_KEY => $response_body );
		}
	);

	# Alerts (400)
	$app->renderer->add_helper(
		alert => sub {
			my $self   = shift || confess("Call on an instance of MojoPlugins::Response");
			my $alerts = shift || confess("Please supply a string or an alerts hash like { 'Error #1: ' => 'Error message' }");

			my $builder ||= MojoPlugins::Response::Builder->new( $self, @_ );
			my @alerts_response = $builder->build_alerts($alerts);

			return $self->render( $STATUS_KEY => 400, $JSON_KEY => { $ALERTS_KEY => \@alerts_response } );
		}
	);

	# Success (200) - With a JSON response and a deprecated message
	$app->renderer->add_helper(
		success_deprecate => sub {
			my $self    = shift || confess("Call on an instance of MojoPlugins::Response");
			my $data    = shift || confess("Please supply a response body hash.");

			my $builder ||= MojoPlugins::Response::Builder->new($self, @_);
			my @alerts_response = ({$LEVEL_KEY => $WARNING_LEVEL, $TEXT_KEY => "This endpoint is deprecated"});

			return $self->render( $STATUS_KEY => 200, $JSON_KEY => { $ALERTS_KEY => \@alerts_response, $RESPONSE_KEY => $data } );
		}
	);

	$app->renderer->add_helper(
		deprecation => sub {
			my $self = shift || confess("Call on an instance of MojoPlugins::Response");
			my $code = shift || confess("Please supply a response code e.g. 400");
			my $alternative = shift || confess("Please supply an alternative handler, like 'PUT /api/1.4/user/current'");
			my $response_object = shift;

			my $builder ||= MojoPlugins::Response::Builder->new($self, @_);
			my @alerts_response = ({$LEVEL_KEY => $WARNING_LEVEL, $TEXT_KEY => "This endpoint is deprecated, please use '" . $alternative . "' instead"});

			if (defined($response_object)) {
				return $self->render( $STATUS_KEY => $code, $JSON_KEY => { $ALERTS_KEY => \@alerts_response, $RESPONSE_KEY => $response_object } );
			} else {
				return $self->render( $STATUS_KEY => $code, $JSON_KEY => { $ALERTS_KEY => \@alerts_response } );
			}
		}
	);

	$app->renderer->add_helper(
		with_deprecation => sub {
			my $self = shift || confess("Call on an instance of MojoPlugins::Response");
			my $alert = shift || confess("Please supply an alert string");
			my $level = shift || confess("Please supply an alert level such as 'error' or 'warning'");
			my $code = shift || confess("Please supply a response code e.g. 400");
			my $alternative = shift || confess("Please supply an alternative handler, like 'PUT /api/1.4/user/current'");
			my $response_object = shift;

			my $builder ||= MojoPlugins::Response::Builder->new($self, @_);
			my @alerts_response = ({$LEVEL_KEY => $level, $TEXT_KEY => $alert}, {$LEVEL_KEY => $WARNING_LEVEL, $TEXT_KEY => "This endpoint is deprecated, please use '" . $alternative . "' instead"});

			if (defined($response_object)) {
				return $self->render( $STATUS_KEY => $code, $JSON_KEY => { $ALERTS_KEY => \@alerts_response, $RESPONSE_KEY => $response_object } );
			} else {
				return $self->render( $STATUS_KEY => $code, $JSON_KEY => { $ALERTS_KEY => \@alerts_response } );
			}
		}
	);

	$app->renderer->add_helper(
		deprecation_with_no_alternative => sub {
			my $self = shift || confess("Call on an instance of MojoPlugins::Response");
			my $code = shift || confess("Please supply a response code e.g. 400");
			my $response_object = shift;

			my $builder ||= MojoPlugins::Response::Builder->new($self, @_);
			my @alerts_response = ({$LEVEL_KEY => $WARNING_LEVEL, $TEXT_KEY => "This endpoint and its functionality is deprecated, and will be removed in the future"});

			if (defined($response_object)) {
				return $self->render( $STATUS_KEY => $code, $JSON_KEY => { $ALERTS_KEY => \@alerts_response, $RESPONSE_KEY => $response_object } );
			} else {
				return $self->render( $STATUS_KEY => $code, $JSON_KEY => { $ALERTS_KEY => \@alerts_response } );
			}
		}
	);

	$app->renderer->add_helper(
		with_deprecation_with_no_alternative => sub {
			my $self = shift || confess("Call on an instance of MojoPlugins::Response");
			my $alert = shift || confess("Please supply an alert string");
			my $level = shift || confess("Please supply an alert level such as 'error' or 'warning'");
			my $code = shift || confess("Please supply a response code e.g. 400");
			my $response_object = shift;

			my $builder ||= MojoPlugins::Response::Builder->new($self, @_);
			my @alerts_response = ({$LEVEL_KEY => $level, $TEXT_KEY => $alert}, {$LEVEL_KEY => $WARNING_LEVEL, $TEXT_KEY => "This endpoint and its functionality is deprecated, and will be removed in the future"});

			if (defined($response_object)) {
				return $self->render( $STATUS_KEY => $code, $JSON_KEY => { $ALERTS_KEY => \@alerts_response, $RESPONSE_KEY => $response_object } );
			} else {
				return $self->render( $STATUS_KEY => $code, $JSON_KEY => { $ALERTS_KEY => \@alerts_response } );
			}
		}
	);

	$app->renderer->add_helper(
		with_deprecation_with_custom_message => sub {
			my $self = shift || confess("Call on an instance of MojoPlugins::Response");
			my $alert = shift || confess("Please supply an alert string");
			my $level = shift || confess("Please supply an alert level such as 'error' or 'warning'");
			my $code = shift || confess("Please supply a response code e.g. 400");
			my $custom_message = shift || confess("Please supply a custom deprecation message");
			my $response_object = shift;

			my $builder ||= MojoPlugins::Response::Builder->new($self, @_);
			my @alerts_response = ({$LEVEL_KEY => $level, $TEXT_KEY => $alert}, {$LEVEL_KEY => $WARNING_LEVEL, $TEXT_KEY => $custom_message});

			if (defined($response_object)) {
				return $self->render( $STATUS_KEY => $code, $JSON_KEY => { $ALERTS_KEY => \@alerts_response, $RESPONSE_KEY => $response_object } );
			} else {
				return $self->render( $STATUS_KEY => $code, $JSON_KEY => { $ALERTS_KEY => \@alerts_response } );
			}
		}
	);

	# Alerts (500)
	$app->renderer->add_helper(
		internal_server_error => sub {
			my $self   = shift || confess("Call on an instance of MojoPlugins::Response");
			my $alerts = shift || confess("Please supply a string or an alerts hash like { 'Error #1: ' => 'Error message' }");

			my $builder ||= MojoPlugins::Response::Builder->new( $self, @_ );
			my @alerts_response = $builder->build_alerts($alerts);

			return $self->render( $STATUS_KEY => 500, $JSON_KEY => { $ALERTS_KEY => \@alerts_response } );
		}
	);

	# Unauthorized (401)
	$app->renderer->add_helper(
		unauthorized => sub {
			my $self = shift || confess("Call on an instance of MojoPlugins::Response");

			my $response_body =
				{ $ALERTS_KEY => [ { $LEVEL_KEY => $ERROR_LEVEL, $TEXT_KEY => "Unauthorized, please log in." } ] };
			return $self->render( $STATUS_KEY => 401, $JSON_KEY => $response_body );
		}
	);

	# Invalid Username or Password (401)
	$app->renderer->add_helper(
		invalid_username_or_password => sub {
			my $self = shift || confess("Call on an instance of MojoPlugins::Response");

			my $response_body =
				{ $ALERTS_KEY => [ { $LEVEL_KEY => $ERROR_LEVEL, $TEXT_KEY => "Invalid username or password." } ] };
			return $self->render( $STATUS_KEY => 401, $JSON_KEY => $response_body );
		}
	);

	# Invalid token (401)
	$app->renderer->add_helper(
		invalid_token => sub {
			my $self = shift || confess("Call on an instance of MojoPlugins::Response");

			my $response_body = { $ALERTS_KEY => [ { $LEVEL_KEY => $ERROR_LEVEL, $TEXT_KEY => "Invalid token. Please contact your administrator." } ] };
			return $self->render( $STATUS_KEY => 401, $JSON_KEY => $response_body );
		}
	);

	# Forbidden (403)
	$app->renderer->add_helper(
		forbidden => sub {
			my $self = shift || confess("Call on an instance of MojoPlugins::Response");
			my $message = shift || "Forbidden";

			my $response_body = { $ALERTS_KEY => [ { $LEVEL_KEY => $ERROR_LEVEL, $TEXT_KEY => $message } ] };
			return $self->render( $STATUS_KEY => 403, $JSON_KEY => $response_body );
		}
	);

	# Not Found (404)
	$app->renderer->add_helper(
		not_found => sub {
			my $self = shift || confess("Call on an instance of MojoPlugins::Response");

			my $response_body = { $ALERTS_KEY => [ { $LEVEL_KEY => $ERROR_LEVEL, $TEXT_KEY => "Resource not found." } ] };
			return $self->render( $STATUS_KEY => 404, $JSON_KEY => $response_body );
		}
	);

	# Deprecate will insert an 'info' message for old APIs
	$app->renderer->add_helper(
		deprecate => sub {
			my $self = shift;
			my $data = shift;

			# this parameter allows the ability to "append" to the info or "overwrite" keys from the defaults"
			my $info = shift;

			my $info_details = merge( $info, { deprecated => 'true', message => 'Expires in version 1.2', "api_doc" => '/api/1.1/docs' } );

			my @response = unshift( @$data, { info => $info_details } );

			return \@response;
		}
	);

}

package MojoPlugins::Response::Builder;

use Mojo::Base -strict;
use Scalar::Util;
use Carp ();
use Validate::Tiny;

sub new {
	my $class = shift;
	my ( $c, $object ) = @_;
	my $self = bless {
		c       => $c,
		object  => $object,
		checks  => [],
		filters => []
	}, $class;

	Scalar::Util::weaken $self->{c};
	$self;
}

# Build the Alerts response
sub build_alerts {
	my $self   = shift;
	my $result = shift;

	my @alerts;
	if ( ref($result) eq 'HASH' ) {
		my %response = %{$result};
		foreach my $msg_key ( keys %response ) {
			my %alert;
			if ( defined( $response{$msg_key} ) ) {
				my $alert_text = $msg_key . " " . $response{$msg_key};
				%alert = ( $LEVEL_KEY => $ERROR_LEVEL, $TEXT_KEY => $alert_text );
				push( @alerts, \%alert );
			}
		}
	}

	# If no key/value pair is passed just push out the error message as defined.
	else {
		my %alert = ( $LEVEL_KEY => $ERROR_LEVEL, $TEXT_KEY => $result );
		push( @alerts, \%alert );
	}
	return @alerts;
}

1;
