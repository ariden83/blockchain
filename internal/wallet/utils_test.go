package wallet

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func Test_Password(t *testing.T) {
	var (
		err         error
		passwordBis []byte
		passwordTer string
		passwordOne []byte
	)
	code := "123456"
	password := []byte(code)

	passwordSaved := []byte("$2a$04$7R950HWXZnjFurpuVMB.wOQBcEeCxfDbRmuGPOEp032Awze1VDksC")

	t.Run("test encode password", func(t *testing.T) {
		passwordOne, err = encryptPassword(password)
		assert.NoError(t, err)

		passwordBis, err = encryptPassword(password)
		assert.NoError(t, err)

		p, err := encryptPassword(password)
		assert.NoError(t, err)
		passwordTer = fmt.Sprintf("%s", string(p))

		assert.NotEqual(t, password, passwordBis)
	})

	t.Run("test decode password", func(t *testing.T) {
		err = bcrypt.CompareHashAndPassword(passwordOne, password)
		assert.NoError(t, err)

		err = bcrypt.CompareHashAndPassword(passwordBis, password)
		assert.NoError(t, err)

		err = bcrypt.CompareHashAndPassword(passwordSaved, password)
		assert.NoError(t, err)

		err = bcrypt.CompareHashAndPassword([]byte(passwordTer), password)
		assert.NoError(t, err)
	})
}
