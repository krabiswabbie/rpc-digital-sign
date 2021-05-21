package main

import (
	"crypto/ecdsa"
	"errors"
	"net"
	"os"

	eth "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rpc"
)

var server *rpc.Server

var privKey *ecdsa.PrivateKey

type TestService struct{}

type EchoMsg struct {
	Data      []byte
	Signature []byte
	PubKey    []byte
}

func (s *TestService) Echo(msg *EchoMsg) (*EchoMsg, error) {
	if Decode(msg) {
		return Encode("pong", privKey)
	}

	return nil, errors.New("Failed to verify signed message")
}

func newServer(addr string) error {
	os.Remove(addr) // just for simplicity

	var err error
	privKey, err = eth.GenerateKey()
	if err != nil {
		return err
	}

	// RPC server
	server = rpc.NewServer()
	err = server.RegisterName("service", new(TestService))
	if err != nil {
		return err
	}

	ls, err := net.Listen("unix", addr)
	if err != nil {
		return err
	}
	defer ls.Close()

	return server.ServeListener(ls)
}
