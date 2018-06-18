package signingKey

import (
	"crypto/rand"
	"encoding/json"
	"time"

	scrypt "github.com/Bit-Nation/panthalassa/crypto/scrypt"
	ed25519 "golang.org/x/crypto/ed25519"
)

type SingingKey struct {
	Name       string            `json:"name"`
	PublicKey  ed25519.PublicKey `json:"public_key"`
	PrivateKey scrypt.CipherText `json:"private_key"`
	privateKey ed25519.PrivateKey
	CreateAt   time.Time `json:"created_at"`
	Version    uint8     `json:"version"`
}

// sign data with private key
func (k *SingingKey) Sign(data []byte) []byte {
	return ed25519.Sign(k.privateKey, data)
}

// create new key
func New(name string) (SingingKey, error) {

	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return SingingKey{}, err
	}

	return SingingKey{
		Name:       name,
		PublicKey:  pub,
		privateKey: priv,
		CreateAt:   time.Now(),
		Version:    1,
	}, nil

}

// export DApp Key
func (k SingingKey) Export(password []byte) ([]byte, error) {

	// encrypt private key
	ct, err := scrypt.NewCipherText(k.privateKey, password)
	if err != nil {
		return nil, err
	}

	k.PrivateKey = ct
	return json.MarshalIndent(k, "", "    ")

}

// decrypt dapp key
func Decrypt(data, password []byte) (SingingKey, error) {

	// unmarshal DApp
	k := SingingKey{}
	if err := json.Unmarshal(data, &k); err != nil {
		return SingingKey{}, err
	}

	// decrypt encrypted private key
	privateKey, err := scrypt.DecryptCipherText(k.PrivateKey, password)
	if err != nil {
		return SingingKey{}, err
	}

	var pk []byte = privateKey
	k.privateKey = pk

	return k, nil

}
