package account

import (
	"sync"
	"time"
	"github.com/mixbee/mixbee-crypto/keypair"
	s "github.com/mixbee/mixbee-crypto/signature"
	"fmt"
	"encoding/hex"
	"strings"
)


type unlockAccountInfo struct {
	acc        *Account
	unlockTime time.Time
	expiredAt  int //s
}


// wallet date
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
	if len(passwd) == 0 {
		return nil, fmt.Errorf("password cannot empty")
	}
	prvkey, pubkey, err := keypair.GenerateKeyPair(typeCode, curveCode)
	if err != nil {
		return nil, fmt.Errorf("generateKeyPair error:%s", err)
	}
	address := AddressFromPubKey(pubkey)
	addressBase58 := address.ToBase58()

	// buf :=  keypair.SerializePublicKey(pubkey)
	// codeHash,_ := common.ToCodeHash(buf)
	// address ,_ := codeHash.ToAddress()

	prvSecret, err := keypair.EncryptPrivateKey(prvkey, addressBase58, passwd)
	if err != nil {
		return nil, fmt.Errorf("encryptPrivateKey error:%s", err)
	}
	accData := &AccountData{}
	accData.Label = label
	accData.SetKeyPair(prvSecret)
	accData.SigSch = sigScheme.Name()
	accData.PubKey = hex.EncodeToString(keypair.SerializePublicKey(pubkey))

	// account date store to file
	err = this.addAccountData(accData)
	if err != nil {
		return nil, err
	}
	return &Account{
		PrivateKey: prvkey,
		PublicKey:  pubkey,
		Address:    address,
		SigScheme:  sigScheme,
	}, nil
}

//ImportAccount import a already exist account to wallet
func (this *ClientImpl) ImportAccount(accMeta *AccountMetadata) error {
	accData := &AccountData{}
	accData.Label = accMeta.Label
	accData.PubKey = accMeta.PubKey
	accData.SigSch = accMeta.SigSch
	accData.Key = accMeta.Key
	accData.Alg = accMeta.KeyType
	accData.Address = accMeta.Address
	accData.EncAlg = accMeta.EncAlg
	accData.Hash = accMeta.Hash
	accData.Salt = accMeta.Salt
	accData.Param = map[string]string{"curve": accMeta.Curve}

	oldAccMeta := this.GetAccountMetadataByLabel(accData.Label)
	if oldAccMeta != nil {
		//rename
		accData.Label = fmt.Sprintf("%s_1", accData.Label)
	}
	return this.addAccountData(accData)
}

//GetAccountByAddress return account object by address
func (this *ClientImpl) GetAccountByAddress(address string, passwd []byte) (*Account, error) {
	this.lock.Lock()
	defer this.lock.Unlock()
	accData, ok := this.accAddrs[address]
	if !ok {
		return nil, nil
	}
	return this.getAccount(accData, passwd)
}

/**
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
*/

//GetAccountMetadataByLabel return account Metadata info by label
func (this *ClientImpl) GetAccountMetadataByLabel(label string) *AccountMetadata{
	if label == "" {
		return nil
	}
	this.lock.RLock()
	defer this.lock.RUnlock()
	accData, ok := this.accLabels[label]
	if !ok {
		return nil
	}
	return this.getAccountMetadata(accData)
}


func (this *ClientImpl) getAccountMetadata(accData *AccountData) *AccountMetadata {
	accMeta := &AccountMetadata{}
	accMeta.Label = accData.Label
	accMeta.KeyType = accData.Alg
	accMeta.SigSch = accData.SigSch
	accMeta.Key = accData.Key
	accMeta.Address = accData.Address
	accMeta.IsDefault = accData.IsDefault
	accMeta.PubKey = accData.PubKey
	accMeta.EncAlg = accData.EncAlg
	accMeta.Hash = accData.Hash
	accMeta.Curve = accData.Param["curve"]
	accMeta.Salt = accData.Salt
	return accMeta
}

func (this *ClientImpl) addAccountData(accData *AccountData) error {
	if !this.checkSigScheme(accData.Alg, accData.SigSch) {
		return fmt.Errorf("sigScheme:%s does not match KeyType:%s", accData.SigSch, accData.Alg)
	}
	this.lock.Lock()
	defer this.lock.Unlock()
	label := accData.Label
	if label != "" {
		_, ok := this.accLabels[label]
		if ok {
			return fmt.Errorf("duplicate label")
		}
	}
	if len(this.walletData.Accounts) == 0 {
		accData.IsDefault = true
	}
	this.walletData.AddAccount(accData)
	err := this.save()
	if err != nil {
		this.walletData.DelAccount(accData.Address)
		return fmt.Errorf("save error:%s", err)
	}
	this.accAddrs[accData.Address] = accData
	if accData.IsDefault {
		this.defaultAcc = accData
	}
	if label != "" {
		this.accLabels[label] = accData
	}
	return nil
}

func (this *ClientImpl) checkSigScheme(keyType, sigScheme string) bool {
	switch strings.ToUpper(keyType) {
	case "ECDSA":
		switch strings.ToUpper(sigScheme) {
		case "SHA224WITHECDSA":
		case "SHA256WITHECDSA":
		case "SHA384WITHECDSA":
		case "SHA512WITHECDSA":
		case "SHA3-224WITHECDSA":
		case "SHA3-256WITHECDSA":
		case "SHA3-384WITHECDSA":
		case "SHA3-512WITHECDSA":
		case "RIPEMD160WITHECDSA":
		default:
			return false
		}
	case "SM2":
		switch strings.ToUpper(sigScheme) {
		case "SM3WITHSM2":
		default:
			return false
		}
	case "ED25519":
		switch strings.ToUpper(sigScheme) {
		case "SHA512WITHEDDSA":
		default:
			return false
		}
	default:
		return false
	}
	return true
}

func (this *ClientImpl) save() error {
	return this.walletData.Save(this.path)
}

func (this *ClientImpl) getAccount(accData *AccountData, passwd []byte) (*Account, error) {
	privateKey, err := keypair.DecryptWithCustomScrypt(&accData.ProtectedKey, passwd, this.walletData.Scrypt)
	if err != nil {
		return nil, fmt.Errorf("decrypt PrivateKey error:%s", err)
	}
	publicKey := privateKey.Public()
	address := AddressFromPubKey(publicKey)
	//addr := types.AddressFromPubKey(publicKey)
	scheme, err := s.GetScheme(accData.SigSch)
	if err != nil {
		return nil, fmt.Errorf("signature scheme error:%s", err)
	}
	return &Account{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		Address:    address,
		SigScheme:  scheme,
	}, nil
}



