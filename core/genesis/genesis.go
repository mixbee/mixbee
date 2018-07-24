

package genesis

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/mixbee/mixbee-crypto/keypair"
	"github.com/mixbee/mixbee/common"
	"github.com/mixbee/mixbee/common/config"
	"github.com/mixbee/mixbee/common/constants"
	"github.com/mixbee/mixbee/consensus/vbft/config"
	"github.com/mixbee/mixbee/core/types"
	"github.com/mixbee/mixbee/core/utils"
	"github.com/mixbee/mixbee/smartcontract/service/native/global_params"
	"github.com/mixbee/mixbee/smartcontract/service/native/governance"
	"github.com/mixbee/mixbee/smartcontract/service/native/ont"
	nutils "github.com/mixbee/mixbee/smartcontract/service/native/utils"
	"github.com/mixbee/mixbee/smartcontract/service/neovm"
)

const (
	BlockVersion uint32 = 0
	GenesisNonce uint64 = 2083236893
)

var (
	ONTToken   = newGoverningToken()
	ONGToken   = newUtilityToken()
	ONTTokenID = ONTToken.Hash()
	ONGTokenID = ONGToken.Hash()
)

var GenBlockTime = (config.DEFAULT_GEN_BLOCK_TIME * time.Second)

var INIT_PARAM = map[string]string{
	"gasPrice": "0",
}

var GenesisBookkeepers []keypair.PublicKey

// BuildGenesisBlock returns the genesis block with default consensus bookkeeper list
func BuildGenesisBlock(defaultBookkeeper []keypair.PublicKey, genesisConfig *config.GenesisConfig) (*types.Block, error) {
	//getBookkeeper
	GenesisBookkeepers = defaultBookkeeper
	nextBookkeeper, err := types.AddressFromBookkeepers(defaultBookkeeper)
	if err != nil {
		return nil, fmt.Errorf("[Block],BuildGenesisBlock err with GetBookkeeperAddress: %s", err)
	}
	conf := bytes.NewBuffer(nil)
	if genesisConfig.VBFT != nil {
		genesisConfig.VBFT.Serialize(conf)
	}
	govConfig := newGoverConfigInit(conf.Bytes())
	consensusPayload, err := vconfig.GenesisConsensusPayload(govConfig.Hash(), 0)
	if err != nil {
		return nil, fmt.Errorf("consensus genesus init failed: %s", err)
	}
	//blockdata
	genesisHeader := &types.Header{
		Version:          BlockVersion,
		PrevBlockHash:    common.Uint256{},
		TransactionsRoot: common.Uint256{},
		Timestamp:        constants.GENESIS_BLOCK_TIMESTAMP,
		Height:           uint32(0),
		ConsensusData:    GenesisNonce,
		NextBookkeeper:   nextBookkeeper,
		ConsensusPayload: consensusPayload,

		Bookkeepers: nil,
		SigData:     nil,
	}

	//block
	ont := newGoverningToken()
	ong := newUtilityToken()
	param := newParamContract()
	oid := deployOntIDContract()
	auth := deployAuthContract()
	config := newConfig()

	genesisBlock := &types.Block{
		Header: genesisHeader,
		Transactions: []*types.Transaction{
			ont,
			ong,
			param,
			oid,
			auth,
			config,
			newGoverningInit(),
			newUtilityInit(),
			newParamInit(),
			govConfig,
		},
	}
	genesisBlock.RebuildMerkleRoot()
	return genesisBlock, nil
}

func newGoverningToken() *types.Transaction {
	tx := utils.NewDeployTransaction(nutils.OntContractAddress[:], "ONT", "1.0",
		"Mixbee Team", "jianlongcui23@gmail.com", "Mixbee Network ONT Token", true)
	return tx
}

func newUtilityToken() *types.Transaction {
	tx := utils.NewDeployTransaction(nutils.OngContractAddress[:], "ONG", "1.0",
		"Mixbee Team", "jianlongcui23@gmail.com", "Mixbee Network ONG Token", true)
	return tx
}

func newParamContract() *types.Transaction {
	tx := utils.NewDeployTransaction(nutils.ParamContractAddress[:],
		"ParamConfig", "1.0", "Mixbee Team", "jianlongcui23t@gmail.com",
		"Chain Global Environment Variables Manager ", true)
	return tx
}

func newConfig() *types.Transaction {
	tx := utils.NewDeployTransaction(nutils.GovernanceContractAddress[:], "CONFIG", "1.0",
		"Mixbee Team", "jianlongcui23@gmail.com", "Mixbee Network Consensus Config", true)
	return tx
}

func deployAuthContract() *types.Transaction {
	tx := utils.NewDeployTransaction(nutils.AuthContractAddress[:], "AuthContract", "1.0",
		"Mixbee Team", "jianlongcui23@gmail.com", "Mixbee Network Authorization Contract", true)
	return tx
}

func deployOntIDContract() *types.Transaction {
	tx := utils.NewDeployTransaction(nutils.OntIDContractAddress[:], "OID", "1.0",
		"Mixbee Team", "jianlongcui23@gmail.com", "Mixbee Network ONT ID", true)
	return tx
}

func newGoverningInit() *types.Transaction {
	bookkeepers, _ := config.DefConfig.GetBookkeepers()

	var addr common.Address
	if len(bookkeepers) == 1 {
		addr = types.AddressFromPubKey(bookkeepers[0])
	} else {
		m := (5*len(bookkeepers) + 6) / 7
		temp, err := types.AddressFromMultiPubKeys(bookkeepers, m)
		if err != nil {
			panic(fmt.Sprint("wrong bookkeeper config, caused by", err))
		}
		addr = temp
	}

	distribute := []struct {
		addr  common.Address
		value uint64
	}{{addr, constants.ONT_TOTAL_SUPPLY}}

	args := bytes.NewBuffer(nil)
	nutils.WriteVarUint(args, uint64(len(distribute)))
	for _, part := range distribute {
		nutils.WriteAddress(args, part.addr)
		nutils.WriteVarUint(args, part.value)
	}

	return utils.BuildNativeTransaction(nutils.OntContractAddress, ont.INIT_NAME, args.Bytes())
}

func newUtilityInit() *types.Transaction {
	return utils.BuildNativeTransaction(nutils.OngContractAddress, ont.INIT_NAME, []byte{})
}

func newParamInit() *types.Transaction {
	params := new(global_params.Params)
	var s []string
	for k, _ := range INIT_PARAM {
		s = append(s, k)
	}

	neovm.GAS_TABLE.Range(func(key, value interface{}) bool {
		INIT_PARAM[key.(string)] = strconv.FormatUint(value.(uint64), 10)
		s = append(s, key.(string))
		return true
	})

	sort.Strings(s)
	for _, v := range s {
		params.SetParam(global_params.Param{Key: v, Value: INIT_PARAM[v]})
	}
	bf := new(bytes.Buffer)
	params.Serialize(bf)

	bookkeepers, _ := config.DefConfig.GetBookkeepers()
	var addr common.Address
	if len(bookkeepers) == 1 {
		addr = types.AddressFromPubKey(bookkeepers[0])
	} else {
		m := (5*len(bookkeepers) + 6) / 7
		temp, err := types.AddressFromMultiPubKeys(bookkeepers, m)
		if err != nil {
			panic(fmt.Sprint("wrong bookkeeper config, caused by", err))
		}
		addr = temp
	}
	nutils.WriteAddress(bf, addr)

	return utils.BuildNativeTransaction(nutils.ParamContractAddress, global_params.INIT_NAME, bf.Bytes())
}

func newGoverConfigInit(config []byte) *types.Transaction {
	return utils.BuildNativeTransaction(nutils.GovernanceContractAddress, governance.INIT_CONFIG, config)
}
