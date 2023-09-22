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
	"context"
	"errors"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/test"

	"github.com/jmoiron/sqlx"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var (
	legacyUser = defineUserCurrent()
)

func defineUserCurrent() *tc.UserCurrent {
	u := tc.UserCurrent{
		UserName:  util.Ptr("testy"),
		LocalUser: util.Ptr(true),
		RoleName:  util.Ptr("role"),
	}
	u.AddressLine1 = util.Ptr("line1")
	u.AddressLine2 = util.Ptr("line2")
	u.City = util.Ptr("city")
	u.Company = util.Ptr("company")
	u.Country = util.Ptr("country")
	u.Email = util.Ptr("test@email.com")
	u.FullName = util.Ptr("Testy Mctestface")
	u.GID = nil
	u.ID = util.Ptr(1)
	u.LastUpdated = nil
	u.PhoneNumber = util.Ptr("999-999-9999")
	u.PostalCode = util.Ptr("11111-1111")
	u.PublicSSHKey = nil
	u.Role = util.Ptr(1)
	u.StateOrProvince = util.Ptr("state")
	u.Tenant = nil
	u.TenantID = util.Ptr(0)
	u.Token = nil
	u.UID = nil

	return &u
}
func addUserRow(rows *sqlmock.Rows, users ...*tc.UserV4) {
	if rows == nil {
		return
	}
	for _, addUser := range users {
		if addUser == nil {
			continue
		}
		rows.AddRow(
			user.AddressLine1,
			user.AddressLine2,
			user.ChangeLogCount,
			user.City,
			user.Company,
			user.Country,
			user.Email,
			user.FullName,
			user.GID,
			user.ID,
			user.LastAuthenticated,
			user.LastUpdated,
			user.NewUser,
			user.PhoneNumber,
			user.PostalCode,
			user.PublicSSHKey,
			user.RegistrationSent,
			user.Role,
			user.StateOrProvince,
			user.Tenant,
			user.TenantID,
			user.UCDN,
			user.UID,
			user.Username,
		)
	}
}

func addLegacyUserRow(rows *sqlmock.Rows, users ...*tc.UserCurrent) {
	if rows == nil {
		return
	}
	for _, addUser := range users {
		if addUser == nil {
			continue
		}
		rows.AddRow(
			legacyUser.AddressLine1,
			legacyUser.AddressLine2,
			legacyUser.City,
			legacyUser.Company,
			legacyUser.Country,
			legacyUser.Email,
			legacyUser.FullName,
			legacyUser.GID,
			legacyUser.ID,
			legacyUser.LastUpdated,
			legacyUser.LocalUser,
			legacyUser.NewUser,
			legacyUser.PhoneNumber,
			legacyUser.PostalCode,
			legacyUser.PublicSSHKey,
			legacyUser.Role,
			legacyUser.RoleName,
			legacyUser.StateOrProvince,
			legacyUser.Tenant,
			legacyUser.TenantID,
			legacyUser.UID,
			legacyUser.UserName,
		)
	}
}

func TestUser(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	cols := test.ColsFromStructByTagExclude("db", tc.UserV4{}, []string{"local_passwd", "token"})
	cols = test.InsertAtStr(cols, map[string][]string{
		"full_name": {
			"gid",
		},
		"state_or_province": {
			"tenant",
		},
		"tenant_id": {
			"ucdn",
			"uid",
		},
	})
	getUserRows := sqlmock.NewRows(cols)
	addUserRow(getUserRows, user)
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT .*").WillReturnError(errors.New("bad"))

	dbCtx, cancelTx := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancelTx()
	tx, err := db.BeginTxx(dbCtx, nil)
	if err != nil {
		t.Fatalf("creating transaction: %v", err)
	}

	_, err = getUser(tx.Tx, *user.ID)
	if err == nil {
		t.Fatal("expected error but got none")
	}

	mock.ExpectQuery("SELECT .*").WillReturnRows(getUserRows)

	newUser, err := getUser(tx.Tx, *user.ID)
	if err != nil {
		t.Fatalf("unable to get user: %s", err)
	}

	newUser.LocalPassword = user.LocalPassword
	if newUser != *user {
		t.Fatal("returned user is not the same as db user")
	}

	updateUserRows := sqlmock.NewRows(cols)
	addUserRow(updateUserRows, user)
	mock.ExpectQuery("UPDATE .* RETURNING .*").WillReturnRows(updateUserRows)
	mock.ExpectExec("UPDATE .*").WithArgs(*user.LocalPassword, *user.ID).WillReturnResult(sqlmock.NewResult(int64(*user.ID), 1))

	err = updateUser(user, tx.Tx, true)
	if err != nil {
		t.Fatalf("unable to update user: %s", err)
	}
}

func TestLegacyUser(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	cols := []string{
		"address_line1",
		"address_line2",
		"city",
		"company",
		"country",
		"email",
		"full_name",
		"gid",
		"id",
		"new_user",
		"phone_number",
		"postal_code",
		"public_ssh_key",
		"role",
		"role_name",
		"state_or_province",
		"tenant",
		"tenant_id",
		"uid",
		"username",
		"last_updated",
		"local_passwd",
	}
	getUserRows := sqlmock.NewRows(cols)
	addLegacyUserRow(getUserRows, legacyUser)
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT .*").WillReturnError(errors.New("bad"))

	dbCtx, cancelTx := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancelTx()
	tx, err := db.BeginTxx(dbCtx, nil)
	if err != nil {
		t.Fatalf("creating transaction: %v", err)
	}

	_, err = getLegacyUser(tx.Tx, *legacyUser.ID)
	if err == nil {
		t.Fatal("expected error but got none")
	}

	mock.ExpectQuery("SELECT .*").WillReturnRows(getUserRows)

	newUser, err := getLegacyUser(tx.Tx, *legacyUser.ID)
	if err != nil {
		t.Fatalf("unable to get user: %s", err)
	}

	if newUser != *legacyUser {
		t.Fatal("returned user is not the same as db user")
	}
}
