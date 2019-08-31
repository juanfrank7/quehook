package handlers

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}

type mockStorage struct {
	putFileErr  error
	getFilesOut map[string]io.Reader
	getFilesErr error
	getPathsOut []string
	getPathsErr error
}

func (m *mockStorage) PutFile(int, int, int, int, io.Reader) error {
	return m.putFileErr
}

func (m *mockStorage) GetPaths() ([]string, error) {
	return m.getPathsOut, m.getPathsErr
}
