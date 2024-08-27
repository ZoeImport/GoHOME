package cluster

import "go-redis/interface/resp"

func makeRouter() map[string]CmdFunc {
	routerMap := make(map[string]CmdFunc)
	routerMap["exists"] = defaultFunc
	routerMap["type"] = defaultFunc
	routerMap["set"] = defaultFunc
	routerMap["setnx"] = defaultFunc
	routerMap["get"] = defaultFunc
	routerMap["getset"] = defaultFunc
	routerMap["rename"] = rename
	routerMap["renamenx"] = rename
	routerMap["ping"] = ping
	routerMap["flushdb"] = flushdb
	routerMap["del"] = del
	routerMap["select"] = execSelect
	return routerMap
}

func defaultFunc(clusterDatabase *ClusterDatabase, connection resp.Connection, cmdArgs [][]byte) resp.Reply {
	key := string(cmdArgs[1])
	peer := clusterDatabase.peerPicker.PickNode(key)
	return clusterDatabase.relay(peer, connection, cmdArgs)
}
