package test_tx

import (
	"encoding/hex"
	"github.com/jchavannes/btcd/chaincfg/chainhash"
	"github.com/jchavannes/btcd/wire"
	"github.com/memocash/index/ref/bitcoin/memo"
	"github.com/memocash/index/ref/bitcoin/tx/build"
	"github.com/memocash/index/ref/bitcoin/tx/gen"
	"github.com/memocash/index/ref/bitcoin/wallet"
)

func GetHexBytes(hash string) []byte {
	h, _ := hex.DecodeString(hash)
	return h
}

func GetHashBytes(hash string) []byte {
	h, _ := chainhash.NewHashFromStr(hash)
	return h.CloneBytes()
}

func GetAddressPkHash(address string) []byte {
	return wallet.GetAddressFromString(address).GetPkHash()
}

func GetPrivateKey(wif string) wallet.PrivateKey {
	key, _ := wallet.ImportPrivateKey(wif)
	return key
}

func GetAddress(address string) wallet.Address {
	return wallet.GetAddressFromString(address)
}

func GetBlockHeader(raw string) wire.BlockHeader {
	r, _ := hex.DecodeString(raw)
	header, _ := memo.GetBlockHeaderFromRaw(r)
	return *header
}

func GetKeyWallet(key *wallet.PrivateKey, utxos []memo.UTXO) build.Wallet {
	return build.Wallet{
		Getter: gen.GetWrapper(&TestGetter{
			UTXOs: utxos,
		}, key.GetPkHash()),
		KeyRing: wallet.KeyRing{
			Keys: []wallet.PrivateKey{*key},
		},
		Address: key.GetAddress(),
	}
}

func CopyTestWallet(wallet build.Wallet) build.Wallet {
	if len(wallet.KeyRing.Keys) == 0 {
		return build.Wallet{}
	}
	var key = wallet.KeyRing.Keys[0]
	var utxos []memo.UTXO
	if igw, ok := wallet.Getter.(*gen.InputGetterWrapper); ok {
		utxos = igw.UTXOs
		if tg, ok := igw.Old.(*TestGetter); ok {
			utxos = append(utxos, tg.UTXOs...)
		}
	}
	return build.Wallet{
		Getter: gen.GetWrapper(&TestGetter{
			UTXOs: utxos,
		}, key.GetPkHash()),
		KeyRing: wallet.KeyRing,
		Address: key.GetAddress(),
	}
}
