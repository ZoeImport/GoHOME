package handler

import (
	"context"
	"errors"
	"go-redis/cluster"
	"go-redis/config"
	"go-redis/database"
	dataBaseface "go-redis/interface/database"
	"go-redis/lib/logger"
	"go-redis/lib/sync/atomic"
	"go-redis/resp/connection"
	"go-redis/resp/parser"
	"go-redis/resp/reply"
	"io"
	"net"
	"strings"
	"sync"
)

var (
	unKnownErrReplyBytes = []byte("-ERR unknown\r\n")
)

type RespHandler struct {
	activeCount sync.Map
	db          dataBaseface.Database
	closing     atomic.Boolean
}

func MakeHandler() *RespHandler {
	var db1 dataBaseface.Database
	if config.Properties.Self != "" && len(config.Properties.Peers) > 0 {
		db1 = cluster.MakeNewClusterDatabase()
	} else {
		db1 = database.NewStandAloneDatabase()
	}
	return &RespHandler{
		db: db1,
	}
}

func (r *RespHandler) closeClient(client *connection.Connection) {
	_ = client.Close()
	r.db.AfterClientClose(client)
	r.activeCount.Delete(client)
}

func (r *RespHandler) Handle(ctx context.Context, conn net.Conn) {
	if r.closing.Get() {
		_ = conn.Close()
	}
	client := connection.NewConnection(conn)
	r.activeCount.Store(client, struct{}{})
	ch := parser.ParseStream(conn)
	for payload := range ch {
		if payload.Err != nil {
			if payload.Err == io.EOF || errors.Is(io.ErrUnexpectedEOF, payload.Err) ||
				strings.Contains(payload.Err.Error(), "use of closed network connection") {
				r.closeClient(client)
				logger.Info("connection closed" + client.RemoteAddr().String())
				return
			}
			//protocol error
			errorReply := reply.MakeStandardErrorReply(payload.Err.Error())
			err := client.Write(errorReply.ToBytes())
			if err != nil {
				r.closeClient(client)
				logger.Info("connection closed" + client.RemoteAddr().String())
				return
			}
			continue
		}
		if payload.Data == nil {
			continue
		}
		bulkReply, ok := payload.Data.(*reply.MultiBulkReply)
		if !ok {
			logger.Error("require multi bulk reply")
			continue
		}
		result := r.db.Exec(client, bulkReply.Args)
		if result != nil {
			_ = client.Write(result.ToBytes())
		} else {
			_ = client.Write(unKnownErrReplyBytes)
		}
	}
}

func (r *RespHandler) Close() error {
	logger.Info("handler shutting down ")
	r.closing.Set(true)
	r.activeCount.Range(func(key, value any) bool {
		client := key.(*connection.Connection)
		_ = client.Close()
		return true
	})
	r.db.Close()
	return nil
}
