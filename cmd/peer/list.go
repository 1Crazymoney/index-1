package peer

import (
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/jlog"
	"github.com/memocash/server/admin/client"
	"github.com/memocash/server/db/item"
	"github.com/spf13/cobra"
	"net"
)

var listCmd = &cobra.Command{
	Use: "list",
	Run: func(cmd *cobra.Command, args []string) {
		peers, err := item.GetPeers(0, nil)
		if err != nil {
			jerr.Get("fatal error getting peers", err).Fatal()
		}
		jlog.Logf("Peers: %d\n", len(peers))
		for i := 0; i < len(peers) && i < 10; i++ {
			jlog.Logf("Peer: %s:%d - %d\n", net.IP(peers[i].Ip), peers[i].Port, peers[i].Services)
		}
	},
}

var listFoundPeersCmd = &cobra.Command{
	Use: "list-found-peers",
	Run: func(cmd *cobra.Command, args []string) {
		foundPeers, err := item.GetFoundPeers(0, nil)
		if err != nil {
			jerr.Get("fatal error getting found peers", err).Fatal()
		}
		jlog.Logf("Found peers: %d\n", len(foundPeers))
		for i := 0; i < len(foundPeers) && i < 10; i++ {
			jlog.Logf("Found peer: %s:%d - %s:%d\n", net.IP(foundPeers[i].Ip), foundPeers[i].Port,
				net.IP(foundPeers[i].FoundIp), foundPeers[i].FoundPort)
		}
	},
}

var listPeerFoundsCmd = &cobra.Command{
	Use: "list-peer-founds",
	Run: func(cmd *cobra.Command, args []string) {
		foundPeers, err := item.GetPeerFounds(0, nil)
		if err != nil {
			jerr.Get("fatal error getting peer founds", err).Fatal()
		}
		jlog.Logf("Peer founds: %d\n", len(foundPeers))
		for i := 0; i < len(foundPeers) && i < 10; i++ {
			jlog.Logf("Peer founds: %s:%d - %s:%d\n", net.IP(foundPeers[i].Ip), foundPeers[i].Port,
				net.IP(foundPeers[i].FinderIp), foundPeers[i].FinderPort)
		}
	},
}

var getCmd = &cobra.Command{
	Use: "get",
	Run: func(cmd *cobra.Command, args []string) {
		peerGet := client.NewPeerGet()
		if err := peerGet.Get(); err != nil {
			jerr.Get("fatal error getting admin index", err).Fatal()
		}
		jlog.Logf("peerGet.Message: %s\n", peerGet.Message)
	},
}
