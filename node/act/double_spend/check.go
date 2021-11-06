package double_spend

import (
	"bytes"
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/server/db/item"
	"github.com/memocash/server/ref/bitcoin/memo"
	"github.com/memocash/server/ref/bitcoin/tx/hs"
	"github.com/memocash/server/ref/bitcoin/tx/script"
	"time"
)

type DoubleSpendCheck struct {
	ParentTxHash  []byte
	ParentTxIndex uint32
	LockHash      []byte
	Spends        []*DoubleSpendCheckSpend
}

func (c DoubleSpendCheck) IsWinnerSpend(spendCheck *DoubleSpendCheckSpend) (bool, error) {
	for _, spend := range c.Spends {
		if bytes.Equal(spend.TxHash, spendCheck.TxHash) {
			continue
		}
		if len(spend.BlockHash) > 0 && len(spendCheck.BlockHash) == 0 {
			return false, nil
		}
		if len(spend.BlockHash) == 0 && len(spendCheck.BlockHash) > 0 {
			return true, nil
		}
		return spendCheck.FirstSeen.Before(spend.FirstSeen), nil
	}
	return false, jerr.Newf("error no spend found to compare against")
}

type DoubleSpendCheckSpend struct {
	TxHash    []byte
	Index     uint32
	FirstSeen time.Time
	BlockHash []byte
}

func AttachAllToDoubleSpendChecks(doubleSpendChecks []*DoubleSpendCheck) error {
	if err := AttachSeensToSpendCheckSpends(doubleSpendChecks); err != nil {
		return jerr.Get("error attaching seens to spend check spends", err)
	}
	if err := AttachBlocksToSpendCheckSpends(doubleSpendChecks); err != nil {
		return jerr.Get("error attaching blocks to spend check spends", err)
	}
	if err := AttachLockHashesToSpendChecks(doubleSpendChecks); err != nil {
		return jerr.Get("error attaching lock hashes to spend check spends", err)
	}
	return nil
}

func AttachSeensToSpendCheckSpends(doubleSpendChecks []*DoubleSpendCheck) error {
	var txHashes [][]byte
	for _, doubleSpendCheck := range doubleSpendChecks {
		for _, spend := range doubleSpendCheck.Spends {
			txHashes = append(txHashes, spend.TxHash)
		}
	}
	txSeens, err := item.GetTxSeens(txHashes)
	if err != nil {
		return jerr.Get("error getting tx seens for double spend check spends", err)
	}
	for _, doubleSpendCheck := range doubleSpendChecks {
		for _, spend := range doubleSpendCheck.Spends {
			for _, txSeen := range txSeens {
				if bytes.Equal(txSeen.TxHash, spend.TxHash) {
					spend.FirstSeen = txSeen.Timestamp
					break
				}
			}
		}
	}
	return nil
}

// AttachBlocksToSpendCheckSpends
// TODO: Handle block hash already set, also include confirmation count
func AttachBlocksToSpendCheckSpends(doubleSpendChecks []*DoubleSpendCheck) error {
	var txHashes [][]byte
	for _, doubleSpendCheck := range doubleSpendChecks {
		for _, spend := range doubleSpendCheck.Spends {
			if len(spend.BlockHash) == 0 {
				txHashes = append(txHashes, spend.TxHash)
			}
		}
	}
	txBlocks, err := item.GetTxBlocks(txHashes)
	if err != nil {
		return jerr.Get("error getting tx blocks for double spend check spends", err)
	}
	for _, doubleSpendCheck := range doubleSpendChecks {
		for _, spend := range doubleSpendCheck.Spends {
			for _, txBlock := range txBlocks {
				if bytes.Equal(txBlock.TxHash, spend.TxHash) {
					spend.BlockHash = txBlock.BlockHash
					break
				}
			}
		}
	}
	return nil
}

// AttachLockHashesToSpendChecks assumes blocks attached before
func AttachLockHashesToSpendChecks(doubleSpendChecks []*DoubleSpendCheck) error {
	var txHashes [][]byte
	for _, doubleSpendCheck := range doubleSpendChecks {
		txHashes = append(txHashes, doubleSpendCheck.ParentTxHash)
	}
	txBlocks, err := item.GetTxBlocks(txHashes)
	if err != nil {
		return jerr.Get("error getting tx blocks for double spend lock hashes", err)
	}
	var mempoolTxHashes [][]byte
Loop:
	for _, txHash := range txHashes {
		for _, txBlock := range txBlocks {
			if bytes.Equal(txBlock.TxHash, txHash) {
				continue Loop
			}
		}
		mempoolTxHashes = append(mempoolTxHashes, txHash)
	}
	txBlockRaws, err := item.GetRawTxBlocksByHashes(txBlocks)
	if err != nil {
		return jerr.Get("error getting tx blocks for double spend check spends", err)
	}
	for _, doubleSpendCheck := range doubleSpendChecks {
		for _, txBlockRaw := range txBlockRaws {
			if bytes.Equal(txBlockRaw.TxHash, doubleSpendCheck.ParentTxHash) {
				msgTx, err := memo.GetMsgFromRaw(txBlockRaw.Raw)
				if err != nil {
					return jerr.Getf(err, "error parsing raw msg for tx: %s",
						hs.GetTxString(doubleSpendCheck.ParentTxHash))
				}
				doubleSpendCheck.LockHash = script.GetLockHash(msgTx.TxOut[doubleSpendCheck.ParentTxIndex].PkScript)
				break
			}
		}
	}
	return nil
}
