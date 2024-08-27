package cluster

import "go-redis/interface/resp"

func ping(cluster *ClusterDatabase, connection resp.Connection, cmdArgs [][]byte) resp.Reply {
	return cluster.db.Exec(connection, cmdArgs)
}
