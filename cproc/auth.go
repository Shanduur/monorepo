package cproc

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/term"
)

// SafeProcessor TODO
type SafeProcessor struct {
	Cp             *CommandProcessor
	PasswdFilePath string
	PasswdFile     *os.File
}

var fd = int(os.Stdin.Fd())

// NewSafeProcessor TODO
func NewSafeProcessor(cp *CommandProcessor, filepath string) (*SafeProcessor, error) {
	f, err := os.OpenFile(filepath, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, fmt.Errorf("unable to open file: %s", err.Error())
	}

	fi, err := f.Stat()
	if fi.Size() < 1 {
		if err := CreateUser(cp, f); err != nil {
			return nil, fmt.Errorf("unable to create user: %s", err.Error())
		}
	}

	sp := &SafeProcessor{
		Cp:             cp,
		PasswdFilePath: filepath,
		PasswdFile:     f,
	}

	cp.AddCommand("new-user",
		func(s ...string) {},
		"built-in",
		"command that allows to add new user",
		`Using this command you can call user creation prompt`).SetArgs(ArgIgnore)

	return sp, nil
}

// Close TODO
func (sp *SafeProcessor) Close() {
	sp.PasswdFile.Close()
}

// CreateUser TODO
func CreateUser(cp *CommandProcessor, f *os.File) error {
	for {
		cp.Printf("New username: ")
		uname, err := cp.Reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("unable to read username: %s", err.Error())
		}
		uname = strings.TrimRight(uname, "\r\n")

	Repeat:
		cp.Printf("New password: ")
		pwd, err := term.ReadPassword(fd)
		cp.Printf("\n")
		if err != nil {
			return fmt.Errorf("unable to read password: %s", err.Error())
		}

		cp.Printf("Repeat password: ")
		pwd2, err := term.ReadPassword(fd)
		cp.Printf("\n")
		if err != nil {
			return fmt.Errorf("unable to read password: %s", err.Error())
		}

		if bytes.Compare(pwd, pwd2) == 0 {
			cp.Printf("Do you want to create user? (Y/n) ")
			r, err := cp.Reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("unable to read input: %s", err.Error())
			}
			r = strings.TrimRight(r, "\n\r")
			if r != "y" && r != "Y" {
				continue
			}
		} else {
			cp.Printf("Pasword does not match, try again.\n")
			goto Repeat
		}

		b, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("unable to generate bcrypt: %s", err.Error())
		}

		if _, err := f.WriteString(uname + "|"); err != nil {
			return fmt.Errorf("unable to save username to file: %s", err.Error())
		}

		if _, err := f.Write(b); err != nil {
			return fmt.Errorf("unable to save bcrypt password to file: %s", err.Error())
		}

		if _, err := f.WriteString("\n"); err != nil {
			return fmt.Errorf("unable to save newline to file: %s", err.Error())
		}

		if err := f.Sync(); err != nil {
			return fmt.Errorf("unable to sync file: %s", err.Error())
		}
		break
	}

	return nil
}

// Login TODO
func (sp *SafeProcessor) Login(uname string, password []byte) error {
	if _, err := sp.PasswdFile.Seek(0, 0); err != nil {
		return fmt.Errorf("unable to seek the begining: %s", err.Error())
	}

	r := bufio.NewReader(sp.PasswdFile)
	var s []string

	for {
		line, _, err := r.ReadLine()
		if errors.Is(err, io.EOF) {
			return fmt.Errorf("user does not exist")
		} else if err != nil {
			return fmt.Errorf("encountered error: %s", err.Error())
		}

		s = strings.SplitN(string(line), "|", 2)
		if len(s) != 2 {
			continue
		} else if s[0] == uname {
			break
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(s[1]), []byte(password)); err != nil {
		return fmt.Errorf("encountered error: %s", err.Error())
	}

	return nil
}

// LoginPrompt TODO
func (sp *SafeProcessor) LoginPrompt() (string, []byte, error) {
	sp.Cp.Printf("Username: ")
	uname, err := sp.Cp.Reader.ReadString('\n')
	if errors.Is(err, io.EOF) {
		return "", nil, io.EOF
	} else if err != nil {
		return "", nil, fmt.Errorf("unable to read: %s", err.Error())
	}
	uname = strings.TrimRight(uname, "\r\n")

	sp.Cp.Printf("Password: ")
	pwd, err := term.ReadPassword(fd)
	sp.Cp.Printf("\n")
	if errors.Is(err, io.EOF) {
		return uname, pwd, io.EOF
	} else if err != nil {
		return "", nil, fmt.Errorf("unable to read: %s", err.Error())
	}
	return uname, pwd, nil
}

// Run TODO
func (sp *SafeProcessor) Run() {
	for {
		user, password, err := sp.LoginPrompt()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			sp.Cp.Log.Fatalf("failed to create login prompt: %s", err.Error())
			continue
		}

		err = sp.Login(user, password)
		if err != nil {
			sp.Cp.Log.Printf("failed to login: %s", err.Error())
			continue
		}
		sp.Cp.Run()
	}
}
