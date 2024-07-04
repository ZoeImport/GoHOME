package database

import (
	"go-redis/interface/resp"
	"go-redis/lib/wildcard"
	"go-redis/resp/reply"
)

//DEL

func execDel(db *DB, args [][]byte) resp.Reply {
	keys := make([]string, len(args))
	for i, v := range args {
		keys[i] = string(v)
	}
	deleted := db.Removes(keys...)
	return reply.MakeIntReply(int64(deleted))
}

//EXISTS

func execExists(db *DB, args [][]byte) resp.Reply {
	result := int64(0)
	for _, arg := range args {
		key := string(arg)
		_, existed := db.GetEntity(key)
		if existed {
			result++
		}
	}
	return reply.MakeIntReply(result)
}

//FLUSHDB

func execFlushDB(db *DB, args [][]byte) resp.Reply {
	db.Flush()
	return reply.MakeOkReply()
}

//TYPE

func execType(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	entity, existed := db.GetEntity(key)
	if !existed {
		return reply.MakeStatusReply("none")
	}
	switch entity.Data.(type) {
	case []byte:
		return reply.MakeStatusReply("string")
	}
	//TODO : other types
	return reply.MakeUnknownReply()
}

// RENAME

func execRename(db *DB, args [][]byte) resp.Reply {
	srcKey := string(args[0])
	destKey := string(args[1])
	entity, existed := db.GetEntity(srcKey)
	if !existed {
		return reply.MakeStandardErrorReply("no such key")
	}
	db.PutEntity(destKey, entity)
	db.Remove(srcKey)
	return reply.MakeOkReply()
}

// RENAMENX

func execRenamenx(db *DB, args [][]byte) resp.Reply {
	srcKey := string(args[0])
	destKey := string(args[1])

	_, ok := db.GetEntity(destKey)
	if ok {
		return reply.MakeIntReply(0)
	}
	entity, existed := db.GetEntity(srcKey)
	if !existed {
		return reply.MakeStandardErrorReply("no such key")
	}
	db.PutEntity(destKey, entity)
	db.Remove(srcKey)
	return reply.MakeIntReply(1)
}

//KEYS

func execKeys(db *DB, args [][]byte) resp.Reply {
	pattern := wildcard.CompilePattern(string(args[0]))
	result := make([][]byte, 0)
	db.data.ForEach(func(key string, val interface{}) bool {
		if pattern.IsMatch(key) {
			result = append(result, []byte(key))
		}
		return true
	})
	return reply.MakeMultiBulkReply(result)
}

func init() {
	RegisterCommand("DEL", execDel, -2)
	RegisterCommand("EXISTS", execExists, -2)
	RegisterCommand("flushdb", execFlushDB, -1)
	RegisterCommand("Type", execType, 2)
	RegisterCommand("RENAME", execRename, 3)
	RegisterCommand("RENAMENX", execRenamenx, 3)
	RegisterCommand("KEYS", execKeys, 2)
}
