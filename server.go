package main

import (
	"crypto/tls"
	"github.com/mpetavy/common"
)

type Server struct {
	tlsConfig *tls.Config
	endpoint  common.Endpoint
	connector common.EndpointConnector
}

func NewServer(address string, tls bool) (*Server, error) {
	var err error

	server := &Server{}

	server.tlsConfig, err = common.NewTlsConfigFromFlags()

	server.endpoint, server.connector, err = common.NewEndpoint(address, false, server.tlsConfig)
	if common.Error(err) {
		return nil, err
	}

	return server, nil
}

func (server *Server) Run() error {
	err := server.endpoint.Start()
	if common.Error(err) {
		return err
	}

	defer func() {
		common.Error(server.endpoint.Stop())
	}()

	for {
		con, err := server.connector()
		if common.Error(err) {
			break
		}

		session, err := NewSession(con)
		if common.Error(err) {
			break
		}

		go session.Run()
	}

	return nil
}
