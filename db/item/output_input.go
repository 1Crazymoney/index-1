package item

import (
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/jutil"
	"github.com/memocash/index/db/client"
	"github.com/memocash/index/ref/bitcoin/memo"
	"github.com/memocash/index/ref/config"
)

type OutputInput struct {
	PrevHash  []byte
	PrevIndex uint32
	Hash      []byte
	Index     uint32
}

func (t OutputInput) GetUid() []byte {
	return jutil.CombineBytes(
		jutil.ByteReverse(t.PrevHash),
		jutil.GetUint32Data(t.PrevIndex),
		jutil.ByteReverse(t.Hash),
		jutil.GetUint32Data(t.Index),
	)
}

func (t *OutputInput) SetUid(uid []byte) {
	if len(uid) != 72 {
		return
	}
	t.PrevHash = jutil.ByteReverse(uid[:32])
	t.PrevIndex = jutil.GetUint32(uid[32:36])
	t.Hash = jutil.ByteReverse(uid[36:68])
	t.Index = jutil.GetUint32(uid[68:72])
}

func (t OutputInput) GetShard() uint {
	return client.GetByteShard(t.PrevHash)
}

func (t OutputInput) GetTopic() string {
	return TopicOutputInput
}

func (t OutputInput) Serialize() []byte {
	return nil
}

func (t *OutputInput) Deserialize([]byte) {}

func GetOutputInput(out memo.Out) ([]*OutputInput, error) {
	shard := GetShardByte32(out.TxHash)
	shardConfig := config.GetShardConfig(shard, config.GetQueueShards())
	db := client.NewClient(shardConfig.GetHost())
	prefix := jutil.CombineBytes(jutil.ByteReverse(out.TxHash), jutil.GetUint32Data(out.Index))
	if err := db.GetByPrefix(TopicOutputInput, prefix); err != nil {
		return nil, jerr.Get("error getting by prefix for output input", err)
	}
	var outputInputs = make([]*OutputInput, len(db.Messages))
	for i := range db.Messages {
		outputInputs[i] = new(OutputInput)
		outputInputs[i].SetUid(db.Messages[i].Uid)
		outputInputs[i].Deserialize(db.Messages[i].Message)
	}
	return outputInputs, nil
}

func GetOutputInputs(outs []memo.Out) ([]*OutputInput, error) {
	var shardOutGroups = make(map[uint32][]memo.Out)
	for _, out := range outs {
		shard := GetShardByte32(out.TxHash)
		shardOutGroups[shard] = append(shardOutGroups[shard], out)
	}
	var outputInputs []*OutputInput
	for shard, outGroup := range shardOutGroups {
		shardConfig := config.GetShardConfig(shard, config.GetQueueShards())
		db := client.NewClient(shardConfig.GetHost())
		var prefixes = make([][]byte, len(outGroup))
		for i := range outGroup {
			prefixes[i] = jutil.CombineBytes(
				jutil.ByteReverse(outGroup[i].TxHash),
				jutil.GetUint32Data(outGroup[i].Index),
			)
		}
		if err := db.GetByPrefixes(TopicOutputInput, prefixes); err != nil {
			return nil, jerr.Get("error getting by prefixes for output inputs", err)
		}
		for i := range db.Messages {
			var outputInput = new(OutputInput)
			outputInput.SetUid(db.Messages[i].Uid)
			outputInput.Deserialize(db.Messages[i].Message)
			outputInputs = append(outputInputs, outputInput)
		}
	}
	return outputInputs, nil
}

func GetOutputInputsForTxHashes(txHashes [][]byte) ([]*OutputInput, error) {
	var shardOutGroups = make(map[uint32][][]byte)
	for _, txHash := range txHashes {
		shard := GetShardByte32(txHash)
		shardOutGroups[shard] = append(shardOutGroups[shard], txHash)
	}
	var outputInputs []*OutputInput
	for shard, outGroup := range shardOutGroups {
		shardConfig := config.GetShardConfig(shard, config.GetQueueShards())
		db := client.NewClient(shardConfig.GetHost())
		var prefixes = make([][]byte, len(outGroup))
		for i := range outGroup {
			prefixes[i] = jutil.ByteReverse(outGroup[i])
		}
		if err := db.GetByPrefixes(TopicOutputInput, prefixes); err != nil {
			return nil, jerr.Get("error getting by prefixes for output inputs by tx hashes", err)
		}
		for i := range db.Messages {
			var outputInput = new(OutputInput)
			outputInput.SetUid(db.Messages[i].Uid)
			outputInput.Deserialize(db.Messages[i].Message)
			outputInputs = append(outputInputs, outputInput)
		}
	}
	return outputInputs, nil
}
