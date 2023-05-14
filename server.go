package main

import (
	"crypto/tls"
	"github.com/mpetavy/common"
)

type Server struct {
	Config    *ServerConfig
	Endpoint  common.Endpoint
	Connector common.EndpointConnector
}

func NewServer(config *ServerConfig) (*Server, error) {
	var err error

	server := &Server{
		Config:    config,
		Endpoint:  nil,
		Connector: nil,
	}

	var tlsConfig *tls.Config

	if config.UseTls {
		tlsConfig, err = common.NewTlsConfigFromFlags()
		if common.Error(err) {
			return nil, err
		}
	}

	server.Endpoint, server.Connector, err = common.NewEndpoint(config.Address, false, tlsConfig)
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
