package decoder

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// https://pilabor.com/blog/2021/05/js-gcm-encrypt-dotnet-decrypt/
func Test_KeyGeneration(t *testing.T) {
	t.Run("generate key", func(t *testing.T) {
		key := GetPrivateKey()
		assert.NotEmpty(t, key)
	})
}

func Test_Password(t *testing.T) {
	for name, test := range map[string]struct {
		ciphertext string
	}{
		"~NB8CcOL#J!H?|Yr": {
			ciphertext: "~NB8CcOL#J!H?|Yr",
		},
		"mnenomic": {
			ciphertext: "couple robot escape silent main once smoke check good basket mimic similar",
		},
		"password": {
			ciphertext: "123456",
		},
	} {
		t.Run(name, func(t *testing.T) {
			passwordKey := GetPrivateKey()
			ciphertext, err := Encrypt([]byte(test.ciphertext), passwordKey)
			require.NoError(t, err)
			password, err := Decrypt(ciphertext, passwordKey)
			require.NoError(t, err)
			assert.Equal(t, test.ciphertext, string(password))
		})
	}

	t.Run("with js data", func(t *testing.T) {
		passwordKey := "fJFRbnYboSfxCZLwYAgOVg=="
		cipher := "mKMSqbZl//nM2UaagSspP6rxMq8FLmKOkE5DKw/PORTBFylNXT6JflcewtQ0xFMl88H8FyaKrQeUxdWIM7nkbp+uHCchcmJcUbVK1o8EKERQMYuF9RGTz4qxugDmvHAzkzEnlLfc"
		password, err := Decrypt(cipher, passwordKey)
		require.NoError(t, err)
		assert.Equal(t, "couple robot escape silent main once smoke check good basket mimic similar", string(password))
	})

	t.Run("with js data 2", func(t *testing.T) {
		passwordKey := "h/zH7+A/F1rwCQfOk5DldXk6zRfdLO8yvlZN2s36HXc="
		cipher := "kqI2NuMryyound0lH8mV7tc7V7sG+M4/QNpkcAq6MWX7jNWkYrBHqw5WIYL/1GC6InizAGkfqRJlyFDnNkjKtMlJMqiQ0nLx0EQtMsbGvKQVDwDO4y8spMkMGXUx09DqWM0MKWkG"
		password, err := Decrypt(cipher, passwordKey)
		require.NoError(t, err)
		assert.Equal(t, "couple robot escape silent main once smoke check good basket mimic similar", string(password))
	})

	t.Run("with js data 3", func(t *testing.T) {
		passwordKey := "FjTOwFMyBNFpHS4YGXfitm7+V31BrYVEUNQXuOnlcpI="
		cipher := "76SJNWbT3Nkbqc2ZrBEtCz/3C3uZlGmoyHjP9wRmzn4FgwlXmGQ6Oe7M/7gffkMS/8p/zSi+Vn28MO0of9jnpwEEx7xd6L8AkQ71awtQad5evI5B90aMezoLp2P/IzhE9IGhHTsC"
		password, err := Decrypt(cipher, passwordKey)
		require.NoError(t, err)
		assert.Equal(t, "couple robot escape silent main once smoke check good basket mimic similar", string(password))
	})
}
