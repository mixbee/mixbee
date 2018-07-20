package account

import (
	"github.com/mixbee/mixbee-crypto/keypair"
	s "github.com/mixbee/mixbee-crypto/signature"
	"github.com/mixbee/mixbee/common"
)

type Account struct {
	PrivateKey keypair.PrivateKey
	PublicKey  keypair.PublicKey
	ProgramHash common.Uint160
	SigScheme  s.SignatureScheme
}

//AccountMetadata all account info without private key
type AccountMetadata struct {
	IsDefault bool   //Is default account
	Label     string //Lable of account
	KeyType   string //KeyType ECDSA,SM2 or EDDSA
	Curve     string //Curve of key type
	Address   string //Address(base58) of account
	PubKey    string //Public  key
	SigSch    string //Signature scheme
	Salt      []byte //Salt
	Key       []byte //PrivateKey in encrypted
	EncAlg    string //Encrypt alg of private key
	Hash      string //Hash alg
}

func NewAccount() *Account  {
	var pkAlgorithm = keypair.PK_ECDSA
	var params = keypair.P256
	var scheme s.SignatureScheme
	scheme = s.SHA256withECDSA

	prk, pub, _ := keypair.GenerateKeyPair(pkAlgorithm, params)
	buf :=  keypair.SerializePublicKey(pub)
	codeHash,_ := common.ToCodeHash(buf)

	return &Account{
		PrivateKey: prk,
		PublicKey:  pub,
		ProgramHash: codeHash,
		SigScheme: scheme,
	}
}



func (this *Account) PrivKey() keypair.PrivateKey {
	return this.PrivateKey
}

func (this *Account) PubKey() keypair.PublicKey {
	return this.PublicKey
}

func (this *Account) getSigScheme() s.SignatureScheme  {
	return this.SigScheme
}









