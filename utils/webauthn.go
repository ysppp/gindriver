package utils

import (
	"gindriver/config"
	"gindriver/session"
	"github.com/duo-labs/webauthn/webauthn"
)

var (
	WebAuthn             *webauthn.WebAuthn
	WebAuthnSessionStore *session.Store
)

func InitWebAuthn() (err error) {
	WebAuthn, err = webauthn.New(&webauthn.Config{
		RPID:          config.Config.RPID,
		RPOrigin:      config.Config.RPOrigin,
		RPDisplayName: config.Config.RPDisplayName,
	})
	WebAuthnSessionStore, err = session.NewStore()
	return err
}
