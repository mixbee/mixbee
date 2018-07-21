package account

import (
	"github.com/mixbee/mixbee-crypto/keypair"
	"testing"
	"encoding/hex"
	"os"
	"github.com/stretchr/testify/assert"
	"sort"
	"github.com/mixbee/mixbee/core/types"

)

func genAccountData() (*AccountData, *keypair.ProtectedKey)  {
	var acc = new(AccountData)
	prvkey, pubkey, _ := keypair.GenerateKeyPair(keypair.PK_ECDSA, keypair.P256)
	ta := types.AddressFromPubKey(pubkey)
	address := ta.ToBase58()
	password := []byte("123456")
	prvSectet, _ := keypair.EncryptPrivateKey(prvkey, address, password)
	acc.SetKeyPair(prvSectet)
	acc.SigSch = "SHA256withECDSA"
	acc.PubKey = hex.EncodeToString(keypair.SerializePublicKey(pubkey))
	return acc, prvSectet
}

func TestAccountData(t *testing.T) {
	account,proKey := genAccountData()
	t.Log("result account:", account)
	t.Log("proKey account", proKey)
}

func TestWalletSave(t *testing.T)  {
	walletFile := "w.data"
	defer func() {
		os.Remove(walletFile)
		os.RemoveAll("Log/")
	}()

	wallet := NewWalletData()
	size := 10
	for i := 0; i < size; i++ {
		acc, _ := genAccountData()
		wallet.AddAccount(acc)
		// valletdate store to file
		err := wallet.Save(walletFile)
		if err != nil {
			t.Errorf("Save error:%s", err)
			return
		}
	}

	// load valletdate from file
	wallet2 := NewWalletData()
	err := wallet2.Load(walletFile)
	if err != nil {
		t.Errorf("Load error:%s", err)
		return
	}
	// t.Log("wallet2", wallet2)
	assert.Equal(t, len(wallet2.Accounts), len(wallet.Accounts))

}
func TestWalletDel(t *testing.T)  {
	walletDate := NewWalletData()
	size := 10
	accList := make([]string, 0, size)
	for i:=0; i<size; i++ {
		acc,_ := genAccountData()
		walletDate.AddAccount(acc)
		accList = append(accList, acc.Address)
	}
	sort.Strings(accList)
	for _,address := range accList {
		walletDate.DelAccount(address)
		_, index := walletDate.GetAccountByAddress(address)
		if !assert.Equal(t, -1, index) {
			return
		}
	}
}

