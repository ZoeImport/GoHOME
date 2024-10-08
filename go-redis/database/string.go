package database

import (
	"go-redis/interface/database"
	"go-redis/interface/resp"
	"go-redis/lib/utils"
	"go-redis/resp/reply"
)

//*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n
//*2\r\n$6\r\nselect\r\n$1\r\n3\r\n
//*2\r\n$3\r\nGET\r\n$3\r\nkey\r\n

//GET

func execGet(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	entity, existed := db.GetEntity(key)
	if !existed {
		return reply.MakeNullBulkReply()
	}
	bytes := entity.Data.([]byte)
	return reply.MakeBulkReply(bytes)
}

//SET

func execSet(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	value := args[1]
	entity := &database.DataEntity{
		Data: value,
	}
	db.PutEntity(key, entity)
	db.addAof(utils.ToCmdLine2("set", args...))
	return reply.MakeOkReply()

}

//SETNX

func execSetnx(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	value := args[1]
	entity := &database.DataEntity{
		Data: value,
	}
	result := db.PutIfAbsent(key, entity)
	db.addAof(utils.ToCmdLine2("setnx", args...))
	return reply.MakeIntReply(int64(result))

}

//GETSET

func execGetset(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	value := args[1]
	entity, existed := db.GetEntity(key)
	db.PutEntity(key, &database.DataEntity{
		Data: value,
	})
	if !existed {
		return reply.MakeNullBulkReply()
	}
	db.addAof(utils.ToCmdLine2("getset", args...))
	return reply.MakeBulkReply(entity.Data.([]byte))
}

//STRLEN

func execStrLen(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	entity, existed := db.GetEntity(key)
	if !existed {
		return reply.MakeNullBulkReply()
	}
	bytes := entity.Data.([]byte)
	return reply.MakeIntReply(int64(len(bytes)))
}

//init

func init() {
	RegisterCommand("Get", execGet, 2)
	RegisterCommand("Set", execSet, 3)
	RegisterCommand("SetNx", execSetnx, 3)
	RegisterCommand("GetSet", execGetset, 3)
	RegisterCommand("StrLen", execStrLen, 2)
}
