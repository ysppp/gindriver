package utils

import (
	"github.com/duo-labs/webauthn/webauthn"
)

var WebAuthn *webauthn.WebAuthn

func InitWebAuthn() (err error) {
	WebAuthn, err = webauthn.New(&webauthn.Config{
		RPID:          "localhost",
		RPDisplayName: "GinDriver",
	})
	return err
}
