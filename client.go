package main

import (
	"crypto/tls"
	"github.com/chzyer/readline"
	"github.com/mpetavy/common"
)

type Client struct {
	TlsConfig *tls.Config
	Endpoint  common.Endpoint
	Connector common.EndpointConnector
	Con       common.EndpointConnection
}

func NewClient(address string, useTls bool) (*Client, error) {
	var err error

	client := &Client{}

	if useTls {
		client.TlsConfig, err = common.NewTlsConfigFromFlags()
		if common.Error(err) {
			return nil, err
		}
	}

	client.Endpoint, client.Connector, err = common.NewEndpoint(address, true, client.TlsConfig)
	if common.Error(err) {
		return nil, err
	}

	return client, nil
}

func (client *Client) Run() error {
	err := client.Endpoint.Start()
	if common.Error(err) {
		return err
	}

	client.Con, err = client.Connector()
	if common.Error(err) {
		return err
	}

	defer func() {
		common.WarnError(client.Endpoint.Stop())
	}()

	rl, err := readline.New("> ")
	if err != nil {
		panic(err)
	}
	defer func() {
		common.WarnError(rl.Close())
	}()

loop:
	for {
		cmd, err := rl.Readline()
		if common.Error(err) {
			break
		}

		_, err = client.Con.Write([]byte(cmd + "\n"))
		if common.Error(err) {
			break
		}

		args := common.Split(cmd, " ")

		switch args[0] {
		case Quit:
			break loop
		}
	}

	return nil
}
