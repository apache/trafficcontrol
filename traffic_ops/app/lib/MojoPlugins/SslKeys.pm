package MojoPlugins::SslKeys;
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
use MIME::Base64;
use Net::DNS;
use MIME::Base64;
use Crypt::OpenSSL::RSA;
use Crypt::OpenSSL::Bignum;
use Crypt::OpenSSL::Random;
use Net::DNS::SEC::Private;
use JSON;
my $TMP_LOCATION = "/var/tmp";

sub register {
	my ( $self, $app, $conf ) = @_;

	$app->renderer->add_helper(
		generate_ssl_keys => sub {
			my $self     = shift;
			my $record = shift;
			my $hostname = $record->{hostname};
			my $country  = $record->{country};
			my $city     = $record->{city};
			my $state    = $record->{state};
			my $org      = $record->{org};
			my $unit     = $record->{unit};;
			my $version  = $record->{version};
			my $key      = $record->{key};
			my $key_type = "ssl";
			my $response;

			#generate key and csr
			my $result = UI::Utils->exec_command(
				"openssl req -nodes -newkey rsa:2048 -keyout $TMP_LOCATION/$hostname.key -out $TMP_LOCATION/$hostname.csr -subj /C=\"$country\"/ST=\"$state\"/L=\"$city\"/O=\"$org\"/OU=\"$unit\"/CN=$hostname"
			);
			if ( $result != 0 ) {
				$response = { _rc => 400, _content => "Error creating key and csr. Result is $result" };
				return $response;
			}

			#remove passphrase
			$result = UI::Utils->exec_command( "openssl", "rsa", "-in", "$TMP_LOCATION/$hostname.key", "-out", "$TMP_LOCATION/$hostname.key" );
			if ( $result != 0 ) {
				$response = { _rc => 400, _content => "Error removing passphrase from key file. Result is $result" };
				return $response;
			}

			#generate self signed cert
			$result = UI::Utils->exec_command(
				"openssl",                     "x509",     "-req",                        "-days", "365", "-in",
				"$TMP_LOCATION/$hostname.csr", "-signkey", "$TMP_LOCATION/$hostname.key", "-out",  "$TMP_LOCATION/$hostname.crt"
			);
			if ( $result != 0 ) {
				$response = { _rc => 400, _content => "Error creating crt. Result is $result" };
				return $response;
			}

			#convert to base64
			$record->{certificate}->{key} = $self->convert_file_to_base64("$TMP_LOCATION/$hostname.key");
			$record->{certificate}->{csr} = $self->convert_file_to_base64("$TMP_LOCATION/$hostname.csr");
			$record->{certificate}->{crt} = $self->convert_file_to_base64("$TMP_LOCATION/$hostname.crt");

			#delete key files on fs
			$result = UI::Utils->exec_command("rm $TMP_LOCATION/$hostname.*");

			#store with the version provided
			$self->riak_put( $key_type, "$key-$version", encode_json($record) );

			#store as latest...there is probably a better way to do it
			#TODO DN figure out linking
			$response = $self->riak_put( $key_type, "$key-latest", encode_json($record) );

			return $response;
		}
	);
	$app->renderer->add_helper(
		convert_file_to_base64 => sub {
			my $self       = shift;
			my $input_file = shift;

			local $/ = undef;

			# $self->app->log->debug("input file is $input_file");
			open FILE, "$input_file" or die "Couldn't open file: $!";
			my $text = <FILE>;
			close FILE;

			my $encoded = encode_base64($text);

			#trim
			$encoded =~ s/\s+$//;
			return $encoded;
		}
	);
	$app->renderer->add_helper(
		add_ssl_keys_to_riak => sub {
			my $self     = shift;
			my $key_type = "ssl";
			my $record = shift;
			my $key = $record->{key};
			my $version = $record->{version};

			#convert to base64
			my $crt = encode_base64($record->{certificate}->{crt});
			my $csr = encode_base64($record->{certificate}->{csr});
			my $priv_key = encode_base64($record->{certificate}->{key});

			#trim
			$priv_key =~ s/\s+$//;
			$crt =~ s/\s+$//;
			$csr =~ s/\s+$//;

			$record->{certificate}->{crt} = $crt;
			$record->{certificate}->{csr} = $csr;
			$record->{certificate}->{key} = $priv_key;

			#store with the version provided
			my $response = $self->riak_put( $key_type, "$key-$version", encode_json($record) );

			if ( $response->{'response'}->{_rc} == 204 ) {

				#store as latest
				$response = $self->riak_put( $key_type, "$key-latest", encode_json($record) );
			}
			return $response;
		}
	);

}
1;
