package main

import (
	"encoding/hex"
	"errors"
	"math/big"

	"github.com/enkhalifapro/go-web3/shh"
)

// WhisperHelper contains messaging functions
type WhisperHelper struct {
	shh *shh.SHH
}

// Message - whisper message DTO
type Message struct {
	From    string
	To      string
	Topic   string
	Content string
	TTL     int64
}

// NewWhisperHelper constructs whisperHelper
func NewWhisperHelper(shh *shh.SHH) *WhisperHelper {
	return &WhisperHelper{shh}
}

// SendAsymMsg sends a message with asymmetric encryption
func (w *WhisperHelper) SendAsymMsg(msg *Message) error {
	// validate sender key
	if msg.From == "" {
		return errors.New("sender key is empty")
	}
	_, err := hex.DecodeString(msg.From)
	if err != nil {
		return errors.New("invalid sender key")
	}

	// validate recipient key
	if msg.To == "" {
		return errors.New("recipient key is empty")
	}
	_, err = hex.DecodeString(msg.To)
	if err != nil {
		return errors.New("invalid recipient key")
	}

	_, err = w.shh.AsymPost(msg.From, msg.To, msg.Topic, msg.Content, big.NewInt(msg.TTL))
	return err

}

// SendSymMsg sends a message with symmetric encryption
func (w *WhisperHelper) SendSymMsg(password string, msg *Message) error {
	if password == "" {
		return errors.New("Invalid send key")
	}
	symKey, err := w.shh.GenerateSymKeyFromPassword(password)
	if err != nil {
		return err
	}
	_, err = w.shh.SymPost(symKey, msg.To, msg.Topic, msg.Content, big.NewInt(msg.TTL))
	return err
}
