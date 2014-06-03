package SocketIO

import (
  "testing"
  "fmt"
  "gopkg.in/redis.v1"
)

func TestPublish(t *testing.T) {
  emitter := NewEmitter(&EmitterOpts{
    Host:"localhost", 
    Port:6379,
  })

  if emitter == nil {
    t.Error("emitter is nil")
  }
  client := redis.NewTCPClient(&redis.Options{ Addr: "localhost:6379"})

  defer client.Close()
  pubsub := client.PubSub()
  defer pubsub.Close()
  err := pubsub.Subscribe("socket.io#emitter")
  _ = err
  msg, err := pubsub.Receive()
  fmt.Println(msg, err)

  emitter.Emit("broadcast event", "hogefuga")
  msg, err = pubsub.Receive()
  fmt.Println(msg, err)
}
