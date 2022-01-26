package decoder

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_Password(t *testing.T) {
	passwordKey := "~NB8CcOL#J!H?|Yr"
	for name, test := range map[string]struct {
		ciphertext string
		iv         string
		key        string
		waiting    string
	}{
		/*"mnenomic": {
			ciphertext: "4JcNCZ2I2v7xMb1YmC9VUfaio8zc5RaAmlNYvJQpNVpVWiHbPVGHwCiooG0DzCj5uXHs8CxALK/UFulVCSqal+JBNB9wGtNH86Uv1CToTLgfMn7a+PSbI+s8dUz2gyby9QRusofiFAGmJZmLzHldPA==",
			iv:         "XAoOABKWiwFC5KpMVDX+HQ==",
			key:        passwordKey,
			waiting: "123456",
		},*/
		"password": {
			ciphertext: "cueX1sBCpr2E3pGb7701+g==",
			iv:         "mUwNWiTN9GyASnDUMpRbjA==",
			key:        passwordKey,
			waiting:    "123456",
		},
	} {
		t.Run(name, func(t *testing.T) {
			str, err := Password(test.ciphertext, test.iv, test.key)
			require.NoError(t, err)
			assert.NotEmptyf(t, str, "error message %s", "formatted")
			assert.Equal(t, test.waiting, str)
		})
	}

	t.Run("encrypt", func(t *testing.T) {
		ciphertext, iv, err := Encrypt([]byte("123456"), []byte(passwordKey))
		fmt.Println(fmt.Sprintf("**************** %+v", err))
		require.NoError(t, err)
		assert.NotEmptyf(t, ciphertext, "error message %s", "formatted")
		assert.NotEmptyf(t, iv, "error message %s", "formatted")
		//assert.Equal(t, test.waiting, str)
	})
}
