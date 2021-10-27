package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/jchavannes/btcd/chaincfg/chainhash"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/jutil"
	"github.com/memocash/server/admin/graph/generated"
	"github.com/memocash/server/admin/graph/model"
	"github.com/memocash/server/db/item"
)

func (r *mutationResolver) Null(ctx context.Context) (*int, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Tx(ctx context.Context, hash string) (*model.Tx, error) {
	chainHash, err := chainhash.NewHashFromStr(hash)
	if err != nil {
		return nil, jerr.Get("error getting tx hash from hash", err)
	}
	txHash := chainHash.CloneBytes()
	preloads := GetPreloads(ctx)
	var raw string
	if jutil.StringInSlice("raw", preloads) {
		txBlocks, err := item.GetSingleTxBlocks(txHash)
		if err != nil {
			return nil, jerr.Get("error getting tx blocks from items", err)
		}
		if len(txBlocks) == 0 {
			mempoolTxRaw, err := item.GetMempoolTxRawByHash(txHash)
			if err != nil {
				return nil, jerr.Get("error getting mempool tx raw", err)
			}
			raw = hex.EncodeToString(mempoolTxRaw.Raw)
		} else {
			txRaw, err := item.GetRawBlockTxByHash(txBlocks[0].BlockHash, txHash)
			if err != nil {
				return nil, jerr.Get("error getting block tx by hash", err)
			}
			raw = hex.EncodeToString(txRaw.Raw)
		}
	}
	return &model.Tx{
		Hash: chainHash.String(),
		Raw:  raw,
	}, nil
}

func (r *queryResolver) Txs(ctx context.Context) ([]*model.Tx, error) {
	return []*model.Tx{{
		Hash: "123",
		Raw:  "456",
	}}, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
