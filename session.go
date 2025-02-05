package main

import (
	"bufio"
	"fmt"
	"github.com/mpetavy/common"
	"slices"
)

type Session struct {
	Con  common.EndpointConnection
	User *User
}

func NewSession(con common.EndpointConnection) (*Session, error) {
	var user *User

	if !Anonymous.Locked {
		var err error
		user, err = common.Clone(Anonymous)
		if common.Error(err) {
			return nil, err
		}
	}

	return &Session{
		Con:  con,
		User: user,
	}, nil
}

func (session *Session) Send(line string) error {
	_, err := session.Con.Write([]byte(line + "\n"))
	if common.Error(err) {
		return err
	}

	return nil
}

func (session *Session) UserCmd(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("no username provided")
	}

	userName := args[1]

	found := slices.IndexFunc(server.Config.Users, func(u *User) bool {
		return userName == u.Login
	})

	if found == -1 {
		return fmt.Errorf("unknown user %s", userName)
	}

	var err error

	session.User, err = common.Clone(serverConfig.Users[found])
	if common.Error(err) {
		return err
	}

	session.User.IsAuthenticated = session.User.Password == ""

	if session.User.IsAuthenticated {
		session.Send("successfull login")
	} else {
		session.Send("please provide password")
	}

	return nil
}

func (session *Session) PassCmd(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("no password provided")
	}

	password := args[1]

	if session.User.Password != password {
		return fmt.Errorf("invalid password")
	} else {
		session.Send("successfull login")
	}

	session.User.IsAuthenticated = true

	return nil
}

func (session *Session) Run() {
	defer func() {
		common.WarnError(session.Con.Close())
	}()

	scanner := bufio.NewScanner(session.Con)
loop:
	for scanner.Scan() {
		common.Debug("common received: %s", scanner.Text())

		args := common.Split(scanner.Text(), " ")

		cmd := args[0]

		var err error

		switch cmd {
		case QuitCmd:
			break loop
		case UserCmd:
			err = session.UserCmd(args)
		case PassCmd:
			err = session.PassCmd(args)
		}

		if common.Error(err) {
			session.Send(err.Error())
		}
	}
}
