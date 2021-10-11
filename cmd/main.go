package cmd

import (
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/jlog"
	admin "github.com/memocash/server/admin/server"
	"github.com/memocash/server/api"
	"github.com/memocash/server/cmd/peer"
	"github.com/memocash/server/cmd/test"
	db "github.com/memocash/server/db/server"
	"github.com/memocash/server/node"
	"github.com/memocash/server/ref/config"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run Server",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if err := config.Init(cmd); err != nil {
			jerr.Get("fatal error initializing config", err).Fatal()
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		var errorHandler = make(chan error)
		nodeGroup := node.NewGroup()
		apiServer := api.NewServer()
		adminServer := admin.NewServer(nodeGroup)
		queueServer0 := db.NewServer(config.DefaultShard0Port, 0)
		queueServer1 := db.NewServer(config.DefaultShard1Port, 1)
		go func() {
			err := apiServer.Run()
			errorHandler <- jerr.Get("error running api server", err)
		}()
		go func() {
			err := adminServer.Run()
			errorHandler <- jerr.Get("error running admin server", err)
		}()
		go func() {
			err := queueServer0.Run()
			errorHandler <- jerr.Get("error running db queue server shard 0", err)
		}()
		go func() {
			err := queueServer1.Run()
			errorHandler <- jerr.Get("error running db queue server shard 1", err)
		}()
		jlog.Logf("API server started on port: %d...\n", apiServer.Port)
		jlog.Logf("Admin server started on port: %d...\n", adminServer.Port)
		jlog.Logf("Queue server 0 started on port: %d...\n", queueServer0.Port)
		jlog.Logf("Queue server 1 started on port: %d...\n", queueServer1.Port)
		jerr.Get("fatal memo server error encountered", <-errorHandler).Fatal()
	},
}

func Execute() error {
	serverCmd.AddCommand(
		test.GetCommand(),
		peer.GetCommand(),
	)
	if err := serverCmd.Execute(); err != nil {
		return jerr.Get("error executing server command", err)
	}
	return nil
}
