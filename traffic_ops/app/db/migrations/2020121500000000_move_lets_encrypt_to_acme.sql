/*
	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at
		http://www.apache.org/licenses/LICENSE-2.0
	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE IF NOT EXISTS acme_account (
  email text NOT NULL,
  private_key text NOT NULL,
  provider text NOT NULL,
  uri text NOT NULL,
  PRIMARY KEY (email, provider)
);

INSERT INTO acme_account(
	email,
	private_key,
	provider,
	uri
)
SELECT
	lets_encrypt_account.email,
	lets_encrypt_account.private_key,
	'Lets Encrypt',
	lets_encrypt_account.uri
FROM lets_encrypt_account;

DROP TABLE IF EXISTS lets_encrypt_account;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

CREATE TABLE IF NOT EXISTS lets_encrypt_account (
  email text NOT NULL,
  private_key text NOT NULL,
  uri text NOT NULL,
  PRIMARY KEY (email)
);

INSERT INTO lets_encrypt_account(
	email,
	private_key,
	uri
)
SELECT
	acme_account.email,
	acme_account.private_key,
	acme_account.uri
FROM acme_account WHERE acme_account.provider = 'Lets Encrypt';

DROP TABLE IF EXISTS acme_account;
