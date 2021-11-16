package saver

import (
	"github.com/jchavannes/btcd/wire"
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/server/db/client"
	"github.com/memocash/server/db/item"
	"github.com/memocash/server/ref/config"
)

type ClearSuspect struct {
	Verbose bool
}

func (s *ClearSuspect) SaveTxs(block *wire.MsgBlock) error {
	if block == nil {
		return jerr.Newf("error nil block")
	}
	saveBlockHash := block.BlockHash()
	saveBlockHeight, err := item.GetBlockHeight(saveBlockHash.CloneBytes())
	if err != nil {
		return jerr.Get("error getting block height for clear suspect", err)
	}
	blocksToConfirm := config.GetBlocksToConfirm()
	if saveBlockHeight == nil || saveBlockHeight.Height < int64(blocksToConfirm) {
		return nil
	}
	confirmedHeightBlocks, err := item.GetHeightBlock(saveBlockHeight.Height - int64(blocksToConfirm))
	if err != nil {
		return jerr.Get("error getting height block for confirm to clear suspect", err)
	}
	if len(confirmedHeightBlocks) != 1 {
		return jerr.Newf("error unexpected number of height blocks returned for clear suspect: %d",
			len(confirmedHeightBlocks))
	}
	const limit = client.DefaultLimit
	var blockHash = confirmedHeightBlocks[0].BlockHash
	var startUid []byte
	for {
		blockTxes, err := item.GetBlockTxes(item.BlockTxesRequest{
			BlockHash: blockHash,
			StartUid:  startUid,
			Limit:     limit,
		})
		if err != nil {
			return jerr.Get("error getting block txs for clear suspect", err)
		}
		var txHashes = make([][]byte, len(blockTxes))
		for i := range blockTxes {
			txHashes[i] = blockTxes[i].TxHash
		}
		doubleSpendInputs, err := item.GetDoubleSpendInputsByTxHashes(txHashes)
		if err != nil {
			return jerr.Get("error getting double spend inputs by tx hashes", err)
		}
		var inputTxsToClear = make([][]byte, len(doubleSpendInputs))
		for i := range doubleSpendInputs {
			inputTxsToClear[i] = doubleSpendInputs[i].TxHash
		}
		if err := s.ClearSuspectAndDescendants(inputTxsToClear); err != nil {
			return jerr.Get("error clearing suspect and descendants", err)
		}
		if len(blockTxes) < limit {
			break
		}
		startUid = item.GetBlockTxUid(blockHash, blockTxes[len(blockTxes)-1].TxHash)
	}
	return nil
}

func (s *ClearSuspect) ClearSuspectAndDescendants(txHashes [][]byte) error {
	// TODO: Recursively go through descendants and clear suspect (could be lots)
	return nil
}

func NewClearSuspect(verbose bool) *ClearSuspect {
	return &ClearSuspect{
		Verbose: verbose,
	}
}
