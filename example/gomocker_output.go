// Code generated by gomocker v1.2.1. DO NOT EDIT.
//
// For more information see github.com/gomocker/gomocker

package main

import (
	io "io"
	rand "math/rand"
)

// Implementations
type (
	ioReaderImpl     struct{ behavior IoReaderBehavior }
	ioWriterImpl     struct{ behavior IoWriterBehavior }
	ioReadWriterImpl struct{ behavior IoReadWriterBehavior }
	randSourceImpl   struct{ behavior RandSourceBehavior }
	databaseImpl     struct{ behavior DatabaseBehavior }
)

// Check
var (
	_ io.Reader     = &ioReaderImpl{}
	_ io.Writer     = &ioWriterImpl{}
	_ io.ReadWriter = &ioReadWriterImpl{}
	_ rand.Source   = &randSourceImpl{}
	_ Database      = &databaseImpl{}
)

// NewIoReader creates mocked implementation of io.Reader interface
func NewIoReaderMock(behavior IoReaderBehavior) io.Reader {
	return &ioReaderImpl{
		behavior: behavior,
	}
}

// NewIoWriter creates mocked implementation of io.Writer interface
func NewIoWriterMock(behavior IoWriterBehavior) io.Writer {
	return &ioWriterImpl{
		behavior: behavior,
	}
}

// NewIoReadWriter creates mocked implementation of io.ReadWriter interface
func NewIoReadWriterMock(behavior IoReadWriterBehavior) io.ReadWriter {
	return &ioReadWriterImpl{
		behavior: behavior,
	}
}

// NewRandSource creates mocked implementation of rand.Source interface
func NewRandSourceMock(behavior RandSourceBehavior) rand.Source {
	return &randSourceImpl{
		behavior: behavior,
	}
}

// NewDatabase creates mocked implementation of Database interface
func NewDatabaseMock(behavior DatabaseBehavior) Database {
	return &databaseImpl{
		behavior: behavior,
	}
}

type IoReaderBehavior struct {
	Read func(p []byte) (n int, err error)
}

type IoWriterBehavior struct {
	Write func(p []byte) (n int, err error)
}

type IoReadWriterBehavior struct {
	Read  func(p []byte) (n int, err error)
	Write func(p []byte) (n int, err error)
}

type RandSourceBehavior struct {
	Int63 func() (out1 int64)
	Seed  func(seed int64)
}

type DatabaseBehavior struct {
	GetPasswordByLogin func(login string) (password string, err error)
}

func (aeiouy *ioReaderImpl) Read(p []byte) (n int, err error) {
	if aeiouy.behavior.Read != nil {
		return aeiouy.behavior.Read(p)
	}
	return
}

func (aeiouy *ioWriterImpl) Write(p []byte) (n int, err error) {
	if aeiouy.behavior.Write != nil {
		return aeiouy.behavior.Write(p)
	}
	return
}

func (aeiouy *ioReadWriterImpl) Read(p []byte) (n int, err error) {
	if aeiouy.behavior.Read != nil {
		return aeiouy.behavior.Read(p)
	}
	return
}

func (aeiouy *ioReadWriterImpl) Write(p []byte) (n int, err error) {
	if aeiouy.behavior.Write != nil {
		return aeiouy.behavior.Write(p)
	}
	return
}

func (aeiouy *randSourceImpl) Int63() (out1 int64) {
	if aeiouy.behavior.Int63 != nil {
		return aeiouy.behavior.Int63()
	}
	return
}

func (aeiouy *randSourceImpl) Seed(seed int64) {
	if aeiouy.behavior.Seed != nil {
		aeiouy.behavior.Seed(seed)
	}
}

func (aeiouy *databaseImpl) GetPasswordByLogin(login string) (password string, err error) {
	if aeiouy.behavior.GetPasswordByLogin != nil {
		return aeiouy.behavior.GetPasswordByLogin(login)
	}
	return
}
