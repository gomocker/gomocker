package main

import (
	"errors"
	"fmt"
	"log"
	"time"
)

type Database interface {
	GetPasswordByLogin(login string) (string, error)
}

type realAuthService struct{}

var (
	ErrInvalidLoginData = errors.New("invalid login data")
	ErrUserNotFound     = errors.New("user not found")
)

func (auth *realAuthService) GetPasswordByLogin(login string) (string, error) {
	// make real call to database
	time.Sleep(5 * time.Second)

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

func NewApplication(auth Database) *Application {
	return &Application{
		Database: auth,
	}
}

func main() {
	app := NewApplication(&realAuthService{})

	if err := app.Login("admin", "admin"); err != nil {
		log.Fatalln(err)
	}

	fmt.Println("ok")
}
