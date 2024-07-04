package database

import (
	"go-redis/datastruct/dict"
	"go-redis/interface/database"
	"go-redis/interface/resp"
	"go-redis/resp/reply"
	"strings"
)

type DB struct {
	index int
	data  dict.Dict
}

type ExecFunc func(db *DB, args [][]byte) resp.Reply

type CmdLine = [][]byte

func MakeDB() *DB {
	db := &DB{
		data: dict.MakeSyncDict(),
	}
	return db
}

func (db *DB) Exec(c resp.Connection, cmdline CmdLine) resp.Reply {
	cmdName := strings.ToLower(string(cmdline[0]))
	cmd, ok := cmdTable[cmdName]
	if !ok {
		return reply.MakeStandardErrorReply("ERR unknow command " + cmdName)
	}
	if !validateArities(cmd.arity, cmdline) {
		return reply.MakeArgNumberErrorReply(cmdName)
	}
	fun := cmd.exector
	return fun(db, cmdline[1:])
}

// SET K V -> arity = 3
// EXISTS K1 k2 arity = -2
func validateArities(arity int, args [][]byte) bool {
	argNumber := len(args)
	if argNumber >= 0 {
		return argNumber == arity
	}
	return argNumber >= -arity
}

func (db *DB) GetEntity(key string) (*database.DataEntity, bool) {
	raw, ok := db.data.Get(key)
	if !ok {
		return nil, false
	}
	entity, _ := raw.(*database.DataEntity)
	return entity, true
}

func (db *DB) PutEntity(key string, entity *database.DataEntity) int {
	return db.data.Put(key, entity)
}

func (db *DB) PubIfExists(key string, entity *database.DataEntity) int {
	return db.data.PutIfExists(key, entity)
}

func (db *DB) PutIfAbsent(key string, entity *database.DataEntity) int {
	return db.data.PutIfAbsent(key, entity)
}

func (db *DB) Remove(key string) int {
	return db.data.Remove(key)
}

func (db *DB) Removes(keys ...string) int {
	var deleted = 0
	for _, key := range keys {
		_, exists := db.data.Get(key)
		if exists {
			db.Remove(key)
			deleted++
		}
	}
	return deleted
}

func (db *DB) Flush() {
	db.data.Clear()
}
