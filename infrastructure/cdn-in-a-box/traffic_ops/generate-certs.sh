#!/bin/bash
# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.

X509_CA_DEFAULT_NAME="ca"
X509_CA_DEFAULT_COUNTRY="ZZ"
X509_CA_DEFAULT_STATE="SomeState"
X509_CA_DEFAULT_CITY="SomeCity"
X509_CA_DEFAULT_COMPANY="SomeCompany"
X509_CA_DEFAULT_ORG="SomeOrganization"
X509_CA_DEFAULT_ORGUNIT="SomeOrgUnit"
X509_CA_DEFAULT_EMAIL="no-reply@some.host.test"
X509_CA_DEFAULT_DIGEST="sha256"
X509_CA_DEFAULT_DURATION_DAYS="365"
X509_CA_DEFAULT_KEYTYPE="rsa"
X509_CA_DEFAULT_KEYSIZE=4096 
X509_CA_DEFAULT_UMASK=0002

x509v3_usage()
{
	local prog="$1"

	echo "Usage: $prog [output-dir]"
}

x509v3_init()
{
	local output_dir="$1"

	X509_CA_DIR="${X509_CA_DIR:-$output_dir}"	
  X509_CA_NAME="${X509_CA_NAME:-$X509_CA_DEFAULT_NAME}"
  X509_CA_COUNTRY="${X509_CA_COUNTRY:-$X509_CA_DEFAULT_COUNTRY}"
  X509_CA_STATE="${X509_CA_STATE:-$X509_CA_DEFAULT_STATE}"
  X509_CA_CITY="${X509_CA_CITY:-$X509_CA_DEFAULT_CITY}"
  X509_CA_COMPANY="${X509_CA_COMPANY:-$X509_CA_DEFAULT_COMPANY}"
  X509_CA_ORG="${X509_CA_ORG:-$X509_CA_DEFAULT_ORG}"
  X509_CA_ORGUNIT="${X509_CA_ORGUNIT:-$X509_CA_DEFAULT_ORGUNIT}"
  X509_CA_EMAIL="${X509_CA_EMAIL:--$X509_CA_DEFAULT_EMAIL}"
  X509_CA_DIGEST="${X509_CA_DIGEST:-$X509_CA_DEFAULT_DIGEST}"
  X509_CA_DURATION_DAYS="${X509_CA_DURATION_DAYS:-$X509_CA_DEFAULT_DURATION_DAYS}"
  X509_CA_KEYTYPE="${X509_CA_KEYTYPE:-$X509_CA_DEFAULT_KEYTYPE}"
  X509_CA_KEYSIZE="${X509_CA_KEYSIZE:-$X509_CA_DEFAULT_KEYSIZE}"
  X509_CA_UMASK="${X509_CA_UMASK:-$X509_CA_DEFAULT_UMASK}"

  # Runtime
	X509_CA_CERT_FILE="$X509_CA_DIR/$X509_CA_NAME.crt"
	X509_CA_KEY_FILE="$X509_CA_DIR/$X509_CA_NAME.key"
	X509_CA_CONFIG_FILE="$X509_CA_DIR/$X509_CA_NAME.config"
	X509_CA_DB_FILE="$X509_CA_DIR/$X509_CA_NAME.db"
	X509_CA_REVOKE_FILE="$X509_CA_DIR/$X509_CA_NAME.crl"
	X509_CA_SERIAL_FILE="$X509_CA_DIR/$X509_CA_NAME.serial"
	X509_CA_ENV_FILE="$X509_CA_DIR/environment"
	X509_CA_DONE_FILE="$X509_CA_DIR/completed"
  
  # Set the Umask
  umask $X509_CA_DEFAULT_UMASK
	
  # If no X509_CA directory exists, create it
  if [ ! -d "$X509_CA_DIR" ] ; then
		mkdir -p "$X509_CA_DIR"
		x509v3_create_ca
  else 
		echo "Previous X509v3 CA Directory Exists..., skipping..."
  fi

}

x509v3_create_ca()
{
	# Use today's epoch date for the first serial number
	echo "$(date +%s)" > "$X509_CA_SERIAL_FILE"

	# Create the CA index file
 	touch "$X509_CA_DB_FILE"

	# Create the CA environment file
	touch "$X509_CA_ENV_FILE"

	local cn="$X509_CA_DEFAULT_ORG Root CA (CA)"

	local ca_config=""`
		`"[ca]\n"`
    `"default_ca = $X509_CA_DEFAULT_NAME\n\n"`
		`"[$X509_CA_DEFAULT_NAME]\n"`
		`"new_certs_dir = $X509_CA_DIR\n"`
		`"certificate = $X509_CA_CERT_FILE\n"`
		`"private_key = $X509_CA_KEY_FILE\n"`
		`"serial = $X509_CA_SERIAL_FILE\n"`
		`"database = $X509_CA_DB_FILE\n"`
		`"default_md = $X509_CA_DEFAULT_DIGEST\n"`
		`"default_days = $X509_CA_DEFAULT_DURATION_DAYS\n"`
		`"prompt = no\n"`
		`"preserve = no\n\n"`
    `"[policy_match]\n"`
    `"countryName = match\n"`
    `"stateOrProvinceName = match\n"`
    `"organizationName = match\n"`
    `"organizationalUnitName = optional\n"`
    `"commonName = supplied\n"`
    `"emailAddress = optional\n\n"`
    `"[policy_anything]\n"`
    `"countryName = optional\n"`
    `"stateOrProvinceName = optional\n"`
    `"localityName = optional\n"`
    `"organizationName = optional\n"`
    `"organizationalUnitName = optional\n"`
    `"commonName = supplied\n"`
		`"emailAddress = optional\n\n"`
    `"[req]\n"`
    `"default_bits = $X509_CA_DEFAULT_KEYSIZE\n"`
    `"distinguished_name = req_dn\n"`
    `"string_mask = nombstr\n"`
    `"x509_extensions = v3_ca\n"`
    `"[req_dn]\n"`
    `"countryName = Country Name (2 letter code)\n"`
    `"countryName_default = $X509_CA_DEFAULT_COUNTRY\n"`
    `"countryName_min = 2\n"`
    `"countryName_max = 2\n"`
    `"stateOrProvinceName = State or Province Name (full name)\n"`
    `"stateOrProvinceName_default = $X509_CA_DEFAULT_STATE\n"`
    `"localityName = Locality Name (eg, city)\n"`
    `"localityName_default = $X509_CA_DEFAULT_CITY\n"`
    `"0.organizationName = Organization Name (eg, company)\n"`
    `"0.organizationName_default = $X509_CA_DEFAULT_ORG\n"`
    `"organizationalUnitName = Organizational Unit Name (eg, section)\n"`
    `"organizationalUnitName_default = $X509_CA_DEFAULT_ORGUNIT\n"`
    `"commonName = Common Name (eg, YOUR name)\n"`
    `"commonName_max = 64\n"`
    `"emailAddress = Email Address\n"`
    `"emailAddress_max = 64\n"`
    `"emailAddress_default = $X509_CA_DEFAULT_EMAIL\n\n"`
		`"[v3_ca]\n"`
		`"basicConstraints = CA:TRUE\n"`
		`"subjectKeyIdentifier = hash\n"`
		`"keyUsage = cRLSign, keyCertSign\n"`
		`"extendedKeyUsage = serverAuth, clientAuth\n\n"

		##"req_extensions = v3_req\n"
		#"[v3_req]\n"
		#"basicConstraints = CA:FALSE\n"
    #"keyUsage = nonRepudiation, digitalSignature, keyEncipherment\n"
		#"extendedKeyUsage = serverAuth, clientAuth, codeSigning, emailProtection\n"
		#"subjectKeyIdentifier = hash\n"
		#"authorityKeyIdentifier = keyid,issuer\n

		echo "Writing CA openssl CA Config File"
	  echo -e "$ca_config" > "$X509_CA_CONFIG_FILE"

		# Create new CA Certificate and Key
 		openssl req -x509 -nodes -extensions v3_ca \
			-config "$X509_CA_CONFIG_FILE" \
      -newkey "$X509_CA_DEFAULT_KEYTYPE:$X509_CA_DEFAULT_KEYSIZE" \
			-keyout "$X509_CA_KEY_FILE" \
      -out "$X509_CA_CERT_FILE" \
  		-subj "/C=$X509_CA_DEFAULT_COUNTRY/ST=$X509_CA_DEFAULT_STATE/L=$X509_CA_DEFAULT_CITY/O=$X509_CA_DEFAULT_ORG/OU=$X509_CA_DEFAULT_ORG/CN=$cn/emailAddress=$X509_CA_DEFAULT_EMAIL/" 
		retcode=$?

		echo "CreateX509 Cert RetCode=$retcode"
	
		return $retcode
}

x509v3_gen_alt_names()
{
	local names="$1"

	[ -z "$name" ] && return 1

	local alt_names_text=""`
		`"[alt_names]\n"

	local count=1
	for n in $names 
	do 
		alt_names_text="${alt_names_text}DNS.$count = $n"
		((count+=1))
	done

	echo "$alt_names_text"

	return 0
}

x509v3_create_cert()
{
  local name="$1"
  local cn="$2"
	local alt_names="$3"

	echo "name=[$name]"
	echo "cn=[$cn]"
	echo "alt_names=[$alt_names]"
	
	# TODO: Error Checking
	
	local config_file="$X509_CA_DIR/$cn.config"
	local exten_file="$X509_CA_DIR/$cn.exten"
	local cert_file="$X509_CA_DIR/$cn.crt"
	local key_file="$X509_CA_DIR/$cn.key"
	local request_file="$X509_CA_DIR/$cn.csr"
	local env_name=$(echo -e "$name" | tr '[a-z]' '[A-Z]' | sed 's/\./_/g')

	# CN is always included in SAN list as it is required by all modern web browsers.
	cn="*.${cn}"		
  alt_names=$(x509v3_gen_alt_names "$cn $san_list") 	

	local request_config=""`
		`"[req]\n"`
		`"encrypt_key = no\n"`
		`"prompt = no\n"`
		`"utf8 = yes\n"`
		`"default_md = $X509_CA_DEFAULT_DIGEST\n"`
		`"default_bits = $X509_CERT_KEYSIZE\n"`
		`"distinguished_name = dn\n"`
		`"req_extensions = req_ext\n\n"`
		`"[dn]\n"`
		`"C = $X509_CA_DEFAULT_COUNTRY\n"`
		`"ST = $X509_CA_DEFAULT_STATE\n"`
		`"L = $X509_CA_DEFAULT_CITY\n"`
		`"O  = $X509_CA_DEFAULT_ORG\n"`
		`"CN = $cn\n\n"`
		`"[req_ext]\n"`
		`"basicConstraints=CA:FALSE\n"`
		`"subjectKeyIdentifier = hash\n"`
		`"subjectAltName=@alt_names\n\n"

	local exten_config=""`
		`"[v3_ext]\n"`
		`"basicConstraints=CA:FALSE\n"`
		`"subjectKeyIdentifier = hash\n"`
		`"subjectAltName=@alt_names\n\n"

  echo "Creating x509v3 request for cn=$cn type $type..."

	# Create OpenSSL config file this request
	echo -e "${request_config}${alt_names}" > "$config_file"
	echo -e "${exten_config}${alt_names}" > "$exten_file"

	# Create the x509 request config file 
  openssl req -new -nodes \
	 -config "$config_file" \
   -newkey "$X509_CA_DEFAULT_KEYTYPE:$X509_CA_DEFAULT_KEYSIZE" \
	 -keyout "$key_file" \
   -out "$request_file" 

  retcode=$?
	echo "Create X509 Req RetCode=$retcode"

  # Sign with the CA
  openssl ca -batch \
   -policy policy_anything \
   -config "$X509_CA_CONFIG_FILE" \
   -out "$cert_file" \
	 -extensions v3_ext \
	 -extfile "$exten_file" \
   -infiles "$request_file" 

  retcode=$?
	echo "Sign X509 Req RetCode=$retcode"

	echo "X509_${env_name}_CONFIG_FILE=\"$config_file\"" >> "$X509_CA_ENV_FILE"
	echo "X509_${env_name}_EXTEN_FILE=\"$exten_file\"" >> "$X509_CA_ENV_FILE"
	echo "X509_${env_name}_CERT_FILE=\"$cert_file\"" >> "$X509_CA_ENV_FILE"
	echo "X509_${env_name}_KEY_FILE=\"$key_file\"" >> "$X509_CA_ENV_FILE"
	echo "X509_${env_name}_REQUEST_FILE=\"$request_file\"" >> "$X509_CA_ENV_FILE"
}

x509v3_dump_env()
{
	tmp_file="$(mktemp)"
	cat "$X509_CA_ENV_FILE" > "$tmp_file"
	env | grep -E '^X509_' >> "$tmp_file"
	set | grep -E '^X509_' >> "$tmp_file"
	sort "$tmp_file" | uniq | sed 's/^/export /' > "$X509_CA_ENV_FILE"
	touch "$X509_CA_DONE_FILE"
}
