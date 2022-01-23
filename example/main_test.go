package main

import (
	"testing"
)

func TestLogin(t *testing.T) {
	type testCase struct {
		database Database
		login    string
		password string
		err      error
	}

	testCases := map[string]testCase{}

	testCases["normal case"] = testCase{
		database: NewDatabaseMock(DatabaseBehavior{
			GetPasswordByLogin: func(login string) (password string, err error) {
				return "admin", nil
			},
		}),
		login:    "admin",
		password: "admin",
		err:      nil,
	}

	testCases["user not found in database"] = testCase{
		database: NewDatabaseMock(DatabaseBehavior{
			GetPasswordByLogin: func(login string) (password string, err error) {
				return "", ErrUserNotFound
			},
		}),
		login:    "admin",
		password: "admin",
		err:      ErrUserNotFound,
	}

	testCases["invalid password"] = testCase{
		database: NewDatabaseMock(DatabaseBehavior{
			GetPasswordByLogin: func(login string) (password string, err error) {
				return "admin", ErrInvalidLoginData
			},
		}),
		login:    "admin",
		password: "invalid_password",
		err:      ErrInvalidLoginData,
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			app := NewApplication(tc.database)

			err := app.Login(tc.login, tc.password)

			if err != tc.err {
				t.Errorf("Expected %q, got %q; Test case %q", tc.err, err, name)
			}
		})
	}
}
