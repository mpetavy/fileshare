package main

import (
	"flag"
	"fmt"
	"github.com/mpetavy/common"
)

func init() {
	common.Init("test", "0.0.0", "", "", "2018", "test", "mpetavy", fmt.Sprintf("https://github.com/mpetavy/%s", common.Title()), common.APACHE, nil, nil, nil, run, 0)
}

var (
	serverAddress = flag.String("s", "", "Server address")
	clientAddress = flag.String("c", "", "Client address")
	useTls        = flag.Bool("tls", true, "TLS usage")
	rootDir       = flag.String("p", ".", "Root directory")
)

func startServer() error {
	server, err := NewServer(*serverAddress, *useTls)
	if common.Error(err) {
		return err
	}

	err = server.Run()
	if common.Error(err) {
		return err
	}

	return nil
}

func startClient() error {
	client, err := NewClient(*clientAddress, *useTls)
	if common.Error(err) {
		return err
	}

	err = client.Run()
	if common.Error(err) {
		return err
	}

	return nil
}

func run() error {
	if *serverAddress != "" {
		err := startServer()
		if common.Error(err) {
			return err
		}
	} else {
		err := startClient()
		if common.Error(err) {
			return err
		}
	}

	return nil
}

func main() {
	common.Run([]string{"c|s"})
}
