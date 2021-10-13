package item

import (
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/jutil"
	"github.com/memocash/server/db/client"
	"github.com/memocash/server/ref/config"
	"net"
	"time"
)

type PeerConnectionStatus int

func (s PeerConnectionStatus) String() string {
	switch s {
	case PeerConnectionStatusFail:
		return "fail"
	case PeerConnectionStatusSuccess:
		return "success"
	default:
		return "unknown"
	}
}

const (
	PeerConnectionStatusFail    PeerConnectionStatus = 0
	PeerConnectionStatusSuccess PeerConnectionStatus = 1
)

type PeerConnection struct {
	Ip     []byte
	Port   uint16
	Time   time.Time
	Status PeerConnectionStatus
}

func (p PeerConnection) GetUid() []byte {
	return jutil.CombineBytes(
		jutil.BytePadPrefix(p.Ip, IpBytePadSize),
		jutil.GetUintData(uint(p.Port)),
		jutil.GetTimeByte(p.Time),
	)
}

func (p PeerConnection) GetShard() uint {
	return client.GetByteShard(p.Ip)
}

func (p PeerConnection) GetTopic() string {
	return TopicPeerConnection
}

func (p PeerConnection) Serialize() []byte {
	return jutil.GetIntData(int(p.Status))
}

func (p *PeerConnection) SetUid(uid []byte) {
	if len(uid) != IpBytePadSize+12 {
		return
	}
	p.Ip = jutil.ByteUnPad(uid[:IpBytePadSize])
	p.Port = uint16(jutil.GetUint(uid[IpBytePadSize : IpBytePadSize+4]))
	p.Time = jutil.GetByteTime(uid[IpBytePadSize+4:])
}

func (p *PeerConnection) Deserialize(data []byte) {
	p.Status = PeerConnectionStatus(jutil.GetInt(data))
}

type PeerConnectionsRequest struct {
	Shard   uint32
	StartId []byte
	Ip      []byte
	Port    uint32
}

func (r PeerConnectionsRequest) GetShard() uint32 {
	if len(r.Ip) > 0 {
		return client.GetByteShard32(r.Ip)
	}
	return r.Shard
}

func GetPeerConnections(request PeerConnectionsRequest) ([]*PeerConnection, error) {
	shardConfig := config.GetShardConfig(request.GetShard(), config.GetQueueShards())
	dbClient := client.NewClient(shardConfig.GetHost())
	var startIdBytes []byte
	if len(request.StartId) > 0 {
		startIdBytes = request.StartId
	}
	var prefixes [][]byte
	if len(request.Ip) > 0 {
		prefixes = [][]byte{jutil.CombineBytes(
			jutil.BytePadPrefix(request.Ip, IpBytePadSize),
			jutil.GetUintData(uint(request.Port)),
		)}
	}
	err := dbClient.GetWOpts(client.Opts{
		Topic:    TopicPeerConnection,
		Start:    startIdBytes,
		Prefixes: prefixes,
	})
	if err != nil {
		return nil, jerr.Get("error getting peer connections from queue client", err)
	}
	var peerConnections = make([]*PeerConnection, len(dbClient.Messages))
	for i := range dbClient.Messages {
		peerConnections[i] = new(PeerConnection)
		peerConnections[i].SetUid(dbClient.Messages[i].Uid)
		peerConnections[i].Deserialize(dbClient.Messages[i].Message)
	}
	return peerConnections, nil
}

func GetPeerConnectionLast(ip []byte, port uint16) (*PeerConnection, error) {
	shardConfig := config.GetShardConfig(client.GetByteShard32(ip), config.GetQueueShards())
	dbClient := client.NewClient(shardConfig.GetHost())
	err := dbClient.GetWOpts(client.Opts{
		Topic: TopicPeerConnection,
		Max:   1,
		Prefixes: [][]byte{jutil.CombineBytes(
			jutil.BytePadPrefix(ip, IpBytePadSize),
			jutil.GetUintData(uint(port)),
		)},
	})
	if err != nil {
		return nil, jerr.Getf(err, "error getting peer connection last for: %s:%d", net.IP(ip), port)
	}
	if len(dbClient.Messages) == 0 {
		return nil, jerr.Get("error no peer connection last found", client.EntryNotFoundError)
	}
	var peerConnection = new(PeerConnection)
	peerConnection.SetUid(dbClient.Messages[0].Uid)
	peerConnection.Deserialize(dbClient.Messages[0].Message)
	return peerConnection, nil
}