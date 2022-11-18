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

func compareIntPtrs(t *testing.T, name string, want, got *int, operation string) {
	if want == nil {
		t.Error("incorrect calling of compareStrPtrs - want must not be nil")
		return
	}

	if got == nil {
		t.Errorf("wrong %s after %s; want: %d, got: nil pointer", name, operation, *want)
	} else if want == got {
		t.Errorf("expected %s to be deeply copied, but it was a pointer to the original struct's field", name)
	} else if *want != *got {
		t.Errorf("wrong %s after %s; want: %d, got: %d", name, operation, *want, *got)
	}
}

func compareStrPtrs(t *testing.T, name string, want, got *string, operation string) {
	if want == nil {
		t.Error("incorrect calling of compareStrPtrs - want must not be nil")
		return
	}

	if got == nil {
		t.Errorf("wrong %s after %s; want: '%s', got: nil pointer", name, operation, *want)
	} else if want == got {
		t.Errorf("expected %s to be deeply copied, but it was a pointer to the original struct's field", name)
	} else if *want != *got {
		t.Errorf("wrong %s after %s; want: '%s', got: '%s'", name, operation, *want, *got)
	}
}

func TestUserV4_ToLegacyCurrentUser(t *testing.T) {
	addressLine1 := "Address Line 1"
	addressLine2 := "Address Line 2"
	city := "City"
	company := "Company"
	country := "Country"
	email := "em@i.l"
	fullName := "Full Name"
	gid := 1
	id := 2
	lastAuthenticated := time.Time{}
	lastUpdated := time.Now()
	localPassword := "LocalPasswd"
	localUser := true
	newUser := true
	phoneNumber := "555-555-5555"
	postalCode := "55555"
	publicSSHKey := "Public SSH Key"
	registrationSent, _ := time.Parse(time.RFC3339, "2000-01-02T03:04:05Z")
	role := "Role Name"
	roleID := 3
	stateOrProvince := "State or Province"
	tenant := "Tenant"
	tenantID := 4
	token := "Token"
	uid := 5
	username := "Username"

	user := UserV4{
		AddressLine1:      &addressLine1,
		AddressLine2:      &addressLine2,
		City:              &city,
		Company:           &company,
		Country:           &country,
		Email:             &email,
		FullName:          &fullName,
		GID:               &gid,
		ID:                &id,
		LastAuthenticated: &lastAuthenticated,
		LastUpdated:       lastUpdated,
		LocalPassword:     &localPassword,
		NewUser:           newUser,
		PhoneNumber:       &phoneNumber,
		PostalCode:        &postalCode,
		PublicSSHKey:      &publicSSHKey,
		RegistrationSent:  &registrationSent,
		Role:              role,
		StateOrProvince:   &stateOrProvince,
		Tenant:            &tenant,
		TenantID:          tenantID,
		Token:             &token,
		UCDN:              "UCDN",
		UID:               &uid,
		Username:          username,
	}

	currentUser := user.ToLegacyCurrentUser(roleID, localUser)
	compareStrPtrs(t, "AddressLine1", user.AddressLine1, currentUser.AddressLine1, "downgrade")
	compareStrPtrs(t, "AddressLine2", user.AddressLine2, currentUser.AddressLine2, "downgrade")
	compareStrPtrs(t, "City", user.City, currentUser.City, "downgrade")
	compareStrPtrs(t, "Company", user.Company, currentUser.Company, "downgrade")
	compareStrPtrs(t, "Country", user.Country, currentUser.Country, "downgrade")
	compareStrPtrs(t, "Email", user.Email, currentUser.Email, "downgrade")
	compareStrPtrs(t, "FullName", user.FullName, currentUser.FullName, "downgrade")
	compareIntPtrs(t, "GID", user.GID, currentUser.GID, "downgrade")
	compareIntPtrs(t, "ID", user.ID, currentUser.ID, "downgrade")
	if currentUser.LastUpdated == nil {
		t.Errorf("wrong LastUpdated after downgrade; want: '%s', got: nil pointer", lastUpdated)
	} else if !currentUser.LastUpdated.Time.Equal(lastUpdated) {
		t.Errorf("wrong LastUpdated after downgrade; want: '%s', got: '%s'", lastUpdated, currentUser.LastUpdated.Time)
	}
	if currentUser.LocalUser == nil {
		t.Errorf("wrong LocalUser after downgrade; want: %t, got: nil pointer", localUser)
	} else if *currentUser.LocalUser != localUser {
		t.Errorf("wrong LocalUser after downgrade; want: %t, got: %t", localUser, *currentUser.LocalUser)
	}
	if currentUser.NewUser == nil {
		t.Errorf("wrong NewUser after downgrade; want: %t, got: nil pointer", newUser)
	} else if *currentUser.NewUser != newUser {
		t.Errorf("wrong NewUser after downgrade; want: %t, got: %t", newUser, *currentUser.NewUser)
	}
	compareStrPtrs(t, "PhoneNumber", user.PhoneNumber, currentUser.PhoneNumber, "downgrade")
	compareStrPtrs(t, "PostalCode", user.PostalCode, currentUser.PostalCode, "downgrade")
	compareStrPtrs(t, "PublicSSHKey", user.PublicSSHKey, currentUser.PublicSSHKey, "downgrade")
	if currentUser.Role == nil {
		t.Errorf("wrong Role after downgrade; want: %d, got: nil pointer", roleID)
	} else if *currentUser.Role != roleID {
		t.Errorf("wrong Role after downgrade; want: %d, got: %d", roleID, *currentUser.Role)

	}
	if currentUser.RoleName == nil {
		t.Errorf("wrong RoleName after downgrade; want: '%s', got: nil pointer", role)
	} else if *currentUser.RoleName != role {
		t.Errorf("wrong RoleName after downgrade; want: '%s', got: '%s'", role, *currentUser.RoleName)
	}
	compareStrPtrs(t, "StateOrProvince", user.StateOrProvince, currentUser.StateOrProvince, "downgrade")
	compareStrPtrs(t, "Tenant", user.Tenant, currentUser.Tenant, "downgrade")
	if currentUser.TenantID == nil {
		t.Errorf("wrong TenantID after downgrade; want: %d, got: nil pointer", tenantID)
	} else if *currentUser.TenantID != tenantID {
		t.Errorf("wrong TenantID after downgrade; want: %d, got: %d", tenantID, *currentUser.TenantID)
	}
	compareStrPtrs(t, "Token", user.Token, currentUser.Token, "downgrade")
	compareIntPtrs(t, "UID", user.UID, currentUser.UID, "downgrade")
	if currentUser.UserName == nil {
		t.Errorf("wrong UserName after downgrade; want: '%s', got: nil pointer", username)
	} else if *currentUser.UserName != username {
		t.Errorf("wrong UserName after downgrade; want: '%s', got: '%s'", username, *currentUser.UserName)
	}
}

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
	compareStrPtrs(t, "AddressLine1", user.AddressLine1, upgraded.AddressLine1, "upgrade")
	compareStrPtrs(t, "AddressLine2", user.AddressLine2, upgraded.AddressLine2, "upgrade")
	compareStrPtrs(t, "City", user.City, upgraded.City, "upgrade")
	compareStrPtrs(t, "Company", user.Company, upgraded.Company, "upgrade")
	compareStrPtrs(t, "Country", user.Country, upgraded.Country, "upgrade")
	compareStrPtrs(t, "Email", user.Email, upgraded.Email, "upgrade")
	compareStrPtrs(t, "FullName", user.FullName, upgraded.FullName, "upgrade")
	compareIntPtrs(t, "GID", user.GID, upgraded.GID, "upgrade")
	compareIntPtrs(t, "ID", user.ID, upgraded.ID, "upgrade")
	if !upgraded.LastUpdated.Equal(lastUpdated.Time) {
		t.Errorf("Incorrect LastUpdated after upgrade; want: %v, got: %v", lastUpdated.Time, upgraded.LastUpdated)
	}
	compareStrPtrs(t, "LocalPassword", user.LocalPassword, upgraded.LocalPassword, "upgrade")
	if upgraded.NewUser != newUser {
		t.Errorf("Incorrect NewUser after upgrade; want: %t, got: %t", newUser, upgraded.NewUser)
	}
	compareStrPtrs(t, "PhoneNumber", user.PhoneNumber, upgraded.PhoneNumber, "upgrade")
	compareStrPtrs(t, "PostalCode", user.PostalCode, upgraded.PostalCode, "upgrade")
	compareStrPtrs(t, "PublicSSHKey", user.PublicSSHKey, upgraded.PublicSSHKey, "upgrade")
	if upgraded.RegistrationSent == nil {
		t.Error("RegistrationSent became nil after upgrade")
	} else if !upgraded.RegistrationSent.Equal(registrationSent.Time) {
		t.Errorf("Incorrect RegistrationSent after upgrade; want: %v, got: %v", registrationSent.Time, *upgraded.RegistrationSent)
	}
	if upgraded.Role != role {
		t.Errorf("Incorrect Role after upgrade; want: '%s', got: '%s'", role, upgraded.Role)
	}
	compareStrPtrs(t, "StateOrProvince", user.StateOrProvince, upgraded.StateOrProvince, "upgrade")
	compareStrPtrs(t, "Tenant", user.Tenant, upgraded.Tenant, "upgrade")
	if upgraded.TenantID != tenantID {
		t.Errorf("Incorrect TenantID after upgrade; want: %d, got: %d", tenantID, upgraded.TenantID)
	}
	compareStrPtrs(t, "Token", user.Token, upgraded.Token, "upgrade")
	compareIntPtrs(t, "UID", user.UID, upgraded.UID, "upgrade")
	if upgraded.Username != username {
		t.Errorf("Incorrect Username after upgrade; want: '%s', got: '%s'", username, upgraded.Username)
	}

	user = upgraded.Downgrade()
	if user.Role != nil {
		t.Errorf("Expected Role to be nil after downgrade, got: %d", *user.Role)
	}
}

func TestUserV4_Downgrade(t *testing.T) {
	addressLine1 := "Address Line 1"
	addressLine2 := "Address Line 2"
	city := "City"
	company := "Company"
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
	} else if *downgraded.ConfirmLocalPassword != localPassword {
		t.Errorf("Incorrect ConfirmLocalPassword after downgrade; want: '%s', got: '%s'", localPassword, *downgraded.ConfirmLocalPassword)
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
