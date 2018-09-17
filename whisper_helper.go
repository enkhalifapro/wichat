package main

import (
	"math/big"

	"github.com/enkhalifapro/go-web3/shh"
)

// IWhisperHelper whisper functions abstraction
type IWhisperHelper interface {
}

type message struct {
	From    string
	To      string
	Topic   string
	Content string
	TTL     int64
}

func sendAsymMsg(shh *shh.SHH, msg *message) error {
	_, err := shh.AsymPost(msg.From, msg.To, msg.Topic, msg.Content, big.NewInt(msg.TTL))
	return err

}

func sendSymMsg(shh *shh.SHH, password string, msg *message) error {
	symKey, err := shh.GenerateSymKeyFromPassword(password)
	if err != nil {
		return err
	}
	_, err = shh.SymPost(symKey, msg.To, msg.Topic, msg.Content, big.NewInt(msg.TTL))
	return err
}
