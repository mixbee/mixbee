package account

import (
	"sync"
	"time"
	"github.com/mixbee/mixbee-crypto/keypair"
	s "github.com/mixbee/mixbee-crypto/signature"

)


type unlockAccountInfo struct {
	acc        *Account
	unlockTime time.Time
	expiredAt  int //s
}

type ClientImpl struct {
	path       string
	accAddrs   map[string]*AccountData //Map Address(base58) => Account
	accLabels  map[string]*AccountData //Map Label => Account
	defaultAcc *AccountData
	walletData *WalletData
	unlockAccs map[string]*unlockAccountInfo //Map Address(base58) => unlockAccountInfo
	lock       sync.RWMutex
}



//NewAccount create a new account.
func (this *ClientImpl) NewAccount(label string, typeCode keypair.KeyType, curveCode byte, sigScheme s.SignatureScheme, passwd []byte) (*Account, error){

}

//ImportAccount import a already exist account to wallet
func (this *ClientImpl) ImportAccount(accMeta *AccountMetadata) error {

}

//GetAccountByAddress return account object by address
func (this *ClientImpl) GetAccountByAddress(address string, passwd []byte) (*Account, error) {

}

//GetAccountByLabel return account object by label
func (this *ClientImpl) GetAccountByLabel(label string, passwd []byte) (*Account, error){

}

//GetAccountByIndex return account object by index. Index start from 1
func (this *ClientImpl) GetAccountByIndex(index int, passwd []byte) (*Account, error){

}


//GetDefaultAccount return default account
func (this *ClientImpl) GetDefaultAccount(passwd []byte) (*Account, error){

}

//GetAccountMetadataByIndex return account Metadata info by address
func (this *ClientImpl) GetAccountMetadataByAddress(address string) *AccountMetadata {

}

//GetAccountMetadataByLabel return account Metadata info by label
func (this *ClientImpl) GetAccountMetadataByLabel(label string) *AccountMetadata{

}


//GetAccountMetadataByIndex return account Metadata info by index. Index start from 1
func (this *ClientImpl) GetAccountMetadataByIndex(index int) *AccountMetadata{

}

//GetDefaultAccountMetadata return default account Metadata info
func (this *ClientImpl) GetDefaultAccountMetadata() *AccountMetadata{

}

//GetAccountNum return total account number
func (this *ClientImpl) GetAccountNum() int {

}

//DeleteAccount delete account
func (this *ClientImpl) DeleteAccount(address string, passwd []byte) (*Account, error) {

}


//UnLockAccount can get account without password in expire time
func (this *ClientImpl) UnLockAccount(address string, expiredAt int, passwd []byte) error {

}

//LockAccount lock unlock account
func (this *ClientImpl) LockAccount(address string) {

}

//GetUnlockAccount return account which was unlock and in expired time
func (this *ClientImpl) GetUnlockAccount(address string) *Account {

}


//Set a new account to default account
func (this *ClientImpl) SetDefaultAccount(address string, passwd []byte) error {

}


//Set a new label to accont
func (this *ClientImpl) SetLabel(address, label string, passwd []byte) error {

}


//Change pasword to account
func (this *ClientImpl) ChangePassword(address string, oldPasswd, newPasswd []byte) error {

}


//Change sig scheme to account
func (this *ClientImpl) ChangeSigScheme(address string, sigScheme s.SignatureScheme, passwd []byte) error {

}


//Get the underlying wallet data
func (this *ClientImpl) GetWalletData() *WalletData {

}



