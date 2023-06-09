package secure

import (
	"github.com/fernet/fernet-go"
	"time"
)

type Shield struct {
	EncryptionKey *fernet.Key
}

func NewShield(key string) *Shield {
	k := fernet.MustDecodeKeys(key)
	return &Shield{
		EncryptionKey: k[0],
	}
}

func (s *Shield) DecryptMessage(cipher string) string {
	tok := fernet.VerifyAndDecrypt([]byte(cipher), 0*time.Second, []*fernet.Key{s.EncryptionKey})
	return string(tok)
}

func (s *Shield) DecryptMessageLink(cipher string) *string {
	tok := fernet.VerifyAndDecrypt([]byte(cipher), 0*time.Second, []*fernet.Key{s.EncryptionKey})
	sTok := string(tok)
	return &sTok
}

func (s *Shield) EncryptMessage(message string) (string, error) {
	tok, err := fernet.EncryptAndSign([]byte(message), s.EncryptionKey)
	return string(tok), err
}
