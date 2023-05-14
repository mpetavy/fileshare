package main

import (
	"bufio"
	"github.com/mpetavy/common"
)

type Session struct {
	con  common.EndpointConnection
	user *User
}

func NewSession(con common.EndpointConnection) (*Session, error) {
	var user *User

	if !Anonymous.Locked {
		user = Anonymous
	}

	return &Session{
		con:  con,
		user: user,
	}, nil
}

func (session *Session) Run() {
	defer func() {
		common.WarnError(session.con.Close())
	}()

	scanner := bufio.NewScanner(session.con)
loop:
	for scanner.Scan() {
		common.Debug("common received: %s", scanner.Text())

		args := common.Split(scanner.Text(), " ")

		switch args[0] {
		case "quit":
			break loop
		}
	}

	common.Debug("session quit")
}
