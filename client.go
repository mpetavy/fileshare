package main

import (
	"crypto/tls"
	"github.com/chzyer/readline"
	"github.com/mpetavy/common"
)

type Client struct {
	tlsConfig *tls.Config
	endpoint  common.Endpoint
	connector common.EndpointConnector
	con       common.EndpointConnection
}

func NewClient(address string, tls bool) (*Client, error) {
	var err error

	client := &Client{}

	client.tlsConfig, err = common.NewTlsConfigFromFlags()

	client.endpoint, client.connector, err = common.NewEndpoint(address, true, client.tlsConfig)
	if common.Error(err) {
		return nil, err
	}

	return client, nil
}

func (client *Client) Run() error {
	err := client.endpoint.Start()
	if common.Error(err) {
		return err
	}

	client.con, err = client.connector()
	if common.Error(err) {
		return err
	}

	defer func() {
		common.Error(client.endpoint.Stop())
	}()

	rl, err := readline.New("> ")
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	for {
		line, err := rl.Readline()
		if err != nil { // io.EOF
			break
		}

		client.con.Write([]byte(line + "\n"))

	}

	return nil
}
