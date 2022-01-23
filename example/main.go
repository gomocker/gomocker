package main

import (
	"errors"
	"fmt"
	"log"
	"time"
)

type Database interface {
	GetPasswordByLogin(login string) (password string, err error)
}

type realDatabase struct{}

var (
	ErrInvalidLoginData = errors.New("invalid login data")
	ErrUserNotFound     = errors.New("user not found")
)

func (*realDatabase) GetPasswordByLogin(login string) (string, error) {
	time.Sleep(5 * time.Second)
	if login != "admin" {
		return "", ErrUserNotFound
	}
	return "admin", nil
}

type Application struct {
	Database Database
}

func (a *Application) Login(login string, password string) error {
	realPassword, err := a.Database.GetPasswordByLogin(login)
	if err != nil {
		return err
	}
	if realPassword != password {
		return ErrInvalidLoginData
	}
	return nil
}

func NewApplication(db Database) *Application {
	return &Application{
		Database: db,
	}
}

func main() {
	app := NewApplication(&realDatabase{})

	if err := app.Login("admin", "admin"); err != nil {
		log.Fatalln(err)
	}

	fmt.Println("ok")
}
