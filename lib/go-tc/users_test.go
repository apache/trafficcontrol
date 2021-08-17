package tc

/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

import (
	"testing"
	"time"
)

func TestUser_UpgradeFromLegacyUser(t *testing.T) {
	addressLine1 := "Address Line 1"
	addressLine2 := "Address Line 2"
	city := "City"
	company := "Company"
	confirmLocalPassword := "Confirm LocalPasswd"
	country := "Country"
	email := "em@i.l"
	fullName := "Full Name"
	gid := 1
	id := 2
	lastUpdated := NewTimeNoMod()
	localPassword := "LocalPasswd"
	newUser := true
	phoneNumber := "555-555-5555"
	postalCode := "55555"
	publicSSHKey := "Public SSH Key"
	registrationSent := NewTimeNoMod()
	role := "Role Name"
	stateOrProvince := "State or Province"
	tenant := "Tenant"
	tenantID := 3
	token := "Token"
	uid := 4
	username := "Username"

	var user User
	user.AddressLine1 = &addressLine1
	user.AddressLine2 = &addressLine2
	user.City = &city
	user.Company = &company
	user.ConfirmLocalPassword = &confirmLocalPassword
	user.Country = &country
	user.Email = &email
	user.FullName = &fullName
	user.GID = &gid
	user.ID = &id
	user.LastUpdated = lastUpdated
	user.LocalPassword = &localPassword
	user.NewUser = &newUser
	user.PhoneNumber = &phoneNumber
	user.PostalCode = &postalCode
	user.PublicSSHKey = &publicSSHKey
	user.RegistrationSent = registrationSent
	user.Role = new(int)
	user.RoleName = &role
	user.StateOrProvince = &stateOrProvince
	user.Tenant = &tenant
	user.TenantID = &tenantID
	user.Token = &token
	user.UID = &uid
	user.Username = &username

	upgraded := user.Upgrade()
	if upgraded.AddressLine1 == nil {
		t.Error("AddressLine1 became nil after upgrade")
	} else if *upgraded.AddressLine1 != addressLine1 {
		t.Errorf("Incorrect AddressLine1 after upgrade; want: '%s', got: '%s'", addressLine1, *upgraded.AddressLine1)
	}
	if upgraded.AddressLine2 == nil {
		t.Error("AddressLine2 became nil after upgrade")
	} else if *upgraded.AddressLine2 != addressLine2 {
		t.Errorf("Incorrect AddressLine2 after upgrade; want: '%s', got: '%s'", addressLine2, *upgraded.AddressLine2)
	}
	if upgraded.City == nil {
		t.Error("City became nil after upgrade")
	} else if *upgraded.City != city {
		t.Errorf("Incorrect City after upgrade; want: '%s', got: '%s'", city, *upgraded.City)
	}
	if upgraded.Company == nil {
		t.Error("Company became nil after upgrade")
	} else if *upgraded.Company != company {
		t.Errorf("Incorrect Company after upgrade; want: '%s', got: '%s'", company, *upgraded.Company)
	}
	if upgraded.ConfirmLocalPassword == nil {
		t.Error("ConfirmLocalPassword became nil after upgrade")
	} else if *upgraded.ConfirmLocalPassword != confirmLocalPassword {
		t.Errorf("Incorrect ConfirmLocalPassword after upgrade; want: '%s', got: '%s'", confirmLocalPassword, *upgraded.ConfirmLocalPassword)
	}
	if upgraded.Country == nil {
		t.Error("Country became nil after upgrade")
	} else if *upgraded.Country != country {
		t.Errorf("Incorrect Country after upgrade; want: '%s', got: '%s'", country, *upgraded.Country)
	}
	if upgraded.Email == nil {
		t.Error("Email became nil after upgrade")
	} else if *upgraded.Email != email {
		t.Errorf("Incorrect Email after upgrade; want: '%s', got: '%s'", email, *upgraded.Email)
	}
	if upgraded.FullName == nil {
		t.Error("Fullname became nil after upgrade")
	} else if *upgraded.FullName != fullName {
		t.Errorf("Incorrect FullName after upgrade; want: '%s', got: '%s'", fullName, *upgraded.FullName)
	}
	if upgraded.GID == nil {
		t.Error("GID became nil after upgrade")
	} else if *upgraded.GID != gid {
		t.Errorf("Incorrect GID after upgrade; want: %d, got: %d", gid, *upgraded.GID)
	}
	if upgraded.ID == nil {
		t.Error("ID became nil after upgrade")
	} else if *upgraded.ID != id {
		t.Errorf("Incorrect ID after upgrade; want: %d, got: %d", id, *upgraded.ID)
	}
	if !upgraded.LastUpdated.Equal(lastUpdated.Time) {
		t.Errorf("Incorrect LastUpdated after upgrade; want: %v, got: %v", lastUpdated.Time, upgraded.LastUpdated)
	}
	if upgraded.LocalPassword == nil {
		t.Error("LocalPassword became nil after upgrade")
	} else if *upgraded.LocalPassword != localPassword {
		t.Errorf("Incorrect LocalPassword after upgrade; want: '%s', got: '%s'", localPassword, *upgraded.LocalPassword)
	}
	if upgraded.NewUser != newUser {
		t.Errorf("Incorrect NewUser after upgrade; want: %t, got: %t", newUser, upgraded.NewUser)
	}
	if upgraded.PhoneNumber == nil {
		t.Error("PhoneNumber became nil after upgrade")
	} else if *upgraded.PhoneNumber != phoneNumber {
		t.Errorf("Incorrect PhoneNumber after upgrade; want: '%s', got: '%s'", phoneNumber, *upgraded.PhoneNumber)
	}
	if upgraded.PostalCode == nil {
		t.Error("PostalCode became nil after upgrade")
	} else if *upgraded.PostalCode != postalCode {
		t.Errorf("Incorrect PostalCode after upgrade; want: '%s', got: '%s'", postalCode, *upgraded.PostalCode)
	}
	if upgraded.PublicSSHKey == nil {
		t.Error("PublicSSHKey became nil after upgrade")
	} else if *upgraded.PublicSSHKey != publicSSHKey {
		t.Errorf("Incorrect PublicSSHKey after upgrade; want: '%s', got: '%s'", publicSSHKey, *upgraded.PublicSSHKey)
	}
	if upgraded.RegistrationSent == nil {
		t.Error("RegistrationSent became nil after upgrade")
	} else if !upgraded.RegistrationSent.Equal(registrationSent.Time) {
		t.Errorf("Incorrect RegistrationSent after upgrade; want: %v, got: %v", registrationSent.Time, *upgraded.RegistrationSent)
	}
	if upgraded.Role != role {
		t.Errorf("Incorrect Role after upgrade; want: '%s', got: '%s'", role, upgraded.Role)
	}
	if upgraded.StateOrProvince == nil {
		t.Error("StateOrProvince became nil after upgrade")
	} else if *upgraded.StateOrProvince != stateOrProvince {
		t.Errorf("Incorrect StateOrProvince after upgrade; want: '%s', got: '%s'", stateOrProvince, *upgraded.StateOrProvince)
	}
	if upgraded.Tenant == nil {
		t.Error("Tenant became nil after upgrade")
	} else if *upgraded.Tenant != tenant {
		t.Errorf("Incorrect Tenant after upgrade; want: '%s', got: '%s'", tenant, *upgraded.Tenant)
	}
	if upgraded.TenantID != tenantID {
		t.Errorf("Incorrect TenantID after upgrade; want: %d, got: %d", tenantID, upgraded.TenantID)
	}
	if upgraded.Token == nil {
		t.Error("Token became nil after upgrade")
	} else if *upgraded.Token != token {
		t.Errorf("Incorrect Token after upgrade; want: '%s', got: '%s'", token, *upgraded.Token)
	}
	if upgraded.UID == nil {
		t.Error("UID became nil after upgrade")
	} else if *upgraded.UID != uid {
		t.Errorf("Incorrect UID after upgrade; want: %d, got: %d", uid, *upgraded.UID)
	}
	if upgraded.Username != username {
		t.Errorf("Incorrect Username after upgrade; want: '%s', got: '%s'", username, upgraded.Username)
	}

	user = upgraded.DowngradeToLegacyUser()
	if user.Role != nil {
		t.Errorf("Expected Role to be nil after downgrade, got: %d", *user.Role)
	}
}

func TestUserV4_Downgrade(t *testing.T) {
	addressLine1 := "Address Line 1"
	addressLine2 := "Address Line 2"
	city := "City"
	company := "Company"
	confirmLocalPassword := "Confirm LocalPasswd"
	country := "Country"
	email := "em@i.l"
	fullName := "Full Name"
	gid := 1
	id := 2
	lastUpdated := time.Now()
	localPassword := "LocalPasswd"
	newUser := true
	phoneNumber := "555-555-5555"
	postalCode := "55555"
	publicSSHKey := "Public SSH Key"
	registrationSent := time.Now().Add(time.Second)
	role := "Role Name"
	stateOrProvince := "State or Province"
	tenant := "Tenant"
	tenantID := 3
	token := "Token"
	uid := 4
	username := "Username"

	user := UserV4{
		FullName:    &fullName,
		LastUpdated: lastUpdated,
		NewUser:     newUser,
		Role:        role,
		TenantID:    tenantID,
		Username:    username,
	}
	user.AddressLine1 = &addressLine1
	user.AddressLine2 = &addressLine2
	user.City = &city
	user.Company = &company
	user.ConfirmLocalPassword = &confirmLocalPassword
	user.Country = &country
	user.Email = &email
	user.GID = &gid
	user.ID = &id
	user.LocalPassword = &localPassword
	user.PhoneNumber = &phoneNumber
	user.PostalCode = &postalCode
	user.PublicSSHKey = &publicSSHKey
	user.RegistrationSent = &registrationSent
	user.StateOrProvince = &stateOrProvince
	user.Tenant = &tenant
	user.Token = &token
	user.UID = &uid

	downgraded := user.Downgrade()
	if downgraded.AddressLine1 == nil {
		t.Error("AddressLine1 became nil after downgrade")
	} else if *downgraded.AddressLine1 != addressLine1 {
		t.Errorf("Incorrect AddressLine1 after downgrade; want: '%s', got: '%s'", addressLine1, *downgraded.AddressLine1)
	}
	if downgraded.AddressLine2 == nil {
		t.Error("AddressLine2 became nil after downgrade")
	} else if *downgraded.AddressLine2 != addressLine2 {
		t.Errorf("Incorrect AddressLine2 after downgrade; want: '%s', got: '%s'", addressLine2, *downgraded.AddressLine2)
	}
	if downgraded.City == nil {
		t.Error("City became nil after downgrade")
	} else if *downgraded.City != city {
		t.Errorf("Incorrect City after downgrade; want: '%s', got: '%s'", city, *downgraded.City)
	}
	if downgraded.Company == nil {
		t.Error("Company became nil after downgrade")
	} else if *downgraded.Company != company {
		t.Errorf("Incorrect Company after downgrade; want: '%s', got: '%s'", company, *downgraded.Company)
	}
	if downgraded.ConfirmLocalPassword == nil {
		t.Error("ConfirmLocalPassword became nil after downgrade")
	} else if *downgraded.ConfirmLocalPassword != confirmLocalPassword {
		t.Errorf("Incorrect ConfirmLocalPassword after downgrade; want: '%s', got: '%s'", confirmLocalPassword, *downgraded.ConfirmLocalPassword)
	}
	if downgraded.Country == nil {
		t.Error("Country became nil after downgrade")
	} else if *downgraded.Country != country {
		t.Errorf("Incorrect Country after downgrade; want: '%s', got: '%s'", country, *downgraded.Country)
	}
	if downgraded.Email == nil {
		t.Error("Email became nil after downgrade")
	} else if *downgraded.Email != email {
		t.Errorf("Incorrect Email after downgrade; want: '%s', got: '%s'", email, *downgraded.Email)
	}
	if downgraded.FullName == nil {
		t.Error("FullName became nil after downgrade")
	} else if *downgraded.FullName != fullName {
		t.Errorf("Incorrect FullName after downgrade; want: '%s', got: '%s'", fullName, *downgraded.FullName)
	}
	if downgraded.GID == nil {
		t.Error("GID became nil after downgrade")
	} else if *downgraded.GID != gid {
		t.Errorf("Incorrect GID after downgrade; want: %d, got: %d", gid, *downgraded.GID)
	}
	if downgraded.ID == nil {
		t.Error("ID became nil after downgrade")
	} else if *downgraded.ID != id {
		t.Errorf("Incorrect ID after downgrade; want: %d, got: %d", id, *downgraded.ID)
	}
	if downgraded.LastUpdated == nil {
		t.Error("LastUpdated became nil after downgrade")
	} else if !downgraded.LastUpdated.Time.Equal(lastUpdated) {
		t.Errorf("Incorrect LastUpdated after downgrade; want: %v, got: %v", lastUpdated, downgraded.LastUpdated.Time)
	}
	if downgraded.LocalPassword == nil {
		t.Error("LocalPassword became nil after downgrade")
	} else if *downgraded.LocalPassword != localPassword {
		t.Errorf("Incorrect LocalPassword after downgrade; want: '%s', got: '%s'", localPassword, *downgraded.LocalPassword)
	}
	if downgraded.NewUser == nil {
		t.Error("NewUser became nil after downgrade")
	} else if *downgraded.NewUser != newUser {
		t.Errorf("Incorrect NewUser after downgrade; want: %t, got: %t", newUser, *downgraded.NewUser)
	}
	if downgraded.PhoneNumber == nil {
		t.Error("PhoneNumber became nil after downgrade")
	} else if *downgraded.PhoneNumber != phoneNumber {
		t.Errorf("Incorrect PhoneNumber after downgrade; want: '%s', got: '%s'", phoneNumber, *downgraded.PhoneNumber)
	}
	if downgraded.PostalCode == nil {
		t.Error("PostalCode became nil after downgrade")
	} else if *downgraded.PostalCode != postalCode {
		t.Errorf("Incorrect PostalCode after downgrade; want: '%s', got: '%s'", postalCode, *downgraded.PostalCode)
	}
	if downgraded.PublicSSHKey == nil {
		t.Error("PublicSSHKey became nil after downgrade")
	} else if *downgraded.PublicSSHKey != publicSSHKey {
		t.Errorf("Incorrect PublicSSHKey after downgrade; want: '%s', got: '%s'", publicSSHKey, *downgraded.PublicSSHKey)
	}
	if downgraded.RegistrationSent == nil {
		t.Error("RegistrationSent became nil after downgrade")
	} else if !downgraded.RegistrationSent.Time.Equal(registrationSent) {
		t.Errorf("Incorrect RegistrationSent after downgrade; want: %v, got: %v", registrationSent, downgraded.RegistrationSent.Time)
	}
	if downgraded.RoleName == nil {
		t.Error("RoleName became nil after downgrade")
	} else if *downgraded.RoleName != role {
		t.Errorf("Incorrect RoleName after downgrade; want: '%s', got: '%s'", role, *downgraded.RoleName)
	}
	if downgraded.StateOrProvince == nil {
		t.Error("StateOrProvince became nil after downgrade")
	} else if *downgraded.StateOrProvince != stateOrProvince {
		t.Errorf("Incorrect StateOrProvince after downgrade; want: '%s', got: '%s'", stateOrProvince, *downgraded.StateOrProvince)
	}
	if downgraded.Tenant == nil {
		t.Error("Tenant became nil after downgrade")
	} else if *downgraded.Tenant != tenant {
		t.Errorf("Incorrect Tenant after downgrade; want: '%s', got: '%s'", tenant, *downgraded.Tenant)
	}
	if downgraded.TenantID == nil {
		t.Error("TenantID became nil after downgrade")
	} else if *downgraded.TenantID != tenantID {
		t.Errorf("Incorrect TenantID after downgrade; want: %d, got: %d", tenantID, *downgraded.TenantID)
	}
	if downgraded.Token == nil {
		t.Error("Token became nil after downgrade")
	} else if *downgraded.Token != token {
		t.Errorf("Incorrect Token after downgrade; want: '%s', got: '%s'", token, *downgraded.Token)
	}
	if downgraded.UID == nil {
		t.Error("UID became nil after downgrade")
	} else if *downgraded.UID != uid {
		t.Errorf("Incorrect UID after downgrade; want: %d, got: %d", uid, *downgraded.UID)
	}
	if downgraded.Username == nil {
		t.Error("Username became nil after downgrade")
	} else if *downgraded.Username != username {
		t.Errorf("Incorrect Username after downgrade; want: '%s', got: '%s'", username, *downgraded.Username)
	}
}

func TestCopyUtilities(t *testing.T) {
	var s *string
	copiedS := copyStringIfNotNil(s)
	if copiedS != nil {
		t.Errorf("Copying a nil string should've given nil, got: %s", *copiedS)
	}
	s = new(string)
	*s = "test string"
	copiedS = copyStringIfNotNil(s)
	if copiedS == nil {
		t.Errorf("Copied pointer to '%s' was nil", *s)
	} else {
		if *copiedS != *s {
			t.Errorf("Incorrectly copied string pointer; expected: '%s', got: '%s'", *s, *copiedS)
		}
		*s = "a different test string"
		if *copiedS == *s {
			t.Error("Expected copy to be 'deep' but modifying the original string changed the copy")
		}
	}

	var i *int
	copiedI := copyIntIfNotNil(i)
	if copiedI != nil {
		t.Errorf("Copying a nil int should've given nil, got: %d", *copiedI)
	}
	i = new(int)
	*i = 9000
	copiedI = copyIntIfNotNil(i)
	if copiedI == nil {
		t.Errorf("Copied pointer to %d was nil", *i)
	} else {
		if *copiedI != *i {
			t.Errorf("Incorrectly copied int pointer; expected: %d, got: %d", *i, *copiedI)
		}
		*i = 9001
		if *copiedI == *i {
			t.Error("Expected copy to be 'deep' but modifying the original int changed the copy")
		}
	}
}
