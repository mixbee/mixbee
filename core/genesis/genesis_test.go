
package genesis

import (
	"github.com/mixbee/mixbee-crypto/keypair"
	"github.com/mixbee/mixbee/common/log"
	"github.com/mixbee/mixbee/common"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"encoding/hex"
	"github.com/mixbee/mixbee/common/config"
	"encoding/json"
)

func TestMain(m *testing.M) {
	log.InitLog(0, log.Stdout)
	m.Run()
	os.RemoveAll("./ActorLog")
}

func TestGenesisBlockInit(t *testing.T) {
	_, pub1, _ := keypair.GenerateKeyPair(keypair.PK_ECDSA, keypair.P256)
	_, pub2, _ := keypair.GenerateKeyPair(keypair.PK_ECDSA, keypair.P256)
	_, pub3, _ := keypair.GenerateKeyPair(keypair.PK_ECDSA, keypair.P256)
	_, pub4, _ := keypair.GenerateKeyPair(keypair.PK_ECDSA, keypair.P256)

	t.Log("pub1", 	hex.EncodeToString(keypair.SerializePublicKey(pub1)))
	t.Log("pub2", 	hex.EncodeToString(keypair.SerializePublicKey(pub1)))
	t.Log("pub3", 	hex.EncodeToString(keypair.SerializePublicKey(pub1)))
	t.Log("pub4", 	hex.EncodeToString(keypair.SerializePublicKey(pub1)))

	dbConf := &config.DBFTConfig {
		Bookkeepers: []string{
			hex.EncodeToString(keypair.SerializePublicKey(pub1)),
			hex.EncodeToString(keypair.SerializePublicKey(pub2)),
			hex.EncodeToString(keypair.SerializePublicKey(pub3)),
			hex.EncodeToString(keypair.SerializePublicKey(pub4)),
		},
	}
	conf := &config.GenesisConfig{
		SeedList: []string{
		"node1.example.com:20338",
		"node2.example.com:20338",
		"node3.example.com:20338",
		"node4.example.com:20338",},
		ConsensusType: "dbft",
		DBFT: dbConf,
	}


	block, err := BuildGenesisBlock([]keypair.PublicKey{pub1,pub2,pub3,pub4}, conf)
	b, err := json.Marshal(block)
	t.Log("block", string(b))
	assert.Nil(t, err)
	assert.NotNil(t, block)
	assert.NotEqual(t, block.Header.TransactionsRoot, common.UINT256_EMPTY)
}

func TestNewParamDeployAndInit(t *testing.T) {
	deployTx := newParamContract()
	initTx := newParamInit()
	assert.NotNil(t, deployTx)
	assert.NotNil(t, initTx)
}
