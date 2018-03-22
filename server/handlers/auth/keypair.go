package auth

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"

	"github.com/sean-/seed"
	"golang.org/x/crypto/ssh"
)

func init() {
	seed.MustInit()
}

type KeyPair struct {
	PublicKey      ssh.PublicKey
	PrivateKey     *rsa.PrivateKey
	FingerprintMD5 string

	publicKeyBase64 string
	privateKeyPEM   string
}

func NewKeyPair(bits int) (*KeyPair, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, err
	}

	sshPublicKey, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, err
	}
	fingerprintMD5 := ssh.FingerprintLegacyMD5(sshPublicKey)

	return &KeyPair{
		PrivateKey:     privateKey,
		PublicKey:      sshPublicKey,
		FingerprintMD5: fingerprintMD5,
	}, nil
}

func DecodeKeyPair(material string) (*KeyPair, error) {
	privateKey, err := x509.ParsePKCS1PrivateKey([]byte(material))
	if err != nil {
		return nil, err
	}

	sshPublicKey, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, err
	}
	fingerprintMD5 := ssh.FingerprintLegacyMD5(sshPublicKey)

	return &KeyPair{
		PrivateKey:     privateKey,
		PublicKey:      sshPublicKey,
		FingerprintMD5: fingerprintMD5,
	}, nil
}

func (kp *KeyPair) genPrivateKeyPEM() error {
	privateKeyBuff := &bytes.Buffer{}
	privatePEM := bufio.NewWriter(privateKeyBuff)

	privatePEMBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(kp.PrivateKey),
	}
	if err := pem.Encode(privatePEM, privatePEMBlock); err != nil {
		return err
	}
	privatePEM.Flush()

	kp.privateKeyPEM = fmt.Sprintf("%s", privateKeyBuff)

	return nil
}

func (kp *KeyPair) PrivateKeyPEM() string {
	if kp.privateKeyPEM == "" {
		err := kp.genPrivateKeyPEM()
		if err != nil {
			return ""
		}
	}

	return kp.privateKeyPEM
}

func (kp *KeyPair) genPublicKeyBase64() {
	publicKeyBase := base64.StdEncoding.EncodeToString(kp.PublicKey.Marshal())
	kp.publicKeyBase64 = "ssh-rsa " + string(publicKeyBase)
}

func (kp *KeyPair) PublicKeyBase64() string {
	if kp.publicKeyBase64 == "" {
		kp.genPublicKeyBase64()
	}

	return kp.publicKeyBase64
}
