package user

/*
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import (
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
)

var (
	user = &tc.UserV4{
		AddressLine1:      util.Ptr("line1"),
		AddressLine2:      util.Ptr("line2"),
		ChangeLogCount:    8,
		City:              util.Ptr("city"),
		Company:           util.Ptr("company"),
		Country:           util.Ptr("country"),
		Email:             util.Ptr("test@email.com"),
		FullName:          util.Ptr("Testy Mctestface"),
		GID:               nil,
		ID:                util.Ptr(1),
		LastAuthenticated: util.Ptr(time.Now()),
		LastUpdated:       time.Now(),
		LocalPassword:     util.Ptr("pw"),
		NewUser:           false,
		PhoneNumber:       util.Ptr("999-999-9999"),
		PostalCode:        util.Ptr("11111-1111"),
		PublicSSHKey:      nil,
		RegistrationSent:  util.Ptr(time.Now()),
		Role:              "role",
		StateOrProvince:   util.Ptr("state"),
		Tenant:            nil,
		TenantID:          0,
		Token:             nil,
		UCDN:              "",
		UID:               nil,
		Username:          "testy",
	}
	oldUser = defineUser()
)

func defineUser() tc.User {
	u := tc.User{
		Username:             util.Ptr(user.Username),
		RegistrationSent:     nil,
		LocalPassword:        user.LocalPassword,
		ConfirmLocalPassword: user.LocalPassword,
		RoleName:             util.Ptr(user.Role),
	}
	u.AddressLine1 = user.AddressLine1
	u.AddressLine2 = user.AddressLine2
	u.City = user.City
	u.Company = user.Company
	u.Country = user.Country
	u.Email = user.Email
	u.FullName = user.FullName
	u.GID = user.GID
	u.ID = user.ID
	u.LastUpdated = tc.TimeNoModFromTime(user.LastUpdated)
	u.NewUser = util.Ptr(user.NewUser)
	u.PhoneNumber = user.PhoneNumber
	u.PostalCode = user.PostalCode
	u.PublicSSHKey = user.PublicSSHKey
	u.Role = nil
	u.StateOrProvince = user.StateOrProvince
	u.Tenant = user.Tenant
	u.TenantID = util.Ptr(user.TenantID)
	u.Token = user.Token
	u.UID = user.UID
	return u
}

func TestDowngrade(t *testing.T) {
	old := user.Downgrade()

	if *old.FullName != *oldUser.FullName {
		t.Fatalf("expected FullName to be equal, got %s and %s", *old.FullName, *oldUser.FullName)
	}
	if *old.NewUser != *oldUser.NewUser {
		t.Fatalf("expected NewUser to be equal, got %t and %t", *old.NewUser, *oldUser.NewUser)
	}
	if *old.RoleName != *oldUser.RoleName {
		t.Fatalf("expected RoleName to be equal, got %s and %s", *old.RoleName, *oldUser.RoleName)
	}
	if *old.TenantID != *oldUser.TenantID {
		t.Fatalf("expected TenantID to be equal, got %d and %d", *old.TenantID, *oldUser.TenantID)
	}
	if *old.Username != *oldUser.Username {
		t.Fatalf("expected Username to be equal, got %s and %s", *old.Username, *oldUser.Username)
	}
	if *old.AddressLine1 != *oldUser.AddressLine1 {
		t.Fatalf("expected AddressLine1 to be equal, got %s and %s", *old.AddressLine1, *oldUser.AddressLine1)
	}
	if *old.AddressLine2 != *oldUser.AddressLine2 {
		t.Fatalf("expected AddressLine2 to be equal, got %s and %s", *old.AddressLine2, *oldUser.AddressLine2)
	}
	if *old.City != *oldUser.City {
		t.Fatalf("expected City to be equal, got %s and %s", *old.City, *oldUser.City)
	}
	if *old.Company != *oldUser.Company {
		t.Fatalf("expected Company to be equal, got %s and %s", *old.Company, *oldUser.Company)
	}
	if old.ConfirmLocalPassword != nil {
		t.Fatalf("expected ConfirmLocalPassword to be nil, got %s", *old.ConfirmLocalPassword)
	}
	if *old.Country != *oldUser.Country {
		t.Fatalf("expected Country to be equal, got %s and %s", *old.Country, *oldUser.Country)
	}
	if *old.Email != *oldUser.Email {
		t.Fatalf("expected Email to be equal, got %s and %s", *old.Email, *oldUser.Email)
	}
	if old.GID != nil {
		t.Fatalf("expected GID to be null, got %d", *old.GID)
	}
	if *old.ID != *oldUser.ID {
		t.Fatalf("expected ID to be equal, got %d and %d", *old.ID, *oldUser.ID)
	}
	if old.LocalPassword != nil {
		t.Fatalf("expected LocalPassword to be nil, got %s", *old.LocalPassword)
	}
	if *old.PhoneNumber != *oldUser.PhoneNumber {
		t.Fatalf("expected PhoneNumber to be equal, got %s and %s", *old.PhoneNumber, *oldUser.PhoneNumber)
	}
	if old.PublicSSHKey != nil {
		t.Fatalf("expected PublicSSHKey to be null, got %s", *old.PublicSSHKey)
	}
	if *old.StateOrProvince != *oldUser.StateOrProvince {
		t.Fatalf("expected StateOrProvince to be equal, got %s and %s", *old.StateOrProvince, *oldUser.StateOrProvince)
	}
	if old.Tenant != nil {
		t.Fatalf("expected Tenant to be null, got %s", *old.Tenant)
	}
	if old.Token != nil {
		t.Fatalf("expected Token to be null, got %s", *old.Token)
	}
	if old.UID != nil {
		t.Fatalf("expected UID to be null, got %d", *old.UID)
	}
}

func TestUpgrade(t *testing.T) {
	user.LocalPassword = util.Ptr("pw")
	newUser := oldUser.Upgrade()

	if *user.FullName != *newUser.FullName {
		t.Fatalf("expected FullName to be equal, got %s and %s", *user.FullName, *newUser.FullName)
	}
	if user.NewUser != newUser.NewUser {
		t.Fatalf("expected NewUser to be equal, got %t and %t", user.NewUser, newUser.NewUser)
	}
	if user.Role != newUser.Role {
		t.Fatalf("expected Role to be equal, got %s and %s", user.Role, newUser.Role)
	}
	if user.TenantID != newUser.TenantID {
		t.Fatalf("expected TenantID to be equal, got %d and %d", user.TenantID, newUser.TenantID)
	}
	if user.Username != newUser.Username {
		t.Fatalf("expected Username to be equal, got %s and %s", user.Username, newUser.Username)
	}
	if *user.AddressLine1 != *newUser.AddressLine1 {
		t.Fatalf("expected AddressLine1 to be equal, got %s and %s", *user.AddressLine1, *newUser.AddressLine1)
	}
	if *user.AddressLine2 != *newUser.AddressLine2 {
		t.Fatalf("expected AddressLine2 to be equal, got %s and %s", *user.AddressLine2, *newUser.AddressLine2)
	}
	if *user.City != *newUser.City {
		t.Fatalf("expected City to be equal, got %s and %s", *user.City, *newUser.City)
	}
	if *user.Company != *newUser.Company {
		t.Fatalf("expected Company to be equal, got %s and %s", *user.Company, *newUser.Company)
	}
	if *user.Country != *newUser.Country {
		t.Fatalf("expected Country to be equal, got %s and %s", *user.Country, *newUser.Country)
	}
	if *user.Email != *newUser.Email {
		t.Fatalf("expected Email to be equal, got %s and %s", *user.Email, *newUser.Email)
	}
	if user.GID != nil {
		t.Fatalf("expected GID to be null, got %d", *user.GID)
	}
	if *user.ID != *newUser.ID {
		t.Fatalf("expected ID to be equal, got %d and %d", *user.ID, *newUser.ID)
	}
	if *user.LocalPassword != *newUser.LocalPassword {
		t.Fatalf("expected LocalPassword to be equal, got %s and %s", *user.LocalPassword, *newUser.LocalPassword)
	}
	if *user.PhoneNumber != *newUser.PhoneNumber {
		t.Fatalf("expected PhoneNumber to be equal, got %s and %s", *user.PhoneNumber, *newUser.PhoneNumber)
	}
	if user.PublicSSHKey != nil {
		t.Fatalf("expected PublicSSHKey to be null, got %s", *user.PublicSSHKey)
	}
	if *user.StateOrProvince != *newUser.StateOrProvince {
		t.Fatalf("expected StateOrProvince to be equal, got %s and %s", *user.StateOrProvince, *newUser.StateOrProvince)
	}
	if user.Tenant != nil {
		t.Fatalf("expected Tenant to be null, got %s", *user.Tenant)
	}
	if user.Token != nil {
		t.Fatalf("expected Token to be null, got %s", *user.Token)
	}
	if user.UID != nil {
		t.Fatalf("expected UID to be null, got %d", *user.UID)
	}
}

func TestUpgradeCurrent(t *testing.T) {
	newUser := legacyUser.Upgrade(user.RegistrationSent, user.LastAuthenticated, user.UCDN, user.ChangeLogCount)

	if *user.FullName != *newUser.FullName {
		t.Fatalf("expected FullName to be equal, got %s and %s", *user.FullName, *newUser.FullName)
	}
	if user.NewUser != newUser.NewUser {
		t.Fatalf("expected NewUser to be equal, got %t and %t", user.NewUser, newUser.NewUser)
	}
	if user.Role != newUser.Role {
		t.Fatalf("expected Role to be equal, got %s and %s", user.Role, newUser.Role)
	}
	if user.TenantID != newUser.TenantID {
		t.Fatalf("expected TenantID to be equal, got %d and %d", user.TenantID, newUser.TenantID)
	}
	if user.Username != newUser.Username {
		t.Fatalf("expected Username to be equal, got %s and %s", user.Username, newUser.Username)
	}
	if *user.AddressLine1 != *newUser.AddressLine1 {
		t.Fatalf("expected AddressLine1 to be equal, got %s and %s", *user.AddressLine1, *newUser.AddressLine1)
	}
	if *user.AddressLine2 != *newUser.AddressLine2 {
		t.Fatalf("expected AddressLine2 to be equal, got %s and %s", *user.AddressLine2, *newUser.AddressLine2)
	}
	if *user.City != *newUser.City {
		t.Fatalf("expected City to be equal, got %s and %s", *user.City, *newUser.City)
	}
	if *user.Company != *newUser.Company {
		t.Fatalf("expected Company to be equal, got %s and %s", *user.Company, *newUser.Company)
	}
	if *user.Country != *newUser.Country {
		t.Fatalf("expected Country to be equal, got %s and %s", *user.Country, *newUser.Country)
	}
	if *user.Email != *newUser.Email {
		t.Fatalf("expected Email to be equal, got %s and %s", *user.Email, *newUser.Email)
	}
	if user.GID != nil {
		t.Fatalf("expected GID to be null, got %d", *user.GID)
	}
	if *user.ID != *newUser.ID {
		t.Fatalf("expected ID to be equal, got %d and %d", *user.ID, *newUser.ID)
	}
	if newUser.LocalPassword != nil {
		t.Fatalf("expected LocalPassword to be nil, got %s", *newUser.LocalPassword)
	}
	if *user.PhoneNumber != *newUser.PhoneNumber {
		t.Fatalf("expected PhoneNumber to be equal, got %s and %s", *user.PhoneNumber, *newUser.PhoneNumber)
	}
	if user.PublicSSHKey != nil {
		t.Fatalf("expected PublicSSHKey to be null, got %s", *user.PublicSSHKey)
	}
	if *user.StateOrProvince != *newUser.StateOrProvince {
		t.Fatalf("expected StateOrProvince to be equal, got %s and %s", *user.StateOrProvince, *newUser.StateOrProvince)
	}
	if user.Tenant != nil {
		t.Fatalf("expected Tenant to be null, got %s", *user.Tenant)
	}
	if user.Token != nil {
		t.Fatalf("expected Token to be null, got %s", *user.Token)
	}
	if user.UID != nil {
		t.Fatalf("expected UID to be null, got %d", *user.UID)
	}
}

func TestDowngradeCurrent(t *testing.T) {
	old := user.ToLegacyCurrentUser(*legacyUser.Role, *legacyUser.LocalUser)

	if *old.FullName != *oldUser.FullName {
		t.Fatalf("expected FullName to be equal, got %s and %s", *old.FullName, *oldUser.FullName)
	}
	if *old.NewUser != *oldUser.NewUser {
		t.Fatalf("expected NewUser to be equal, got %t and %t", *old.NewUser, *oldUser.NewUser)
	}
	if *old.RoleName != *oldUser.RoleName {
		t.Fatalf("expected RoleName to be equal, got %s and %s", *old.RoleName, *oldUser.RoleName)
	}
	if *old.TenantID != *oldUser.TenantID {
		t.Fatalf("expected TenantID to be equal, got %d and %d", *old.TenantID, *oldUser.TenantID)
	}
	if *old.UserName != *oldUser.Username {
		t.Fatalf("expected Username to be equal, got %s and %s", *old.UserName, *oldUser.Username)
	}
	if *old.AddressLine1 != *oldUser.AddressLine1 {
		t.Fatalf("expected AddressLine1 to be equal, got %s and %s", *old.AddressLine1, *oldUser.AddressLine1)
	}
	if *old.AddressLine2 != *oldUser.AddressLine2 {
		t.Fatalf("expected AddressLine2 to be equal, got %s and %s", *old.AddressLine2, *oldUser.AddressLine2)
	}
	if *old.City != *oldUser.City {
		t.Fatalf("expected City to be equal, got %s and %s", *old.City, *oldUser.City)
	}
	if *old.Company != *oldUser.Company {
		t.Fatalf("expected Company to be equal, got %s and %s", *old.Company, *oldUser.Company)
	}
	if *old.Country != *oldUser.Country {
		t.Fatalf("expected Country to be equal, got %s and %s", *old.Country, *oldUser.Country)
	}
	if *old.Email != *oldUser.Email {
		t.Fatalf("expected Email to be equal, got %s and %s", *old.Email, *oldUser.Email)
	}
	if old.GID != nil {
		t.Fatalf("expected GID to be null, got %d", *old.GID)
	}
	if *old.ID != *oldUser.ID {
		t.Fatalf("expected ID to be equal, got %d and %d", *old.ID, *oldUser.ID)
	}
	if *old.PhoneNumber != *oldUser.PhoneNumber {
		t.Fatalf("expected PhoneNumber to be equal, got %s and %s", *old.PhoneNumber, *oldUser.PhoneNumber)
	}
	if old.PublicSSHKey != nil {
		t.Fatalf("expected PublicSSHKey to be null, got %s", *old.PublicSSHKey)
	}
	if *old.StateOrProvince != *oldUser.StateOrProvince {
		t.Fatalf("expected StateOrProvince to be equal, got %s and %s", *old.StateOrProvince, *oldUser.StateOrProvince)
	}
	if old.Tenant != nil {
		t.Fatalf("expected Tenant to be null, got %s", *old.Tenant)
	}
	if old.Token != nil {
		t.Fatalf("expected Token to be null, got %s", *old.Token)
	}
	if old.UID != nil {
		t.Fatalf("expected UID to be null, got %d", *old.UID)
	}
}
