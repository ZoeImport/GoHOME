package cluster

import (
	"go-redis/interface/resp"
	"go-redis/lib/logger"
	"go-redis/resp/reply"
)

func flushdb(cluster *ClusterDatabase, connection resp.Connection, cmdArgs [][]byte) resp.Reply {
	replies := cluster.broadcast(connection, cmdArgs)
	var errReply reply.ErrorReply
	for _, rep := range replies {
		if reply.IsErrorReply(rep) {
			errReply = rep.(reply.ErrorReply)
			break
		}
	}
	if errReply == nil {
		return reply.MakeOkReply()
	}
	logger.Error("error" + errReply.Error())
	return reply.MakeStandardErrorReply("error " + errReply.Error())
}
