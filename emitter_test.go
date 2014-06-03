package SocketIO

import (
  "testing"
  "strings"
  _ "fmt"
  "github.com/garyburd/redigo/redis"
)

func TestPublish(t *testing.T) {
  emitter, _ := NewEmitter(&EmitterOpts{
    Host:"localhost", 
    Port:6379,
  })

  if emitter == nil {
    t.Error("emitter is nil")
  }
  c, _ := redis.Dial("tcp", "localhost:6379")
  defer c.Close()
  psc := redis.PubSubConn{Conn: c}
  psc.Subscribe("socket.io#emitter")
  emitter.Emit("time", "hogefuga")
for {
    switch v := psc.Receive().(type) {
    case redis.Message:
        isContain := strings.Contains(string(v.Data), "hogefuga")
        if !isContain {
          t.Errorf("%s not contains hogefuga", v.Data)
          return
        } else {
          return
        }
    }
}
}
