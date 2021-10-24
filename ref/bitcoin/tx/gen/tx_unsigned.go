package gen

import (
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/server/ref/bitcoin/memo"
)

func TxUnsigned(request TxRequest) (*memo.Tx, error) {
	create := Create{
		Request:     request,
		InputsToUse: request.InputsToUse,
		Outputs:     request.Outputs,
	}
	msgTx, err := create.Build()
	if err != nil {
		return nil, jerr.Get("error building tx", err)
	}
	memoTx := getMemoTx(msgTx, create.InputsToUse, create.Outputs)
	return memoTx, nil
}
