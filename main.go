package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/mpetavy/common"
	"os"
)

func init() {
	common.Init("fileshare", "", "", "", "2018", "test", "mpetavy", fmt.Sprintf("https://github.com/mpetavy/%s", common.Title()), common.APACHE, nil, nil, nil, run, 0)
}

type ServerConfig struct {
	common.Configuration
	Address string
	UseTls  bool
	RootDir string
	Users   []*User
}

var (
	serverAddress = flag.String("s", "", "Server address")
	clientAddress = flag.String("c", "", "Client address")
	useTls        = flag.Bool("tls", true, "TLS usage")
	rootDir       = flag.String("p", ".", "Root directory")

	server       *Server
	serverConfig *ServerConfig
)

func startServer() error {
	if common.FileExists(*serverAddress) {
		ba, err := os.ReadFile(*serverAddress)
		if common.Error(err) {
			return err
		}

		serverConfig = &ServerConfig{}

		err = json.Unmarshal(ba, serverConfig)
		if common.Error(err) {
			return err
		}
	} else {
		serverConfig = &ServerConfig{
			Address: *serverAddress,
			UseTls:  *useTls,
			RootDir: *rootDir,
		}
	}

	if len(serverConfig.Users) == 0 {
		serverConfig.Users = []*User{Anonymous}
	}

	var err error

	server, err = NewServer(serverConfig)
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
