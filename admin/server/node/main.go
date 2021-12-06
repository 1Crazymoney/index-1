package node

import "github.com/memocash/index/admin/admin"

func GetRoutes() []admin.Route {
	return []admin.Route{
		connectRoute,
		connectDefaultRoute,
		connectNextRoute,
		disconnectRoute,
		loopingEnableRoute,
		loopingDisableRoute,
		listConnectionsRoute,
		historyRoute,
		foundPeersRoute,
		getAddrsRoute,
		peersRoute,
		peerReportRoute,
	}
}
