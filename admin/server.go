package admin

import (
	"encoding/json"
	"fmt"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/jlog"
	"github.com/memocash/server/db/item"
	"github.com/memocash/server/node"
	"github.com/memocash/server/ref/config"
	"io/ioutil"
	"net"
	"net/http"
)

type NodeDisconnectRequest struct {
	NodeId string
}

type NodeConnectRequest struct {
	Ip   []byte
	Port uint16
}

const (
	UrlIndex               = "/"
	UrlNodeGetAddrs        = "/node/get_addrs"
	UrlNodeConnect         = "/node/connect"
	UrlNodeConnectDefault  = "/node/connect_default"
	UrlNodeConnectNext     = "/node/connect_next"
	UrlNodeListConnections = "/node/list_connections"
	UrlNodeDisconnect      = "/node/disconnect"
	UrlNodeHistory         = "/node/history"
	UrlNodeLoopingEnable   = "/node/looping_enable"
	UrlNodeLoopingDisable  = "/node/looping_disable"
)

type Server struct {
	Nodes *node.Group
}

func (s *Server) Run() error {
	mux := http.NewServeMux()
	mux.HandleFunc(UrlIndex, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Memo Admin 0.1")
	})
	mux.HandleFunc(UrlNodeGetAddrs, func(w http.ResponseWriter, r *http.Request) {
		jlog.Log("Node get addrs request")
		for _, serverNode := range s.Nodes.Nodes {
			serverNode.GetAddr()
		}
	})
	mux.HandleFunc(UrlNodeConnect, func(w http.ResponseWriter, r *http.Request) {
		jlog.Log("Node connect")
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			jerr.Get("error reading node connect body", err).Print()
			return
		}
		var connectRequest = new(NodeConnectRequest)
		if err := json.Unmarshal(body, connectRequest); err != nil {
			jerr.Get("error unmarshalling node connect request", err).Print()
			return
		}
		s.Nodes.AddNode(connectRequest.Ip, connectRequest.Port)
	})
	mux.HandleFunc(UrlNodeConnectDefault, func(w http.ResponseWriter, r *http.Request) {
		jlog.Log("Node connect default")
		s.Nodes.AddDefaultNode()
	})
	mux.HandleFunc(UrlNodeConnectNext, func(w http.ResponseWriter, r *http.Request) {
		jlog.Log("Node connect next")
		s.Nodes.AddNextNode()
	})
	mux.HandleFunc(UrlNodeLoopingEnable, func(w http.ResponseWriter, r *http.Request) {
		jlog.Log("Node looping enable")
		if s.Nodes.Looping {
			return
		}
		s.Nodes.Looping = true
		if !s.Nodes.HasActive() {
			s.Nodes.AddNextNode()
		}
	})
	mux.HandleFunc(UrlNodeLoopingDisable, func(w http.ResponseWriter, r *http.Request) {
		jlog.Log("Node looping disabled")
		s.Nodes.Looping = false
	})
	mux.HandleFunc(UrlNodeListConnections, func(w http.ResponseWriter, r *http.Request) {
		jlog.Log("Node list connections")
		for id, serverNode := range s.Nodes.Nodes {
			fmt.Fprintf(w, "%s - %s:%d (%t)\n", id, net.IP(serverNode.Ip), serverNode.Port,
				serverNode.Peer != nil && serverNode.Peer.Connected())
		}
	})
	mux.HandleFunc(UrlNodeHistory, func(w http.ResponseWriter, r *http.Request) {
		jlog.Log("Node list history")
		peerConnections, err := item.GetPeerConnections(0, nil)
		if err != nil {
			jerr.Get("fatal error getting peer connections", err).Fatal()
		}
		fmt.Fprintf(w, "History peer connections: %d\n", len(peerConnections))
		for i := 0; i < len(peerConnections) && i < 10; i++ {
			fmt.Fprintf(w, "Peer connection: %s:%d - %s - %d\n", net.IP(peerConnections[i].Ip), peerConnections[i].Port,
				peerConnections[i].Time.Format("2006-01-02 15:04:05"), peerConnections[i].Status)
		}
	})
	mux.HandleFunc(UrlNodeDisconnect, func(w http.ResponseWriter, r *http.Request) {
		jlog.Log("Node disconnect")
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			jerr.Get("error reading node disconnect body", err).Print()
			return
		}
		var disconnectRequest = new(NodeDisconnectRequest)
		if err := json.Unmarshal(body, disconnectRequest); err != nil {
			jerr.Get("error unmarshalling node disconnect request", err).Print()
			return
		}
		for id, serverNode := range s.Nodes.Nodes {
			if id == disconnectRequest.NodeId {
				serverNode.Disconnect()
				fmt.Fprint(w, "Server disconnected")
				return
			}
		}
	})
	server := http.Server{
		Addr:    config.GetHost(config.GetAdminPort()),
		Handler: mux,
	}
	if err := server.ListenAndServe(); err != nil {
		return jerr.Get("error listening and serving admin server", err)
	}
	return nil
}

func NewServer(group *node.Group) *Server {
	return &Server{
		Nodes: group,
	}
}
