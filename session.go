package main

import (
	"bufio"
	"github.com/mpetavy/common"
)

type Session struct {
	con common.EndpointConnection
}

func NewSession(con common.EndpointConnection) (*Session, error) {
	return &Session{con: con}, nil
}

func (session *Session) Run() {
	for {
		scanner := bufio.NewScanner(session.con)
		for scanner.Scan() {
			cmd := scanner.Text()

			common.Debug("common received: %s", cmd)
		}
	}
}
