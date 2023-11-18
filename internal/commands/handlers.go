package commands

import (
	"sync"

	"github.com/vlaner/rediska/internal/resp"
)

var mu = sync.RWMutex{}
var storage = map[string]string{}

var Handlers = map[string]func([]resp.Value) resp.Value{
	"PING": ping,
	"SET":  set,
	"GET":  get,
}

func ping(args []resp.Value) resp.Value {
	if len(args) < 1 {
		return resp.Value{Typ: "string", Str: "PONG"}
	}
	return resp.Value{Typ: "string", Str: args[0].Bulk}
}

func set(args []resp.Value) resp.Value {
	if len(args) != 2 {
		return resp.Value{Typ: "error", Str: "invalid arguments"}
	}

	key, val := args[0], args[1]
	mu.Lock()
	defer mu.Unlock()

	storage[key.Bulk] = val.Bulk
	return resp.Value{Typ: "string", Str: "OK"}
}

func get(args []resp.Value) resp.Value {
	if len(args) != 1 {
		return resp.Value{Typ: "error", Str: "invalid arguments"}
	}

	key := args[0]

	mu.RLock()
	defer mu.RUnlock()

	val, ok := storage[key.Bulk]
	if !ok {
		return resp.Value{Typ: "null"}
	}

	return resp.Value{Typ: "bulk", Bulk: val}
}
