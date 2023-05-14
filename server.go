package main

import (
	"crypto/tls"
	"github.com/mpetavy/common"
)

type Server struct {
	TlsConfig *tls.Config
	Endpoint  common.Endpoint
	Connector common.EndpointConnector
}

func NewServer(address string, uswTls bool) (*Server, error) {
	var err error

	server := &Server{
		TlsConfig: nil,
		Endpoint:  nil,
		Connector: nil,
	}

	if uswTls {
		server.TlsConfig, err = common.NewTlsConfigFromFlags()
		if common.Error(err) {
			return nil, err
		}
	}

	server.Endpoint, server.Connector, err = common.NewEndpoint(address, false, server.TlsConfig)
	if common.Error(err) {
		return nil, err
	}

	return server, nil
}

func (server *Server) Run() error {
	err := server.Endpoint.Start()
	if common.Error(err) {
		return err
	}

	defer func() {
		common.WarnError(server.Endpoint.Stop())
	}()

	for {
		con, err := server.Connector()
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
