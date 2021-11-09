package main

import (
	"testing"
)

func TestLogin(t *testing.T) {
	type testCase struct {
		databaseBehavior DatabaseBehavior
		login            string
		password         string
		err              error
	}

	tt := map[string]testCase{}

	tt["normal case"] = testCase{
		databaseBehavior: DatabaseBehavior{
			GetPasswordByLogin: func(login string) (out1 string, out2 error) {
				return "admin", nil
			},
		},
		login:    "admin",
		password: "admin",
		err:      nil,
	}

	tt["invalid login; user not found"] = testCase{
		databaseBehavior: DatabaseBehavior{
			GetPasswordByLogin: func(login string) (out1 string, out2 error) {
				return "admin", ErrUserNotFound
			},
		},
		login:    "invalid_login",
		password: "admin",
		err:      ErrUserNotFound,
	}

	tt["invalid password"] = testCase{
		databaseBehavior: DatabaseBehavior{
			GetPasswordByLogin: func(login string) (out1 string, out2 error) {
				return "admin", nil
			},
		},
		login:    "admin",
		password: "invalid_password",
		err:      ErrInvalidLoginData,
	}

	for name, tc := range tt {
		app := NewApplication(NewDatabaseMock(tc.databaseBehavior))

		err := app.Login(tc.login, tc.password)

		if err != tc.err {
			t.Errorf("Expected %q, got %q; Test case %q", tc.err, err, name)
		}
	}
}
