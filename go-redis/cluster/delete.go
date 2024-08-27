package cluster

import (
	"go-redis/interface/resp"
	"go-redis/resp/reply"
)

// del ...
func del(clusterDatabase *ClusterDatabase, connection resp.Connection, cmdArgs [][]byte) resp.Reply {
	replies := clusterDatabase.broadcast(connection, cmdArgs)
	var errReply reply.ErrorReply
	var deleted int64 = 0
	for _, rep := range replies {
		if reply.IsErrorReply(rep) {
			errReply = rep.(reply.ErrorReply)
			break
		}
		intReply, ok := rep.(*reply.IntReply)
		if !ok {
			errReply = reply.MakeStandardErrorReply("error")
		}
		deleted += intReply.Code
	}
	if errReply == nil {
		return reply.MakeIntReply(deleted)
	}
	return reply.MakeStandardErrorReply("error : " + errReply.Error())
}
