package models

import (
	"crypto/rand"
	"encoding/binary"
	"gindriver/utils"
	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"
)

type User struct {
	Id          uint64                `gorm:"type:bigint(64) unsigned;primary_key"`
	Name        string                `gorm:"type:varchar(100)"`
	DisplayName string                `gorm:"type:varchar(100)"`
	PublicKey   string                `gorm:"type:varchar(100)"`
	credentials []webauthn.Credential `gorm:"-"`
}

func NewUser(name string, displayName string) *User {
	user := &User{}
	user.Id = randomUint64()
	user.Name = name
	user.DisplayName = displayName
	// user.credentials = []webauthn.Credential{}

	return user
}

func randomUint64() uint64 {
	buf := make([]byte, 8)
	rand.Read(buf)
	return binary.LittleEndian.Uint64(buf)
}

// WebAuthnID returns the user's ID
func (u User) WebAuthnID() []byte {
	buf := make([]byte, binary.MaxVarintLen64)
	binary.PutUvarint(buf, uint64(u.Id))
	return buf
}

// WebAuthnName returns the user's username
func (u User) WebAuthnName() string {
	return u.Name
}

// WebAuthnDisplayName returns the user's display name
func (u User) WebAuthnDisplayName() string {
	return u.DisplayName
}

// WebAuthnIcon is not (yet) implemented
func (u User) WebAuthnIcon() string {
	return ""
}

// AddCredential associates the credential to the user
func (u *User) AddCredential(cred webauthn.Credential) {
	u.credentials = append(u.credentials, cred)
}

// WebAuthnCredentials returns credentials owned by the user
func (u User) WebAuthnCredentials() []webauthn.Credential {
	return u.credentials
}

// CredentialExcludeList returns a CredentialDescriptor array filled
// with all the user's credentials
func (u User) CredentialExcludeList() []protocol.CredentialDescriptor {

	credentialExcludeList := []protocol.CredentialDescriptor{}
	for _, cred := range u.credentials {
		descriptor := protocol.CredentialDescriptor{
			Type:         protocol.PublicKeyCredentialType,
			CredentialID: cred.ID,
		}
		credentialExcludeList = append(credentialExcludeList, descriptor)
	}

	return credentialExcludeList
}

func (u User) Insert() (user *User, err error) {
	result := utils.Database.Create(&u)

	user = &u

	if result.Error != nil {
		err = result.Error
		return
	}

	return
}

func GetUserByName(username string) (u User, err error) {
	if err = utils.Database.
		Where("name = ?", username).
		First(&u).Error; err != nil {
		return
	}
	return
}
