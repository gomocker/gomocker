// Code generated by gomocker. DO NOT EDIT.
// versions:
//     gomocker v1.0.0

package main

import "errors"

// Impls
type (
	databaseImpl struct{ behavior DatabaseBehavior }
)

var ErrMock = errors.New("")

// Check
var (
	_ Database = &databaseImpl{}
)

// NewDatabase creates mocked implementation of Database interface
func NewDatabaseMock(behavior DatabaseBehavior) Database {
	return &databaseImpl{
		behavior: behavior,
	}
}

type DatabaseBehavior struct {
	GetPasswordByLogin func(login string) (out1 string, out2 error)
}

func (aeiouy *databaseImpl) GetPasswordByLogin(login string) (out1 string, out2 error) {
	if aeiouy.behavior.GetPasswordByLogin != nil {
		return aeiouy.behavior.GetPasswordByLogin(login)
	}
	return
}