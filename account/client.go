package account

import (
	"github.com/mixbee/mixbee-crypto/keypair"
	s "github.com/mixbee/mixbee-crypto/signature"
)


// about wallet client
type Client interface {

	//NewAccount create a new account.
	NewAccount(label string, typeCode keypair.KeyType, curveCode byte, sigScheme s.SignatureScheme, passwd []byte) (*Account, error)
	//ImportAccount import a already exist account to wallet
	ImportAccount(accMeta *AccountMetadata) error
	//GetAccountByAddress return account object by address
	GetAccountByAddress(address string, passwd []byte) (*Account, error)
	//GetAccountByLabel return account object by label
	GetAccountByLabel(label string, passwd []byte) (*Account, error)
	//GetAccountByIndex return account object by index. Index start from 1
	GetAccountByIndex(index int, passwd []byte) (*Account, error)
	//GetDefaultAccount return default account
	GetDefaultAccount(passwd []byte) (*Account, error)
	//GetAccountMetadataByIndex return account Metadata info by address
	GetAccountMetadataByAddress(address string) *AccountMetadata
	//GetAccountMetadataByLabel return account Metadata info by label
	GetAccountMetadataByLabel(label string) *AccountMetadata
	//GetAccountMetadataByIndex return account Metadata info by index. Index start from 1
	GetAccountMetadataByIndex(index int) *AccountMetadata
	//GetDefaultAccountMetadata return default account Metadata info
	GetDefaultAccountMetadata() *AccountMetadata
	//GetAccountNum return total account number
	GetAccountNum() int
	//DeleteAccount delete account
	DeleteAccount(address string, passwd []byte) (*Account, error)
	//UnLockAccount can get account without password in expire time
	UnLockAccount(address string, expiredAt int, passwd []byte) error
	//LockAccount lock unlock account
	LockAccount(address string)
	//GetUnlockAccount return account which was unlock and in expired time
	GetUnlockAccount(address string) *Account
	//Set a new account to default account
	SetDefaultAccount(address string, passwd []byte) error
	//Set a new label to accont
	SetLabel(address, label string, passwd []byte) error
	//Change pasword to account
	ChangePassword(address string, oldPasswd, newPasswd []byte) error
	//Change sig scheme to account
	ChangeSigScheme(address string, sigScheme s.SignatureScheme, passwd []byte) error
	//Get the underlying wallet data
	GetWalletData() *WalletData
}
