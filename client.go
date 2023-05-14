package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"github.com/mpetavy/common"
	"os"
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

func (client *Client) Send(line string) error {
	_, err := client.Con.Write([]byte(line + "\n"))
	if common.Error(err) {
		return err
	}

	return nil
}

func (client *Client) QuitCmd(args []string) error {
	err := client.Send(QuitCmd)
	if common.Error(err) {
		return err
	}

	return nil
}

func (client *Client) UserCmd() error {
	err := client.Send(UserCmd)
	if common.Error(err) {
		return err
	}

	return nil
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

	conScanner := bufio.NewScanner(client.Con)
	inputScanner := bufio.NewScanner(os.Stdin)

loop:
	for {
		if !inputScanner.Scan() {
			break loop
		}

		line := inputScanner.Text()

		err = client.Send(line)
		if common.Error(err) {
			return err
		}

		args := common.Split(line, " ")

		cmd := args[0]

		if cmd == QuitCmd {
			break loop
		}

		if !conScanner.Scan() {
			break
		}

		fmt.Printf("%s\n", conScanner.Text())
	}

	return nil
}
