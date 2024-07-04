package cluster

import "go-redis/interface/resp"

func execSelect(cluster *ClusterDatabase, connection resp.Connection, cmdArgs [][]byte) resp.Reply {
	return cluster.db.Exec(connection, cmdArgs)
}
