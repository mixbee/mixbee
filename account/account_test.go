package account

import ("testing"
	"github.com/stretchr/testify/assert"
	"log"
)

func TestNewAccount(t  *testing.T)  {
	account := NewAccount()

	assert.NotNil(t, account)
	log.Println("account", account)
	log.Println("PrivateKey", account.PrivateKey)
	log.Println("PublicKey", account.PublicKey)

}

