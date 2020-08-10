package utils

import (
	"gindriver/config"
	"github.com/duo-labs/webauthn/webauthn"
)

var WebAuthn *webauthn.WebAuthn

func InitWebAuthn() (err error) {
	WebAuthn, err = webauthn.New(&webauthn.Config{
		RPID:          config.Config.RPID,
		RPDisplayName: config.Config.RPDisplayName,
	})
	return err
}
