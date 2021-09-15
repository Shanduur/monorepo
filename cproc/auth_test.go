package cproc_test

import (
	"os"
	"strings"
	"testing"

	"github.com/shanduur/cproc"
)

var (
	user     = "user"
	password = "passw0rd"
	yes      = "y"
)

func TestNewSafeProcessor(t *testing.T) {
	defer os.Remove("pwd.test")
	sp, err := cproc.NewSafeProcessor(cproc.NewProcessor(
		os.Stdout,
		os.Stderr,
		os.Stdin,
		"> ",
	), "pwd.test")
	if err != nil {
		t.Errorf("failed: %s", err.Error())
	}
	defer sp.Close()
}

func TestLogin(t *testing.T) {
	defer os.Remove("pwd.test")

	sp, err := cproc.NewSafeProcessor(cproc.NewProcessor(os.Stdout,
		os.Stderr,
		strings.NewReader(user+"\n"+password+"\n"+password+"\n"+yes+"\n"),
		"> ",
	), "pwd.test")
	if err != nil {
		t.Errorf("failed: %s", err.Error())
	}
	defer sp.Close()

	if err := sp.Login(user, []byte(password)); err != nil {
		t.Errorf("login failed: %s", err.Error())
	}
}
