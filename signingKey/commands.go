package signingKey

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	pui "github.com/manifoldco/promptui"
	cli "github.com/urfave/cli"
	ed25519 "golang.org/x/crypto/ed25519"
)

var KeyNew = cli.Command{
	Name: "sk:new",
	Action: func(c *cli.Context) error {

		// ask for name of signing key
		p := pui.Prompt{
			Label: "Enter name of new signing key",
		}
		name, err := p.Run()
		if err != nil {
			panic(err)
		}
		if name == "" {
			return errors.New("please enter a name for this signing key")
		}

		// ask for password
		p = pui.Prompt{
			Label: "Enter password for signing key encryption",
			Mask:  '*',
		}
		pw, err := p.Run()
		if err != nil {
			panic(err)
		}

		// ask for password confirmation
		p = pui.Prompt{
			Label: "Confirm password",
			Mask:  '*',
		}
		pwConfirm, err := p.Run()
		if err != nil {
			panic(err)
		}
		if pw != pwConfirm {
			return errors.New("failed to confirm password")
		}

		// create key store
		pub, priv, err := ed25519.GenerateKey(rand.Reader)
		sk := SingingKey{
			Name:       name,
			PublicKey:  pub[:],
			privateKey: priv[:],
			CreateAt:   time.Now(),
			Version:    uint8(1),
		}

		// export signing key
		data, err := sk.Export([]byte(pw))
		if err != nil {
			return err
		}

		fileName := fmt.Sprintf("%s_%d.signing_key.json", name, time.Now().Unix())
		if err := ioutil.WriteFile(fileName, data, 0775); err != nil {
			return err
		}

		fmt.Println(fmt.Sprintf("wrote signing key to: %s", fileName))

		return nil

	},
}
