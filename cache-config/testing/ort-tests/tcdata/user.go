package tcdata

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

import (
	"strings"
	"testing"
)

func (r *TCData) CreateTestUsers(t *testing.T) {
	for _, user := range r.TestData.Users {

		resp, _, err := TOSession.CreateUser(&user)
		if err != nil {
			t.Errorf("could not CREATE user: %v", err)
		}
		t.Logf("Alerts: %+v", resp.Alerts)
	}
}

// ForceDeleteTestUsers forcibly deletes the users from the db.
func (r *TCData) ForceDeleteTestUsers(t *testing.T) {

	// NOTE: Special circumstances!  This should *NOT* be done without a really good reason!
	//  Connects directly to the DB to remove users rather than going thru the client.
	//  This is required here because the DeleteUser action does not really delete users,  but disables them.
	db, err := r.OpenConnection()
	if err != nil {
		t.Error("cannot open db")
	}
	defer db.Close()

	var usernames []string
	for _, user := range r.TestData.Users {
		usernames = append(usernames, `'`+*user.Username+`'`)
	}

	// there is a constraint that prevents users from being deleted when they have a log
	q := `DELETE FROM log WHERE NOT tm_user = (SELECT id FROM tm_user WHERE username = 'admin')`
	err = execSQL(db, q)
	if err != nil {
		t.Errorf("cannot execute SQL: %s; SQL is %s", err.Error(), q)
	}

	q = `DELETE FROM tm_user WHERE username IN (` + strings.Join(usernames, ",") + `)`
	err = execSQL(db, q)
	if err != nil {
		t.Errorf("cannot execute SQL: %s; SQL is %s", err.Error(), q)
	}
}

func (r *TCData) DeleteTestUsers(t *testing.T) {
	for _, user := range r.TestData.Users {

		resp, _, err := TOSession.GetUserByUsername(*user.Username)
		if err != nil {
			t.Errorf("cannot GET user by name: %s - %v", *user.Username, err)
		}

		if resp != nil {
			respUser := resp[0]

			_, _, err := TOSession.DeleteUserByID(*respUser.ID)
			if err != nil {
				t.Errorf("cannot DELETE user by name: '%s' %v", *respUser.Username, err)
			}

			// Make sure it got deleted
			resp, _, err := TOSession.GetUserByUsername(*user.Username)
			if err != nil {
				t.Errorf("error deleting user by name: %v", err)
			}
			if len(resp) > 0 {
				t.Errorf("expected user: %s to be deleted", *user.Username)
			}
		}
	}
}
