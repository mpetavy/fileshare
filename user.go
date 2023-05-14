package main

import "github.com/mpetavy/common"

type User struct {
	Login    string
	Password string
	RootDir  string
	Locked   bool
	ReadOnly bool
}

var (
	Anonymous *User
)

func init() {
	common.Events.AddListener(common.EventFlagsParsed{}, func(ev common.Event) {
		var err error

		Anonymous, err = NewUser("anonymous", "", common.CleanPath(*rootDir), false)
		common.Panic(err)
	})
}

func NewUser(login string, password string, rootDir string, locked bool) (*User, error) {
	rootDir = common.CleanPath(rootDir)

	return &User{
		Login:    login,
		Password: "",
		RootDir:  rootDir,
		Locked:   locked,
	}, nil
}
