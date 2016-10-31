.. 
.. 
.. Licensed under the Apache License, Version 2.0 (the "License");
.. you may not use this file except in compliance with the License.
.. You may obtain a copy of the License at
.. 
..     http://www.apache.org/licenses/LICENSE-2.0
.. 
.. Unless required by applicable law or agreed to in writing, software
.. distributed under the License is distributed on an "AS IS" BASIS,
.. WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
.. See the License for the specific language governing permissions and
.. limitations under the License.
.. 

Traffic Vault
=============

Traffic Vault is a keystore used for storing the following types of information:

* SSL Certificates
	- Private Key
	- CRT
	- CSR
* DNSSEC Keys
	- Key Signing Key
		- private key
		- public key
	- Zone Signing Key
		- private key
		- public key
* URL Signing Keys

As the name suggests, Traffic Vault is meant to be a "vault" of private keys that only certain users are allowed to access.  In order to create, add, and retrieve keys a user must have admin privileges.  Keys can be created via the Traffic Ops UI, but they can only be retrieved via the Traffic Ops API.  The keystore used by Traffic Vault is `Riak <http://basho.com/riak/>`_.  Traffic ops accesses Riak via https on port 8088.  Traffic ops uses Riak's rest API with username/password authentication.  Information on the API can be found `here <http://docs.basho.com/riak/latest/dev/references/http/>`_.


