

package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/mixbee/mixbee/common/log"
	"github.com/mixbee/mixbee/common/serialization"
	"github.com/mixbee/mixbee/core/ledger"
	"github.com/mixbee/mixbee/core/types"
	"io"
	"os"
)

func ImportBlocks(importFile string, targetHeight uint32) error {
	currBlockHeight := ledger.DefLedger.GetCurrentBlockHeight()
	if targetHeight > 0 && currBlockHeight >= targetHeight {
		log.Infof("No blocks to import.")
		return nil
	}

	ifile, err := os.OpenFile(importFile, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer ifile.Close()
	fReader := bufio.NewReader(ifile)

	metadata := NewExportBlockMetadata()
	err = metadata.Deserialize(fReader)
	if err != nil {
		return fmt.Errorf("Metadata deserialize error:%s", err)
	}
	endBlockHeight := metadata.BlockHeight
	if endBlockHeight <= currBlockHeight {
		log.Infof("No blocks to import.\n")
		return nil
	}
	if targetHeight == 0 {
		targetHeight = endBlockHeight
	}
	if targetHeight < endBlockHeight {
		endBlockHeight = targetHeight
	}

	log.Infof("Start import blocks")
	log.Infof("Current block height:%d TotalBlocks:%d", currBlockHeight, endBlockHeight-currBlockHeight)

	for i := uint32(0); i <= endBlockHeight; i++ {
		size, err := serialization.ReadUint32(fReader)
		if err != nil {
			return fmt.Errorf("Read block height:%d error:%s", i, err)
		}
		compressData := make([]byte, size)

		_, err = io.ReadFull(fReader, compressData)
		if err != nil {
			return fmt.Errorf("Read block data height:%d error:%s", i, err)
		}
		if i <= currBlockHeight {
			continue
		}

		blockData, err := DecompressBlockData(compressData, metadata.CompressType)
		if err != nil {
			return fmt.Errorf("block height:%d decompress error:%s", i, err)
		}

		block := &types.Block{}
		err = block.Deserialize(bytes.NewReader(blockData))
		if err != nil {
			return fmt.Errorf("block height:%d deserialize error:%s", i, err)
		}

		err = ledger.DefLedger.AddBlock(block)
		if err != nil {
			return fmt.Errorf("add block height:%d error:%s", i, err)
		}
	}
	log.Infof("Import block complete, current block height:%d", ledger.DefLedger.GetCurrentBlockHeight())
	return nil
}
