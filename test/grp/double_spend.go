package grp

import (
	"github.com/jchavannes/btcd/wire"
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/server/node/obj/get"
	"github.com/memocash/server/node/obj/saver"
	"github.com/memocash/server/ref/bitcoin/memo"
	"github.com/memocash/server/ref/bitcoin/tx/build"
	"github.com/memocash/server/ref/bitcoin/tx/gen"
	"github.com/memocash/server/ref/bitcoin/tx/script"
	"github.com/memocash/server/ref/bitcoin/util/testing/test_tx"
	"github.com/memocash/server/ref/dbi"
)

const (
	FundingValue = 1e8
	SendAmount   = 1e5
)

type DoubleSpend struct {
	TxSaver         dbi.TxSave
	FundingTx       *memo.Tx
	FundingPkScript []byte
}

func (s *DoubleSpend) Init(wallet *build.Wallet) error {
	s.TxSaver = saver.CombinedTxSaver(false)
	var err error
	if s.FundingTx, err = test_tx.GetFundingTx(wallet.Address, FundingValue); err != nil {
		return jerr.Get("error getting funding tx for address", err)
	}
	if s.FundingPkScript, err = s.FundingTx.Outputs[0].Script.Get(); err != nil {
		return jerr.Get("error getting output script", err)
	}
	wallet.Getter.AddChangeUTXO(script.GetOutputUTXOs(s.FundingTx)[0])
	return nil
}

func (s *DoubleSpend) Create(output *memo.Output, wallet build.Wallet) (*memo.Tx, error) {
	var txRequest = gen.TxRequest{
		Outputs: []*memo.Output{output},
		Getter:  wallet.Getter,
		Change:  wallet.GetChange(),
		KeyRing: wallet.KeyRing,
	}
	tx, err := gen.Tx(txRequest)
	if err != nil {
		return nil, jerr.Get("error generating transaction", err)
	}
	if err := s.TxSaver.SaveTxs(memo.GetBlockFromTxs([]*wire.MsgTx{tx.MsgTx}, nil)); err != nil {
		return nil, jerr.Get("error saving tx", err)
	}
	return tx, nil
}

func (s *DoubleSpend) SaveBlock(tx *memo.Tx) error {
	txBlock := memo.GetBlockFromTxs([]*wire.MsgTx{s.FundingTx.MsgTx, tx.MsgTx}, &test_tx.Block1Header)
	if err := s.TxSaver.SaveTxs(txBlock); err != nil {
		return jerr.Get("error adding tx1 tx3 block1 to network", err)
	}
	return nil
}

func (s *DoubleSpend) GetAddressBalance(address string) (int64, error) {
	newBalance2, err := get.NewBalanceFromAddress(address)
	if err != nil {
		return 0, jerr.Get("error getting address 2 from string for balance", err)
	}
	if err := newBalance2.GetUtxos(); err != nil {
		return 0, jerr.Get("error getting address 2 balance from network", err)
	}
	return newBalance2.Balance, nil
}
