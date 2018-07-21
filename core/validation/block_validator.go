

package validation

import (
	"errors"
	"fmt"

	"github.com/mixbee/mixbee/core/ledger"
	"github.com/mixbee/mixbee/core/signature"
	"github.com/mixbee/mixbee/core/types"
	ontErrors "github.com/mixbee/mixbee/errors"
)

// VerifyBlock checks whether the block is valid
func VerifyBlock(block *types.Block, ld *ledger.Ledger, completely bool) error {
	header := block.Header
	if header.Height == 0 {
		return nil
	}

	m := len(header.Bookkeepers) - (len(header.Bookkeepers)-1)/3
	hash := block.Hash()
	err := signature.VerifyMultiSignature(hash[:], header.Bookkeepers, m, header.SigData)
	if err != nil {
		return err
	}

	prevHeader, err := ld.GetHeaderByHash(block.Header.PrevBlockHash)
	if err != nil {
		return fmt.Errorf("[BlockValidator], can not find prevHeader: %s", err)
	}

	err = VerifyHeader(block.Header, prevHeader)
	if err != nil {
		return err
	}

	//verfiy block's transactions
	if completely {
		/*
			//TODO: NextBookkeeper Check.
			bookkeeperaddress, err := ledger.GetBookkeeperAddress(ld.Blockchain.GetBookkeepersByTXs(block.Transactions))
			if err != nil {
				return errors.New(fmt.Sprintf("GetBookkeeperAddress Failed."))
			}
			if block.Header.NextBookkeeper != bookkeeperaddress {
				return errors.New(fmt.Sprintf("Bookkeeper is not validate."))
			}
		*/
		for _, txVerify := range block.Transactions {
			if errCode := VerifyTransaction(txVerify); errCode != ontErrors.ErrNoError {
				return errors.New(fmt.Sprintf("VerifyTransaction failed when verifiy block"))
			}

			if errCode := VerifyTransactionWithLedger(txVerify, ld); errCode != ontErrors.ErrNoError {
				return errors.New(fmt.Sprintf("VerifyTransaction failed when verifiy block"))
			}
		}
	}

	return nil
}

func VerifyHeader(header, prevHeader *types.Header) error {
	if header.Height == 0 {
		return nil
	}

	if prevHeader == nil {
		return errors.New("[BlockValidator], can not find previous block.")
	}

	if prevHeader.Height+1 != header.Height {
		return errors.New("[BlockValidator], block height is incorrect.")
	}

	if prevHeader.Timestamp >= header.Timestamp {
		return errors.New("[BlockValidator], block timestamp is incorrect.")
	}

	address, err := types.AddressFromBookkeepers(header.Bookkeepers)
	if err != nil {
		return err
	}

	if prevHeader.NextBookkeeper != address {
		return fmt.Errorf("bookkeeper address error")
	}

	return nil
}
