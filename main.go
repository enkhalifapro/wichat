package main

import (
	"fmt"
	"math/big"

	"github.com/enkhalifapro/go-web3/utils"
	"github.com/manifoldco/promptui"

	"github.com/carlescere/scheduler"

	"github.com/enkhalifapro/go-web3/dto"
	"github.com/enkhalifapro/go-web3/providers"
	"github.com/enkhalifapro/go-web3/shh"
)

type config struct {
	NickName   string
	IsAsym     bool
	PrivateKey string
	Password   string
}

type message struct {
	From    string
	To      string
	Topic   string
	Content string
	TTL     int64
}

func readConfig() (*config, error) {
	config := &config{IsAsym: false}
	nickNameprompt := promptui.Prompt{
		Label: "Enter Nickname: ",
	}
	nickName, err := nickNameprompt.Run()

	if err != nil {
		return nil, err
	}
	config.NickName = nickName

	encTypePrompt := promptui.Select{
		Label: "Encryption type: ",
		Items: []string{"symmetric", "asymmetric"},
	}

	_, encType, err := encTypePrompt.Run()

	if err != nil {
		return nil, err
	}

	if encType == "asymmetric" {
		config.IsAsym = true
		// get privateKey
		privateKeyPrompt := promptui.Prompt{
			Label: "Private key: ",
		}

		privateKey, err := privateKeyPrompt.Run()

		if err != nil {
			return nil, err
		}

		config.PrivateKey = privateKey
	} else {
		// get password
		passPrompt := promptui.Prompt{
			Label: "Password: ",
		}

		pass, err := passPrompt.Run()

		if err != nil {
			return nil, err
		}

		config.Password = pass
	}

	return config, nil
}

func sendMsg(shh *shh.SHH, msg *message) error {
	_, err := shh.AsymPost(msg.From, msg.To, msg.Topic, msg.Content, big.NewInt(msg.TTL))
	return err
}

func main() {

	config, err := readConfig()
	if err != nil {
		panic(err)
	}

	provider := providers.NewHTTPProvider("127.0.0.1:8545", 10, false)
	shh := shh.NewSHH(provider)

	keyID := config.PrivateKey

	// 2- create a message filter
	filterID, err := shh.NewMsgFilter(&dto.SHHSubscribeParam{
		PrivateKeyID: keyID,
		Topics:       []string{"0xdeadbeef"},
	})
	if err != nil {
		panic(err)
	}

	getMsgs := func() {
		msgs := shh.GetFilterMsgs(filterID)
		if len(msgs) > 0 {
			for _, msg := range msgs {
				fmt.Println("in rec")
				fmt.Println(utils.DecodeHex(msg.Payload))
			}
		}
	}

	scheduler.Every(1).Seconds().Run(getMsgs)

	for {
		msgPrompt := promptui.Prompt{
			Label: "",
		}
		msg, _ := msgPrompt.Run()
		fakeRecpient := "0x0477e7a5e6215d00df2c19fbfc4241973984e5ab114a10346e894e37699c41186b4ada203b925dd37a3dcb4df609c1d3b8151d38a98a87307624a7108648450008"
		err = sendMsg(shh, &message{From: config.PrivateKey,
			To:      fakeRecpient,
			Topic:   "0xdeadbeef",
			Content: msg,
			TTL:     7,
		})
		if err != nil {
			panic(err)
		}
	}
}
