package main

import (
	"crypto/ecdsa"
	"errors"
	"time"

	eth "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rpc"
)

func Encode(msg string, priv *ecdsa.PrivateKey) (*EchoMsg, error) {
	hash := eth.Keccak256([]byte(msg))
	sign, err := eth.Sign(hash, priv)
	if err != nil {
		return nil, err
	}

	pub, err := eth.Ecrecover(hash, sign)
	if err != nil {
		return nil, err
	}

	return &EchoMsg{
		Data:      []byte(msg),
		Signature: sign,
		PubKey:    pub,
	}, nil
}

func Decode(msg *EchoMsg) bool {
	sign := msg.Signature[:len(msg.Signature)-1] // remove recovery id
	hash := eth.Keccak256(msg.Data)
	return eth.VerifySignature(msg.PubKey, hash, sign)
}

func callEcho(c *rpc.Client, msg *EchoMsg) (*EchoMsg, error) {
	resp := new(EchoMsg)
	err := c.Call(resp, "service_echo", msg)
	return resp, err
}

func newClient(addr string) error {
	privKey, err := eth.GenerateKey()
	if err != nil {
		return err
	}

	c, err := rpc.Dial(addr)
	if err != nil {
		return err
	}
	defer c.Close()

	for {
		msg, err := Encode("ping", privKey)
		if err != nil {
			return err
		}

		resp, err := callEcho(c, msg)
		if err != nil {
			return err
		}

		verified := Decode(resp)
		if verified != true {
			return errors.New("Failed to verify server message")
		}

		time.Sleep(time.Second)
	}
}
