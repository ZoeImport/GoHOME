package cluster

import (
	"go-redis/interface/resp"
	"go-redis/resp/reply"
)

// RENAME k1 k2
func rename(cluster *ClusterDatabase, connection resp.Connection, cmdArgs [][]byte) resp.Reply {
	if len(cmdArgs) != 3 {
		return reply.MakeStandardErrorReply("ERR Wrong number of arguments.")
	}
	src := string(cmdArgs[1])
	dest := string(cmdArgs[2])
	srcPeer := cluster.peerPicker.PickNode(src)
	destPeer := cluster.peerPicker.PickNode(dest)
	if srcPeer != destPeer {
		return reply.MakeStandardErrorReply("ERR rename must within on peer")
	}
	return cluster.relay(srcPeer, connection, cmdArgs)
}
