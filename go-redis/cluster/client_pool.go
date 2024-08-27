package cluster

import (
	"context"
	"errors"
	pool "github.com/jolestar/go-commons-pool/v2"
	"go-redis/lib/logger"
	"go-redis/resp/client"
)

type connectionFactory struct {
	Peer string
}

func (conn *connectionFactory) MakeObject(ctx context.Context) (*pool.PooledObject, error) {
	c, err := client.MakeClient(conn.Peer)
	if err != nil {
		return nil, err
	}
	c.Start()
	return pool.NewPooledObject(c), nil
}

func (conn *connectionFactory) DestroyObject(ctx context.Context, object *pool.PooledObject) error {
	c, ok := object.Object.(*client.Client)
	if !ok {
		logger.Error("factory type mismatch")
		return errors.New("type mismatch")
	}
	c.Close()
	return nil
}

func (conn *connectionFactory) ValidateObject(ctx context.Context, object *pool.PooledObject) bool {
	return true
}

func (conn *connectionFactory) ActivateObject(ctx context.Context, object *pool.PooledObject) error {
	return nil
}

func (conn *connectionFactory) PassivateObject(ctx context.Context, object *pool.PooledObject) error {
	return nil
}
