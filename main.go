package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"

	"github.com/enkhalifapro/go-web3/utils"
	"github.com/ethereum/go-ethereum/common/hexutil"
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
	Topic      string
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

	topicPrompt := promptui.Prompt{
		Label: "Topic: ",
	}
	topic, err := topicPrompt.Run()
	if err != nil {
		return nil, err
	}
	config.Topic = hexutil.Encode([]byte(topic))
	return config, nil
}

// excute command
func run(command string, arguments ...string) error {
	cmd := exec.Command(command, arguments...)

	// Connect pipe to read Stderr
	stderr, err := cmd.StderrPipe()

	if err != nil {
		// Failed to connect pipe
		return fmt.Errorf("%q failed to connect stderr pipe: %v", command, err)
	}

	// Do not use cmd.Run()
	if err := cmd.Start(); err != nil {
		// Problem while copying stdin, stdout, or stderr
		return fmt.Errorf("%q failed: %v", command, err)
	}

	// Zero exit status
	// Darwin: launchctl can fail with a zero exit status,
	// so check for emtpy stderr
	if command == "launchctl" {
		slurp, _ := ioutil.ReadAll(stderr)
		if len(slurp) > 0 {
			return fmt.Errorf("%q failed with stderr: %s", command, slurp)
		}
	}

	return nil
}

func main() {

	config, err := readConfig()
	if err != nil {
		panic(err)
	}

	// run geth
	err = run("geth", "--testnet", "--light", "--rpc", "--shh", "--rpcport", "8545", "--rpcaddr", "127.0.0.1", "--rpccorsdomain", "*")
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

	// messages listner
	getMsgs := func() {
		msgs := shh.GetFilterMsgs(filterID)
		if len(msgs) > 0 {
			for _, msg := range msgs {
				fmt.Println(utils.DecodeHex(msg.Payload))
			}
		}
	}

	scheduler.Every(1).Seconds().Run(getMsgs)

	// check new messages to send
	for {
		newMsgPrompt := promptui.Prompt{
			Label: "",
		}
		msgContent, _ := newMsgPrompt.Run()
		fakeRecpient := "0x0477e7a5e6215d00df2c19fbfc4241973984e5ab114a10346e894e37699c41186b4ada203b925dd37a3dcb4df609c1d3b8151d38a98a87307624a7108648450008"
		msg := &message{From: config.PrivateKey,
			To:      fakeRecpient,
			Topic:   "0xdeadbeef",
			Content: msgContent,
			TTL:     7,
		}

		if config.PrivateKey != "" {
			// send asym
			err = sendAsymMsg(shh, msg)
			if err != nil {
				panic(err)
			}
		} else {
			// send asym
			err = sendSymMsg(shh, config.Password, msg)
			if err != nil {
				panic(err)
			}
		}
	}
}
